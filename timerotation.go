package logging

import (
	"os"
	"strings"
	"time"
)

type TimeRotationHandler struct {
	BaseHandler
	LocalData map[string]string
}

func NewTimeRotationHandler(shortfile string, suffix string) (*TimeRotationHandler, error) {
	r := &TimeRotationHandler{}
	fullfile := strings.Join([]string{shortfile, time.Now().Format(suffix)}, ".")
	file, err := r.OpenFile(fullfile, shortfile)
	if err != nil {
		return nil, err
	}
	r.BaseHandler = *NewBaseHandler(file, DEBUG, DefaultTimeLayout, DefaultFormat)
	r.LocalData = make(map[string]string)
	r.LocalData["oldfilepath"] = fullfile
	r.LocalData["linkpath"] = shortfile
	r.LocalData["suffix"] = suffix
	r.PredoFunc = r.Rotate
	return r, nil
}

func (r *TimeRotationHandler) OpenFile(filepath, linkpath string) (*os.File, error) {
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(filepath); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
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

func (r *TimeRotationHandler) Rotate() {
	oldfilepath := r.LocalData["oldfilepath"]
	linkpath := r.LocalData["linkpath"]
	suffix := r.LocalData["suffix"]
	filepath := strings.Join([]string{linkpath, time.Now().Format(suffix)}, ".")
	if filepath != oldfilepath {
		r.Writer.Close()
		file, err := r.OpenFile(filepath, linkpath)
		if err != nil {
			return
		}
		r.Writer = file
		r.LocalData["oldfilepath"] = filepath
	}
}
