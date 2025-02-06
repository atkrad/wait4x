// Copyright 2023 The Wait4X Authors
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

package dns

import (
	"errors"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	dns "wait4x.dev/v2/checker/dns/txt"
	"wait4x.dev/v2/waiter"
)

// NewTXTCommand creates the DNS TXT command
func NewTXTCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "TXT ADDRESS [--command [args...]]",
		Aliases: []string{"txt"},
		Short:   "Check DNS TXT records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
  # Check TXT records existence
  wait4x dns TXT wait4x.dev

  # Check TXT records with expected values
  wait4x dns TXT wait4x.dev --expected-value 'include:_spf.mx.cloudflare.net'

  # Check TXT records by defined nameserver
  wait4x dns TXT wait4x.dev --expected-value 'include:_spf.mx.cloudflare.net' -n gordon.ns.cloudflare.com
`,
		RunE: runTXT,
	}

	command.Flags().StringArray("expect-value", nil, "Expect values.")

	return command
}

func runTXT(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")
	nameserver, _ := cmd.Flags().GetString("nameserver")
	expectValues, _ := cmd.Flags().GetStringArray("expect-value")

	logger, err := logr.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	address := args[0]

	dc := dns.New(
		address,
		dns.WithExpectedValues(expectValues),
		dns.WithNameServer(nameserver),
	)

	return waiter.WaitContext(cmd.Context(),
		dc,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(logger),
	)
}
