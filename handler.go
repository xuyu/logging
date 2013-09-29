package logging

type Handler interface {
	SetLevel(LogLevel)
	GetLevel() LogLevel
	SetTimeLayout(string)
	GetTimeLayout() string
	SetFormat(string) error
	Emit(LogLevel, string, ...interface{})
}

type Record struct {
	TimeString string
	Level      LogLevel
	Message    string
}
