package logging

import (
	"os"
	"path"
	"testing"
)

func TestSingleFileHandler(t *testing.T) {
	f, err := NewSingleFileHandler(path.Join(os.TempDir(), "sf.log"))
	if err != nil {
		t.Fatal(err)
	}
	AddHandler("file", f)
	Debug("%d, %s", 1, "OK")
}
