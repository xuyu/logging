package logging

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	LAYOUT = "2006-01-02 15:04:05"
)

type BaseLogger struct {
	mutex   sync.Mutex
	prefix  string
	out     io.WriteCloser
	lv      level
	layout  string
	handler func(log string)
}

func NewBaseLogger(out io.WriteCloser, prefix string, lv level, layout string) *BaseLogger {
	return &BaseLogger{out: out, prefix: prefix, lv: lv, layout: layout}
}

func (l *BaseLogger) SetPrefix(prefix string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.prefix = prefix
}

func (l *BaseLogger) SetLevel(lv level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.lv = lv
}

func (l *BaseLogger) SetWriter(out io.WriteCloser) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.out = out
}

func (l *BaseLogger) SetLayout(layout string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.layout = layout
}

func (l *BaseLogger) timestamp() string {
	return time.Now().Format(l.layout)
}

func (l *BaseLogger) head(lv string) string {
	a := []string{"[", l.timestamp(), "]"}
	if l.prefix != "" {
		a = append(a, " ", l.prefix)
	}
	a = append(a, " ", lv)
	return strings.Join(a, "")
}

func (l *BaseLogger) log(lv string, format string, v ...interface{}) {
	log := l.head(lv) + " " + fmt.Sprintf(strings.TrimRight(format, "\n")+"\n", v...)
	if l.handler != nil {
		l.handler(log)
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	fmt.Fprint(l.out, log)
}

func (l *BaseLogger) Debug(format string, v ...interface{}) {
	l.do(DEBUG, format, v...)
}

func (l *BaseLogger) Info(format string, v ...interface{}) {
	l.do(INFO, format, v...)
}

func (l *BaseLogger) Warning(format string, v ...interface{}) {
	l.do(WARNING, format, v...)
}

func (l *BaseLogger) Error(format string, v ...interface{}) {
	l.do(ERROR, format, v...)
}

func (l *BaseLogger) do(lv level, format string, v ...interface{}) {
	if l.lv <= lv {
		l.log(lv.String(), format, v...)
	}
}
