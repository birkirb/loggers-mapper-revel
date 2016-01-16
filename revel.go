package revel

import (
	"fmt"
	"strings"

	"github.com/revel/revel"
	"gopkg.in/birkirb/loggers.v1"
	"gopkg.in/birkirb/loggers.v1/mappers"
)

// Logger is a Contextual logger wrapper over Revel's logger.
type Logger struct{}

// NewLogger returns a Contextual Logger for revel's internal logger.
// Note that Revel's loggers must be initialized before any logging can be made.
func NewLogger() loggers.Contextual {
	var l *Logger
	var a = mappers.NewContextualMap(l)
	a.Info("Now using Revel's logger package (via loggers/mappers/revel).")
	return a
}

// LevelPrint is a Mapper method
func (l *Logger) LevelPrint(lev mappers.Level, i ...interface{}) {
	getRevelLevel(lev).Print(i...)
}

// LevelPrintf is a Mapper method
func (l *Logger) LevelPrintf(lev mappers.Level, format string, i ...interface{}) {
	getRevelLevel(lev).Printf(format, i...)
}

// LevelPrintln is a Mapper method
func (l *Logger) LevelPrintln(lev mappers.Level, i ...interface{}) {
	getRevelLevel(lev).Println(i...)
}

// WithField returns an advanced logger with a pre-set field.
func (l *Logger) WithField(key string, value interface{}) loggers.Advanced {
	return l.WithFields(key, value)
}

// WithFields returns an advanced logger with pre-set fields.
func (l *Logger) WithFields(fields ...interface{}) loggers.Advanced {
	s := make([]string, len(fields)/2)
	for i := 0; i+1 < len(fields); i = i + 2 {
		key := fields[i]
		value := fields[i+1]
		s = append(s, fmt.Sprint(key, "=", value))
	}

	r := revelPostfixLogger{strings.Join(s, " ")}
	return mappers.NewAdvancedMap(&r)
}

type revelPostfixLogger struct {
	postfix string
}

func (r *revelPostfixLogger) LevelPrint(lev mappers.Level, i ...interface{}) {
	i = append(i, r.postfix)
	getRevelLevel(lev).Print(i...)
}

func (r *revelPostfixLogger) LevelPrintf(lev mappers.Level, format string, i ...interface{}) {
	if len(r.postfix) > 0 {
		format = format + " %s"
		i = append(i, r.postfix)
	}
	getRevelLevel(lev).Printf(format, i...)
}

func (r *revelPostfixLogger) LevelPrintln(lev mappers.Level, i ...interface{}) {
	i = append(i, r.postfix)
	getRevelLevel(lev).Println(i...)
}

func getRevelLevel(lev mappers.Level) loggers.Standard {
	switch lev {
	case mappers.LevelDebug:
		return revel.TRACE
	case mappers.LevelInfo:
		return revel.INFO
	case mappers.LevelWarn:
		return revel.WARN
	case mappers.LevelError:
		return revel.ERROR
	case mappers.LevelFatal:
		return revel.ERROR
	case mappers.LevelPanic:
		return revel.ERROR
	default:
		panic("unreachable")
	}
}
