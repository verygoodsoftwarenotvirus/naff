package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_randTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := randTestDotGo()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

var _ secretGenerator = (*mockSecretGenerator)(nil)

type mockSecretGenerator struct {
	mock.Mock
}

func (m *mockSecretGenerator) GenerateTwoFactorSecret() (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func (m *mockSecretGenerator) GenerateSalt() ([]byte, error) {
	args := m.Called()

	return args.Get(0).([]byte), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandTestMockSecretGenerator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildRandTestMockSecretGenerator()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

var _ secretGenerator = (*mockSecretGenerator)(nil)

type mockSecretGenerator struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandTestMockSecretGeneratorGenerateTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildRandTestMockSecretGeneratorGenerateTwoFactorSecret()

		expected := `
package example

import ()

func (m *mockSecretGenerator) GenerateTwoFactorSecret() (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandTestMockSecretGeneratorGenerateSalt(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildRandTestMockSecretGeneratorGenerateSalt()

		expected := `
package example

import ()

func (m *mockSecretGenerator) GenerateSalt() ([]byte, error) {
	args := m.Called()

	return args.Get(0).([]byte), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
