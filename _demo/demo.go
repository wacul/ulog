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
	ulog.Info(ctx, "this is function f")

	// add key-value
	ctx = ulog.With(ctx,
		ulog.Field("key1", 1),
		ulog.Field("key2", 2))

	// and more
	ctx2 := ulog.With(ctx, ulog.Field("key3", 3))

	ulog.Info(ctx, "called with ctx")   // with key1, key2
	ulog.Info(ctx2, "called with ctx2") // with key1, key2, key3
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

/* output sample

2017/01/14 00:49:57 info	this is function f
2017/01/14 00:49:57 info	called with ctx		key1=1	key2=2
2017/01/14 00:49:57 info	called with ctx2		key1=1	key2=2	key3=3
connector(fallback) info this is function f null
connector(fallback) info called with ctx [{"Key":"key1","Value":1},{"Key":"key2","Value":2}]
connector(fallback) info called with ctx2 [{"Key":"key1","Value":1},{"Key":"key2","Value":2},{"Key":"key3","Value":3}]
connector2 info this is function f null
connector2 info called with ctx [{"Key":"key1","Value":1},{"Key":"key2","Value":2}]
connector2 info called with ctx2 [{"Key":"key1","Value":1},{"Key":"key2","Value":2},{"Key":"key3","Value":3}]

*/
