package config

import "strings"

// Project is the top-level configuration for a NAFF-generated project.
type Project struct {
	ProjectMeta ProjectMeta        `yaml:"project"`
	Features    Features           `yaml:"features"`
	Targets     Targets            `yaml:"targets"`
	Domains     []Domain           `yaml:"domains"`
}

// ProjectMeta contains project-level metadata.
type ProjectMeta struct {
	Name           string `yaml:"name"`
	Module         string `yaml:"module"`
	PlatformModule string `yaml:"platform_module"`
	IOSBundleID    string `yaml:"ios_bundle_id,omitempty"`
	IOSModuleName  string `yaml:"ios_module_name,omitempty"`
}

// Targets controls which generation targets are enabled.
//
// Backend and AsyncMessageHandler are *bool so "unset" (default on) is
// distinguishable from "explicitly false" (opt-out). IOS is a plain bool
// (default off) — iOS scaffolding is opt-in.
type Targets struct {
	Backend             *bool `yaml:"backend"`
	IOS                 bool  `yaml:"ios"`
	AsyncMessageHandler *bool `yaml:"async_message_handler"`
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
	Webhooks      *bool `yaml:"webhooks"`
	IssueReports  *bool `yaml:"issuereports"`
	Comments      *bool `yaml:"comments"`
	Notifications *bool `yaml:"notifications"`
	Payments      *bool `yaml:"payments"`
	Waitlists     *bool `yaml:"waitlists"`
	UploadedMedia *bool `yaml:"uploadedmedia"`
	DataPrivacy   *bool `yaml:"dataprivacy"`
	Settings      *bool `yaml:"settings"`
	ConsumerApp   *bool `yaml:"consumer_app"`
	AdminApp      *bool `yaml:"admin_app"`
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
	Entities []Entity `yaml:"entities"`
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
	BelongsTo        string  `yaml:"belongs_to"`
	BelongsToAccount bool    `yaml:"belongs_to_account"`
	CreatedByUser    bool    `yaml:"created_by_user"`
	Searchable       bool    `yaml:"searchable"`
	ConsumerEditable bool    `yaml:"consumer_editable"`
	Fields           []Field `yaml:"fields"`
}

// Field represents a single field on an entity.
type Field struct {
	Name      string  `yaml:"name"`
	Type      string  `yaml:"type"`
	Required  *bool   `yaml:"required"`
	Creatable *bool   `yaml:"creatable"`
	Editable  *bool   `yaml:"editable"`
	Omitempty bool    `yaml:"omitempty"`
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
