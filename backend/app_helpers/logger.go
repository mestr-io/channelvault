package app_helpers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

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
func InitWailsLog(ctx context.Context, level slog.Level) {
	if appLogger == nil {
		appLogger = &WailsAppLogger{ctx: ctx}
		l, ok := slogLevelToWailsLoggerLevel[level]
		if !ok {
			log.Fatalf("Invalid slog log level: %v", level)
		}
		wruntime.LogSetLogLevel(ctx, l)
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
	wf func(ctx context.Context, message string),
	format string,
	v ...any) {
	where := GetLineAndNo(5)
	s := fmt.Sprintf(format, v...)

	if al.ctx == nil {
		panic("AppLogger context is nil. Call InitLog(ctx, level) first.")
	}

	ss := fmt.Sprintf("%s: %s", where, s)
	wf(al.ctx, ss)
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
func coreLogf(label string, wf func(string, ...any), lf func(format string, v ...any), format string, v ...any) {
	if appLogger == nil {
		where := GetLineAndNo(3)
		s := fmt.Sprintf(format, v...)
		ss := fmt.Sprintf("%s | %s: %s", label, where, s)
		lf(ss)

		return
	}
	wf(format, v...)
}
