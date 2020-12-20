package log

// Logger is the interface that wraps the basic logging methods.
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoWithFields(msg string, fields map[string]interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Fatal(args ...interface{})
}
