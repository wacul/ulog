package ulog

import "context"

// LoggerAdapter is interface
type LoggerAdapter interface {
	Handle(LogEntry)
}

var defaultAdapter LoggerAdapter = &stdlogAdapter{}

func currentAdapter(ctx context.Context) LoggerAdapter {
	if ctx == nil {
		return defaultAdapter
	}
	lc := adapterFromContext(ctx)
	if lc == nil {
		return defaultAdapter
	}
	return lc
}

// SetDefaultAdapter sets the LoggerAdapter c for default adapter, that will be used when the adapter is not set or context is nil.
func SetDefaultAdapter(c LoggerAdapter) {
	if c == nil {
		panic("passed adapter is nil")
	}
	defaultAdapter = c
}

// LogField is key-value pair to log
type LogField struct {
	Key   string
	Value interface{}
}

// LogEntry is dataset to log, passed to LoggerAdapter's Handle method
type LogEntry struct {
	Context   context.Context
	Level     Level
	Message   string
	CallDepth int
}

// Fields returns log fields binded with context
func (e *LogEntry) Fields() []LogField {
	if e.Context == nil {
		return nil
	}
	return fieldsFromContext(e.Context)
}
