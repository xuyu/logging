package logging

import (
	"bytes"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"
	"time"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
	DefaultFormat     = "[{{.TimeString}}] {{.Level}} {{.Message}}\n"
	FormatNoTime      = "{{.Level}} {{.Message}}\n"
	FormatNoLevel     = "[{{.TimeString}}] {{.Message}}\n"
	FormatOnlyMessage = "{{.Message}}\n"
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
}

type Record struct {
	Time       time.Time
	TimeString string
	Level      LogLevel
	Message    string
}

type BaseHandler struct {
	Mutex      sync.Mutex
	State      bool
	LastError  error
	Writer     io.Writer
	Level      LogLevel
	LRange     *LevelRange
	TimeLayout string
	Tmpl       *template.Template
	Buffer     chan *Record
	BufSize    int
	Filter     func(*Record) bool
	Before     func(*Record, io.ReadWriter)
	After      func(*Record, int64)
}

func NewBaseHandler(out io.Writer, level LogLevel, layout, format string) (*BaseHandler, error) {
	h := &BaseHandler{
		State:      true,
		Writer:     out,
		Level:      level,
		TimeLayout: layout,
	}
	if err := h.SetFormat(format); err != nil {
		return nil, err
	}
	h.BufSize = DefaultBufSize
	h.Buffer = make(chan *Record, h.BufSize)
	go h.notify()
	go h.WriteRecord()
	return h, nil
}

func (h *BaseHandler) SetBufSize(size int) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.BufSize = size
	h.Buffer <- nil
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
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Buffer <- &rd
}

func (h *BaseHandler) upgrade_buffer() {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	close(h.Buffer)
	buffer := make(chan *Record, h.BufSize)
	for {
		remain, ok := <-h.Buffer
		if remain == nil || !ok {
			break
		}
		buffer <- remain
	}
	h.Buffer = buffer
}

func (h *BaseHandler) set_state(state bool, err error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.State = state
	h.LastError = err
}

func (h *BaseHandler) get_state() (bool, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.State, h.LastError
}

func (h *BaseHandler) handle_record(rd *Record, buf *bytes.Buffer) {
	defer func() {
		if err := recover(); err != nil {
			h.set_state(false, err.(error))
		}
	}()
	if state, _ := h.get_state(); !state {
		return
	}
	if h.Filter != nil && h.Filter(rd) {
		return
	}
	rd.TimeString = rd.Time.Format(h.TimeLayout)
	buf.Reset()
	if err := h.Tmpl.Execute(buf, rd); err != nil {
		h.set_state(false, err)
		return
	}
	if h.Before != nil {
		h.Before(rd, buf)
	}
	n, err := io.Copy(h.Writer, buf)
	if err != nil {
		h.set_state(false, err)
	}
	if h.After != nil {
		h.After(rd, int64(n))
	}
}

func (h *BaseHandler) WriteRecord() {
	rd := &Record{}
	buf := bytes.NewBuffer(nil)
	for {
		rd = <-h.Buffer
		if rd == nil {
			h.upgrade_buffer()
			go h.WriteRecord()
			break
		}
		h.handle_record(rd, buf)
	}
}

func (h *BaseHandler) notify() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP)
	for {
		<-c
		h.set_state(true, nil)
	}

}
