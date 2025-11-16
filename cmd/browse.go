package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Browse flags
	browseType     string
	browseTags     []string
	browseSortBy   string
	browseSortDesc bool
	browseLimit    int
	browseJSON     bool
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse and explore tools in the registry",
	Long: `Browse all available tools in the registry with sorting and filtering options.

This command provides an interactive way to discover tools with support for:
  - Sorting by downloads, updated date, created date, or name
  - Filtering by type (agent, command, skill)
  - Filtering by tags
  - Limiting results

Examples:
  cntm browse                        # Browse all tools
  cntm browse --sort downloads       # Sort by most downloaded
  cntm browse --sort downloads --desc # Sort by downloads descending
  cntm browse --sort updated         # Sort by recently updated
  cntm browse --type agent           # Browse only agents
  cntm browse --tag git              # Browse tools with "git" tag
  cntm browse --limit 10             # Show top 10 tools
  cntm browse --json                 # Output in JSON format`,
	RunE: runBrowse,
}

// trendingCmd represents the trending command (alias for browse --sort downloads --desc)
var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Show trending tools (most downloaded)",
	Long: `Show trending tools sorted by download count.

This is a convenience command equivalent to:
  cntm browse --sort downloads --desc

Examples:
  cntm trending              # Show all trending tools
  cntm trending --limit 10   # Show top 10 trending tools
  cntm trending --type agent # Show trending agents only`,
	RunE: runTrending,
}

func init() {
	rootCmd.AddCommand(browseCmd)
	rootCmd.AddCommand(trendingCmd)

	// Browse flags
	browseCmd.Flags().StringVarP(&browseType, "type", "t", "", "filter by tool type (agent, command, skill)")
	browseCmd.Flags().StringSliceVar(&browseTags, "tag", []string{}, "filter by tags (can specify multiple)")
	browseCmd.Flags().StringVar(&browseSortBy, "sort", "name", "sort by field (name, created, updated, downloads)")
	browseCmd.Flags().BoolVar(&browseSortDesc, "desc", false, "sort in descending order")
	browseCmd.Flags().IntVarP(&browseLimit, "limit", "l", 0, "limit number of results (0 for all)")
	browseCmd.Flags().BoolVarP(&browseJSON, "json", "j", false, "output in JSON format")

	// Trending flags (subset of browse flags)
	trendingCmd.Flags().StringVarP(&browseType, "type", "t", "", "filter by tool type (agent, command, skill)")
	trendingCmd.Flags().StringSliceVar(&browseTags, "tag", []string{}, "filter by tags (can specify multiple)")
	trendingCmd.Flags().IntVarP(&browseLimit, "limit", "l", 10, "limit number of results (default 10)")
	trendingCmd.Flags().BoolVarP(&browseJSON, "json", "j", false, "output in JSON format")
}

func runBrowse(cmd *cobra.Command, args []string) error {
	return executeBrowse(browseSortBy, browseSortDesc)
}

func runTrending(cmd *cobra.Command, args []string) error {
	// Override sort settings for trending
	return executeBrowse("downloads", true)
}

func executeBrowse(sortBy string, sortDesc bool) error {
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
		Tags:     browseTags,
		SortDesc: sortDesc,
		Limit:    browseLimit,
	}

	// Parse tool type if provided
	if browseType != "" {
		filter.Type = models.ToolType(browseType)
	}

	// Parse sort field
	switch sortBy {
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
		return fmt.Errorf("invalid sort field: %s (must be: name, created, updated, downloads)", sortBy)
	}

	// Validate filter
	if err := filter.Validate(); err != nil {
		return fmt.Errorf("invalid filter: %w", err)
	}

	// Show progress message
	if !browseJSON && verbose {
		fmt.Fprintln(os.Stderr, "Fetching registry...")
	}

	// List tools
	results, err := registryService.ListTools(filter)
	if err != nil {
		return fmt.Errorf("browse failed: %w", err)
	}

	// Display results
	if browseJSON {
		return outputJSON(results)
	}

	return displayBrowseTable(results, sortBy)
}

// displayBrowseTable displays tools in a detailed browse view
func displayBrowseTable(tools []*models.ToolInfo, sortBy string) error {
	if len(tools) == 0 {
		fmt.Println("No tools found in the registry.")
		fmt.Println("Check your registry configuration or internet connection.")
		return nil
	}

	// Prepare table data
	headers := []string{"Name", "Type", "Version", "Author", "Downloads", "Updated", "Description"}
	var rows [][]string

	for _, tool := range tools {
		description := tool.Description
		if len(description) > 60 {
			description = description[:57] + "..."
		}

		// Format updated date
		updatedStr := formatRelativeTime(tool.UpdatedAt)

		rows = append(rows, []string{
			tool.Name,
			string(tool.Type),
			tool.LatestVersion,
			tool.Author,
			fmt.Sprintf("%d", tool.Downloads),
			updatedStr,
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

	// Display summary
	sortDesc := ""
	if browseSortDesc {
		sortDesc = " (descending)"
	}
	fmt.Printf("\nShowing %d tool(s), sorted by %s%s\n", len(tools), sortBy, sortDesc)

	return nil
}

// formatRelativeTime formats a time as a relative string (e.g., "2 days ago")
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	// Less than a minute
	if diff < time.Minute {
		return "just now"
	}

	// Less than an hour
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	// Less than a day
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	// Less than a week
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	// Less than a month
	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	// Less than a year
	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	// More than a year
	years := int(diff.Hours() / 24 / 365)
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}
