package logging

import (
	"testing"
	"time"
)

func TestTimeRotationHandler(t *testing.T) {
	r, err := NewTimeRotationHandler("/tmp/tr.log", "060102-15:04:05")
	if err != nil {
		t.Fatal(err)
	}
	r.SetLevel(INFO)
	r.Panic(true)
	AddHandler("rotation", r)
	for i := 0; i < 3; i++ {
		Debug("%d, %s", 1, "OK")
		time.Sleep(time.Second)
		Info("%d, %s", 2, "OK")
		time.Sleep(time.Second)
		Warning("%d, %s", 3, "OK")
		time.Sleep(time.Second)
		Error("%d, %s", 4, "OK")
	}
}
