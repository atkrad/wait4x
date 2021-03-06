// Copyright 2020 Mohammad Abdolirad
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
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/internal/pkg/waiter"
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/spf13/cobra"
)

// NewRedisCommand creates the redis sub-command
func NewRedisCommand() *cobra.Command {
	redisCommand := &cobra.Command{
		Use:   "redis ADDRESS",
		Short: "Check Redis connection or key existence.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("ADDRESS is required argument for the redis command")
			}

			return nil
		},
		Example: `
  # Checking Redis connection
  wait4x redis redis://127.0.0.1:6379

  # Specify username, password and db
  wait4x redis redis://user:password@localhost:6379/1

  # Checking Redis connection over unix socket
  wait4x redis unix://user:password@/path/to/redis.sock?db=1

  # Checking a key existence
  wait4x redis redis://127.0.0.1:6379 --expect-key FOO

  # Checking a key existence and matching the value
  wait4x redis redis://127.0.0.1:6379 --expect-key "FOO=^b[A-Z]r$"
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			timeout, _ := cmd.Flags().GetDuration("timeout")

			conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
			expectKey, _ := cmd.Flags().GetString("expect-key")

			rc := checker.NewRedis(args[0], expectKey, conTimeout)
			rc.SetLogger(Logger)

			return waiter.Wait(rc.Check, timeout, interval)
		},
	}

	redisCommand.Flags().Duration("connection-timeout", time.Second*5, "Dial timeout for establishing new connections.")
	redisCommand.Flags().String("expect-key", "", "Checking key existence.")

	return redisCommand
}
