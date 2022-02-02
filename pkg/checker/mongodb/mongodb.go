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
	"github.com/atkrad/wait4x/pkg/checker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB represents MongoDB checker
type MongoDB struct {
	dsn string
	*checker.LogAware
}

// New creates the MongoDB checker
func New(dsn string) checker.Checker {
	i := &MongoDB{
		dsn:      dsn,
		LogAware: &checker.LogAware{},
	}

	return i
}

// Check checks MongoDB connection
func (m *MongoDB) Check(ctx context.Context) bool {
	m.Logger().Info("Checking MongoDB connection ...")

	// Create a new c and connect to the server
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(m.dsn))
	if err != nil {
		m.Logger().Debug(err)

		return false
	}

	defer func() {
		if err := c.Disconnect(ctx); err != nil {
			m.Logger().Debug(err)
		}
	}()

	// Ping the primary
	err = c.Ping(ctx, readpref.Primary())
	if err != nil {
		m.Logger().Debug(err)

		return false
	}

	return true
}
