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
	"fmt"
	"github.com/atkrad/wait4x/v2/pkg/checker"
	errors2 "github.com/atkrad/wait4x/v2/pkg/checker/errors"
	"github.com/go-logr/logr"
	"reflect"
	"sync"
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

// WaitParallel waits for end up all of checks execution.
func WaitParallel(checkers []checker.Checker, opts ...Option) error {
	return WaitParallelContext(context.Background(), checkers, opts...)
}

// WaitParallelContext waits for end up all of checks execution.
func WaitParallelContext(ctx context.Context, checkers []checker.Checker, opts ...Option) error {
	// Make channels to pass wgErrors in WaitGroup
	wgErrors := make(chan error)
	wgDone := make(chan bool)

	var wg sync.WaitGroup

	for _, chr := range checkers {
		wg.Add(1)

		go func(chr checker.Checker) {
			defer wg.Done()

			err := WaitContext(ctx, chr, opts...)
			if err != nil {
				wgErrors <- err
			}
		}(chr)
	}

	// Important final goroutine to wait until WaitGroup is done
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		return nil
	case err := <-wgErrors:
		close(wgErrors)

		return err
	}
}

// Wait waits for end up of check execution.
func Wait(checker checker.Checker, opts ...Option) error {
	return WaitContext(context.Background(), checker, opts...)
}

// WaitWithContext waits for end up of check execution.
// Deprecated: The function will be removed in v3.0.0, please use the WaitContext.
func WaitWithContext(ctx context.Context, checker checker.Checker, opts ...Option) error {
	return WaitContext(ctx, checker, opts...)
}

// WaitContext waits for end up of check execution.
func WaitContext(ctx context.Context, checker checker.Checker, opts ...Option) error {
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

	var chkName string
	if t := reflect.TypeOf(checker); t.Kind() == reflect.Ptr {
		chkName = t.Elem().Name()
	} else {
		chkName = t.Name()
	}

	chkID, err := checker.Identity()
	if err != nil {
		return err
	}

	for {
		if options.logger != nil {
			options.logger.Info(fmt.Sprintf("[%s] Checking the %s ...", chkName, chkID))
		}

		err := checker.Check(ctx)
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
