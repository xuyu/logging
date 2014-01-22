package logging

import (
	"bytes"
	"io"
	"text/template"
	"time"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
	DefaultFormat     = "[{{.TimeString}}] {{.Level}} {{.Message}}\n"
	FormatNoTime      = "{{.Level}} {{.Message}}\n"
	FormatNoLevel     = "[{{.TimeString}}] {{.Message}}\n"
	FormatOnlyMessage = "{{.Message}}\n"
)

var (
	DefaultBufSize = 1024
)

type Handler interface {
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
	Writer     io.Writer
	Level      LogLevel
	LRange     *LevelRange
	TimeLayout string
	Tmpl       *template.Template
	Buffer     chan *Record
	Filter     func(*Record) bool
	Before     func(*Record, io.ReadWriter)
	After      func(*Record, int64)
}

func NewBaseHandler(out io.Writer, level LogLevel, layout, format string) (*BaseHandler, error) {
	h := &BaseHandler{
		Writer:     out,
		Level:      level,
		TimeLayout: layout,
	}
	if err := h.SetFormat(format); err != nil {
		return nil, err
	}
	h.Buffer = make(chan *Record, DefaultBufSize)
	go h.WriteRecord()
	return h, nil
}

func (h *BaseHandler) SetLevel(level LogLevel) {
	h.Level = level
}

func (h *BaseHandler) SetLevelString(s string) {
	h.SetLevel(StringToLogLevel(s))
}

func (h *BaseHandler) SetLevelRange(min_level, max_level LogLevel) {
	h.LRange = &LevelRange{min_level, max_level}
}

func (h *BaseHandler) SetLevelRangeString(smin, smax string) {
	h.SetLevelRange(StringToLogLevel(smin), StringToLogLevel(smax))
}

func (h *BaseHandler) SetTimeLayout(layout string) {
	h.TimeLayout = layout
}

func (h *BaseHandler) SetFormat(format string) error {
	tmpl, err := template.New("tmpl").Parse(format)
	if err != nil {
		return err
	}
	h.Tmpl = tmpl
	return nil
}

func (h *BaseHandler) SetFilter(f func(*Record) bool) {
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

func (h *BaseHandler) handle_record(rd *Record, buf *bytes.Buffer) {
	if h.Filter != nil && h.Filter(rd) {
		return
	}
	rd.TimeString = rd.Time.Format(h.TimeLayout)
	buf.Reset()
	if err := h.Tmpl.Execute(buf, rd); err != nil {
		return
	}
	if h.Before != nil {
		h.Before(rd, buf)
	}
	n, err := io.Copy(h.Writer, buf)
	if err != nil {
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
		if rd != nil {
			h.handle_record(rd, buf)
		}
	}
}
