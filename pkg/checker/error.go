package checker

type Error struct {
	msg    string
	level  string
	fields map[string]interface{}
	err    error
}

func NewError(msg string, lvl string) *Error {
	return &Error{
		msg:    msg,
		level:  lvl,
		fields: make(map[string]interface{}),
	}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) WithWrap(err error) error {
	e.err = err

	return e
}

func (e *Error) WithFields(fields map[string]interface{}) error {
	e.fields = fields

	return e
}

func (e *Error) WithField(name string, value interface{}) error {
	e.fields[name] = value

	return e
}
