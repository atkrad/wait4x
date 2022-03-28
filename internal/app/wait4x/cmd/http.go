// Copyright 2019 Mohammad Abdolirad
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

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	net_http "net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"github.com/atkrad/wait4x/pkg/checker/http"
	"github.com/atkrad/wait4x/pkg/waiter"

	"github.com/spf13/cobra"
)

// NewHTTPCommand creates the http sub-command
func NewHTTPCommand() *cobra.Command {
	httpCommand := &cobra.Command{
		Use:   "http ADDRESS [flags] [-- command [args...]]",
		Short: "Check HTTP connection.",
		Long:  "",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the http command")
			}

			_, err := url.Parse(args[0])
			if err != nil {
				return err
			}

			return nil
		},
		Example: `
  # If you want checking just http connection
  wait4x http https://ifconfig.co

  # If you want checking http connection and expect specify http status code
  wait4x http https://ifconfig.co --expect-status-code 200

  # If you want to check a http response header
  # NOTE: the value in the expected header is regex.
  # Sample response header: Authorization Token 1234ABCD
  # You can match it by these ways:

  # Full key value:
  wait4x http https://ifconfig.co --expect-header "Authorization=Token 1234ABCD"

  # Value starts with:
  wait4x http https://ifconfig.co --expect-header "Authorization=Token"

  # Regex value:
  wait4x http https://ifconfig.co --expect-header "Authorization=Token\s.+"

  # Body JSON
  wait4x http https://ifconfig.co/json --expect-body-json "user_agent.product"
  To know more about JSON syntax https://github.com/tidwall/gjson/blob/master/SYNTAX.md

  # Body XPath
  wait4x http https://www.kernel.org/ --expect-body-xpath "//*[@id="tux-gear"]"

  # Request headers:
  wait4x http https://ifconfig.co --request-header "Content-Type: application/json" --request-header "Authorization: Token 123"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			invertCheck, _ := cmd.Flags().GetBool("invert-check")

			expectStatusCode, _ := cmd.Flags().GetInt("expect-status-code")
			expectBodyRegex, _ := cmd.Flags().GetString("expect-body-regex")
			expectBody, _ := cmd.Flags().GetString("expect-body")
			expectBodyJSON, _ := cmd.Flags().GetString("expect-body-json")
			expectBodyXPath, _ := cmd.Flags().GetString("expect-body-xpath")
			expectHeader, _ := cmd.Flags().GetString("expect-header")
			requestRawHeaders, _ := cmd.Flags().GetStringArray("request-header")
			connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
			insecureSkipTLSVerify, _ := cmd.Flags().GetBool("insecure-skip-tls-verify")

			if len(expectBody) != 0 {
				expectBodyRegex = expectBody
			}

			// Convert raw headers (e.g. 'a: b') into a http Header.
			var requestHeaders net_http.Header
			if len(requestRawHeaders) > 0 {
				rawHTTPHeaders := strings.Join(requestRawHeaders, "\r\n")
				tpReader := textproto.NewReader(bufio.NewReader(strings.NewReader(rawHTTPHeaders)))
				MIMEHeaders, err := tpReader.ReadMIMEHeader()
				if err != nil && err != io.EOF {
					return fmt.Errorf("can't parse the request header: %w", err)
				}
				requestHeaders = net_http.Header(MIMEHeaders)
			}

			hc := http.New(args[0],
				http.WithExpectStatusCode(expectStatusCode),
				http.WithExpectBodyRegex(expectBodyRegex),
				http.WithExpectBodyJSON(expectBodyJSON),
				http.WithExpectBodyXPath(expectBodyXPath),
				http.WithExpectHeader(expectHeader),
				http.WithRequestHeaders(requestHeaders),
				http.WithTimeout(connectionTimeout),
				http.WithInsecureSkipTLSVerify(insecureSkipTLSVerify),
			)

			return waiter.WaitContext(
				cmd.Context(),
				hc.Check,
				waiter.WithTimeout(timeout),
				waiter.WithInterval(interval),
				waiter.WithInvertCheck(invertCheck),
				waiter.WithLogger(&Logger),
			)
		},
	}

	httpCommand.Flags().Int("expect-status-code", 0, "Expect response code e.g. 200, 204, ... .")
	httpCommand.Flags().String("expect-body", "", "Expect response body pattern.")
	httpCommand.Flags().MarkDeprecated("expect-body", "This flag will be removed in v3.0.0, please use --expect-body-regex.")
	httpCommand.Flags().String("expect-body-regex", "", "Expect response body pattern.")
	httpCommand.Flags().String("expect-body-json", "", "Expect response body JSON pattern.")
	httpCommand.Flags().String("expect-body-xpath", "", "Expect response body XPath pattern.")
	httpCommand.Flags().String("expect-header", "", "Expect response header pattern.")
	httpCommand.Flags().StringArray("request-header", nil, "User request headers.")
	httpCommand.Flags().Duration("connection-timeout", time.Second*5, "Http connection timeout, The timeout includes connection time, any redirects, and reading the response body.")
	httpCommand.Flags().Bool("insecure-skip-tls-verify", false, "Skips tls certificate checks for the HTTPS request.")

	return httpCommand
}
