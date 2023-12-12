// Copyright 2020 The Wait4X Authors
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

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"wait4x.dev/v2/checker"
	// Needed for the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"regexp"
)

var hidePasswordRegexp = regexp.MustCompile(`^([^:]+):[^:@]+@`)

// MySQL represents MySQL checker
type MySQL struct {
	dsn string
}

// New creates the MySQL checker
func New(dsn string) checker.Checker {
	m := &MySQL{
		dsn: dsn,
	}

	return m
}

// Identity returns the identity of the checker
func (m *MySQL) Identity() (string, error) {
	cfg, err := mysql.ParseDSN(m.dsn)
	if err != nil {
		return "", fmt.Errorf("can't retrieve the checker identity: %w", err)
	}

	return cfg.Addr, nil
}

// Check checks MySQL connection
func (m *MySQL) Check(ctx context.Context) (err error) {
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		if dberr := db.Close(); dberr != nil {
			err = dberr
		}
	}(db)

	err = db.PingContext(ctx)
	if err != nil {
		if checker.IsConnectionRefused(err) {
			return checker.NewExpectedError(
				"failed to establish a connection to the mysql server", err,
				"dsn", hidePasswordRegexp.ReplaceAllString(m.dsn, `$1:***@`),
			)
		}

		return err
	}

	return nil
}
