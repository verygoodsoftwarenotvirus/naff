package renderer

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
)

func TestGoZeroValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typ  string
		want string
	}{
		{"string", `""`},
		{"bool", "false"},
		{"int", "0"},
		{"int64", "0"},
		{"uint32", "0"},
		{"float64", "0"},
		{"*string", "nil"},
		{"*int", "nil"},
	}

	for _, tc := range tests {
		t.Run(tc.typ, func(t *testing.T) {
			t.Parallel()
			if got := goZeroValue(tc.typ); got != tc.want {
				t.Errorf("goZeroValue(%q) = %q, want %q", tc.typ, got, tc.want)
			}
		})
	}
}

func TestSqlType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typ  string
		want string
	}{
		{"string", "TEXT"},
		{"bool", "BOOLEAN"},
		{"int", "INTEGER"},
		{"int64", "BIGINT"},
		{"int16", "SMALLINT"},
		{"uint32", "INTEGER"},
		{"float32", "REAL"},
		{"float64", "DOUBLE PRECISION"},
		{"*string", "TEXT"},
		{"*int64", "BIGINT"},
	}

	for _, tc := range tests {
		t.Run(tc.typ, func(t *testing.T) {
			t.Parallel()
			if got := sqlType(tc.typ); got != tc.want {
				t.Errorf("sqlType(%q) = %q, want %q", tc.typ, got, tc.want)
			}
		})
	}
}

func TestProtoType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typ  string
		want string
	}{
		{"string", "string"},
		{"bool", "bool"},
		{"int32", "int32"},
		{"int64", "int64"},
		{"uint32", "uint32"},
		{"uint64", "uint64"},
		{"float32", "float"},
		{"float64", "double"},
		{"*string", "string"},
	}

	for _, tc := range tests {
		t.Run(tc.typ, func(t *testing.T) {
			t.Parallel()
			if got := protoType(tc.typ); got != tc.want {
				t.Errorf("protoType(%q) = %q, want %q", tc.typ, got, tc.want)
			}
		})
	}
}

func boolPtr(b bool) *bool { return &b }

func TestFieldFilters(t *testing.T) {
	t.Parallel()

	fields := []config.Field{
		{Name: "Name", Type: "string"},                                               // all defaults: required, creatable, editable
		{Name: "Secret", Type: "string", Creatable: boolPtr(false)},                  // not creatable
		{Name: "Immutable", Type: "string", Editable: boolPtr(false)},                // not editable
		{Name: "Optional", Type: "*string", Required: boolPtr(false)},                // not required
		{Name: "ReadOnly", Type: "string", Creatable: boolPtr(false), Editable: boolPtr(false)}, // neither
	}

	t.Run("creatableFields", func(t *testing.T) {
		t.Parallel()
		result := creatableFields(fields)
		if len(result) != 3 {
			t.Errorf("expected 3 creatable fields, got %d", len(result))
		}
	})

	t.Run("editableFields", func(t *testing.T) {
		t.Parallel()
		result := editableFields(fields)
		if len(result) != 3 {
			t.Errorf("expected 3 editable fields, got %d", len(result))
		}
	})

	t.Run("requiredFields", func(t *testing.T) {
		t.Parallel()
		result := requiredFields(fields)
		if len(result) != 4 {
			t.Errorf("expected 4 required fields, got %d", len(result))
		}
	})
}

func TestJsonTag(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		f := config.Field{Name: "IssueType", Type: "string"}
		if got := jsonTag(f); got != "issueType" {
			t.Errorf("expected issueType, got %q", got)
		}
	})

	t.Run("with omitempty", func(t *testing.T) {
		t.Parallel()
		f := config.Field{Name: "RelevantTable", Type: "string", Omitempty: true}
		if got := jsonTag(f); got != "relevantTable,omitempty" {
			t.Errorf("expected relevantTable,omitempty, got %q", got)
		}
	})
}
