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
	alwaysTrue := func() bool { return true }
	err := Wait(alwaysTrue, time.Second*10, time.Second, false)

	assert.Nil(t, err)
}

func TestWaitTimedOut(t *testing.T) {
	alwaysFalse := func() bool { return false }
	err := Wait(alwaysFalse, time.Second, time.Second, false)

	assert.Equal(t, errors.NewTimedOutError(), err)
}

func TestWaitInvertCheck(t *testing.T) {
	alwaysTrue := func() bool { return true }
	alwaysFalse := func() bool { return false }

	err := Wait(alwaysTrue, time.Second, time.Second, true)
	assert.Equal(t, errors.NewTimedOutError(), err)

	err = Wait(alwaysFalse, time.Second, time.Second, true)
	assert.Nil(t, err)
}
