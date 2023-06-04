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
	"strings"
	"wait4x.dev/v2/checker"
	"wait4x.dev/v2/checker/cassandra"
	"wait4x.dev/v2/waiter"
)

var (
	ErrCassandraInsufficientHosts = errors.New("hosts is required argument for the cassandra sub-command")
)

// NewCassandraCommand creates the cassandra sub-command
func NewCassandraCommand() *cobra.Command {
	cassandraCommand := &cobra.Command{
		Use:   "cassandra ADDRESSES... [flags] [-- command [args...]]",
		Short: "Check cassandra connection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return ErrCassandraInsufficientHosts
			}
			return nil
		},
		Example: `
	# Checking cassandra connection
	wait4x cassandra '127.0.0.1:9042'`,
		RunE: runCassandra,
	}

	cassandraCommand.Flags().String("username", "", "Cassandra cluster username")
	cassandraCommand.Flags().String("password", "", "Cassandra cluster password")

	return cassandraCommand
}

func runCassandra(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")
	backoffPolicy, _ := cmd.Flags().GetString("backoff-policy")
	backoffExpMaxInterval, _ := cmd.Flags().GetDuration("backoff-exponential-max-interval")
	backoffCoefficient, _ := cmd.Flags().GetFloat64("backoff-exponential-coefficient")

	cassandraUsername, _ := cmd.Flags().GetString("username")
	cassandraPassword, _ := cmd.Flags().GetString("password")

	cassandraHosts := strings.Split(args[0], ",")
	connectionsParams := cassandra.ConnectionParams{
		Hosts:    cassandraHosts,
		Username: cassandraUsername,
		Password: cassandraPassword,
	}

	cc := cassandra.New(connectionsParams)
	return waiter.WaitParallelContext(
		cmd.Context(),
		[]checker.Checker{cc},
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithBackoffCoefficient(backoffCoefficient),
		waiter.WithBackoffPolicy(backoffPolicy),
		waiter.WithBackoffExponentialMaxInterval(backoffExpMaxInterval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(Logger),
	)
}
