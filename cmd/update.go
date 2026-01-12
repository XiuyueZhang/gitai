package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xyue92/gitai/internal/updater"
)

var (
	checkOnly bool
	forceFlag bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update gitai to the latest version",
	Long: `Check for updates and automatically install the latest version of gitai.

The update command will:
  1. Check GitHub releases for the latest version
  2. Download the appropriate binary for your platform
  3. Verify the binary using SHA256 checksums
  4. Replace the current binary with the new version

Example:
  gitai update              # Update to latest version
  gitai update --check      # Only check for updates without installing
  gitai update --force      # Force update even if already on latest version`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "Only check for updates without installing")
	updateCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force update even if already on latest version")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := rootCmd.Version
	if currentVersion == "" {
		currentVersion = "dev"
	}

	// Create updater
	upd := updater.New(currentVersion)

	// Show current version
	fmt.Printf("Current version: %s\n\n", color.CyanString(currentVersion))

	// Check for updates
	fmt.Println("Checking for updates...")
	release, needsUpdate, err := upd.CheckForUpdate()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !needsUpdate && !forceFlag {
		fmt.Println(color.GreenString("✓ You are already on the latest version!"))
		return nil
	}

	// Display update information
	fmt.Printf("\nLatest version: %s\n", color.GreenString(release.TagName))

	if needsUpdate {
		fmt.Println(color.YellowString("→ A new version is available!"))
	} else {
		fmt.Println(color.YellowString("→ Forcing update to reinstall current version"))
	}

	// If check-only mode, exit here
	if checkOnly {
		if needsUpdate {
			fmt.Println("\nRun 'gitai update' to install the latest version.")
		}
		return nil
	}

	// Confirm update
	if !forceFlag {
		fmt.Printf("\nDo you want to update to %s? [Y/n] ", release.TagName)
		var response string
		fmt.Scanln(&response)

		if response != "" && response != "y" && response != "Y" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	// Perform update
	fmt.Println("\nDownloading update...")
	if err := upd.Update(release.TagName); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// Success message
	fmt.Println()
	fmt.Println(color.GreenString("✓ Successfully updated to %s!", release.TagName))
	fmt.Println("\nRun 'gitai --version' to verify the new version.")

	return nil
}

// SetVersion sets the version for the root command
func SetVersion(version string) {
	if version == "" {
		version = "dev"
	}
	rootCmd.Version = version
}

// GetVersion returns the current version
func GetVersion() string {
	if rootCmd.Version == "" {
		return "dev"
	}
	return rootCmd.Version
}
