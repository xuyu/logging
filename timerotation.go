package logging

import (
	"io"
	"os"
	"strings"
	"time"
)

type TimeRotationHandler struct {
	*BaseHandler
	LocalData map[string]string
}

func NewTimeRotationHandler(shortfile string, suffix string) (*TimeRotationHandler, error) {
	h := &TimeRotationHandler{}
	fullfile := strings.Join([]string{shortfile, time.Now().Format(suffix)}, ".")
	file, err := h.OpenFile(fullfile, shortfile)
	if err != nil {
		return nil, err
	}
	bh, err := NewBaseHandler(file, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		return nil, err
	}
	h.BaseHandler = bh
	h.Before = h.Rotate
	h.LocalData = make(map[string]string)
	h.LocalData["oldfilepath"] = fullfile
	h.LocalData["linkpath"] = shortfile
	h.LocalData["suffix"] = suffix
	return h, nil
}

func (h *TimeRotationHandler) OpenFile(filepath, linkpath string) (*os.File, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(filepath); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	os.Remove(linkpath)
	var fn string
	if err := os.Symlink(filepath, linkpath); err != nil {
		fn = filepath
	} else {
		fn = linkpath
	}
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (h *TimeRotationHandler) Rotate(*Record, io.ReadWriter) {
	oldfilepath := h.LocalData["oldfilepath"]
	linkpath := h.LocalData["linkpath"]
	suffix := h.LocalData["suffix"]
	filepath := strings.Join([]string{linkpath, time.Now().Format(suffix)}, ".")
	if filepath != oldfilepath {
		h.Writer.(io.Closer).Close()
		file, err := h.OpenFile(filepath, linkpath)
		if err != nil {
			h.set_state(false, err)
			return
		}
		h.Writer = file
		h.LocalData["oldfilepath"] = filepath
	}
}
