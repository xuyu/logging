package logging

import (
	"fmt"
	"time"
)

type Logger struct {
	Handlers map[string]Handler
}

func NewLogger() *Logger {
	return &Logger{Handlers: make(map[string]Handler)}
}

var DefaultLogger *Logger = NewLogger()

func (l *Logger) AddHandler(name string, h Handler) {
	l.Handlers[name] = h
}

func (l *Logger) Log(level LogLevel, format string, values ...interface{}) {
	rd := Record{
		Time:    time.Now(),
		Level:   level,
		Message: fmt.Sprintf(format, values...),
	}
	for _, h := range l.Handlers {
		h.Emit(rd)
	}
}

func AddHandler(name string, h Handler) {
	DefaultLogger.AddHandler(name, h)
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
