package ulog

import "context"

type (
	fieldsKeyType    struct{}
	connectorKeyType struct{}
	fieldKey         string

	optionFunc func(ctx context.Context) context.Context
)

var (
	fieldsKey    fieldsKeyType
	connectorKey connectorKeyType
)

func With(ctx context.Context, opts ...optionFunc) context.Context {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		ctx = opt(ctx)
	}

	return ctx
}

func withField(ctx context.Context, key string, value interface{}) context.Context {
	fk := fieldKey(key)
	ctx = context.WithValue(ctx, fk, value)
	fks, _ := ctx.Value(fieldsKey).([]fieldKey)
	ctx = context.WithValue(ctx, fieldsKey, append(fks, fk))
	return ctx
}

func fields(ctx context.Context) []ConnectorField {
	keys, _ := ctx.Value(fieldsKey).([]fieldKey)
	if len(keys) == 0 {
		return nil
	}
	fs := make([]ConnectorField, len(keys))
	for i := range keys {
		fs[i] = ConnectorField{
			Key:   string(keys[i]),
			Value: ctx.Value(keys[i]),
		}
	}
	return fs
}

func withConnector(ctx context.Context, lc LoggerConnector) context.Context {
	return context.WithValue(ctx, connectorKey, lc)
}

func connector(ctx context.Context) LoggerConnector {
	lc, _ := ctx.Value(connectorKey).(LoggerConnector)
	return lc
}

func Field(key string, value interface{}) optionFunc {
	return func(ctx context.Context) context.Context {
		return withField(ctx, key, value)
	}
}

func Connector(lc LoggerConnector) optionFunc {
	return func(ctx context.Context) context.Context {
		return withConnector(ctx, lc)
	}
}
