package app_helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Define log labels for different log levels
// NOTE: wails logger uses its own labels which we mimic here

const (
	//logLevel used before InitLog is called
	defaultLogLevel = slog.LevelDebug

	//log labels used in slog
	fatalLabel   = "FATAL"
	infoLabel    = "INF"
	warningLabel = "WAR"
	errorLabel   = "ERR"
	debugLabel   = "DEB"
)

// slogLevelToWailsLoggerLevel maps slog levels to Wails logger levels.
var slogLevelToWailsLoggerLevel = map[slog.Level]logger.LogLevel{
	slog.LevelDebug: logger.DEBUG,
	slog.LevelInfo:  logger.INFO,
	slog.LevelWarn:  logger.WARNING,
	slog.LevelError: logger.ERROR,
}

// SlogLogLevelToWailsLogLevel translates log level constants from slog library to wails log
// library. If no mapping is found logger.DEBUG is returned.
func SlogLogLevelToWailsLogLevel(in slog.Level) (out logger.LogLevel) {
	out, ok := slogLevelToWailsLoggerLevel[in]
	if !ok {
		return logger.DEBUG
	}
	return out
}

// appLogger is a non-exported singleton instance of the AppLogger
// It is initialized with the context of the application and provides methods for logging errors and messages.
var appLogger *WailsAppLogger

// WailsAppLogger is a custom logger that uses slog for structured logging.
type WailsAppLogger struct {
	ctx context.Context
}

func init() {
	th := NewCustomTextHandler(os.Stdout, &slog.HandlerOptions{Level: defaultLogLevel, AddSource: false})
	logger := slog.New(th)
	slog.SetDefault(logger)
}

// InitWailsLog initializes the Wails logger with the provided context and log level.
func InitWailsLog(ctx context.Context) {
	if appLogger == nil {
		appLogger = &WailsAppLogger{ctx: ctx}

		//NOTE: LogSetLogLevel is called by CreateApp with the level
		//according to the build type, to override with "level" (type slog.Level):
		//	SlogLogLevelToWailsLogLevel(level)
		//	wruntime.LogSetLogLevel(ctx, l)
	}
}

// Fatal logs a fatal error message and exits the application.
func (al *WailsAppLogger) Fatalf(format string, v ...any) {
	al.coreLogf(wruntime.LogFatal, format, v...)
}

// Infof logs an informational message.
func (al *WailsAppLogger) Infof(format string, v ...any) {
	al.coreLogf(wruntime.LogInfo, format, v...)
}

// Warningf logs a warning message.
func (al *WailsAppLogger) Warningf(format string, v ...any) {
	al.coreLogf(wruntime.LogWarning, format, v...)
}

// Errorf logs an error message.
func (al *WailsAppLogger) Errorf(format string, v ...any) {
	al.coreLogf(wruntime.LogError, format, v...)
}

// Debugf logs a debug message.
func (al *WailsAppLogger) Debugf(format string, v ...any) {
	al.coreLogf(wruntime.LogDebug, format, v...)
}

// coreLogf is a helper function that formats the log message with file and line number information
func (al *WailsAppLogger) coreLogf(
	wailsLogFun func(ctx context.Context, message string),
	format string,
	v ...any) {
	where := GetLineAndNo(5)
	s := fmt.Sprintf(format, v...)

	if al.ctx == nil {
		panic("AppLogger context is nil. Call InitLog(ctx, level) first.")
	}

	ss := fmt.Sprintf("%s: %s", where, s)
	wailsLogFun(al.ctx, ss)
}

/*
 Following methods can be used at any time. If Wails logger has been initialized it will use it,
  otherwise it will use standard slog package methods with Wails default format)
*/

// Fatalf logs a fatal error message and exits the application.
func Fatalf(format string, v ...any) {
	coreLogf(fatalLabel, appLogger.Fatalf, log.Fatalf, format, v...)
}

// Infof logs an informational message.
func Infof(format string, v ...any) {
	coreLogf(infoLabel, appLogger.Infof, slog.Info, format, v...)
}

// Warningf logs a warning message.
func Warningf(format string, v ...any) {
	coreLogf(warningLabel, appLogger.Warningf, slog.Warn, format, v...)
}

// Errorf logs an error message.
func Errorf(format string, v ...any) {
	coreLogf(errorLabel, appLogger.Errorf, slog.Error, format, v...)
}

// Debugf logs a debug message.
func Debugf(format string, v ...any) {
	coreLogf(debugLabel, appLogger.Debugf, slog.Debug, format, v...)
}

// coreLogf is a helper function that formats the log message with file and line number information
func coreLogf(label string, wailsLogFun func(string, ...any), nonWailsLogFun func(format string, v ...any), format string, v ...any) {
	if appLogger == nil {
		where := GetLineAndNo(3)
		s := fmt.Sprintf(format, v...)
		ss := fmt.Sprintf("%s | %s: %s", label, where, s)
		nonWailsLogFun(ss)

		return
	}
	wailsLogFun(format, v...)
}

func DebugJSONUnmarshal(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err == nil {
		return nil
	}

	// Check if it's a Type Mismatch error
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		fmt.Println("---------------------------------------------------")
		fmt.Printf("❌ JSON TYPE MISMATCH:\n")
		fmt.Printf("   • Field:    %s\n", typeErr.Field) // The struct field, e.g., "fields.name"
		fmt.Printf("   • Expected: %v\n", typeErr.Type)  // What Go wanted, e.g., "int"
		fmt.Printf("   • Got:      %s\n", typeErr.Value) // What JSON had, e.g., "string"

		// Grab context around the error (e.g., 20 chars before/after)
		start := typeErr.Offset - 100
		if start < 0 {
			start = 0
		}
		end := typeErr.Offset + 100
		if end > int64(len(data)) {
			end = int64(len(data))
		}

		snippet := string(data[start:end])
		// Add a marker to show exactly where the error is
		pointer := strings.Repeat(" ", int(typeErr.Offset-start)) + "^"

		fmt.Printf("   • Context:  ...%s...\n", snippet)
		fmt.Printf("                  %s\n", pointer)
		fmt.Println("---------------------------------------------------")
	} else {
		// Handle syntax errors (missing braces, commas, etc.)
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("❌ JSON SYNTAX ERROR at byte %d: %s\n", syntaxErr.Offset, syntaxErr.Error())
		} else {
			fmt.Printf("❌ JSON ERROR: %s\n", err.Error())
		}
	}

	return err
}

func DescribeJSONError(err error, data []byte) string {
	if err == nil {
		return ""
	}

	var sb strings.Builder

	// Check if it's a Type Mismatch error (the most common "Field" error)
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		fmt.Fprintf(&sb, "❌ JSON TYPE MISMATCH:\n")
		fmt.Fprintf(&sb, "   • Field:    %s\n", typeErr.Field)
		fmt.Fprintf(&sb, "   • Expected: %v\n", typeErr.Type)
		fmt.Fprintf(&sb, "   • Got:      %s\n", typeErr.Value)

		// Calculate context window (20 chars before/after)
		start := int(typeErr.Offset) - 100
		if start < 0 {
			start = 0
		}
		end := int(typeErr.Offset) + 100
		if end > len(data) {
			end = len(data)
		}

		// Create the snippet and the pointer line
		snippet := string(data[start:end])
		// Calculate spaces needed for the pointer: (Offset - start)
		ptrOffset := int(typeErr.Offset) - start
		pointer := strings.Repeat(" ", ptrOffset) + "^"

		fmt.Fprintf(&sb, "   • Context:  ...%s...\n", snippet)
		fmt.Fprintf(&sb, "                  %s\n", pointer)

	} else if syntaxErr, ok := err.(*json.SyntaxError); ok {
		// Handle Syntax Errors (e.g., missing commas)
		fmt.Fprintf(&sb, "❌ JSON SYNTAX ERROR at byte %d: %s\n", syntaxErr.Offset, syntaxErr.Error())
	} else {
		// Generic errors
		fmt.Fprintf(&sb, "❌ JSON ERROR: %s\n", err.Error())
	}

	return sb.String()
}
