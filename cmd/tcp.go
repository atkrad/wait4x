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
	Run: func (cmd *cobra.Command, args []string) {
		ticker := time.NewTicker(Interval)
		defer ticker.Stop()

		go func() {
			connectionTimeout, _ := cmd.Flags().GetDuration("connection-timeout")
			for ; true; <-ticker.C {
				log.Info("Checking TCP connection ...")

				d := net.Dialer{Timeout: connectionTimeout}
				_, err := d.Dial("tcp", args[0])
				if err != nil {
					log.Debug(err)

					continue
				} else {
					os.Exit(EXIT_SUCCESS)
				}
			}
		}()

		time.Sleep(Timeout)
		log.Info("Operation Timed Out")

		os.Exit(EXIT_TIMEDOUT)
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)

	tcpCmd.Flags().Duration("connection-timeout", time.Second*5, "Timeout is the maximum amount of time a dial will wait for a connection to complete.")
}
