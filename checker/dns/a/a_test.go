package a

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceA(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectA(t *testing.T) {
	d := New(server, WithExpectedIPV4s([]string{"172.67.154.180", "127.0.0.1"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectA(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedIPV4s([]string{"127.0.0.1"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectA(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedIPV4s([]string{"172.67.154.180"}))
	assert.Nil(t, d.Check(context.Background()))
}
