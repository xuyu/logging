package logging

import (
	"os"
)

type FileLogger struct {
	BaseLogger
}

func NewFileLogger(file string) (*FileLogger, error) {
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	l := &FileLogger{}
	l.BaseLogger = *NewBaseLogger(fp, "", DEBUG, LAYOUT)
	return l, nil
}
