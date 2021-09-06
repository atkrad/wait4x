package waiter

import (
	"context"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
)

// Check represents the checker's check method.
type Check func() bool

// Wait waits for end up of check execution.
func Wait(check Check, timeout time.Duration, interval time.Duration, invertCheck bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	checking := check
	if invertCheck == true {
		checking = func() bool { return !check() }
	}

	for !checking() {
		select {
		case <-ctx.Done():
			return errors.NewTimedOutError()
		case <-time.After(interval):
		}
	}

	return nil
}
