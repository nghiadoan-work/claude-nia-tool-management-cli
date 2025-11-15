package version

// Build-time variables, set via ldflags during compilation
// Example: go build -ldflags "-X github.com/yourusername/claude-nia-tool-management-cli/pkg/version.GitCommit=$(git rev-parse HEAD)"
var (
	// Version is the semantic version of the application
	Version = "1.0.0"

	// GitCommit is the git commit hash (set during build)
	GitCommit = "dev"

	// BuildDate is the build timestamp (set during build)
	BuildDate = "unknown"

	// GoVersion is the Go version used to build
	GoVersion = "unknown"
)

// Info returns a structured representation of version information
type Info struct {
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"git_commit" yaml:"git_commit"`
	BuildDate string `json:"build_date" yaml:"build_date"`
	GoVersion string `json:"go_version" yaml:"go_version"`
}

// GetInfo returns the version information
func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
	}
}

// String returns a human-readable version string
func (i Info) String() string {
	return i.Version
}

// LongString returns a detailed version string
func (i Info) LongString() string {
	return "cntm version " + i.Version +
		"\nGit commit: " + i.GitCommit +
		"\nBuild date: " + i.BuildDate +
		"\nGo version: " + i.GoVersion
}
