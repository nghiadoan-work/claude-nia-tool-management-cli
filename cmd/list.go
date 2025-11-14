package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// List flags
	listRemote   bool
	listType     string
	listTags     []string
	listAuthor   string
	listSortBy   string
	listSortDesc bool
	listLimit    int
	listJSON     bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tools from the registry",
	Long: `List tools from the remote registry with optional filtering and sorting.

By default, this command lists all available tools in the registry.
You can filter by type, tags, author, and apply sorting.

Examples:
  cntm list --remote                      # List all remote tools
  cntm list --remote --type agent         # List only agents
  cntm list --remote --tag git            # List tools with "git" tag
  cntm list --remote --sort-by downloads  # Sort by download count
  cntm list --remote --sort-desc          # Sort in descending order
  cntm list --remote --limit 10           # Limit to 10 results
  cntm list --remote --json               # Output in JSON format`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// List flags
	listCmd.Flags().BoolVar(&listRemote, "remote", false, "list remote tools from registry")
	listCmd.Flags().StringVarP(&listType, "type", "t", "", "filter by tool type (agent, command, skill)")
	listCmd.Flags().StringSliceVar(&listTags, "tag", []string{}, "filter by tags (can specify multiple)")
	listCmd.Flags().StringVarP(&listAuthor, "author", "a", "", "filter by author")
	listCmd.Flags().StringVar(&listSortBy, "sort-by", "name", "sort by field (name, created, updated, downloads)")
	listCmd.Flags().BoolVar(&listSortDesc, "sort-desc", false, "sort in descending order")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 0, "limit number of results (0 for all)")
	listCmd.Flags().BoolVarP(&listJSON, "json", "j", false, "output in JSON format")
}

func runList(cmd *cobra.Command, args []string) error {
	if !listRemote {
		return fmt.Errorf("currently only --remote listing is supported\nHint: Use 'cntm list --remote' to list tools from the registry")
	}

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

	// Build list filter
	filter := &models.ListFilter{
		Tags:     listTags,
		Author:   listAuthor,
		SortDesc: listSortDesc,
		Limit:    listLimit,
	}

	// Parse tool type if provided
	if listType != "" {
		filter.Type = models.ToolType(listType)
	}

	// Parse sort field
	switch listSortBy {
	case "name":
		filter.SortBy = models.SortByName
	case "created":
		filter.SortBy = models.SortByCreated
	case "updated":
		filter.SortBy = models.SortByUpdated
	case "downloads":
		filter.SortBy = models.SortByDownloads
	case "":
		filter.SortBy = models.SortByName
	default:
		return fmt.Errorf("invalid sort field: %s (must be: name, created, updated, downloads)", listSortBy)
	}

	// Validate filter
	if err := filter.Validate(); err != nil {
		return fmt.Errorf("invalid list filter: %w", err)
	}

	// Show progress message
	if !listJSON && verbose {
		fmt.Fprintln(os.Stderr, "Fetching registry...")
	}

	// List tools
	results, err := registryService.ListTools(filter)
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}

	// Display results
	if listJSON {
		return outputJSON(results)
	}

	return displayToolsTable(results)
}
