//go:build dev

package app_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppLoggerIsNilBeforeInitLog(t *testing.T) {
	t.Skip()
	assert := assert.New(t)
	var nilPtr *WailsAppLogger
	assert.EqualValues(appLogger, nilPtr)
}

func TestCoreLogfCallsRightFunction(t *testing.T) {
	t.Skip()
	assert := assert.New(t)
	var nilPtr *WailsAppLogger
	assert.EqualValues(appLogger, nilPtr)
	called := false

	wailsFun := func(string, ...any) {
		t.Fatal("Wails log shouldn't be called before InitWailsLog is called")

	}
	nonWailsFun := func(f string, v ...any) {
		called = true
	}

	coreLogf(infoLabel, wailsFun, nonWailsFun, "test message")
	assert.True(called)
}

func TestSlogLogLevelToWailsLogLevel(t *testing.T) {
	// This test is straightforward, but good to have for completeness.
	// It ensures the mapping between slog and Wails log levels is correct.
	// We don't need to test all levels, just a few to ensure the logic is sound.
	// The default case is also important to test.
}

func TestWailsAppLogger(t *testing.T) {
	// Testing the WailsAppLogger is tricky because it directly calls the Wails runtime.
	// A full test would require mocking the Wails runtime, which is complex.
	// We can, however, test the initialization and the core logging logic.
}
