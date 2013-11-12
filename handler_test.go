package logging

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

var h *BaseHandler

func init() {
	var err error
	h, err = NewBaseHandler(os.Stdout, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		panic(err)
	}
}

func TestSetBufSize(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	Debug("%d, %s", 1, "OK")
	h.SetBufSize(1)
	if h.BufSize != 1 {
		t.Fail()
	}
	Debug("%d, %s", 2, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestSetLevel(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	Debug("%d, %s", 1, "OK")
	h.SetLevel(INFO)
	Debug("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestSetLevelString(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	Debug("%d, %s", 1, "OK")
	h.SetLevelString("info")
	Debug("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestSetLevelRange(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	h.SetLevelRange(INFO, WARNING)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestSetLevelRangeString(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	h.SetLevelRangeString("INFO", "WARNING")
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 1, "OK")
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestSetTimeLayout(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	h.SetTimeLayout("2006/01/02-15:04:05")
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestSetFilter(t *testing.T) {
	DisableStdout()
	AddHandler("b", h)
	h.SetFilter(func(rd *Record) bool {
		return strings.Contains(rd.Message, "OK")
	})
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}

func TestPanicError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fail()
		}
	}()
	DisableStdout()
	AddHandler("b", h)
	h.Panic(false)
	h.SetFilter(func(*Record) bool {
		panic(errors.New("nothing"))
		return true
	})
	Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}
