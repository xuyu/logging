package logging

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var h *Handler
var b *bytes.Buffer

func init() {
	b = bytes.NewBuffer(nil)
	var err error
	h, err = NewHandler(b, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		panic(err)
	}
	DisableStdout()
	AddHandler("b", h)
}

func TestSetLevel(t *testing.T) {
	b.Reset()
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

func TestAsyncHandler(t *testing.T) {
	b.Reset()
	h.Async = false
	h.SetLevel(DEBUG)
	Error("%d, %s", 1, "OK")
	if b.Len() != 34 {
		t.Fail()
	}
}
