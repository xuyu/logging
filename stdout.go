package logging

import (
	"os"
)

const (
	DefaultStdoutHandlerName = "stdout"
)

func (l *Logger) DisableDefaultStdout() {
	l.DelHandler(DefaultStdoutHandlerName)
}

func (l *Logger) EnableDefaultStdout() Handler {
	h := l.GetHandler(DefaultStdoutHandlerName)
	if h == nil {
		h = StdoutHandler()
		l.AddHandler(DefaultStdoutHandlerName, h)
	}
	return h
}

func StdoutHandler() Handler {
	return NewBaseHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
}
