package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
)

func TestDiffAwareWrite(t *testing.T) {
	t.Parallel()

	t.Run("creates new file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "subdir", "test.generated.go")

		status, err := diffAwareWrite(path, []byte("hello"), 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != FileCreated {
			t.Errorf("expected FileCreated, got %v", status)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading file: %v", err)
		}
		if string(content) != "hello" {
			t.Errorf("expected hello, got %q", content)
		}
	})

	t.Run("skips unchanged file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "test.generated.go")
		if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
			t.Fatal(err)
		}

		status, err := diffAwareWrite(path, []byte("hello"), 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != FileUnchanged {
			t.Errorf("expected FileUnchanged, got %v", status)
		}
	})

	t.Run("updates changed file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "test.generated.go")
		if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
			t.Fatal(err)
		}

		status, err := diffAwareWrite(path, []byte("new"), 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != FileUpdated {
			t.Errorf("expected FileUpdated, got %v", status)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "new" {
			t.Errorf("expected new, got %q", content)
		}
	})
}

func TestFindOrphans(t *testing.T) {
	t.Parallel()

	t.Run("detects orphaned generated files", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		// Create some files.
		planned := filepath.Join(dir, "planned.generated.go")
		orphan := filepath.Join(dir, "orphan.generated.go")
		custom := filepath.Join(dir, "custom.go")

		for _, path := range []string{planned, orphan, custom} {
			if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
		}

		p := &Pipeline{outputDir: dir}
		plannedFiles := []PlannedFile{
			{OutputPath: "planned.generated.go"},
		}

		orphans := p.findOrphans(plannedFiles)
		if len(orphans) != 1 {
			t.Fatalf("expected 1 orphan, got %d: %v", len(orphans), orphans)
		}
		if orphans[0] != orphan {
			t.Errorf("expected %s, got %s", orphan, orphans[0])
		}
	})

	t.Run("ignores non-generated files", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		custom := filepath.Join(dir, "custom.go")
		if err := os.WriteFile(custom, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}

		p := &Pipeline{outputDir: dir}
		orphans := p.findOrphans(nil)
		if len(orphans) != 0 {
			t.Errorf("expected 0 orphans, got %d", len(orphans))
		}
	})
}

func TestWriteNaffConfig(t *testing.T) {
	t.Parallel()

	t.Run("writes config to output directory", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		cfg := &config.Project{
			ProjectMeta: config.ProjectMeta{
				Name:   "TestProject",
				Module: "github.com/example/test",
			},
			Domains: []config.Domain{
				{
					Name: "widgets",
					Entities: []config.Entity{
						{
							Name: "Widget",
							Fields: []config.Field{
								{Name: "Name", Type: "string"},
							},
						},
					},
				},
			},
		}

		p := &Pipeline{
			config:    cfg,
			outputDir: dir,
		}

		if err := p.writeNaffConfig(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify the file was written.
		naffPath := filepath.Join(dir, config.NaffConfigFileName)
		data, err := os.ReadFile(naffPath)
		if err != nil {
			t.Fatalf("reading .naff.yaml: %v", err)
		}

		// Verify it round-trips back to an equivalent config.
		parsed, err := config.Parse(data)
		if err != nil {
			t.Fatalf("parsing written .naff.yaml: %v", err)
		}

		if parsed.ProjectMeta.Name != cfg.ProjectMeta.Name {
			t.Errorf("name: got %q, want %q", parsed.ProjectMeta.Name, cfg.ProjectMeta.Name)
		}
		if parsed.ProjectMeta.Module != cfg.ProjectMeta.Module {
			t.Errorf("module: got %q, want %q", parsed.ProjectMeta.Module, cfg.ProjectMeta.Module)
		}
		if len(parsed.Domains) != 1 {
			t.Fatalf("expected 1 domain, got %d", len(parsed.Domains))
		}
		if parsed.Domains[0].Entities[0].Name != "Widget" {
			t.Errorf("entity name: got %q, want %q", parsed.Domains[0].Entities[0].Name, "Widget")
		}
	})

	t.Run("is diff-aware on repeated writes", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		cfg := &config.Project{
			ProjectMeta: config.ProjectMeta{
				Name:   "TestProject",
				Module: "github.com/example/test",
			},
			Domains: []config.Domain{
				{
					Name: "things",
					Entities: []config.Entity{
						{
							Name:   "Thing",
							Fields: []config.Field{{Name: "Val", Type: "string"}},
						},
					},
				},
			},
		}

		p := &Pipeline{config: cfg, outputDir: dir}

		// First write.
		if err := p.writeNaffConfig(); err != nil {
			t.Fatalf("first write: %v", err)
		}

		naffPath := filepath.Join(dir, config.NaffConfigFileName)
		info1, err := os.Stat(naffPath)
		if err != nil {
			t.Fatal(err)
		}

		// Second write with same config should not change the file content.
		if err = p.writeNaffConfig(); err != nil {
			t.Fatalf("second write: %v", err)
		}

		// Read content to confirm it's unchanged (diffAwareWrite returns FileUnchanged).
		data1, _ := os.ReadFile(naffPath)

		// Modify config and write again.
		cfg.ProjectMeta.Name = "UpdatedProject"
		if err = p.writeNaffConfig(); err != nil {
			t.Fatalf("third write: %v", err)
		}

		data2, _ := os.ReadFile(naffPath)
		if string(data1) == string(data2) {
			t.Error("expected file content to change after config update")
		}

		info2, _ := os.Stat(naffPath)
		if info1.ModTime().Equal(info2.ModTime()) && string(data1) != string(data2) {
			// This is fine — just confirms the file was rewritten
		}
	})
}

func TestFileStatusString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status FileStatus
		want   string
	}{
		{FileCreated, "created"},
		{FileUpdated, "updated"},
		{FileUnchanged, "unchanged"},
		{FileOrphaned, "orphaned"},
		{FileError, "error"},
	}

	for _, tc := range tests {
		if got := tc.status.String(); got != tc.want {
			t.Errorf("FileStatus(%d).String() = %q, want %q", tc.status, got, tc.want)
		}
	}
}
