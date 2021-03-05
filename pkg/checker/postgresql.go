// Copyright 2020 Mohammad Abdolirad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checker

import (
	"database/sql"

	// Needed for the PostgreSQL driver
	_ "github.com/lib/pq"
)

// PostgreSQL represents PostgreSQL checker
type PostgreSQL struct {
	dsn string
	*LogAware
}

// NewPostgreSQL creates the PostgreSQL checker
func NewPostgreSQL(dsn string) Checker {
	p := &PostgreSQL{
		dsn:      dsn,
		LogAware: &LogAware{},
	}

	return p
}

// Check checks PostgreSQL connection
func (p *PostgreSQL) Check() bool {
	p.logger.Info("Checking PostgreSQL connection ...")
	db, err := sql.Open("postgres", p.dsn)
	if err != nil {
		p.logger.Debug(err)

		return false
	}

	defer func() {
		if err := db.Close(); err != nil {
			p.logger.Debug(err)
		}
	}()

	err = db.Ping()
	if err != nil {
		p.logger.Debug(err)

		return false
	}

	return true
}
