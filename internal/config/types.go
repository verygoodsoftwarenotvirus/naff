package config

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

// PlatformModule is the Go module path for the platform package that every
// naff-generated project depends on. It is fixed — naff does not support
// pointing at an alternative platform implementation.
const PlatformModule = "github.com/primandproper/platform"

// Project is the top-level configuration for a NAFF-generated project.
type Project struct {
	ProjectMeta ProjectMeta `yaml:"project"`
	Features    Features    `yaml:"features,omitempty"`
	Targets     Targets     `yaml:"targets,omitempty"`
	Domains     []Domain    `yaml:"domains,omitempty"`
}

// ProjectMeta contains project-level metadata.
type ProjectMeta struct {
	Name          string `yaml:"name"`
	Module        string `yaml:"module"`
	IOSBundleID   string `yaml:"ios_bundle_id,omitempty"`
	IOSModuleName string `yaml:"ios_module_name,omitempty"`
}

// Name-form methods. These exist so templates can reference the project Name
// in whichever casing a given string literal needs, without the template
// author having to know anything about casing conventions. Input is always
// the natural display form from `project.name` in .naff.yaml, e.g.
// "Dinner Done Better" or "Kitchen Sink".

// Title returns the project name in its natural display form
// (e.g. "Dinner Done Better").
func (p ProjectMeta) Title() string { return p.Name }

// Pascal returns the project name in PascalCase (e.g. "DinnerDoneBetter").
func (p ProjectMeta) Pascal() string { return strcase.ToCamel(p.Name) }

// Kebab returns the project name in kebab-case (e.g. "dinner-done-better").
func (p ProjectMeta) Kebab() string { return strcase.ToKebab(p.Name) }

// Snake returns the project name in snake_case (e.g. "dinner_done_better").
func (p ProjectMeta) Snake() string { return strcase.ToSnake(p.Name) }

// TitleKebab returns the project name with its original word casing,
// whitespace-joined with hyphens (e.g. "Dinner-Done-Better"). Used for
// HTTP header names like `X-Dinner-Done-Better-Signature`.
func (p ProjectMeta) TitleKebab() string {
	return strings.Join(strings.Fields(p.Name), "-")
}

// Abbrev returns the uppercase initials of the whitespace-separated words
// in the project name (e.g. "DDB" for "Dinner Done Better", "KS" for
// "Kitchen Sink", "K" for "Kitchen").
func (p ProjectMeta) Abbrev() string {
	var b strings.Builder
	for _, word := range strings.Fields(p.Name) {
		r := []rune(word)[0]
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

// LowerAbbrev returns the lowercase form of Abbrev (e.g. "ddb").
func (p ProjectMeta) LowerAbbrev() string { return strings.ToLower(p.Abbrev()) }

// Targets controls which generation targets are enabled.
//
// Backend and AsyncMessageHandler are *bool so "unset" (default on) is
// distinguishable from "explicitly false" (opt-out). IOS is a plain bool
// (default off) — iOS scaffolding is opt-in.
type Targets struct {
	Backend             *bool `yaml:"backend,omitempty"`
	IOS                 bool  `yaml:"ios,omitempty"`
	AsyncMessageHandler *bool `yaml:"async_message_handler,omitempty"`
}

// BackendEnabled reports whether backend generation is enabled (default true).
func (t Targets) BackendEnabled() bool {
	if t.Backend == nil {
		return true
	}
	return *t.Backend
}

// AsyncMessageHandlerEnabled reports whether cmd/functions/async_message_handler
// should be emitted (default true). Projects that don't ship a hand-rolled
// internal/build/functions/data_change_message_handler package should set
// this to false explicitly — otherwise the generated binary won't compile.
func (t Targets) AsyncMessageHandlerEnabled() bool {
	if t.AsyncMessageHandler == nil {
		return true
	}
	return *t.AsyncMessageHandler
}

// Features controls which built-in infrastructure domains are enabled.
type Features struct {
	Webhooks      *bool `yaml:"webhooks,omitempty"`
	IssueReports  *bool `yaml:"issuereports,omitempty"`
	Comments      *bool `yaml:"comments,omitempty"`
	Notifications *bool `yaml:"notifications,omitempty"`
	Payments      *bool `yaml:"payments,omitempty"`
	Waitlists     *bool `yaml:"waitlists,omitempty"`
	UploadedMedia *bool `yaml:"uploadedmedia,omitempty"`
	DataPrivacy   *bool `yaml:"dataprivacy,omitempty"`
	Settings      *bool `yaml:"settings,omitempty"`
	ConsumerApp   *bool `yaml:"consumer_app,omitempty"`
	AdminApp      *bool `yaml:"admin_app,omitempty"`
}

// FeatureEnabled returns whether a feature is enabled, defaulting to the given value.
func (f Features) FeatureEnabled(val *bool, defaultVal bool) bool {
	if val == nil {
		return defaultVal
	}
	return *val
}

// Domain represents a group of related entities.
type Domain struct {
	Name     string   `yaml:"name"`
	TypeName string   `yaml:"type_name,omitempty"`
	Entities []Entity `yaml:"entities,omitempty"`
}

// RepoTypeName returns the PascalCase identifier used to build type/function
// names derived from this domain (e.g. "ProvideXRepository"). It returns
// TypeName when set; otherwise it falls back to title-casing Name. Falling
// back produces "Audit" or "Webhooks" which is fine for single-word domains,
// but multi-word or irregular domains (audit→"AuditLog", mealplanning→
// "MealPlanning") should set TypeName explicitly.
func (d Domain) RepoTypeName() string {
	if d.TypeName != "" {
		return d.TypeName
	}
	if len(d.Name) == 0 {
		return ""
	}
	// Uppercase the first rune; leave the rest. Multi-word lowercase names
	// (e.g. "mealplanning") will render as "Mealplanning" — set TypeName to
	// override.
	return strings.ToUpper(d.Name[:1]) + d.Name[1:]
}

// Entity represents a single CRUD-able type.
type Entity struct {
	Name             string  `yaml:"name"`
	BelongsTo        string  `yaml:"belongs_to,omitempty"`
	BelongsToAccount bool    `yaml:"belongs_to_account,omitempty"`
	CreatedByUser    bool    `yaml:"created_by_user,omitempty"`
	Searchable       bool    `yaml:"searchable,omitempty"`
	ConsumerEditable bool    `yaml:"consumer_editable,omitempty"`
	Fields           []Field `yaml:"fields,omitempty"`
}

// Field represents a single field on an entity.
type Field struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Required  *bool  `yaml:"required,omitempty"`
	Creatable *bool  `yaml:"creatable,omitempty"`
	Editable  *bool  `yaml:"editable,omitempty"`
	Omitempty bool   `yaml:"omitempty,omitempty"`
}

// IsRequired returns whether the field is required (defaults to true).
func (f Field) IsRequired() bool {
	if f.Required == nil {
		return true
	}
	return *f.Required
}

// IsCreatable returns whether the field is included in creation input (defaults to true).
func (f Field) IsCreatable() bool {
	if f.Creatable == nil {
		return true
	}
	return *f.Creatable
}

// IsEditable returns whether the field is included in update input (defaults to true).
func (f Field) IsEditable() bool {
	if f.Editable == nil {
		return true
	}
	return *f.Editable
}

// IsPointer returns whether the field's Go type is a pointer.
func (f Field) IsPointer() bool {
	return len(f.Type) > 0 && f.Type[0] == '*'
}

// BaseType returns the type without a pointer prefix.
func (f Field) BaseType() string {
	if f.IsPointer() {
		return f.Type[1:]
	}
	return f.Type
}
