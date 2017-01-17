package ulog

import "testing"

func TestLevel(t *testing.T) {
	levels := []Level{
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
	}

	for i := 0; i < len(levels)-1; i++ {
		if levels[i] >= levels[i+1] {
			t.Errorf("%s level must be grater than %s level", levels[i].String(), levels[i+1].String())
		}
	}
}
