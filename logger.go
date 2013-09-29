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

func (l *Logger) log(lv level, format string, values ...interface{}) {
	for _, h := range l.handlers {
		if h.GetLevel() > lv {
			return
		}
		h.Emit(lv, format, values...)
	}
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

type Handler interface {
	SetLevel(level)
	GetLevel() level
	SetTimeLayout(string)
	GetTimeLayout() string
	SetFormat(string)
	Emit(level, string, ...interface{})
}

type Formatter struct {
	TimeString string
	Level      level
	Message    string
}

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
