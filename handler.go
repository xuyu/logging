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
	SetLevel(LogLevel)
	SetLevelRange(LogLevel, LogLevel)
	SetTimeLayout(string)
	SetFormat(string) error
	Emit(Record)
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
	RecordChan chan *Record
	PredoFunc  func(io.ReadWriter)
	WriteN     func(int64)
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
	h.RecordChan = make(chan *Record, DefaultBufSize)
	h.GotError = h.PanicError
	go h.WriteRecord()
	return h, nil
}

func (h *BaseHandler) SetLevel(level LogLevel) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Level = level
}

func (h *BaseHandler) SetLevelRange(min_level, max_level LogLevel) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.LRange = &LevelRange{min_level, max_level}
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

func (h *BaseHandler) Emit(rd Record) {
	if h.LRange != nil {
		if !h.LRange.Contain(rd.Level) {
			return
		}
	} else if h.Level > rd.Level {
		return
	}
	h.RecordChan <- &rd
}

func (h *BaseHandler) PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func (h *BaseHandler) WriteRecord() {
	rd := &Record{}
	buf := bytes.NewBuffer(nil)
	for {
		rd = <-h.RecordChan
		if h.Writer == nil {
			continue
		}
		buf.Reset()
		rd.TimeString = rd.Time.Format(h.TimeLayout)
		if err := h.Tmpl.Execute(buf, rd); err != nil {
			h.GotError(err)
			continue
		}
		if h.PredoFunc != nil {
			h.PredoFunc(buf)
		}
		n, err := io.Copy(h.Writer, buf)
		if err != nil {
			h.GotError(err)
		}
		if h.WriteN != nil {
			h.WriteN(int64(n))
		}
	}
}
