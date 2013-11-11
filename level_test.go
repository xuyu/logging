package logging

import (
	"testing"
)

func TestLevelToString(t *testing.T) {
	data := map[LogLevel]string{
		DEBUG:         "DEBUG",
		INFO:          "INFO",
		WARNING:       "WARNING",
		ERROR:         "ERROR",
		LogLevel(250): "DISABLE",
	}
	for level, str := range data {
		if level.String() != str {
			t.Error(level.String())
		}
	}
}

func TestStringToLevel(t *testing.T) {
	data := map[string]LogLevel{
		"DEBUG":   DEBUG,
		"INFO":    INFO,
		"WARNING": WARNING,
		"ERROR":   ERROR,
		"":        DISABLE,
	}
	for str, level := range data {
		if StringToLogLevel(str) != level {
			t.Error(str)
		}
	}
}

func TestLevelRange(t *testing.T) {
	lr := LevelRange{INFO, WARNING}
	if lr.Contain(ERROR) {
		t.Error("TestLevelRange Fail")
	}
	if lr.Contain(DEBUG) {
		t.Error("TestLevelRange Fail")
	}
	if !lr.Contain(INFO) {
		t.Error("TestLevelRange Fail")
	}
}
