package ulog

import "golang.org/x/net/context"

// Adapter is interface
type Adapter interface {
	Handle(Entry)
}

var defaultAdapter Adapter = &stdlogAdapter{}

func currentAdapter(ctx context.Context) Adapter {
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
func SetDefaultAdapter(c Adapter) {
	if c == nil {
		panic("passed adapter is nil")
	}
	defaultAdapter = c
}

// LogField is key-value pair to log
type Field struct {
	Key   string
	Value interface{}
}

// Entry is dataset to log, passed to LoggerAdapter's Handle method
type Entry struct {
	Context   context.Context
	Level     Level
	Message   string
	CallDepth int
}

// Fields returns log fields binded with context
func (e *Entry) Fields() []Field {
	if e.Context == nil {
		return nil
	}
	return fieldsFromContext(e.Context)
}
