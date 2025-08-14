package app_helpers

import (
	//"context"
	//"log/slog"
	"testing"
	//"embed"
	"github.com/stretchr/testify/assert"
	//"github.com/wailsapp/wails/v2"
	//"github.com/wailsapp/wails/v2/pkg/options"
	//"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func TestAppLoggerIsNilBeforeInitLog(t *testing.T) {
	assert := assert.New(t)
	var nilPtr *WailsAppLogger
	assert.EqualValues(appLogger, nilPtr)
}

/* WIP
func TestAppLoggerIsNotNilAfterInitLog(t *testing.T) {
	assert := assert.New(t)
	var nilPtr *WailsAppLogger

	app := CallWailsRun(t)
	ctx := app.GetContext()
	InitWailsLog(ctx, slog.LevelDebug)
	assert.NotEqualValues(appLogger, nilPtr)
}

type FakeApp struct {
	ctx        context.Context
}

func (a *FakeApp) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *FakeApp) GetContext() context.Context {
	return a.ctx
}

func (a *FakeApp) GetViewState() string {
	return string("setup")
}

//go:embed all:frontend/dist
var assets embed.FS

func CallWailsRun(t *testing.T) *FakeApp {
	app := FakeApp{}
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
		Frameless:        false,
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
	})

	if err != nil {
		t.Fatalf("Error starting Wails app: %s", err.Error())
	}
	return &app
}
*/
