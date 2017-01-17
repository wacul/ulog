package main

import (
	"context"
	"encoding/json"
	"fmt"
	stdlog "log"

	"github.com/tutuming/ulog"
)

type demoAdapter string

func (c demoAdapter) Handle(e ulog.Entry) {
	b, _ := json.Marshal(e.Fields())
	fmt.Println(c, e.Level.String(), e.Message, string(b))
}

func wrapLog(ctx context.Context, msg string) {
	ulog.Logger(ctx).WithCallDepth(1).Infof("wrapped : %s", msg)
}

func f(ctx context.Context) {
	logger := ulog.Logger(ctx)
	logger.Info("this is function f")

	// log with  key-value
	logger.WithField("key1", 1).Warnf("warning! %s", "message")

	// Logger implement context.Context and holds key-values
	ctx = ulog.Logger(ctx).WithField("key1", 1).WithField("key2", 2)

	wrapLog(ctx, "show this line?")
}

func main() {
	stdlog.SetFlags(stdlog.Llongfile)

	// default logger will be called
	ctx := context.Background()
	f(ctx)

	// custom logger will be called
	ulog.SetDefaultAdapter(demoAdapter("adapter(fallback)"))
	f(ctx)

	// custom logger will be called under this context
	ctx = ulog.Logger(ctx).WithAdapter(demoAdapter("adapter2"))
	f(ctx)
}
