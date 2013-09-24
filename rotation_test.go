package logging

import (
	"testing"
)

func TestRotationLogger(t *testing.T) {
	l, err := NewRotationLogger("/tmp/rotation.log", "060102-15")
	if err != nil {
		t.Fatal(err)
	}
	SetDefaultLogger(l)
	SetLevel(INFO)
	SetPrefix("Prefix")
	Debug("%d, %s\n", 1, "OK")
	Info("%d, %s\n", 2, "OK")
	Warning("%d, %s\n", 3, "OK")
	Error("%d, %s\n", 4, "OK")
}
