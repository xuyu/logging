package logging

import (
	"os"
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
	AddHandler("b", h)
}

func TestSetBufSize(t *testing.T) {
	DisableStdout()
	Debug("%d, %s", 1, "OK")
	h.SetBufSize(1)
	if h.BufSize != 1 {
		t.Fail()
	}
	Debug("%d, %s", 2, "OK")
	time.Sleep(100 * time.Millisecond)
}
