package logging

import (
	"bytes"
	"io"
	"os"
	"text/template"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
	DefaultFormat     = "[{{.TimeString}}] {{.Level}} {{.Message}}\n"
)

var (
	FileCreateFlag             = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	FileCreatePerm os.FileMode = 0640
	DefaultBufSize             = 1024
)

type Handler struct {
	Async      bool
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

func NewHandler(out io.Writer, level LogLevel, layout, format string) (*Handler, error) {
	h := &Handler{
		Async:      true,
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

func (h *Handler) SetLevel(level LogLevel) {
	h.Level = level
}

func (h *Handler) SetLevelString(s string) {
	h.SetLevel(StringToLogLevel(s))
}

func (h *Handler) SetLevelRange(minLevel, maxLevel LogLevel) {
	h.LRange = &LevelRange{minLevel, maxLevel}
}

func (h *Handler) SetLevelRangeString(smin, smax string) {
	h.SetLevelRange(StringToLogLevel(smin), StringToLogLevel(smax))
}

func (h *Handler) SetTimeLayout(layout string) {
	h.TimeLayout = layout
}

func (h *Handler) SetFormat(format string) error {
	tmpl, err := template.New("tmpl").Parse(format)
	if err != nil {
		return err
	}
	h.Tmpl = tmpl
	return nil
}

func (h *Handler) SetFilter(f func(*Record) bool) {
	h.Filter = f
}

func (h *Handler) Emit(rd Record) {
	if h.LRange != nil {
		if !h.LRange.Contain(rd.Level) {
			return
		}
	} else if h.Level > rd.Level {
		return
	}
	if h.Async {
		h.Buffer <- &rd
	} else {
		h.handleRecord(&rd, bytes.NewBuffer(nil))
	}
}

func (h *Handler) handleRecord(rd *Record, buf *bytes.Buffer) {
	if h.Writer == nil {
		return
	}
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
		return
	}
	if h.After != nil {
		h.After(rd, int64(n))
	}
}

func (h *Handler) WriteRecord() {
	rd := &Record{}
	buf := bytes.NewBuffer(nil)
	for {
		rd = <-h.Buffer
		if rd != nil {
			h.handleRecord(rd, buf)
		}
	}
}
