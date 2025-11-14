package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutdatedCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "with json flag",
			args: []string{"--json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := outdatedCmd
			assert.NotNil(t, cmd)
			assert.Equal(t, "outdated", cmd.Use)
			assert.NotEmpty(t, cmd.Short)
			assert.NotEmpty(t, cmd.Long)
			assert.NotNil(t, cmd.RunE)

			// Check flags exist
			jsonFlag := cmd.Flags().Lookup("json")
			assert.NotNil(t, jsonFlag)
		})
	}
}
