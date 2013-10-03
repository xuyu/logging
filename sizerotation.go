package logging

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type FileNameInfo struct {
	FileName string
	FileInfo os.FileInfo
}

type FileNameInfoSlice struct {
	files []FileNameInfo
}

func (f *FileNameInfoSlice) Append(fi FileNameInfo) {
	f.files = append(f.files, fi)
}

func (f *FileNameInfoSlice) Len() int {
	return len(f.files)
}

func (f *FileNameInfoSlice) Less(i, j int) bool {
	return f.files[i].FileInfo.ModTime().Before(f.files[j].FileInfo.ModTime())
}

func (f *FileNameInfoSlice) Swap(i, j int) {
	f.files[j], f.files[i] = f.files[i], f.files[j]
}

func (f *FileNameInfoSlice) Sort() {
	sort.Sort(f)
}

func (f *FileNameInfoSlice) RemoveBefore(n int) {
	files := []FileNameInfo{}
	for i := 0; i < f.Len(); i++ {
		if i < n {
			os.Remove(f.files[i].FileName)
		} else {
			files = append(files, f.files[i])
		}
	}
	f.files = files
}

func (f *FileNameInfoSlice) RenameIndex(prefix string) {
	for index, fi := range f.files {
		newname := prefix + "." + strconv.Itoa(index+1)
		os.Rename(fi.FileName, newname)
	}
}

type SizeRotationHandler struct {
	*BaseHandler
	FileName    string
	CurFileSize uint64
	MaxFileSize uint64
	MaxFiles    uint32
}

func NewSizeRotationHandler(fn string, size uint64, count uint32) (*SizeRotationHandler, error) {
	h := &SizeRotationHandler{FileName: fn, MaxFileSize: size, MaxFiles: count}
	fp, err := h.OpenCreateFile(fn)
	if err != nil {
		return nil, err
	}
	h.CurFileSize, err = h.FileSize()
	if err != nil {
		fp.Close()
		return nil, err
	}
	bh, err := NewBaseHandler(fp, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		return nil, err
	}
	h.BaseHandler = bh
	h.Before = h.Rotate
	h.WriteN = h.AfterWrite
	return h, nil
}

func (h *SizeRotationHandler) OpenCreateFile(fn string) (*os.File, error) {
	return os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
}

func (h *SizeRotationHandler) FileSize() (uint64, error) {
	info, err := os.Stat(h.FileName)
	if err != nil {
		return 0, err
	}
	return uint64(info.Size()), nil
}

func (h *SizeRotationHandler) ReleaseFiles() (string, error) {
	pattern := h.FileName + ".*"
	fs, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	re, err2 := regexp.Compile("[0-9]+")
	if err2 != nil {
		return "", err2
	}
	files := &FileNameInfoSlice{}
	for _, name := range fs {
		suf := strings.TrimPrefix(name+".", h.FileName)
		if re.MatchString(suf) {
			if fileinfo, err := os.Stat(name); err == nil {
				files.Append(FileNameInfo{name, fileinfo})
			}
		}
	}
	files.Sort()
	files.RemoveBefore(files.Len() - int(h.MaxFiles))
	files.RenameIndex(h.FileName)
	release := h.FileName + "." + strconv.Itoa(files.Len()+1)
	return release, nil
}

func (h *SizeRotationHandler) AfterWrite(n int64) {
	h.CurFileSize += uint64(n)
}

func (h *SizeRotationHandler) Rotate(io.ReadWriter) {
	if h.CurFileSize < h.MaxFileSize {
		return
	}
	h.CurFileSize = 0
	h.Writer.Close()
	name, err := h.ReleaseFiles()
	if err != nil {
		h.GotError(err)
		return
	}
	if err := os.Rename(h.FileName, name); err != nil {
		h.GotError(err)
		return
	}
	fp, err := h.OpenCreateFile(h.FileName)
	if err != nil {
		h.GotError(err)
		return
	}
	h.Writer = fp
}
