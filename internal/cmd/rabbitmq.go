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
	"wait4x.dev/v2/checker/rabbitmq"
	"wait4x.dev/v2/waiter"
)

// NewRabbitMQCommand creates the rabbitmq sub-command
func NewRabbitMQCommand() *cobra.Command {
	rabbitmqCommand := &cobra.Command{
		Use:   "rabbitmq DSN... [flags] [-- command [args...]]",
		Short: "Check RabbitMQ connection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("DSN is required argument for the rabbitmq sub-command")
			}

			return nil
		},
		Example: `
  # Checking RabbitMQ connection
  wait4x rabbitmq 'amqp://127.0.0.1:5672'

  # Checking RabbitMQ connection with credentials and vhost
  wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672/vhost'
`,
		RunE: runRabbitMQ,
	}

	rabbitmqCommand.Flags().Duration("connection-timeout", rabbitmq.DefaultConnectionTimeout, "Timeout is the maximum amount of time a dial will wait for a connection to complete.")
	rabbitmqCommand.Flags().Bool("insecure-skip-tls-verify", rabbitmq.DefaultInsecureSkipTLSVerify, "InsecureSkipTLSVerify controls whether a client verifies the server's certificate chain and hostname.")

	return rabbitmqCommand
}

func runRabbitMQ(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")
	backoffPoclicy, _ := cmd.Flags().GetString("backoff-policy")
	backoffExpMaxInterval, _ := cmd.Flags().GetDuration("backoff-exponential-max-interval")
	backoffCoefficient, _ := cmd.Flags().GetFloat64("backoff-exponential-coefficient")

	conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
	insecureSkipTLSVerify, _ := cmd.Flags().GetBool("insecure-skip-tls-verify")

	// ArgsLenAtDash returns -1 when -- was not specified
	if i := cmd.ArgsLenAtDash(); i != -1 {
		args = args[:i]
	} else {
		args = args[:]
	}

	checkers := make([]checker.Checker, 0)
	for _, arg := range args {
		rc := rabbitmq.New(
			arg,
			rabbitmq.WithTimeout(conTimeout),
			rabbitmq.WithInsecureSkipTLSVerify(insecureSkipTLSVerify),
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
