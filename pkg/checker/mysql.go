package checker

import (
	"database/sql"

	"github.com/atkrad/wait4x/pkg/log"
	// Needed for the MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dsn    string
	logger log.Logger
}

// NewMysqlCommand creates the mysql sub-command
func NewMysql(dsn string) Checker {
	m := &MySQL{
		dsn: dsn,
	}

	return m
}

func (m *MySQL) SetLogger(logger log.Logger) {
	m.logger = logger
}

func (m *MySQL) Check() bool {
	m.logger.Info("Checking MySQL connection ...")
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		m.logger.Debug(err)

		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		m.logger.Debug(err)

		return false
	}

	m.logger.Info("Connection established successfully.")

	return true
}
