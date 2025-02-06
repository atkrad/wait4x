package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewACommand(t *testing.T) {
	cmd := NewACommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "A ADDRESS [value] [--command [args...]]", cmd.Use)
	assert.Equal(t, []string{"a"}, cmd.Aliases)
}

func TestACommand_NoArgs(t *testing.T) {
	cmd := NewACommand()
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)
	assert.Equal(t, "ADDRESS is required argument for the dns command", err.Error())
}

func TestACommand_WithArgs(t *testing.T) {
	cmd := NewACommand()
	err := cmd.Args(cmd, []string{"example.com"})
	assert.NoError(t, err)
}

func TestRunA(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		flags map[string]string
	}{
		{
			name: "basic check",
			args: []string{"example.com"},
		},
		{
			name: "with expected IP",
			args: []string{"example.com"},
			flags: map[string]string{
				"expect-ip": "93.184.216.34",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewACommand()

			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				assert.NoError(t, err)
			}

			err := cmd.Args(cmd, tt.args)
			assert.NoError(t, err)
		})
	}
}
