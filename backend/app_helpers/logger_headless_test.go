//go:build dev && headless

package app_helpers

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitWailsLog_InitializesLogger(t *testing.T) {
	//TODO Fix this test
	t.Skip()
	waitTeardown, ctx, err := GetSomeWailsContext(t, "")
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

func TestCoreLogf(t *testing.T) {
	var buf bytes.Buffer
	th := NewCustomTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})
	logger := slog.New(th)
	slog.SetDefault(logger)

	// Reset appLogger for testing
	appLogger = nil

	t.Run("Logging INF before InitWailsLog", func(t *testing.T) {
		buf.Reset()
		Infof("test message %d", 1)
		logOutput := buf.String()
		assert.Contains(t, logOutput, "INF")
		assert.Contains(t, logOutput, "test message 1")
		assert.Contains(t, logOutput, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging DEB before InitWailsLog", func(t *testing.T) {
		buf.Reset()
		Debugf("test message %d", 2)
		logOutput := buf.String()
		assert.Contains(t, logOutput, "DEB")
		assert.Contains(t, logOutput, "test message 2")
		assert.Contains(t, logOutput, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging WAR before InitWailsLog", func(t *testing.T) {
		buf.Reset()
		Warningf("test message %d", 3)
		logOutput := buf.String()
		assert.Contains(t, logOutput, "WAR")
		assert.Contains(t, logOutput, "test message 3")
		assert.Contains(t, logOutput, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging ERR before InitWailsLog", func(t *testing.T) {
		buf.Reset()
		Errorf("test message %d", 4)
		logOutput := buf.String()
		assert.Contains(t, logOutput, "ERR")
		assert.Contains(t, logOutput, "test message 4")
		assert.Contains(t, logOutput, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging INF after InitWailsLog", func(t *testing.T) {
		// Create a temporary directory for log file
		tempDir, err := os.MkdirTemp("", "test-files")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "logger_headless_test.log")
		waitTeardown, ctx, err := GetSomeWailsContext(t, logPath)
		if err != nil {
			t.Fatalf("Error getting Wails context: %s", err.Error())
		}
		defer waitTeardown()

		InitWailsLog(ctx)
		defer func() { appLogger = nil }()

		Infof("test message 5")

		// Read the file logfile
		contentBytes, err := os.ReadFile(logPath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		// Convert the byte slice to a string
		fileContentString := string(contentBytes)

		//eliminate last \n
		fileContentString = strings.TrimSuffix(fileContentString, "\n")

		//get last line
		lines := strings.Split(fileContentString, "\n")
		lastLine := ""
		for _, l := range lines {
			lastLine = l
		}

		//log line checks
		assert.Contains(t, lastLine, "INF")
		assert.Contains(t, lastLine, "test message 5")
		assert.Contains(t, lastLine, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging DEB after InitWailsLog", func(t *testing.T) {
		//TODO Fix this test
		t.Skip()
		// Create a temporary directory for log file
		tempDir, err := os.MkdirTemp("", "test-files")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "logger_headless_test.log")
		waitTeardown, ctx, err := GetSomeWailsContext(t, logPath)
		if err != nil {
			t.Fatalf("Error getting Wails context: %s", err.Error())
		}
		defer waitTeardown()

		InitWailsLog(ctx)
		defer func() { appLogger = nil }()

		Debugf("test message 6")

		// Read the file logfile
		contentBytes, err := os.ReadFile(logPath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		// Convert the byte slice to a string
		fileContentString := string(contentBytes)

		//eliminate last \n
		fileContentString = strings.TrimSuffix(fileContentString, "\n")

		//get last line
		lines := strings.Split(fileContentString, "\n")
		lastLine := ""
		for _, l := range lines {
			lastLine = l
		}

		//log line checks
		assert.Contains(t, lastLine, "DEB")
		assert.Contains(t, lastLine, "test message 6")
		assert.Contains(t, lastLine, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging WAR after InitWailsLog", func(t *testing.T) {
		//TODO Fix this test
		t.Skip()
		// Create a temporary directory for log file
		tempDir, err := os.MkdirTemp("", "test-files")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "logger_headless_test.log")
		waitTeardown, ctx, err := GetSomeWailsContext(t, logPath)
		if err != nil {
			t.Fatalf("Error getting Wails context: %s", err.Error())
		}
		defer waitTeardown()

		InitWailsLog(ctx)
		defer func() { appLogger = nil }()

		Warningf("test message 7")

		// Read the file logfile
		contentBytes, err := os.ReadFile(logPath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		// Convert the byte slice to a string
		fileContentString := string(contentBytes)

		//eliminate last \n
		fileContentString = strings.TrimSuffix(fileContentString, "\n")

		//get last line
		lines := strings.Split(fileContentString, "\n")
		lastLine := ""
		for _, l := range lines {
			lastLine = l
		}

		//log line checks
		assert.Contains(t, lastLine, "WAR")
		assert.Contains(t, lastLine, "test message 7")
		assert.Contains(t, lastLine, "logger_headless_test.go") // Check for source location
	})

	t.Run("Logging ERR after InitWailsLog", func(t *testing.T) {
		//TODO Fix this test
		t.Skip()
		// Create a temporary directory for log file
		tempDir, err := os.MkdirTemp("", "test-files")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "logger_headless_test.log")
		waitTeardown, ctx, err := GetSomeWailsContext(t, logPath)
		if err != nil {
			t.Fatalf("Error getting Wails context: %s", err.Error())
		}
		defer waitTeardown()

		InitWailsLog(ctx)
		defer func() { appLogger = nil }()

		Errorf("test message 8")

		// Read the file logfile
		contentBytes, err := os.ReadFile(logPath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		// Convert the byte slice to a string
		fileContentString := string(contentBytes)

		//eliminate last \n
		fileContentString = strings.TrimSuffix(fileContentString, "\n")

		//get last line
		lines := strings.Split(fileContentString, "\n")
		lastLine := ""
		for _, l := range lines {
			lastLine = l
		}

		//log line checks
		assert.Contains(t, lastLine, "ERR")
		assert.Contains(t, lastLine, "test message 8")
		assert.Contains(t, lastLine, "logger_headless_test.go") // Check for source location
	})

}
