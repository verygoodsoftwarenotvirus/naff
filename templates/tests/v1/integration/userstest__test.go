package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := usersTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
	"time"
)

func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}

// randString produces a random string.
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func randString() (string, error) {
	b := make([]byte, 64)
	// Note that err == nil only if we read len(b) bytes
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(b), nil
}

func buildDummyUser(ctx context.Context, t *testing.T) (*v1.UserCreationResponse, *v1.UserCreationInput, *http.Cookie) {
	t.Helper()

	// build user creation route input.
	userInput := fake.BuildFakeUserCreationInput()
	user, err := todoClient.CreateUser(ctx, userInput)
	require.NotNil(t, user)
	require.NoError(t, err)

	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	require.NoError(t, err)
	require.NoError(t, todoClient.VerifyTOTPSecret(ctx, user.ID, token))

	cookie := loginUser(ctx, t, userInput.Username, userInput.Password, user.TwoFactorSecret)

	require.NoError(t, err)
	require.NotNil(t, cookie)

	return user, userInput, cookie
}

func checkUserCreationEquality(t *testing.T, expected *v1.UserCreationInput, actual *v1.UserCreationResponse) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotEmpty(t, actual.TwoFactorSecret)
	assert.NotZero(t, actual.CreatedOn)
	assert.Nil(t, actual.LastUpdatedOn)
	assert.Nil(t, actual.ArchivedOn)
}

func checkUserEquality(t *testing.T, expected *v1.UserCreationInput, actual *v1.User) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotZero(t, actual.CreatedOn)
	assert.Nil(t, actual.LastUpdatedOn)
	assert.Nil(t, actual.ArchivedOn)
}

func TestUsers(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be creatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fake.BuildFakeUserCreationInput()
			actual, err := todoClient.CreateUser(ctx, exampleUserInput)
			checkValueAndError(t, actual, err)

			// Assert user equality.
			checkUserCreationEquality(t, exampleUserInput, actual)

			// Clean up.
			assert.NoError(t, todoClient.ArchiveUser(ctx, actual.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch user.
			actual, err := todoClient.GetUser(ctx, nonexistentID)
			assert.Nil(t, actual)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fake.BuildFakeUserCreationInput()
			premade, err := todoClient.CreateUser(ctx, exampleUserInput)
			checkValueAndError(t, premade, err)
			assert.NotEmpty(t, premade.TwoFactorSecret)

			secretVerificationToken, err := totp.GenerateCode(premade.TwoFactorSecret, time.Now().UTC())
			checkValueAndError(t, secretVerificationToken, err)

			assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, premade.ID, secretVerificationToken))

			// Fetch user.
			actual, err := todoClient.GetUser(ctx, premade.ID)
			if err != nil {
				t.Logf("error encountered trying to fetch user %q: %v\n", premade.Username, err)
			}
			checkValueAndError(t, actual, err)

			// Assert user equality.
			checkUserEquality(t, exampleUserInput, actual)

			// Clean up.
			assert.NoError(t, todoClient.ArchiveUser(ctx, actual.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fake.BuildFakeUserCreationInput()
			u, err := todoClient.CreateUser(ctx, exampleUserInput)
			assert.NoError(t, err)
			assert.NotNil(t, u)

			if u == nil || err != nil {
				t.Log("something has gone awry, user returned is nil")
				t.FailNow()
			}

			// Execute.
			err = todoClient.ArchiveUser(ctx, u.ID)
			assert.NoError(t, err)
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersTestsInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUsersTestsInit()

		expected := `
package example

import (
	"crypto/rand"
)

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

func Test_buildUsersTestsRandString(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUsersTestsRandString()

		expected := `
package example

import (
	"crypto/rand"
	"encoding/base32"
)

// randString produces a random string.
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func randString() (string, error) {
	b := make([]byte, 64)
	// Note that err == nil only if we read len(b) bytes
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

func Test_buildUsersTestsBuildDummyUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildUsersTestsBuildDummyUser(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"testing"
	"time"
)

func buildDummyUser(ctx context.Context, t *testing.T) (*v1.UserCreationResponse, *v1.UserCreationInput, *http.Cookie) {
	t.Helper()

	// build user creation route input.
	userInput := fake.BuildFakeUserCreationInput()
	user, err := todoClient.CreateUser(ctx, userInput)
	require.NotNil(t, user)
	require.NoError(t, err)

	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	require.NoError(t, err)
	require.NoError(t, todoClient.VerifyTOTPSecret(ctx, user.ID, token))

	cookie := loginUser(ctx, t, userInput.Username, userInput.Password, user.TwoFactorSecret)

	require.NoError(t, err)
	require.NotNil(t, cookie)

	return user, userInput, cookie
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersTestsCheckUserCreationEquality(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildUsersTestsCheckUserCreationEquality(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func checkUserCreationEquality(t *testing.T, expected *v1.UserCreationInput, actual *v1.UserCreationResponse) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotEmpty(t, actual.TwoFactorSecret)
	assert.NotZero(t, actual.CreatedOn)
	assert.Nil(t, actual.LastUpdatedOn)
	assert.Nil(t, actual.ArchivedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersTestsCheckUserEquality(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildUsersTestsCheckUserEquality(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func checkUserEquality(t *testing.T, expected *v1.UserCreationInput, actual *v1.User) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.NotZero(t, actual.CreatedOn)
	assert.Nil(t, actual.LastUpdatedOn)
	assert.Nil(t, actual.ArchivedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersTestsTestUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildUsersTestsTestUsers(proj)

		expected := `
package example

import (
	"context"
	totp "github.com/pquerna/otp/totp"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func TestUsers(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be creatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fake.BuildFakeUserCreationInput()
			actual, err := todoClient.CreateUser(ctx, exampleUserInput)
			checkValueAndError(t, actual, err)

			// Assert user equality.
			checkUserCreationEquality(t, exampleUserInput, actual)

			// Clean up.
			assert.NoError(t, todoClient.ArchiveUser(ctx, actual.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch user.
			actual, err := todoClient.GetUser(ctx, nonexistentID)
			assert.Nil(t, actual)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fake.BuildFakeUserCreationInput()
			premade, err := todoClient.CreateUser(ctx, exampleUserInput)
			checkValueAndError(t, premade, err)
			assert.NotEmpty(t, premade.TwoFactorSecret)

			secretVerificationToken, err := totp.GenerateCode(premade.TwoFactorSecret, time.Now().UTC())
			checkValueAndError(t, secretVerificationToken, err)

			assert.NoError(t, todoClient.VerifyTOTPSecret(ctx, premade.ID, secretVerificationToken))

			// Fetch user.
			actual, err := todoClient.GetUser(ctx, premade.ID)
			if err != nil {
				t.Logf("error encountered trying to fetch user %q: %v\n", premade.Username, err)
			}
			checkValueAndError(t, actual, err)

			// Assert user equality.
			checkUserEquality(t, exampleUserInput, actual)

			// Clean up.
			assert.NoError(t, todoClient.ArchiveUser(ctx, actual.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create user.
			exampleUserInput := fake.BuildFakeUserCreationInput()
			u, err := todoClient.CreateUser(ctx, exampleUserInput)
			assert.NoError(t, err)
			assert.NotNil(t, u)

			if u == nil || err != nil {
				t.Log("something has gone awry, user returned is nil")
				t.FailNow()
			}

			// Execute.
			err = todoClient.ArchiveUser(ctx, u.ID)
			assert.NoError(t, err)
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
