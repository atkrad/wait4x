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

package temporal

import (
	"errors"
	"github.com/atkrad/wait4x/v2/pkg/checker/temporal"
	"github.com/atkrad/wait4x/v2/pkg/waiter"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

// NewServerCommand creates the server sub-command
func NewServerCommand() *cobra.Command {
	serverCommand := &cobra.Command{
		Use:   "server TARGET [flags] [-- command [args...]]",
		Short: "Check Temporal server health check",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("TARGET is required argument for the server command")
			}

			return nil
		},
		Example: `
  # Checking just Temporal server health check
  wait4x temporal server 127.0.0.1:7233

  # Checking insecure Temporal server (no TLS)
  wait4x temporal server 127.0.0.1:7233 --insecure-transport
`,
		RunE: runServer,
	}

	return serverCommand
}

func runServer(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
	insecureTransport, _ := cmd.Flags().GetBool("insecure-transport")
	insecureSkipTLSVerify, _ := cmd.Flags().GetBool("insecure-skip-tls-verify")

	logger, err := logr.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	// ArgsLenAtDash returns -1 when -- was not specified
	if i := cmd.ArgsLenAtDash(); i != -1 {
		args = args[:i]
	}

	tc := temporal.New(
		temporal.CheckModeServer,
		args[0],
		temporal.WithTimeout(conTimeout),
		temporal.WithInsecureTransport(insecureTransport),
		temporal.WithInsecureSkipTLSVerify(insecureSkipTLSVerify),
	)

	return waiter.WaitContext(
		cmd.Context(),
		tc,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(logger),
	)
}
