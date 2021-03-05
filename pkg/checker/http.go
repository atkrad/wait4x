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
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// HTTP represents HTTP checker
type HTTP struct {
	address          string
	timeout          time.Duration
	expectStatusCode int
	expectBody       string
	*LogAware
}

// NewHTTP creates the HTTP checker
func NewHTTP(address string, expectStatusCode int, expectBody string, timeout time.Duration) Checker {
	h := &HTTP{
		address:          address,
		expectStatusCode: expectStatusCode,
		expectBody:       expectBody,
		timeout:          timeout,
		LogAware:         &LogAware{},
	}

	return h
}

// Check checks HTTP connection
func (h *HTTP) Check() bool {
	var httpClient = &http.Client{
		Timeout: h.timeout,
	}

	h.logger.Info("Checking HTTP connection ...")

	resp, err := httpClient.Get(h.address)

	if err != nil {
		h.logger.Debug(err)

		return false
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.Debug(err)
		}
	}()

	if h.httpResponseCodeExpectation(h.expectStatusCode, resp) && h.httpResponseBodyExpectation(h.expectBody, resp) {
		return true
	}

	return false
}

func (h *HTTP) httpResponseCodeExpectation(expectStatusCode int, resp *http.Response) bool {
	if expectStatusCode == 0 {
		return true
	}

	h.logger.InfoWithFields("Checking http response code expectation", map[string]interface{}{"actual": resp.StatusCode, "expect": expectStatusCode})

	return expectStatusCode == resp.StatusCode
}

func (h *HTTP) httpResponseBodyExpectation(expectBody string, resp *http.Response) bool {
	if expectBody == "" {
		return true
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.logger.Fatal(err)
	}

	bodyString := string(bodyBytes)

	// TODO: Logging full body response in debug level.

	h.logger.InfoWithFields("Checking http response body expectation", map[string]interface{}{"actual": h.truncateString(bodyString, 50), "expect": expectBody})

	matched, _ := regexp.MatchString(expectBody, bodyString)
	return matched
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
