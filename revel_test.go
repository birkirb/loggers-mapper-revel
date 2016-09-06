package revel

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/revel/revel"
	"gopkg.in/birkirb/loggers.v1"
	"gopkg.in/birkirb/loggers.v1/mappers"
)

func TestRevelInterface(t *testing.T) {
	var _ loggers.Contextual = NewLogger()
}

func TestRevelLevelOutputWithColor(t *testing.T) {
	l, b := newBufferedRevelLog()
	l.Debugln("\x1b[30mThis text will have black color\x1b[0m")
	l.Debugln("This text will have default color")
	var expectedMatch = []string{
		"TRACE.*This text will have black color.+$",
		"TRACE.*This text will have default color",
	}
	actual := b.String()
	lines := strings.Split(actual, "\n")
	k := 1 // Offset for lines before expected

	for i, expected := range expectedMatch {
		if ok, _ := regexp.Match(expected, []byte(lines[i+k])); !ok {
			t.Errorf("Log output mismatch `%s` (actual) != `%s` (expected)", lines[i+k], expected)
		}
	}
}

func TestRevelLevelOutput(t *testing.T) {
	l, b := newBufferedRevelLog()
	l.Info("This is a test")

	expectedMatch := "INFO.*This is a test\n"
	actual := b.String()
	if ok, _ := regexp.Match(expectedMatch, []byte(actual)); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, expectedMatch)
	}
}

func TestRevelLevelfOutput(t *testing.T) {
	l, b := newBufferedRevelLog()
	l.Errorf("This is %s test", "a")

	expectedMatch := "ERROR.*This is a test\n"
	actual := b.String()
	if ok, _ := regexp.Match(expectedMatch, []byte(actual)); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, expectedMatch)
	}
}

func TestRevelLevellnOutput(t *testing.T) {
	l, b := newBufferedRevelLog()
	l.Debugln("This is a test.", "So is this.")

	expectedMatch := "TRACE.*This is a test. So is this.\n"
	actual := b.String()
	if ok, _ := regexp.Match(expectedMatch, []byte(actual)); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, expectedMatch)
	}
}

func TestRevelWithFieldsOutput(t *testing.T) {
	l, b := newBufferedRevelLog()
	l.WithFields("test", true).Warn("This is a message.")

	expectedMatch := "WARN.*This is a message. test=true\n"
	actual := b.String()
	if ok, _ := regexp.Match(expectedMatch, []byte(actual)); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, expectedMatch)
	}
}

func TestRevelWithFieldsfOutput(t *testing.T) {
	l, b := newBufferedRevelLog()
	l.WithFields("test", true, "Error", "serious").Errorf("This is a %s.", "message")

	expectedMatch := "ERROR.*This is a message.   test=true Error=serious\n"
	actual := b.String()
	if ok, _ := regexp.Match(expectedMatch, []byte(actual)); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, expectedMatch)
	}
}

func newBufferedRevelLog() (loggers.Contextual, *bytes.Buffer) {
	var b []byte
	var bb = bytes.NewBuffer(b)

	// Loggers
	revel.TRACE = log.New(bb, "TRACE ", log.Ldate|log.Ltime|log.Lshortfile)
	revel.INFO = log.New(bb, "INFO  ", log.Ldate|log.Ltime|log.Lshortfile)
	revel.WARN = log.New(bb, "WARN  ", log.Ldate|log.Ltime|log.Lshortfile)
	revel.ERROR = log.New(bb, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)
	return NewLogger(), bb
}

func TestStackTrace(t *testing.T) {
	l, b := newBufferedRevelLog()
	EnableTrace(mappers.LevelError)
	defer DisableTrace(mappers.LevelError)
	l.Error("an error")

	mustContain := "runtime/debug.Stack"
	actual := b.String()
	if ok := strings.Contains(actual, mustContain); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, mustContain)
	}
}

func TestStackTraceF(t *testing.T) {
	l, b := newBufferedRevelLog()
	EnableTrace(mappers.LevelError)
	defer DisableTrace(mappers.LevelError)
	l.Errorf("an error: %s", "value")

	mustContain := "runtime/debug.Stack"
	actual := b.String()
	if ok := strings.Contains(actual, mustContain); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, mustContain)
	}
}

func TestStackTraceLn(t *testing.T) {
	l, b := newBufferedRevelLog()
	EnableTrace(mappers.LevelError)
	defer DisableTrace(mappers.LevelError)
	l.Errorln("an error")

	mustContain := "runtime/debug.Stack"
	actual := b.String()
	if ok := strings.Contains(actual, mustContain); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, mustContain)
	}
}
