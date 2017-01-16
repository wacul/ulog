package ulog

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

const defaultCallDepth = 4

// LoggerContext is context bounded logger
type LoggerContext struct {
	context.Context
}

// Logger returns the logger object contains context
func Logger(ctx context.Context) *LoggerContext {
	return &LoggerContext{ctx}
}

// Error level log
func (ctx *LoggerContext) Error(args ...interface{}) {
	logLevel(ctx, ErrorLevel, args...)
}

// Errorf level log with format
func (ctx *LoggerContext) Errorf(format string, args ...interface{}) {
	logLevelf(ctx, ErrorLevel, format, args...)
}

// Warn level log
func (ctx *LoggerContext) Warn(args ...interface{}) {
	logLevel(ctx, WarnLevel, args...)
}

// Warnf level log with format
func (ctx *LoggerContext) Warnf(format string, args ...interface{}) {
	logLevelf(ctx, WarnLevel, format, args...)
}

// Info level log
func (ctx *LoggerContext) Info(args ...interface{}) {
	logLevel(ctx, InfoLevel, args...)
}

// Infof level log with format
func (ctx *LoggerContext) Infof(format string, args ...interface{}) {
	logLevelf(ctx, InfoLevel, format, args...)
}

// Debug level log
func (ctx *LoggerContext) Debug(args ...interface{}) {
	logLevel(ctx, DebugLevel, args...)
}

// Debugf level log with format
func (ctx *LoggerContext) Debugf(format string, args ...interface{}) {
	logLevelf(ctx, DebugLevel, format, args...)
}

// Top level functions

// Error level log
func Error(ctx context.Context, args ...interface{}) {
	logLevel(ctx, ErrorLevel, args...)
}

// Errorf level log with format
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, ErrorLevel, format, args...)
}

// Warn level log
func Warn(ctx context.Context, args ...interface{}) {
	logLevel(ctx, WarnLevel, args...)
}

// Warnf level log with format
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, WarnLevel, format, args...)
}

// Info level log
func Info(ctx context.Context, args ...interface{}) {
	logLevel(ctx, InfoLevel, args...)
}

// Infof level log with format
func Infof(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, InfoLevel, format, args...)
}

// Debug level log
func Debug(ctx context.Context, args ...interface{}) {
	logLevel(ctx, DebugLevel, args...)
}

// Debugf level log with format
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logLevelf(ctx, DebugLevel, format, args...)
}

// utility functions

func logLevelf(ctx context.Context, lv Level, format string, args ...interface{}) {
	logLevelMessage(ctx, lv, fmt.Sprintf(format, args...))
}

func logLevel(ctx context.Context, lv Level, args ...interface{}) {
	logLevelMessage(ctx, lv, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

func logLevelMessage(ctx context.Context, lv Level, msg string) {
	entry := LogEntry{
		Context:   ctx,
		Level:     lv,
		Message:   msg,
		CallDepth: defaultCallDepth + callDepthFromContext(ctx),
	}
	currentAdapter(ctx).Handle(entry)
}
