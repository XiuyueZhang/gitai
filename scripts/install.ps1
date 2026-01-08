# GitAI Installer Script for Windows
# Run with: iwr -useb https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.ps1 | iex

$ErrorActionPreference = 'Stop'

$REPO = "xyue92/gitai"
$BINARY_NAME = "gitai.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\Programs\GitAI"

# Colors
function Write-Info {
    param($Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Error-Custom {
    param($Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Write-Warning-Custom {
    param($Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64" -or $arch -eq "x86_64") {
        return "amd64"
    }
    elseif ($arch -eq "ARM64") {
        return "arm64"
    }
    else {
        Write-Error-Custom "Unsupported architecture: $arch"
        exit 1
    }
}

# Get latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest release..."
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
        $version = $response.tag_name
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error-Custom "Failed to get latest version: $_"
        exit 1
    }
}

# Download binary
function Download-Binary {
    param($Version, $Arch)

    $downloadFilename = "gitai-windows-$Arch.exe"
    $downloadUrl = "https://github.com/$REPO/releases/download/$Version/$downloadFilename"
    $checksumUrl = "$downloadUrl.sha256"

    $tempDir = [System.IO.Path]::GetTempPath()
    $tempFile = Join-Path $tempDir $downloadFilename
    $tempChecksum = "$tempFile.sha256"

    Write-Info "Downloading from: $downloadUrl"

    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
    }
    catch {
        Write-Error-Custom "Failed to download binary: $_"
        exit 1
    }

    # Download and verify checksum
    try {
        Invoke-WebRequest -Uri $checksumUrl -OutFile $tempChecksum -UseBasicParsing
        Write-Info "Verifying checksum..."

        $checksumContent = Get-Content $tempChecksum -Raw
        $expectedHash = ($checksumContent -split '\s+')[0]
        $actualHash = (Get-FileHash $tempFile -Algorithm SHA256).Hash

        if ($expectedHash.ToLower() -ne $actualHash.ToLower()) {
            Write-Error-Custom "Checksum verification failed"
            Write-Error-Custom "Expected: $expectedHash"
            Write-Error-Custom "Got: $actualHash"
            Remove-Item $tempFile, $tempChecksum -Force
            exit 1
        }

        Write-Info "Checksum verified successfully"
        Remove-Item $tempChecksum -Force
    }
    catch {
        Write-Warning-Custom "Checksum verification skipped: $_"
    }

    return $tempFile
}

# Install binary
function Install-Binary {
    param($TempFile)

    # Create install directory if it doesn't exist
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }

    $targetFile = Join-Path $INSTALL_DIR $BINARY_NAME

    # Remove old version if exists
    if (Test-Path $targetFile) {
        Write-Info "Removing old version..."
        Remove-Item $targetFile -Force
    }

    # Move new binary
    Move-Item $TempFile $targetFile -Force
    Write-Info "Binary installed to: $targetFile"
}

# Add to PATH
function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($currentPath -notlike "*$INSTALL_DIR*") {
        Write-Info "Adding to PATH..."
        $newPath = "$currentPath;$INSTALL_DIR"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")

        # Update current session PATH
        $env:Path = "$env:Path;$INSTALL_DIR"

        Write-Info "Added to PATH. You may need to restart your terminal."
    }
    else {
        Write-Info "Already in PATH"
    }
}

# Verify installation
function Test-Installation {
    $gitaiPath = Join-Path $INSTALL_DIR $BINARY_NAME

    if (Test-Path $gitaiPath) {
        Write-Host ""
        Write-Host "✓ GitAI has been installed successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Installation location: $gitaiPath" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Try it out:" -ForegroundColor Yellow
        Write-Host "  gitai --help"
        Write-Host ""
        Write-Host "Note: If 'gitai' command is not found, please restart your terminal." -ForegroundColor Yellow
    }
    else {
        Write-Error-Custom "Installation verification failed"
        exit 1
    }
}

# Main installation flow
function Main {
    Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║   GitAI Installer for Windows          ║" -ForegroundColor Cyan
    Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
    Write-Host ""

    $arch = Get-Architecture
    Write-Info "Detected architecture: $arch"

    $version = Get-LatestVersion
    $tempFile = Download-Binary -Version $version -Arch $arch
    Install-Binary -TempFile $tempFile
    Add-ToPath
    Test-Installation

    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Green
    Write-Host "1. Install Ollama from: https://ollama.com/download"
    Write-Host "2. Pull an AI model: ollama pull qwen2.5-coder:7b"
    Write-Host "3. Navigate to a git repository and run: gitai commit"
    Write-Host ""
}

# Run installation
Main
