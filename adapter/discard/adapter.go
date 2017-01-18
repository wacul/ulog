package discard

import "github.com/wacul/ulog"

// Discard all logs
var Discard = ulog.AdapterFunc(func(ulog.Entry) {})
