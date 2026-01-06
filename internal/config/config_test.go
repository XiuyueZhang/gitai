package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Model == "" {
		t.Error("Default model should not be empty")
	}

	if cfg.Language == "" {
		t.Error("Default language should not be empty")
	}

	if len(cfg.Types) == 0 {
		t.Error("Default types should not be empty")
	}

	if cfg.Template == "" {
		t.Error("Default template should not be empty")
	}
}

func TestGetTypeByName(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name     string
		want     bool
		typeName string
	}{
		{"feat exists", true, "feat"},
		{"fix exists", true, "fix"},
		{"nonexistent", false, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cfg.GetTypeByName(tt.typeName)
			if tt.want && result == nil {
				t.Errorf("Expected to find type %s but got nil", tt.typeName)
			}
			if !tt.want && result != nil {
				t.Errorf("Expected not to find type %s but got %v", tt.typeName, result)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "gitcommit-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Create and save config
	cfg := DefaultConfig()
	cfg.Model = "test-model"
	cfg.Language = "zh"

	if err := cfg.Save(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := loadFromFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loaded.Model != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", loaded.Model)
	}

	if loaded.Language != "zh" {
		t.Errorf("Expected language 'zh', got '%s'", loaded.Language)
	}
}
