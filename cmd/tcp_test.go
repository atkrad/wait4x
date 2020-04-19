package cmd

import (
	"github.com/atkrad/wait4x/internal/errors"
	"github.com/atkrad/wait4x/internal/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)

	os.Exit(m.Run())
}

func TestTcpCommandInvalidArgument(t *testing.T) {
	wait4xCommand := NewWait4X()
	wait4xCommand.AddCommand(NewTcpCommand())

	_, err := test.ExecuteCommand(wait4xCommand, "tcp")

	assert.Equal(t, "ADDRESS is required argument for the tcp command", err.Error())
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
