package cmd

import (
	"context"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/spf13/cobra"
)

// NewTCPCommand creates the tcp sub-command
func NewTCPCommand() *cobra.Command {
	tcpCommand := &cobra.Command{
		Use:   "tcp ADDRESS",
		Short: "Check TCP connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("ADDRESS is required argument for the tcp command")
			}

			return nil
		},
		Example: `
  # If you want checking just tcp connection
  wait4x tcp 127.0.0.1:9090
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, _ := cmd.Flags().GetDuration("connection-timeout")

			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			tc := checker.NewTCP(args[0], timeout)
			tc.SetLogger(Logger)

			for !tc.Check() {
				select {
				case <-ctx.Done():
					return errors.NewTimedOutError()
				case <-time.After(Interval):
				}
			}

			return nil
		},
	}

	tcpCommand.Flags().Duration("connection-timeout", time.Second*5, "Timeout is the maximum amount of time a dial will wait for a connection to complete.")

	return tcpCommand
}
