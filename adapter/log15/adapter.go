package log15

import (
	_log15 "github.com/inconshreveable/log15"
	"github.com/tutuming/ulog"
)

//Log15Adapter is ulog adapter for Log15
type Log15Adapter struct {
	Logger _log15.Logger
}

// New Log15Adapter
func New(logger _log15.Logger) *Log15Adapter {
	return &Log15Adapter{
		Logger: logger,
	}
}

// Handle handles ulog entry
func (c *Log15Adapter) Handle(e ulog.LogEntry) {
	var l _log15.Logger = c.Logger
	for _, f := range e.Fields() {
		l = l.New(f.Key, f.Value)
	}
	switch e.Level {
	case ulog.ErrorLevel:
		l.Error(e.Message)
	case ulog.WarnLevel:
		l.Warn(e.Message)
	case ulog.InfoLevel:
		l.Info(e.Message)
	case ulog.DebugLevel:
		l.Debug(e.Message)
	}
}
