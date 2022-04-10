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
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/atkrad/wait4x/pkg/checker/mysql"
	"github.com/atkrad/wait4x/pkg/waiter"
	"github.com/spf13/cobra"
)

// NewMysqlCommand creates the mysql sub-command
func NewMysqlCommand() *cobra.Command {
	mysqlCommand := &cobra.Command{
		Use:   "mysql DSN... [flags] [-- command [args...]]",
		Short: "Check MySQL connection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("DSN is required argument for the mysql command")
			}

			return nil
		},
		Example: `
  # Checking MySQL TCP connection
  wait4x mysql user:password@tcp(localhost:5555)/dbname?tls=skip-verify

  # Checking MySQL UNIX Socket existence
  wait4x mysql username:password@unix(/tmp/mysql.sock)/myDatabase
`,
		RunE: runMysql,
	}

	return mysqlCommand
}

func runMysql(cmd *cobra.Command, args []string) error {
	interval, _ := cmd.Flags().GetDuration("interval")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	invertCheck, _ := cmd.Flags().GetBool("invert-check")

	// ArgsLenAtDash returns -1 when -- was not specified
	if i := cmd.ArgsLenAtDash(); i != -1 {
		args = args[:i]
	} else {
		args = args[:]
	}

	checkers := make([]checker.Checker, 0)
	for _, arg := range args {
		mc := mysql.New(arg)

		checkers = append(checkers, mc)
	}

	return waiter.WaitParallelContext(
		cmd.Context(),
		checkers,
		waiter.WithTimeout(timeout),
		waiter.WithInterval(interval),
		waiter.WithInvertCheck(invertCheck),
		waiter.WithLogger(&Logger),
	)
}
