package ulog

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

const defaultCallDepth = 4

// type assertion
var _ context.Context = &LoggerContext{context.Background()}

// LoggerContext is context bounded logger
type LoggerContext struct {
	context.Context
}

// Logger returns LoggerContext
func Logger(ctx context.Context) *LoggerContext {
	return &LoggerContext{Context: ctx}
}

// WithField returns new LoggerContext with field key-value
func (ctx *LoggerContext) WithField(key string, value interface{}) *LoggerContext {
	return &LoggerContext{
		Context: withField(ctx, key, value),
	}
}

// WithAdapter returns a new LoggerContext that holds LoggerAdapter used to logging.
func (ctx *LoggerContext) WithAdapter(la Adapter) *LoggerContext {
	return &LoggerContext{
		Context: withAdapter(ctx, la),
	}
}

// WithCallDepth returns a new conte
func (ctx *LoggerContext) WithCallDepth(depth int) *LoggerContext {
	return &LoggerContext{
		Context: withAddingCallDepth(ctx, depth),
	}
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

// utility functions

func logLevelf(ctx context.Context, lv Level, format string, args ...interface{}) {
	logLevelMessage(ctx, lv, fmt.Sprintf(format, args...))
}

func logLevel(ctx context.Context, lv Level, args ...interface{}) {
	logLevelMessage(ctx, lv, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

func logLevelMessage(ctx context.Context, lv Level, msg string) {
	entry := Entry{
		Context:   ctx,
		Level:     lv,
		Message:   msg,
		CallDepth: defaultCallDepth + callDepthFromContext(ctx),
	}
	currentAdapter(ctx).Handle(entry)
}
