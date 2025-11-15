# install.ps1 - PowerShell installation script for cntm (Windows)
# Usage: iwr -useb https://raw.githubusercontent.com/USER/REPO/main/scripts/install.ps1 | iex

param(
    [string]$Version = "1.0.0",
    [string]$Repo = "yourusername/claude-nia-tool-management-cli",
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\cntm"
)

$ErrorActionPreference = "Stop"

# Colors
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

Write-ColorOutput "================================================" "Blue"
Write-ColorOutput "  Installing cntm v$Version" "Blue"
Write-ColorOutput "================================================" "Blue"
Write-Host ""

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$Platform = "windows-$Arch"

Write-ColorOutput "Detected platform: $Platform" "Green"

# Create install directory
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-ColorOutput "  ✓ Created install directory" "Green"
}

# Create temp directory
$TmpDir = Join-Path $env:TEMP "cntm-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null

try {
    # Download archive
    $ArchiveName = "cntm-$Version-$Platform.zip"
    $DownloadUrl = "https://github.com/$Repo/releases/download/v$Version/$ArchiveName"
    $ArchivePath = Join-Path $TmpDir $ArchiveName

    Write-ColorOutput "Downloading from: $DownloadUrl" "Blue"

    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchivePath -UseBasicParsing
        Write-ColorOutput "  ✓ Downloaded" "Green"
    }
    catch {
        Write-ColorOutput "Error: Failed to download cntm" "Red"
        Write-ColorOutput "Please check that version $Version exists at:" "Yellow"
        Write-ColorOutput "  https://github.com/$Repo/releases/tag/v$Version" "Yellow"
        exit 1
    }

    # Download checksums (optional)
    $ChecksumUrl = "https://github.com/$Repo/releases/download/v$Version/checksums.txt"
    $ChecksumPath = Join-Path $TmpDir "checksums.txt"

    try {
        Write-Host "Downloading checksums..."
        Invoke-WebRequest -Uri $ChecksumUrl -OutFile $ChecksumPath -UseBasicParsing

        Write-Host "Verifying checksum..."

        # Calculate SHA256
        $Hash = Get-FileHash -Path $ArchivePath -Algorithm SHA256
        $ActualChecksum = $Hash.Hash.ToLower()

        # Read expected checksum
        $ChecksumContent = Get-Content $ChecksumPath
        $BinaryName = "cntm-$Platform.exe"
        $ExpectedLine = $ChecksumContent | Where-Object { $_ -match $BinaryName }

        if ($ExpectedLine) {
            $ExpectedChecksum = ($ExpectedLine -split '\s+')[0].ToLower()

            if ($ActualChecksum -eq $ExpectedChecksum) {
                Write-ColorOutput "  ✓ Checksum verified" "Green"
            }
            else {
                Write-ColorOutput "Error: Checksum verification failed" "Red"
                Write-ColorOutput "Expected: $ExpectedChecksum" "Yellow"
                Write-ColorOutput "Actual:   $ActualChecksum" "Yellow"
                exit 1
            }
        }
        else {
            Write-ColorOutput "Warning: Could not find checksum for $BinaryName" "Yellow"
        }
    }
    catch {
        Write-ColorOutput "Warning: Could not verify checksum" "Yellow"
    }

    # Extract archive
    Write-Host "Extracting archive..."
    Expand-Archive -Path $ArchivePath -DestinationPath $TmpDir -Force
    Write-ColorOutput "  ✓ Extracted" "Green"

    # Install binary
    $BinaryName = "cntm-$Platform.exe"
    $SourcePath = Join-Path $TmpDir $BinaryName
    $DestPath = Join-Path $InstallDir "cntm.exe"

    Write-Host "Installing to $InstallDir..."
    Copy-Item -Path $SourcePath -Destination $DestPath -Force
    Write-ColorOutput "  ✓ Installed" "Green"

    # Add to PATH if not already present
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        Write-Host "Adding to PATH..."
        $NewPath = "$UserPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        Write-ColorOutput "  ✓ Added to PATH" "Green"
        Write-ColorOutput "Note: Restart your terminal for PATH changes to take effect" "Yellow"
    }

    # Verify installation
    Write-Host ""
    Write-ColorOutput "================================================" "Green"
    Write-ColorOutput "  cntm installed successfully!" "Green"
    Write-ColorOutput "================================================" "Green"
    Write-Host "Location: $DestPath"
    Write-Host ""
    Write-Host "Get started:"
    Write-Host "  cntm init              # Initialize your project"
    Write-Host "  cntm search <query>    # Search for tools"
    Write-Host "  cntm install <name>    # Install a tool"
    Write-Host ""
    Write-ColorOutput "Note: Restart your terminal for PATH changes to take effect" "Yellow"
    Write-Host ""
}
finally {
    # Cleanup
    if (Test-Path $TmpDir) {
        Remove-Item -Path $TmpDir -Recurse -Force
    }
}
