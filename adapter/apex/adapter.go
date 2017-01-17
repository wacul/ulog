package apex

import (
	_apex "github.com/apex/log"
	"github.com/tutuming/ulog"
)

// assert type
var _ ulog.Adapter = &ApexAdapter{}

//ApexAdapter is ulog adapter for apex
type ApexAdapter struct {
	Logger _apex.Interface
}

// New ApexAdapter
func New(logger _apex.Interface) *ApexAdapter {
	return &ApexAdapter{
		Logger: logger,
	}
}

// Handle handles ulog entry
func (c *ApexAdapter) Handle(e ulog.Entry) {
	var l _apex.Interface = c.Logger
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
