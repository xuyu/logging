package logging

import (
	"bytes"
	"fmt"
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
	GetLevel() LogLevel
	SetTimeLayout(string)
	GetTimeLayout() string
	SetFormat(string) error
	Emit(LogLevel, string, ...interface{})
}

type Record struct {
	TimeString string
	Level      LogLevel
	Message    string
}

type BaseHandler struct {
	Mutex      sync.Mutex
	Writer     io.WriteCloser
	Level      LogLevel
	TimeLayout string
	Tmpl       *template.Template
	RecordChan chan *Record
	PredoFunc  func(io.ReadWriter)
	WriteN     func(int64)
	GotError   func(error)
}

func NewBaseHandler(out io.WriteCloser, level LogLevel, layout, format string) *BaseHandler {
	h := &BaseHandler{
		Writer:     out,
		Level:      level,
		TimeLayout: layout,
	}
	h.SetFormat(format)
	h.RecordChan = make(chan *Record, DefaultBufSize)
	h.GotError = h.PanicError
	go h.BackendWriteRecord()
	return h
}

func (h *BaseHandler) SetLevel(level LogLevel) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Level = level
}

func (h *BaseHandler) GetLevel() LogLevel {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.Level
}

func (h *BaseHandler) SetTimeLayout(layout string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.TimeLayout = layout
}

func (h *BaseHandler) GetTimeLayout() string {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.TimeLayout
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

func (h *BaseHandler) Emit(level LogLevel, f string, values ...interface{}) {
	if h.GetLevel() > level {
		return
	}
	rd := &Record{
		TimeString: time.Now().Format(h.GetTimeLayout()),
		Level:      level,
		Message:    fmt.Sprintf(f, values...),
	}
	h.RecordChan <- rd
}

func (h *BaseHandler) PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func (h *BaseHandler) BackendWriteRecord() {
	rd := &Record{}
	buf := bytes.NewBuffer(nil)
	for {
		rd = <-h.RecordChan
		if h.Writer == nil {
			continue
		}
		buf.Reset()
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
