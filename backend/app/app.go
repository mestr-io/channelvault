package app

import (
	"archive/zip"
	"context"      // For the context.Context type
	"database/sql" // For embedding frontend assets
	"fmt"          // For formatting output, like in LogPrintf
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3" // For checking the build mode

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
	ctx context.Context
	db  *sql.DB
}

func getSqliteDatabasePath(appName string) (string, error) {
	// xdg.DataFile will create a path like: <OS_DATA_DIR>/<AppName>/<filename>
	// This is the most reliable and recommended method.
	dbPath, err := xdg.DataFile(filepath.Join(appName, "data.db"))
	if err != nil {
		return "", fmt.Errorf("could not resolve XDG data path: %w", err)
	}
	return dbPath, nil
}

func NewApp(appName string) *App {
	dbPath, err := getSqliteDatabasePath(appName)
	if err != nil {
		log.Panicf("FATAL: Could not determine database path: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("FATAL: Could not determine database path: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Panicf("FATAL: fail to open database: %v", err)
	}

	return &App{
		db: db,
	}
}

func (a *App) GetViewState() string {
	return string(ViewStateDashboard)
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
	go func(p string) {
		a.Import(p)
	}(paths[0])
}

func (a *App) Import(path string) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Panicf("FATAL: fail to open the zip export file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		select {
		case <-a.ctx.Done():
			runtime.LogInfo(a.ctx, "Import operation was canceled.")
			return
		default:
			// Context is not canceled, so we can continue our work.
			runtime.LogInfof(a.ctx, "Processing file: %s", f.Name)
			runtime.EventsEmit(a.ctx, "import:progress", f.Name)

			// parse users.json -> users.parquet
			// parse channels.json -> channels.parquet
			// iterate over channels
			// for each channel in channels
			// ---- enter channel directory and parse messages and files -> messages/chanel_idxxx.parquet  files/channel_idxxx.parquet
			// -----download binaries
		}
	}
}
