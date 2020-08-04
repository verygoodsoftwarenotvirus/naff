package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mockDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
)

var _ auth.Authenticator = (*Authenticator)(nil)

// Authenticator is a mock Authenticator.
type Authenticator struct {
	mock.Mock
}

// ValidateLogin satisfies our authenticator interface.
func (m *Authenticator) ValidateLogin(
	ctx context.Context,
	hashedPassword,
	providedPassword,
	twoFactorSecret,
	twoFactorCode string,
	salt []byte,
) (valid bool, err error) {
	args := m.Called(
		ctx,
		hashedPassword,
		providedPassword,
		twoFactorSecret,
		twoFactorCode,
		salt,
	)
	return args.Bool(0), args.Error(1)
}

// PasswordIsAcceptable satisfies our authenticator interface.
func (m *Authenticator) PasswordIsAcceptable(password string) bool {
	return m.Called(password).Bool(0)
}

// HashPassword satisfies our authenticator interface.
func (m *Authenticator) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

// PasswordMatches satisfies our authenticator interface.
func (m *Authenticator) PasswordMatches(
	ctx context.Context,
	hashedPassword,
	providedPassword string,
	salt []byte,
) bool {
	args := m.Called(
		ctx,
		hashedPassword,
		providedPassword,
		salt,
	)
	return args.Bool(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInterfaceImplementationDeclaration(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildInterfaceImplementationDeclaration(proj)

		expected := `
package example

import (
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
)

var _ auth.Authenticator = (*Authenticator)(nil)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockAuthenticator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockAuthenticator()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// Authenticator is a mock Authenticator.
type Authenticator struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockValidateLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockValidateLogin()

		expected := `
package example

import (
	"context"
)

// ValidateLogin satisfies our authenticator interface.
func (m *Authenticator) ValidateLogin(
	ctx context.Context,
	hashedPassword,
	providedPassword,
	twoFactorSecret,
	twoFactorCode string,
	salt []byte,
) (valid bool, err error) {
	args := m.Called(
		ctx,
		hashedPassword,
		providedPassword,
		twoFactorSecret,
		twoFactorCode,
		salt,
	)
	return args.Bool(0), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockPasswordIsAcceptable(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockPasswordIsAcceptable()

		expected := `
package example

import ()

// PasswordIsAcceptable satisfies our authenticator interface.
func (m *Authenticator) PasswordIsAcceptable(password string) bool {
	return m.Called(password).Bool(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockHashPassword(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockHashPassword()

		expected := `
package example

import (
	"context"
)

// HashPassword satisfies our authenticator interface.
func (m *Authenticator) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockPasswordMatches(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMockPasswordMatches()

		expected := `
package example

import (
	"context"
)

// PasswordMatches satisfies our authenticator interface.
func (m *Authenticator) PasswordMatches(
	ctx context.Context,
	hashedPassword,
	providedPassword string,
	salt []byte,
) bool {
	args := m.Called(
		ctx,
		hashedPassword,
		providedPassword,
		salt,
	)
	return args.Bool(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
