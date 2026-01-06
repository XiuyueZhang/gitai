# GitAI Installation Guide

Complete installation instructions for all platforms.

## Prerequisites

### 1. Go Programming Language

**Minimum Version**: Go 1.21 or higher

#### Check if Go is installed:
```bash
go version
```

#### Install Go:

**macOS (Homebrew):**
```bash
brew install go
```

**macOS (Official):**
1. Download from [go.dev/dl](https://go.dev/dl/)
2. Open the `.pkg` file
3. Follow installer

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install golang-go
```

**Linux (Official):**
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Windows:**
1. Download installer from [go.dev/dl](https://go.dev/dl/)
2. Run the `.msi` file
3. Follow installer

---

### 2. Ollama (Local AI Runtime)

**Required**: For running AI models locally

#### macOS/Linux:
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

#### Windows:
1. Download from [ollama.com/download](https://ollama.com/download)
2. Run installer
3. Ollama will start automatically

#### Verify Installation:
```bash
ollama --version
```

#### Start Ollama Service:
```bash
ollama serve
```

Keep this terminal window open, or run as background service.

---

### 3. AI Model

**Recommended**: qwen2.5-coder:7b (4.7 GB)

```bash
ollama pull qwen2.5-coder:7b
```

**Alternative Models:**
```bash
# Faster, smaller
ollama pull mistral:7b

# Good for code
ollama pull codellama:7b

# Larger, more accurate (if you have RAM)
ollama pull qwen2.5-coder:14b
```

**Check installed models:**
```bash
ollama list
```

---

## Installing GitAI

### Option 1: Build from Source (Recommended)

#### Step 1: Navigate to Project
```bash
cd /Users/yue/go-source/gitai
```

#### Step 2: Download Dependencies
```bash
go mod download
```

#### Step 3: Build
```bash
make build
# or manually:
# go build -o gitai .
```

#### Step 4: Install to System (Optional)
```bash
# Install to /usr/local/bin
sudo make install

# Or copy manually
sudo cp gitai /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)
```

#### Verify Installation:
```bash
gitai --version
```

---

### Option 2: Go Install (If Published)

```bash
go install github.com/yourusername/gitai@latest
```

Make sure `$GOPATH/bin` is in your PATH:
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

---

### Option 3: Download Pre-built Binary (Future)

Once releases are available:

1. Go to [GitHub Releases](https://github.com/yourusername/gitai/releases)
2. Download binary for your platform:
   - `gitai-darwin-amd64` (macOS Intel)
   - `gitai-darwin-arm64` (macOS Apple Silicon)
   - `gitai-linux-amd64` (Linux)
   - `gitai-windows-amd64.exe` (Windows)
3. Make executable (macOS/Linux):
   ```bash
   chmod +x gitai-*
   sudo mv gitai-* /usr/local/bin/gitai
   ```
4. Windows: Move to a directory in PATH

---

## Platform-Specific Instructions

### macOS

#### Apple Silicon (M1/M2/M3):
```bash
cd /Users/yue/go-source/gitai
GOARCH=arm64 go build -o gitai
sudo mv gitai /usr/local/bin/
```

#### Intel:
```bash
cd /Users/yue/go-source/gitai
GOARCH=amd64 go build -o gitai
sudo mv gitai /usr/local/bin/
```

#### Homebrew Installation (Future):
```bash
brew install gitai
```

---

### Linux

#### Ubuntu/Debian:
```bash
# Install prerequisites
sudo apt update
sudo apt install golang-go git

# Build GitAI
cd /path/to/gitai
make build
sudo make install
```

#### Fedora/RHEL:
```bash
sudo dnf install golang git
cd /path/to/gitai
make build
sudo make install
```

#### Arch Linux:
```bash
sudo pacman -S go git
cd /path/to/gitai
make build
sudo make install
```

---

### Windows

#### Using PowerShell:

1. **Install Go**:
   - Download from [go.dev/dl](https://go.dev/dl/)
   - Run installer

2. **Install Ollama**:
   - Download from [ollama.com](https://ollama.com/download)
   - Run installer

3. **Build GitAI**:
   ```powershell
   cd C:\path\to\gitai
   go build -o gitai.exe
   ```

4. **Add to PATH**:
   - Copy `gitai.exe` to `C:\Windows\System32`
   - Or add gitai directory to PATH

5. **Pull AI Model**:
   ```powershell
   ollama pull qwen2.5-coder:7b
   ```

---

## Verification

### Test Installation:

```bash
# Check GitAI is installed
gitai --version

# Check Ollama is running
curl http://localhost:11434/api/tags

# Check model is available
ollama list
```

### Quick Test:

```bash
# Create test repo
mkdir test-gitai
cd test-gitai
git init

# Make a change
echo "# Test" > README.md
git add README.md

# Run GitAI
gitai commit --dry-run
```

If you see the interactive prompt and message generation, it works!

---

## Configuration

### Create Config File:
```bash
cd ~/your-project
gitai config --init
```

### Edit Configuration:
```bash
vim .gitcommit.yaml
```

### Example Config:
```yaml
model: "qwen2.5-coder:7b"
language: "en"
scopes:
  - "frontend"
  - "backend"
  - "api"
```

---

## Troubleshooting Installation

### Go Not Found:
```bash
# Check Go installation
which go
go version

# Add to PATH if needed
export PATH=$PATH:/usr/local/go/bin
```

### Ollama Connection Failed:
```bash
# Start Ollama
ollama serve

# In another terminal
curl http://localhost:11434/api/tags
```

### Model Not Found:
```bash
# List installed models
ollama list

# Pull model if missing
ollama pull qwen2.5-coder:7b
```

### Permission Denied (macOS/Linux):
```bash
# Make binary executable
chmod +x gitai

# Use sudo for system install
sudo mv gitai /usr/local/bin/
```

### Build Errors:
```bash
# Clean and retry
make clean
go clean -modcache
go mod download
make build
```

---

## Updating GitAI

### If Built from Source:
```bash
cd /path/to/gitai
git pull  # If tracking git repo
make clean
make build
sudo make install
```

### If Installed with Go:
```bash
go install github.com/yourusername/gitai@latest
```

---

## Uninstallation

### Remove Binary:
```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/gitai

# If installed via Go
rm $(go env GOPATH)/bin/gitai
```

### Remove Config (Optional):
```bash
rm ~/.gitcommit.yaml
```

### Keep Ollama and Models:
Ollama and models can be used by other tools, so you may want to keep them.

### Remove Everything:
```bash
# Remove GitAI
sudo rm /usr/local/bin/gitai

# Remove Ollama (optional)
# macOS
brew uninstall ollama

# Linux
sudo rm /usr/local/bin/ollama

# Windows: Use "Add or Remove Programs"
```

---

## Next Steps

1. ✅ GitAI installed
2. ✅ Ollama running
3. ✅ Model downloaded
4. 📖 Read [QUICKSTART.md](QUICKSTART.md)
5. 🚀 Start using: `gitai commit`

---

## Getting Help

- **Documentation**: See [README.md](README.md)
- **Examples**: See [EXAMPLES.md](EXAMPLES.md)
- **Quick Start**: See [QUICKSTART.md](QUICKSTART.md)
- **Issues**: Check your Ollama and Git setup first

---

## System Requirements

### Minimum:
- OS: macOS 10.15+, Linux (any modern distro), Windows 10+
- RAM: 8 GB (for 7B models)
- Disk: 10 GB free (for model + binary)
- Go: 1.21+

### Recommended:
- OS: macOS 12+, Ubuntu 22.04+, Windows 11
- RAM: 16 GB
- Disk: 20 GB free
- Go: 1.21+
- SSD for faster model loading

### For Larger Models (13B+):
- RAM: 32 GB+
- Disk: 30 GB+

---

Happy Installing! 🚀
