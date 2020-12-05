package cmd

import (
	"errors"
	"os"
	"time"

	errs "github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Interval represents time between each loop
	Interval time.Duration
	// Timeout represents the maximum amount of time that Wait4X will wait for a checking operation
	Timeout time.Duration
	// LogLevel represents logging level e.g. info, warn, error, debug
	LogLevel string
)

// NewWait4X creates the root command
func NewWait4X() *cobra.Command {
	wait4x := &cobra.Command{
		Use:   "wait4x",
		Short: "Wait4X allows waiting for a port or a service to enter into specify state",
		Long:  `Wait4X allows waiting for a port to enter into specify state or waiting for a service e.g. redis, mysql, postgres, ... to enter inter ready state`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Prevent showing usage when subcommand return error.
			cmd.SilenceUsage = true

			lvl, err := logrus.ParseLevel(LogLevel)
			if err != nil {
				return err
			}

			logrus.SetOutput(os.Stdout)
			logrus.SetLevel(lvl)

			return nil
		},
	}

	wait4x.PersistentFlags().DurationVarP(&Interval, "interval", "i", 1*time.Second, "Interval time between each loop.")
	wait4x.PersistentFlags().DurationVarP(&Timeout, "timeout", "t", 10*time.Second, "Timeout is the maximum amount of time that Wait4X will wait for a checking operation.")
	wait4x.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", logrus.InfoLevel.String(), "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"fatal\")")

	return wait4x
}

// Execute run Wait4X application
func Execute() {
	wait4x := NewWait4X()
	wait4x.AddCommand(NewTCPCommand())
	wait4x.AddCommand(NewHTTPCommand())
	wait4x.AddCommand(NewMysqlCommand())
	wait4x.AddCommand(NewRedisCommand())
	wait4x.AddCommand(NewVersionCommand())

	if err := wait4x.Execute(); err != nil {
		var commandError *errs.CommandError
		if errors.As(err, &commandError) {
			os.Exit(commandError.ExitCode)
		}

		os.Exit(errs.ExitError)
	}
}
