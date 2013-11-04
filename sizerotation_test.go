package logging

import (
	"testing"
	"time"
)

func TestSizeRotationHandler(t *testing.T) {
	h, err := NewSizeRotationHandler("/tmp/sr.log", 1024, 5)
	if err != nil {
		t.Fatal(err)
	}
	h.SetLevel(INFO)
	h.Panic(true)
	AddHandler("sr", h)
	for i := 0; i < 100; i++ {
		Debug("%d, %s", i, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		Info("%d, %s", i, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		Warning("%d, %s", i, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		Error("%d, %s", i, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		time.Sleep(100 * time.Millisecond)
	}
}
