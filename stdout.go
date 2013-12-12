package logging

import (
	"io"
	"os"
)

const (
	StdoutHandlerName = "stdout"
)

var (
	StdoutHandler Handler
)

func init() {
	bh, err := NewBaseHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		panic(err)
	}
	StdoutHandler = bh
	EnableStdout()
	EnableColorful()
}

func DisableStdout() {
	delete(DefaultLogger.Handlers, StdoutHandlerName)
}

func EnableStdout() {
	DefaultLogger.Handlers[StdoutHandlerName] = StdoutHandler
}

func EnableColorful() {
	StdoutHandler.(*BaseHandler).Before = func(rd *Record, buf io.ReadWriter) {
		colorful(rd.Level)
	}
	StdoutHandler.(*BaseHandler).After = func(*Record, int64) {
		resetColorful()
	}
}

func DisableColorful() {
	resetColorful()
	StdoutHandler.(*BaseHandler).Before = nil
	StdoutHandler.(*BaseHandler).After = nil
}
