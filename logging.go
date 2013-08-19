package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	DISABLE
	LAYOUT = "2006-01-02 15:04:05"
)

type Logger struct {
	mutex  sync.Mutex
	prefix string
	out    io.Writer
	level  int
	layout string
}

func New(out io.Writer, prefix string, level int, layout string) *Logger {
	return &Logger{out: out, prefix: prefix, level: level, layout: layout}
}

var std = New(os.Stdout, "", DEBUG, LAYOUT)

func (l *Logger) SetPrefix(prefix string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.prefix = prefix
}

func (l *Logger) SetLevel(level int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

func (l *Logger) SetWriter(out io.Writer) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.out = out
}

func (l *Logger) SetLayout(layout string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.layout = layout
}

func (l *Logger) Level() int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.level
}

func (l *Logger) Prefix() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.prefix
}

func (l *Logger) timestamp() string {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return time.Now().Format(l.layout)
}

func (l *Logger) head(level string) string {
	a := []string{"[", l.timestamp(), "] ", level, " ", l.Prefix(), " "}
	return strings.Join(a, "")
}

func (l *Logger) log(level string, format string, v ...interface{}) {
	head := l.head(level)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	fmt.Fprint(l.out, head)
	fmt.Fprintf(l.out, strings.TrimRight(format, "\n")+"\n", v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.Level() <= DEBUG {
		l.log("DEBUG", format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.Level() <= INFO {
		l.log("INFO", format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if l.Level() <= WARNING {
		l.log("WARNING", format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.Level() <= ERROR {
		l.log("ERROR", format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	std.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	std.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	std.Warning(format, v...)
}

func Error(format string, v ...interface{}) {
	std.Error(format, v...)
}
