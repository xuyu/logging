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
	level      logLevel
	lRange     *levelRange
	timeLayout string
	tmpl       *template.Template
	filter     func(*Record) bool

	Before func(*Record, io.ReadWriter)
	After  func(*Record, int64)
}

func NewHandler(out io.Writer, level logLevel, layout, format string) (*Handler, error) {
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

func (h *Handler) Close() error {
	var err error
	h.mutex.Lock()
	closer, ok := h.writer.(io.Closer)
	if ok {
		err = closer.Close()
	}
	h.writer = nil
	h.mutex.Unlock()
	return err
}

func (h *Handler) SetLevel(level logLevel) {
	h.level = level
}

func (h *Handler) SetLevelString(s string) {
	h.SetLevel(stringToLogLevel(s))
}

func (h *Handler) SetLevelRange(minLevel, maxLevel logLevel) {
	h.lRange = &levelRange{minLevel, maxLevel}
}

func (h *Handler) SetLevelRangeString(smin, smax string) {
	h.SetLevelRange(stringToLogLevel(smin), stringToLogLevel(smax))
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
		if !h.lRange.contains(rd.Level) {
			return
		}
	} else if h.level > rd.Level {
		return
	}
	h.handleRecord(&rd)
}

func (h *Handler) handleRecord(rd *Record) {
	if h.filter != nil && h.filter(rd) {
		return
	}
	rd.TimeString = rd.Time.Format(h.timeLayout)
	h.mutex.Lock()
	if h.writer == nil {
		return
	}
	h.buffer.Reset()
	if err := h.tmpl.Execute(h.buffer, rd); err != nil {
		h.mutex.Unlock()
		return
	}
	if h.Before != nil {
		h.Before(rd, h.buffer)
	}
	n, err := io.Copy(h.writer, h.buffer)
	if err != nil {
		h.mutex.Unlock()
		return
	}
	h.mutex.Unlock()
	if h.After != nil {
		h.After(rd, int64(n))
	}
}
