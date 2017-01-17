package ulog

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
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	default:
		return ""
	}
}
