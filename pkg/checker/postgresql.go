package checker

import (
	"database/sql"

	"github.com/atkrad/wait4x/pkg/log"
	// Needed for the PostgreSQL driver
	_ "github.com/lib/pq"
)

// PostgreSQL represents PostgreSQL checker
type PostgreSQL struct {
	dsn    string
	logger log.Logger
}

// NewPostgreSQL creates the PostgreSQL checker
func NewPostgreSQL(dsn string) Checker {
	p := &PostgreSQL{
		dsn: dsn,
	}

	return p
}

// SetLogger sets default logger
func (p *PostgreSQL) SetLogger(logger log.Logger) {
	p.logger = logger
}

// Check checks PostgreSQL connection
func (p *PostgreSQL) Check() bool {
	p.logger.Info("Checking PostgreSQL connection ...")
	db, err := sql.Open("postgres", p.dsn)
	if err != nil {
		p.logger.Debug(err)

		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		p.logger.Debug(err)

		return false
	}

	return true
}
