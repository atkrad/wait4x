package cname

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceCNAME(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectCNAME(t *testing.T) {
	d := New(server, WithExpectedDomains([]string{"wait4x.dev"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectCNAME(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedDomains([]string{"something wrong"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectCNAME(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedDomains([]string{"wait4x.dev"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestRegexCorrectCNAME(t *testing.T) {
	d := New(server, WithExpectedDomains([]string{".*wait4.*"}))
	assert.Nil(t, d.Check(context.Background()))
}
