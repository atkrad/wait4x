package cmd

import (
	"context"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/spf13/cobra"
)

// NewPostgresqlCommand creates the postgresql sub-command
func NewPostgresqlCommand() *cobra.Command {
	postgresqlCommand := &cobra.Command{
		Use:   "postgresql DSN",
		Aliases: []string{"postgres", "postgre"},
		Short: "Check PostgreSQL connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("DSN is required argument for the postgresql command")
			}

			return nil
		},
		Example: `
  # Checking PostgreSQL TCP connection
  wait4x postgresql postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			mc := checker.NewPostgreSQL(args[0])
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

	return postgresqlCommand
}
