package checker

import (
	"net"
	"time"

	"github.com/atkrad/wait4x/pkg/log"
)

// TCP represents TCP checker
type TCP struct {
	logger  log.Logger
	address string
	timeout time.Duration
}

// NewTCP creates the TCP checker
func NewTCP(address string, timeout time.Duration) Checker {
	t := &TCP{
		address: address,
		timeout: timeout,
	}

	return t
}

// SetLogger sets default logger
func (t *TCP) SetLogger(logger log.Logger) {
	t.logger = logger
}

// Check checks TCP connection
func (t *TCP) Check() bool {
	d := net.Dialer{Timeout: t.timeout}

	if t.logger != nil {
		t.logger.Info("Checking TCP connection ...")
	}
	_, err := d.Dial("tcp", t.address)
	if err != nil {
		if t.logger != nil {
			t.logger.Debug(err)
		}

		return false
	}

	return true
}
