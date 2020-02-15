package wordsmith

import (
	"fmt"
	"strings"

	"github.com/codemodus/kace"
	pluralize "github.com/gertd/go-pluralize"
)

func FromSingularPascalCase(input string) SuperPalabra {
	if len(input) > 0 {
		return &SuperWord{
			meta:       kace.Pascal(input),
			pluralizer: pluralize.NewClient(),
		}
	}
	return nil
}

type SuperPalabra interface {
	Singular() string
	Plural() string
	RouteName() string
	KebabName() string
	PluralRouteName() string
	UnexportedVarName() string
	PluralUnexportedVarName() string
	PackageName() string
	SingularPackageName() string
	SingularCommonName() string
	ProperSingularCommonNameWithPrefix() string
	PluralCommonName() string
	SingularCommonNameWithPrefix() string
	PluralCommonNameWithPrefix() string
}

type SuperWord struct {
	meta       string
	pluralizer *pluralize.Client
}

func (s *SuperWord) Singular() string {
	return kace.Pascal(s.meta)
}

func (s *SuperWord) Plural() string {
	return s.pluralizer.Plural(s.meta)
}

func (s *SuperWord) RouteName() string {
	return kace.Snake(s.meta)
}

func (s *SuperWord) KebabName() string {
	return kace.Kebab(s.meta)
}

func (s *SuperWord) PluralRouteName() string {
	return kace.Snake(s.pluralizer.Plural(s.meta))
}

func (s *SuperWord) UnexportedVarName() string {
	x := strings.ToLower(s.meta)
	switch x {
	case "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func", "go", "goto", "if", "iota", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "type", "var":
		return kace.Camel(fmt.Sprintf("_%s", s.meta))
	default:
		return kace.Camel(s.meta)
	}
}

func (s *SuperWord) PluralUnexportedVarName() string {
	return kace.Camel(s.pluralizer.Plural(s.meta))
}

func (s *SuperWord) PackageName() string {
	return strings.ToLower(s.Plural())
}

func (s *SuperWord) SingularPackageName() string {
	return strings.ToLower(s.Singular())
}

func (s *SuperWord) SingularCommonName() string {
	return strings.Join(strings.Split(s.RouteName(), "_"), " ")
}

func (s *SuperWord) ProperSingularCommonNameWithPrefix() string {
	return fmt.Sprintf("%s %s", AOrAn(s.Singular()), strings.Title(strings.Join(strings.Split(s.RouteName(), "_"), " ")))
}

func (s *SuperWord) PluralCommonName() string {
	return strings.Join(strings.Split(s.PluralRouteName(), "_"), " ")
}

func (s *SuperWord) SingularCommonNameWithPrefix() string {
	return fmt.Sprintf("%s %s", AOrAn(s.Singular()), s.SingularCommonName())
}

func (s *SuperWord) PluralCommonNameWithPrefix() string {
	return fmt.Sprintf("%s %s", AOrAn(s.Singular()), s.PluralCommonName())
}

// AOrAn return "a" or "an" depending on the input
func AOrAn(input string) string {
	if len(input) > 0 {
		switch input[0] {
		case 'a', 'A', 'e', 'E', 'i', 'I', 'o', 'O', 'u', 'U':
			return "an"
		}
	}
	return "a"
}
