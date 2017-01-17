package glog

import (
	_glog "github.com/golang/glog"
	"github.com/tutuming/ulog"
)

// assert type
var _ ulog.Adapter = &GlogAdapter{}

// GlogAdapter is ulog adapter for glog
type GlogAdapter struct{}

// New GlogAdapter
func New() *GlogAdapter {
	return &GlogAdapter{}
}

// Handle handles ulog entry
func (c *GlogAdapter) Handle(e ulog.Entry) {
	depth := e.CallDepth - 1
	msg := e.Message
	switch e.Level {
	case ulog.ErrorLevel:
		_glog.ErrorDepth(depth, msg)
	case ulog.WarnLevel:
		_glog.WarningDepth(depth, msg)
	case ulog.InfoLevel, ulog.DebugLevel: // glog doesn't have debug level
		_glog.InfoDepth(depth, msg)
	}
}
