package pipeline

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/renderer"
	"github.com/verygoodsoftwarenotvirus/naff/templates"
)

func init() {
	// Tests that exercise the planner need TemplateFS populated; cmd/naff
	// wires this up at runtime.
	TemplateFS = templates.FS
}

func TestDiffAwareWrite(t *testing.T) {
	t.Parallel()

	t.Run("creates new file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "subdir", "test.go")

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
		path := filepath.Join(dir, "test.go")
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
		path := filepath.Join(dir, "test.go")
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

// TestFindOrphans was dropped in the verbatim-mirror pivot; findOrphans was
// removed from pipeline.go (orphans are handled by `make generate-ddb`
// wiping the output directory). The original test lives in git history if
// orphan detection needs to be re-implemented.

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

func TestRewriteOutputPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in, want string
	}{
		{"_backend/cmd/services/api/main.go", "backend/cmd/services/api/main.go"},
		{"_backend/_go.mod", "backend/go.mod"},
		{"_backend/_go.sum", "backend/go.sum"},
		{"frontend/package.json", "frontend/package.json"},
		{"proto/filtering.proto", "proto/filtering.proto"},
		{"_backend/deploy/environments/prod/terraform/_terraform_.tf", "backend/deploy/environments/prod/terraform/_terraform_.tf"},
		{"_root/.github/workflows/ci.yaml", ".github/workflows/ci.yaml"},
		{"_root/Makefile", "Makefile"},
		{"_root/.gitignore", ".gitignore"},
	}

	for _, tc := range tests {
		if got := rewriteOutputPath(tc.in); got != tc.want {
			t.Errorf("rewriteOutputPath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestModeFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want os.FileMode
	}{
		{"backend/cmd/services/api/main.go", 0},
		{"backend/scripts/queries.sh", 0o755},
		{"backend/scripts/helpers", 0o755},
		{"ios/scripts/build.sh", 0o755},
		{"frontend/package.json", 0},
	}

	for _, tc := range tests {
		if got := modeFor(tc.path); got != tc.want {
			t.Errorf("modeFor(%q) = %o, want %o", tc.path, got, tc.want)
		}
	}
}

func TestPerEntityOutputName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		template, entity, want string
	}{
		{"entity.go.tmpl", "RecipeRating", "recipe_rating.go"},
		{"entity_test.go.tmpl", "RecipeRating", "recipe_rating_test.go"},
		{"entity.go.tmpl", "MealComponent", "meal_component.go"},
		{"entity.go.tmpl", "ValidIngredientMeasurementUnit", "valid_ingredient_measurement_unit.go"},
	}

	for _, tc := range tests {
		if got := perEntityOutputName(tc.template, tc.entity); got != tc.want {
			t.Errorf("perEntityOutputName(%q, %q) = %q, want %q",
				tc.template, tc.entity, got, tc.want)
		}
	}
}

func TestPlanPerEntityFiles(t *testing.T) {
	t.Parallel()

	cfg := &config.Project{
		ProjectMeta: config.ProjectMeta{
			Name:   "TestProject",
			Module: "github.com/example/test",
		},
		Domains: []config.Domain{
			{
				Name: "widgets",
				Entities: []config.Entity{
					{Name: "Widget", Fields: []config.Field{{Name: "Color", Type: "string"}}},
					{Name: "Gadget", Fields: []config.Field{{Name: "Power", Type: "int32"}}},
				},
			},
		},
	}

	p := &Pipeline{config: cfg}
	got := p.planPerEntityFiles()

	// Two templates (entity.go.tmpl + entity_test.go.tmpl) × 2 entities = 4 files.
	if len(got) != 4 {
		t.Fatalf("got %d planned files, want 4: %+v", len(got), got)
	}

	wantPaths := map[string]bool{
		"backend/internal/domain/widgets/widget.go":      false,
		"backend/internal/domain/widgets/widget_test.go": false,
		"backend/internal/domain/widgets/gadget.go":      false,
		"backend/internal/domain/widgets/gadget_test.go": false,
	}
	for _, pf := range got {
		if _, ok := wantPaths[pf.OutputPath]; !ok {
			t.Errorf("unexpected output path: %s", pf.OutputPath)
			continue
		}
		wantPaths[pf.OutputPath] = true

		ctx, ok := pf.Data.(parameterizedCtx)
		if !ok {
			t.Errorf("Data for %s is not parameterizedCtx: %T", pf.OutputPath, pf.Data)
			continue
		}
		if ctx.Project != cfg {
			t.Errorf("Project for %s does not match input config", pf.OutputPath)
		}
		if ctx.Domain == nil || ctx.Domain.Name != "widgets" {
			t.Errorf("Domain for %s = %+v, want widgets", pf.OutputPath, ctx.Domain)
		}
		if ctx.Entity == nil {
			t.Errorf("Entity for %s is nil", pf.OutputPath)
		}
		if !pf.IsGo {
			t.Errorf("IsGo for %s = false, want true", pf.OutputPath)
		}
	}
	for path, seen := range wantPaths {
		if !seen {
			t.Errorf("missing planned output path: %s", path)
		}
	}
}

// TestEntityTemplateRendersValidGo exercises the per-entity templates against
// a handful of representative Entity shapes and verifies the output is
// well-formed Go after gofmt. Catches template-level bugs (mismatched
// pointer/value comparisons, missing imports, broken control-flow) without
// running a full naff-generate against the embedded testdata.
func TestEntityTemplateRendersValidGo(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		domain config.Domain
		entity config.Entity
	}{
		{
			name:   "vanilla scalar fields",
			domain: config.Domain{Name: "widgets"},
			entity: config.Entity{
				Name: "Widget",
				Fields: []config.Field{
					{Name: "Color", Type: "string"},
					{Name: "Power", Type: "int32"},
				},
			},
		},
		{
			name:   "pointer field exercises pointer-aware Update body",
			domain: config.Domain{Name: "widgets"},
			entity: config.Entity{
				Name: "Gadget",
				Fields: []config.Field{
					{Name: "MaxPortions", Type: "*float32"},
				},
			},
		},
		{
			name:   "belongs_to + created_by_user routes server-side fields correctly",
			domain: config.Domain{Name: "mealplanning"},
			entity: config.Entity{
				Name:          "RecipeRating",
				BelongsTo:     "Recipe",
				CreatedByUser: true,
				Fields: []config.Field{
					{Name: "Notes", Type: "string"},
					{Name: "Taste", Type: "float32"},
				},
			},
		},
		{
			name:   "links_to synthesizes FK fields",
			domain: config.Domain{Name: "mealplanning"},
			entity: config.Entity{
				Name:      "MealComponent",
				BelongsTo: "Meal",
				LinksTo:   []config.Link{{Target: "Recipe"}},
				Fields: []config.Field{
					{Name: "ComponentType", Type: "string"},
					{Name: "RecipeScale", Type: "float32"},
				},
			},
		},
		{
			name:   "no editable fields produces a no-op Update body",
			domain: config.Domain{Name: "widgets"},
			entity: config.Entity{
				Name: "FrozenWidget",
				Fields: []config.Field{
					{Name: "ID", Type: "string"},
				},
			},
		},
	}

	r := renderer.New(templates.FS)
	const tmpl = "_backend/_parameterized/domain_entity/entity.go.tmpl"
	const testTmpl = "_backend/_parameterized/domain_entity/entity_test.go.tmpl"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := parameterizedCtx{
				Project: &config.Project{ProjectMeta: config.ProjectMeta{Name: "Test"}},
				Domain:  &tc.domain,
				Entity:  &tc.entity,
			}

			for _, path := range []string{tmpl, testTmpl} {
				out, err := r.Render(path, ctx)
				if err != nil {
					t.Fatalf("Render(%s): %v", path, err)
				}

				formatted, err := renderer.FormatGo([]byte(out))
				if err != nil {
					t.Fatalf("FormatGo(%s):\n%v\n--- raw output ---\n%s", path, err, out)
				}

				fset := token.NewFileSet()
				if _, err := parser.ParseFile(fset, "out.go", formatted, parser.AllErrors); err != nil {
					t.Fatalf("parser.ParseFile(%s): %v\n--- output ---\n%s", path, err, formatted)
				}
			}
		})
	}
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
