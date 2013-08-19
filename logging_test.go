package logging

import (
	"testing"
)

func TestLogging(t *testing.T) {
	Debug("%d, %s\n", 1, "OK")
	Info("%d, %s\n", 1, "OK")
	Warning("%d, %s\n", 1, "OK")
	Error("%d, %s\n", 1, "OK")
}
