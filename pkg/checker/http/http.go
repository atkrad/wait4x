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

package http

import (
	"context"
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/atkrad/wait4x/pkg/checker/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// Option configures an HTTP.
type Option func(h *HTTP)

// HTTP represents HTTP checker
type HTTP struct {
	address          string
	timeout          time.Duration
	expectBody       string
	expectStatusCode int
}

// New creates the HTTP checker
func New(address string, opts ...Option) checker.Checker {
	h := &HTTP{
		address: address,
		timeout: time.Second * 5,
	}

	// apply the list of options to HTTP
	for _, opt := range opts {
		opt(h)
	}

	return h
}

// WithTimeout configures a time limit for requests made by the HTTP client
func WithTimeout(timeout time.Duration) Option {
	return func(h *HTTP) {
		h.timeout = timeout
	}
}

// WithExpectBody configures response body expectation
func WithExpectBody(body string) Option {
	return func(h *HTTP) {
		h.expectBody = body
	}
}

// WithExpectStatusCode configures response status code expectation
func WithExpectStatusCode(code int) Option {
	return func(h *HTTP) {
		h.expectStatusCode = code
	}
}

// Check checks HTTP connection
func (h *HTTP) Check(ctx context.Context) (err error) {
	var httpClient = &http.Client{
		Timeout: h.timeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", h.address, nil)
	if err != nil {
		return errors.Wrap(err, errors.DebugLevel)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, errors.DebugLevel)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			err = errors.Wrap(err, errors.DebugLevel)
		}
	}()

	if h.expectStatusCode != 0 && h.expectStatusCode != resp.StatusCode {
		return errors.New(
			"the status code doesn't expect",
			errors.InfoLevel,
			errors.WithFields("actual", resp.StatusCode, "expect", h.expectStatusCode),
		)
	}

	if h.expectBody != "" {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, errors.DebugLevel)
		}

		bodyString := string(bodyBytes)
		matched, _ := regexp.MatchString(h.expectBody, bodyString)

		if !matched {
			return errors.New(
				"the body doesn't expect",
				errors.InfoLevel,
				errors.WithFields("actual", h.truncateString(bodyString, 50), "expect", h.expectBody),
			)
		}
	}

	return nil
}

func (h *HTTP) truncateString(str string, num int) string {
	truncatedStr := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		truncatedStr = str[0:num] + "..."
	}

	return truncatedStr
}
