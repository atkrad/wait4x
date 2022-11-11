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
	"errors"
	"fmt"

	"github.com/atkrad/wait4x/v2/pkg/checker/dns"
	"github.com/atkrad/wait4x/v2/pkg/waiter"
	"github.com/spf13/cobra"
)

func NewDNSCommand() *cobra.Command {
	dnsCommand := &cobra.Command{
		Use:   "dns COMMAND [flags] [--command [args...]]",
		Short: "Check DNS records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		RunE: runDNS,
	}

	dnsACommand := &cobra.Command{
		Use:   "A ADDRESS [value] [--command [args...]]",
		Short: "Check DNS A records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
# Check A existence
wait4x dns A wait4x.dev

# Check A is wait4x.dev
wait4x dns A wait4x.dev 172.67.154.180

# Check A by defined nameserver
wait4x dns A wait4x.dev 172.67.154.180 -n gordon.ns.cloudflare.com
`,
		RunE: runA,
	}

	dnsAAAACommand := &cobra.Command{
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

	dnsCNAMECommand := &cobra.Command{
		Use:   "CNAME ADDRESS [value] [--command [args...]]",
		Short: "Check DNS CNAME records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
# Check CNAME existence
wait4x dns CNAME 172.67.154.180

# Check CNAME is wait4x.dev
wait4x dns CNAME 172.67.154.180 wait4x.dev

# Check CNAME by defined nameserver
wait4x dns CNAME 172.67.154.180 -n gordon.ns.cloudflare.com
`,
		RunE: runCNAME,
	}

	dnsTXTCommand := &cobra.Command{
		Use:   "TXT ADDRESS [value] [--command [args...]]",
		Short: "Check DNS TXT records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
# Check TXT existence
wait4x dns TXT wait4x.dev

# Check TXT is wait4x.dev
wait4x dns TXT wait4x.dev 'include:_spf.mx.cloudflare.net'

# Check TXT by defined nameserver
wait4x dns TXT wait4x.dev 'include:_spf.mx.cloudflare.net' -n gordon.ns.cloudflare.com
`,
		RunE: runTXT,
	}

	dnsMXCommand := &cobra.Command{
		Use:   "MX ADDRESS [value] [--command [args...]]",
		Short: "Check DNS MX records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
# Check MX existence
wait4x dns MX wait4x.dev

# Check MX is wait4x.dev
wait4x dns MX wait4x.dev 'route1.mx.cloudflare.net'

# Check MX by defined nameserver
wait4x dns MX wait4x.dev 'route1.mx.cloudflare.net.' -n gordon.ns.cloudflare.com
`,
		RunE: runMX,
	}

	dnsNSCommand := &cobra.Command{
		Use:   "NS ADDRESS [value] [--command [args...]]",
		Short: "Check DNS NS records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
		Example: `
# Check NS existence
wait4x dns NS wait4x.dev

# Check NS is wait4x.dev
wait4x dns NS wait4x.dev 'emma.ns.cloudflare.com'

# Check NS by defined nameserver
wait4x dns NS wait4x.dev 'emma.ns.cloudflare.com' -n gordon.ns.cloudflare.com
`,
		RunE: runNS,
	}

	dnsCommand.Flags().String("nameserver", "", " Address of the nameserver to send packets to")

	dnsCommand.AddCommand(dnsACommand)
	dnsCommand.AddCommand(dnsAAAACommand)
	dnsCommand.AddCommand(dnsCNAMECommand)
	dnsCommand.AddCommand(dnsMXCommand)
	dnsCommand.AddCommand(dnsTXTCommand)
	dnsCommand.AddCommand(dnsNSCommand)

	return dnsCommand
}

func runDNS(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("command not found")
}

func runA(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	nameserver, _ := cmd.Flags().GetString("nameserver")

	address := args[0]
	var expectedValue string
	if len(args) == 2 {
		expectedValue = args[1]
	}

	dc := dns.New(dns.A,
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

func runCNAME(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	nameserver, _ := cmd.Flags().GetString("nameserver")

	address := args[0]
	var expectedValue string
	if len(args) == 2 {
		expectedValue = args[1]
	}

	dc := dns.New(dns.CNAME,
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

func runMX(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	nameserver, _ := cmd.Flags().GetString("nameserver")

	address := args[0]
	var expectedValue string
	if len(args) == 2 {
		expectedValue = args[1]
	}

	dc := dns.New(dns.MX,
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

func runTXT(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	nameserver, _ := cmd.Flags().GetString("nameserver")

	address := args[0]
	var expectedValue string
	if len(args) == 2 {
		expectedValue = args[1]
	}

	dc := dns.New(dns.TXT,
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

func runNS(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	nameserver, _ := cmd.Flags().GetString("nameserver")

	address := args[0]
	var expectedValue string
	if len(args) == 2 {
		expectedValue = args[1]
	}

	dc := dns.New(dns.NS,
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
