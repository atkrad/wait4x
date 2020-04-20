package errors

// CommandError represents command errors
type CommandError struct {
	Message  string
	ExitCode int
}

// General error exit code.
const ExitError = 1

// Timed out error exit code
const ExitTimedOut = 124

// Timed out error message
const TimedOutErrorMessage = "Timed Out"

func (e *CommandError) Error() string {
	return e.Message
}

// NewCommandError creates the general error
func NewCommandError(msg string) error {
	return &CommandError{
		Message:  msg,
		ExitCode: ExitError,
	}
}

// NewTimedOutError creates the timed out error
func NewTimedOutError() error {
	return &CommandError{
		Message:  TimedOutErrorMessage,
		ExitCode: ExitTimedOut,
	}
}
