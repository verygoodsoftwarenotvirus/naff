package wordsmith

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAOrAn(T *testing.T) {
	T.Parallel()

	T.Run("without vowel", func(t *testing.T) {
		t.Parallel()

		expected := "a"
		actual := AOrAn("computer")

		assert.Equal(t, expected, actual)
	})

	T.Run("with vowel", func(t *testing.T) {
		t.Parallel()

		expected := "an"
		actual := AOrAn("arachnophobe")

		assert.Equal(t, expected, actual)
	})
}

func TestFromSingularPascalCase(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, FromSingularPascalCase("Things"))
	})

	T.Run("with empty string", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, FromSingularPascalCase(""))
	})
}

func TestSuperWord_Abbreviation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "B"
		actual := w.Abbreviation()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_KebabName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blah"
		actual := w.KebabName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_LowercaseAbbreviation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "b"
		actual := w.LowercaseAbbreviation()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_PackageName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blahs"
		actual := w.PackageName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_Plural(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "Blahs"
		actual := w.Plural()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_PluralCommonName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blahs"
		actual := w.PluralCommonName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_PluralRouteName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blahs"
		actual := w.PluralRouteName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_PluralUnexportedVarName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blahs"
		actual := w.PluralUnexportedVarName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_ProperSingularCommonNameWithPrefix(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "a Blah"
		actual := w.ProperSingularCommonNameWithPrefix()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_RouteName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blah"
		actual := w.RouteName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_Singular(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "Blah"
		actual := w.Singular()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_SingularCommonName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blah"
		actual := w.SingularCommonName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_SingularCommonNameWithPrefix(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "a blah"
		actual := w.SingularCommonNameWithPrefix()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_SingularPackageName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blah"
		actual := w.SingularPackageName()

		assert.Equal(t, expected, actual)
	})
}

func TestSuperWord_UnexportedVarName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Blah")

		expected := "blah"
		actual := w.UnexportedVarName()

		assert.Equal(t, expected, actual)
	})

	T.Run("with reserved word", func(t *testing.T) {
		t.Parallel()

		w := FromSingularPascalCase("Var")

		expected := "_var"
		actual := w.UnexportedVarName()

		assert.Equal(t, expected, actual)
	})
}
