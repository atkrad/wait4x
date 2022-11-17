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
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	dnsCmd "github.com/atkrad/wait4x/v2/internal/app/wait4x/cmd/dns"
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
			noColor, _ := cmd.Flags().GetBool("no-color")
			quiet, _ := cmd.Flags().GetBool("quiet")

			// Prevent showing error when the quiet mode enabled.
			cmd.SilenceErrors = quiet
			lvl := zerolog.Disabled

			if !quiet {
				logLevel, _ := cmd.Flags().GetString("log-level")
				lvl, err = zerolog.ParseLevel(logLevel)
				if err != nil {
					return err
				}
			}

			// Prevent showing usage when subcommand return error.
			cmd.SilenceUsage = true

			zl := zerolog.New(
				zerolog.ConsoleWriter{
					Out:        os.Stderr,
					NoColor:    color.NoColor || noColor,
					TimeFormat: time.RFC3339,
				},
			).Level(lvl).
				With().
				Timestamp().
				Logger()
			Logger = zerologr.New(&zl)
			// VerbosityFieldName (v) is not emitted.
			zerologr.VerbosityFieldName = ""

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.ArgsLenAtDash() != -1 && (len(args)-cmd.ArgsLenAtDash()) > 0 {
				command := args[cmd.ArgsLenAtDash():][0]
				arguments := args[cmd.ArgsLenAtDash():][1:]
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
	rootCmd.PersistentFlags().DurationP("timeout", "t", 10*time.Second, "Timeout is the maximum amount of time that Wait4X will wait for a checking operation, 0 is unlimited.")
	rootCmd.PersistentFlags().BoolP("invert-check", "v", false, "Invert the sense of checking.")
	rootCmd.PersistentFlags().StringP("log-level", "l", zerolog.InfoLevel.String(), "Set the logging level (\"trace\"|\"debug\"|\"info\")")
	rootCmd.PersistentFlags().Bool("no-color", false, "If specified, output won't contain any color.")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet or silent mode. Do not show logs or error messages.")

	return rootCmd
}

// Execute run Wait4X application
func Execute() {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTCPCommand())
	rootCmd.AddCommand(dnsCmd.NewDNSCommand(rootCmd))
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
