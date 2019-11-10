package cmd

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"os"
	"time"
)

var tcpCmd = &cobra.Command{
	Use:   "tcp ADDRESS",
	Short: "Check TCP connection.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("ADDRESS is required argument for the tcp command")
		}

		return nil
	},
	Example: `
  # If you want checking just tcp connection 
  wait4x tcp 127.0.0.1:9090
`,
	Run: func(cmd *cobra.Command, args []string) {
		timeout, _ := cmd.Flags().GetDuration("timeout")

		var i = 1
		for i <= RetryCount {
			log.Info("Checking tcp connection")

			d := net.Dialer{Timeout: timeout}
			_, err := d.Dial("tcp", args[0])
			if err != nil {
				log.Debug(err)

				time.Sleep(Sleep)
				i += 1
				continue
			} else {
				os.Exit(0)
			}
		}

		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)
	tcpCmd.Flags().Duration("timeout", time.Second*10, "Timeout is the maximum amount of time a dial will wait for a connect to complete.")
}
