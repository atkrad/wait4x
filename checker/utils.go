// Copyright 2022 The Wait4X Authors
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
	"net"
	"net/url"
	"syscall"
)

// IsConnectionRefused attempts to determine if the given error was caused by a failure to establish a connection.
func IsConnectionRefused(err error) bool {
	switch t := err.(type) {
	case *url.Error:
		return IsConnectionRefused(t.Err)
	case *net.OpError:
		if t.Op == "dial" || t.Op == "read" {
			return true
		}
		return IsConnectionRefused(t.Err)
	case syscall.Errno:
		return t == syscall.ECONNREFUSED
	}

	return false
}
