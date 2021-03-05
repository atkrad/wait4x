package checker

import (
	"github.com/atkrad/wait4x/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestTcpValidAddress(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	go func() {
		defer ln.Close()
		_, _ = ln.Accept()
	}()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	tc := NewTCP(ln.Addr().String(), time.Second*5)
	tc.SetLogger(logger)

	assert.Equal(t, true, tc.Check())
}

func TestTcpInvalidAddress(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	go func() {
		defer ln.Close()
		_, _ = ln.Accept()
	}()

	logger, _ := log.NewLogrus(logrus.DebugLevel.String(), ioutil.Discard)

	tc := NewTCP(ln.Addr().String()+"0", time.Second*5)
	tc.SetLogger(logger)

	assert.Equal(t, false, tc.Check())
}
