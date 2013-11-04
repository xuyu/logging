package logging

import (
	"strings"
)

type LogLevel uint8

const (
	DEBUG   LogLevel = 1
	INFO    LogLevel = 2
	WARNING LogLevel = 3
	ERROR   LogLevel = 4
	DISABLE LogLevel = 255
)

func StringToLogLevel(s string) LogLevel {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	default:
		return DISABLE
	}
}

func (level *LogLevel) String() string {
	switch *level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}

type LevelRange struct {
	MinLevel LogLevel
	MaxLevel LogLevel
}

func (lr *LevelRange) Contain(level LogLevel) bool {
	return level >= lr.MinLevel && level <= lr.MaxLevel
}
