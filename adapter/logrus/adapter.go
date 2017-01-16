package logrus

import (
	_logrus "github.com/Sirupsen/logrus"
	"github.com/tutuming/ulog"
)

//LogrusAdapter is ulog adapter for logrus
type LogrusAdapter struct {
	Logger _logrus.FieldLogger
}

// New LogrusAdapter
func New(logger _logrus.FieldLogger) *LogrusAdapter {
	return &LogrusAdapter{
		Logger: logger,
	}
}

// Handle handles ulog entry
func (c *LogrusAdapter) Handle(e ulog.LogEntry) {
	var l _logrus.FieldLogger = c.Logger
	for _, f := range e.Fields() {
		l = l.WithField(f.Key, f.Value)
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
