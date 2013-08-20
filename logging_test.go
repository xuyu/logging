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
