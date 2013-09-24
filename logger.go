package logging

import (
	"io"
)

type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warning(string, ...interface{})
	Error(string, ...interface{})

	SetPrefix(string)
	SetLevel(level)
	SetLayout(string)
	SetWriter(io.WriteCloser)
}
