package logging

import (
	"os"
)

const (
	DefaultStdoutHandlerName = "stdout"
)

func DisableDefaultStdout() {
	delete(DefaultLogger.Handlers, DefaultStdoutHandlerName)
}

func EnableDefaultStdout() Handler {
	h, exists := DefaultLogger.Handlers[DefaultStdoutHandlerName]
	if !exists {
		h = StdoutHandler()
		DefaultLogger.Handlers[DefaultStdoutHandlerName] = h
	}
	return h
}

func StdoutHandler() Handler {
	h := NewBaseHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
	h.GotError = h.PanicError
	return h
}
