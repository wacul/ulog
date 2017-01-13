package ulog

import (
	"bytes"
	"fmt"
	stdlog "log"
)

type defaultConnector struct{}

func (c *defaultConnector) Handle(e ConnectorEntry) {
	var buf bytes.Buffer
	for _, f := range e.Fields {
		buf.Write([]byte(fmt.Sprintf("\t%s=%v", f.Key, f.Value)))
	}
	stdlog.Printf("%s\t%s\t%s", e.Level, e.Message, buf.String())
}
