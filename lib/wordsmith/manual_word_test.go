package wordsmith

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManualWord_Abbreviation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			AbbreviationStr: expected,
		}

		actual := mw.Abbreviation()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_KebabName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			KebabNameStr: expected,
		}

		actual := mw.KebabName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_LowercaseAbbreviation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			LowercaseAbbreviationStr: expected,
		}

		actual := mw.LowercaseAbbreviation()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_PackageName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			PackageNameStr: expected,
		}

		actual := mw.PackageName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_Plural(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			PluralStr: expected,
		}

		actual := mw.Plural()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_PluralCommonName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			PluralCommonNameStr: expected,
		}

		actual := mw.PluralCommonName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_PluralRouteName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			PluralRouteNameStr: expected,
		}

		actual := mw.PluralRouteName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_PluralUnexportedVarName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			PluralUnexportedVarNameStr: expected,
		}

		actual := mw.PluralUnexportedVarName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_ProperSingularCommonNameWithPrefix(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			ProperSingularCommonNameWithPrefixStr: expected,
		}

		actual := mw.ProperSingularCommonNameWithPrefix()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_RouteName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			RouteNameStr: expected,
		}

		actual := mw.RouteName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_Singular(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			SingularStr: expected,
		}

		actual := mw.Singular()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_SingularCommonName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			SingularCommonNameStr: expected,
		}

		actual := mw.SingularCommonName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_SingularCommonNameWithPrefix(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			SingularCommonNameWithPrefixStr: expected,
		}

		actual := mw.SingularCommonNameWithPrefix()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_SingularPackageName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			SingularPackageNameStr: expected,
		}

		actual := mw.SingularPackageName()

		assert.Equal(t, expected, actual)
	})
}

func TestManualWord_UnexportedVarName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "blah"
		mw := &ManualWord{
			UnexportedVarNameStr: expected,
		}

		actual := mw.UnexportedVarName()

		assert.Equal(t, expected, actual)
	})
}
