package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/ui"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/version"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	versionOutput string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for cntm, including:
  - Semantic version
  - Git commit hash
  - Build date
  - Go version used to build

Examples:
  cntm version              # Show version
  cntm version --output json   # JSON output
  cntm version --output yaml   # YAML output`,
	RunE: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionOutput, "output", "o", "text", "Output format: text, json, yaml")
}

func runVersion(cmd *cobra.Command, args []string) error {
	info := version.GetInfo()

	switch versionOutput {
	case "json":
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}
		fmt.Println(string(data))

	case "yaml":
		data, err := yaml.Marshal(info)
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}
		fmt.Print(string(data))

	case "text":
		fallthrough
	default:
		ui.PrintSuccess(info.LongString())
	}

	return nil
}
