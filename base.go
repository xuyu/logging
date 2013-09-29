package logging

import (
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
	mutex  sync.Mutex
	out    io.WriteCloser
	lv     level
	layout string
	tmpl   *template.Template
	c      chan *Formatter
	predo  func()
}

func NewBaseHandler(out io.WriteCloser, lv level, layout string, format string) *BaseHandler {
	b := &BaseHandler{
		out:    out,
		lv:     lv,
		layout: layout,
	}
	b.SetFormat(format)
	b.c = make(chan *Formatter, DefaultBufSize)
	go b.work()
	return b
}

func (b *BaseHandler) SetLevel(lv level) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.lv = lv
}

func (b *BaseHandler) GetLevel() level {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.lv
}

func (b *BaseHandler) SetTimeLayout(layout string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.layout = layout
}

func (b *BaseHandler) GetTimeLayout() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.layout
}

func (b *BaseHandler) SetFormat(format string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	tmpl, err := template.New("tmpl").Parse(format)
	if err != nil {
		panic(err)
	}
	b.tmpl = tmpl
}

func (b *BaseHandler) Emit(lv level, f string, values ...interface{}) {
	if b.GetLevel() > lv {
		return
	}
	fm := &Formatter{
		TimeString: time.Now().Format(b.GetTimeLayout()),
		Level:      lv,
		Message:    fmt.Sprintf(f, values...),
	}
	b.c <- fm
}

func (b *BaseHandler) work() {
	var fm *Formatter
	for {
		fm = <-b.c
		if b.predo != nil {
			b.predo()
		}
		b.tmpl.Execute(b.out, fm)
	}
}
