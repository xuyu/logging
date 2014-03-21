package logging

import (
	"bytes"
	"io"
	"os"
	"sync"
	"text/template"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
	DefaultFormat     = "[{{.TimeString}}] {{.Level}} {{.Message}}\n"
)

var (
	FileCreateFlag             = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	FileCreatePerm os.FileMode = 0640
)

type Handler struct {
	mutex      sync.Mutex
	buffer     *bytes.Buffer
	writer     io.Writer
	level      LogLevel
	lRange     *LevelRange
	timeLayout string
	tmpl       *template.Template
	filter     func(*Record) bool

	Before func(*Record, io.ReadWriter)
	After  func(*Record, int64)
}

func NewHandler(out io.Writer, level LogLevel, layout, format string) (*Handler, error) {
	h := &Handler{
		buffer:     bytes.NewBuffer(nil),
		writer:     out,
		level:      level,
		timeLayout: layout,
	}
	if err := h.SetFormat(format); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Handler) SetLevel(level LogLevel) {
	h.level = level
}

func (h *Handler) SetLevelString(s string) {
	h.SetLevel(StringToLogLevel(s))
}

func (h *Handler) SetLevelRange(minLevel, maxLevel LogLevel) {
	h.lRange = &LevelRange{minLevel, maxLevel}
}

func (h *Handler) SetLevelRangeString(smin, smax string) {
	h.SetLevelRange(StringToLogLevel(smin), StringToLogLevel(smax))
}

func (h *Handler) SetTimeLayout(layout string) {
	h.timeLayout = layout
}

func (h *Handler) SetFormat(format string) error {
	tmpl, err := template.New("tmpl").Parse(format)
	if err != nil {
		return err
	}
	h.tmpl = tmpl
	return nil
}

func (h *Handler) SetFilter(f func(*Record) bool) {
	h.filter = f
}

func (h *Handler) Emit(rd Record) {
	if h.lRange != nil {
		if !h.lRange.Contain(rd.Level) {
			return
		}
	} else if h.level > rd.Level {
		return
	}
	h.mutex.Lock()
	h.buffer.Reset()
	h.handleRecord(&rd, h.buffer)
	h.mutex.Unlock()
}

func (h *Handler) handleRecord(rd *Record, buf *bytes.Buffer) {
	if h.writer == nil {
		return
	}
	if h.filter != nil && h.filter(rd) {
		return
	}
	rd.TimeString = rd.Time.Format(h.timeLayout)
	if err := h.tmpl.Execute(buf, rd); err != nil {
		return
	}
	if h.Before != nil {
		h.Before(rd, buf)
	}
	n, err := io.Copy(h.writer, buf)
	if err != nil {
		return
	}
	if h.After != nil {
		h.After(rd, int64(n))
	}
}
