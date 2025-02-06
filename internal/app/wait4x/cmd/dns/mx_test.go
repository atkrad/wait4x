package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMXCommand(t *testing.T) {
	cmd := NewMXCommand()

	assert.Equal(t, "MX ADDRESS [--command [args...]]", cmd.Use)
	assert.Equal(t, []string{"mx"}, cmd.Aliases)
	assert.Equal(t, "Check DNS MX records", cmd.Short)

	err := cmd.Args(cmd, []string{})
	assert.EqualError(t, err, "ADDRESS is required argument for the dns command")

	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	flags := cmd.Flags()
	expectDomain, err := flags.GetStringArray("expect-domain")
	assert.NoError(t, err)
	assert.Empty(t, expectDomain)
}
