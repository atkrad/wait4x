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
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/atkrad/wait4x/pkg/checker/errors"
	"github.com/tidwall/gjson"
)

// Option configures an HTTP.
type Option func(h *HTTP)

// HTTP represents HTTP checker
type HTTP struct {
	address          string
	timeout          time.Duration
	expectBodyRegex  string
	expectBodyJSON   string
	expectHeader     string
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

// WithExpectBodyRegex configures response body expectation
func WithExpectBodyRegex(body string) Option {
	return func(h *HTTP) {
		h.expectBodyRegex = body
	}
}

// WithExpectBodyJSON configures response json expectation
func WithExpectBodyJSON(json string) Option {
	return func(h *HTTP) {
		h.expectBodyJSON = json
	}
}

// WithExpectHeader configures response header expectation
func WithExpectHeader(header string) Option {
	return func(h *HTTP) {
		h.expectHeader = header
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

	if h.expectStatusCode != 0 {
		err := h.checkingStatusCodeExpectation(resp)

		if err != nil {
			return err
		}
	}

	if h.expectBodyRegex != "" {
		err := h.checkingBodyExpectation(resp)

		if err != nil {
			return err
		}
	}

	if h.expectBodyJSON != "" {
		err := h.checkingJSONExpectation(resp)

		if err != nil {
			return err
		}
	}

	if h.expectHeader != "" {
		return h.checkingHeaderExpectation(resp)
	}

	return nil
}

func (h *HTTP) checkingStatusCodeExpectation(resp *http.Response) error {
	if h.expectStatusCode != resp.StatusCode {
		return errors.New(
			"the status code doesn't expect",
			errors.InfoLevel,
			errors.WithFields("actual", resp.StatusCode, "expect", h.expectStatusCode),
		)
	}

	return nil
}

func (h *HTTP) checkingBodyExpectation(resp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, errors.DebugLevel)
	}

	bodyString := string(bodyBytes)
	matched, _ := regexp.MatchString(h.expectBodyRegex, bodyString)

	if !matched {
		return errors.New(
			"the body doesn't expect",
			errors.InfoLevel,
			errors.WithFields("actual", h.truncateString(bodyString, 50), "expect", h.expectBodyRegex),
		)
	}

	return nil
}

func (h *HTTP) checkingJSONExpectation(resp *http.Response) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, errors.DebugLevel)
	}

	bodyString := string(bodyBytes)
	value := gjson.Get(bodyString, h.expectBodyJSON)

	if !value.Exists() {
		return errors.New(
			"the JSON doesn't match",
			errors.InfoLevel,
			errors.WithFields("actual", h.truncateString(bodyString, 50), "expect", h.expectBodyJSON),
		)
	}

	return nil
}

func (h *HTTP) checkingHeaderExpectation(resp *http.Response) error {
	// Key value. e.g. Content-Type=application/json
	expectedHeaderParsed := strings.SplitN(h.expectHeader, "=", 2)
	if len(expectedHeaderParsed) == 2 {
		headerValue := resp.Header.Get(expectedHeaderParsed[0])
		matched, _ := regexp.MatchString(expectedHeaderParsed[1], headerValue)
		if !matched {
			return errors.New(
				"the http header key and value doesn't expect",
				errors.InfoLevel,
				errors.WithFields("actual", headerValue, "expect", h.expectHeader),
			)
		}
	}

	// Only key.
	if _, ok := resp.Header[expectedHeaderParsed[0]]; !ok {
		return errors.New(
			"the http header key doesn't expect",
			errors.InfoLevel,
			errors.WithFields("actual", resp.Header, "expect", h.expectHeader),
		)
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
