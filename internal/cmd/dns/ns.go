// Copyright 2019-2025 The Wait4X Authors
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
	dns "wait4x.dev/v2/checker/dns/ns"
	"wait4x.dev/v2/waiter"
)

// NewNSCommand creates the DNS NS command
func NewNSCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "NS ADDRESS [-- command [args...]]",
		Aliases: []string{"ns"},
		Short:   "Check DNS NS records for a given domain",
		Long:    "Check DNS NS records for a given domain name and verify nameserver records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the NS command")
			}

			return nil
		},
		Example: `
  # Check NS records existence
  wait4x dns NS wait4x.dev

  # Check NS records with expected nameservers
  wait4x dns NS wait4x.dev --expect-nameserver 'emma.ns.cloudflare.com'

  # Check NS records with multiple expected nameservers
  wait4x dns NS wait4x.dev --expect-nameserver 'emma.ns.cloudflare.com' --expect-nameserver 'gordon.ns.cloudflare.com'

  # Check NS records using a specific nameserver
  wait4x dns NS wait4x.dev --nameserver '8.8.8.8:53'

  # Check NS records with timeout and interval
  wait4x dns NS wait4x.dev --timeout 60s --interval 5s

  # Invert the check (wait until NS records don't match)
  wait4x dns NS wait4x.dev --expect-nameserver 'emma.ns.cloudflare.com' --invert-check`,
		RunE: runNS,
	}

	command.Flags().StringArray("expect-nameserver", nil, "Expect nameservers.")

	return command
}

func runNS(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")
	nameserver, _ := cmd.Flags().GetString("nameserver")
	expectNameservers, _ := cmd.Flags().GetStringArray("expect-nameserver")

	logger, err := logr.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	address := args[0]

	dc := dns.New(
		address,
		dns.WithExpectedNameservers(expectNameservers),
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
