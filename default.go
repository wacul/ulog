package ulog

import (
	"bytes"
	"fmt"
	stdlog "log"
)

type stdlogAdapter struct{}

func (a *stdlogAdapter) Handle(e LogEntry) {
	var buf bytes.Buffer
	for _, f := range e.Fields() {
		buf.Write([]byte(fmt.Sprintf("\t%s=%v", f.Key, f.Value)))
	}
	stdlog.Output(1+e.CallDepth, fmt.Sprintf("%s\t%s\t%s", e.Level, e.Message, buf.String()))
}
