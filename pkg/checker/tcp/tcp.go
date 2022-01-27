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

package tcp

import (
	"context"
	"github.com/atkrad/wait4x/pkg/checker"
	"net"
	"time"
)

// Option configures a TCP.
type Option func(s *TCP)

// TCP represents TCP checker
type TCP struct {
	address string
	timeout time.Duration
	*checker.LogAware
}

// New creates the TCP checker
func New(address string, opts ...Option) checker.Checker {
	t := &TCP{
		address:  address,
		timeout:  time.Second * 5,
		LogAware: &checker.LogAware{},
	}

	// apply the list of options to TCP
	for _, opt := range opts {
		opt(t)
	}

	return t
}

// WithTimeout configures a timeout for maximum amount of time a dial will wait for a connection to complete
func WithTimeout(timeout time.Duration) Option {
	return func(h *TCP) {
		h.timeout = timeout
	}
}

// Check checks TCP connection
func (t *TCP) Check(ctx context.Context) bool {
	d := net.Dialer{Timeout: t.timeout}

	t.Logger().Info("Checking TCP connection ...")

	_, err := d.DialContext(ctx, "tcp", t.address)
	if err != nil {
		t.Logger().Debug(err)

		return false
	}

	return true
}
