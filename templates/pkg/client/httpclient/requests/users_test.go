package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := usersDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strconv"
)

const usersBasePath = "users"

// BuildGetUserRequest builds an HTTP request for fetching a user.
func (c *V1Client) BuildGetUserRequest(ctx context.Context, userID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetUserRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath, strconv.FormatUint(userID, 10))

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetUser retrieves a user.
func (c *V1Client) GetUser(ctx context.Context, userID uint64) (user *v1.User, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetUser")
	defer span.End()

	req, err := c.BuildGetUserRequest(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &user)
	return user, err
}

// BuildGetUsersRequest builds an HTTP request for fetching a user.
func (c *V1Client) BuildGetUsersRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetUsersRequest")
	defer span.End()

	uri := c.buildVersionlessURL(filter.ToValues(), usersBasePath)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetUsers retrieves a list of users.
func (c *V1Client) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUsers")
	defer span.End()

	users := &v1.UserList{}

	req, err := c.BuildGetUsersRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &users)
	return users, err
}

// BuildCreateUserRequest builds an HTTP request for creating a user.
func (c *V1Client) BuildCreateUserRequest(ctx context.Context, body *v1.UserCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateUserRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath)

	return c.buildDataRequest(ctx, http.MethodPost, uri, body)
}

// CreateUser creates a new user.
func (c *V1Client) CreateUser(ctx context.Context, input *v1.UserCreationInput) (*v1.UserCreationResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateUser")
	defer span.End()

	user := &v1.UserCreationResponse{}

	req, err := c.BuildCreateUserRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeUnauthenticatedDataRequest(ctx, req, &user)
	return user, err
}

// BuildArchiveUserRequest builds an HTTP request for updating a user.
func (c *V1Client) BuildArchiveUserRequest(ctx context.Context, userID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveUserRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath, strconv.FormatUint(userID, 10))

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveUser archives a user.
func (c *V1Client) ArchiveUser(ctx context.Context, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveUser")
	defer span.End()

	req, err := c.BuildArchiveUserRequest(ctx, userID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}

// BuildLoginRequest builds an authenticating HTTP request.
func (c *V1Client) BuildLoginRequest(ctx context.Context, input *v1.UserLoginInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildLoginRequest")
	defer span.End()

	if input == nil {
		return nil, errors.New("nil input provided")
	}

	body, err := createBodyFromStruct(&input)
	if err != nil {
		return nil, fmt.Errorf("building request body: %w", err)
	}

	uri := c.buildVersionlessURL(nil, usersBasePath, "login")
	return c.buildDataRequest(ctx, http.MethodPost, uri, body)
}

// Login will, when provided the correct credentials, fetch a login cookie.
func (c *V1Client) Login(ctx context.Context, input *v1.UserLoginInput) (*http.Cookie, error) {
	ctx, span := tracing.StartSpan(ctx, "Login")
	defer span.End()

	if input == nil {
		return nil, errors.New("nil input provided")
	}

	req, err := c.BuildLoginRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error building login request: %w", err)
	}

	res, err := c.plainClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered error executing login request: %w", err)
	}
	c.closeResponseBody(res)

	cookies := res.Cookies()
	if len(cookies) > 0 {
		return cookies[0], nil
	}

	return nil, errors.New("no cookies returned from request")
}

// BuildVerifyTOTPSecretRequest builds a request to validate a TOTP secret.
func (c *V1Client) BuildVerifyTOTPSecretRequest(ctx context.Context, userID uint64, token string) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildVerifyTOTPSecretRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath, "totp_secret", "verify")

	return c.buildDataRequest(ctx, http.MethodPost, uri, &v1.TOTPSecretVerificationInput{
		TOTPToken: token,
		UserID:    userID,
	})
}

// VerifyTOTPSecret executes a request to verify a TOTP secret.
func (c *V1Client) VerifyTOTPSecret(ctx context.Context, userID uint64, token string) error {
	ctx, span := tracing.StartSpan(ctx, "VerifyTOTPSecret")
	defer span.End()

	req, err := c.BuildVerifyTOTPSecretRequest(ctx, userID, token)
	if err != nil {
		return fmt.Errorf("error building TOTP validation request: %w", err)
	}

	res, err := c.executeRawRequest(ctx, c.plainClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	c.closeResponseBody(res)

	if res.StatusCode == http.StatusBadRequest {
		return ErrInvalidTOTPToken
	} else if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("erroneous response code when validating TOTP secret: %d", res.StatusCode)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetUserRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildGetUserRequest builds an HTTP request for fetching a user.
func (c *V1Client) BuildGetUserRequest(ctx context.Context, userID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetUserRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath, strconv.FormatUint(userID, 10))

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetUser(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUser retrieves a user.
func (c *V1Client) GetUser(ctx context.Context, userID uint64) (user *v1.User, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetUser")
	defer span.End()

	req, err := c.BuildGetUserRequest(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &user)
	return user, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUsersRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetUsersRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildGetUsersRequest builds an HTTP request for fetching a user.
func (c *V1Client) BuildGetUsersRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetUsersRequest")
	defer span.End()

	uri := c.buildVersionlessURL(filter.ToValues(), usersBasePath)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetUsers(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUsers retrieves a list of users.
func (c *V1Client) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUsers")
	defer span.End()

	users := &v1.UserList{}

	req, err := c.BuildGetUsersRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.retrieve(ctx, req, &users)
	return users, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateUserRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildCreateUserRequest builds an HTTP request for creating a user.
func (c *V1Client) BuildCreateUserRequest(ctx context.Context, body *v1.UserCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateUserRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath)

	return c.buildDataRequest(ctx, http.MethodPost, uri, body)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateUser(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateUser creates a new user.
func (c *V1Client) CreateUser(ctx context.Context, input *v1.UserCreationInput) (*v1.UserCreationResponse, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateUser")
	defer span.End()

	user := &v1.UserCreationResponse{}

	req, err := c.BuildCreateUserRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeUnauthenticatedDataRequest(ctx, req, &user)
	return user, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildArchiveUserRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildArchiveUserRequest builds an HTTP request for updating a user.
func (c *V1Client) BuildArchiveUserRequest(ctx context.Context, userID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveUserRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath, strconv.FormatUint(userID, 10))

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildArchiveUser(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveUser archives a user.
func (c *V1Client) ArchiveUser(ctx context.Context, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveUser")
	defer span.End()

	req, err := c.BuildArchiveUserRequest(ctx, userID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildLoginRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildLoginRequest(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildLoginRequest builds an authenticating HTTP request.
func (c *V1Client) BuildLoginRequest(ctx context.Context, input *v1.UserLoginInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildLoginRequest")
	defer span.End()

	if input == nil {
		return nil, errors.New("nil input provided")
	}

	body, err := createBodyFromStruct(&input)
	if err != nil {
		return nil, fmt.Errorf("building request body: %w", err)
	}

	uri := c.buildVersionlessURL(nil, usersBasePath, "login")
	return c.buildDataRequest(ctx, http.MethodPost, uri, body)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildLogin(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// Login will, when provided the correct credentials, fetch a login cookie.
func (c *V1Client) Login(ctx context.Context, input *v1.UserLoginInput) (*http.Cookie, error) {
	ctx, span := tracing.StartSpan(ctx, "Login")
	defer span.End()

	if input == nil {
		return nil, errors.New("nil input provided")
	}

	req, err := c.BuildLoginRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error building login request: %w", err)
	}

	res, err := c.plainClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered error executing login request: %w", err)
	}
	c.closeResponseBody(res)

	cookies := res.Cookies()
	if len(cookies) > 0 {
		return cookies[0], nil
	}

	return nil, errors.New("no cookies returned from request")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildVerifyTOTPSecretRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildVerifyTOTPSecretRequest(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildVerifyTOTPSecretRequest builds a request to validate a TOTP secret.
func (c *V1Client) BuildVerifyTOTPSecretRequest(ctx context.Context, userID uint64, token string) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildVerifyTOTPSecretRequest")
	defer span.End()

	uri := c.buildVersionlessURL(nil, usersBasePath, "totp_secret", "verify")

	return c.buildDataRequest(ctx, http.MethodPost, uri, &v1.TOTPSecretVerificationInput{
		TOTPToken: token,
		UserID:    userID,
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVerifyTOTPSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildVerifyTOTPSecret(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
)

// VerifyTOTPSecret executes a request to verify a TOTP secret.
func (c *V1Client) VerifyTOTPSecret(ctx context.Context, userID uint64, token string) error {
	ctx, span := tracing.StartSpan(ctx, "VerifyTOTPSecret")
	defer span.End()

	req, err := c.BuildVerifyTOTPSecretRequest(ctx, userID, token)
	if err != nil {
		return fmt.Errorf("error building TOTP validation request: %w", err)
	}

	res, err := c.executeRawRequest(ctx, c.plainClient, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	c.closeResponseBody(res)

	if res.StatusCode == http.StatusBadRequest {
		return ErrInvalidTOTPToken
	} else if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("erroneous response code when validating TOTP secret: %d", res.StatusCode)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
