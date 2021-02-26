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
	"os"
	"time"

	errs "github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/pkg/log"
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
	// Logger is the global logger.
	Logger log.Logger
)

// NewRootCommand creates the root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wait4x",
		Short: "Wait4X allows waiting for a port or a service to enter into specify state",
		Long:  `Wait4X allows waiting for a port to enter into specify state or waiting for a service e.g. redis, mysql, postgres, ... to enter inter ready state`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			// Prevent showing usage when subcommand return error.
			cmd.SilenceUsage = true

			Logger, err = log.NewLogrus(LogLevel, os.Stdout)
			if err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().DurationVarP(&Interval, "interval", "i", 1*time.Second, "Interval time between each loop.")
	rootCmd.PersistentFlags().DurationVarP(&Timeout, "timeout", "t", 10*time.Second, "Timeout is the maximum amount of time that Wait4X will wait for a checking operation.")
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", logrus.InfoLevel.String(), "Set the logging level (\"debug\"|\"info\"|\"warn\"|\"error\"|\"fatal\")")

	return rootCmd
}

// Execute run Wait4X application
func Execute() {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTCPCommand())
	rootCmd.AddCommand(NewHTTPCommand())
	rootCmd.AddCommand(NewPostgresqlCommand())
	rootCmd.AddCommand(NewMysqlCommand())
	rootCmd.AddCommand(NewRedisCommand())
	rootCmd.AddCommand(NewVersionCommand())

	if err := rootCmd.Execute(); err != nil {
		var commandError *errs.CommandError
		if errors.As(err, &commandError) {
			os.Exit(commandError.ExitCode)
		}

		os.Exit(errs.ExitError)
	}
}
