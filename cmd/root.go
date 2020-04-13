package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	Interval time.Duration
	Timeout time.Duration
	rootCmd = &cobra.Command{
		Use:   "wait4x",
		Short: "wait4x allows waiting for a port or a service to enter into specify state",
		Long: `wait4x allows waiting for a port to enter into specify state or waiting for a service e.g. redis, mysql, postgres, ... to enter inter ready state`,
	}
)

const EXIT_SUCCESS = 0
const EXIT_ERROR = 1
const EXIT_TIMEDOUT = 124

func init() {
	rootCmd.PersistentFlags().DurationVarP(&Interval, "interval", "i", 1 * time.Second, "Interval time between each loop.")
	rootCmd.PersistentFlags().DurationVarP(&Timeout, "timeout", "t", 10 * time.Second, "Timeout is the maximum amount of time that Wait4X will wait for a checking operation.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(EXIT_ERROR)
	}
}
