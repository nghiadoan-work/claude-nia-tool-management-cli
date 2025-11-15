# Installation Guide

Multiple ways to install cntm on your system.

## Quick Install (Recommended)

### Using npx (No Installation Required)

Run cntm directly without installing:

```bash
npx cntm search "code review"
npx cntm install code-reviewer
npx cntm --help
```

### Using npx from GitHub

```bash
npx github:nghiadt/claude-nia-tool-management-cli
```

## Install with npm

### Global Installation

Install once, use everywhere:

```bash
npm install -g cntm
cntm --help
```

### Local Project Installation

Install in your Node.js project:

```bash
npm install cntm
npx cntm init
```

Add to package.json scripts:

```json
{
  "scripts": {
    "tools:init": "cntm init",
    "tools:install": "cntm install",
    "tools:update": "cntm update --all"
  }
}
```

## Install with curl/wget

### macOS and Linux

```bash
curl -fsSL https://raw.githubusercontent.com/nghiadt/claude-nia-tool-management-cli/main/scripts/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/nghiadt/claude-nia-tool-management-cli/main/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/nghiadt/claude-nia-tool-management-cli/main/scripts/install.ps1 | iex
```

## Download Pre-built Binaries

1. Go to [GitHub Releases](https://github.com/nghiadt/claude-nia-tool-management-cli/releases)
2. Download the appropriate file for your platform:
   - macOS Intel: `cntm-1.0.0-darwin-amd64.tar.gz`
   - macOS Apple Silicon: `cntm-1.0.0-darwin-arm64.tar.gz`
   - Linux x64: `cntm-1.0.0-linux-amd64.tar.gz`
   - Linux ARM64: `cntm-1.0.0-linux-arm64.tar.gz`
   - Windows x64: `cntm-1.0.0-windows-amd64.zip`
3. Extract the archive:

   ```bash
   # macOS/Linux
   tar -xzf cntm-1.0.0-darwin-arm64.tar.gz
   
   # Windows
   # Use Windows Explorer or:
   tar -xf cntm-1.0.0-windows-amd64.zip
   ```

4. Move to PATH:

   ```bash
   # macOS/Linux
   sudo mv cntm-darwin-arm64 /usr/local/bin/cntm
   chmod +x /usr/local/bin/cntm
   
   # Windows (PowerShell as Admin)
   Move-Item cntm-windows-amd64.exe C:\Windows\System32\cntm.exe
   ```

5. Verify installation:

   ```bash
   cntm version
   ```

## Build from Source

### Prerequisites

- Go 1.21 or later
- Git

### Steps

```bash
# Clone repository
git clone https://github.com/nghiadt/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli

# Build
go build -o cntm .

# Move to PATH (optional)
sudo mv cntm /usr/local/bin/

# Verify
cntm version
```

### Build for Specific Platform

```bash
# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o cntm-darwin-arm64

# Linux x64
GOOS=linux GOARCH=amd64 go build -o cntm-linux-amd64

# Windows x64
GOOS=windows GOARCH=amd64 go build -o cntm-windows-amd64.exe
```

## Package Managers

### Homebrew (Coming Soon)

```bash
brew install cntm
```

### Scoop (Windows - Coming Soon)

```bash
scoop install cntm
```

### AUR (Arch Linux - Coming Soon)

```bash
yay -S cntm
```

## Docker (Optional)

Run cntm in Docker:

```bash
docker run -it --rm \
  -v $(pwd):/workspace \
  -w /workspace \
  ghcr.io/nghiadt/cntm:latest \
  --help
```

## Verify Installation

After installation, verify it works:

```bash
# Check version
cntm version

# Check help
cntm --help

# Try a command
cntm search test
```

## Update cntm

### NPM Installation

```bash
npm update -g cntm
```

### Manual Installation

Download the latest release and replace your existing binary.

### Built from Source

```bash
cd claude-nia-tool-management-cli
git pull
go build -o cntm .
```

## Uninstall

### NPM Installation

```bash
npm uninstall -g cntm
```

### Manual Installation

```bash
# macOS/Linux
sudo rm /usr/local/bin/cntm

# Windows
Remove-Item C:\Windows\System32\cntm.exe
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Check if the binary is in your PATH:
   ```bash
   which cntm  # macOS/Linux
   where cntm  # Windows
   ```

2. Add to PATH manually:
   ```bash
   # macOS/Linux (add to ~/.bashrc or ~/.zshrc)
   export PATH="$PATH:/usr/local/bin"
   
   # Windows (PowerShell as Admin)
   $env:Path += ";C:\path\to\cntm"
   ```

### Permission Denied (macOS/Linux)

```bash
chmod +x /usr/local/bin/cntm
```

### Windows Security Warning

Windows might block the executable. To allow it:
1. Right-click the .exe file
2. Properties → Unblock → OK

Or use PowerShell:
```powershell
Unblock-File cntm-windows-amd64.exe
```

### NPM Installation Fails

If npm installation fails:

1. Try with sudo (macOS/Linux):
   ```bash
   sudo npm install -g cntm
   ```

2. Configure npm to use a different directory:
   ```bash
   mkdir ~/.npm-global
   npm config set prefix '~/.npm-global'
   export PATH=~/.npm-global/bin:$PATH
   npm install -g cntm
   ```

3. Use npx instead (no installation needed):
   ```bash
   npx cntm --help
   ```

## Platform Support

| Platform | Architecture | Supported | Method |
|----------|--------------|-----------|--------|
| macOS | x64 (Intel) | ✅ | npm, curl, binary |
| macOS | ARM64 (M1/M2) | ✅ | npm, curl, binary |
| Linux | x64 | ✅ | npm, curl, binary |
| Linux | ARM64 | ✅ | npm, curl, binary |
| Windows | x64 | ✅ | npm, PowerShell, binary |

## Next Steps

After installation:

1. **Initialize a project:**
   ```bash
   cntm init
   ```

2. **Search for tools:**
   ```bash
   cntm search "your query"
   ```

3. **Install a tool:**
   ```bash
   cntm install tool-name
   ```

4. **Read the docs:**
   ```bash
   cntm help
   ```

For more information, see:
- [Command Reference](./COMMANDS.md)
- [Configuration Guide](./CONFIGURATION.md)
- [Troubleshooting](./TROUBLESHOOTING.md)
