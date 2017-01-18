package glog

import (
	_glog "github.com/golang/glog"
	"github.com/tutuming/ulog"
)

// Adapter for glog
var Adapter = ulog.AdapterFunc(func(e ulog.Entry) {
	depth := e.CallDepth() - 1
	msg := e.Message
	switch e.Level {
	case ulog.ErrorLevel:
		_glog.ErrorDepth(depth, msg)
	case ulog.WarnLevel:
		_glog.WarningDepth(depth, msg)
	case ulog.InfoLevel, ulog.DebugLevel: // glog doesn't have debug level
		_glog.InfoDepth(depth, msg)
	}
})
