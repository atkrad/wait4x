package checker

import "fmt"

// ExpectedError defines the expectation error
type ExpectedError struct {
	msg     string
	cause   error
	details []any
}

// NewExpectedError creates the ExpectedError
func NewExpectedError(msg string, cause error, details ...any) error {
	ee := &ExpectedError{
		msg:     msg,
		cause:   cause,
		details: details,
	}

	return ee
}

// Details returns the error details
func (ee *ExpectedError) Details() []any {
	return ee.details
}

func (ee *ExpectedError) Unwrap() error {
	return ee.cause
}

func (ee *ExpectedError) Error() string {
	if ee.cause != nil {
		return fmt.Sprintf("%s, caused by: %s", ee.msg, ee.cause.Error())
	}

	return ee.msg
}
