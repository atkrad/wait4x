// Copyright 2019 Mohammad Abdolirad
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

package dns

import (
	"errors"

	"github.com/spf13/cobra"
)

func NewDNSCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "dns COMMAND [flags] [--command [args...]]",
		Short: "Check DNS records",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("ADDRESS is required argument for the dns command")
			}

			return nil
		},
	}

	command.Flags().String("nameserver", "", " Address of the nameserver to send packets to")

	command.AddCommand(NewDNSACommand())
	command.AddCommand(NewDNSAAAACommand())
	command.AddCommand(NewDNSCNAMECommand())
	command.AddCommand(NewDNSMXCommand())
	command.AddCommand(NewDNSTXTCommand())
	command.AddCommand(NewDNSNSCommand())

	return command
}
