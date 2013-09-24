package logging

import (
	"os"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	SetLevel(INFO)
	SetPrefix("Prefix")
	SetWriter(os.Stderr)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 2, "OK")
	Warning("%d, %s", 3, "OK")
	Error("%d, %s", 4, "OK")
}
