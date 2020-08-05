package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_randDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := randDotGo()

		expected := `
package example

import (
	"crypto/rand"
	"encoding/base32"
)

const (
	saltSize         = 16
	randomSecretSize = 64
)

// this function tests that we have appropriate access to crypto/rand
func init() {
	b := make([]byte, randomSecretSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}

var _ secretGenerator = (*standardSecretGenerator)(nil)

type standardSecretGenerator struct{}

func (g *standardSecretGenerator) GenerateTwoFactorSecret() (string, error) {
	b := make([]byte, randomSecretSize)

	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(b), nil
}

func (g *standardSecretGenerator) GenerateSalt() ([]byte, error) {
	b := make([]byte, saltSize)

	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRandConstantDefs()

		expected := `
package example

import ()

const (
	saltSize         = 16
	randomSecretSize = 64
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRandInit()

		expected := `
package example

import (
	"crypto/rand"
)

// this function tests that we have appropriate access to crypto/rand
func init() {
	b := make([]byte, randomSecretSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandStandardSecretGeneratorGenerateTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRandStandardSecretGeneratorGenerateTwoFactorSecret()

		expected := `
package example

import (
	"crypto/rand"
	"encoding/base32"
)

var _ secretGenerator = (*standardSecretGenerator)(nil)

type standardSecretGenerator struct{}

func (g *standardSecretGenerator) GenerateTwoFactorSecret() (string, error) {
	b := make([]byte, randomSecretSize)

	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(b), nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandStandardSecretGeneratorGenerateSalt(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRandStandardSecretGeneratorGenerateSalt()

		expected := `
package example

import (
	"crypto/rand"
)

func (g *standardSecretGenerator) GenerateSalt() ([]byte, error) {
	b := make([]byte, saltSize)

	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
