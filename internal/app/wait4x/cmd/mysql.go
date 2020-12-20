package cmd

import (
	"context"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/spf13/cobra"
)

// NewMysqlCommand creates the mysql sub-command
func NewMysqlCommand() *cobra.Command {
	mysqlCommand := &cobra.Command{
		Use:   "mysql DSN",
		Short: "Check MySQL connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("DSN is required argument for the mysql command")
			}

			return nil
		},
		Example: `
  # Checking MySQL TCP connection
  wait4x mysql user:password@tcp(localhost:5555)/dbname?tls=skip-verify

  # Checking MySQL UNIX Socket existence
  wait4x mysql usernname:password@unix(/tmp/mysql.sock)/myDatabase
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			mc := checker.NewMySQL(args[0])
			mc.SetLogger(Logger)

			for !mc.Check() {
				select {
				case <-ctx.Done():
					return errors.NewTimedOutError()
				case <-time.After(Interval):
				}
			}

			return nil
		},
	}

	return mysqlCommand
}
