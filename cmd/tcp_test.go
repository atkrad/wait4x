package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/atkrad/wait4x/internal/test"
	"github.com/atkrad/wait4x/internal/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"io/ioutil"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	os.Exit(m.Run())
}

func TestTcpConnectionSuccess(t *testing.T) {
	wait4xCommand := NewWait4X()
	wait4xCommand.AddCommand(NewTcpCommand())

	_, err := test.ExecuteCommand(wait4xCommand, "tcp", "1.1.1.1:53")

	assert.Nil(t, err)
}

func TestTcpConnectionFail(t *testing.T) {
	wait4xCommand := NewWait4X()
	wait4xCommand.AddCommand(NewTcpCommand())

	_, err := test.ExecuteCommand(wait4xCommand, "tcp", "127.0.0.1:8080", "-t", "2s")

	assert.Equal(t, errors.TIMED_OUT_ERROR_MESSAGE, err.Error())
}
