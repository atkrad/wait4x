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

package mysql

import (
	"context"
	"database/sql"
	"github.com/atkrad/wait4x/pkg/checker"

	// Needed for the MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// MySQL represents MySQL checker
type MySQL struct {
	dsn string
	*checker.LogAware
}

// New creates the MySQL checker
func New(dsn string) checker.Checker {
	m := &MySQL{
		dsn:      dsn,
		LogAware: &checker.LogAware{},
	}

	return m
}

// Check checks MySQL connection
func (m *MySQL) Check(ctx context.Context) bool {
	m.Logger().Info("Checking MySQL connection ...")
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		m.Logger().Debug(err)

		return false
	}

	defer func() {
		if err := db.Close(); err != nil {
			m.Logger().Debug(err)
		}
	}()

	err = db.PingContext(ctx)
	if err != nil {
		m.Logger().Debug(err)

		return false
	}

	m.Logger().Info("Connection established successfully.")

	return true
}
