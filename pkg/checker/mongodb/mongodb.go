// Copyright 2022 Mohammad Abdolirad
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

package mongodb

import (
	"context"
	"errors"
	"github.com/atkrad/wait4x/pkg/checker"
	errs "github.com/atkrad/wait4x/pkg/checker/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
)

// MongoDB represents MongoDB checker
type MongoDB struct {
	dsn string
}

// New creates the MongoDB checker
func New(dsn string) checker.Checker {
	i := &MongoDB{
		dsn: dsn,
	}

	return i
}

// Identity returns the identity of the checker
func (m MongoDB) Identity() (string, error) {
	cops := options.Client().ApplyURI(m.dsn)
	if len(cops.Hosts) == 0 {
		return "", errors.New("can't retrieve the checker identity")
	}

	return strings.Join(cops.Hosts, ","), nil
}

// Check checks MongoDB connection
func (m *MongoDB) Check(ctx context.Context) (err error) {
	// Creates a new Client and then initializes it using the Connect method.
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(m.dsn))
	if err != nil {
		return errs.Wrap(err, errs.DebugLevel)
	}

	defer func() {
		if err := c.Disconnect(ctx); err != nil {
			err = errs.Wrap(err, errs.DebugLevel)
		}
	}()

	// Ping the primary
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		return errs.Wrap(err, errs.DebugLevel)
	}

	return nil
}
