package logging

import (
	"os"
	"testing"
)

func TestLogging(t *testing.T) {
	SetLevel(INFO)
	SetPrefix("Prefix")
	SetWriter(os.Stderr)
	Debug("%d, %s\n", 1, "OK")
	Info("%d, %s\n", 1, "OK")
	Warning("%d, %s\n", 1, "OK")
	Error("%d, %s\n", 1, "OK")
}

func TestRotationLogger(t *testing.T) {
	l, err := NewRotationLogger("test.log", "/tmp", "060102-15")
	if err != nil {
		t.Fatal(err.Error())
	}
	SetDefaultLogger(l)
	SetLevel(INFO)
	SetPrefix("Prefix")
	Debug("%d, %s\n", 1, "OK")
	Info("%d, %s\n", 1, "OK")
	Warning("%d, %s\n", 1, "OK")
	Error("%d, %s\n", 1, "OK")
}
