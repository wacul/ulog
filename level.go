package ulog

import (
	"fmt"
	"strings"
)

// Level type
type Level uint8

// Logging levels
// numeric order (Debug < Info < Warn < Error) is guaranteed
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return ""
	}
}

var levelMap = map[string]Level{
	"debug":   DebugLevel,
	"info":    InfoLevel,
	"warn":    WarnLevel,
	"warning": WarnLevel,
	"error":   ErrorLevel,
}

// ParseLevel parses string to Level
func ParseLevel(s string) (Level, error) {
	l, ok := levelMap[strings.ToLower(s)]
	if !ok {
		return 0, fmt.Errorf("unknown level %s", s)
	}
	return l, nil
}

// MustLevel parses string to Level.
// Panics when unknwon level given
func MustLevel(s string) Level {
	l, err := ParseLevel(s)
	if err != nil {
		panic(err)
	}
	return l
}
