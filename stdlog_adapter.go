package ulog

import (
	"bytes"
	"fmt"
	stdlog "log"
)

// type assersion
var _ Adapter = &StdLogAdapter{}

// StdLogAdapter is default ulog.Adapter
type StdLogAdapter struct {
	Level Level
}

// Handle handles ulog.Entry
func (a *StdLogAdapter) Handle(e Entry) {
	if a.Level > e.Level {
		return
	}
	var buf bytes.Buffer
	for _, f := range e.Fields() {
		buf.Write([]byte(fmt.Sprintf("\t%s=%v", f.Key, f.Value)))
	}
	stdlog.Output(1+e.CallDepth, fmt.Sprintf("%s\t%s\t%s", e.Level, e.Message, buf.String()))
}
