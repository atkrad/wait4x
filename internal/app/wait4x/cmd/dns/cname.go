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
	dns "wait4x.dev/v2/checker/dns/cname"
	"wait4x.dev/v2/waiter"
)

// NewCNAMECommand creates a new Cobra command for the "dns CNAME" subcommand. This command
// checks DNS CNAME records and optionally verifies expected domain names. It supports various
// configuration options such as timeout, interval, nameserver, and expected domains.
func NewCNAMECommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "CNAME ADDRESS [-- command [args...]]",
		Aliases: []string{"cname"},
		Short:   "Check DNS CNAME records for a given domain",
		Long:    "Check DNS CNAME records and optionally verify expected domain names",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the CNAME command")
			}
			return nil
		},
		Example: `
  # Check CNAME record existence
  wait4x dns CNAME example.com

  # Check CNAME records with expected domains
  wait4x dns CNAME example.com --expected-domain target.example.com

  # Check CNAME record using a specific nameserver
  wait4x dns CNAME example.com --expected-domain target.example.com -n 8.8.8.8

  # Check CNAME record with custom timeout and interval
  wait4x dns CNAME example.com --timeout 30s --interval 5s
`,
		RunE: runCNAME,
	}

	command.Flags().StringArray("expect-domain", nil, "Expect domains.")

	return command
}

// runCNAME is the command handler for the "dns CNAME" command. It checks DNS CNAME records and
// optionally verifies expected domain names. It uses the provided command flags to configure the
// DNS check, such as timeout, interval, nameserver, and expected domains. The function returns
// an error if the DNS check fails.
func runCNAME(cmd *cobra.Command, args []string) error {
	interval, err := cmd.Flags().GetDuration("interval")
	if err != nil {
		return err
	}

	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return err
	}

	invertCheck, err := cmd.Flags().GetBool("invert-check")
	if err != nil {
		return err
	}

	nameserver, err := cmd.Flags().GetString("nameserver")
	if err != nil {
		return err
	}

	expectDomains, err := cmd.Flags().GetStringArray("expect-domain")
	if err != nil {
		return err
	}

	logger, err := logr.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	address := args[0]

	dc := dns.New(
		address,
		dns.WithExpectedDomains(expectDomains),
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
