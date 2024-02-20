package txt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceTXT(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectTXT(t *testing.T) {
	d := New(server, WithExpectedValues([]string{"v=spf1 include:_spf.mx.cloudflare.net ~all"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectTXT(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedValues([]string{"127.0.0.1"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectTXT(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedValues([]string{"v=spf1 include:_spf.mx.cloudflare.net ~all"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestRegexCorrectTXT(t *testing.T) {
	d := New(server, WithExpectedValues([]string{".* include:.*"}))
	assert.Nil(t, d.Check(context.Background()))
}
