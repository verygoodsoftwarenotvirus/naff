package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffAwareWrite(t *testing.T) {
	t.Parallel()

	t.Run("creates new file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "subdir", "test.generated.go")

		status, err := diffAwareWrite(path, []byte("hello"))
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

		status, err := diffAwareWrite(path, []byte("hello"))
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

		status, err := diffAwareWrite(path, []byte("new"))
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
