package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := usersTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestV1Client_BuildGetUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		exampleUser := fake.BuildFakeUser()

		actual, err := c.BuildGetUserRequest(ctx, exampleUser.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleUser.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		// the hashed password is never transmitted over the wire.
		exampleUser.HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		exampleUser.TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		exampleUser.TwoFactorSecretVerifiedOn = nil

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleUser.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/users/%d", exampleUser.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleUser))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetUser(ctx, exampleUser.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleUser, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleUser.Salt = nil
		exampleUser.HashedPassword = ""

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUser(ctx, exampleUser.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetUsersRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetUsersRequest(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetUsers(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUserList := fake.BuildFakeUserList()
		// the hashed password is never transmitted over the wire.
		exampleUserList.Users[0].HashedPassword = ""
		exampleUserList.Users[1].HashedPassword = ""
		exampleUserList.Users[2].HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		exampleUserList.Users[0].TwoFactorSecret = ""
		exampleUserList.Users[1].TwoFactorSecret = ""
		exampleUserList.Users[2].TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		exampleUserList.Users[0].TwoFactorSecretVerifiedOn = nil
		exampleUserList.Users[1].TwoFactorSecretVerifiedOn = nil
		exampleUserList.Users[2].TwoFactorSecretVerifiedOn = nil

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/users", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleUserList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetUsers(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleUserList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUsers(ctx, nil)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateUserRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		expected := fake.BuildDatabaseCreationResponse(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/users", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.UserCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(expected))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateUser(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateUser(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		exampleUser := fake.BuildFakeUser()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveUserRequest(ctx, exampleUser.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleUser.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/users/%d", exampleUser.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveUser(ctx, exampleUser.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		err := buildTestClientWithInvalidURL(t).ArchiveUser(ctx, exampleUser.ID)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildLoginRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		req, err := c.BuildLoginRequest(ctx, exampleInput)
		require.NotNil(t, req)
		assert.Equal(t, req.Method, http.MethodPost)
		assert.NoError(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		req, err := c.BuildLoginRequest(ctx, nil)
		assert.Nil(t, req)
		assert.Error(t, err)
	})
}

func TestV1Client_Login(T *testing.T) {
	T.Parallel()

	const expectedPath = "/users/login"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					http.SetCookie(res, &http.Cookie{Name: exampleUser.Username})
				},
			),
		)
		c := buildTestClient(t, ts)

		cookie, err := c.Login(ctx, exampleInput)
		require.NotNil(t, cookie)
		assert.NoError(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		cookie, err := c.Login(ctx, nil)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		c := buildTestClientWithInvalidURL(t)

		cookie, err := c.Login(ctx, exampleInput)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})

	T.Run("with timeout", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
					time.Sleep(10 * time.Hour)
				},
			),
		)
		c := buildTestClient(t, ts)
		c.plainClient.Timeout = 500 * time.Microsecond

		cookie, err := c.Login(ctx, exampleInput)
		require.Nil(t, cookie)
		assert.Error(t, err)
	})

	T.Run("with missing cookie", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
				},
			),
		)
		c := buildTestClient(t, ts)

		cookie, err := c.Login(ctx, exampleInput)
		require.Nil(t, cookie)
		assert.Error(t, err)
	})
}

func TestV1Client_BuildValidateTOTPSecretRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		req, err := c.BuildVerifyTOTPSecretRequest(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.NoError(t, err)
		require.NotNil(t, req)
		assert.Equal(t, req.Method, http.MethodPost)
	})
}

func TestV1Client_ValidateTOTPSecret(T *testing.T) {
	T.Parallel()

	const expectedPath = "/users/totp_secret/verify"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					res.WriteHeader(http.StatusAccepted)
				},
			),
		)
		c := buildTestClient(t, ts)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.NoError(t, err)
	})

	T.Run("with bad request response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					res.WriteHeader(http.StatusBadRequest)
				},
			),
		)
		c := buildTestClient(t, ts)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTOTPToken, err)
	})

	T.Run("with otherwise invalid status code response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					res.WriteHeader(http.StatusInternalServerError)
				},
			),
		)
		c := buildTestClient(t, ts)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		c := buildTestClientWithInvalidURL(t)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	T.Run("with timeout", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					time.Sleep(10 * time.Minute)

					res.WriteHeader(http.StatusAccepted)
				},
			),
		)
		c := buildTestClient(t, ts)
		c.plainClient.Timeout = time.Millisecond

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildGetUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildGetUserRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildGetUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		exampleUser := fake.BuildFakeUser()

		actual, err := c.BuildGetUserRequest(ctx, exampleUser.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleUser.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientGetUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientGetUser(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		// the hashed password is never transmitted over the wire.
		exampleUser.HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		exampleUser.TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		exampleUser.TwoFactorSecretVerifiedOn = nil

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleUser.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/users/%d", exampleUser.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleUser))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetUser(ctx, exampleUser.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleUser, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleUser.Salt = nil
		exampleUser.HashedPassword = ""

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUser(ctx, exampleUser.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildGetUsersRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestV1ClientBuildGetUsersRequest()

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildGetUsersRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetUsersRequest(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientGetUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientGetUsers(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_GetUsers(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUserList := fake.BuildFakeUserList()
		// the hashed password is never transmitted over the wire.
		exampleUserList.Users[0].HashedPassword = ""
		exampleUserList.Users[1].HashedPassword = ""
		exampleUserList.Users[2].HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		exampleUserList.Users[0].TwoFactorSecret = ""
		exampleUserList.Users[1].TwoFactorSecret = ""
		exampleUserList.Users[2].TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		exampleUserList.Users[0].TwoFactorSecretVerifiedOn = nil
		exampleUserList.Users[1].TwoFactorSecretVerifiedOn = nil
		exampleUserList.Users[2].TwoFactorSecretVerifiedOn = nil

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/users", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleUserList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetUsers(ctx, nil)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleUserList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUsers(ctx, nil)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildCreateUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildCreateUserRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildCreateUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateUserRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientCreateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientCreateUser(proj)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)
		expected := fake.BuildDatabaseCreationResponse(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, "/users", "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.UserCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(expected))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateUser(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, expected, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserCreationInputFromUser(exampleUser)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateUser(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildArchiveUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildArchiveUserRequest(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildArchiveUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		exampleUser := fake.BuildFakeUser()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveUserRequest(ctx, exampleUser.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleUser.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientArchiveUser(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/users/%d", exampleUser.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveUser(ctx, exampleUser.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		err := buildTestClientWithInvalidURL(t).ArchiveUser(ctx, exampleUser.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildLoginRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildLoginRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildLoginRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		req, err := c.BuildLoginRequest(ctx, exampleInput)
		require.NotNil(t, req)
		assert.Equal(t, req.Method, http.MethodPost)
		assert.NoError(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		req, err := c.BuildLoginRequest(ctx, nil)
		assert.Nil(t, req)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientLogin(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestV1Client_Login(T *testing.T) {
	T.Parallel()

	const expectedPath = "/users/login"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					http.SetCookie(res, &http.Cookie{Name: exampleUser.Username})
				},
			),
		)
		c := buildTestClient(t, ts)

		cookie, err := c.Login(ctx, exampleInput)
		require.NotNil(t, cookie)
		assert.NoError(t, err)
	})

	T.Run("with nil input", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		cookie, err := c.Login(ctx, nil)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		c := buildTestClientWithInvalidURL(t)

		cookie, err := c.Login(ctx, exampleInput)
		assert.Nil(t, cookie)
		assert.Error(t, err)
	})

	T.Run("with timeout", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
					time.Sleep(10 * time.Hour)
				},
			),
		)
		c := buildTestClient(t, ts)
		c.plainClient.Timeout = 500 * time.Microsecond

		cookie, err := c.Login(ctx, exampleInput)
		require.Nil(t, cookie)
		assert.Error(t, err)
	})

	T.Run("with missing cookie", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserLoginInputFromUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)
				},
			),
		)
		c := buildTestClient(t, ts)

		cookie, err := c.Login(ctx, exampleInput)
		require.Nil(t, cookie)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildValidateTOTPSecretRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientBuildValidateTOTPSecretRequest(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildValidateTOTPSecretRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		req, err := c.BuildVerifyTOTPSecretRequest(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.NoError(t, err)
		require.NotNil(t, req)
		assert.Equal(t, req.Method, http.MethodPost)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientValidateTOTPSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestV1ClientValidateTOTPSecret(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestV1Client_ValidateTOTPSecret(T *testing.T) {
	T.Parallel()

	const expectedPath = "/users/totp_secret/verify"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					res.WriteHeader(http.StatusAccepted)
				},
			),
		)
		c := buildTestClient(t, ts)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.NoError(t, err)
	})

	T.Run("with bad request response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					res.WriteHeader(http.StatusBadRequest)
				},
			),
		)
		c := buildTestClient(t, ts)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTOTPToken, err)
	})

	T.Run("with otherwise invalid status code response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					res.WriteHeader(http.StatusInternalServerError)
				},
			),
		)
		c := buildTestClient(t, ts)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		c := buildTestClientWithInvalidURL(t)

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})

	T.Run("with timeout", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeTOTPSecretValidationInputForUser(exampleUser)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					time.Sleep(10 * time.Minute)

					res.WriteHeader(http.StatusAccepted)
				},
			),
		)
		c := buildTestClient(t, ts)
		c.plainClient.Timeout = time.Millisecond

		err := c.VerifyTOTPSecret(ctx, exampleUser.ID, exampleInput.TOTPToken)
		assert.Error(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
