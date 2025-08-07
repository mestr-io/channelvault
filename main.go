package main

import (

	// For the context.Context type

	"embed" // For embedding frontend assets
	"fmt"
	"log"
	"path/filepath"

	// For formatting output, like in LogPrintf

	"github.com/adrg/xdg"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mestr-io/slack-export-viewer/backend/app"
	"github.com/mestr-io/slack-export-viewer/backend/repository"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"             // For checking the build mode
	"github.com/wailsapp/wails/v2/pkg/options/assetserver" // Ensure this is imported for assetserver.Options
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "SlackImportViewer"

func getSqliteDatabasePath(appName string) (string, error) {
	// xdg.DataFile will create a path like: <OS_DATA_DIR>/<AppName>/<filename>
	// This is the most reliable and recommended method.
	dbPath, err := xdg.DataFile(filepath.Join(appName, "data.db"))
	if err != nil {
		return "", fmt.Errorf("could not resolve XDG data path: %w", err)
	}
	return dbPath, nil
}

func main() {
	dbPath, err := getSqliteDatabasePath(AppName)
	if err != nil {
		log.Panicf("FATAL: Could not determine database path: %v", err)
	}

	r, err := repository.New(dbPath)
	if err != nil {
		log.Panicf("FATAL: Could not setup sqlite DB: %v", err)
	}
	app := app.New(r)

	err = wails.Run(&options.App{
		Title:  "slack-viewer", // Use your application's title
		Width:  1024,
		Height: 768,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		OnStartup: app.Startup,

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
