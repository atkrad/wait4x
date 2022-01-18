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
	"github.com/atkrad/wait4x/pkg/checker/tcp"
	"github.com/atkrad/wait4x/pkg/waiter"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/spf13/cobra"
)

// NewTCPCommand creates the tcp sub-command
func NewTCPCommand() *cobra.Command {
	tcpCommand := &cobra.Command{
		Use:   "tcp ADDRESS",
		Short: "Check TCP connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("ADDRESS is required argument for the tcp command")
			}

			return nil
		},
		Example: `
  # If you want checking just tcp connection
  wait4x tcp 127.0.0.1:9090
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			invertCheck, _ := cmd.Flags().GetBool("invert-check")

			conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")

			tc := tcp.NewTCP(args[0], tcp.WithTimeout(conTimeout))
			tc.SetLogger(Logger)

			return waiter.Wait(
				tc.Check,
				waiter.WithTimeout(timeout),
				waiter.WithInterval(interval),
				waiter.WithInvertCheck(invertCheck),
			)
		},
	}

	tcpCommand.Flags().Duration("connection-timeout", time.Second*5, "Timeout is the maximum amount of time a dial will wait for a connection to complete.")

	return tcpCommand
}
