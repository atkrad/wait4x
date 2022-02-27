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

package cmd

import (
	"errors"
	"github.com/atkrad/wait4x/pkg/checker/rabbitmq"
	"github.com/atkrad/wait4x/pkg/waiter"
	"github.com/spf13/cobra"
)

// NewRabbitMQCommand creates the rabbitmq sub-command
func NewRabbitMQCommand() *cobra.Command {
	rabbitmqCommand := &cobra.Command{
		Use:   "rabbitmq DSN",
		Short: "Check RabbitMQ connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("DSN is required argument for the rabbitmq sub-command")
			}

			return nil
		},
		Example: `
  # Checking RabbitMQ connection
  wait4x rabbitmq 'amqp://127.0.0.1:5672'

  # Checking RabbitMQ connection with credentials
  wait4x rabbitmq 'amqp://guest:guest@127.0.0.1:5672'
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			invertCheck, _ := cmd.Flags().GetBool("invert-check")

			conTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
			insecureSkipTLSVerify, _ := cmd.Flags().GetBool("insecure-skip-tls-verify")

			rc := rabbitmq.New(
				args[0],
				rabbitmq.WithTimeout(conTimeout),
				rabbitmq.WithInsecureSkipTLSVerify(insecureSkipTLSVerify),
			)

			return waiter.WaitWithContext(
				cmd.Context(),
				rc.Check,
				waiter.WithTimeout(timeout),
				waiter.WithInterval(interval),
				waiter.WithInvertCheck(invertCheck),
				waiter.WithLogger(&Logger),
			)
		},
	}

	rabbitmqCommand.Flags().Duration("connection-timeout", rabbitmq.DefaultConnectionTimeout, "Timeout is the maximum amount of time a dial will wait for a connection to complete.")
	rabbitmqCommand.Flags().Bool("insecure-skip-tls-verify", rabbitmq.DefaultInsecureSkipTLSVerify, "InsecureSkipTLSVerify controls whether a client verifies the server's certificate chain and hostname.")

	return rabbitmqCommand
}
