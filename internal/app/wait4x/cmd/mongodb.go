// Copyright 2022 The Wait4X Authors
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
	"wait4x.dev/v2/checker/mongodb"
	"wait4x.dev/v2/waiter"
)

// NewMongoDBCommand creates the mongodb sub-command
func NewMongoDBCommand() *cobra.Command {
	mongodbCommand := &cobra.Command{
		Use:   "mongodb DSN... [flags] [-- command [args...]]",
		Short: "Check MongoDB connection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("DSN is required argument for the mongodb command")
			}

			return nil
		},
		Example: `
  # Checking MongoDB connection
  wait4x mongodb 'mongodb://127.0.0.1:27017'

  # Checking MongoDB connection with credentials and options
  wait4x mongodb 'mongodb://user:pass@127.0.0.1:27017/?maxPoolSize=20&w=majority'
`,
		RunE: runMongoDB,
	}

	return mongodbCommand
}

func runMongoDB(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")
	backoffPoclicy, _ := cmd.Flags().GetString("backoff-policy")
	backoffExpMaxInterval, _ := cmd.Flags().GetDuration("backoff-exponential-max-interval")
	backoffCoefficient, _ := cmd.Flags().GetFloat64("backoff-exponential-coefficient")

	// ArgsLenAtDash returns -1 when -- was not specified
	if i := cmd.ArgsLenAtDash(); i != -1 {
		args = args[:i]
	} else {
		args = args[:]
	}

	checkers := make([]checker.Checker, 0)
	for _, arg := range args {
		mc := mongodb.New(arg)

		checkers = append(checkers, mc)
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
