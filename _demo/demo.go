package main

import (
	"context"
	stdlog "log"

	"github.com/tutuming/ulog"
)

func doSomething(ctx context.Context) {
	logger := ulog.Logger(ctx)
	logger.Info("this is function f")

	// log with  key-value
	logger.WithField("key1", 1).Warnf("warning! %s", "message")
}

func main() {
	// ulog uses go's standard log as default
	stdlog.SetFlags(stdlog.Lshortfile)

	ctx := context.Background()
	doSomething(ctx)

	// ulog.Logger returns type ulog.LoggerContext that also implements context.Context
	ctx = ulog.Logger(ctx).WithField("module", "app1")
	// so you can pass as context to other function
	doSomething(ctx)
}
