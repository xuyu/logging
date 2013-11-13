package logging

import (
	"testing"
)

func TestStdoutHandler(t *testing.T) {
	EnableStdout()
	Debug("%d, %s", 1, "OK")
}
