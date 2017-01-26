package stdlog

import (
	"bytes"
	"fmt"
	"log"

	"github.com/wacul/ulog"
)

// type assersion
var _ ulog.Adapter = &Adapter{}

// Adapter is ulog adapter for go's standard log
type Adapter struct {
	Level ulog.Level
}

// Handle handles ulog.Entry
func (a *Adapter) Handle(e ulog.Entry) {
	if a.Level > e.Level {
		return
	}

	level := e.Level.String()

	var b bytes.Buffer
	fmt.Fprintf(&b, "%5s %-25s", level, e.Message)

	for _, f := range e.Fields() {
		fmt.Fprintf(&b, " %s=%v", f.Key, f.Value)
	}

	var buf bytes.Buffer
	for _, f := range e.Fields() {
		buf.Write([]byte(fmt.Sprintf("\t%s=%v", f.Key, f.Value)))
	}
	log.Output(1+e.CallDepth(), b.String())
}
