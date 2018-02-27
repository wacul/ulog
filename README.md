# ulog - Structured and context oriented logging interface for Go

[![Build Status](https://travis-ci.org/wacul/ulog.svg?branch=master)](https://travis-ci.org/wacul/ulog) [![GoDoc](https://godoc.org/github.com/wacul/ulog?status.svg)](https://godoc.org/github.com/wacul/ulog)


Package ulog provides a simple way to handle structured and context oriented logging and decouples package from specific log implementation.


## Pain in Logging

If you write some module or function in your project, you might want to do logging.

```go
func doSomething() {
   log.Println("start do something")
   // ...
   log.Println("end do something")
}
```

But your code may be called from server's request handler or some batch process, so it should be logged with authenticated user's ID or batch process name.

Some libraries like logrus, log15 provides a way to bind the context value to a logger object and output structured log.
So now you can write like

```go
func doSomething(logger logrus.FieldLogger) {
   logger.Info("start do something")
   // ...
   logger.Info("end do something")
}

```
```go
// from serer
doSomething(logrus.WithField("userID" : "123"))

// from batch process
doSomething(logrus.WithField("processName" : "someBatch"))
```

This looks pretty good! but there're some problem

* Everywhere you want to output log should take the logger
* Your code is strongly coupled with some logger library.
    * One project uses glog but the other new project may use logrus

## Use `context.Context` with ulog

ulog provides a way to carry a context using Go's `context` package (for 1.6 or older version `golang.org/x/net/context`).

You can write

```go
func doSomething(ctx context.Context) {
	logger := ulog.Logger(ctx)
    logger.Info("start do something")
    // ...
    logger.Info("end do something")
}
```

```go
// from server
import "github.com/wacul/ulog/adapter/glog"
// ...
// set glog adapter and fields to context
ctx = ulog.Logger(ctx).
		WithAdapter(glog.Adapter).
		WithField("userID" : "123")

doSomething(ctx)
```

```go
// from batch process
import (
    "github.com/sirupsen/logrus"
    logrusadapter "github.com/wacul/adapter/logrus"
)
//...
// set logrus adapter and fields to context
ctx = ulog.Logger(ctx).
		WithAdapter(logrusadapter.New(logrus.New())).
		WithField("processName" : "someBatch")

doSomething(ctx)
```

`ulog.Logger(ctx)` returns `ulog.LoggerContext` interface that also implements `context.Context` interface.

## Adapters

ulog has no output handler itself. By default, all logs are discarded.

To output logs, ulog provides adapter implementations for popular loggers.
(see the [adapter](./adapter) directory)

* stdlog (go's standard library log)
* logrus
* log15
* glog
* apex/log
* discard
    * discards all logs
* tee
    * splits logs to multiple adapters

There're two ways to set an adapter.

```go
// set a global adapter
// used when context has no adapter
ulog.SetDefaultAdapter(adapter)

// set a adapter used with child context
ctx = ulog.Logger(ctx).WithAdapter(adapter)
```

You can implement a custom adapter implementing simple `ulog.Adapter` interface.

## Code Example

```go
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
```

**Outoput**

```
demo.go:13:  info Start doSomething        
demo.go:16:  warn warning! message          key1=1
demo.go:18:  info End doSomething          
demo.go:16:  warn warning! message          module=app1 key1=1
```
