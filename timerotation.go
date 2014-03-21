package logging

import (
	"io"
	"os"
	"strings"
	"time"
)

type TimeRotationHandler struct {
	*Handler
	LocalData map[string]string
}

func NewTimeRotationHandler(shortfile string, suffix string) (*TimeRotationHandler, error) {
	h := &TimeRotationHandler{}
	fullfile := strings.Join([]string{shortfile, time.Now().Format(suffix)}, ".")
	file, err := h.openFile(fullfile, shortfile)
	if err != nil {
		return nil, err
	}
	bh, err := NewHandler(file, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		return nil, err
	}
	h.Handler = bh
	h.Before = h.rotate
	h.LocalData = make(map[string]string)
	h.LocalData["oldfilepath"] = fullfile
	h.LocalData["linkpath"] = shortfile
	h.LocalData["suffix"] = suffix
	return h, nil
}

func (h *TimeRotationHandler) openFile(filepath, linkpath string) (*os.File, error) {
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
	file, err := os.OpenFile(fn, FileCreateFlag, FileCreatePerm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (h *TimeRotationHandler) rotate(*Record, io.ReadWriter) {
	filepath := h.LocalData["linkpath"] + "." + time.Now().Format(h.LocalData["suffix"])
	if filepath != h.LocalData["oldfilepath"] {
		h.writer.(io.Closer).Close()
		file, err := h.openFile(filepath, h.LocalData["linkpath"])
		if err != nil {
			return
		}
		h.writer = file
		h.LocalData["oldfilepath"] = filepath
	}
}
