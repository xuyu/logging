package logging

import (
	"bytes"
	"io"
	"sync"
	"text/template"
	"time"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
	DefaultFormat     = "[{{.TimeString}}] {{.Level}} {{.Message}}\n"
	DefaultBufSize    = 1024
)

type Handler interface {
	SetBufSize(int)
	SetLevel(LogLevel)
	SetLevelString(string)
	SetLevelRange(LogLevel, LogLevel)
	SetLevelRangeString(string, string)
	SetTimeLayout(string)
	SetFormat(string) error
	SetFilter(func(*Record) bool)
	Emit(Record)
	Panic(bool)
}

type Record struct {
	Time       time.Time
	TimeString string
	Level      LogLevel
	Message    string
}

type BaseHandler struct {
	Mutex      sync.Mutex
	Writer     io.WriteCloser
	Level      LogLevel
	LRange     *LevelRange
	TimeLayout string
	Tmpl       *template.Template
	Buffer     chan *Record
	BufSize    int
	Filter     func(*Record) bool
	Before     func(io.ReadWriter)
	After      func(int64)
	GotError   func(error)
}

func NewBaseHandler(out io.WriteCloser, level LogLevel, layout, format string) (*BaseHandler, error) {
	h := &BaseHandler{
		Writer:     out,
		Level:      level,
		TimeLayout: layout,
	}
	if err := h.SetFormat(format); err != nil {
		return nil, err
	}
	h.BufSize = DefaultBufSize
	go h.WriteRecord()
	return h, nil
}

func (h *BaseHandler) SetBufSize(size int) {
	h.BufSize = size
	close(h.Buffer)
}

func (h *BaseHandler) SetLevel(level LogLevel) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Level = level
}

func (h *BaseHandler) SetLevelString(s string) {
	h.SetLevel(StringToLogLevel(s))
}

func (h *BaseHandler) SetLevelRange(min_level, max_level LogLevel) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.LRange = &LevelRange{min_level, max_level}
}

func (h *BaseHandler) SetLevelRangeString(smin, smax string) {
	h.SetLevelRange(StringToLogLevel(smin), StringToLogLevel(smax))
}

func (h *BaseHandler) SetTimeLayout(layout string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.TimeLayout = layout
}

func (h *BaseHandler) SetFormat(format string) error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	tmpl, err := template.New("tmpl").Parse(format)
	if err != nil {
		return err
	}
	h.Tmpl = tmpl
	return nil
}

func (h *BaseHandler) SetFilter(f func(*Record) bool) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Filter = f
}

func (h *BaseHandler) Emit(rd Record) {
	if h.LRange != nil {
		if !h.LRange.Contain(rd.Level) {
			return
		}
	} else if h.Level > rd.Level {
		return
	}
	h.Buffer <- &rd
}

func (h *BaseHandler) PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func (h *BaseHandler) IgnoreError(error) {
}

func (h *BaseHandler) Panic(b bool) {
	if b {
		h.GotError = h.PanicError
	} else {
		h.GotError = h.IgnoreError
	}
}

func (h *BaseHandler) WriteRecord() {
	rd := &Record{}
	buf := bytes.NewBuffer(nil)
	h.Buffer = make(chan *Record, h.BufSize)
	for {
		rd = <-h.Buffer
		if rd == nil {
			go h.WriteRecord()
			break
		}
		if h.Filter != nil && h.Filter(rd) {
			continue
		}
		if h.Writer == nil {
			continue
		}
		buf.Reset()
		rd.TimeString = rd.Time.Format(h.TimeLayout)
		if err := h.Tmpl.Execute(buf, rd); err != nil {
			h.GotError(err)
			continue
		}
		if h.Before != nil {
			h.Before(buf)
		}
		n, err := io.Copy(h.Writer, buf)
		if err != nil {
			h.GotError(err)
		}
		if h.After != nil {
			h.After(int64(n))
		}
	}
}
