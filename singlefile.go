package logging

import (
	"os"
)

type SingleFileHandler struct {
	*BaseHandler
}

func NewSingleFileHandler(file string) (*SingleFileHandler, error) {
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	f := &SingleFileHandler{}
	f.BaseHandler = NewBaseHandler(fp, DEBUG, DefaultTimeLayout, DefaultFormat)
	return f, nil
}
