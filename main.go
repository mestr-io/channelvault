package main

import (
	"archive/zip"
	"context" // For the context.Context type
	"database/sql"
	"embed" // For embedding frontend assets
	"fmt"   // For formatting output, like in LogPrintf
	"log"
	"path/filepath"
	"runtime/debug"

	"github.com/adrg/xdg"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"             // For checking the build mode
	"github.com/wailsapp/wails/v2/pkg/options/assetserver" // Ensure this is imported for assetserver.Options
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "SlackImportViewer"

type App struct {
	ctx    context.Context
	sqlite *sql.DB // Store the database connection pool
	duckdb *sql.DB
}

func handlePanic(ctx context.Context) {
	if r := recover(); r != nil {
		// Get the stack trace
		stack := string(debug.Stack())

		// Log the full panic and stack trace for debugging
		errorMessage := fmt.Sprintf("FATAL: Application panicked: %v\n%s", r, stack)

		// Determine the message to show to the user
		userMessage := "An unexpected error occurred. The application needs to close."

		// In development, show more details to the developer
		envInfo := runtime.Environment(ctx)
		if envInfo.BuildType != "production" {
			userMessage = errorMessage

		}
		// Show a dialog to the user
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Unhandled Application Error",
			Message: userMessage,
		})

		// Since this is a fatal error, we should quit.
		// runtime.LogFatal already does this, but we can be explicit.
		runtime.Quit(ctx)
	}
}

func NewApp() *App {
	return &App{}
}

func (a *App) getSqliteDatabasePath() (string, error) {
	// xdg.DataFile will create a path like: <OS_DATA_DIR>/<AppName>/<filename>
	// This is the most reliable and recommended method.
	dbPath, err := xdg.DataFile(filepath.Join(AppName, "data.db"))
	if err != nil {
		return "", fmt.Errorf("could not resolve XDG data path: %w", err)
	}
	return dbPath, nil
}

func (a *App) getDuckDBPath() (string, error) {
	// xdg.DataFile will create a path like: <OS_DATA_DIR>/<AppName>/<filename>
	// This is the most reliable and recommended method.
	dbPath, err := xdg.DataFile(filepath.Join(AppName, "data.duck"))
	if err != nil {
		return "", fmt.Errorf("could not resolve XDG data path: %w", err)
	}
	return dbPath, nil
}

func (a *App) GetViewState() string {
	return "setup"
}

func (a *App) SetViewState() {

}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	defer handlePanic(a.ctx)

	dbPath, err := a.getSqliteDatabasePath()
	if err != nil {
		log.Panicf("FATAL: Could not determine database path: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{ // Correct struct name
			Type:    runtime.ErrorDialog, // Correct constant name
			Title:   "Fatal Error",
			Message: fmt.Sprintf("Could not determine database path: %v", err),
		})
		log.Panicf("Could not determine database path: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Panicf("FATAL: fail to open database: %v", err)
	}

	a.sqlite = db
	runtime.LogInfof(ctx, "Database connection established at %s", dbPath)

	dbPath, err = a.getDuckDBPath()
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("duckdb", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	a.duckdb = db

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
	runtime.EventsEmit(a.ctx, "backend:viewState", "loading")
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

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "slack-viewer", // Use your application's title
		Width:  1024,
		Height: 768,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		OnStartup: app.startup,

		Bind: []any{
			app,
		},

		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},

		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.DEBUG,
	})

	if err != nil {
		println("Error starting Wails app:", err.Error())
	}
}
