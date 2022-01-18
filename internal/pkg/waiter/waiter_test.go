// Copyright 2021 Mohammad Abdolirad
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
	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestWaitSuccessful(t *testing.T) {
	alwaysTrue := func() bool {
		time.Sleep(3 * time.Second)
		return true
	}
	err := Wait(alwaysTrue, WithInterval(time.Second))

	assert.Nil(t, err)
}

func TestWaitTimedOut(t *testing.T) {
	alwaysFalse := func() bool { return false }
	err := Wait(alwaysFalse, WithTimeout(time.Second))

	assert.Equal(t, errors.NewTimedOutError(), err)
}

func TestWaitInvertCheck(t *testing.T) {
	alwaysTrue := func() bool { return true }
	alwaysFalse := func() bool { return false }

	err := Wait(alwaysTrue, WithTimeout(time.Second), WithInvertCheck(true))
	assert.Equal(t, errors.NewTimedOutError(), err)

	err = Wait(alwaysFalse, WithTimeout(time.Second), WithInvertCheck(true))
	assert.Nil(t, err)
}
