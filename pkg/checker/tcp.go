package checker

import (
	"net"
	"time"

	"github.com/atkrad/wait4x/pkg/log"
)

type TCP struct {
	logger log.Logger
	address string
	timeout time.Duration
}

func NewTCP(address string, timeout time.Duration) Checker {
	t := &TCP{
		address: address,
		timeout: timeout,
	}

	return t
}

func (t *TCP) SetLogger(logger log.Logger) {
	t.logger = logger
}

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
