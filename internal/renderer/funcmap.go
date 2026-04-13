package renderer

import (
	"fmt"
	"strings"
	"text/template"

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
		"fakeValue":           fakeValue,
		"pointerConversionFunc": pointerConversionFunc,

		// Field filtering
		"creatableFields": creatableFields,
		"editableFields":  editableFields,
		"requiredFields":  requiredFields,
		"columnFields":    columnFields,

		// JSON tag helpers
		"jsonTag": jsonTag,

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
	case "float32", "float64":
		return "Float64PointerFromNullFloat64"
	default:
		return "StringPointerFromNullString"
	}
}

// fakeValue returns a gofakeit expression for generating fake data for a field.
func fakeValue(f config.Field) string {
	if f.IsPointer() {
		// For pointer fields in fakes, we don't generate pointer values directly.
		// The fake builder will take the address.
		return fakeValueForType(f.BaseType())
	}
	return fakeValueForType(f.Type)
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
	default:
		return `""`
	}
}

// jsonTag returns the JSON struct tag value for a field.
func jsonTag(f config.Field) string {
	tag := naming.New(f.Name).Camel()
	if f.Omitempty {
		tag += ",omitempty"
	}
	return tag
}
