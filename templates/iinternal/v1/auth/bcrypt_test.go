package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_bcryptDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := bcryptDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	totp "github.com/pquerna/otp/totp"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	bcrypt "golang.org/x/crypto/bcrypt"
	"math"
)

const (
	bcryptCostCompensation     = 2
	defaultMinimumPasswordSize = 16

	// DefaultBcryptHashCost is what it says on the tin.
	DefaultBcryptHashCost = BcryptHashCost(bcrypt.DefaultCost + bcryptCostCompensation)
)

var (
	_ Authenticator = (*BcryptAuthenticator)(nil)

	// ErrCostTooLow indicates that a password has too low a Bcrypt cost.
	ErrCostTooLow = errors.New("stored password's cost is too low")
)

type (
	// BcryptAuthenticator is our bcrypt-based authenticator
	BcryptAuthenticator struct {
		logger              v1.Logger
		hashCost            uint
		minimumPasswordSize uint
	}

	// BcryptHashCost is an arbitrary type alias for dependency injection's sake.
	BcryptHashCost uint
)

// ProvideBcryptAuthenticator returns a bcrypt powered Authenticator.
func ProvideBcryptAuthenticator(hashCost BcryptHashCost, logger v1.Logger) Authenticator {
	ba := &BcryptAuthenticator{
		logger:              logger.WithName("bcrypt"),
		hashCost:            uint(math.Min(float64(DefaultBcryptHashCost), float64(hashCost))),
		minimumPasswordSize: defaultMinimumPasswordSize,
	}
	return ba
}

// HashPassword takes a password and hashes it using bcrypt.
func (b *BcryptAuthenticator) HashPassword(ctx context.Context, password string) (string, error) {
	_, span := tracing.StartSpan(ctx, "HashPassword")
	defer span.End()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), int(b.hashCost))
	return string(hashedPass), err
}

// ValidateLogin validates a login attempt by:
// 1. checking that the provided password matches the stored hashed password
// 2. checking that the temporary one-time password provided jives with the stored two factor secret
// 3. checking that the provided hashed password isn't too weak, and returning an error otherwise
func (b *BcryptAuthenticator) ValidateLogin(
	ctx context.Context,
	hashedPassword,
	providedPassword,
	twoFactorSecret,
	twoFactorCode string,
	_ []byte,
) (passwordMatches bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidateLogin")
	defer span.End()

	passwordMatches = b.PasswordMatches(ctx, hashedPassword, providedPassword, nil)
	tooWeak := b.hashedPasswordIsTooWeak(ctx, hashedPassword)

	if !totp.Validate(twoFactorCode, twoFactorSecret) {
		b.logger.WithValues(map[string]interface{}{
			"password_matches": passwordMatches,
			"2fa_secret":       twoFactorSecret,
			"provided_code":    twoFactorCode,
		}).Debug("invalid code provided")

		return passwordMatches, ErrInvalidTwoFactorCode
	}

	if tooWeak {
		// NOTE: this can end up with a return set where passwordMatches is true and the err is not nil.
		// This is the valid case in the event the user has logged in with a valid password, but the
		// bcrypt cost has been raised since they last logged in.
		return passwordMatches, ErrCostTooLow
	}

	return passwordMatches, nil
}

// PasswordMatches validates whether or not a bcrypt-hashed password matches a provided password
func (b *BcryptAuthenticator) PasswordMatches(ctx context.Context, hashedPassword, providedPassword string, _ []byte) bool {
	_, span := tracing.StartSpan(ctx, "PasswordMatches")
	defer span.End()

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword)) == nil
}

// hashedPasswordIsTooWeak determines if a given hashed password was hashed with too weak a bcrypt cost.
func (b *BcryptAuthenticator) hashedPasswordIsTooWeak(ctx context.Context, hashedPassword string) bool {
	_, span := tracing.StartSpan(ctx, "hashedPasswordIsTooWeak")
	defer span.End()

	cost, err := bcrypt.Cost([]byte(hashedPassword))

	return err != nil || uint(cost) < b.hashCost
}

// PasswordIsAcceptable takes a password and returns whether or not it satisfies the authenticator.
func (b *BcryptAuthenticator) PasswordIsAcceptable(pass string) bool {
	return uint(len(pass)) >= b.minimumPasswordSize
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBcryptConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBcryptConstDeclarations()

		expected := `
package example

import (
	bcrypt "golang.org/x/crypto/bcrypt"
)

const (
	bcryptCostCompensation     = 2
	defaultMinimumPasswordSize = 16

	// DefaultBcryptHashCost is what it says on the tin.
	DefaultBcryptHashCost = BcryptHashCost(bcrypt.DefaultCost + bcryptCostCompensation)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBcryptVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBcryptVarDeclarations()

		expected := `
package example

import (
	"errors"
)

var (
	_ Authenticator = (*BcryptAuthenticator)(nil)

	// ErrCostTooLow indicates that a password has too low a Bcrypt cost.
	ErrCostTooLow = errors.New("stored password's cost is too low")
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBcryptTypeDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBcryptTypeDeclarations()

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

type (
	// BcryptAuthenticator is our bcrypt-based authenticator
	BcryptAuthenticator struct {
		logger              v1.Logger
		hashCost            uint
		minimumPasswordSize uint
	}

	// BcryptHashCost is an arbitrary type alias for dependency injection's sake.
	BcryptHashCost uint
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideBcryptAuthenticator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideBcryptAuthenticator()

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"math"
)

// ProvideBcryptAuthenticator returns a bcrypt powered Authenticator.
func ProvideBcryptAuthenticator(hashCost BcryptHashCost, logger v1.Logger) Authenticator {
	ba := &BcryptAuthenticator{
		logger:              logger.WithName("bcrypt"),
		hashCost:            uint(math.Min(float64(DefaultBcryptHashCost), float64(hashCost))),
		minimumPasswordSize: defaultMinimumPasswordSize,
	}
	return ba
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildHashPassword(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildHashPassword(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	bcrypt "golang.org/x/crypto/bcrypt"
)

// HashPassword takes a password and hashes it using bcrypt.
func (b *BcryptAuthenticator) HashPassword(ctx context.Context, password string) (string, error) {
	_, span := tracing.StartSpan(ctx, "HashPassword")
	defer span.End()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), int(b.hashCost))
	return string(hashedPass), err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildValidateLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildValidateLogin(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ValidateLogin validates a login attempt by:
// 1. checking that the provided password matches the stored hashed password
// 2. checking that the temporary one-time password provided jives with the stored two factor secret
// 3. checking that the provided hashed password isn't too weak, and returning an error otherwise
func (b *BcryptAuthenticator) ValidateLogin(
	ctx context.Context,
	hashedPassword,
	providedPassword,
	twoFactorSecret,
	twoFactorCode string,
	_ []byte,
) (passwordMatches bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidateLogin")
	defer span.End()

	passwordMatches = b.PasswordMatches(ctx, hashedPassword, providedPassword, nil)
	tooWeak := b.hashedPasswordIsTooWeak(ctx, hashedPassword)

	if !totp.Validate(twoFactorCode, twoFactorSecret) {
		b.logger.WithValues(map[string]interface{}{
			"password_matches": passwordMatches,
			"2fa_secret":       twoFactorSecret,
			"provided_code":    twoFactorCode,
		}).Debug("invalid code provided")

		return passwordMatches, ErrInvalidTwoFactorCode
	}

	if tooWeak {
		// NOTE: this can end up with a return set where passwordMatches is true and the err is not nil.
		// This is the valid case in the event the user has logged in with a valid password, but the
		// bcrypt cost has been raised since they last logged in.
		return passwordMatches, ErrCostTooLow
	}

	return passwordMatches, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildPasswordMatches(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildPasswordMatches(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	bcrypt "golang.org/x/crypto/bcrypt"
)

// PasswordMatches validates whether or not a bcrypt-hashed password matches a provided password
func (b *BcryptAuthenticator) PasswordMatches(ctx context.Context, hashedPassword, providedPassword string, _ []byte) bool {
	_, span := tracing.StartSpan(ctx, "PasswordMatches")
	defer span.End()

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword)) == nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildHashedPasswordIsTooWeak(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildHashedPasswordIsTooWeak(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	bcrypt "golang.org/x/crypto/bcrypt"
)

// hashedPasswordIsTooWeak determines if a given hashed password was hashed with too weak a bcrypt cost.
func (b *BcryptAuthenticator) hashedPasswordIsTooWeak(ctx context.Context, hashedPassword string) bool {
	_, span := tracing.StartSpan(ctx, "hashedPasswordIsTooWeak")
	defer span.End()

	cost, err := bcrypt.Cost([]byte(hashedPassword))

	return err != nil || uint(cost) < b.hashCost
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildPasswordIsAcceptable(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildPasswordIsAcceptable()

		expected := `
package example

import ()

// PasswordIsAcceptable takes a password and returns whether or not it satisfies the authenticator.
func (b *BcryptAuthenticator) PasswordIsAcceptable(pass string) bool {
	return uint(len(pass)) >= b.minimumPasswordSize
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
