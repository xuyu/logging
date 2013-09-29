package logging

import (
	"os"
)

type FileHandler struct {
	BaseHandler
}

func NewFileHandler(file string) (*FileHandler, error) {
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	f := &FileHandler{}
	f.BaseHandler = *NewBaseHandler(fp, DEBUG, DefaultTimeLayout, DefaultFormat)
	return f, nil
}
