package main

import (
	"context"
	stdlog "log"

	"github.com/wacul/ulog"
	stdlog_adapter "github.com/wacul/ulog/adapter/stdlog"
)

func doSomething(ctx context.Context) {
	logger := ulog.Logger(ctx)
	logger.Info("Start doSomething")

	// log with  key-value
	logger.WithField("key1", 1).Warnf("warning! %s", "message")

	logger.Info("End doSomething")
}

func main() {
	// ulog discards all logs by default
	stdlog.SetFlags(stdlog.Lshortfile)

	ctx := context.Background()
	ctx = ulog.Logger(ctx).WithAdapter(&stdlog_adapter.Adapter{})
	doSomething(ctx)

	// ulog.Logger returns type ulog.LoggerContext that also implements context.Context
	ctx = ulog.Logger(ctx).
		// set field for children
		WithField("module", "app1").
		// and set log adapter for children
		WithAdapter(&stdlog_adapter.Adapter{Level: ulog.WarnLevel})

	// so you can pass as context to other function
	doSomething(ctx)
}
