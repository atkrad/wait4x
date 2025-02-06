package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTXTCommand(t *testing.T) {
	cmd := NewTXTCommand()

	assert.Equal(t, "TXT ADDRESS [--command [args...]]", cmd.Use)
	assert.Equal(t, []string{"txt"}, cmd.Aliases)
	assert.Equal(t, "Check DNS TXT records", cmd.Short)

	err := cmd.Args(cmd, []string{})
	assert.EqualError(t, err, "ADDRESS is required argument for the dns command")

	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	flags := cmd.Flags()
	expectValue, err := flags.GetStringArray("expect-value")
	assert.NoError(t, err)
	assert.Empty(t, expectValue)
}
