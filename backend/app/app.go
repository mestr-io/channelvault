package app

import (
	"context" // For the context.Context type
	// For embedding frontend assets
	// For formatting output, like in LogPrintf
	"time"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3" // For checking the build mode
	ah "github.com/mestr-io/slack-export-viewer/backend/app_helpers"
	"github.com/mestr-io/slack-export-viewer/backend/importer"
	"github.com/mestr-io/slack-export-viewer/backend/repository"

	// Ensure this is imported for assetserver.Options
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ViewState string

const (
	ViewStateSetup     ViewState = "setup"
	ViewStateDashboard ViewState = "dashboard"
	VIewStateLoading   ViewState = "loading"
)

type App struct {
	ctx         context.Context
	repository  repository.Repository
	state       ViewState
	filesFolder string
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
}

func (a *App) onFileDrop(x, y int, paths []string) {
	runtime.LogPrintf(a.ctx, "Go-side OnFileDrop received: %v at (%d,%d)", paths, x, y)
	if len(paths) != 1 {
		ah.Fatalf("internal error. import part incorrect: %v", paths)
	}

	a.processImport(paths[0])
}

func (a *App) processImport(path string) {
	a.setViewState(VIewStateLoading)
	record, err := a.repository.StoreImport(a.ctx, path)
	if err != nil {
		ah.Errorf("Import failed: %v", err)
		a.setViewState(ViewStateSetup) // Go back to setup on failure
		return
	}

	progressChan := make(chan importer.Progress)
	imp := importer.NewImporter(a.ctx, path, a.repository)

	var importErr error
	go func() {
		importErr = imp.Import(progressChan)
	}()

	var lastUpdateTime time.Time

	for p := range progressChan {
		if time.Since(lastUpdateTime) >= (time.Millisecond * 250) {
			runtime.EventsEmit(a.ctx, "importProgress", p)
			lastUpdateTime = time.Now()
		}
		ah.Infof("processed: %v\n", p)
	}

	if importErr != nil {
		ah.Errorf("Import failed: %v", importErr)
		runtime.EventsEmit(a.ctx, "importError", importErr.Error())
		a.setViewState(ViewStateSetup) // Go back to setup on failure
		return
	}

	ah.Infof("Import completed successfully!!")
	runtime.EventsEmit(a.ctx, "importComplete")
	a.repository.FinishedImport(a.ctx, record.ID)

	a.setViewState(ViewStateDashboard) // Move to the dashboard
}
