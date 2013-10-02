package logging

import (
	"os"
)

const (
	StdoutHandlerName = "stdout"
)

var (
	StdoutHandler = NewBaseHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
)

func DisableStdout() {
	delete(DefaultLogger.Handlers, StdoutHandlerName)
}

func EnableStdout() {
	DefaultLogger.Handlers[StdoutHandlerName] = StdoutHandler
}
