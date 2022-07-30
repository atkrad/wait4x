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

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/atkrad/wait4x/v2/pkg/checker"
	"github.com/atkrad/wait4x/v2/pkg/checker/errors"
	"net/url"

	// Needed for the PostgreSQL driver
	_ "github.com/lib/pq"
)

// PostgreSQL represents PostgreSQL checker
type PostgreSQL struct {
	dsn string
}

// New creates the PostgreSQL checker
func New(dsn string) checker.Checker {
	p := &PostgreSQL{
		dsn: dsn,
	}

	return p
}

// Identity returns the identity of the checker
func (p PostgreSQL) Identity() (string, error) {
	u, err := url.Parse(p.dsn)
	if err != nil {
		return "", fmt.Errorf("can't retrieve the checker identity: %w", err)
	}

	return u.Host, nil
}

// Check checks PostgreSQL connection
func (p *PostgreSQL) Check(ctx context.Context) (err error) {
	db, err := sql.Open("postgres", p.dsn)
	if err != nil {
		return errors.Wrap(err, errors.DebugLevel)
	}

	defer func() {
		if err := db.Close(); err != nil {
			err = errors.Wrap(err, errors.DebugLevel)
		}
	}()

	err = db.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, errors.DebugLevel)
	}

	return nil
}
