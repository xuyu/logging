package logging

import (
	"testing"
)

func TestFileLogger(t *testing.T) {
	l, err := NewFileLogger("/tmp/file.log")
	if err != nil {
		t.Fatal(err)
	}
	SetDefaultLogger(l)
	SetLevel(INFO)
	SetPrefix("Prefix")
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 2, "OK")
	Warning("%d, %s", 3, "OK")
	Error("%d, %s", 4, "OK")
}
