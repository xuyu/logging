package logging

import (
	"bytes"
	"errors"
	"os"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

var h *BaseHandler
var b *bytes.Buffer

func init() {
	b = bytes.NewBuffer(nil)
	var err error
	h, err = NewBaseHandler(b, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		panic(err)
	}
	DisableStdout()
}

func TestSetBufSize(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	Debug("%d, %s", 1, "OK")
	h.SetBufSize(128)
	if h.BufSize != 128 {
		t.Fail()
	}
	Debug("%d, %s", 2, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 68 {
		t.Fail()
	}
}

func TestSetLevel(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	Debug("%d, %s", 1, "OK")
	h.SetLevel(INFO)
	Debug("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 34 {
		t.Fail()
	}
}

func TestSetLevelString(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	Debug("%d, %s", 1, "OK")
	h.SetLevelString("info")
	Debug("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 34 {
		t.Fail()
	}
}

func TestSetLevelRange(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	h.SetLevelRange(INFO, WARNING)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 34*4-2 {
		t.Fail()
	}
	h.LRange = nil
}

func TestSetLevelRangeString(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	h.SetLevelRangeString("INFO", "WARNING")
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 34*4-2 {
		t.Fail()
	}
	h.LRange = nil
}

func TestSetTimeLayout(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	h.SetTimeLayout("2006/01/02-15:04:05")
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 34 {
		t.Fail()
	}
	h.SetTimeLayout(DefaultTimeLayout)
}

func TestSetFilter(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	h.SetFilter(func(rd *Record) bool {
		return strings.Contains(rd.Message, "OK")
	})
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if b.Len() != 0 {
		t.Fail()
	}
	h.SetFilter(nil)
}

func TestPanicError(t *testing.T) {
	b.Reset()
	AddHandler("b", h)
	h.SetLevel(DEBUG)
	h.SetFilter(func(*Record) bool {
		panic(errors.New("nothing"))
	})
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
	if state, _ := h.get_state(); state {
		t.Fail()
	}
	h.set_state(true, nil)
	h.SetFilter(nil)
}

func TestNotify(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}
	h.set_state(false, nil)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGHUP)
	time.Sleep(100 * time.Millisecond)
	if state, _ := h.get_state(); !state {
		t.Fail()
	}
	h.set_state(true, nil)
}
