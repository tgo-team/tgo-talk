package tgo

// Level type
type LogLevel uint32

var AllLevels = []LogLevel{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

const (
	PanicLevel LogLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)
type Log interface {
	Info(format string,a ...interface{})
	Error(format string,a ...interface{})
	Debug(format string,a ...interface{})
	Warn(format string,a ...interface{})
	Fatal(format string,a ...interface{})
}