package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("parses valid config", func(t *testing.T) {
		t.Parallel()

		input := []byte(`
project:
  name: "TestProject"
  module: "github.com/example/test"
  platform_module: "github.com/example/platform"

features:
  webhooks: true
  payments: false

domains:
  - name: "widgets"
    entities:
      - name: "Widget"
        belongs_to_account: true
        created_by_user: true
        fields:
          - name: "Name"
            type: "string"
          - name: "Count"
            type: "int"
            required: false
          - name: "Active"
            type: "bool"
            creatable: false
`)

		cfg, err := Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ProjectMeta.Name != "TestProject" {
			t.Errorf("expected name TestProject, got %q", cfg.ProjectMeta.Name)
		}
		if cfg.ProjectMeta.Module != "github.com/example/test" {
			t.Errorf("expected module github.com/example/test, got %q", cfg.ProjectMeta.Module)
		}
		if cfg.ProjectMeta.PlatformModule != "github.com/example/platform" {
			t.Errorf("expected platform_module github.com/example/platform, got %q", cfg.ProjectMeta.PlatformModule)
		}

		if len(cfg.Domains) != 1 {
			t.Fatalf("expected 1 domain, got %d", len(cfg.Domains))
		}
		d := cfg.Domains[0]
		if d.Name != "widgets" {
			t.Errorf("expected domain name widgets, got %q", d.Name)
		}
		if len(d.Entities) != 1 {
			t.Fatalf("expected 1 entity, got %d", len(d.Entities))
		}

		e := d.Entities[0]
		if e.Name != "Widget" {
			t.Errorf("expected entity name Widget, got %q", e.Name)
		}
		if !e.BelongsToAccount {
			t.Error("expected belongs_to_account to be true")
		}
		if !e.CreatedByUser {
			t.Error("expected created_by_user to be true")
		}
		if len(e.Fields) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(e.Fields))
		}

		// Name field: all defaults (required=true, creatable=true, editable=true)
		nameField := e.Fields[0]
		if nameField.Name != "Name" {
			t.Errorf("expected field Name, got %q", nameField.Name)
		}
		if !nameField.IsRequired() {
			t.Error("Name should be required by default")
		}
		if !nameField.IsCreatable() {
			t.Error("Name should be creatable by default")
		}
		if !nameField.IsEditable() {
			t.Error("Name should be editable by default")
		}

		// Count field: required=false override
		countField := e.Fields[1]
		if countField.IsRequired() {
			t.Error("Count should not be required (explicitly set false)")
		}
		if !countField.IsCreatable() {
			t.Error("Count should be creatable by default")
		}

		// Active field: creatable=false override
		activeField := e.Fields[2]
		if !activeField.IsRequired() {
			t.Error("Active should be required by default")
		}
		if activeField.IsCreatable() {
			t.Error("Active should not be creatable (explicitly set false)")
		}
	})

	t.Run("parses features with defaults", func(t *testing.T) {
		t.Parallel()

		input := []byte(`
project:
  name: "Test"
  module: "github.com/example/test"
domains:
  - name: "x"
    entities:
      - name: "X"
        fields:
          - name: "Val"
            type: "string"
`)

		cfg, err := Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// All feature pointers should be nil (unset)
		if cfg.Features.Webhooks != nil {
			t.Error("expected webhooks to be nil (unset)")
		}

		// FeatureEnabled with nil should return default
		if !cfg.Features.FeatureEnabled(cfg.Features.Webhooks, true) {
			t.Error("expected unset feature to default to true")
		}
		if cfg.Features.FeatureEnabled(cfg.Features.Webhooks, false) {
			t.Error("expected unset feature to default to false")
		}
	})

	t.Run("returns error on invalid YAML", func(t *testing.T) {
		t.Parallel()

		_, err := Parse([]byte("{{invalid"))
		if err == nil {
			t.Error("expected error on invalid YAML")
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	t.Run("loads valid file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "project.yaml")

		content := []byte(`
project:
  name: "FileTest"
  module: "github.com/example/filetest"
domains:
  - name: "things"
    entities:
      - name: "Thing"
        fields:
          - name: "Value"
            type: "string"
`)
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatalf("writing test file: %v", err)
		}

		cfg, err := LoadFromFile(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.ProjectMeta.Name != "FileTest" {
			t.Errorf("expected name FileTest, got %q", cfg.ProjectMeta.Name)
		}
	})

	t.Run("returns error for missing file", func(t *testing.T) {
		t.Parallel()

		_, err := LoadFromFile("/nonexistent/path.yaml")
		if err == nil {
			t.Error("expected error for missing file")
		}
	})
}
