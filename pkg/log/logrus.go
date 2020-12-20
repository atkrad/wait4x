package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logrus is the implementation for a Logger using Logrus.
type Logrus struct {
	logger *logrus.Logger
}

// NewLogrus creates Logrus logger
func NewLogrus(level string, output io.Writer) (Logger, error) {
	log := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	log.SetLevel(lvl)
	log.SetOutput(output)

	l := &Logrus{
		logger: log,
	}

	return l, nil
}

// Info logging a new message with info level.
func (l *Logrus) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logging a new message with info level and custom format.
func (l *Logrus) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// InfoWithFields logging a new message with info level and extra fields.
func (l *Logrus) InfoWithFields(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}

// Debug logging a new message with debug level.
func (l *Logrus) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf logging a new message with debug level and custom format.
func (l *Logrus) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Fatal logging a new message with fatal level.
func (l *Logrus) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
