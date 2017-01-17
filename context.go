package ulog

import "golang.org/x/net/context"

type (
	fieldsKeyType    struct{}
	adapterKeyType   struct{}
	callDepthKeyType struct{}
	fieldKey         string
)

var (
	fieldsKey    fieldsKeyType
	adapterKey   adapterKeyType
	callDepthKey callDepthKeyType
)

// withField creates a context that holds key-value pair to log.
func withField(ctx context.Context, key string, value interface{}) context.Context {
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

func fieldsFromContext(ctx context.Context) []Field {
	keys, _ := ctx.Value(fieldsKey).([]fieldKey)
	if len(keys) == 0 {
		return nil
	}
	fs := make([]Field, len(keys))
	for i := range keys {
		fs[i] = Field{
			Key:   string(keys[i]),
			Value: ctx.Value(keys[i]),
		}
	}
	return fs
}

// withAdapter returns a new context that holds LoggerAdapter that is used to logging.
func withAdapter(ctx context.Context, lc Adapter) context.Context {
	return context.WithValue(ctx, adapterKey, lc)
}

func adapterFromContext(ctx context.Context) Adapter {
	lc, _ := ctx.Value(adapterKey).(Adapter)
	return lc
}

func callDepthFromContext(ctx context.Context) int {
	cd, _ := ctx.Value(callDepthKey).(int)
	return cd
}

// withAddingCallDepth returns a new context that has incremented call depth to log. Used with wrapped or utilized logger functions.
func withAddingCallDepth(ctx context.Context, depth int) context.Context {
	return context.WithValue(ctx, callDepthKey, callDepthFromContext(ctx)+depth)
}
