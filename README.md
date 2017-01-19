# ulog - Structured and context oriented logging interface for Go

[![Build Status](https://semaphoreci.com/api/v1/tutuming/ulog/branches/master/badge.svg)](https://semaphoreci.com/tutuming/ulog) [![GoDoc](https://godoc.org/github.com/wacul/ulog?status.svg)](https://godoc.org/github.com/wacul/ulog)


Package ulog provides a simple way to handle structured and context oriented logging.

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

## Use `contet.Context`!

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

```
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
    "github.com/Sirupsen/logrus"
    logrusadapter "github.com/wacul/adapter/logrus"
)
//...
// set logrus adapter and fields to context
ctx = ulog.Logger(ctx).
		WithAdapter(logrusadapter.New(logrus.New())).
		WithField("processName" : "someBatch")

doSomething(ctx)
```

## Adapters

ulog has no output handler itself. As default, all logs are output via go's standard library log.

Instead, ulog provides adapter implementations for popular loggers.
(see the [adapter](./adapter) directory)

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
)

func doSomething(ctx context.Context) {
	logger := ulog.Logger(ctx)
	logger.Info("Start doSomething")

	// log with  key-value
	logger.WithField("key1", 1).Warnf("warning! %s", "message")

	logger.Info("End doSomething")
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
```

**Outoput**

```
demo.go:12:  info Start doSomething        
demo.go:15:  warn warning! message          key1=1
demo.go:17:  info End doSomething          
demo.go:12:  info Start doSomething         module=app1
demo.go:15:  warn warning! message          module=app1 key1=1
demo.go:17:  info End doSomething           module=app1
```
