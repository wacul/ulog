# ulog - structured and context oriented logging

[![Build Status](https://semaphoreci.com/api/v1/tutuming/ulog/branches/master/badge.svg)](https://semaphoreci.com/tutuming/ulog) [![GoDoc](https://godoc.org/github.com/wacul/ulog?status.svg)](https://godoc.org/github.com/wacul/ulog)


Package ulog provides a simple way to handle structured and context oriented logging.

## Examples

```go
package main

import (
	"context"
	stdlog "log"

	"github.com/wacul/ulog"
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
```

**Outoput**

```
demo.go:12:  info this is function f       
demo.go:15:  warn warning! message          key1=1
demo.go:12:  info this is function f        module=app1
demo.go:15:  warn warning! message          module=app1 key1=1
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

You can implement a custom adapter implementing simple `ulog.Adapter` interface.

There're two ways to set an adapter.

```go
// set a global adapter
// used when context has no adapter
ulog.SetDefaultAdapter(adapter)

// set a adapter used with child context
ctx = ulog.Logger(ctx).WithAdapter(adapter)
```
