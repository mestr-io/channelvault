//go:build dev && headless

package app_helpers

import (
	"context"
	"testing"
	"time"
)

func TestInitWailsLog_InitializesLogger(t *testing.T) {
	ch := make(chan string)
	app := &FakeApp{}

	//schedule test tear down

	go func() {
		time.Sleep(time.Duration(2500) * time.Millisecond)
		TearDownFakeApp(app.GetContext(), ch)
	}()

	err := StartupFakeApp(app)
	if err != nil {
		t.Fatalf("Error starting fake app: %s", err.Error())
	}
	
	// Wait for app context to be ready, polling with timeout instead of fixed sleep
	var ctx context.Context
	timeout := time.After(5 * time.Second)
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
	
	// Initialize the logger with the app context	
	InitWailsLog(ctx)
	defer func(){
		//undo InitWailsLog
		appLogger = nil
	}()

	if appLogger == nil {
		t.Fatal("appLogger should be initialized")
	}

	if appLogger.ctx != ctx {
		t.Fatal("appLogger context should match input context")
	}

	WaitTestToEnd(t, ch)
}
