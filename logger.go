package ulog

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

const defaultCallDepth = 4

// LoggerContext is context bounded logger
type LoggerContext interface {
	context.Context
	WithField(key string, value interface{}) LoggerContext
	WithAdapter(Adapter) LoggerContext
	WithCallDepth(depth int) LoggerContext

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

// type assertion
var _ LoggerContext = &loggerContext{}

// loggerContext is context bounded logger
type loggerContext struct {
	context.Context
}

// Logger returns LoggerContext
func Logger(ctx context.Context) LoggerContext {
	// do not wrap when the ctx is loggerContext
	if l, ok := ctx.(*loggerContext); ok {
		return l
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &loggerContext{Context: ctx}
}

// WithField returns new LoggerContext with field key-value
func (ctx *loggerContext) WithField(key string, value interface{}) LoggerContext {
	return &loggerContext{
		Context: withField(ctx, key, value),
	}
}

// WithAdapter returns a new LoggerContext that holds LoggerAdapter used to logging.
func (ctx *loggerContext) WithAdapter(la Adapter) LoggerContext {
	return &loggerContext{
		Context: withAdapter(ctx, la),
	}
}

// WithCallDepth returns a new conte
func (ctx *loggerContext) WithCallDepth(depth int) LoggerContext {
	return &loggerContext{
		Context: withAddingCallDepth(ctx, depth),
	}
}

// Error level log
func (ctx *loggerContext) Error(args ...interface{}) {
	logLevel(ctx, ErrorLevel, args...)
}

// Errorf level log with format
func (ctx *loggerContext) Errorf(format string, args ...interface{}) {
	logLevelf(ctx, ErrorLevel, format, args...)
}

// Warn level log
func (ctx *loggerContext) Warn(args ...interface{}) {
	logLevel(ctx, WarnLevel, args...)
}

// Warnf level log with format
func (ctx *loggerContext) Warnf(format string, args ...interface{}) {
	logLevelf(ctx, WarnLevel, format, args...)
}

// Info level log
func (ctx *loggerContext) Info(args ...interface{}) {
	logLevel(ctx, InfoLevel, args...)
}

// Infof level log with format
func (ctx *loggerContext) Infof(format string, args ...interface{}) {
	logLevelf(ctx, InfoLevel, format, args...)
}

// Debug level log
func (ctx *loggerContext) Debug(args ...interface{}) {
	logLevel(ctx, DebugLevel, args...)
}

// Debugf level log with format
func (ctx *loggerContext) Debugf(format string, args ...interface{}) {
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
	entry := Entry{
		Context: withAddingCallDepth(ctx, defaultCallDepth),
		Level:   lv,
		Message: msg,
	}
	currentAdapter(ctx).Handle(entry)
}
