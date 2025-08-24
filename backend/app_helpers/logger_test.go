//go:build dev

package app_helpers

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed all:testdata/dist
var assets embed.FS //neded by FakeApp

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
	nonWailsFun := func(f string, v ...any) {
		called = true
	}

	coreLogf(infoLabel, wailsFun, nonWailsFun, "test message")
	assert.True(called)
}