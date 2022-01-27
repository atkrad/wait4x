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
	"github.com/atkrad/wait4x/pkg/checker/http"
	"github.com/atkrad/wait4x/pkg/waiter"
	"net/url"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/spf13/cobra"
)

// NewHTTPCommand creates the http sub-command
func NewHTTPCommand() *cobra.Command {
	httpCommand := &cobra.Command{
		Use:   "http ADDRESS",
		Short: "Check HTTP connection.",
		Long:  "",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("ADDRESS is required argument for the http command")
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
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			invertCheck, _ := cmd.Flags().GetBool("invert-check")

			expectStatusCode, _ := cmd.Flags().GetInt("expect-status-code")
			expectBody, _ := cmd.Flags().GetString("expect-body")
			connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")

			hc := http.NewHTTP(args[0],
				http.WithExpectStatusCode(expectStatusCode),
				http.WithExpectBody(expectBody),
				http.WithTimeout(connectionTimeout),
			)
			hc.SetLogger(Logger)

			return waiter.Wait(
				hc.Check,
				waiter.WithTimeout(timeout),
				waiter.WithInterval(interval),
				waiter.WithInvertCheck(invertCheck),
			)
		},
	}

	httpCommand.Flags().Int("expect-status-code", 0, "Expect response code e.g. 200, 204, ... .")
	httpCommand.Flags().String("expect-body", "", "Expect response body pattern.")
	httpCommand.Flags().Duration("connection-timeout", time.Second*5, "Http connection timeout, The timeout includes connection time, any redirects, and reading the response body.")

	return httpCommand
}
