package logging

import (
	"io"
	"os"
)

const (
	stdoutHandlerName = "stdout"
)

var (
	StdoutHandler *Handler
)

func init() {
	h, err := NewHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		panic(err)
	}
	StdoutHandler = h
	EnableStdout()
	EnableColorful()
}

func DisableStdout() {
	delete(DefaultLogger.Handlers, stdoutHandlerName)
}

func EnableStdout() {
	DefaultLogger.Handlers[stdoutHandlerName] = StdoutHandler
}

func EnableColorful() {
	StdoutHandler.Before = func(rd *Record, buf io.ReadWriter) {
		colorful(rd.Level)
	}
	StdoutHandler.After = func(*Record, int64) {
		resetColorful()
	}
}

func DisableColorful() {
	resetColorful()
	StdoutHandler.Before = nil
	StdoutHandler.After = nil
}
