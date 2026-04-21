package naming

import (
	"strings"
	"unicode"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

var pluralizer = pluralize.NewClient()

// Name provides all case/plurality transformations for an entity or field name.
type Name struct {
	singular string
	plural   string
}

// New creates a Name from a PascalCase singular form (e.g. "IssueReport").
// The plural is auto-computed but can be overridden with NewWithPlural.
func New(singular string) Name {
	return Name{
		singular: singular,
		plural:   pluralizer.Plural(singular),
	}
}

// NewWithPlural creates a Name with an explicit plural override.
func NewWithPlural(singular, plural string) Name {
	return Name{
		singular: singular,
		plural:   plural,
	}
}

// Singular returns the PascalCase singular form: "IssueReport"
func (n Name) Singular() string { return n.singular }

// Plural returns the PascalCase plural form: "IssueReports"
func (n Name) Plural() string { return n.plural }

// Snake returns the snake_case singular form: "issue_report"
func (n Name) Snake() string { return strcase.ToSnake(n.singular) }

// PluralSnake returns the snake_case plural form: "issue_reports"
func (n Name) PluralSnake() string { return strcase.ToSnake(n.plural) }

// Camel returns the camelCase singular form: "issueReport"
func (n Name) Camel() string { return strcase.ToLowerCamel(n.singular) }

// PluralCamel returns the camelCase plural form: "issueReports"
func (n Name) PluralCamel() string { return strcase.ToLowerCamel(n.plural) }

// Kebab returns the kebab-case singular form: "issue-report"
func (n Name) Kebab() string { return strcase.ToKebab(n.singular) }

// PluralKebab returns the kebab-case plural form: "issue-reports"
func (n Name) PluralKebab() string { return strcase.ToKebab(n.plural) }

// PackageName returns the lowercase plural form suitable for a Go package name: "issuereports"
func (n Name) PackageName() string {
	return strings.ToLower(n.plural)
}

// Abbreviation returns the uppercase initials: "IR" for "IssueReport"
func (n Name) Abbreviation() string {
	var b strings.Builder
	for _, r := range n.singular {
		if unicode.IsUpper(r) {
			b.WriteRune(r)
		}
	}
	result := b.String()
	if result == "" && len(n.singular) > 0 {
		return strings.ToUpper(n.singular[:1])
	}
	return result
}

// LowerAbbreviation returns the lowercase initials: "ir" for "IssueReport"
func (n Name) LowerAbbreviation() string {
	return strings.ToLower(n.Abbreviation())
}

// HumanReadable returns the lowercased, space-separated form: "issue report"
func (n Name) HumanReadable() string {
	return strings.ToLower(strcase.ToDelimited(n.singular, ' '))
}

// PluralHumanReadable returns the lowercased, space-separated plural form: "issue reports"
func (n Name) PluralHumanReadable() string {
	return strings.ToLower(strcase.ToDelimited(n.plural, ' '))
}

// uConsonantSoundPrefixes enumerates lowercased prefixes of words that begin with
// the letter 'u' but are pronounced with a /j/ ("yoo") consonant sound, and thus
// take "a" rather than "an" (e.g. "user", "unit", "universe", "utility").
var uConsonantSoundPrefixes = []string{
	"ubi",
	"uni",
	"usa",
	"use",
	"uti",
}

// hSilentPrefixes enumerates lowercased prefixes of words that begin with a
// silent 'h' and thus take "an" rather than "a" (e.g. "hour", "honor", "heir").
var hSilentPrefixes = []string{
	"heir",
	"hones",
	"honor",
	"hour",
}

// Article returns "a" or "an" based on the pronounced first sound of the
// human-readable form. It handles the common English exceptions where letter
// and sound diverge: /j/-sound u-words take "a" (a user), and silent-h words
// take "an" (an hour).
func (n Name) Article() string {
	hr := n.HumanReadable()
	if len(hr) == 0 {
		return "a"
	}
	for _, p := range hSilentPrefixes {
		if strings.HasPrefix(hr, p) {
			return "an"
		}
	}
	for _, p := range uConsonantSoundPrefixes {
		if strings.HasPrefix(hr, p) {
			return "a"
		}
	}
	switch unicode.ToLower(rune(hr[0])) {
	case 'a', 'e', 'i', 'o', 'u':
		return "an"
	default:
		return "a"
	}
}

// WithArticle returns the article + human-readable form: "an issue report"
func (n Name) WithArticle() string {
	return n.Article() + " " + n.HumanReadable()
}
