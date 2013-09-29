package logging

import (
	"os"
)

var defaultLogger Logger = NewBaseLogger(os.Stdout, DEBUG, defaultTimeLayout)

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

func SetLevel(lv level) {
	defaultLogger.SetLevel(lv)
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}
