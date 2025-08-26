//go:build dev

package app_helpers

import (
	"context"
	"testing"
	"time"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:testdata/dist
var assets embed.FS //needed by FakeApp

type FakeApp struct {
	ctx context.Context
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

func StartupFakeApp(app *FakeApp) error {
	return wails.Run(&options.App{
		Title:  "fake-app",
		Width:  1024,
		Height: 768,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnDomReady: app.Startup,
		Bind: []any{
			app,
		},
		Frameless:        false,
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.DEBUG,
		StartHidden:        true,
	})
}
func TearDownFakeApp(ctx context.Context, ch chan string) {
	//fmt.Printf("After %d second(s) ending test", afterSeconds)
	runtime.Quit(ctx)
	ch <- "test ends"
}

func WaitTestToEnd(t *testing.T, ch chan string) {
	select {
	case /*msg := */ <-ch:
		/*fmt.Println(msg)*/
	case <-time.After(2 * time.Second):
		t.Fatalf("Timeout")
	}
}
