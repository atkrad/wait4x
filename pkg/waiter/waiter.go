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
	"errors"
	errors2 "github.com/atkrad/wait4x/pkg/checker/errors"
	"github.com/go-logr/logr"
	"time"
)

// Check represents the checker's check method.
type Check func(ctx context.Context) error

// Option configures an options
type Option func(s *options)

// options represents waiter options
type options struct {
	timeout     time.Duration
	interval    time.Duration
	invertCheck bool
	logger      *logr.Logger
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

// WithLogger configures waiter logger
func WithLogger(logger *logr.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Wait waits for end up of check execution.
func Wait(check Check, opts ...Option) error {
	return WaitContext(context.Background(), check, opts...)
}

// WaitWithContext waits for end up of check execution.
// Deprecated: The function will be removed in v3.0.0, please use the WaitContext.
func WaitWithContext(ctx context.Context, check Check, opts ...Option) error {
	return WaitContext(ctx, check, opts...)
}

// WaitContext waits for end up of check execution.
func WaitContext(ctx context.Context, check Check, opts ...Option) error {
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

	for {
		if options.logger != nil {
			options.logger.Info("Checking the service ...")
		}

		err := check(ctx)
		if options.logger != nil {
			var cError *errors2.Error
			if errors.As(err, &cError) {
				options.logger.V(int(cError.Level())).
					WithValues(cError.Fields()...).
					Info(err.Error())
			}
		}

		if options.invertCheck == true {
			if err == nil {
				goto CONTINUE
			}

			break
		}

		if err == nil {
			break
		}

	CONTINUE:

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(options.interval):
		}
	}

	return nil
}
