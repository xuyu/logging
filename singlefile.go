package logging

import (
	"os"
)

type SingleFileHandler struct {
	*Handler
}

func NewSingleFileHandler(file string) (*SingleFileHandler, error) {
	fp, err := os.OpenFile(file, FileCreateFlag, FileCreatePerm)
	if err != nil {
		return nil, err
	}
	bh, err := NewHandler(fp, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		return nil, err
	}
	return &SingleFileHandler{bh}, nil
}
