package cmd

import (
	"context"
	"github.com/atkrad/wait4x/internal/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"time"
)

func NewTcpCommand() *cobra.Command {
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
			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			for !checkingTcp(cmd, args) {
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

func checkingTcp(cmd *cobra.Command, args []string) bool {
	connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
	log.Info("Checking TCP connection ...")

	d := net.Dialer{Timeout: connectionTimeout}
	_, err := d.Dial("tcp", args[0])
	if err != nil {
		log.Debug(err)

		return false
	}

	return true
}
