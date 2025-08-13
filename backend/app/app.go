package app

import (
	"context" // For the context.Context type
	// For embedding frontend assets
	// For formatting output, like in LogPrintf
	"log/slog"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3" // For checking the build mode
	. "github.com/mestr-io/slack-export-viewer/backend/app_helpers"
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
	ctx        context.Context
	repository *repository.Repository
	state      ViewState
}

func New(repository *repository.Repository) *App {
	return &App{
		repository: repository,
	}
}

func (a *App) GetViewState() string {
	return string(ViewStateDashboard)
}

func (a *App) setViewState(state ViewState) {
	a.state = state
	runtime.EventsEmit(a.ctx, "backend:viewState", a.state)
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	// register callback
	runtime.OnFileDrop(ctx, a.onFileDrop)
	InitLog(ctx, slog.LevelInfo)
	Infof("Wails application started successfully!")
}

func (a *App) onFileDrop(x, y int, paths []string) {
	runtime.LogPrintf(a.ctx, "Go-side OnFileDrop received: %v at (%d,%d)", paths, x, y)
	if len(paths) != 1 {
		Fatalf("FATAL: internal error. import part incorrect: %v", paths)
	}

	a.processImport(paths[0])
}

func (a *App) processImport(path string) {
	a.setViewState(VIewStateLoading)

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
		Infof("processed: %v\n", p)
	}

	if importErr != nil {
		Errorf("Import failed: %v", importErr)
		runtime.EventsEmit(a.ctx, "importError", importErr.Error())
		a.setViewState(ViewStateSetup) // Go back to setup on failure
		return
	}

	Infof("Import completed successfully!")
	runtime.EventsEmit(a.ctx, "importComplete")
	a.setViewState(ViewStateDashboard) // Move to the dashboard
}
