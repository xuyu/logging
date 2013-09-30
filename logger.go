package logging

type Logger struct {
	Handlers map[string]Handler
}

func NewLogger() *Logger {
	return &Logger{
		Handlers: make(map[string]Handler),
	}
}

var DefaultLogger *Logger = NewLogger()

func (l *Logger) AddHandler(name string, h Handler) {
	l.Handlers[name] = h
}

func (l *Logger) log(level LogLevel, format string, values ...interface{}) {
	for _, h := range l.Handlers {
		if h.GetLevel() > level {
			continue
		}
		h.Emit(level, format, values...)
	}
}

func AddHandler(name string, h Handler) {
	DefaultLogger.AddHandler(name, h)
}

func Debug(format string, values ...interface{}) {
	DefaultLogger.log(DEBUG, format, values...)
}

func Info(format string, values ...interface{}) {
	DefaultLogger.log(INFO, format, values...)
}

func Warning(format string, values ...interface{}) {
	DefaultLogger.log(WARNING, format, values...)
}

func Error(format string, values ...interface{}) {
	DefaultLogger.log(ERROR, format, values...)
}

type LogLevel uint8

const (
	DEBUG   LogLevel = 1
	INFO    LogLevel = 2
	WARNING LogLevel = 3
	ERROR   LogLevel = 4
	DISABLE LogLevel = 255
)

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
