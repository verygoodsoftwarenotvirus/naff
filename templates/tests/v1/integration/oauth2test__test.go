package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2TestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2TestDotGo(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
	"time"
)

func mustBuildCode(t *testing.T, totpSecret string) string {
	t.Helper()
	code, err := totp.GenerateCode(totpSecret, time.Now().UTC())
	require.NoError(t, err)
	return code
}

func buildDummyOAuth2ClientInput(t *testing.T, username, password, totpToken string) *v1.OAuth2ClientCreationInput {
	t.Helper()

	x := &v1.OAuth2ClientCreationInput{
		UserLoginInput: v1.UserLoginInput{
			Username:  username,
			Password:  password,
			TOTPToken: mustBuildCode(t, totpToken),
		},
		Scopes:      []string{"*"},
		RedirectURI: "http://localhost",
	}

	return x
}

func convertInputToClient(input *v1.OAuth2ClientCreationInput) *v1.OAuth2Client {
	return &v1.OAuth2Client{
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
}

func checkOAuth2ClientEquality(t *testing.T, expected, actual *v1.OAuth2Client) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.NotEmpty(t, actual.ClientID)
	assert.NotEmpty(t, actual.ClientSecret)
	assert.Equal(t, expected.RedirectURI, actual.RedirectURI)
	assert.Equal(t, expected.Scopes, actual.Scopes)
	assert.NotZero(t, actual.CreatedOn)
	assert.Nil(t, actual.ArchivedOn)
}

func TestOAuth2Clients(test *testing.T) {
	_ctx := context.Background()

	// create user.
	x, y, cookie := buildDummyUser(_ctx, test)
	assert.NotNil(test, cookie)

	input := buildDummyOAuth2ClientInput(test, x.Username, y.Password, x.TwoFactorSecret)
	premade, err := todoClient.CreateOAuth2Client(_ctx, cookie, input)
	checkValueAndError(test, premade, err)

	testClient, err := http.NewClient(
		_ctx,
		premade.ClientID,
		premade.ClientSecret,
		todoClient.URL,
		noop.ProvideNoopLogger(),
		todoClient.PlainClient(),
		premade.Scopes,
		debug,
	)
	require.NoError(test, err, "error setting up auxiliary client")

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be creatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Client.
			actual, err := testClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(t, actual, err)

			// Assert oauth2Client equality.
			checkOAuth2ClientEquality(t, convertInputToClient(input), actual)

			// Clean up.
			err = testClient.ArchiveOAuth2Client(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read one that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch oauth2Client.
			_, err := testClient.GetOAuth2Client(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Client.
			input := buildDummyOAuth2ClientInput(t, x.Username, y.Password, x.TwoFactorSecret)
			c, err := testClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(t, c, err)

			// Fetch oauth2Client.
			actual, err := testClient.GetOAuth2Client(ctx, c.ID)
			checkValueAndError(t, actual, err)

			// Assert oauth2Client equality.
			checkOAuth2ClientEquality(t, convertInputToClient(input), actual)

			// Clean up.
			err = testClient.ArchiveOAuth2Client(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Client.
			input := buildDummyOAuth2ClientInput(t, x.Username, y.Password, x.TwoFactorSecret)
			premade, err := testClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(t, premade, err)

			// Clean up.
			err = testClient.ArchiveOAuth2Client(ctx, premade.ID)
			assert.NoError(t, err)
		})

		T.Run("should be unable to authorize after being deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// create user.
			createdUser, createdUserInput, _ := buildDummyUser(ctx, test)
			assert.NotNil(test, cookie)

			input := buildDummyOAuth2ClientInput(test, createdUserInput.Username, createdUserInput.Password, createdUser.TwoFactorSecret)
			premade, err := todoClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(test, premade, err)

			// archive oauth2Client.
			require.NoError(t, testClient.ArchiveOAuth2Client(ctx, premade.ID))

			c2, err := http.NewClient(
				ctx,
				premade.ClientID,
				premade.ClientSecret,
				todoClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				premade.Scopes,
				true,
			)
			checkValueAndError(test, c2, err)

			_, err = c2.GetOAuth2Clients(ctx, nil)
			assert.Error(t, err, "expected error from what should be an unauthorized client")
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Clients.
			var expected []*v1.OAuth2Client
			for i := 0; i < 5; i++ {
				input := buildDummyOAuth2ClientInput(t, x.Username, y.Password, x.TwoFactorSecret)
				oac, err := testClient.CreateOAuth2Client(ctx, cookie, input)
				checkValueAndError(t, oac, err)
				expected = append(expected, oac)
			}

			// Assert oauth2Client list equality.
			actual, err := testClient.GetOAuth2Clients(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(actual.Clients)-len(expected) > 0,
				"expected %d - %d to be > 0",
				len(actual.Clients),
				len(expected),
			)

			for _, oAuth2Client := range expected {
				clientFound := false
				for _, c := range actual.Clients {
					if c.ID == oAuth2Client.ID {
						clientFound = true
						break
					}
				}
				assert.True(t, clientFound, "expected oAuth2Client ID %d to be present in results", oAuth2Client.ID)
			}

			// Clean up.
			for _, oa2c := range expected {
				err = testClient.ArchiveOAuth2Client(ctx, oa2c.ID)
				assert.NoError(t, err, "error deleting client %d: %v", oa2c.ID, err)
			}
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2TestMustBuildCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2TestMustBuildCode()

		expected := `
package example

import (
	totp "github.com/pquerna/otp/totp"
	require "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func mustBuildCode(t *testing.T, totpSecret string) string {
	t.Helper()
	code, err := totp.GenerateCode(totpSecret, time.Now().UTC())
	require.NoError(t, err)
	return code
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2TestBuildDummyOAuth2ClientInput(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildOAuth2TestBuildDummyOAuth2ClientInput(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func buildDummyOAuth2ClientInput(t *testing.T, username, password, totpToken string) *v1.OAuth2ClientCreationInput {
	t.Helper()

	x := &v1.OAuth2ClientCreationInput{
		UserLoginInput: v1.UserLoginInput{
			Username:  username,
			Password:  password,
			TOTPToken: mustBuildCode(t, totpToken),
		},
		Scopes:      []string{"*"},
		RedirectURI: "http://localhost",
	}

	return x
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2TestConvertInputToClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildOAuth2TestConvertInputToClient(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

func convertInputToClient(input *v1.OAuth2ClientCreationInput) *v1.OAuth2Client {
	return &v1.OAuth2Client{
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2TestCheckOAuth2ClientEquality(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildOAuth2TestCheckOAuth2ClientEquality(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func checkOAuth2ClientEquality(t *testing.T, expected, actual *v1.OAuth2Client) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.NotEmpty(t, actual.ClientID)
	assert.NotEmpty(t, actual.ClientSecret)
	assert.Equal(t, expected.RedirectURI, actual.RedirectURI)
	assert.Equal(t, expected.Scopes, actual.Scopes)
	assert.NotZero(t, actual.CreatedOn)
	assert.Nil(t, actual.ArchivedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2TestTestOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildOAuth2TestTestOAuth2Clients(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func TestOAuth2Clients(test *testing.T) {
	_ctx := context.Background()

	// create user.
	x, y, cookie := buildDummyUser(_ctx, test)
	assert.NotNil(test, cookie)

	input := buildDummyOAuth2ClientInput(test, x.Username, y.Password, x.TwoFactorSecret)
	premade, err := todoClient.CreateOAuth2Client(_ctx, cookie, input)
	checkValueAndError(test, premade, err)

	testClient, err := http.NewClient(
		_ctx,
		premade.ClientID,
		premade.ClientSecret,
		todoClient.URL,
		noop.ProvideNoopLogger(),
		todoClient.PlainClient(),
		premade.Scopes,
		debug,
	)
	require.NoError(test, err, "error setting up auxiliary client")

	test.Run("Creating", func(T *testing.T) {
		T.Run("should be creatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Client.
			actual, err := testClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(t, actual, err)

			// Assert oauth2Client equality.
			checkOAuth2ClientEquality(t, convertInputToClient(input), actual)

			// Clean up.
			err = testClient.ArchiveOAuth2Client(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read one that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch oauth2Client.
			_, err := testClient.GetOAuth2Client(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Client.
			input := buildDummyOAuth2ClientInput(t, x.Username, y.Password, x.TwoFactorSecret)
			c, err := testClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(t, c, err)

			// Fetch oauth2Client.
			actual, err := testClient.GetOAuth2Client(ctx, c.ID)
			checkValueAndError(t, actual, err)

			// Assert oauth2Client equality.
			checkOAuth2ClientEquality(t, convertInputToClient(input), actual)

			// Clean up.
			err = testClient.ArchiveOAuth2Client(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Client.
			input := buildDummyOAuth2ClientInput(t, x.Username, y.Password, x.TwoFactorSecret)
			premade, err := testClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(t, premade, err)

			// Clean up.
			err = testClient.ArchiveOAuth2Client(ctx, premade.ID)
			assert.NoError(t, err)
		})

		T.Run("should be unable to authorize after being deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// create user.
			createdUser, createdUserInput, _ := buildDummyUser(ctx, test)
			assert.NotNil(test, cookie)

			input := buildDummyOAuth2ClientInput(test, createdUserInput.Username, createdUserInput.Password, createdUser.TwoFactorSecret)
			premade, err := todoClient.CreateOAuth2Client(ctx, cookie, input)
			checkValueAndError(test, premade, err)

			// archive oauth2Client.
			require.NoError(t, testClient.ArchiveOAuth2Client(ctx, premade.ID))

			c2, err := http.NewClient(
				ctx,
				premade.ClientID,
				premade.ClientSecret,
				todoClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				premade.Scopes,
				true,
			)
			checkValueAndError(test, c2, err)

			_, err = c2.GetOAuth2Clients(ctx, nil)
			assert.Error(t, err, "expected error from what should be an unauthorized client")
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create oauth2Clients.
			var expected []*v1.OAuth2Client
			for i := 0; i < 5; i++ {
				input := buildDummyOAuth2ClientInput(t, x.Username, y.Password, x.TwoFactorSecret)
				oac, err := testClient.CreateOAuth2Client(ctx, cookie, input)
				checkValueAndError(t, oac, err)
				expected = append(expected, oac)
			}

			// Assert oauth2Client list equality.
			actual, err := testClient.GetOAuth2Clients(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(actual.Clients)-len(expected) > 0,
				"expected %d - %d to be > 0",
				len(actual.Clients),
				len(expected),
			)

			for _, oAuth2Client := range expected {
				clientFound := false
				for _, c := range actual.Clients {
					if c.ID == oAuth2Client.ID {
						clientFound = true
						break
					}
				}
				assert.True(t, clientFound, "expected oAuth2Client ID %d to be present in results", oAuth2Client.ID)
			}

			// Clean up.
			for _, oa2c := range expected {
				err = testClient.ArchiveOAuth2Client(ctx, oa2c.ID)
				assert.NoError(t, err, "error deleting client %d: %v", oa2c.ID, err)
			}
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
