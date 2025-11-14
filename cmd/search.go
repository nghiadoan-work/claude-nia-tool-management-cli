package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Search flags
	searchType          string
	searchTags          []string
	searchAuthor        string
	searchMinDownloads  int
	searchRegex         bool
	searchCaseSensitive bool
	searchJSON          bool
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for tools in the registry",
	Long: `Search for tools in the remote registry by name, description, tags, or author.

The search query will match against:
  - Tool name
  - Tool description
  - Tool tags
  - Tool author

Examples:
  cntm search "code review"           # Search for code review tools
  cntm search git --type agent        # Search for git agents
  cntm search test --tag testing      # Search tools with "testing" tag
  cntm search --author john           # Search tools by author "john"
  cntm search "^code" --regex         # Search using regex pattern
  cntm search tool --json             # Output in JSON format`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Search flags
	searchCmd.Flags().StringVarP(&searchType, "type", "t", "", "filter by tool type (agent, command, skill)")
	searchCmd.Flags().StringSliceVar(&searchTags, "tag", []string{}, "filter by tags (can specify multiple)")
	searchCmd.Flags().StringVarP(&searchAuthor, "author", "a", "", "filter by author")
	searchCmd.Flags().IntVar(&searchMinDownloads, "min-downloads", 0, "filter by minimum downloads")
	searchCmd.Flags().BoolVarP(&searchRegex, "regex", "r", false, "use regex for pattern matching")
	searchCmd.Flags().BoolVar(&searchCaseSensitive, "case-sensitive", false, "case-sensitive search")
	searchCmd.Flags().BoolVarP(&searchJSON, "json", "j", false, "output in JSON format")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Load config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse GitHub URL to get owner and repo
	owner, repo, err := parseGitHubURL(cfg.Registry.URL)
	if err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}

	// Initialize services
	githubClient := services.NewGitHubClient(services.GitHubClientConfig{
		Owner:     owner,
		Repo:      repo,
		Branch:    cfg.Registry.Branch,
		AuthToken: cfg.Registry.AuthToken,
	})

	cacheManager, err := data.NewCacheManager(basePath, 3600*time.Second) // 1 hour TTL
	if err != nil {
		return fmt.Errorf("failed to create cache manager: %w", err)
	}
	registryService := services.NewRegistryService(githubClient, cacheManager)

	// Build search filter
	filter := &models.SearchFilter{
		Query:         query,
		Tags:          searchTags,
		Author:        searchAuthor,
		MinDownloads:  searchMinDownloads,
		Regex:         searchRegex,
		CaseSensitive: searchCaseSensitive,
	}

	// Parse tool type if provided
	if searchType != "" {
		filter.Type = models.ToolType(searchType)
	}

	// Validate filter
	if err := filter.Validate(); err != nil {
		return fmt.Errorf("invalid search filter: %w", err)
	}

	// Show progress message
	if !searchJSON && verbose {
		fmt.Fprintln(os.Stderr, "Searching registry...")
	}

	// Search tools
	results, err := registryService.SearchTools(filter)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Display results
	if searchJSON {
		return outputJSON(results)
	}

	return displayToolsTable(results)
}

func displayToolsTable(tools []*models.ToolInfo) error {
	if len(tools) == 0 {
		fmt.Println("No tools found matching your search criteria.")
		return nil
	}

	// Prepare table data
	headers := []string{"Name", "Type", "Version", "Author", "Downloads", "Description"}
	var rows [][]string

	for _, tool := range tools {
		description := tool.Description
		if len(description) > 80 {
			description = description[:77] + "..."
		}

		rows = append(rows, []string{
			tool.Name,
			string(tool.Type),
			tool.Version,
			tool.Author,
			fmt.Sprintf("%d", tool.Downloads),
			description,
		})
	}

	// Create table with new API
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(headers),
	)

	// Add rows
	for _, row := range rows {
		table.Append(row)
	}

	// Render table
	table.Render()
	fmt.Printf("\nFound %d tool(s)\n", len(tools))

	return nil
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
