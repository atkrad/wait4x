package cmd

import (
	"context"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
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
			timeout, _ := cmd.Flags().GetDuration("connection-timeout")
			expectKey, _ := cmd.Flags().GetString("expect-key")

			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			rc := checker.NewRedis(args[0], expectKey, timeout)
			rc.SetLogger(Logger)

			for !rc.Check() {
				select {
				case <-ctx.Done():
					return errors.NewTimedOutError()
				case <-time.After(Interval):
				}
			}

			return nil
		},
	}

	redisCommand.Flags().Duration("connection-timeout", time.Second*5, "Dial timeout for establishing new connections.")
	redisCommand.Flags().String("expect-key", "", "Checking key existence.")

	return redisCommand
}
