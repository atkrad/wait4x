package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var RetryCount int
var Sleep time.Duration
var rootCmd = &cobra.Command{
	Use:   "wait4x",
	Short: "wait4x allows waiting for a port or a service to enter into specify state",
	Long: `wait4x allows waiting for a port to enter into specify state or waiting for a service e.g. redis, mysql, postgres, ... to enter inter ready state`,
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&RetryCount, "count", "c", 10, "Maximum retry counts.")
	rootCmd.PersistentFlags().DurationVarP(&Sleep, "sleep", "s", time.Second, "Sleep time between each loop.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
