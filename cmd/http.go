package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

var expectResponseCode int

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Check HTTP connection.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		expectResponseCode, _ := cmd.Flags().GetString("expect-response-code")
		expectResponseBody, _ := cmd.Flags().GetString("expect-response-body")

		var i = 1
		for i <= RetryCount {
			fmt.Print(".")

			resp, err := http.Get(address)

			if err != nil {
				time.Sleep(Sleep)
				i += 1
				continue
			} else {
				if (expectResponseCode != "" && resp.Status == expectResponseCode) ||
					(expectResponseBody != "" && resp. == expectResponseCode) {

				}
				defer resp.Body.Close()
				os.Exit(0)
			}
		}

		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(tcpCmd)
	httpCmd.Flags().String("address", "http://127.0.0.1", "Http address.")
	httpCmd.Flags().String("expect-response-code", "", "Expect response code e.g. 200, 204, ... .")
	httpCmd.Flags().String("expect-response-body", "", "Expect response body.")
}
