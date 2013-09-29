package logging

import (
	"testing"
)

func TestFileHandler(t *testing.T) {
	f, err := NewFileHandler("/tmp/file.log")
	if err != nil {
		t.Fatal(err)
	}
	f.SetLevel(INFO)
	DefaultLogger.AddHandler("file", f)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 2, "OK")
	Warning("%d, %s", 3, "OK")
	Error("%d, %s", 4, "OK")
}
