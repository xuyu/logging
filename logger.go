package logging

import (
	"fmt"
	"runtime"
	"time"
)

func init() {
	if runtime.GOMAXPROCS(0) <= 1 {
		runtime.GOMAXPROCS(2)
	}
}

type Record struct {
	Time       time.Time
	TimeString string
	Level      LogLevel
	Message    string
}

type Emitter interface {
	Emit(Record)
}

type Logger struct {
	Handlers map[string]Emitter
}

func NewLogger() *Logger {
	return &Logger{Handlers: make(map[string]Emitter)}
}

var DefaultLogger = NewLogger()

func (l *Logger) AddHandler(name string, h Emitter) {
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

func AddHandler(name string, h Emitter) {
	DefaultLogger.AddHandler(name, h)
}

func Log(level LogLevel, format string, values ...interface{}) {
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
