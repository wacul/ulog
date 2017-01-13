package ulog

import (
	"context"
	"fmt"
)

func Debug(ctx context.Context, args ...interface{}) {
	logLevel(ctx, DebugLevel, args...)
}
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, DebugLevel, format, args...)
}
func Error(ctx context.Context, args ...interface{}) {
	logLevel(ctx, ErrorLevel, args...)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, ErrorLevel, format, args...)
}
func Fatal(ctx context.Context, args ...interface{}) {
	logLevel(ctx, FatalLevel, args...)
}
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, FatalLevel, format, args...)
}
func Info(ctx context.Context, args ...interface{}) {
	logLevel(ctx, InfoLevel, args...)
}
func Infof(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, InfoLevel, format, args...)
}
func Warn(ctx context.Context, args ...interface{}) {
	logLevel(ctx, WarnLevel, args...)
}
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, WarnLevel, format, args...)
}

func logLevelf(ctx context.Context, lv Level, format string, args ...interface{}) {
	logLevelMessage(ctx, lv, fmt.Sprintf(format, args...))
}

func logLevel(ctx context.Context, lv Level, args ...interface{}) {
	logLevelMessage(ctx, lv, fmt.Sprint(args...))
}

func logLevelMessage(ctx context.Context, lv Level, msg string) {
	entry := ConnectorEntry{
		Level:   lv,
		Message: msg,
	}
	if ctx != nil {
		entry.Fields = fields(ctx)
	}
	currentConnector(ctx).Handle(entry)
}
