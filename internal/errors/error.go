package errors

type CommandError struct {
	Message  string
	ExitCode int
}

const ExitError = 1
const ExitTimedOut = 124

const TimedOutErrorMessage = "Timed Out"

func (e *CommandError) Error() string {
	return e.Message
}

func NewCommandError(msg string) error {
	return &CommandError{
		Message:  msg,
		ExitCode: ExitError,
	}
}

func NewTimedOutError() error {
	return &CommandError{
		Message:  TimedOutErrorMessage,
		ExitCode: ExitTimedOut,
	}
}
