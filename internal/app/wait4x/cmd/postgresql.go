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
	"errors"
	"github.com/atkrad/wait4x/pkg/checker/postgresql"
	"github.com/atkrad/wait4x/pkg/waiter"
	"github.com/spf13/cobra"
)

// NewPostgresqlCommand creates the postgresql sub-command
func NewPostgresqlCommand() *cobra.Command {
	postgresqlCommand := &cobra.Command{
		Use:     "postgresql DSN [flags] [-- command [args...]]",
		Aliases: []string{"postgres", "postgre"},
		Short:   "Check PostgreSQL connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("DSN is required argument for the postgresql command")
			}

			return nil
		},
		Example: `
  # Checking PostgreSQL TCP connection
  wait4x postgresql postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full
`,
		RunE: runPostgresql,
	}

	return postgresqlCommand
}

func runPostgresql(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	pc := postgresql.New(args[0])

	return waiter.WaitContext(
		cmd.Context(),
		pc.Check,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(&Logger),
	)
}
