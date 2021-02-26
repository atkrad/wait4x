// Copyright 2020 Mohammad Abdolirad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
