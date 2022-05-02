package logger

import (
	"context"
	"fmt"
	"time"
)

const (
	levelNone int8 = iota
	levelFatal
	levelError
	levelWarn
	levelInfo
	levelDebug
	levelAll
)
const timeFormat = "2006-01-02T15:04:05.000000"

type Logger struct {
	v0 string
	v1 string
	v2 string
	v3 string
}

type columns struct {
	V0     string `json:"v0"`
	V1     string `json:"v1"`
	V2     string `json:"v2"`
	V3     string `json:"v3"`
	Level  int8   `json:"level"`
	Time   string `json:"time"`
	Msg    string `json:"msg"`
	Detail []any  `json:"detail"`
}

func convert(val any) any {
	switch v := val.(type) {
	case error:
		return v.Error()
	case fmt.Stringer:
		return Compress(v.String())
	case string:
		return Compress(v)
	case []byte:
		return Compress(v)
	default:
		return v
	}
}
func (l *Logger) stash(level int8, msg string, detail []any) {
	if level > globalLevel {
		return
	}
	for i, v := range detail {
		detail[i] = convert(v)
	}
	c := &columns{
		V0:     l.v0,
		V1:     l.v1,
		V2:     l.v2,
		V3:     l.v3,
		Level:  level,
		Time:   time.Now().Format(timeFormat),
		Msg:    msg,
		Detail: detail,
	}
	handle(c)
}

func New(v0, v1, v2, v3 string) *Logger {
	return &Logger{v0, v1, v2, v3}
}
func FromContext(c context.Context) *Logger {
	l := &Logger{}
	l.v0, _ = c.Value("v0").(string)
	l.v1, _ = c.Value("v1").(string)
	l.v2, _ = c.Value("v2").(string)
	l.v3, _ = c.Value("v3").(string)
	return l
}

func (l *Logger) Fatal(msg string, detail ...any) { l.stash(levelFatal, msg, detail) }
func (l *Logger) Error(msg string, detail ...any) { l.stash(levelError, msg, detail) }
func (l *Logger) Warn(msg string, detail ...any)  { l.stash(levelWarn, msg, detail) }
func (l *Logger) Info(msg string, detail ...any)  { l.stash(levelInfo, msg, detail) }
func (l *Logger) Debug(msg string, detail ...any) { l.stash(levelDebug, msg, detail) }
