# Implementation Roadmap - Claude Nia Tool Management CLI (cntm)

## Overview
Build a package manager CLI for Claude Code tools with GitHub registry backend. Estimated total time: 8-10 weeks.

**Project:** `claude-nia-tool-management-cli`
**CLI Command:** `cntm`

---

## Phase 1: Foundation & Core Models (Week 1)

### Milestone 1.1: Project Setup
- [ ] Initialize Go module
- [ ] Install dependencies (cobra, go-github, etc.)
- [ ] Set up project structure
- [ ] Configure testing framework
- [ ] Set up GitHub test repository

**Dependencies**:
```bash
go get github.com/spf13/cobra@latest
go get github.com/google/go-github/v56/github
go get golang.org/x/oauth2
go get github.com/schollz/progressbar/v3
go get github.com/olekukonko/tablewriter
go get gopkg.in/yaml.v3
go get github.com/stretchr/testify
```

**Deliverables**:
- Project scaffolding complete
- Basic build system working
- Can run `go build` successfully

### Milestone 1.2: Core Models
- [ ] Define all models in `pkg/models/`
  - [ ] ToolInfo, ToolType, Registry
  - [ ] InstalledTool, LockFile
  - [ ] ToolMetadata
  - [ ] SearchFilter, ListFilter
  - [ ] Config structures
- [ ] Write model validation methods
- [ ] Write model tests
- [ ] Create example JSON files for testing

**Deliverables**:
- Complete `pkg/models/models.go`
- Unit tests for all models
- Example registry.json and lock file

### Milestone 1.3: Configuration Management
- [ ] Implement ConfigService
- [ ] Load from YAML files
- [ ] Environment variable support
- [ ] Config precedence logic (ENV > Project > Global > Defaults)
- [ ] Write config tests

**Deliverables**:
- `internal/config/config.go` complete
- Can load and merge configs
- Tests passing

---

## Phase 2: GitHub Integration & Registry (Week 2-3)

### Milestone 2.1: GitHub Client
- [ ] Implement GitHubClient service
- [ ] Fetch file from repo
- [ ] Download large files with progress
- [ ] Authentication with PAT
- [ ] Rate limit handling
- [ ] Retry logic with backoff
- [ ] Write GitHub client tests (mocked)

**Deliverables**:
- `internal/services/github.go` complete
- Can fetch files from GitHub
- Handles rate limiting gracefully

### Milestone 2.2: Registry Service
- [ ] Implement RegistryService
- [ ] Fetch and parse registry.json
- [ ] Search tools in registry
- [ ] Get specific cntm info
- [ ] List tools with filtering
- [ ] Write registry tests

**Deliverables**:
- `internal/services/registry.go` complete
- Can fetch and parse registry
- Search works correctly

### Milestone 2.3: Cache Manager
- [ ] Implement CacheManager
- [ ] Cache registry locally
- [ ] TTL-based expiration
- [ ] Cache invalidation
- [ ] Write cache tests

**Deliverables**:
- `internal/data/cache.go` complete
- Registry caching works
- Reduces API calls

### Milestone 2.4: CLI Search & List Commands
- [ ] Implement `cntm search` command
- [ ] Implement `cntm list --remote` command
- [ ] Implement `cntm info` command
- [ ] Add table formatting for output
- [ ] Add JSON output option
- [ ] Write CLI integration tests

**Deliverables**:
- Can search registry from CLI
- Can list available tools
- Can view tool details
- Nice formatted output

**Phase 2 Demo**:
```bash
cntm search code-review
cntm list --remote --type agent
cntm info code-reviewer
```

---

## Phase 3: Installation System (Week 4-5)

### Milestone 3.1: File System Manager
- [ ] Implement FSManager
- [ ] Extract ZIP safely
- [ ] Create ZIP from directory
- [ ] Path validation (prevent traversal)
- [ ] ZIP bomb protection
- [ ] Calculate integrity hashes
- [ ] Write FS tests

**Deliverables**:
- `internal/data/fs.go` complete
- Safe ZIP operations
- Security validated

### Milestone 3.2: Lock File Service
- [ ] Implement LockFileService
- [ ] Read/write .claude-lock.json
- [ ] Add/remove/update tools
- [ ] Atomic operations
- [ ] Write lock file tests

**Deliverables**:
- `internal/services/lockfile.go` complete
- Lock file management works
- No race conditions

### Milestone 3.3: Installer Service
- [ ] Implement InstallerService
- [ ] Install single tool
- [ ] Install multiple tools
- [ ] Verify installation
- [ ] Progress tracking
- [ ] Error handling and rollback
- [ ] Write installer tests

**Installation Flow**:
1. Get cntm info from registry
2. Download ZIP from GitHub with progress bar
3. Verify integrity (SHA256)
4. Extract to `.claude/<type>/<name>/`
5. Update lock file
6. Cleanup temp files

**Deliverables**:
- `internal/services/installer.go` complete
- Can install tools successfully
- Progress bars work
- Rollback on failure

### Milestone 3.4: CLI Install Commands
- [ ] Implement `cntm install <name>` command
- [ ] Support version pinning (`name@version`)
- [ ] Support multiple installs
- [ ] Add `--path` flag for custom directory
- [ ] Add `--force` flag for reinstall
- [ ] Write CLI integration tests

**Deliverables**:
- Install command fully working
- Can install from real GitHub registry
- Good UX with progress indication

**Phase 3 Demo**:
```bash
cntm install code-reviewer
cntm install git-helper@1.0.0
cntm install agent1 agent2 agent3
cntm list  # show installed tools
```

---

## Phase 4: Update System (Week 6)

### Milestone 4.1: Updater Service
- [ ] Implement UpdaterService
- [ ] Check for outdated tools
- [ ] Update specific tool
- [ ] Update all tools
- [ ] Version comparison logic
- [ ] Write updater tests

**Deliverables**:
- `internal/services/updater.go` complete
- Can detect outdated tools
- Can update tools

### Milestone 4.2: CLI Update Commands
- [ ] Implement `cntm outdated` command
- [ ] Implement `cntm update <name>` command
- [ ] Implement `cntm update --all` command
- [ ] Show changelog/changes summary
- [ ] Add confirmation prompts
- [ ] Write CLI integration tests

**Deliverables**:
- Update commands working
- Nice table showing outdated tools
- Safe updates with confirmation

**Phase 4 Demo**:
```bash
cntm outdated
cntm update code-reviewer
cntm update --all
```

---

## Phase 5: Publishing System (Week 7-8)

### Milestone 5.1: Publisher Service Core
- [ ] Implement PublisherService
- [ ] Validate tool locally
- [ ] Generate metadata.json
- [ ] Create ZIP from tool directory
- [ ] Calculate integrity hash
- [ ] Write publisher tests

**Deliverables**:
- `internal/services/publisher.go` started
- Can create tool ZIPs
- Can generate metadata

### Milestone 5.2: GitHub Publishing
- [ ] Fork registry repo (if needed)
- [ ] Create branch
- [ ] Upload ZIP and metadata
- [ ] Update registry.json programmatically
- [ ] Create commit
- [ ] Create pull request
- [ ] Write publishing integration tests

**Deliverables**:
- Can create PRs to registry repo
- Registry updates correctly
- Good error messages

### Milestone 5.3: Create Tool Locally
- [ ] Implement create command
- [ ] Interactive prompts for metadata
- [ ] Basic template support
- [ ] Create directory structure
- [ ] Write create tests

**Deliverables**:
- `cntm create` command working
- Can create agent/command/skill locally
- Good interactive experience

### Milestone 5.4: CLI Publish Commands
- [ ] Implement `cntm create <type> <name>` command
- [ ] Implement `cntm publish <name>` command
- [ ] Version bumping logic
- [ ] Changelog prompts
- [ ] Add `--force` flag
- [ ] Write CLI integration tests

**Deliverables**:
- Publish command fully working
- Can publish to real registry
- Creates valid PRs

**Phase 5 Demo**:
```bash
cntm create agent my-agent
cntm publish my-agent --version 1.0.0
# Creates PR to registry
```

---

## Phase 6: Enhanced Features (Week 9)

### Milestone 6.1: Browse & Discovery
- [ ] Implement browse command
- [ ] Trending tools logic
- [ ] Sort by downloads/recent
- [ ] Tag-based filtering
- [ ] Nice UI for browsing
- [ ] Write browse tests

**Deliverables**:
- Browse command working
- Good discovery experience

### Milestone 6.2: Remove Tool
- [ ] Implement remove/uninstall command
- [ ] Confirmation prompts
- [ ] Update lock file
- [ ] Clean removal
- [ ] Write remove tests

**Deliverables**:
- Remove command working
- Safe deletion

### Milestone 6.3: Init Command
- [ ] Implement init command
- [ ] Create .claude directory structure
- [ ] Initialize lock file
- [ ] Set up config
- [ ] Write init tests

**Deliverables**:
- Init command working
- Easy project setup

**Phase 6 Demo**:
```bash
cntm init
cntm browse --sort downloads
cntm trending
cntm remove old-agent
```

---

## Phase 7: Polish & Documentation (Week 10)

### Milestone 7.1: Error Handling & UX
- [ ] Improve all error messages
- [ ] Add helpful hints
- [ ] Better progress indication
- [ ] Color output
- [ ] Spinner animations
- [ ] Confirmation prompts

**Deliverables**:
- Professional UX
- Clear error messages
- Beautiful output

### Milestone 7.2: Testing & Bug Fixes
- [ ] Achieve 80%+ code coverage
- [ ] Integration tests for all workflows
- [ ] End-to-end tests
- [ ] Fix all known bugs
- [ ] Cross-platform testing (macOS, Linux, Windows)

**Deliverables**:
- High test coverage
- No known bugs
- Works on all platforms

### Milestone 7.3: Documentation
- [ ] Comprehensive README
- [ ] Command reference
- [ ] Configuration guide
- [ ] Publishing guide
- [ ] Troubleshooting guide
- [ ] Example workflows
- [ ] Setup registry guide

**Deliverables**:
- Complete documentation
- Easy to get started
- Clear publishing guide

### Milestone 7.4: Setup Example Registry
- [ ] Create example GitHub registry repo
- [ ] Add sample tools
- [ ] Create registry.json
- [ ] Add metadata for each tool
- [ ] Document registry structure
- [ ] Add PR template

**Deliverables**:
- Working example registry
- Users can fork and customize

---

## Phase 8: Release (Week 11)

### Milestone 8.1: Release Preparation
- [ ] Version tagging (v1.0.0)
- [ ] Release notes
- [ ] Build for multiple platforms
- [ ] Create install scripts
- [ ] GitHub release

**Platforms**:
- macOS (amd64, arm64)
- Linux (amd64, arm64)
- Windows (amd64)

**Deliverables**:
- Binary releases for all platforms
- Installation instructions
- Release announcement

### Milestone 8.2: Distribution
- [ ] Homebrew formula (optional)
- [ ] Install script (curl | bash)
- [ ] Docker image (optional)
- [ ] AUR package (optional)

**Deliverables**:
- Easy installation methods
- Wide platform support

---

## Testing Strategy

### Unit Tests
```
internal/services/*_test.go
internal/data/*_test.go
internal/config/*_test.go
pkg/models/*_test.go
```

### Integration Tests
```
tests/integration/
├── install_test.go
├── update_test.go
├── publish_test.go
└── search_test.go
```

### E2E Tests
- Test with real GitHub test repository
- Full workflows (install, update, publish)
- Rate limit handling

---

## Dependencies & Prerequisites

### Required
- Go 1.21+
- GitHub account with PAT
- Git

### Optional for Testing
- GitHub test repository for registry
- Multiple test repositories for edge cases

---

## Risk Management

### Technical Risks

1. **GitHub Rate Limiting**
   - **Mitigation**: Implement caching, support authenticated requests (5000 req/hr)

2. **Large ZIP Files**
   - **Mitigation**: Stream downloads, show progress, implement timeouts

3. **Network Failures**
   - **Mitigation**: Retry logic, resume downloads, clear error messages

4. **Concurrent Installations**
   - **Mitigation**: File locking, atomic operations

### Project Risks

1. **Scope Creep**
   - **Mitigation**: Stick to MVP, defer advanced features to v2.0

2. **GitHub API Changes**
   - **Mitigation**: Use stable API version, add version checks

---

## Success Metrics

### v1.0 Goals
- Can search and browse tools in registry
- Can install tools from GitHub
- Can update installed tools
- Can publish tools via PR
- Lock file tracks installations
- Works reliably with good UX
- Good documentation
- Test coverage >80%

### Performance Targets
- Search: <500ms
- Install: <5s on good connection
- Update check: <1s with cache
- Publish: <10s (excluding PR creation time)

---

## Post-v1.0 Enhancements

### v1.1 - Multiple Registries
- Support multiple registry sources
- Private registries
- Registry priorities

### v1.2 - Dependencies
- Tool dependencies
- Automatic dependency installation
- Dependency graph

### v1.3 - Advanced Features
- Tool ratings/reviews
- Usage analytics
- Auto-updates
- Rollback installations
- Tool aliases

### v2.0 - Web UI
- Web-based registry browser
- Online tool editor
- Community features
- CI/CD integrations

---

## Getting Started

1. **Set up development environment**:
   ```bash
   # Install Go
   brew install go

   # Clone and setup
   git clone <repo>
   cd claude-nia-tool-management-cli
   go mod init github.com/yourusername/claude-nia-tool-management-cli
   ```

2. **Create GitHub test registry**:
   - Create new GitHub repo
   - Add registry.json
   - Add sample tools

3. **Start with Phase 1**:
   - Implement models
   - Set up config
   - Create foundation

4. **Follow milestones sequentially**:
   - Complete each milestone before moving on
   - Test thoroughly at each step
   - Keep documentation updated

---

## Summary Timeline

| Phase | Duration | Key Deliverables |
|-------|----------|-----------------|
| 1. Foundation | 1 week | Models, Config, Project setup |
| 2. GitHub & Registry | 2 weeks | GitHub client, Registry service, Search |
| 3. Installation | 2 weeks | Installer service, Install command |
| 4. Updates | 1 week | Updater service, Update commands |
| 5. Publishing | 2 weeks | Publisher service, Publish command |
| 6. Enhanced Features | 1 week | Browse, Remove, Init |
| 7. Polish | 1 week | Testing, Docs, UX improvements |
| 8. Release | 1 week | Build, Package, Distribute |

**Total: 10-11 weeks to v1.0**

Good luck with the implementation! 🚀
