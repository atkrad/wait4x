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
