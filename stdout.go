package logging

import (
	"os"
)

const (
	DefaultStdoutHandlerName = "stdout"
)

func DisableDefaultStdout() {
	DelHandler(DefaultStdoutHandlerName)
}

func EnableDefaultStdout() Handler {
	h := GetHandler(DefaultStdoutHandlerName)
	if h == nil {
		h = StdoutHandler()
		AddHandler(DefaultStdoutHandlerName, h)
	}
	return h
}

func StdoutHandler() Handler {
	return NewBaseHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
}