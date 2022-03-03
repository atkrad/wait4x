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
	"context"
	"errors"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Logger is the global logger.
var Logger logr.Logger

// ExitError is the error exit code
const ExitError = 1

// ExitTimedOut is the timed out exit code
const ExitTimedOut = 124

// NewRootCommand creates the root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wait4x",
		Short: "Wait4X allows waiting for a port or a service to enter into specify state",
		Long:  `Wait4X allows waiting for a port to enter into specify state or waiting for a service e.g. redis, mysql, postgres, ... to enter inter ready state`,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			logLevel, _ := cmd.Flags().GetString("log-level")
			lvl, err := zerolog.ParseLevel(logLevel)
			if err != nil {
				return err
			}

			// Prevent showing usage when subcommand return error.
			cmd.SilenceUsage = true

			zl := zerolog.New(
				zerolog.ConsoleWriter{
					Out:        os.Stderr,
					NoColor:    false,
					TimeFormat: time.RFC3339,
				},
			).Level(lvl).
				With().
				Timestamp().
				Logger()
			Logger = zerologr.New(&zl)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				command := args[1:][0]
				arguments := args[1:][1:]
				for i, arg := range arguments {
					arguments[i] = os.ExpandEnv(arg)
				}

				c := exec.CommandContext(cmd.Context(), command, arguments...)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr

				return c.Run()
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().DurationP("interval", "i", 1*time.Second, "Interval time between each loop.")
	rootCmd.PersistentFlags().DurationP("timeout", "t", 10*time.Second, "Timeout is the maximum amount of time that Wait4X will wait for a checking operation.")
	rootCmd.PersistentFlags().BoolP("invert-check", "v", false, "Invert the sense of checking.")
	rootCmd.PersistentFlags().StringP("log-level", "l", zerolog.InfoLevel.String(), "Set the logging level (\"trace\"|\"debug\"|\"info\")")

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
	rootCmd.AddCommand(NewInfluxDBCommand())
	rootCmd.AddCommand(NewMongoDBCommand())
	rootCmd.AddCommand(NewRabbitMQCommand())
	rootCmd.AddCommand(NewVersionCommand())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			os.Exit(ExitTimedOut)
		}

		os.Exit(ExitError)
	}
}
