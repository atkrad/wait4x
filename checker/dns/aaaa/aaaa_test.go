package aaaa

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceAAAA(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectAAAA(t *testing.T) {
	d := New(server, WithExpectedIPV6s([]string{"2606:4700:3034::6815:591"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectAAAA(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedIPV6s([]string{"127.0.0.1"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectAAAA(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedIPV6s([]string{"2606:4700:3034::6815:591"}))
	assert.Nil(t, d.Check(context.Background()))
}
