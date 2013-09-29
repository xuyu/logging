package logging

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimeLayout = "2006-01-02 15:04:05"
	defaultBufSize    = 1024
)

type BaseLogger struct {
	mutex  sync.Mutex
	out    io.WriteCloser
	lv     level
	layout string
	c      chan string
	predo  func()
}

func NewBaseLogger(out io.WriteCloser, lv level, layout string) *BaseLogger {
	l := &BaseLogger{
		out:    out,
		lv:     lv,
		layout: layout,
	}
	l.c = make(chan string, defaultBufSize)
	go l.work()
	return l
}

func (l *BaseLogger) SetLevel(lv level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.lv = lv
}

func (l *BaseLogger) emit(lv level, format string, v ...interface{}) {
	l.mutex.Lock()
	local_lv := l.lv
	l.mutex.Unlock()
	if local_lv > lv {
		return
	}
	s := fmt.Sprintf("[%s] %s ", time.Now().Format(l.layout), lv.String())
	s += fmt.Sprintf(strings.TrimRight(format, "\r\n")+"\n", v...)
	l.c <- s
}

func (l *BaseLogger) work() {
	var s string
	for {
		s = <-l.c
		if l.predo != nil {
			l.predo()
		}
		io.WriteString(l.out, s)
	}
}

func (l *BaseLogger) Debug(format string, v ...interface{}) {
	l.emit(DEBUG, format, v...)
}

func (l *BaseLogger) Info(format string, v ...interface{}) {
	l.emit(INFO, format, v...)
}

func (l *BaseLogger) Warning(format string, v ...interface{}) {
	l.emit(WARNING, format, v...)
}

func (l *BaseLogger) Error(format string, v ...interface{}) {
	l.emit(ERROR, format, v...)
}
