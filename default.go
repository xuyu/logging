package logging

import (
	"io"
	"os"
)

var defaultLogger Logger = NewBaseLogger(os.Stdout, "", DEBUG, LAYOUT)

func Debug(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warning(format string, v ...interface{}) {
	defaultLogger.Warning(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func SetPrefix(prefix string) {
	defaultLogger.SetPrefix(prefix)
}

func SetLevel(lv level) {
	defaultLogger.SetLevel(lv)
}

func SetLayout(layout string) {
	defaultLogger.SetLayout(layout)
}

func SetWriter(out io.WriteCloser) {
	defaultLogger.SetWriter(out)
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}
