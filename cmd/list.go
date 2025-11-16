package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/olekukonko/tablewriter"
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
	Short: "List installed tools or tools from the registry",
	Long: `List locally installed tools or tools from the remote registry.

By default, this command lists locally installed tools.
Use --remote to list all available tools in the registry.

Examples:
  cntm list                               # List locally installed tools
  cntm list --json                        # List local tools in JSON format
  cntm list --remote                      # List all remote tools
  cntm list --remote --type agent         # List only agents
  cntm list --remote --tag git            # List tools with "git" tag
  cntm list --remote --sort-by downloads  # Sort by download count
  cntm list --remote --sort-desc          # Sort in descending order
  cntm list --remote --limit 10           # Limit to 10 results`,
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
		// List locally installed tools
		return runListLocal()
	}

	// List remote tools from registry
	return runListRemote()
}

// runListLocal lists locally installed tools
func runListLocal() error {
	// Load lock file service
	lockFilePath := filepath.Join(basePath, ".claude-lock.json")
	lockFileService, err := services.NewLockFileService(lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to create lock file service: %w", err)
	}

	// Get installed tools
	tools, err := lockFileService.ListTools()
	if err != nil {
		return fmt.Errorf("failed to list installed tools: %w\nHint: No tools installed yet? Use 'cntm install <name>' to install tools", err)
	}

	if len(tools) == 0 {
		fmt.Println("No tools installed yet.")
		fmt.Println("Use 'cntm install <name>' to install tools from the registry.")
		return nil
	}

	// Display results
	if listJSON {
		return outputJSON(tools)
	}

	return displayInstalledTools(tools)
}

// runListRemote lists tools from the remote registry
func runListRemote() error {
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

// displayInstalledTools displays locally installed tools in a table format
func displayInstalledTools(tools map[string]*models.InstalledTool) error {
	if len(tools) == 0 {
		fmt.Println("No tools installed yet.")
		return nil
	}

	// Convert map to sorted slice for consistent ordering
	type toolEntry struct {
		name string
		tool *models.InstalledTool
	}
	var entries []toolEntry
	for name, tool := range tools {
		entries = append(entries, toolEntry{name: name, tool: tool})
	}

	// Sort by name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].name < entries[j].name
	})

	// Prepare table data
	headers := []string{"Name", "Version", "Type", "Installed At"}
	var rows [][]string

	for _, entry := range entries {
		rows = append(rows, []string{
			entry.name,
			entry.tool.Version,
			string(entry.tool.Type),
			entry.tool.InstalledAt.Format("2006-01-02 15:04:05"),
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
	fmt.Printf("\n%d tool(s) installed\n", len(tools))

	return nil
}
