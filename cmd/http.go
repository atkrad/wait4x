package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Check HTTP connection.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("address")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		expectStatusCode, _ := cmd.Flags().GetInt("expect-status-code")
		expectBody, _ := cmd.Flags().GetString("expect-body")

		var httpClient = &http.Client{
			Timeout: timeout,
		}

		var i = 1
		for i <= RetryCount {
			log.Info("Checking HTTP connection")

			resp, err := httpClient.Get(addr)

			if err != nil {
				log.Debug(err)
				time.Sleep(Sleep)

				i += 1
				continue
			} else {
				defer resp.Body.Close()

				if httpResponseCodeExpectation(expectStatusCode, resp) && httpResponseBodyExpectation(expectBody, resp) {
					os.Exit(0)
				} else {
					i += 1
					continue
				}

				os.Exit(0)
			}
		}

		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().String("address", "http://127.0.0.1", "Http address.")
	httpCmd.Flags().Int("expect-status-code", 0, "Expect response code e.g. 200, 204, ... .")
	httpCmd.Flags().String("expect-body", "", "Expect response body pattern.")
	httpCmd.Flags().Duration("timeout", time.Second*10, "The timeout includes connection time, any redirects, and reading the response body.")
}

func httpResponseCodeExpectation(expectStatusCode int, resp *http.Response) bool {
	if expectStatusCode == 0 {
		return true
	}

	log.WithFields(log.Fields{
		"actual": resp.StatusCode,
		"expect": expectStatusCode,
	}).Info("Checking http response code expectation")

	return expectStatusCode == resp.StatusCode
}

func httpResponseBodyExpectation(expectBody string, resp *http.Response) bool {
	if expectBody == "" {
		return true
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)

	log.WithFields(log.Fields{
		"response": bodyString,
	}).Debugf("Full response of request to '%s'", resp.Request.Host)

	log.WithFields(log.Fields{
		"actual": truncateString(bodyString, 50),
		"expect": expectBody,
	}).Info("Checking http response body expectation")

	matched, _ := regexp.MatchString(expectBody, bodyString)
	return matched
}

func truncateString(str string, num int) string {
	truncatedStr := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		truncatedStr = str[0:num] + "..."
	}

	return truncatedStr
}
