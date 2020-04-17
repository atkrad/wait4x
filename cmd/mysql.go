package cmd

import (
	"github.com/atkrad/wait4x/internal/errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"context"
	"time"
)

func NewMysqlCommand() *cobra.Command {
	mysqlCommand := &cobra.Command{
		Use:   "mysql DSN",
		Short: "Check MySQL connection.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("DSN is required argument for the mysql command")
			}

			return nil
		},
		Example: `
  # Checking MySQL TCP connection
  wait4x mysql user:password@tcp(localhost:5555)/dbname?tls=skip-verify

  # Checking MySQL UNIX Socket existence
  wait4x mysql usernname:password@unix(/tmp/mysql.sock)/myDatabase
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), Timeout)
			defer cancel()

			for !checkingMysql(cmd, args) {
				select {
				case <-ctx.Done():
					return errors.NewTimedOutError()
				case <-time.After(Interval):
				}
			}

			return nil
		},
	}

	return mysqlCommand
}

func checkingMysql(cmd *cobra.Command, args []string) bool {
	log.Info("Checking MySQL connection ...")
	db, err := sql.Open("mysql", args[0])
	if err != nil {
		log.Warn("Validating DSN data has error.")
		log.Debug(err)

		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Warn("Pinging MySQL has error.")
		log.Debug(err)

		return false
	}

	log.Info("Connection established successfully.")

	return true
}
