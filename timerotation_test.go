package logging

import (
	"testing"
)

func TestTimeRotationHandler(t *testing.T) {
	r, err := NewTimeRotationHandler("/tmp/tr.log", "060102-15")
	if err != nil {
		t.Fatal(err)
	}
	r.SetLevel(INFO)
	AddHandler("rotation", r)
	Debug("%d, %s", 1, "OK")
	Info("%d, %s", 2, "OK")
	Warning("%d, %s", 3, "OK")
	Error("%d, %s", 4, "OK")
}
