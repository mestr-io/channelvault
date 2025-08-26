//go:build dev

package app_helpers

import (
	"context"
	"testing"
	"time"
	"embed"
	
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/application"
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

func StartupFakeApp(app *FakeApp, loggerFile string) (mainApp *application.Application) {

	var thatLogger logger.Logger

	if loggerFile == "" {
		thatLogger = logger.NewDefaultLogger()
	} else {
		thatLogger = logger.NewFileLogger(loggerFile)
	}

	mainApp = application.NewWithOptions(&options.App{
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
		Logger:				thatLogger,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.DEBUG,
		StartHidden:        true,
	})

	go func() {
		err := mainApp.Run()
		if err != nil {
			panic(err)
		}
	}()

	return
}

func TearDownFakeApp(app *application.Application, ch chan string) {
	app.Quit()
	ch <- "test ends"
}

func TearDownFakeAppWithCtx(ctx context.Context, ch chan string) {
	runtime.Quit(ctx)
	ch <- "test ends"
}

func WaitTestToEnd(t *testing.T, ch chan string) {
	select {
	case /*msg := */ <-ch:
		/*fmt.Println(msg)*/
	case <-time.After(5 * time.Second):
		t.Fatalf("Timeout")
	}
}

func GetSomeWailsContext(t *testing.T, loggerFile string) (waitTeardown func(), 
										ctx context.Context,
										err error){
	ch := make(chan string)
	app := &FakeApp{}

	waitTeardown = func(){
		WaitTestToEnd(t, ch)
	}
		
	wailsApp := StartupFakeApp(app, loggerFile) // needs teardown scheduled before...	

	//schedule test tear down
	go func() {		
		time.Sleep(time.Duration(2500) * time.Millisecond)
		TearDownFakeApp(wailsApp, ch)
	}()	
	
	// Wait for app context to be ready, polling with timeout instead of fixed sleep
	timeout := time.After(2 * time.Second)
	tick := time.Tick(50 * time.Millisecond)

	ever := true
	for ever {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for app context to be ready")
		case <-tick:
			ctx = app.GetContext()
			if ctx != nil {
				ever = false				
			}
		}
	}

	return
}
