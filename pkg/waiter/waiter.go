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
	"context"
	"github.com/atkrad/wait4x/internal/pkg/errors"
	"time"
)

// Check represents the checker's check method.
type Check func(ctx context.Context) error

// Option configures an options
type Option func(s *options)

// options represents waiter options
type options struct {
	timeout             time.Duration
	interval            time.Duration
	invertCheck         bool
	checkerErrorChannel chan error
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

// WithCheckerErrorChannel configures checker errors
func WithCheckerErrorChannel(errChan chan error) Option {
	return func(o *options) {
		o.checkerErrorChannel = errChan
	}
}

// Wait waits for end up of check execution.
func Wait(check Check, opts ...Option) error {
	return WaitWithContext(context.Background(), check, opts...)
}

// WaitWithContext waits for end up of check execution.
func WaitWithContext(ctx context.Context, check Check, opts ...Option) error {
	options := &options{
		timeout:     10 * time.Second,
		interval:    time.Second,
		invertCheck: false,
	}

	// apply the list of options to waiter
	for _, opt := range opts {
		opt(options)
	}

	ctx, cancel := context.WithTimeout(ctx, options.timeout)
	defer cancel()

	checking := check
	//if options.invertCheck == true {
	//	checking = func(ctx context.Context) error { return !check(ctx) }
	//}

	for {
		err := checking(ctx)
		if err == nil {
			break
		}

		options.checkerErrorChannel <- err

		select {
		case <-ctx.Done():
			//return ctx.Err()
			return errors.NewTimedOutError()
		case <-time.After(options.interval):
		}
	}

	return nil
}
