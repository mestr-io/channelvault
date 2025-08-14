package app_helpers

import (
	"os"
	"context"
	"fmt"
	"log"
	"log/slog"
	"github.com/wailsapp/wails/v2/pkg/logger"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Define log labels for different log levels
// NOTE: wails logger uses its own labels which we mimic here

const (
	//logLevel used before InitLog is called
	defaultLogLevel = slog.LevelDebug
	fatalLabel   = "FATAL"
	infoLabel    = "INF"
	warningLabel = "WAR"
	errorLabel   = "ERR"
	debugLabel   = "DEB"
)

var slogLevelToWailsLoggerLevel = map[slog.Level]logger.LogLevel{
	slog.LevelDebug: logger.DEBUG,
	slog.LevelInfo:  logger.INFO,
	slog.LevelWarn:  logger.WARNING,
	slog.LevelError: logger.ERROR,
}

// appLogger is a non-exported singleton instance of the AppLogger
// It is initialized with the context of the application and provides methods for logging errors and messages.
var appLogger *AppLogger

type AppLogger struct {
	ctx context.Context
}

func init() {
	th := NewCustomTextHandler(os.Stdout, &slog.HandlerOptions{Level: defaultLogLevel, AddSource: false})
	logger := slog.New(th)
	slog.SetDefault(logger)
}

func InitLog(ctx context.Context, level slog.Level) {
	if appLogger == nil {
		appLogger = &AppLogger{ctx: ctx}
		l, ok := slogLevelToWailsLoggerLevel[level]
		if !ok {
			log.Fatalf("Invalid slog log level: %v", level)
		}
		wruntime.LogSetLogLevel(ctx, l)
	}	
	slog.SetLogLoggerLevel(level)
}

// Fatal logs a fatal error message and exits the application.
func (al *AppLogger) Fatalf(format string, v ...any) {
	al.coreLogf(wruntime.LogFatal, format, v...)
}

func (al *AppLogger) Infof(format string, v ...any) {
	al.coreLogf(wruntime.LogInfo, format, v...)
}

func (al *AppLogger) Warningf(format string, v ...any) {
	al.coreLogf(wruntime.LogWarning, format, v...)
}

func (al *AppLogger) Errorf(format string, v ...any) {
	al.coreLogf(wruntime.LogError, format, v...)
}

func (al *AppLogger) Debugf(format string, v ...any) {
	al.coreLogf(wruntime.LogDebug, format, v...)
}

func (al *AppLogger) coreLogf(
	wf func(ctx context.Context, message string),	
	format string,
	v ...any) {
	where := getLineAndNo(5)
	s := fmt.Sprintf(format, v...)

	if al.ctx == nil {
		panic("AppLogger context is nil. Call InitLog(ctx, level) first.")
	}

	ss := fmt.Sprintf("%s: %s", where, s)
	wf(al.ctx, ss)
}

func Fatalf(format string, v ...any) {
	coreLogf(fatalLabel, appLogger.Infof, log.Fatalf, format, v...)
}

func Infof(format string, v ...any) {
	coreLogf(infoLabel, appLogger.Infof, slog.Info, format, v...)
}

func Warningf(format string, v ...any) {
	coreLogf(warningLabel, appLogger.Warningf, slog.Warn, format, v...)
}

func Errorf(format string, v ...any) {
	coreLogf(errorLabel, appLogger.Errorf, slog.Error, format, v...)
}

func Debugf(format string, v ...any) {
	coreLogf(debugLabel, appLogger.Debugf, slog.Debug, format, v...)
}

func coreLogf(label string, wf func(string, ...any), lf func(format string, v ...any), format string, v ...any) {
	if appLogger == nil {
		where := getLineAndNo(3)
		s := fmt.Sprintf(format, v...)	
		ss := fmt.Sprintf("%s | %s: %s", label, where, s)
		lf(ss)

		return
	}	
	wf(format, v...)
}
