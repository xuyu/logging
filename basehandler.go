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
	b := &BaseHandler{
		Writer:     out,
		Level:      level,
		TimeLayout: layout,
	}
	b.SetFormat(format)
	b.RecordChan = make(chan *Record, DefaultBufSize)
	b.GotError = b.CloseWhenError
	go b.BackendWriteRecord()
	return b
}

func (b *BaseHandler) SetLevel(level LogLevel) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.Level = level
}

func (b *BaseHandler) GetLevel() LogLevel {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.Level
}

func (b *BaseHandler) SetTimeLayout(layout string) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.TimeLayout = layout
}

func (b *BaseHandler) GetTimeLayout() string {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	return b.TimeLayout
}

func (b *BaseHandler) SetFormat(format string) error {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	tmpl, err := template.New("tmpl").Parse(format)
	if err != nil {
		return err
	}
	b.Tmpl = tmpl
	return nil
}

func (b *BaseHandler) Emit(level LogLevel, f string, values ...interface{}) {
	if b.GetLevel() > level {
		return
	}
	rd := &Record{
		TimeString: time.Now().Format(b.GetTimeLayout()),
		Level:      level,
		Message:    fmt.Sprintf(f, values...),
	}
	b.RecordChan <- rd
}

func (b *BaseHandler) PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func (b *BaseHandler) IgnoreError(error) {
}

func (b *BaseHandler) CloseWhenError(err error) {
	if err != nil {
		b.Writer = nil
	}
}

func (b *BaseHandler) BackendWriteRecord() {
	var rd *Record
	buf := bytes.NewBuffer(nil)
	for {
		rd = <-b.RecordChan
		if b.Writer != nil {
			buf.Reset()
			if err := b.Tmpl.Execute(buf, rd); err != nil {
				b.GotError(err)
				continue
			}
			if b.PredoFunc != nil {
				b.PredoFunc(buf)
			}
			n, err := io.Copy(b.Writer, buf)
			if err != nil {
				b.GotError(err)
			}
			if b.WriteN != nil {
				b.WriteN(int64(n))
			}
		}
	}
}
