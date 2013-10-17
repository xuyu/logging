package logging

import (
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
}

func DisableStdout() {
	delete(DefaultLogger.Handlers, StdoutHandlerName)
}

func EnableStdout() {
	DefaultLogger.Handlers[StdoutHandlerName] = StdoutHandler
}
