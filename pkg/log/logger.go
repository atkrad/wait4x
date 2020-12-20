package log

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoWithFields(msg string, fields map[string]interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Fatal(args ...interface{})
}
