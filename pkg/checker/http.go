package checker

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/atkrad/wait4x/pkg/log"
)

// HTTP represents HTTP checker
type HTTP struct {
	address          string
	timeout          time.Duration
	expectStatusCode int
	expectBody       string
	logger           log.Logger
}

// NewHTTP creates the HTTP checker
func NewHTTP(address string, expectStatusCode int, expectBody string, timeout time.Duration) Checker {
	h := &HTTP{
		address:          address,
		expectStatusCode: expectStatusCode,
		expectBody:       expectBody,
		timeout:          timeout,
	}

	return h
}

func (h *HTTP) SetLogger(logger log.Logger) {
	h.logger = logger
}

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

	defer resp.Body.Close()

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
