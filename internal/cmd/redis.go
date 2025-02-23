// Copyright 2020 The Wait4X Authors
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

	"github.com/spf13/cobra"
	"wait4x.dev/v2/checker"
	"wait4x.dev/v2/checker/redis"
	"wait4x.dev/v2/waiter"
)

// NewRedisCommand creates the redis sub-command
func NewRedisCommand() *cobra.Command {
	redisCommand := &cobra.Command{
		Use:   "redis ADDRESS... [flags] [-- command [args...]]",
		Short: "Check Redis connection or key existence",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the redis command")
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
		RunE: runRedis,
	}

	redisCommand.Flags().Duration("connection-timeout", redis.DefaultConnectionTimeout, "Dial timeout for establishing new connections.")
	redisCommand.Flags().String("expect-key", "", "Checking key existence.")

	return redisCommand
}

func runRedis(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")
	backoffPoclicy, _ := cmd.Flags().GetString("backoff-policy")
	backoffExpMaxInterval, _ := cmd.Flags().GetDuration("backoff-exponential-max-interval")
	backoffCoefficient, _ := cmd.Flags().GetFloat64("backoff-exponential-coefficient")

	conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
	expectKey, _ := cmd.Flags().GetString("expect-key")

	// ArgsLenAtDash returns -1 when -- was not specified
	if i := cmd.ArgsLenAtDash(); i != -1 {
		args = args[:i]
	} else {
		args = args[:]
	}

	checkers := make([]checker.Checker, 0)
	for _, arg := range args {
		rc := redis.New(
			arg,
			redis.WithExpectKey(expectKey),
			redis.WithTimeout(conTimeout),
		)

		checkers = append(checkers, rc)
	}

	return waiter.WaitParallelContext(
		cmd.Context(),
		checkers,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithBackoffCoefficient(backoffCoefficient),
		waiter.WithBackoffPolicy(backoffPoclicy),
		waiter.WithBackoffExponentialMaxInterval(backoffExpMaxInterval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(Logger),
	)
}
