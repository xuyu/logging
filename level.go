package logging

type level int

const (
	DEBUG   level = 1
	INFO    level = 2
	WARNING level = 3
	ERROR   level = 4
	DISABLE level = 5
)

func (l *level) String() string {
	switch *l {
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
