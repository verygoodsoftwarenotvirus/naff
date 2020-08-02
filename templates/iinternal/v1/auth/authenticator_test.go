package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_authenticatorDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := authenticatorDotGo(proj)

		expected := `
package example

import (
	"context"
	"crypto/rand"
	"errors"
	wire "github.com/google/wire"
)

var (
	// ErrInvalidTwoFactorCode indicates that a provided two factor code is invalid.
	ErrInvalidTwoFactorCode = errors.New("invalid two factor code")
	// ErrPasswordHashTooWeak indicates that a provided password hash is too weak.
	ErrPasswordHashTooWeak = errors.New("password's hash is too weak")

	// Providers represents what this package offers to external libraries in the way of constructors.
	Providers = wire.NewSet(
		ProvideBcryptAuthenticator,
		ProvideBcryptHashCost,
	)
)

// ProvideBcryptHashCost provides a BcryptHashCost.
func ProvideBcryptHashCost() BcryptHashCost {
	return DefaultBcryptHashCost
}

type (
	// PasswordHasher hashes passwords.
	PasswordHasher interface {
		PasswordIsAcceptable(password string) bool
		HashPassword(ctx context.Context, password string) (string, error)
		PasswordMatches(ctx context.Context, hashedPassword, providedPassword string, salt []byte) bool
	}

	// Authenticator is a poorly named Authenticator interface.
	Authenticator interface {
		PasswordHasher

		ValidateLogin(
			ctx context.Context,
			HashedPassword,
			ProvidedPassword,
			TwoFactorSecret,
			TwoFactorCode string,
			Salt []byte,
		) (valid bool, err error)
	}
)

// we run this function to ensure that we have no problem reading from crypto/rand
func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthenticatorVariableDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildAuthenticatorVariableDeclarations()

		expected := `
package example

import (
	"errors"
	wire "github.com/google/wire"
)

var (
	// ErrInvalidTwoFactorCode indicates that a provided two factor code is invalid.
	ErrInvalidTwoFactorCode = errors.New("invalid two factor code")
	// ErrPasswordHashTooWeak indicates that a provided password hash is too weak.
	ErrPasswordHashTooWeak = errors.New("password's hash is too weak")

	// Providers represents what this package offers to external libraries in the way of constructors.
	Providers = wire.NewSet(
		ProvideBcryptAuthenticator,
		ProvideBcryptHashCost,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideBcryptHashCost(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildProvideBcryptHashCost()

		expected := `
package example

import ()

// ProvideBcryptHashCost provides a BcryptHashCost.
func ProvideBcryptHashCost() BcryptHashCost {
	return DefaultBcryptHashCost
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAuthenticatorTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildAuthenticatorTypeDefinitions()

		expected := `
package example

import (
	"context"
)

type (
	// PasswordHasher hashes passwords.
	PasswordHasher interface {
		PasswordIsAcceptable(password string) bool
		HashPassword(ctx context.Context, password string) (string, error)
		PasswordMatches(ctx context.Context, hashedPassword, providedPassword string, salt []byte) bool
	}

	// Authenticator is a poorly named Authenticator interface.
	Authenticator interface {
		PasswordHasher

		ValidateLogin(
			ctx context.Context,
			HashedPassword,
			ProvidedPassword,
			TwoFactorSecret,
			TwoFactorCode string,
			Salt []byte,
		) (valid bool, err error)
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInit()

		expected := `
package example

import (
	"crypto/rand"
)

// we run this function to ensure that we have no problem reading from crypto/rand
func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
