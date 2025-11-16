package cmd

import (
	"os"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/version"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	cfgFile  string
	verbose  bool
	basePath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cntm",
	Short: "Claude Nia Tool Management CLI - Package manager for Claude Code tools",
	Long: `cntm is a package manager for Claude Code tools (agents, commands, and skills).

Like npm for Node.js, cntm helps you:
  - Initialize tool configuration
  - Search and discover available tools
  - Install tools from a GitHub registry
  - Update tools to the latest versions
  - Publish your own tools to share with others
  - Remove installed tools

Available commands:
  cntm init                     # Initialize tool configuration
  cntm search "code review"     # Search for tools
  cntm install code-reviewer    # Install a tool
  cntm update --all             # Update all tools
  cntm publish my-agent         # Publish your tool
  cntm remove code-reviewer     # Remove an installed tool`,
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.claude-tools-config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&basePath, "path", "p", ".claude", "path to .claude directory")

	// Local flags
	rootCmd.Flags().BoolP("version", "", false, "version for cntm")
}
