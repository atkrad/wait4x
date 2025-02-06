package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNSCommand(t *testing.T) {
	cmd := NewNSCommand()

	assert.Equal(t, "NS ADDRESS [--command [args...]]", cmd.Use)
	assert.Equal(t, []string{"ns"}, cmd.Aliases)
	assert.Equal(t, "Check DNS NS records", cmd.Short)

	err := cmd.Args(cmd, []string{})
	assert.EqualError(t, err, "ADDRESS is required argument for the dns command")

	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	flags := cmd.Flags()
	expectNameserver, err := flags.GetStringArray("expect-nameserver")
	assert.NoError(t, err)
	assert.Empty(t, expectNameserver)
}
