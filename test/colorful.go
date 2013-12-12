package main

import (
	"time"

	"git.code4.in/logging"
)

func main() {
	logging.EnableStdout()
	logging.EnableColorful()

	println("Colorful logging:")
	logging.Debug("%d, %s", 1, "OK")
	logging.Info("%d, %s", 1, "OK")
	logging.Warning("%d, %s", 1, "OK")
	logging.Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)

	logging.DisableColorful()
	println("\nNormal logging:")
	logging.Debug("%d, %s", 1, "OK")
	logging.Info("%d, %s", 1, "OK")
	logging.Warning("%d, %s", 1, "OK")
	logging.Error("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}
