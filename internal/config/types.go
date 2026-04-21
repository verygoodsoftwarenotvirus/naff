package config

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"

	"github.com/verygoodsoftwarenotvirus/naff/internal/naming"
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
	BelongsToUser    bool    `yaml:"belongs_to_user,omitempty"`
	CreatedByUser    bool    `yaml:"created_by_user,omitempty"`
	Searchable       bool    `yaml:"searchable,omitempty"`
	ConsumerEditable bool    `yaml:"consumer_editable,omitempty"`
	// Nested, when true, signals that this entity is denormalized as a
	// collection on its parent's domain struct and input variants (e.g.
	// MealComponent is nested on Meal as `Components []*MealComponent`).
	// The parent's parameterized converter then iterates that collection
	// and calls this entity's matching converter. Only meaningful when
	// BelongsTo is also set. Defaults false — children are fetched
	// separately by default.
	Nested  bool    `yaml:"nested,omitempty"`
	LinksTo []Link  `yaml:"links_to,omitempty"`
	Fields  []Field `yaml:"fields,omitempty"`
}

// SyntheticLinkFields returns one Field per LinksTo entry, materialized via
// Link.AsField(). Lets templates iterate links through the same helpers used
// for scalar fields (creatableFields, editableFields, jsonTag, etc.).
func (e Entity) SyntheticLinkFields() []Field {
	out := make([]Field, 0, len(e.LinksTo))
	for _, l := range e.LinksTo {
		out = append(out, l.AsField())
	}
	return out
}

// AllFields returns scalar Fields followed by synthetic link Fields. Use
// from per-entity templates that emit struct definitions / input variants
// so links and scalars get identical filtering and formatting.
func (e Entity) AllFields() []Field {
	out := make([]Field, 0, len(e.Fields)+len(e.LinksTo))
	out = append(out, e.Fields...)
	out = append(out, e.SyntheticLinkFields()...)
	return out
}

// Link represents a non-owning foreign-key reference from this entity to another.
// Unlike BelongsTo, Link does not imply URL nesting, lifecycle ownership, or
// cycle restrictions — links can cross domains and refer back to the enclosing
// entity (self-references).
type Link struct {
	// Target is the name of the entity being referenced. Required. Must match
	// an entity declared elsewhere in the config (or this entity itself when
	// Self is true).
	Target string `yaml:"target"`
	// As overrides the column/field prefix. Without it, naming derives from
	// Target (e.g. Target=ValidIngredient → field=ValidIngredientID, column=
	// valid_ingredient_id). With As="ForIngredient" → ForIngredientID /
	// for_ingredient_id. Use when one entity links to the same target twice.
	As string `yaml:"as,omitempty"`
	// DomainAs overrides the *domain struct* accessor name for this link,
	// i.e. the field name under which the linked entity is denormalized onto
	// the owning domain struct. Defaults to As when set, otherwise Target.
	// Use this when the domain accessor should drop a redundant prefix (e.g.
	// Target=ValidPreparation on a ValidIngredientPreparation owner → set
	// DomainAs=Preparation so the field is `Preparation ValidPreparation`
	// rather than `ValidPreparation ValidPreparation`). The FK field name
	// (FieldName) is still derived from As/Target and is unaffected.
	DomainAs string `yaml:"domain_as,omitempty"`
	// Optional marks the FK nullable. Defaults false (required link).
	Optional bool `yaml:"optional,omitempty"`
	// Editable controls whether the link appears in UpdateRequestInput.
	// Defaults false — rows are deleted and re-added rather than re-targeted.
	Editable *bool `yaml:"editable,omitempty"`
	// Self must be true when Target equals the enclosing entity's Name. It
	// documents intent and keeps self-references visible at a glance.
	Self bool `yaml:"self,omitempty"`
	// OnDelete controls the SQL ON DELETE action. Defaults to RESTRICT.
	OnDelete LinkOnDelete `yaml:"on_delete,omitempty"`
}

// LinkOnDelete enumerates the valid ON DELETE actions for a foreign-key link.
type LinkOnDelete string

const (
	LinkOnDeleteRestrict LinkOnDelete = "restrict"
	LinkOnDeleteCascade  LinkOnDelete = "cascade"
	LinkOnDeleteSetNull  LinkOnDelete = "set_null"
)

// UnmarshalYAML accepts either a bare string ("UploadedMedia") as shorthand
// for {target: UploadedMedia, optional: false} or the full mapping form.
func (l *Link) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		l.Target = node.Value
		return nil
	}
	type linkAlias Link
	var a linkAlias
	if err := node.Decode(&a); err != nil {
		return err
	}
	*l = Link(a)
	return nil
}

// IsEditable reports whether the link should appear in UpdateRequestInput.
// Defaults to false — rows are deleted and re-added rather than re-targeted.
func (l Link) IsEditable() bool {
	if l.Editable == nil {
		return false
	}
	return *l.Editable
}

// OnDeleteAction returns the ON DELETE action, defaulting to RESTRICT when unset.
func (l Link) OnDeleteAction() LinkOnDelete {
	if l.OnDelete == "" {
		return LinkOnDeleteRestrict
	}
	return l.OnDelete
}

// prefix returns the naming source for derived field/column names: As when
// set, otherwise Target.
func (l Link) prefix() string {
	if l.As != "" {
		return l.As
	}
	return l.Target
}

// FieldName returns the synthesized Go struct field name for this link
// (e.g. Target=Recipe → "RecipeID"; As=ForIngredient → "ForIngredientID").
// Always suffixed with "ID".
func (l Link) FieldName() string {
	return l.prefix() + "ID"
}

// DomainAccessor returns the name under which the linked entity is
// denormalized onto the owning domain struct. Precedence: DomainAs → As →
// Target. The accessor is NOT suffixed with "ID" — it names a nested struct
// (e.g. "Preparation" for `Preparation ValidPreparation`), not a scalar
// column.
func (l Link) DomainAccessor() string {
	if l.DomainAs != "" {
		return l.DomainAs
	}
	if l.As != "" {
		return l.As
	}
	return l.Target
}

// ColumnName returns the snake-case SQL column name for this link
// (e.g. Target=ValidIngredient → "valid_ingredient_id").
func (l Link) ColumnName() string {
	return naming.New(l.prefix()).Snake() + "_id"
}

// GoType returns the Go type for the link's foreign-key field: "string" when
// required, "*string" when Optional.
func (l Link) GoType() string {
	if l.Optional {
		return "*string"
	}
	return "string"
}

// JSONTag returns the JSON struct tag value for the link's field
// (camelCase, with ",omitempty" appended when Optional).
func (l Link) JSONTag() string {
	tag := naming.New(l.prefix()).Camel() + "ID"
	if l.Optional {
		tag += ",omitempty"
	}
	return tag
}

// SQLOnDelete returns the uppercased SQL keyword for the ON DELETE action
// (e.g. "RESTRICT", "CASCADE", "SET NULL").
func (l Link) SQLOnDelete() string {
	switch l.OnDeleteAction() {
	case LinkOnDeleteCascade:
		return "CASCADE"
	case LinkOnDeleteSetNull:
		return "SET NULL"
	default:
		return "RESTRICT"
	}
}

// AsField materializes the link as a Field so existing field-iterating
// helpers (creatableFields, editableFields, requiredFields, jsonTag, etc.)
// can treat link FKs and scalar fields uniformly. The returned Field has:
//   - Name = FieldName() (e.g. "RecipeID")
//   - Type = GoType() (e.g. "string" or "*string")
//   - Required = !Optional (a required link → required field)
//   - Creatable = true (links are always set at creation time)
//   - Editable = IsEditable() (defaults false; rows are deleted/re-added)
//   - Omitempty = Optional (mirrors JSONTag's ",omitempty" rule)
func (l Link) AsField() Field {
	required := !l.Optional
	creatable := true
	editable := l.IsEditable()
	return Field{
		Name:      l.FieldName(),
		Type:      l.GoType(),
		Required:  &required,
		Creatable: &creatable,
		Editable:  &editable,
		Omitempty: l.Optional,
	}
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
