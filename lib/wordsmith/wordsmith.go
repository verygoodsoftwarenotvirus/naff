package wordsmith

import (
	"github.com/codemodus/kace"
	pluralize "github.com/gertd/go-pluralize"
)

func FromSingularPascalCase(input string) *SuperPalabra {
	if len(input) > 0 {
		return &SuperPalabra{
			meta:       input,
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
	return s.meta
}

func (s *SuperPalabra) Plural() string {
	return s.pluralizer.Plural(s.meta)
}

func (s *SuperPalabra) RouteName() string {
	return kace.Snake(s.meta)
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
