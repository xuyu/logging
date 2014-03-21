package logging

import (
	"fmt"
	"time"
)

type Record struct {
	Time        time.Time
	TimeString  string
	Level       logLevel
	Message     string
	LoggerName  string
	HandlerName string
}

type Emitter interface {
	Emit(Record)
}

type Logger struct {
	Name     string
	Handlers map[string]Emitter
}

func NewLogger() *Logger {
	return &Logger{Handlers: make(map[string]Emitter)}
}

var DefaultLogger = NewLogger()

func (l *Logger) AddHandler(name string, h Emitter) {
	l.Handlers[name] = h
}

func (l *Logger) Log(level logLevel, format string, values ...interface{}) {
	rd := Record{
		Time:       time.Now(),
		Level:      level,
		Message:    fmt.Sprintf(format, values...),
		LoggerName: l.Name,
	}
	for name, h := range l.Handlers {
		rd.HandlerName = name
		h.Emit(rd)
	}
}

func AddHandler(name string, h Emitter) {
	DefaultLogger.AddHandler(name, h)
}

func Log(level logLevel, format string, values ...interface{}) {
	DefaultLogger.Log(level, format, values...)
}

func Debug(format string, values ...interface{}) {
	DefaultLogger.Log(DEBUG, format, values...)
}

func Info(format string, values ...interface{}) {
	DefaultLogger.Log(INFO, format, values...)
}

func Warning(format string, values ...interface{}) {
	DefaultLogger.Log(WARNING, format, values...)
}

func Error(format string, values ...interface{}) {
	DefaultLogger.Log(ERROR, format, values...)
}
