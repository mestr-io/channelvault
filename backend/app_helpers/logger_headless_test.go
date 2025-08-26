//go:build (dev&&headless)

package app_helpers

import (	
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
	
	//a race condition appears if getting context before 2 seconds...
	time.Sleep(2*time.Second)

	ctx := app.GetContext()
	if ctx == nil {
		t.Fatalf("got nil context")
	}
	
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
