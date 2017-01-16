package ulog

import (
	"bufio"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

type testAdapter struct {
	lastEntry      *LogEntry
	lastCallerFile string
	lastCallerLine int
}

func (c *testAdapter) Handle(e LogEntry) {
	c.lastEntry = &e
	_, file, line, ok := runtime.Caller(c.lastEntry.CallDepth)
	if !ok {
		panic("could not get caller")
	}
	c.lastCallerFile = file
	c.lastCallerLine = line
}

func TestEntryLevel(t *testing.T) {
	var c testAdapter
	ctx := context.Background()
	ctx = WithAdapter(ctx, &c)

	testSets := []struct {
		levelFunc  func(ctx context.Context, args ...interface{})
		levelfFunc func(ctx context.Context, format string, args ...interface{})

		wantLevel Level
	}{{
		levelFunc:  Debug,
		levelfFunc: Debugf,
		wantLevel:  DebugLevel,
	}, {
		levelFunc:  Info,
		levelfFunc: Infof,
		wantLevel:  InfoLevel,
	}, {
		levelFunc:  Warn,
		levelfFunc: Warnf,
		wantLevel:  WarnLevel,
	}, {
		levelFunc:  Error,
		levelfFunc: Errorf,
		wantLevel:  ErrorLevel,
	}}

	const testMessage = "test message"
	for _, ts := range testSets {
		ts.levelFunc(ctx, "test", "message")
		if c.lastEntry.Level != ts.wantLevel {
			t.Errorf("invalid level given want %s actual %s", ts.wantLevel, c.lastEntry.Level)
		}
		if c.lastEntry.Message != testMessage {
			t.Errorf("invalid message want %s actual %s", testMessage, c.lastEntry.Message)
		}
		ts.levelfFunc(ctx, "test %s", "message")
		if c.lastEntry.Level != ts.wantLevel {
			t.Errorf("invalid level given want %s actual %s", ts.wantLevel, c.lastEntry.Level)
		}
		if c.lastEntry.Message != testMessage {
			t.Errorf("invalid message want %s actual %s", testMessage, c.lastEntry.Message)
		}
	}
}

func infoRecursive(ctx context.Context, depth int, message string) {
	if depth == 1 {
		Info(ctx, message)
		return
	}
	infoRecursive(ctx, depth-1, message)
}

func TestCallDepth(t *testing.T) {
	var c testAdapter
	ctx := context.Background()
	ctx = WithAdapter(ctx, &c)

	checkDepthAndMarker := func(t *testing.T, wantDepth int, marker string) {
		if wantDepth != c.lastEntry.CallDepth {
			t.Errorf("invalid callDepth want %d, actual %d", wantDepth, c.lastEntry.CallDepth)
		}
		f, err := os.Open(c.lastCallerFile)
		if err != nil {
			t.Fatal(err)
		}
		r := bufio.NewReader(f)
		var aline string
		// read lines
		for i := 0; i < c.lastCallerLine; i++ {
			l, _, err := r.ReadLine()
			if err != nil {
				t.Fatal("error while reading file", err)
			}
			aline = string(l)
		}
		if !strings.Contains(aline, marker) {
			t.Errorf("%s L%d : %s  doesn't contains marker %s", c.lastCallerFile, c.lastCallerLine, aline, marker)
		}
	}

	Info(ctx, "mark normal depth")
	checkDepthAndMarker(t, defaultCallDepth, "mark normal depth")

	infoRecursive(WithAddingCallDepth(ctx, 1), 1, "mark depth 1")
	checkDepthAndMarker(t, defaultCallDepth+1, "mark depth 1")

	infoRecursive(WithAddingCallDepth(ctx, 2), 2, "mark depth 2")
	checkDepthAndMarker(t, defaultCallDepth+2, "mark depth 2")

	infoRecursive(With(ctx, AddingCallDepth(5), AddingCallDepth(7)), 12, "mark depth 5+7")
	checkDepthAndMarker(t, defaultCallDepth+12, "mark depth 5+7")
}

func TestFields(t *testing.T) {
	var c testAdapter
	ctx := context.Background()
	ctx = WithAdapter(ctx, &c)

	testSet := []struct {
		setup func(context.Context) context.Context
		want  []LogField
	}{{
		setup: func(ctx context.Context) context.Context {
			return WithField(ctx, "key1", 123)
		},
		want: []LogField{{"key1", 123}},
	}, {
		setup: func(ctx context.Context) context.Context {
			ctx = WithField(ctx, "key1", 123)
			ctx = WithField(ctx, "key2", 456)
			return ctx
		},
		want: []LogField{{"key1", 123}, {"key2", 456}},
	}, {
		setup: func(ctx context.Context) context.Context {
			ctx = WithField(ctx, "key1", 123)
			ctx = WithField(ctx, "key2", 456)
			ctx = WithField(ctx, "key1", 789)
			return ctx
		},
		want: []LogField{{"key1", 789}, {"key2", 456}},
	}}

	for _, ts := range testSet {
		ctx := ts.setup(ctx)
		Info(ctx, "message")
		if actual := c.lastEntry.Fields(); !reflect.DeepEqual(actual, ts.want) {
			t.Errorf("invalid fields want %v, actual %v", ts.want, actual)
		}
	}
}
