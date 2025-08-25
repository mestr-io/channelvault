//go:build (dev&&headless)
package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ah "github.com/mestr-io/slack-export-viewer/backend/app_helpers"
)

func TestSetViewState(t *testing.T) {
	ch := make(chan string)
	app := &ah.FakeApp{}

	//schedule test tear down

	go func() {
		time.Sleep(time.Duration(2500) * time.Millisecond)
		ah.TearDownFakeApp(app.GetContext(), ch)
	}()

	err := ah.StartupFakeApp(app)
	if err != nil {
		t.Fatalf("Error starting fake app: %s", err.Error())
	}
	
	//a race condition appears if getting context before one second...
	time.Sleep(time.Second)

	ctx := app.GetContext()
	if ctx == nil {
		t.Fatalf("got nil context")
	}


	realApp := New(nil, "")
	realApp.ctx = app.GetContext() // In a real scenario, this would be set by Wails

	// This test is limited because we can't easily check for the emitted event
	// without a more complex setup. We'll just check the state change.
	realApp.setViewState(ViewStateDashboard)
	assert.Equal(t, ViewStateDashboard, realApp.state)

	ah.WaitTestToEnd(t, ch)
}
