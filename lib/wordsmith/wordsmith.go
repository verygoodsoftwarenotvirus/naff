package wordsmith

import (
	"fmt"
	"strings"

	"github.com/codemodus/kace"
	pluralize "github.com/gertd/go-pluralize"
)

func FromSingularPascalCase(input string) *SuperPalabra {
	if len(input) > 0 {
		return &SuperPalabra{
			meta:       kace.Pascal(input),
			pluralizer: pluralize.NewClient(),
		}
	}
	return nil
}

type SuperPalabra struct {
	meta       string
	pluralizer *pluralize.Client
}

func (s *SuperPalabra) Singular() string {
	return kace.Pascal(s.meta)
}

func (s *SuperPalabra) Plural() string {
	return s.pluralizer.Plural(s.meta)
}

func (s *SuperPalabra) RouteName() string {
	return kace.Snake(s.meta)
}

func (s *SuperPalabra) KebabName() string {
	return kace.Kebab(s.meta)
}

func (s *SuperPalabra) PluralRouteName() string {
	return kace.Snake(s.pluralizer.Plural(s.meta))
}

func (s *SuperPalabra) UnexportedVarName() string {
	return kace.Camel(s.meta)
}

func (s *SuperPalabra) PluralUnexportedVarName() string {
	return kace.Camel(s.pluralizer.Plural(s.meta))
}

func (s *SuperPalabra) PackageName() string {
	return strings.ToLower(s.Plural())
}

func (s *SuperPalabra) SingularPackageName() string {
	return strings.ToLower(s.Singular())
}

func (s *SuperPalabra) SingularCommonName() string {
	return strings.Join(strings.Split(s.RouteName(), "_"), " ")
}

func (s *SuperPalabra) ProperSingularCommonNameWithPrefix() string {
	return fmt.Sprintf("%s %s", AOrAn(s.Singular()), strings.Title(strings.Join(strings.Split(s.RouteName(), "_"), " ")))
}

func (s *SuperPalabra) PluralCommonName() string {
	return strings.Join(strings.Split(s.PluralRouteName(), "_"), " ")
}

func (s *SuperPalabra) SingularCommonNameWithPrefix() string {
	return fmt.Sprintf("%s %s", AOrAn(s.Singular()), s.SingularCommonName())
}

func (s *SuperPalabra) PluralCommonNameWithPrefix() string {
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
