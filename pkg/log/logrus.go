package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logrus is the implementation for a Logger using Logrus.
type Logrus struct {
	logger  *logrus.Logger
}

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

func (l *Logrus) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logrus) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logrus) InfoWithFields(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}

func (l *Logrus) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logrus) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logrus) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
