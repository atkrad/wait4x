package cmd

import (
	"context"
	"net/url"
	"time"

	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/pkg/checker"
	"github.com/spf13/cobra"
)

// NewHTTPCommand creates the http sub-command
func NewHTTPCommand() *cobra.Command {
	httpCommand := &cobra.Command{
		Use:   "http ADDRESS",
		Short: "Check HTTP connection.",
		Long:  "",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("ADDRESS is required argument for the http command")
			}

			_, err := url.Parse(args[0])
			if err != nil {
				return err
			}

			return nil
		},
		Example: `
  # If you want checking just http connection
  wait4x http http://ifconfig.co

  # If you want checking http connection and expect specify http status code
  wait4x http http://ifconfig.co --expect-status-code 200
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			expectStatusCode, _ := cmd.Flags().GetInt("expect-status-code")
			expectBody, _ := cmd.Flags().GetString("expect-body")
			connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")

			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			hc := checker.NewHTTP(args[0], expectStatusCode, expectBody, connectionTimeout)
			hc.SetLogger(Logger)

			for !hc.Check() {
				select {
				case <-ctx.Done():
					return errors.NewTimedOutError()
				case <-time.After(Interval):
				}
			}

			return nil
		},
	}

	httpCommand.Flags().Int("expect-status-code", 0, "Expect response code e.g. 200, 204, ... .")
	httpCommand.Flags().String("expect-body", "", "Expect response body pattern.")
	httpCommand.Flags().Duration("connection-timeout", time.Second*5, "Http connection timeout, The timeout includes connection time, any redirects, and reading the response body.")

	return httpCommand
}
