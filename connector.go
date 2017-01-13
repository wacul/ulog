package ulog

import "context"

// Connector is interface
type LoggerConnector interface {
	Handle(ConnectorEntry)
}

var fallbackConnector LoggerConnector = &defaultConnector{}

func currentConnector(ctx context.Context) LoggerConnector {
	if ctx == nil {
		return fallbackConnector
	}
	lc := connector(ctx)
	if lc == nil {
		return fallbackConnector
	}
	return lc
}

func SetFallbackConnector(c LoggerConnector) {
	if c == nil {
		panic("passed connector is nil")
	}
	fallbackConnector = c
}

type ConnectorField struct {
	Key   string
	Value interface{}
}

type ConnectorEntry struct {
	Level   Level
	Message string
	Fields  []ConnectorField
}
