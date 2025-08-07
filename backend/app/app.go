package app

import (
	"context" // For the context.Context type
	// For embedding frontend assets
	// For formatting output, like in LogPrintf
	"log"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3" // For checking the build mode
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
	return string(ViewStateSetup)
}

func (a *App) setViewState(state ViewState) {
	a.state = state
	runtime.EventsEmit(a.ctx, "backend:viewState", a.state)
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	// register callback
	runtime.OnFileDrop(ctx, a.onFileDrop)

	runtime.LogPrint(ctx, "Wails application started successfully!")
}

func (a *App) onFileDrop(x, y int, paths []string) {
	runtime.LogPrintf(a.ctx, "Go-side OnFileDrop received: %v at (%d,%d)", paths, x, y)
	if len(paths) != 1 {
		log.Panicf("FATAL: internal error. import part incorrect: %v", paths)
	}

	a.processImport(paths[0])
}

// processImport is the new method that manages the import lifecycle.
func (a *App) processImport(path string) {
	// 1. Set state to Loading and notify frontend
	a.setViewState(VIewStateLoading)

	// 2. Create the channel and the importer
	progressChan := make(chan importer.Progress)
	imp := importer.NewImporter(a.ctx, path, a.repository)

	// 3. Start the import in its own goroutine
	var importErr error
	go func() {
		importErr = imp.Import(progressChan)
	}()

	// 4. Listen for progress updates until the channel is closed
	var lastUpdateTime time.Time

	for p := range progressChan {
		// For every progress update, emit an event to the Wails frontend
		if time.Since(lastUpdateTime) >= time.Second {
			runtime.EventsEmit(a.ctx, "importProgress", p)
			lastUpdateTime = time.Now()
		}
		log.Printf("processed: %v\n", p)
	}

	// 5. The loop above finishes only when the channel is closed.
	//    Now we can check if the import function returned an error.
	if importErr != nil {
		log.Printf("ERROR: Import failed: %v", importErr)
		runtime.EventsEmit(a.ctx, "importError", importErr.Error())
		a.setViewState(ViewStateSetup) // Go back to setup on failure
		return
	}

	// 6. Success!
	runtime.LogPrint(a.ctx, "Import completed successfully!")
	runtime.EventsEmit(a.ctx, "importComplete")
	a.setViewState(ViewStateDashboard) // Move to the dashboard
}
