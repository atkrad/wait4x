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

package errors

// CommandError represents command errors
type CommandError struct {
	Message  string
	ExitCode int
}

// ExitError represents general error exit code
const ExitError = 1

// ExitTimedOut represents timed out error exit code
const ExitTimedOut = 124

// TimedOutErrorMessage represents timed out error message
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
