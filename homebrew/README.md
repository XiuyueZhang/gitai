# Homebrew Tap Setup Guide

This directory contains the Homebrew formula for GitAI. To enable Homebrew installation, you need to create a separate tap repository.

## Step 1: Create Tap Repository

Create a new GitHub repository named `homebrew-tap` under your account:

```bash
# Repository URL will be: https://github.com/xyue92/homebrew-tap
```

## Step 2: Setup Tap Repository Structure

```bash
# Clone the new repository
git clone https://github.com/xyue92/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir -p Formula

# Copy the formula file
cp /path/to/gitai/homebrew/gitai.rb Formula/gitai.rb

# Commit and push
git add Formula/gitai.rb
git commit -m "Add GitAI formula"
git push origin main
```

## Step 3: Update Formula After Each Release

After creating a new release, update the formula:

1. Download the release binaries and calculate SHA256:
```bash
# For macOS ARM64
wget https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-darwin-arm64
shasum -a 256 gitai-darwin-arm64

# For macOS AMD64
wget https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-darwin-amd64
shasum -a 256 gitai-darwin-amd64

# For Linux ARM64
wget https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-linux-arm64
shasum -a 256 gitai-linux-arm64

# For Linux AMD64
wget https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-linux-amd64
shasum -a 256 gitai-linux-amd64
```

2. Update `Formula/gitai.rb`:
   - Update `version` to the new version number
   - Update all download URLs
   - Replace SHA256 hashes with the calculated values

3. Commit and push:
```bash
git add Formula/gitai.rb
git commit -m "Update GitAI to v1.0.0"
git push origin main
```

## Step 4: Users Can Install

Once the tap is set up, users can install with:

```bash
# Add the tap (first time only)
brew tap xyue92/tap

# Install GitAI
brew install gitai

# Update GitAI
brew upgrade gitai
```

## Automation (Optional)

You can automate formula updates using GitHub Actions in your homebrew-tap repository:

1. Create `.github/workflows/update-formula.yml` in the tap repository
2. Trigger on new releases from the main gitai repository
3. Automatically calculate SHA256 and update the formula

## Testing the Formula

Before publishing:

```bash
# Install from local formula
brew install --build-from-source Formula/gitai.rb

# Test the installation
gitai --help

# Uninstall
brew uninstall gitai
```

## Alternative: Use `brew bump-formula-pr`

Homebrew provides a tool to help update formulas:

```bash
brew bump-formula-pr \
  --url=https://github.com/xyue92/gitai/releases/download/v1.0.0/gitai-darwin-arm64 \
  gitai
```

This will automatically:
- Download the file
- Calculate SHA256
- Update the formula
- Create a pull request (if using homebrew-core)
