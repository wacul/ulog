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
	lastEntry      Entry
	lastCallerFile string
	lastCallerLine int
}

func (c *testAdapter) Handle(e Entry) {
	c.lastEntry = e
	_, file, line, ok := runtime.Caller(c.lastEntry.CallDepth())
	if !ok {
		panic("could not get caller")
	}
	c.lastCallerFile = file
	c.lastCallerLine = line
}

func TestEntryLevel(t *testing.T) {
	var c testAdapter
	ctx := context.Background()
	logger := Logger(ctx).WithAdapter(&c)

	testSets := []struct {
		levelFunc  func(args ...interface{})
		levelfFunc func(format string, args ...interface{})

		wantLevel Level
	}{{
		levelFunc:  logger.Debug,
		levelfFunc: logger.Debugf,
		wantLevel:  DebugLevel,
	}, {
		levelFunc:  logger.Info,
		levelfFunc: logger.Infof,
		wantLevel:  InfoLevel,
	}, {
		levelFunc:  logger.Warn,
		levelfFunc: logger.Warnf,
		wantLevel:  WarnLevel,
	}, {
		levelFunc:  logger.Error,
		levelfFunc: logger.Errorf,
		wantLevel:  ErrorLevel,
	}}

	const testMessage = "test message"
	for _, ts := range testSets {
		ts.levelFunc("test", "message")
		if c.lastEntry.Level != ts.wantLevel {
			t.Errorf("invalid level given. want %s actual %s", ts.wantLevel, c.lastEntry.Level)
		}
		if c.lastEntry.Message != testMessage {
			t.Errorf("invalid message. want %s actual %s", testMessage, c.lastEntry.Message)
		}
		ts.levelfFunc("test %s", "message")
		if c.lastEntry.Level != ts.wantLevel {
			t.Errorf("invalid level given. want %s actual %s", ts.wantLevel, c.lastEntry.Level)
		}
		if c.lastEntry.Message != testMessage {
			t.Errorf("invalid message. want %s actual %s", testMessage, c.lastEntry.Message)
		}
	}
}

func infoRecursive(ctx context.Context, depth int, message string) {
	if depth == 1 {
		Logger(ctx).Info(message)
		return
	}
	infoRecursive(ctx, depth-1, message)
}

func TestCallDepth(t *testing.T) {
	var c testAdapter
	ctx := context.Background()
	logger := Logger(ctx).WithAdapter(&c)

	checkDepthAndMarker := func(t *testing.T, wantDepth int, marker string) {
		if wantDepth != c.lastEntry.CallDepth() {
			t.Errorf("invalid callDepth want %d, actual %d", wantDepth, c.lastEntry.CallDepth())
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

	logger.Info("mark normal depth")
	checkDepthAndMarker(t, defaultCallDepth, "mark normal depth")

	i := logger.WithCallDepth(1).Info
	i("mark closure depth")
	checkDepthAndMarker(t, defaultCallDepth+1, "mark closure depth")

	infoRecursive(logger.WithCallDepth(1), 1, "mark depth 1")
	checkDepthAndMarker(t, defaultCallDepth+1, "mark depth 1")

	infoRecursive(logger.WithCallDepth(2), 2, "mark depth 2")
	checkDepthAndMarker(t, defaultCallDepth+2, "mark depth 2")

	infoRecursive(logger.WithCallDepth(5).WithCallDepth(7), 12, "mark depth 5+7")
	checkDepthAndMarker(t, defaultCallDepth+12, "mark depth 5+7")
}

func TestFields(t *testing.T) {
	var c testAdapter
	ctx := context.Background()
	logger := Logger(ctx).WithAdapter(&c)

	testSet := []struct {
		setup func(l LoggerContext) LoggerContext
		want  []Field
	}{{
		setup: func(l LoggerContext) LoggerContext {
			return l.WithField("key1", 123)
		},
		want: []Field{{"key1", 123}},
	}, {
		setup: func(l LoggerContext) LoggerContext {
			return l.WithField("key1", 123).WithField("key2", 456)
		},
		want: []Field{{"key1", 123}, {"key2", 456}},
	}, {
		setup: func(l LoggerContext) LoggerContext {
			return l.WithField("key1", 123).WithField("key2", 456).WithField("key1", 789)
		},
		want: []Field{{"key1", 789}, {"key2", 456}},
	}}

	for _, ts := range testSet {
		l := ts.setup(logger)
		l.Info("message")
		if actual := c.lastEntry.Fields(); !reflect.DeepEqual(actual, ts.want) {
			t.Errorf("invalid fields want %v, actual %v", ts.want, actual)
		}
	}
}
