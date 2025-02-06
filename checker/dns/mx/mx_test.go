package mx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceMX(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectMX(t *testing.T) {
	d := New(server, WithExpectedDomains([]string{"route1.mx.cloudflare.net", "route2.mx.cloudflare.net"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectMX(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedDomains([]string{"127.0.0.1"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectA(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedDomains([]string{"route1.mx.cloudflare.net"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestRegexCorrectA(t *testing.T) {
	d := New(server, WithExpectedDomains([]string{".*.mx.cloudflare.net"}))
	assert.Nil(t, d.Check(context.Background()))
}
