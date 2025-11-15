package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveCommand_Validation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "single tool",
			args:    []string{"code-reviewer"},
			wantErr: false,
		},
		{
			name:    "multiple tools",
			args:    []string{"tool1", "tool2", "tool3"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test argument validation
			err := removeCmd.Args(removeCmd, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveCommand_Aliases(t *testing.T) {
	// Test that the command has the expected aliases
	assert.Contains(t, removeCmd.Aliases, "uninstall")
	assert.Contains(t, removeCmd.Aliases, "rm")
}

func TestRemoveCommand_Flags(t *testing.T) {
	// Test that the --yes flag exists
	flag := removeCmd.Flags().Lookup("yes")
	assert.NotNil(t, flag)
	assert.Equal(t, "y", flag.Shorthand)

	// Test that the -y shorthand exists
	flag = removeCmd.Flags().ShorthandLookup("y")
	assert.NotNil(t, flag)
	assert.Equal(t, "yes", flag.Name)
}
