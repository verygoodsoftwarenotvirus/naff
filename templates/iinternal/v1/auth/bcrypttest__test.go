package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_bcryptTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := bcryptTestDotGo(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
	"time"
)

const (
	examplePassword             = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
	weaklyHashedExamplePassword = "$2a$04$7G7dHZe7MeWjOMsYKO8uCu/CRKnDMMBHOfXaB6YgyQL/cl8nhwf/2"
	hashedExamplePassword       = "$2a$13$hxMAo/ZRDmyaWcwvIem/vuUJkmeNytg3rwHUj6bRZR1d/cQHXjFvW"
	exampleTwoFactorSecret      = "HEREISASECRETWHICHIVEMADEUPBECAUSEIWANNATESTRELIABLY"
)

func TestBcrypt_HashPassword(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual, err := x.HashPassword(ctx, "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
	})
}

func TestBcrypt_PasswordMatches(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("normal usage", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual := x.PasswordMatches(ctx, hashedExamplePassword, examplePassword, nil)
		assert.True(t, actual)
	})

	T.Run("when passwords don't match", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual := x.PasswordMatches(ctx, hashedExamplePassword, "password", nil)
		assert.False(t, actual)
	})
}

func TestBcrypt_PasswordIsAcceptable(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		assert.True(t, x.PasswordIsAcceptable(examplePassword))
		assert.False(t, x.PasswordIsAcceptable("hi there"))
	})
}

func TestBcrypt_ValidateLogin(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		code, err := totp.GenerateCode(exampleTwoFactorSecret, time.Now().UTC())
		assert.NoError(t, err, "error generating code to validate login")

		valid, err := x.ValidateLogin(
			ctx,
			hashedExamplePassword,
			examplePassword,
			exampleTwoFactorSecret,
			code,
			nil,
		)
		assert.NoError(t, err, "unexpected error encountered validating login: %v", err)
		assert.True(t, valid)
	})

	T.Run("with weak hash", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		code, err := totp.GenerateCode(exampleTwoFactorSecret, time.Now().UTC())
		assert.NoError(t, err, "error generating code to validate login")

		valid, err := x.ValidateLogin(
			ctx,
			weaklyHashedExamplePassword,
			examplePassword,
			exampleTwoFactorSecret,
			code,
			nil,
		)
		assert.Error(t, err, "unexpected error encountered validating login: %v", err)
		assert.True(t, valid)
	})

	T.Run("with non-matching password", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		code, err := totp.GenerateCode(exampleTwoFactorSecret, time.Now().UTC())
		assert.NoError(t, err, "error generating code to validate login")

		valid, err := x.ValidateLogin(
			ctx,
			hashedExamplePassword,
			"examplePassword",
			exampleTwoFactorSecret,
			code,
			nil,
		)
		assert.NoError(t, err, "unexpected error encountered validating login: %v", err)
		assert.False(t, valid)
	})

	T.Run("with invalid code", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		valid, err := x.ValidateLogin(
			ctx,
			hashedExamplePassword,
			examplePassword,
			exampleTwoFactorSecret,
			"CODE",
			nil,
		)
		assert.Error(t, err, "unexpected error encountered validating login: %v", err)
		assert.True(t, valid)
	})
}

func TestProvideBcrypt(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConstDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildConstDefinitions()

		expected := `
package example

import ()

const (
	examplePassword             = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
	weaklyHashedExamplePassword = "$2a$04$7G7dHZe7MeWjOMsYKO8uCu/CRKnDMMBHOfXaB6YgyQL/cl8nhwf/2"
	hashedExamplePassword       = "$2a$13$hxMAo/ZRDmyaWcwvIem/vuUJkmeNytg3rwHUj6bRZR1d/cQHXjFvW"
	exampleTwoFactorSecret      = "HEREISASECRETWHICHIVEMADEUPBECAUSEIWANNATESTRELIABLY"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBcrypt_HashPassword(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestBcrypt_HashPassword(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
)

func TestBcrypt_HashPassword(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual, err := x.HashPassword(ctx, "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBcrypt_PasswordMatches(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestBcrypt_PasswordMatches(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
)

func TestBcrypt_PasswordMatches(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("normal usage", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual := x.PasswordMatches(ctx, hashedExamplePassword, examplePassword, nil)
		assert.True(t, actual)
	})

	T.Run("when passwords don't match", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		actual := x.PasswordMatches(ctx, hashedExamplePassword, "password", nil)
		assert.False(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBcrypt_PasswordIsAcceptable(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestBcrypt_PasswordIsAcceptable(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
)

func TestBcrypt_PasswordIsAcceptable(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		assert.True(t, x.PasswordIsAcceptable(examplePassword))
		assert.False(t, x.PasswordIsAcceptable("hi there"))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBcrypt_ValidateLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestBcrypt_ValidateLogin(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
	"time"
)

func TestBcrypt_ValidateLogin(T *testing.T) {
	T.Parallel()

	x := auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		code, err := totp.GenerateCode(exampleTwoFactorSecret, time.Now().UTC())
		assert.NoError(t, err, "error generating code to validate login")

		valid, err := x.ValidateLogin(
			ctx,
			hashedExamplePassword,
			examplePassword,
			exampleTwoFactorSecret,
			code,
			nil,
		)
		assert.NoError(t, err, "unexpected error encountered validating login: %v", err)
		assert.True(t, valid)
	})

	T.Run("with weak hash", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		code, err := totp.GenerateCode(exampleTwoFactorSecret, time.Now().UTC())
		assert.NoError(t, err, "error generating code to validate login")

		valid, err := x.ValidateLogin(
			ctx,
			weaklyHashedExamplePassword,
			examplePassword,
			exampleTwoFactorSecret,
			code,
			nil,
		)
		assert.Error(t, err, "unexpected error encountered validating login: %v", err)
		assert.True(t, valid)
	})

	T.Run("with non-matching password", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		code, err := totp.GenerateCode(exampleTwoFactorSecret, time.Now().UTC())
		assert.NoError(t, err, "error generating code to validate login")

		valid, err := x.ValidateLogin(
			ctx,
			hashedExamplePassword,
			"examplePassword",
			exampleTwoFactorSecret,
			code,
			nil,
		)
		assert.NoError(t, err, "unexpected error encountered validating login: %v", err)
		assert.False(t, valid)
	})

	T.Run("with invalid code", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		valid, err := x.ValidateLogin(
			ctx,
			hashedExamplePassword,
			examplePassword,
			exampleTwoFactorSecret,
			"CODE",
			nil,
		)
		assert.Error(t, err, "unexpected error encountered validating login: %v", err)
		assert.True(t, valid)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideBcrypt(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestProvideBcrypt(proj)

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	"testing"
)

func TestProvideBcrypt(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		auth.ProvideBcryptAuthenticator(auth.DefaultBcryptHashCost, noop.ProvideNoopLogger())
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
