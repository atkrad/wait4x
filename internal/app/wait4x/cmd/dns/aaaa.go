// Copyright 2022 Mohammad Abdolirad
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

	"github.com/atkrad/wait4x/v2/pkg/checker/dns"
	"github.com/atkrad/wait4x/v2/pkg/waiter"
	"github.com/spf13/cobra"
)

func NewAAAACommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "AAAA ADDRESS [value] [--command [args...]]",
		Short: "Check DNS AAAA records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
# Check AAAA existence
wait4x dns AAAA wait4x.dev

# Check AAAA is wait4x.dev
wait4x dns AAAA wait4x.dev '2606:4700:3033::ac43:9ab4'

# Check AAAA by defined nameserver
wait4x dns AAAA wait4x.dev '2606:4700:3033::ac43:9ab4' -n gordon.ns.cloudflare.com
`,
		RunE: runAAAA,
	}

	return command
}

func runAAAA(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	nameserver, _ := cmd.Flags().GetString("nameserver")

	address := args[0]
	var expectedValue string
	if len(args) == 2 {
		expectedValue = args[1]
	}

	dc := dns.New(dns.AAAA,
		address,
		dns.WithExpectedValue(expectedValue),
		dns.WithNameServer(nameserver),
	)
	return waiter.WaitContext(cmd.Context(),
		dc,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithInvertCheck(invertCheck),
	)
}
