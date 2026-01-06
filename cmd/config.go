package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gitai/internal/config"
)

var (
	configInit bool
	configShow bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gitai configuration",
	Long:  "View or initialize gitai configuration file",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVar(&configInit, "init", false, "Create default config file in current directory")
	configCmd.Flags().BoolVar(&configShow, "show", false, "Show current configuration")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if configInit {
		return initConfig()
	}

	if configShow {
		return showConfig()
	}

	// Default: show help
	return cmd.Help()
}

func initConfig() error {
	cfg := config.DefaultConfig()

	path := ".gitcommit.yaml"
	if err := cfg.Save(path); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("✅ Created default configuration at: %s\n", path)
	fmt.Println("\nEdit this file to customize commit types, scopes, and AI model settings.")

	return nil
}

func showConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	fmt.Println("Current Configuration:")
	fmt.Println("======================")
	fmt.Printf("Model: %s\n", cfg.Model)
	fmt.Printf("Language: %s\n", cfg.Language)
	fmt.Printf("Template: %s\n", cfg.Template)
	fmt.Println("\nCommit Types:")
	for _, t := range cfg.Types {
		fmt.Printf("  %s %s - %s\n", t.Emoji, t.Name, t.Desc)
	}

	if len(cfg.Scopes) > 0 {
		fmt.Println("\nScopes:")
		for _, s := range cfg.Scopes {
			fmt.Printf("  - %s\n", s)
		}
	}

	return nil
}
