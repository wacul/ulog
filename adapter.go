package ulog

import "golang.org/x/net/context"

// Adapter is interface
type Adapter interface {
	Handle(Entry)
}

// AdapterFunc wraps func to Adapter
type AdapterFunc func(Entry)

// Handle calls f
func (f AdapterFunc) Handle(e Entry) {
	// add call depth
	e.Context = withAddingCallDepth(e.Context, 1)
	f(e)
}

// discards all logs by default
var defaultAdapter Adapter = AdapterFunc(func(Entry) {})

// currentAdapter returns an adapter set in current context. if Adapter is not set
// returns a default adapter (set with SetDefaultAdapter)
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

// SetDefaultAdapter sets the Adapter for default adapter which will be used when the adapter is not set in context or context is nil.
func SetDefaultAdapter(a Adapter) {
	if a == nil {
		panic("passed adapter is nil")
	}
	defaultAdapter = a
}

// Field is key-value pair to log
type Field struct {
	Key   string
	Value interface{}
}

// Entry is dataset to log, passed to LoggerAdapter's Handle method
type Entry struct {
	context.Context
	Level   Level
	Message string
}

// CallDepth retruns a numver of caller depth
func (e *Entry) CallDepth() int {
	return callDepthFromContext(e)
}

// Fields returns log fields binded with context
func (e *Entry) Fields() []Field {
	if e.Context == nil {
		return nil
	}
	return fieldsFromContext(e)
}
