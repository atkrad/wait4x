// Copyright 2021 The Wait4X Authors
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

package waiter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tonglil/buflogr"
	"os"
	"testing"
	"time"
	"wait4x.dev/v2/checker"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestWaitSuccessful(t *testing.T) {
	mockChecker := new(checker.MockChecker)
	mockChecker.On("Check", mock.Anything).Return(nil).
		On("Identity").Return("ID", nil)

	err := Wait(mockChecker, WithInterval(time.Second))

	assert.Nil(t, err)
	mockChecker.AssertExpectations(t)
}

func TestWaitTimedOut(t *testing.T) {
	mockChecker := new(checker.MockChecker)
	mockChecker.On("Check", mock.Anything).Return(fmt.Errorf("error")).
		On("Identity").Return("ID", nil)

	err := Wait(mockChecker, WithTimeout(time.Second))

	assert.Equal(t, context.DeadlineExceeded, err)
	mockChecker.AssertExpectations(t)
}

func TestWaitInvalidIdentity(t *testing.T) {
	invalidIdentityError := errors.New("invalid identity")

	mockChecker := new(checker.MockChecker)
	mockChecker.On("Identity").Return(mock.Anything, invalidIdentityError)

	err := Wait(mockChecker)

	assert.Equal(t, invalidIdentityError, err)
	mockChecker.AssertExpectations(t)
}

func TestWaitLogger(t *testing.T) {
	mockChecker := new(checker.MockChecker)
	mockChecker.On("Check", mock.Anything).
		Return(fmt.Errorf("error message")).
		On("Identity").Return("ID", nil)

	var buf bytes.Buffer
	var log logr.Logger = buflogr.NewWithBuffer(&buf)
	// TODO: Change the "WaitWithContext" to "Wait" when we want release v3.0.0
	err := WaitWithContext(context.TODO(), mockChecker, WithLogger(log), WithTimeout(time.Second))

	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Contains(t, buf.String(), "INFO [MockChecker] Checking the ID ...")
	assert.Contains(t, buf.String(), "error message")
	mockChecker.AssertExpectations(t)
}

func TestWaitInvertCheck(t *testing.T) {
	alwaysTrue := new(checker.MockChecker)
	alwaysTrue.On("Check", mock.Anything).Return(nil).
		On("Identity").Return("ID", nil)

	err := Wait(alwaysTrue, WithTimeout(time.Second*3), WithInvertCheck(true))
	assert.Equal(t, context.DeadlineExceeded, err)
	alwaysTrue.AssertExpectations(t)

	alwaysFalse := new(checker.MockChecker)
	alwaysFalse.On("Check", mock.Anything).Return(fmt.Errorf("error")).
		On("Identity").Return("ID", nil)

	err = Wait(alwaysFalse, WithTimeout(time.Second), WithInvertCheck(true))
	assert.Nil(t, err)
	alwaysFalse.AssertExpectations(t)
}

func TestWaitParallelSuccessful(t *testing.T) {
	alwaysTrueFirst := new(checker.MockChecker)
	alwaysTrueFirst.On("Check", mock.Anything).Return(nil).
		On("Identity").Return("ID", nil)

	alwaysTrueSecond := new(checker.MockChecker)
	alwaysTrueSecond.On("Check", mock.Anything).Return(nil).
		On("Identity").Return("ID", nil)

	err := WaitParallel([]checker.Checker{alwaysTrueFirst, alwaysTrueSecond}, WithTimeout(time.Second*3))
	assert.Nil(t, err)
	alwaysTrueFirst.AssertExpectations(t)
	alwaysTrueSecond.AssertExpectations(t)
}

func TestWaitParallelFail(t *testing.T) {
	alwaysTrueFirst := new(checker.MockChecker)
	alwaysTrueFirst.On("Check", mock.Anything).Return(nil).
		On("Identity").Return("ID", nil)

	alwaysTrueSecond := new(checker.MockChecker)
	alwaysTrueSecond.On("Check", mock.Anything).Return(nil).
		On("Identity").Return("ID", nil)

	alwaysError := new(checker.MockChecker)
	alwaysError.On("Check", mock.Anything).Return(fmt.Errorf("error")).
		On("Identity").Return("ID", nil)

	err := WaitParallel([]checker.Checker{alwaysTrueFirst, alwaysTrueSecond, alwaysError}, WithTimeout(time.Second*3))
	assert.Equal(t, context.DeadlineExceeded, err)

	alwaysTrueFirst.AssertExpectations(t)
	alwaysTrueSecond.AssertExpectations(t)
	alwaysError.AssertExpectations(t)
}
