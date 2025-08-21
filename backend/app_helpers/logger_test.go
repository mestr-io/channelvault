//go:build dev
package app_helpers

import (	
	"testing"	
	"embed"
	"time"
	
	"github.com/stretchr/testify/assert"
)

//go:embed all:testdata/dist
var assets embed.FS//neded by FakeApp

func TestAppLoggerIsNilBeforeInitLog(t *testing.T) {
	assert := assert.New(t)
	var nilPtr *WailsAppLogger
	assert.EqualValues(appLogger, nilPtr)
}

func TestCoreLogfCallsRightFunction(t *testing.T) {
	assert := assert.New(t)
	var nilPtr *WailsAppLogger
	assert.EqualValues(appLogger, nilPtr)
	called := false

	wailsFun := func(string, ...any) {
		t.Fatal("Wails log shouldn't be called before InitWailsLog is called")
	
	}
	nonWailsFun := func(f string, v...any) {
		called = true
	}

    coreLogf(infoLabel, wailsFun, nonWailsFun, "test message")    
	assert.True(called)
}


func TestInitWailsLog_InitializesLogger(t *testing.T) {		
	ch := make(chan string)
	app := &FakeApp{}	

	//avoid wails to see go test flags
	//save := backupFlags()
	//resetFlags()
	//defer restoreFlags(save)
	
	//schedule test tear down		
	
	go func() {
		time.Sleep(time.Duration(1500) * time.Millisecond)
		TearDownFakeApp(app.GetContext(), ch)
	}()

	err := StartupFakeApp(app)
	if err != nil {
		t.Fatalf("Error starting fake app: %s", err.Error())
	}

	ctx := app.GetContext()
	if ctx == nil {
		t.Fatalf("got nil context")
	}
	InitWailsLog(ctx)	

    if appLogger == nil {
        t.Fatal("appLogger should be initialized")
    }

    if appLogger.ctx != ctx {
        t.Fatal("appLogger context should match input context")
    }

	WaitTestToEnd(t, ch)	
}

