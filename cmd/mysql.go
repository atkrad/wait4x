package cmd

import (
	"errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql DSN",
	Short: "Check MySQL connection.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("DSN is required argument for the mysql command")
		}

		return nil
	},
	Example: `
  # Checking MySQL TCP connection
  wait4x mysql user:password@tcp(localhost:5555)/dbname?tls=skip-verify

  # Checking MySQL UNIX Socket existence
  wait4x mysqli usernname:password@unix(/tmp/mysql.sock)/myDatabase
`,
	Run: func(cmd *cobra.Command, args []string) {
		ticker := time.NewTicker(Interval)
		defer ticker.Stop()

		go func() {
			for ; true; <-ticker.C {
				log.Info("Checking MySQL connection ...")
				db, err := sql.Open("mysql", args[0])
				if err != nil {
					log.Warn("Validating DSN data has error.")
					log.Debug(err)

					continue
				}
				defer db.Close()

				err = db.Ping()
				if err != nil {
					log.Warn("Pinging MySQL has error.")
					log.Debug(err)

					continue
				} else {
					log.Info("Connection established successfully.")
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
	rootCmd.AddCommand(mysqlCmd)
}
