//go:build (dev&&headless)
package app

import (
	"testing"
	"context"
	"time"

	"github.com/stretchr/testify/assert"

	ah "github.com/mestr-io/slack-export-viewer/backend/app_helpers"
)

func TestSetViewState(t *testing.T) {
	ch := make(chan string)

	// Create a fake app instance to get a valid Wails context
	app := new(ah.FakeApp)

	//schedule test tear down

	go func() {
		time.Sleep(time.Duration(2500) * time.Millisecond)
		ah.TearDownFakeAppWithCtx(app.GetContext(), ch)
	}()

	_ = ah.StartupFakeApp(app, "")


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

	//ctx := app.GetContext()
	//if ctx == nil {
	//	t.Fatalf("got nil context")
	//}

	// Create a real App instance with the fake app context
	realApp := New(nil, "")
	realApp.ctx = app.GetContext() // In a real scenario, this would be set by Wails

	// This test is limited because we can't easily check for the emitted event
	// without a more complex setup. We'll just check the state change.
	realApp.setViewState(ViewStateDashboard)
	assert.Equal(t, ViewStateDashboard, realApp.state)

	ah.WaitTestToEnd(t, ch)
}
