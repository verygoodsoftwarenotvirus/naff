package renderer

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/verygoodsoftwarenotvirus/naff/internal/config"
	"github.com/verygoodsoftwarenotvirus/naff/internal/naming"
)

// FuncMap returns the template function map used by all NAFF templates.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// Naming transformations (take a string, return a string)
		"singular":           func(s string) string { return naming.New(s).Singular() },
		"plural":             func(s string) string { return naming.New(s).Plural() },
		"snake":              func(s string) string { return naming.New(s).Snake() },
		"pluralSnake":        func(s string) string { return naming.New(s).PluralSnake() },
		"camel":              func(s string) string { return naming.New(s).Camel() },
		"pluralCamel":        func(s string) string { return naming.New(s).PluralCamel() },
		"kebab":              func(s string) string { return naming.New(s).Kebab() },
		"pluralKebab":        func(s string) string { return naming.New(s).PluralKebab() },
		"packageName":        func(s string) string { return naming.New(s).PackageName() },
		"abbreviation":       func(s string) string { return naming.New(s).Abbreviation() },
		"lowerAbbreviation":  func(s string) string { return naming.New(s).LowerAbbreviation() },
		"humanReadable":      func(s string) string { return naming.New(s).HumanReadable() },
		"pluralHumanReadable": func(s string) string { return naming.New(s).PluralHumanReadable() },
		"article":            func(s string) string { return naming.New(s).Article() },
		"withArticle":        func(s string) string { return naming.New(s).WithArticle() },

		// String utilities
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      strings.Title, //nolint:staticcheck
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
		"join":       strings.Join,
		"replace":    strings.ReplaceAll,

		// Go type helpers
		"goZeroValue": goZeroValue,
		"isPointer":   func(f config.Field) bool { return f.IsPointer() },
		"baseType":    func(f config.Field) string { return f.BaseType() },
		"sqlType":     sqlType,
		"protoType":   protoType,
		"fakeValue":             fakeValue,
		"fakeValueBase":         fakeValueBase,
		"pointerConversionFunc": pointerConversionFunc,
		"readConversionFunc":    readConversionFunc,
		"writeConversionFunc":   writeConversionFunc,

		// protoFieldName converts a Go field identifier to the name protoc-gen-go
		// emits from the snake-cased proto field. protoc lowercases acronyms,
		// so `ReferencedID` → `referenced_id` → `ReferencedId` (not `ReferencedID`).
		// Use in grpc/converters templates when referencing fields on generated
		// .pb.go structs.
		"protoFieldName": func(s string) string {
			n := naming.New(s)
			return strcase.ToCamel(n.Snake())
		},

		// Field filtering
		"creatableFields": creatableFields,
		"editableFields":  editableFields,
		"requiredFields":  requiredFields,
		"columnFields":    columnFields,

		// Entity ordering: returns a new slice sorted by Name (A→Z) so
		// templates that iterate entities emit deterministic, config-order
		// -independent output.
		"sortedEntities": sortedEntities,

		// Parent-chain helpers: given (domain, entity), return the chain of
		// ancestor entities reached by walking `belongs_to`, outermost-first.
		// Used by per-entity templates to synthesize DataManager method
		// signatures (Exists/Get/List/Archive) that take parent IDs in path
		// order.
		"parentChain":    parentChain,
		"existsArgNames": existsArgNames,
		"listArgNames":   listArgNames,
		"archiveArgNames": archiveArgNames,

		// Nested-collection helpers: enumerate the children of an entity
		// that are marked `nested: true`, and derive the collection field
		// name (plural, parent-prefix-stripped) on the parent struct.
		"nestedChildren":          nestedChildren,
		"nestedCollectionField":   nestedCollectionField,

		// jsonTagForStruct is like jsonTag but also appends ",omitempty"
		// for pointer-typed fields (mirroring the upstream convention on
		// domain structs, where *T fields carry omitempty even when the
		// schema does not explicitly set Omitempty).
		"jsonTagForStruct": jsonTagForStruct,

		// JSON tag helpers
		"jsonTag": jsonTag,

		// sharedColumn reports whether a snake_case field name is one of the
		// shared column constants defined in cmd/tools/codegen/queries/helpers.go.
		// Per-entity query files skip redeclaring these to avoid collisions
		// across the package main namespace.
		"sharedColumn": sharedColumn,

		// Code generation helpers
		"importPath": func(module, pkg string) string {
			return fmt.Sprintf("%s/%s", module, pkg)
		},
		"repeat": strings.Repeat,

		// Conditional helpers
		"and": func(a, b bool) bool { return a && b },
		"or": func(vals ...bool) bool {
			for _, v := range vals {
				if v {
					return true
				}
			}
			return false
		},
		"not": func(a bool) bool { return !a },

		// Struct tag helpers
		"bt": func() string { return "`" },
		"tag": func(parts ...string) string {
			return "`" + strings.Join(parts, " ") + "`"
		},

		// Arithmetic
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },

		// Collection helpers
		"last": func(i, length int) bool { return i == length-1 },

		// TypeScript/frontend helpers
		"tsType":         tsType,
		"tsDefaultValue": tsDefaultValue,
		"htmlInputType":  htmlInputType,

		// Swift/iOS helpers
		"swiftType":      swiftType,
		"swiftZeroValue": swiftZeroValue,
	}
}

// goZeroValue returns the Go zero value for a type string.
func goZeroValue(typ string) string {
	if len(typ) > 0 && typ[0] == '*' {
		return "nil"
	}
	switch typ {
	case "bool":
		return "false"
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "0"
	case "time.Time":
		return "time.Time{}"
	default:
		return `""`
	}
}

// sqlType returns the PostgreSQL column type for a Go type.
func sqlType(typ string) string {
	base := typ
	if len(base) > 0 && base[0] == '*' {
		base = base[1:]
	}

	switch base {
	case "bool":
		return "BOOLEAN"
	case "string":
		return "TEXT"
	case "int", "int32":
		return "INTEGER"
	case "int8", "int16":
		return "SMALLINT"
	case "int64":
		return "BIGINT"
	case "uint", "uint32":
		return "INTEGER"
	case "uint8", "uint16":
		return "SMALLINT"
	case "uint64":
		return "BIGINT"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"
	case "time.Time":
		return "TIMESTAMP WITH TIME ZONE"
	default:
		return "TEXT"
	}
}

// protoType returns the Protobuf type for a Go type.
func protoType(typ string) string {
	base := typ
	if len(base) > 0 && base[0] == '*' {
		base = base[1:]
	}

	switch base {
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "int", "int32":
		return "int32"
	case "int8", "int16":
		return "int32"
	case "int64":
		return "int64"
	case "uint", "uint32":
		return "uint32"
	case "uint8", "uint16":
		return "uint32"
	case "uint64":
		return "uint64"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "time.Time":
		return "google.protobuf.Timestamp"
	default:
		return "string"
	}
}

// creatableFields returns only the fields that are included in creation input.
func creatableFields(fields []config.Field) []config.Field {
	var result []config.Field
	for _, f := range fields {
		if f.IsCreatable() {
			result = append(result, f)
		}
	}
	return result
}

// editableFields returns only the fields that are included in update input.
func editableFields(fields []config.Field) []config.Field {
	var result []config.Field
	for _, f := range fields {
		if f.IsEditable() {
			result = append(result, f)
		}
	}
	return result
}

// requiredFields returns only the fields that are required.
func requiredFields(fields []config.Field) []config.Field {
	var result []config.Field
	for _, f := range fields {
		if f.IsRequired() {
			result = append(result, f)
		}
	}
	return result
}

// columnFields returns all fields (all fields are columns by default).
func columnFields(fields []config.Field) []config.Field {
	return fields
}

// sortedEntities returns a copy of the given slice sorted by Name (A→Z).
// Callers mutate neither the input slice nor the underlying Entity values.
func sortedEntities(entities []config.Entity) []config.Entity {
	out := make([]config.Entity, len(entities))
	copy(out, entities)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// pointerConversionFunc returns the database helper function name for converting nullable DB types.
func pointerConversionFunc(f config.Field) string {
	base := f.BaseType()
	switch base {
	case "string":
		return "StringPointerFromNullString"
	case "bool":
		return "BoolPointerFromNullBool"
	case "int", "int32":
		return "Int32PointerFromNullInt32"
	case "int64":
		return "Int64PointerFromNullInt64"
	case "float32":
		return "Float32PointerFromNullString"
	case "float64":
		return "Float64PointerFromNullString"
	case "time.Time":
		return "TimePointerFromNullTime"
	default:
		return "StringPointerFromNullString"
	}
}

// readConversionFunc returns the database-to-domain conversion helper for a field.
// Returns "" when the field is required and non-pointer (passthrough, no conversion).
// Used by the repository template to pick the correct helper for each of {string, bool,
// int32, int64, float32, float64} × {required, non-required non-pointer, pointer}.
// Note: platform stores floats as NullString columns, so float helpers convert
// to/from NullString, not NullFloat32/NullFloat64.
func readConversionFunc(f config.Field) string {
	base := f.BaseType()
	if f.IsPointer() {
		switch base {
		case "string":
			return "StringPointerFromNullString"
		case "bool":
			return "BoolPointerFromNullBool"
		case "int", "int32":
			return "Int32PointerFromNullInt32"
		case "int64":
			return "Int64PointerFromNullInt64"
		case "float32":
			return "Float32PointerFromNullString"
		case "float64":
			return "Float64PointerFromNullString"
		case "time.Time":
			return "TimePointerFromNullTime"
		default:
			return "StringPointerFromNullString"
		}
	}
	if !f.IsRequired() {
		switch base {
		case "string":
			return "StringFromNullString"
		case "bool":
			return "BoolFromNullBool"
		case "int", "int32":
			return "Int32FromNullInt32"
		case "int64":
			return "Int64FromNullInt64"
		case "float32":
			return "Float32FromNullString"
		case "float64":
			return "Float64FromNullString"
		case "time.Time":
			return "TimeFromNullTime"
		default:
			return "StringFromNullString"
		}
	}
	return ""
}

// writeConversionFunc returns the domain-to-database conversion helper for a field.
// Returns "" when the field is required and non-pointer (passthrough).
// Note: platform stores floats as NullString columns, so float helpers return
// NullString, not NullFloat*.
func writeConversionFunc(f config.Field) string {
	base := f.BaseType()
	if f.IsPointer() {
		switch base {
		case "string":
			return "NullStringFromStringPointer"
		case "bool":
			return "NullBoolFromBoolPointer"
		case "int", "int32":
			return "NullInt32FromInt32Pointer"
		case "int64":
			return "NullInt64FromInt64Pointer"
		case "float32":
			return "NullStringFromFloat32Pointer"
		case "float64":
			return "NullStringFromFloat64Pointer"
		case "time.Time":
			return "NullTimeFromTimePointer"
		default:
			return "NullStringFromStringPointer"
		}
	}
	if !f.IsRequired() {
		switch base {
		case "string":
			return "NullStringFromString"
		case "bool":
			return "NullBoolFromBool"
		case "int", "int32":
			return "NullInt32FromInt32"
		case "int64":
			return "NullInt64FromInt64"
		case "float32":
			return "NullStringFromFloat32"
		case "float64":
			return "NullStringFromFloat64"
		case "time.Time":
			return "NullTimeFromTime"
		default:
			return "NullStringFromString"
		}
	}
	return ""
}

// fakeValue returns a gofakeit expression for generating fake data for a field.
// For pointer fields, wraps the value in a ptr() helper so callers can assign
// directly to struct literal fields. The fakes template must define a generic
// `ptr[T any](v T) *T` helper for this to resolve.
func fakeValue(f config.Field) string {
	if f.IsPointer() {
		return "ptr(" + fakeValueForType(f.BaseType()) + ")"
	}
	return fakeValueForType(f.Type)
}

// fakeValueBase returns the base (non-pointer) fake expression for a field,
// regardless of whether the field is a pointer. Use this in fakes templates
// where the caller takes `&var` separately (e.g. update-request inputs that
// declare an intermediate var).
func fakeValueBase(f config.Field) string {
	return fakeValueForType(f.BaseType())
}

func fakeValueForType(typ string) string {
	switch typ {
	case "string":
		return `fake.Word()`
	case "bool":
		return `fake.Bool()`
	case "int", "int8", "int16", "int32":
		return `int(fake.Int32())`
	case "int64":
		return `fake.Int64()`
	case "uint", "uint8", "uint16", "uint32":
		return `uint(fake.Uint32())`
	case "uint64":
		return `fake.Uint64()`
	case "float32":
		return `fake.Float32()`
	case "float64":
		return `fake.Float64()`
	case "time.Time":
		return `fake.Date()`
	default:
		return `fake.Word()`
	}
}

// tsType returns the TypeScript type for a Go type string.
func tsType(typ string) string {
	base := typ
	isPtr := false
	if len(base) > 0 && base[0] == '*' {
		base = base[1:]
		isPtr = true
	}

	var ts string
	switch base {
	case "bool":
		ts = "boolean"
	case "string":
		ts = "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		ts = "number"
	case "time.Time":
		// Transported as an ISO 8601 string over JSON; callers parse with new Date().
		ts = "string"
	default:
		ts = "string"
	}

	if isPtr {
		return ts + " | null"
	}
	return ts
}

// tsDefaultValue returns the TypeScript default value for a Go type string.
func tsDefaultValue(typ string) string {
	if len(typ) > 0 && typ[0] == '*' {
		return "null"
	}
	switch typ {
	case "bool":
		return "false"
	case "string":
		return "''"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "0"
	case "time.Time":
		return "''"
	default:
		return "''"
	}
}

// htmlInputType returns the HTML input type for a Go type string.
func htmlInputType(typ string) string {
	base := typ
	if len(base) > 0 && base[0] == '*' {
		base = base[1:]
	}
	switch base {
	case "bool":
		return "checkbox"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "time.Time":
		return "datetime-local"
	default:
		return "text"
	}
}

// swiftType returns the Swift type for a Go type string.
func swiftType(typ string) string {
	base := typ
	isPtr := false
	if len(base) > 0 && base[0] == '*' {
		base = base[1:]
		isPtr = true
	}

	var swift string
	switch base {
	case "bool":
		swift = "Bool"
	case "string":
		swift = "String"
	case "int", "int32":
		swift = "Int32"
	case "int8":
		swift = "Int8"
	case "int16":
		swift = "Int16"
	case "int64":
		swift = "Int"
	case "uint", "uint32":
		swift = "UInt32"
	case "uint8":
		swift = "UInt8"
	case "uint16":
		swift = "UInt16"
	case "uint64":
		swift = "UInt"
	case "float32":
		swift = "Float"
	case "float64":
		swift = "Double"
	case "time.Time":
		swift = "Date"
	default:
		swift = "String"
	}

	if isPtr {
		return swift + "?"
	}
	return swift
}

// swiftZeroValue returns the Swift default value for a Go type string.
func swiftZeroValue(typ string) string {
	if len(typ) > 0 && typ[0] == '*' {
		return "nil"
	}
	switch typ {
	case "bool":
		return "false"
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "0"
	case "time.Time":
		return "Date()"
	default:
		return `""`
	}
}

// sharedColumn reports whether the given snake_case column name is declared
// as a constant in cmd/tools/codegen/queries/helpers.go. Per-entity files
// must NOT redeclare these constants (Go would error: "X redeclared in this block").
var sharedColumns = map[string]struct{}{
	"id":                 {},
	"name":               {},
	"plural_name":        {},
	"notes":              {},
	"description":        {},
	"icon_path":          {},
	"slug":               {},
	"created_at":         {},
	"last_updated_at":    {},
	"archived_at":        {},
	"last_indexed_at":    {},
	"belongs_to_account": {},
	"belongs_to_user":    {},
	"created_by_user":    {},
	"content":            {},
	"title":              {},
	"status":             {},
}

func sharedColumn(snakeName string) bool {
	_, ok := sharedColumns[snakeName]
	return ok
}

// jsonTag returns the JSON struct tag value for a field.
func jsonTag(f config.Field) string {
	tag := naming.New(f.Name).Camel()
	if f.Omitempty {
		tag += ",omitempty"
	}
	return tag
}

// jsonTagForStruct returns the JSON struct tag value for a field when used
// on the domain struct or an input variant. It mirrors jsonTag but also
// appends ",omitempty" for pointer-typed fields, matching the upstream
// convention (pointer fields render as T | null in TS and are elided from
// wire payloads when nil).
func jsonTagForStruct(f config.Field) string {
	tag := naming.New(f.Name).Camel()
	if f.Omitempty || f.IsPointer() {
		tag += ",omitempty"
	}
	return tag
}

// parentChain returns the ancestor entities of e reachable via repeated
// `belongs_to` traversal, ordered outermost-first (e.g., for MealPlanOption
// belongs_to MealPlanEvent belongs_to MealPlan it returns [MealPlan,
// MealPlanEvent]). Self is excluded. An entity with no belongs_to returns
// a zero-length slice. Unknown parent names (dangling belongs_to references)
// terminate the chain silently; validation is the config layer's job.
func parentChain(d *config.Domain, e *config.Entity) []*config.Entity {
	var chain []*config.Entity
	cur := e
	// Bound the walk to len(d.Entities) to avoid looping on cycles (the
	// config validator should reject these, but defend anyway).
	for i := 0; cur != nil && cur.BelongsTo != "" && i < len(d.Entities); i++ {
		parent := findEntity(d, cur.BelongsTo)
		if parent == nil {
			break
		}
		chain = append([]*config.Entity{parent}, chain...)
		cur = parent
	}
	return chain
}

// findEntity returns the first entity in d whose Name matches name, or nil.
func findEntity(d *config.Domain, name string) *config.Entity {
	for i := range d.Entities {
		if d.Entities[i].Name == name {
			return &d.Entities[i]
		}
	}
	return nil
}

// existsArgNames returns the ordered list of identifier-argument names used
// for Exists/Get method signatures: every parent's ID (outermost-first)
// followed by the entity's own ID. All entries are camelCase + "ID".
func existsArgNames(d *config.Domain, e *config.Entity) []string {
	names := make([]string, 0, 4)
	for _, p := range parentChain(d, e) {
		names = append(names, naming.New(p.Name).Camel()+"ID")
	}
	names = append(names, naming.New(e.Name).Camel()+"ID")
	return names
}

// listArgNames returns the ordered parent-scope identifier names used for
// the List method signature (Get<Plural>). Nested entities yield the parent
// chain; top-level entities yield a single scope arg (accountID when
// belongs_to_account, otherwise userID when created_by_user, otherwise
// empty for globally-scoped collections).
func listArgNames(d *config.Domain, e *config.Entity) []string {
	parents := parentChain(d, e)
	if len(parents) > 0 {
		names := make([]string, 0, len(parents))
		for _, p := range parents {
			names = append(names, naming.New(p.Name).Camel()+"ID")
		}
		return names
	}
	switch {
	case e.BelongsToAccount:
		return []string{"accountID"}
	case e.CreatedByUser:
		return []string{"userID"}
	default:
		return nil
	}
}

// nestedChildren returns the entities in d that declare belongs_to ==
// parent.Name and are marked nested. Order follows the parent's own
// Entities slice position (i.e. config order), which is what the emitted
// converter's field-iteration order will match. Self is excluded.
func nestedChildren(d *config.Domain, parent *config.Entity) []*config.Entity {
	var out []*config.Entity
	for i := range d.Entities {
		c := &d.Entities[i]
		if c.Name == parent.Name {
			continue
		}
		if c.BelongsTo != parent.Name {
			continue
		}
		if !c.Nested {
			continue
		}
		out = append(out, c)
	}
	return out
}

// nestedCollectionField returns the name of the slice field on the parent
// that holds child instances. Derivation: strip the parent's Name from
// the start of the child's Name (if present), then pluralize. Falls back
// to plural of the child name when the child name doesn't carry the
// parent prefix. Examples:
//
//	parent=Meal,     child=MealComponent    → "Components"
//	parent=RecipeStep, child=RecipeStepIngredient → "Ingredients"
//	parent=Foo,      child=Bar              → "Bars"
func nestedCollectionField(parent, child *config.Entity) string {
	stem := strings.TrimPrefix(child.Name, parent.Name)
	if stem == "" {
		stem = child.Name
	}
	return naming.New(stem).Plural()
}

// archiveArgNames returns the ordered identifier names used for the
// Archive method signature. Nested entities: parent chain + selfID
// (matching Exists/Get). Top-level entities: selfID followed by the
// scope arg (accountID or userID), or just selfID when globally scoped.
// The self-first ordering on top-level mirrors the upstream convention
// (e.g., ArchiveMeal(mealID, userID)).
func archiveArgNames(d *config.Domain, e *config.Entity) []string {
	parents := parentChain(d, e)
	self := naming.New(e.Name).Camel() + "ID"
	if len(parents) > 0 {
		names := make([]string, 0, len(parents)+1)
		for _, p := range parents {
			names = append(names, naming.New(p.Name).Camel()+"ID")
		}
		return append(names, self)
	}
	switch {
	case e.BelongsToAccount:
		return []string{self, "accountID"}
	case e.CreatedByUser:
		return []string{self, "userID"}
	default:
		return []string{self}
	}
}
