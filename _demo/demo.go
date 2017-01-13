package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tutuming/ulog"
)

type demoConnector string

func (c demoConnector) Handle(e ulog.ConnectorEntry) {
	b, _ := json.Marshal(e.Fields)
	fmt.Println(c, e.Level.String(), e.Message, string(b))
}

func f(ctx context.Context) {
	ulog.Info(nil, "a")
	ulog.Info(ctx, "b")
	ulog.Warn(ctx, "c")
	ctx = ulog.With(ctx,
		ulog.Field("key1", 1),
		ulog.Field("key2", 2))
	ctx2 := ulog.With(ctx, ulog.Field("key3", 3))
	ulog.Info(ctx, "abc")
	ulog.Info(ctx2, "abc")
}

func main() {
	// default logger will be called
	ctx := context.Background()
	f(ctx)

	// custom logger will be called
	ulog.SetFallbackConnector(demoConnector("connector(fallback)"))
	f(ctx)

	// custom logger will be called under this context
	ctx = ulog.With(ctx, ulog.Connector(demoConnector("connector2")))
	f(ctx)
}
