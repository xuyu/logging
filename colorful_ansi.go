// +build !windows

package logging

import (
	"os"
)

func resetColorful() {
	os.Stdout.WriteString("x1b[0m")
}

func changeColor(c color) {
	switch c {
	case red:
		os.Stdout.WriteString("\x1b[31;1m")
	case yellow:
		os.Stdout.WriteString("\x1b[33;1m")
	case green:
		os.Stdout.WriteString("\x1b[32;1m")
	}
}
