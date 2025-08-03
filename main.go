package main

import (

	// For the context.Context type

	"embed" // For embedding frontend assets

	// For formatting output, like in LogPrintf

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mestr-io/slack-export-viewer/backend/app"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"             // For checking the build mode
	"github.com/wailsapp/wails/v2/pkg/options/assetserver" // Ensure this is imported for assetserver.Options
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "SlackImportViewer"

func main() {
	app := app.NewApp(AppName)

	err := wails.Run(&options.App{
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
