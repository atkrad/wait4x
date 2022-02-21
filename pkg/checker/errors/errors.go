// Copyright 2022 Mohammad Abdolirad
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

package errors

// Level represents error level
type Level int

// Option configures an Error.
type Option func(s *Error)

const (
	// InfoLevel info error level
	InfoLevel Level = iota
	DebugLevel
	TraceLevel
)

// Error represents checker error
type Error struct {
	msg    string
	err    error
	level  Level
	fields []interface{}
}

// New creates the Error
func New(msg string, level Level, opts ...Option) error {
	e := &Error{
		msg:    msg,
		level:  level,
		fields: make([]interface{}, 0),
	}

	// apply the list of options to Error
	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Wrap wraps an error in the Error
func Wrap(err error, level Level, opts ...Option) *Error {
	e := &Error{
		err:    err,
		level:  level,
		fields: make([]interface{}, 0),
	}

	// apply the list of options to Error
	for _, opt := range opts {
		opt(e)
	}

	return e
}

// WithFields configures the error fields
func WithFields(fields ...interface{}) Option {
	return func(e *Error) {
		e.fields = fields
	}
}

// Level returns the error level
func (e *Error) Level() Level {
	return e.level
}

// Fields returns the error fields
func (e *Error) Fields() []interface{} {
	return e.fields
}

func (e *Error) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return msg
}

func (e *Error) Unwrap() error {
	return e.err
}
