package errors

type CommandError struct {
	Message  string
	ExitCode int
}

const EXIT_ERROR = 1
const EXIT_TIMEDOUT = 124

const TIMED_OUT_ERROR_MESSAGE = "Timed Out"

func (e *CommandError) Error() string {
	return e.Message
}

func NewCommandError(msg string) error {
	return &CommandError{
		Message:  msg,
		ExitCode: EXIT_ERROR,
	}
}

func NewTimedOutError() error {
	return &CommandError{
		Message:  TIMED_OUT_ERROR_MESSAGE,
		ExitCode: EXIT_TIMEDOUT,
	}
}
