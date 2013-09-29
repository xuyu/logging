package logging

import (
	"os"
	"strings"
	"time"
)

type RotationHandler struct {
	BaseHandler
	data map[string]string
}

func NewRotationHandler(shortfile string, suffix string) (*RotationHandler, error) {
	fullfile := strings.Join([]string{shortfile, time.Now().Format(suffix)}, ".")
	file, err := mklogfile(fullfile, shortfile)
	if err != nil {
		return nil, err
	}
	r := &RotationHandler{}
	r.BaseHandler = *NewBaseHandler(file, DEBUG, DefaultTimeLayout, DefaultFormat)
	r.data = make(map[string]string)
	r.data["oldfilepath"] = fullfile
	r.data["linkpath"] = shortfile
	r.data["suffix"] = suffix
	r.predo = r.rotation
	return r, nil
}

func mklogfile(filepath, linkpath string) (*os.File, error) {
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

func (r *RotationHandler) rotation() {
	oldfilepath := r.data["oldfilepath"]
	linkpath := r.data["linkpath"]
	suffix := r.data["suffix"]
	filepath := strings.Join([]string{linkpath, time.Now().Format(suffix)}, ".")
	if filepath != oldfilepath {
		r.out.Close()
		file, err := mklogfile(filepath, linkpath)
		if err != nil {
			return
		}
		r.out = file
		r.data["oldfilepath"] = filepath
	}
}
