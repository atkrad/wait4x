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
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
)

// Check represents the checker's check method.
type Check func(ctx context.Context) bool

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
		checking = func(ctx context.Context) bool { return !check(ctx) }
	}

	for !checking(ctx) {
		select {
		case <-ctx.Done():
			return errors.NewTimedOutError()
		case <-time.After(options.interval):
		}
	}

	return nil
}
