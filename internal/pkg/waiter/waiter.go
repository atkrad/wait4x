package waiter

import (
	"context"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
)

// Check represents the checker's check method.
type Check func() bool

// Option configures an options
type Option func(s *options)

// options represents waiter options
type options struct {
	timeout     time.Duration
	interval    time.Duration
	invertCheck bool
}

// WithTimeout configures a time limit for whole of checking
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithInterval configures time duration for each of checking interval
func WithInterval(interval time.Duration) Option {
	return func(o *options) {
		o.interval = interval
	}
}

// WithInvertCheck configures invert checking
func WithInvertCheck(invertCheck bool) Option {
	return func(o *options) {
		o.invertCheck = invertCheck
	}
}

// Wait waits for end up of check execution.
func Wait(check Check, opts ...Option) error {
	options := &options{
		timeout:     10 * time.Second,
		interval:    time.Second,
		invertCheck: false,
	}

	// apply the list of options to waiter
	for _, opt := range opts {
		opt(options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.timeout)
	defer cancel()

	checking := check
	if options.invertCheck == true {
		checking = func() bool { return !check() }
	}

	for !checking() {
		select {
		case <-ctx.Done():
			return errors.NewTimedOutError()
		case <-time.After(options.interval):
		}
	}

	return nil
}
