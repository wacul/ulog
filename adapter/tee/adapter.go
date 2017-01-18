package tee

import "github.com/tutuming/ulog"

// Tee split log to multiple adapters
func Tee(adapters ...ulog.Adapter) ulog.Adapter {
	return ulog.AdapterFunc(func(e ulog.Entry) {
		e.Context = ulog.Logger(e.Context).WithCallDepth(1)
		for _, a := range adapters {
			a.Handle(e)
		}
	})
}
