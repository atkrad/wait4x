package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAAAACommand(t *testing.T) {
	cmd := NewAAAACommand()

	assert.Equal(t, "AAAA ADDRESS [--command [args...]]", cmd.Use)
	assert.Equal(t, []string{"aaaa"}, cmd.Aliases)
	assert.Equal(t, "Check DNS AAAA records", cmd.Short)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)
	assert.Equal(t, "ADDRESS is required argument for the dns command", err.Error())

	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	flags := cmd.Flags()
	expectIP, err := flags.GetStringArray("expect-ip")
	assert.NoError(t, err)
	assert.Empty(t, expectIP)
}

func TestRunAAAA(t *testing.T) {
	cmd := NewAAAACommand()

	err := cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	cmd.Flags().Set("expect-ip", "2606:4700:3033::ac43:9ab4")
	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	cmd.Flags().Set("nameserver", "8.8.8.8")
	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	cmd.Flags().Set("interval", "1s")
	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	cmd.Flags().Set("timeout", "5s")
	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	cmd.Flags().Set("invert-check", "true")
	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)
}
