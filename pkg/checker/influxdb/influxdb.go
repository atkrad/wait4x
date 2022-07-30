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

package influxdb

import (
	"context"
	"github.com/atkrad/wait4x/v2/pkg/checker"
	"github.com/atkrad/wait4x/v2/pkg/checker/errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// InfluxDB represents InfluxDB checker
type InfluxDB struct {
	serverURL string
}

// New creates the InfluxDB checker
func New(serverURL string) checker.Checker {
	i := &InfluxDB{
		serverURL: serverURL,
	}

	return i
}

// Identity returns the identity of the checker
func (i InfluxDB) Identity() (string, error) {
	return i.serverURL, nil
}

// Check checks InfluxDB connection
func (i *InfluxDB) Check(ctx context.Context) error {
	// InfluxDB doesn't validate authentication params on Ping and Health requests.
	ic := influxdb2.NewClient(i.serverURL, "")
	defer ic.Close()

	res, err := ic.Ping(ctx)
	if res == false {
		return errors.Wrap(err, errors.DebugLevel)
	}

	return nil
}
