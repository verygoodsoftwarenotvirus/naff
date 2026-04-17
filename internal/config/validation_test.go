package config

import (
	"strings"
	"testing"
)

func boolPtr(b bool) *bool { return &b }

func TestValidate(t *testing.T) {
	t.Parallel()

	validProject := func() *Project {
		return &Project{
			ProjectMeta: ProjectMeta{
				Name:   "Test",
				Module: "github.com/example/test",
			},
			Domains: []Domain{
				{
					Name: "widgets",
					Entities: []Entity{
						{
							Name: "Widget",
							Fields: []Field{
								{Name: "Name", Type: "string"},
							},
						},
					},
				},
			},
		}
	}

	t.Run("accepts valid config", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		if err := p.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("rejects empty project name", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.ProjectMeta.Name = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "project.name") {
			t.Errorf("expected error about project.name, got: %v", err)
		}
	})

	t.Run("rejects empty module when targets default", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.ProjectMeta.Module = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "project.module") {
			t.Errorf("expected error about project.module, got: %v", err)
		}
	})

	t.Run("rejects no domains", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains = nil
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "at least one domain") {
			t.Errorf("expected error about domains, got: %v", err)
		}
	})

	t.Run("rejects empty domain name", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Name = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "domain name") {
			t.Errorf("expected error about domain name, got: %v", err)
		}
	})

	t.Run("rejects domain with no entities", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities = nil
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "at least one entity") {
			t.Errorf("expected error about entities, got: %v", err)
		}
	})

	t.Run("rejects empty entity name", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].Name = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "entity name") {
			t.Errorf("expected error about entity name, got: %v", err)
		}
	})

	t.Run("rejects duplicate entity names", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains = append(p.Domains, Domain{
			Name: "other",
			Entities: []Entity{
				{
					Name:   "Widget", // duplicate
					Fields: []Field{{Name: "X", Type: "string"}},
				},
			},
		})
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "duplicate entity") {
			t.Errorf("expected error about duplicate, got: %v", err)
		}
	})

	t.Run("rejects entity with no fields", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].Fields = nil
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "at least one field") {
			t.Errorf("expected error about fields, got: %v", err)
		}
	})

	t.Run("rejects empty field name", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].Fields[0].Name = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "field name") {
			t.Errorf("expected error about field name, got: %v", err)
		}
	})

	t.Run("rejects empty field type", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].Fields[0].Type = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "type is required") {
			t.Errorf("expected error about field type, got: %v", err)
		}
	})

	t.Run("rejects unsupported field type", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].Fields[0].Type = "map[string]string"
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "unsupported type") {
			t.Errorf("expected error about unsupported type, got: %v", err)
		}
	})

	t.Run("accepts pointer types", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].Fields[0].Type = "*string"
		if err := p.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("accepts all valid types", func(t *testing.T) {
		t.Parallel()

		for _, typ := range []string{"bool", "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "time.Time"} {
			p := validProject()
			p.Domains[0].Entities[0].Fields[0].Type = typ
			if err := p.Validate(); err != nil {
				t.Errorf("type %q should be valid, got: %v", typ, err)
			}
		}
	})

	t.Run("rejects unknown belongs_to reference", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities[0].BelongsTo = "Nonexistent"
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "unknown entity") {
			t.Errorf("expected error about unknown entity, got: %v", err)
		}
	})

	t.Run("accepts valid belongs_to reference", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Domains[0].Entities = append(p.Domains[0].Entities, Entity{
			Name:      "WidgetPart",
			BelongsTo: "Widget",
			Fields:    []Field{{Name: "Label", Type: "string"}},
		})
		if err := p.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("defaults backend to true when unset; iOS stays false", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		// Targets zero-valued: Backend nil, IOS false
		if err := p.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !p.Targets.BackendEnabled() {
			t.Error("expected backend target to default to true")
		}
		if p.Targets.IOS {
			t.Error("expected iOS target to stay false when unset")
		}
		if !p.Targets.AsyncMessageHandlerEnabled() {
			t.Error("expected async_message_handler target to default to true")
		}
	})

	t.Run("respects explicit target selection", func(t *testing.T) {
		t.Parallel()

		trueVal := true
		p := validProject()
		p.Targets.Backend = &trueVal
		p.Targets.IOS = false
		if err := p.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !p.Targets.BackendEnabled() {
			t.Error("expected backend target to remain true")
		}
		if p.Targets.IOS {
			t.Error("expected iOS target to remain false")
		}
	})

	t.Run("respects explicit backend opt-out", func(t *testing.T) {
		t.Parallel()

		falseVal := false
		p := validProject()
		p.Targets.Backend = &falseVal
		p.Targets.IOS = true
		p.ProjectMeta.Module = ""
		if err := p.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if p.Targets.BackendEnabled() {
			t.Error("expected explicit backend:false to stay false")
		}
	})

	t.Run("respects explicit async_message_handler opt-out", func(t *testing.T) {
		t.Parallel()

		falseVal := false
		p := validProject()
		p.Targets.AsyncMessageHandler = &falseVal
		if err := p.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Targets.AsyncMessageHandlerEnabled() {
			t.Error("expected explicit async_message_handler:false to stay false")
		}
	})

	t.Run("defaults IOSModuleName to project name", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Targets.IOS = true
		if err := p.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ProjectMeta.IOSModuleName != "Test" {
			t.Errorf("expected IOSModuleName to default to %q, got %q", "Test", p.ProjectMeta.IOSModuleName)
		}
	})

	t.Run("defaults IOSBundleID from project name", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Targets.IOS = true
		if err := p.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ProjectMeta.IOSBundleID != "com.example.test" {
			t.Errorf("expected IOSBundleID to default to %q, got %q", "com.example.test", p.ProjectMeta.IOSBundleID)
		}
	})

	t.Run("preserves explicit iOS fields", func(t *testing.T) {
		t.Parallel()

		p := validProject()
		p.Targets.IOS = true
		p.ProjectMeta.IOSBundleID = "com.custom.app"
		p.ProjectMeta.IOSModuleName = "CustomApp"
		if err := p.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ProjectMeta.IOSBundleID != "com.custom.app" {
			t.Errorf("expected IOSBundleID %q, got %q", "com.custom.app", p.ProjectMeta.IOSBundleID)
		}
		if p.ProjectMeta.IOSModuleName != "CustomApp" {
			t.Errorf("expected IOSModuleName %q, got %q", "CustomApp", p.ProjectMeta.IOSModuleName)
		}
	})

	t.Run("rejects empty module when backend target is enabled", func(t *testing.T) {
		t.Parallel()

		trueVal := true
		p := validProject()
		p.Targets.Backend = &trueVal
		p.Targets.IOS = false
		p.ProjectMeta.Module = ""
		err := p.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "project.module") {
			t.Errorf("expected error about project.module, got: %v", err)
		}
	})

	t.Run("allows empty module when only iOS target is enabled", func(t *testing.T) {
		t.Parallel()

		falseVal := false
		p := validProject()
		p.Targets.Backend = &falseVal
		p.Targets.IOS = true
		p.ProjectMeta.Module = ""
		if err := p.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestCheckOwnershipCycles(t *testing.T) {
	t.Parallel()

	t.Run("detects simple cycle", func(t *testing.T) {
		t.Parallel()

		domains := []Domain{
			{
				Name: "test",
				Entities: []Entity{
					{Name: "A", BelongsTo: "B", Fields: []Field{{Name: "X", Type: "string"}}},
					{Name: "B", BelongsTo: "A", Fields: []Field{{Name: "X", Type: "string"}}},
				},
			},
		}

		err := checkOwnershipCycles(domains)
		if err == nil {
			t.Fatal("expected cycle error")
		}
		if !strings.Contains(err.Error(), "cycle") {
			t.Errorf("expected cycle error, got: %v", err)
		}
	})

	t.Run("detects three-way cycle", func(t *testing.T) {
		t.Parallel()

		domains := []Domain{
			{
				Name: "test",
				Entities: []Entity{
					{Name: "A", BelongsTo: "B", Fields: []Field{{Name: "X", Type: "string"}}},
					{Name: "B", BelongsTo: "C", Fields: []Field{{Name: "X", Type: "string"}}},
					{Name: "C", BelongsTo: "A", Fields: []Field{{Name: "X", Type: "string"}}},
				},
			},
		}

		err := checkOwnershipCycles(domains)
		if err == nil {
			t.Fatal("expected cycle error")
		}
	})

	t.Run("allows valid chain", func(t *testing.T) {
		t.Parallel()

		domains := []Domain{
			{
				Name: "test",
				Entities: []Entity{
					{Name: "Forum", Fields: []Field{{Name: "X", Type: "string"}}},
					{Name: "Thread", BelongsTo: "Forum", Fields: []Field{{Name: "X", Type: "string"}}},
					{Name: "Post", BelongsTo: "Thread", Fields: []Field{{Name: "X", Type: "string"}}},
				},
			},
		}

		if err := checkOwnershipCycles(domains); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestFieldMethods(t *testing.T) {
	t.Parallel()

	t.Run("IsPointer", func(t *testing.T) {
		t.Parallel()

		if !(Field{Type: "*string"}).IsPointer() {
			t.Error("*string should be a pointer")
		}
		if (Field{Type: "string"}).IsPointer() {
			t.Error("string should not be a pointer")
		}
		if (Field{Type: ""}).IsPointer() {
			t.Error("empty type should not be a pointer")
		}
	})

	t.Run("BaseType", func(t *testing.T) {
		t.Parallel()

		if got := (Field{Type: "*string"}).BaseType(); got != "string" {
			t.Errorf("expected string, got %q", got)
		}
		if got := (Field{Type: "int64"}).BaseType(); got != "int64" {
			t.Errorf("expected int64, got %q", got)
		}
	})

	t.Run("defaults are all true", func(t *testing.T) {
		t.Parallel()

		f := Field{Name: "X", Type: "string"}
		if !f.IsRequired() {
			t.Error("expected required default true")
		}
		if !f.IsCreatable() {
			t.Error("expected creatable default true")
		}
		if !f.IsEditable() {
			t.Error("expected editable default true")
		}
	})

	t.Run("explicit false overrides", func(t *testing.T) {
		t.Parallel()

		f := Field{
			Name:      "X",
			Type:      "string",
			Required:  boolPtr(false),
			Creatable: boolPtr(false),
			Editable:  boolPtr(false),
		}
		if f.IsRequired() {
			t.Error("expected required false")
		}
		if f.IsCreatable() {
			t.Error("expected creatable false")
		}
		if f.IsEditable() {
			t.Error("expected editable false")
		}
	})
}
