package ulog

import "golang.org/x/net/context"

type (
	fieldsKeyType    struct{}
	adapterKeyType   struct{}
	callDepthKeyType struct{}
	fieldKey         string

	wrapFunc func(ctx context.Context) context.Context
)

var (
	fieldsKey    fieldsKeyType
	adapterKey   adapterKeyType
	callDepthKey callDepthKeyType
)

// With is a utility function to set context values.
func With(ctx context.Context, opts ...wrapFunc) context.Context {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		ctx = opt(ctx)
	}

	return ctx
}

// WithField creates a context that holds key-value pair to log.
func WithField(ctx context.Context, key string, value interface{}) context.Context {
	fk := fieldKey(key)
	ctx = context.WithValue(ctx, fk, value)
	fks, _ := ctx.Value(fieldsKey).([]fieldKey)
	for _, fk := range fks {
		if string(fk) == key {
			// fields already contains key
			return ctx
		}
	}
	ctx = context.WithValue(ctx, fieldsKey, append(fks, fk))
	return ctx
}

// Field is a utility function passed to With method, works same as WithField
func Field(key string, value interface{}) wrapFunc {
	return func(ctx context.Context) context.Context {
		return WithField(ctx, key, value)
	}
}

func fieldsFromContext(ctx context.Context) []LogField {
	keys, _ := ctx.Value(fieldsKey).([]fieldKey)
	if len(keys) == 0 {
		return nil
	}
	fs := make([]LogField, len(keys))
	for i := range keys {
		fs[i] = LogField{
			Key:   string(keys[i]),
			Value: ctx.Value(keys[i]),
		}
	}
	return fs
}

// WithAdapter returns a new context that holds LoggerAdapter that is used to logging.
func WithAdapter(ctx context.Context, lc LoggerAdapter) context.Context {
	return context.WithValue(ctx, adapterKey, lc)
}

// Adapter is a utility function passed to With method, works same as WithAdapter
func Adapter(lc LoggerAdapter) wrapFunc {
	return func(ctx context.Context) context.Context {
		return WithAdapter(ctx, lc)
	}
}

func adapterFromContext(ctx context.Context) LoggerAdapter {
	lc, _ := ctx.Value(adapterKey).(LoggerAdapter)
	return lc
}

func callDepthFromContext(ctx context.Context) int {
	cd, _ := ctx.Value(callDepthKey).(int)
	return cd
}

// WithAddingCallDepth returns a new context that has incremented call depth to log. Used with wrapped or utilized logger functions.
func WithAddingCallDepth(ctx context.Context, depth int) context.Context {
	return context.WithValue(ctx, callDepthKey, callDepthFromContext(ctx)+depth)
}

// AddingCallDepth is a utility function passed to With method, works same as WithAddingCallDepth
func AddingCallDepth(depth int) wrapFunc {
	return func(ctx context.Context) context.Context {
		return WithAddingCallDepth(ctx, depth)
	}
}
