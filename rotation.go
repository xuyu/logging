package logging

import (
	"os"
	"strings"
	"time"
)

type RotationLogger struct {
	BaseLogger
	data map[string]string
}

func NewRotationLogger(shortfile string, suffix string) (*RotationLogger, error) {
	fullfile := strings.Join([]string{shortfile, time.Now().Format(suffix)}, ".")
	file, err := mklogfile(fullfile, shortfile)
	if err != nil {
		return nil, err
	}
	l := &RotationLogger{}
	l.BaseLogger = *NewBaseLogger(file, DEBUG, defaultTimeLayout)
	l.data = make(map[string]string)
	l.data["oldfilepath"] = fullfile
	l.data["linkpath"] = shortfile
	l.data["suffix"] = suffix
	l.predo = l.rotation
	return l, nil
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

func (l *RotationLogger) rotation() {
	oldfilepath := l.data["oldfilepath"]
	linkpath := l.data["linkpath"]
	suffix := l.data["suffix"]
	filepath := strings.Join([]string{linkpath, time.Now().Format(suffix)}, ".")
	if filepath != oldfilepath {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		l.out.Close()
		file, err := mklogfile(filepath, linkpath)
		if err != nil {
			return
		}
		l.out = file
		l.data["oldfilepath"] = filepath
	}
}
