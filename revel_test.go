package revel

import (
	"bytes"
	"fmt"
	stdlog "log"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/revel/revel"
	"gopkg.in/birkirb/loggers.v1"
	"gopkg.in/birkirb/loggers.v1/log"
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
	revel.TRACE = stdlog.New(bb, "TRACE ", stdlog.Ldate|stdlog.Ltime)
	revel.INFO = stdlog.New(bb, "INFO  ", stdlog.Ldate|stdlog.Ltime)
	revel.WARN = stdlog.New(bb, "WARN  ", stdlog.Ldate|stdlog.Ltime)
	revel.ERROR = stdlog.New(bb, "ERROR ", stdlog.Ldate|stdlog.Ltime)
	return NewLogger(), bb
}

func TestBackTrace(t *testing.T) {
	l, b := newBufferedRevelLog()
	log.Logger = l
	log.Error("an error")
	_, file, line, _ := runtime.Caller(0)

	mustContain := fmt.Sprintf("%s:%d", filepath.Base(file), line-1)
	actual := b.String()
	if ok := strings.Contains(actual, mustContain); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, mustContain)
	}
}

func TestBackTraceF(t *testing.T) {
	l, b := newBufferedRevelLog()
	log.Logger = l
	log.Errorf("an error: %s", "value")
	_, file, line, _ := runtime.Caller(0)

	mustContain := fmt.Sprintf("%s:%d", filepath.Base(file), line-1)
	actual := b.String()
	if ok := strings.Contains(actual, mustContain); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, mustContain)
	}
}

func TestBackTraceLn(t *testing.T) {
	l, b := newBufferedRevelLog()
	log.Logger = l
	log.Errorln("an error")
	_, file, line, _ := runtime.Caller(0)

	mustContain := fmt.Sprintf("%s:%d", filepath.Base(file), line-1)
	actual := b.String()
	if ok := strings.Contains(actual, mustContain); !ok {
		t.Errorf("Log output mismatch %s (actual) != %s (expected)", actual, mustContain)
	}
}
