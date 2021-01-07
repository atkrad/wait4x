package checker

import (
	"database/sql"

	"github.com/atkrad/wait4x/pkg/log"
	// Needed for the MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// MySQL represents MySQL checker
type MySQL struct {
	dsn    string
	logger log.Logger
}

// NewMySQL creates the MySQL checker
func NewMySQL(dsn string) Checker {
	m := &MySQL{
		dsn: dsn,
	}

	return m
}

// SetLogger sets default logger
func (m *MySQL) SetLogger(logger log.Logger) {
	m.logger = logger
}

// Check checks MySQL connection
func (m *MySQL) Check() bool {
	m.logger.Info("Checking MySQL connection ...")
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		m.logger.Debug(err)

		return false
	}

	defer func() {
		if err := db.Close(); err != nil {
			m.logger.Debug(err)
		}
	}()

	err = db.Ping()
	if err != nil {
		m.logger.Debug(err)

		return false
	}

	m.logger.Info("Connection established successfully.")

	return true
}
