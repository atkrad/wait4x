package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCNAMECommand(t *testing.T) {
	cmd := NewCNAMECommand()

	assert.Equal(t, "CNAME ADDRESS [--command [args...]]", cmd.Use)
	assert.Equal(t, []string{"cname"}, cmd.Aliases)
	assert.Equal(t, "Check DNS CNAME records", cmd.Short)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)
	assert.Equal(t, "ADDRESS is required argument for the dns command", err.Error())

	err = cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)

	flags := cmd.Flags()
	expectDomain, err := flags.GetStringArray("expect-domain")
	assert.NoError(t, err)
	assert.Empty(t, expectDomain)
}

func TestRunCNAME(t *testing.T) {
	cmd := NewCNAMECommand()

	cmd.Flags().Duration("interval", 0, "")
	cmd.Flags().Duration("timeout", 0, "")
	cmd.Flags().Bool("invert-check", false, "")
	cmd.Flags().String("nameserver", "", "")

	err := cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)
}
