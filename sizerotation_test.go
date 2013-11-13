package logging

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestSizeRotationHandler(t *testing.T) {
	h, err := NewSizeRotationHandler(path.Join(os.TempDir(), "sr.log"), 64, 3)
	if err != nil {
		t.Fatal(err)
	}
	AddHandler("sr", h)
	Debug("%d, %s", 1, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	Info("%d, %s", 2, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	Warning("%d, %s", 3, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	Error("%d, %s", 4, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	time.Sleep(100 * time.Millisecond)
}
