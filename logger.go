package logging

type Logger struct {
	handlers map[string]Handler
}

func NewLogger() *Logger {
	return &Logger{
		handlers: make(map[string]Handler),
	}
}

var DefaultLogger *Logger = NewLogger()

func (l *Logger) AddHandler(name string, h Handler) {
	l.handlers[name] = h
}

func (l *Logger) DelHandler(name string) {
	delete(l.handlers, name)
}

func (l *Logger) GetHandler(name string) Handler {
	return l.handlers[name]
}

func (l *Logger) log(level LogLevel, format string, values ...interface{}) {
	for _, h := range l.handlers {
		if h.GetLevel() > level {
			continue
		}
		h.Emit(level, format, values...)
	}
}

func AddHandler(name string, h Handler) {
	DefaultLogger.AddHandler(name, h)
}

func DelHandler(name string) {
	DefaultLogger.DelHandler(name)
}

func GetHandler(name string) Handler {
	return DefaultLogger.GetHandler(name)
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

func (l *LogLevel) String() string {
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
