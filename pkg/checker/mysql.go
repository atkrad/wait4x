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
