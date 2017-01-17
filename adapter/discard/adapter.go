package discard

import "github.com/tutuming/ulog"

// assert type
var _ ulog.Adapter = &discardAdapter{}

type discardAdapter struct{}

// Handle handles ulog entry
func (c *discardAdapter) Handle(ulog.Entry) {}

// Discard all logs
var Discard *discardAdapter
