package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tutuming/ulog"
)

type demoAdapter string

func (c demoAdapter) Handle(e ulog.LogEntry) {
	b, _ := json.Marshal(e.Fields())
	fmt.Println(c, e.Level.String(), e.Message, string(b))
}

func wrapLog(ctx context.Context, msg string) {
	ulog.Info(ulog.WithAddingCallDepth(ctx, 1), "wrapped : "+msg)
}

func f(ctx context.Context) {
	ulog.Info(ctx, "this is function f")

	// add key-value
	ctx = ulog.WithField(ctx, "key1", 1)

	// and multiple field
	ctx2 := ulog.With(ctx,
		ulog.Field("key2", 2),
		ulog.Field("key3", 3))

	ulog.Info(ctx, "called with ctx")   // with key1, key2
	ulog.Info(ctx2, "called with ctx2") // with key1, key2, key3

	wrapLog(ctx, "show this line?")

	// get bounded logger
	logger := ulog.Logger(ctx2)
	logger.Info("bounded logger")
}

func main() {
	log.SetFlags(log.Llongfile)
	// default logger will be called
	ctx := context.Background()
	f(ctx)

	// custom logger will be called
	ulog.SetDefaultAdapter(demoAdapter("adapter(fallback)"))
	f(ctx)

	// custom logger will be called under this context
	ctx = ulog.With(ctx, ulog.Adapter(demoAdapter("adapter2")))
	f(ctx)
}
