# Distribution Guide

This document outlines current and planned distribution methods for cntm.

## Current Distribution Methods (v1.0.0)

### 1. Direct Binary Downloads

**Status**: ✅ Available

Pre-built binaries are available for download from GitHub Releases:

- macOS (Intel): `cntm-darwin-amd64`
- macOS (Apple Silicon): `cntm-darwin-arm64`
- Linux (x64): `cntm-linux-amd64`
- Linux (ARM64): `cntm-linux-arm64`
- Windows (x64): `cntm-windows-amd64.exe`

**Installation**:
```bash
# Download from https://github.com/yourusername/claude-nia-tool-management-cli/releases
# Extract and move to PATH
tar -xzf cntm-*.tar.gz
sudo mv cntm-* /usr/local/bin/cntm
```

### 2. Install Scripts

**Status**: ✅ Available

#### Unix/Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.sh | bash
```

Features:
- Automatic platform detection
- SHA256 checksum verification
- PATH management
- Permission handling
- User-friendly error messages

#### Windows (PowerShell)
```powershell
iwr -useb https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.ps1 | iex
```

Features:
- Platform detection
- Checksum verification
- Automatic PATH updates
- Colored output

### 3. Build from Source

**Status**: ✅ Available

```bash
git clone https://github.com/yourusername/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli
go build -o cntm
./cntm version
```

Requirements:
- Go 1.21 or later
- Git

---

## Planned Distribution Methods

### 4. Homebrew (macOS/Linux)

**Status**: 📋 Planned for v1.1

Create a Homebrew formula to enable:
```bash
brew install cntm
```

#### Steps to Implement

1. **Create Formula File** (`cntm.rb`):

```ruby
class Cntm < Formula
  desc "Package manager for Claude Code tools"
  homepage "https://github.com/yourusername/claude-nia-tool-management-cli"
  version "1.0.0"
  license "MIT"

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v1.0.0/cntm-1.0.0-darwin-amd64.tar.gz"
    sha256 "YOUR_SHA256_HERE"
  elsif OS.mac? && Hardware::CPU.arm?
    url "https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v1.0.0/cntm-1.0.0-darwin-arm64.tar.gz"
    sha256 "YOUR_SHA256_HERE"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v1.0.0/cntm-1.0.0-linux-amd64.tar.gz"
    sha256 "YOUR_SHA256_HERE"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v1.0.0/cntm-1.0.0-linux-arm64.tar.gz"
    sha256 "YOUR_SHA256_HERE"
  end

  def install
    bin.install "cntm-#{OS.kernel_name.downcase}-#{Hardware::CPU.arch}" => "cntm"
  end

  test do
    assert_match "cntm version", shell_output("#{bin}/cntm version")
  end
end
```

2. **Submit to homebrew-core**:
   - Fork [homebrew-core](https://github.com/Homebrew/homebrew-core)
   - Add formula to `Formula/cntm.rb`
   - Create pull request

3. **Or create a tap** (recommended for initial distribution):
   ```bash
   # Create homebrew-cntm repository
   # Add formula
   # Users install with:
   brew tap yourusername/cntm
   brew install cntm
   ```

**Resources**:
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Acceptable Formulae](https://docs.brew.sh/Acceptable-Formulae)

### 5. Scoop (Windows)

**Status**: 📋 Planned for v1.1

Enable Windows users to install via Scoop:
```powershell
scoop install cntm
```

#### Steps to Implement

1. **Create Manifest** (`cntm.json`):

```json
{
    "version": "1.0.0",
    "description": "Package manager for Claude Code tools",
    "homepage": "https://github.com/yourusername/claude-nia-tool-management-cli",
    "license": "MIT",
    "architecture": {
        "64bit": {
            "url": "https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v1.0.0/cntm-1.0.0-windows-amd64.zip",
            "hash": "YOUR_SHA256_HERE",
            "extract_dir": ""
        }
    },
    "bin": "cntm-windows-amd64.exe",
    "checkver": {
        "github": "https://github.com/yourusername/claude-nia-tool-management-cli"
    },
    "autoupdate": {
        "architecture": {
            "64bit": {
                "url": "https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v$version/cntm-$version-windows-amd64.zip"
            }
        }
    }
}
```

2. **Submit to scoop bucket**:
   - Fork [scoop-extras](https://github.com/ScoopInstaller/Extras)
   - Add manifest
   - Create pull request

**Resources**:
- [Scoop Documentation](https://github.com/ScoopInstaller/Scoop/wiki)
- [Creating an App Manifest](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)

### 6. AUR (Arch User Repository)

**Status**: 📋 Planned for v1.2

Enable Arch Linux users to install via AUR:
```bash
yay -S cntm
```

#### Steps to Implement

1. **Create PKGBUILD**:

```bash
# Maintainer: Your Name <your@email.com>
pkgname=cntm
pkgver=1.0.0
pkgrel=1
pkgdesc="Package manager for Claude Code tools"
arch=('x86_64' 'aarch64')
url="https://github.com/yourusername/claude-nia-tool-management-cli"
license=('MIT')
depends=()
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/yourusername/claude-nia-tool-management-cli/archive/v$pkgver.tar.gz")
sha256sums=('YOUR_SHA256_HERE')

build() {
  cd "$srcdir/claude-nia-tool-management-cli-$pkgver"
  go build -ldflags "-X github.com/yourusername/claude-nia-tool-management-cli/pkg/version.Version=$pkgver" -o cntm .
}

package() {
  cd "$srcdir/claude-nia-tool-management-cli-$pkgver"
  install -Dm755 cntm "$pkgdir/usr/bin/cntm"
  install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
```

2. **Or use binary package** (PKGBUILD-bin):

```bash
# Maintainer: Your Name <your@email.com>
pkgname=cntm-bin
pkgver=1.0.0
pkgrel=1
pkgdesc="Package manager for Claude Code tools (binary)"
arch=('x86_64' 'aarch64')
url="https://github.com/yourusername/claude-nia-tool-management-cli"
license=('MIT')
provides=('cntm')
conflicts=('cntm')

source_x86_64=("https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v$pkgver/cntm-$pkgver-linux-amd64.tar.gz")
source_aarch64=("https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v$pkgver/cntm-$pkgver-linux-arm64.tar.gz")

sha256sums_x86_64=('YOUR_SHA256_HERE')
sha256sums_aarch64=('YOUR_SHA256_HERE')

package() {
  install -Dm755 cntm-linux-* "$pkgdir/usr/bin/cntm"
}
```

3. **Publish to AUR**:
   - Create account on aur.archlinux.org
   - Clone AUR repository: `git clone ssh://aur@aur.archlinux.org/cntm.git`
   - Add PKGBUILD
   - Push to AUR

**Resources**:
- [AUR Submission Guidelines](https://wiki.archlinux.org/title/AUR_submission_guidelines)
- [Creating Packages](https://wiki.archlinux.org/title/Creating_packages)

### 7. Docker Image

**Status**: 📋 Planned for v1.2

Provide a Docker image for containerized usage:
```bash
docker run -v $(pwd)/.claude:/workspace/.claude ghcr.io/yourusername/cntm:latest search code-review
```

#### Steps to Implement

1. **Create Dockerfile**:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" -o cntm .

FROM alpine:latest

RUN apk --no-cache add ca-certificates git
COPY --from=builder /build/cntm /usr/local/bin/cntm

WORKDIR /workspace
ENTRYPOINT ["cntm"]
CMD ["--help"]
```

2. **Multi-arch build**:

```bash
# Build for multiple architectures
docker buildx build --platform linux/amd64,linux/arm64 \
  -t ghcr.io/yourusername/cntm:1.0.0 \
  -t ghcr.io/yourusername/cntm:latest \
  --push .
```

3. **GitHub Actions for automated builds**:

```yaml
name: Docker

on:
  push:
    tags:
      - 'v*'

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: docker/setup-buildx-action@v2
      - uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v4
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:${{ github.ref_name }}
            ghcr.io/${{ github.repository }}:latest
```

**Resources**:
- [Docker Multi-platform Builds](https://docs.docker.com/build/building/multi-platform/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)

### 8. Snap (Linux)

**Status**: 📋 Planned for v1.3

Enable installation via Snap Store:
```bash
sudo snap install cntm
```

#### Steps to Implement

1. **Create snapcraft.yaml**:

```yaml
name: cntm
version: '1.0.0'
summary: Package manager for Claude Code tools
description: |
  cntm is a package manager for Claude Code tools (agents, commands, and skills).
  Like npm for Node.js, cntm helps you install, update, and manage Claude Code tools.

grade: stable
confinement: strict
base: core22

apps:
  cntm:
    command: cntm
    plugs:
      - home
      - network

parts:
  cntm:
    plugin: go
    source: .
    build-snaps: [go/latest/stable]
    build-packages:
      - git
```

2. **Build and publish**:

```bash
snapcraft
snapcraft upload --release=stable cntm_1.0.0_amd64.snap
```

**Resources**:
- [Snapcraft.io](https://snapcraft.io/docs)
- [Building Go Snaps](https://snapcraft.io/docs/go-applications)

### 9. Chocolatey (Windows)

**Status**: 📋 Planned for v1.3

Enable Windows users to install via Chocolatey:
```powershell
choco install cntm
```

#### Steps to Implement

1. **Create nuspec file** (`cntm.nuspec`):

```xml
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>cntm</id>
    <version>1.0.0</version>
    <title>cntm</title>
    <authors>Your Name</authors>
    <projectUrl>https://github.com/yourusername/claude-nia-tool-management-cli</projectUrl>
    <licenseUrl>https://github.com/yourusername/claude-nia-tool-management-cli/blob/main/LICENSE</licenseUrl>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
    <description>Package manager for Claude Code tools</description>
    <summary>Like npm for Claude Code - manage agents, commands, and skills</summary>
    <tags>cli tools package-manager claude</tags>
  </metadata>
  <files>
    <file src="tools\**" target="tools" />
  </files>
</package>
```

2. **Create install script** (`tools/chocolateyinstall.ps1`):

```powershell
$packageName = 'cntm'
$url64 = 'https://github.com/yourusername/claude-nia-tool-management-cli/releases/download/v1.0.0/cntm-1.0.0-windows-amd64.zip'
$checksum64 = 'YOUR_SHA256_HERE'

Install-ChocolateyZipPackage `
  -PackageName $packageName `
  -Url64bit $url64 `
  -UnzipLocation "$(Split-Path -parent $MyInvocation.MyCommand.Definition)" `
  -Checksum64 $checksum64 `
  -ChecksumType64 'sha256'
```

3. **Submit to Chocolatey**:
   - Create account on chocolatey.org
   - Submit package for moderation

**Resources**:
- [Chocolatey Documentation](https://docs.chocolatey.org/en-us/create/create-packages)
- [Package Creation](https://docs.chocolatey.org/en-us/create/create-packages-quick-start)

---

## Package Manager Comparison

| Method | Platforms | Effort | Auto-Update | Priority |
|--------|-----------|--------|-------------|----------|
| Direct Download | All | Low | No | ✅ v1.0 |
| Install Scripts | All | Low | No | ✅ v1.0 |
| Build from Source | All | Low | No | ✅ v1.0 |
| Homebrew | macOS/Linux | Medium | Yes | 📋 v1.1 |
| Scoop | Windows | Medium | Yes | 📋 v1.1 |
| AUR | Arch Linux | Medium | Yes | 📋 v1.2 |
| Docker | All | Medium | Yes | 📋 v1.2 |
| Snap | Linux | High | Yes | 📋 v1.3 |
| Chocolatey | Windows | High | Yes | 📋 v1.3 |

## Verification Process

All distribution methods should verify:

1. **SHA256 Checksums** - Verify binary integrity
2. **GPG Signatures** (future) - Verify authenticity
3. **Build Reproducibility** - Ensure consistent builds

## Update Notifications

Future versions will include:
- Built-in update checking
- Notification of new versions
- Self-update capability via package managers

## Contributing

To add a new distribution method:

1. Create the necessary configuration files
2. Test the installation process
3. Update this document
4. Submit a pull request

---

**Last Updated**: November 15, 2025
**Version**: 1.0.0
