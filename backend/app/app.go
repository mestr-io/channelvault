package app

import (
	"context" // For the context.Context type
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	// For embedding frontend assets
	// For formatting output, like in LogPrintf
	// For formatting output, like in LogPrintf
	"time"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3" // For checking the build mode
	ah "github.com/mestr-io/slack-export-viewer/backend/app_helpers"
	"github.com/mestr-io/slack-export-viewer/backend/importer"
	"github.com/mestr-io/slack-export-viewer/backend/models"
	"github.com/mestr-io/slack-export-viewer/backend/repository"

	// Ensure this is imported for assetserver.Options
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ViewState string

const (
	ViewStateSetup     ViewState = "setup"
	ViewStateDashboard ViewState = "dashboard"
	ViewStateLoading   ViewState = "loading"
)

type ImportTaskStatus struct {
	ImportID int64
}

type DownloadTaskStatus struct {
	ImportID int64
}

type App struct {
	ctx          context.Context
	repository   repository.Repository
	state        ViewState
	filesFolder  string
	importTask   *ImportTaskStatus
	downloadTask *DownloadTaskStatus
}

func New(repository repository.Repository, filesFolder string) *App {
	return &App{
		repository:  repository,
		filesFolder: filesFolder,
	}
}

func (a *App) GetViewState() string {
	ready, err := a.repository.NoFinishedImports(a.ctx)
	if err != nil {
		ah.Fatalf("FATAL: cannot get current imports: %v", err)
	}
	if ready {
		return string(ViewStateDashboard)
	} else {
		return string(ViewStateSetup)
	}
}

func (a *App) setViewState(state ViewState) {
	a.state = state
	runtime.EventsEmit(a.ctx, "backend:viewState", a.state)
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	// register callback
	runtime.OnFileDrop(ctx, a.onFileDrop)
	ah.InitWailsLog(ctx)
	ah.Infof("Wails application started successfully!")

	// start here any pending task
	imports, err := a.repository.ListImports(a.ctx)
	if err != nil {
		panic(err)
	}
	for _, imp := range imports {
		// check for pending import
		if imp.Finished == "" {
			// import has not finished
			err = a.ResumeImport(imp.ID)
			if err != nil {
				panic(err)
			}
		}

		// check for pending downloads
		if imp.DownloadStarted.Valid && !imp.DownloadFinished.Valid {
			err = a.DownloadBinaries(imp.ID, false)
			if err != nil {
				panic(err)
			}
		}
	}

}

func (a *App) onFileDrop(x, y int, paths []string) {
	runtime.LogPrintf(a.ctx, "Go-side OnFileDrop received: %v at (%d,%d)", paths, x, y)
	if len(paths) != 1 {
		ah.Fatalf("internal error. import part incorrect: %v", paths)
	}

	record, err := a.repository.CreateImport(a.ctx, paths[0], defaultImportName(paths[0]))
	if err != nil {
		ah.Errorf("Import failed: %v", err)
		a.setViewState(ViewStateSetup) // Go back to setup on failure
		return
	}
	a.processImport(record)
}

func defaultImportName(path string) string {
	name, _ := strings.CutSuffix(filepath.Base(path), ".zip")
	return name
}

func (a *App) processImport(currentImport *models.SlackImport) {
	// notify the FE
	currentImport.Error = ""
	currentImport.NumFiles = 0
	runtime.EventsEmit(a.ctx, "startImport", currentImport)

	progressChan := make(chan models.Progress)
	imp := importer.NewImporter(a.ctx, a.repository, currentImport)

	var wg sync.WaitGroup
	wg.Add(1)

	var importErr error
	go func() {
		defer wg.Done()
		importErr = imp.Import(progressChan)
	}()

	var lastUpdateTime time.Time

	for p := range progressChan {
		if time.Since(lastUpdateTime) >= (time.Millisecond * 250) {
			runtime.EventsEmit(a.ctx, "importProgress", p)
			lastUpdateTime = time.Now()
		}
		currentImport.NumFiles++
	}

	// wait for completion or error
	wg.Wait()

	if importErr != nil {
		ah.Errorf("Import failed: %v", importErr)
		runtime.EventsEmit(a.ctx, "importError", importErr.Error())
		currentImport.Error = fmt.Sprintf("%v", importErr)
		// update the import
		err := a.repository.FinishedImportUpdateWithError(a.ctx, currentImport)
		if err != nil {
			ah.Errorf("Failed to update import status: %v", err)
		}
	} else {
		ah.Infof("Import completed successfully!!")
		runtime.EventsEmit(a.ctx, "importComplete")
		// update the import
		err := a.repository.FinishedImportUpdate(a.ctx, currentImport)
		if err != nil {
			ah.Errorf("Failed to update import status: %v", err)
		}
	}

}
