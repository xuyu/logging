package logging

import (
	"fmt"
	"io"
	"os"
	"path"
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
	mutex   sync.Mutex
	prefix  string
	out     io.WriteCloser
	level   int
	layout  string
	handler func(log string)
	data    map[string]string
}

func New(out io.WriteCloser, prefix string, level int, layout string) *Logger {
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

func (l *Logger) SetWriter(out io.WriteCloser) {
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
	return l.prefix
}

func (l *Logger) timestamp() string {
	return time.Now().Format(l.layout)
}

func (l *Logger) head(level string) string {
	a := []string{"[", l.timestamp(), "]"}
	p := l.Prefix()
	if p != "" {
		a = append(a, " ", p)
	}
	a = append(a, " ", level)
	return strings.Join(a, "")
}

func (l *Logger) log(level string, format string, v ...interface{}) {
	log := l.head(level) + " " + fmt.Sprintf(strings.TrimRight(format, "\n")+"\n", v...)
	if l.handler != nil {
		l.handler(log)
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	fmt.Fprint(l.out, log)
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

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func SetLayout(layout string) {
	std.SetLayout(layout)
}

func SetWriter(out io.WriteCloser) {
	std.SetWriter(out)
}

func SetDefaultLogger(l *Logger) {
	std = l
}

func mklogfile(filepath, linkpath string) (*os.File, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(filepath); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var fn string
	if err := os.Symlink(filepath, linkpath); err != nil {
		fn = filepath
	} else {
		fn = linkpath
	}
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (l *Logger) rotation(log string) {
	oldfilepath := l.data["oldfilepath"]
	linkpath := l.data["linkpath"]
	suffix := l.data["suffix"]
	filepath := strings.Join([]string{linkpath, time.Now().Format(suffix)}, ".")
	if filepath != oldfilepath {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		l.out.Close()
		file, err := mklogfile(filepath, linkpath)
		if err != nil {
			return
		}
		l.out = file
		l.data["oldfilepath"] = filepath
	}
}

func NewRotationLogger(filename string, dir string, suffix string) (*Logger, error) {
	linkpath := path.Join(dir, filename)
	filepath := strings.Join([]string{linkpath, time.Now().Format(suffix)}, ".")
	file, err := mklogfile(filepath, linkpath)
	if err != nil {
		return nil, err
	}
	l := New(file, "", DEBUG, LAYOUT)
	l.data = make(map[string]string)
	l.data["oldfilepath"] = filepath
	l.data["linkpath"] = linkpath
	l.data["suffix"] = suffix
	l.handler = l.rotation
	return l, nil
}
