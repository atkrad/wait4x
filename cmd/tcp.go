package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"
	"time"
)

var host string
var port string

var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Check TCP connection.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var i = 1
		for i <= RetryCount {
			fmt.Print(".")

			_, err := net.Dial("tcp", host+":"+port)
			if err != nil {
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
	tcpCmd.Flags().StringVar(&host, "host", "127.0.0.1", "TCP host.")
	tcpCmd.Flags().StringVarP(&port, "port", "p", "", "TCP port.")
}
