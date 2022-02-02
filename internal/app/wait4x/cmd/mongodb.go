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

package cmd

import (
	"github.com/atkrad/wait4x/internal/pkg/errors"
	"github.com/atkrad/wait4x/pkg/checker/mongodb"
	"github.com/atkrad/wait4x/pkg/waiter"
	"github.com/spf13/cobra"
)

// NewMongoDBCommand creates the mongodb sub-command
func NewMongoDBCommand() *cobra.Command {
	mongodbCommand := &cobra.Command{
		Use:   "mongodb DSN",
		Short: "Check MongoDB connection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.NewCommandError("DSN is required argument for the mongodb command")
			}

			return nil
		},
		Example: `
  # Checking MongoDB connection
  wait4x mongodb 'mongodb://127.0.0.1:27017'

  # Checking MongoDB connection with credentials and options
  wait4x mongodb 'mongodb://user:pass@127.0.0.1:27017/?maxPoolSize=20&w=majority'
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			invertCheck, _ := cmd.Flags().GetBool("invert-check")

			mc := mongodb.New(args[0])
			mc.SetLogger(Logger)

			return waiter.Wait(
				mc.Check,
				waiter.WithTimeout(timeout),
				waiter.WithInterval(interval),
				waiter.WithInvertCheck(invertCheck),
			)
		},
	}

	return mongodbCommand
}
