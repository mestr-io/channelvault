//go:build dev

package app_helpers

import (
	"context"
	"testing"
	"time"
	"embed"
	"sync"
	"fmt"
	
	"github.com/wailsapp/wails/v2/pkg/runtime" //para runtime.Quit
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/application"
)

//go:embed all:testdata/dist
var assets embed.FS //needed by FakeApp

type FakeApp struct {
	mu sync.RWMutex
	ctx context.Context
	shutdownHook func(ctx context.Context) 
	beforeCloseHook func(ctx context.Context) bool
}

func (a *FakeApp) Startup(ctx context.Context) {
	a.mu.Lock()
    defer a.mu.Unlock()
	a.ctx = ctx
}

func (a *FakeApp) GetContext() context.Context {
	a.mu.RLock()
    defer a.mu.RUnlock()
	return a.ctx
}

func (a *FakeApp) Shutdown(ctx context.Context) {
	if a.shutdownHook != nil {
		a.shutdownHook(ctx)
	}	
}

func (a *FakeApp) BeforeClose(ctx context.Context) bool {
	if a.beforeCloseHook != nil {
		return a.beforeCloseHook(ctx)
	}
	return false
}

func (a *FakeApp) GetViewState() string {
	return string("setup")
}

func StartupFakeApp(app *FakeApp, loggerFile string, ch chan struct{}) (mainApp *application.Application) {

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
		OnShutdown: app.Shutdown,
		OnBeforeClose: app.BeforeClose,
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
		fmt.Println("Starting main app...", time.Now())
		err := mainApp.Run()
		fmt.Println("Main app stopped.", time.Now())
		if ch != nil {
			close(ch)
		}
		if err != nil {
			panic(err)
		}		
	}()

	return
}

func GetSomeWailsContext(t *testing.T, loggerFile string) (waitTeardown func(),
															ctx context.Context,
															err error,
															app *FakeApp) {
	ch := make(chan struct{})
	
	app = &FakeApp{ shutdownHook: func(ctx context.Context) {		
		fmt.Println("FakeApp Shutdown called.", time.Now())		
	}}	

	waitTeardown = func(){			
		fmt.Println("waiting Done channel")
		tick := time.Tick(1000 * time.Millisecond)
		for {
			select{
			case <-ch:
				return
			case <-tick:
				fmt.Println("tick tick")
			}
		}		
	}
	
	/*wailsApp :=*/ StartupFakeApp(app, loggerFile, ch) // needs teardown scheduled before...
	
	//schedule test tear down
	go func() {		
		//some time for the test to live		
		time.Sleep(time.Duration(5) * time.Second)
		//get context
		c := app.GetContext()
		//kill		
		runtime.Quit(c)
		//wait for shutdown to complete
		//<-c.Done()		
	}()

	// Wait for app context to be ready, polling with timeout instead of fixed sleep
	timeout := time.After(4 * time.Second)
	tick := time.Tick(50 * time.Millisecond)

	ever := true
	for ever {
		select {
		case <-timeout:
			t.Fatalf("GetSomeWailsContext timeout waiting for app context to be ready")
		case <-tick:
			ctx = app.GetContext()
			if ctx != nil {
				ever = false				
			}
		}
	}

	return
}
