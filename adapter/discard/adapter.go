package discard

import "github.com/tutuming/ulog"

// Discard all logs
var Discard = ulog.AdapterFunc(func(ulog.Entry) {})
