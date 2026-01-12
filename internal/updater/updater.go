package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubAPIURL = "https://api.github.com/repos/xyue92/gitai/releases/latest"
	githubRawURL = "https://github.com/xyue92/gitai/releases/download"
	timeout      = 30 * time.Second
)

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
}

// Updater handles self-update logic
type Updater struct {
	CurrentVersion string
	Owner          string
	Repo           string
}

// New creates a new Updater instance
func New(currentVersion string) *Updater {
	return &Updater{
		CurrentVersion: currentVersion,
		Owner:          "xyue92",
		Repo:           "gitai",
	}
}

// CheckForUpdate checks if a newer version is available
func (u *Updater) CheckForUpdate() (*Release, bool, error) {
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Accept header for GitHub API
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Compare versions
	needsUpdate := u.needsUpdate(u.CurrentVersion, release.TagName)

	return &release, needsUpdate, nil
}

// needsUpdate compares current version with latest version
func (u *Updater) needsUpdate(current, latest string) bool {
	// Remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Handle dev version
	if current == "dev" {
		return true
	}

	// Simple version comparison (works for semantic versioning)
	return current != latest
}

// Update downloads and installs the latest version
func (u *Updater) Update(version string) error {
	// Determine platform and architecture
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Build binary name
	binaryName := fmt.Sprintf("gitai-%s-%s", platform, arch)
	if platform == "windows" {
		binaryName += ".exe"
	}

	// Download URLs
	binaryURL := fmt.Sprintf("%s/%s/%s", githubRawURL, version, binaryName)
	checksumURL := fmt.Sprintf("%s/%s/checksums.txt", githubRawURL, version)

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gitai-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download binary
	binaryPath := filepath.Join(tempDir, binaryName)
	if err := u.downloadFile(binaryURL, binaryPath); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}

	// Download checksums
	checksumPath := filepath.Join(tempDir, "checksums.txt")
	if err := u.downloadFile(checksumURL, checksumPath); err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}

	// Verify checksum
	if err := u.verifyChecksum(binaryPath, checksumPath, binaryName); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	// Get current executable path
	currentExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Resolve symlinks if any (e.g., /usr/local/bin/gitai might be a symlink)
	currentExePath, err = filepath.EvalSymlinks(currentExePath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Create backup of current binary
	backupPath := currentExePath + ".backup"
	if err := u.copyFile(currentExePath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Make new binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		os.Remove(backupPath)
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Replace current binary
	if err := u.replaceFile(binaryPath, currentExePath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentExePath)
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Remove backup on success
	os.Remove(backupPath)

	return nil
}

// downloadFile downloads a file from URL to filepath
func (u *Updater) downloadFile(url, filepath string) error {
	client := &http.Client{Timeout: 5 * time.Minute}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// verifyChecksum verifies the SHA256 checksum of a file
func (u *Updater) verifyChecksum(filePath, checksumPath, fileName string) error {
	// Calculate file hash
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	calculatedHash := hex.EncodeToString(hash.Sum(nil))

	// Read checksums file
	checksumData, err := os.ReadFile(checksumPath)
	if err != nil {
		return err
	}

	// Find matching checksum
	lines := strings.Split(string(checksumData), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == fileName {
			expectedHash := parts[0]
			if calculatedHash != expectedHash {
				return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, calculatedHash)
			}
			return nil
		}
	}

	return fmt.Errorf("checksum not found for %s", fileName)
}

// copyFile copies a file from src to dst
func (u *Updater) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// replaceFile replaces dst with src atomically
func (u *Updater) replaceFile(src, dst string) error {
	// On Unix-like systems, we can rename atomically
	// On Windows, we need to remove first
	if runtime.GOOS == "windows" {
		if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return os.Rename(src, dst)
}
