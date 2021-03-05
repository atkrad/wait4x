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

package checker

import (
	"net"
	"time"
)

// TCP represents TCP checker
type TCP struct {
	address string
	timeout time.Duration
	*LogAware
}

// NewTCP creates the TCP checker
func NewTCP(address string, timeout time.Duration) Checker {
	t := &TCP{
		address:  address,
		timeout:  timeout,
		LogAware: &LogAware{},
	}

	return t
}

// Check checks TCP connection
func (t *TCP) Check() bool {
	d := net.Dialer{Timeout: t.timeout}

	t.logger.Info("Checking TCP connection ...")

	_, err := d.Dial("tcp", t.address)
	if err != nil {
		t.logger.Debug(err)

		return false
	}

	return true
}
