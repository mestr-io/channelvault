//go:build dev && headless

package app_helpers

import (
	"path/filepath"
	"log/slog"
	"strings"
	"testing"
	"bytes"
	"fmt"		
	"os"	

	"github.com/stretchr/testify/assert"
)

func TestInitWailsLog_InitializesLogger(t *testing.T) {		
	waitTeardown, ctx, err, _ := GetSomeWailsContext(t, "")
	if err != nil {
		t.Fatalf("Error getting Wails context: %s", err.Error())
	}
	defer waitTeardown()

	// Initialize the logger with the app context
	InitWailsLog(ctx)
	defer func() {
		// Schedule undo InitWailsLog
		appLogger = nil
	}()

	if appLogger == nil {
		t.Fatal("appLogger should be initialized")
	}

	if appLogger.ctx != ctx {
		t.Fatal("appLogger context should match input context")
	}
}


func TestCoreLogfBeforeInitWailsLog(t *testing.T) {
	var buf bytes.Buffer
	th := NewCustomTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})
	logger := slog.New(th)
	slog.SetDefault(logger)

	// Reset appLogger for testing
	appLogger = nil


	tests := []struct {
		name          string
		levelLabel    string
		message       string
		logMethod     func(format string, v ...any)
	}{
		{
			name: "Logging INF before InitWailsLog",
			levelLabel: "INF",
			message: "test message 1",
			logMethod: Infof,
		},
		{
			name: "Logging DEB before InitWailsLog",
			levelLabel: "DEB",
			message: "test message 2",
			logMethod: Debugf,
		},
		{
			name: "Logging WAR before InitWailsLog",
			levelLabel: "WAR",
			message: "test message 3",
			logMethod: Warningf,
		},
		{
			name: "Logging ERR before InitWailsLog",
			levelLabel: "ERR",
			message: "test message 4",
			logMethod: Errorf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logMethod("%s", tt.message)
			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.levelLabel)
			assert.Contains(t, logOutput, tt.message)
			assert.Contains(t, logOutput, "logger_headless_test.go") // Check for source location
		})
	}
}

func TestCoreLogfAfterInitWailsLog(t *testing.T) {
	// Next tests have been commented due to a present race condition on some wails's 
	// Frontend interface implementation on dev mode. Let's wait for next version
	t.Skip()
	var buf bytes.Buffer
	th := NewCustomTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})
	logger := slog.New(th)
	slog.SetDefault(logger)

	// Reset appLogger for testing
	appLogger = nil

	tests := []struct {
		name          string
		levelLabel    string
		message       string
		logMethod     func(format string, v ...any)
	}{
		{
			name: "Logging INF after InitWailsLog",
			levelLabel: "INF",
			message: "test message 5",
			logMethod: Infof,
		},
		{
			name: "Logging DEB after InitWailsLog",
			levelLabel: "DEB",
			message: "test message 6",
			logMethod: Debugf,
		},
		{
			name: "Logging WAR after InitWailsLog",
			levelLabel: "WAR",
			message: "test message 7",
			logMethod: Warningf,			
		},
		{
			name: "Logging ERR after InitWailsLog",
			levelLabel: "ERR",
			message: "test message 8",
			logMethod: Errorf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for log file
			tempDir, err := os.MkdirTemp("", "test-files")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			logPath := filepath.Join(tempDir, "logger_headless_test.log")
			waitAppEnd, ctx, err, _ := GetSomeWailsContext(t, logPath)
			if err != nil {
				t.Fatalf("Error getting Wails context: %s", err.Error())
			}
			defer waitAppEnd()

			InitWailsLog(ctx)
			defer func() { appLogger = nil }()

			tt.logMethod("%s", tt.message)

			// Read the file logfile
			contentBytes, err := os.ReadFile(logPath)
			if err != nil {
				fmt.Printf("Error reading log file: %v\n", err)
				return
			}			
			fileContentString := string(contentBytes)
			fileContentString = strings.TrimSuffix(fileContentString, "\n")

			// Get last line
			lines := strings.Split(fileContentString, "\n")
			lastLine := ""
			for _, l := range lines {
				lastLine = l
			}

			// Log line checks
			assert.Contains(t, lastLine, tt.levelLabel)
			assert.Contains(t, lastLine, tt.message)
			assert.Contains(t, lastLine, "logger_headless_test.go") 			
		})
	}
}
