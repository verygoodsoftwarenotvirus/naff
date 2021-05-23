package querier

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
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

var (
	_ v1.UserDataManager = (*Client)(nil)

	// ErrUserExists is a sentinel error for returning when a username is taken.
	ErrUserExists = errors.New("error: username already exists")
)

// GetUser fetches a user.
func (c *Client) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUser")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetUser called")

	return c.querier.GetUser(ctx, userID)
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret.
func (c *Client) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUserWithUnverifiedTwoFactorSecret")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetUserWithUnverifiedTwoFactorSecret called")

	return c.querier.GetUserWithUnverifiedTwoFactorSecret(ctx, userID)
}

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (c *Client) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "VerifyUserTwoFactorSecret")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("VerifyUserTwoFactorSecret called")

	return c.querier.VerifyUserTwoFactorSecret(ctx, userID)
}

// GetUserByUsername fetches a user by their username.
func (c *Client) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUserByUsername")
	defer span.End()

	tracing.AttachUsernameToSpan(span, username)
	c.logger.WithValue("username", username).Debug("GetUserByUsername called")

	return c.querier.GetUserByUsername(ctx, username)
}

// GetAllUsersCount fetches a count of users from the database that meet a particular filter.
func (c *Client) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllUsersCount")
	defer span.End()

	c.logger.Debug("GetAllUsersCount called")

	return c.querier.GetAllUsersCount(ctx)
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (c *Client) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUsers")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)
	c.logger.WithValue("filter", filter).Debug("GetUsers called")

	return c.querier.GetUsers(ctx, filter)
}

// CreateUser creates a user.
func (c *Client) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateUser")
	defer span.End()

	tracing.AttachUsernameToSpan(span, input.Username)
	c.logger.WithValue("username", input.Username).Debug("CreateUser called")

	return c.querier.CreateUser(ctx, input)
}

// UpdateUser receives a complete User struct and updates its record in the database.
// NOTE: this function uses the ID provided in the input to make its query.
func (c *Client) UpdateUser(ctx context.Context, updated *v1.User) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateUser")
	defer span.End()

	tracing.AttachUsernameToSpan(span, updated.Username)
	c.logger.WithValue("username", updated.Username).Debug("UpdateUser called")

	return c.querier.UpdateUser(ctx, updated)
}

// UpdateUserPassword updates a user's password hash in the database.
func (c *Client) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateUserPassword")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("UpdateUserPassword called")

	return c.querier.UpdateUserPassword(ctx, userID, newHash)
}

// ArchiveUser archives a user.
func (c *Client) ArchiveUser(ctx context.Context, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveUser")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("ArchiveUser called")

	return c.querier.ArchiveUser(ctx, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildVarDeclarations(proj)

		expected := `
package example

import (
	"errors"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

var (
	_ v1.UserDataManager = (*Client)(nil)

	// ErrUserExists is a sentinel error for returning when a username is taken.
	ErrUserExists = errors.New("error: username already exists")
)
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
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetUser fetches a user.
func (c *Client) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUser")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetUser called")

	return c.querier.GetUser(ctx, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetUserWithUnverifiedTwoFactorSecret(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret.
func (c *Client) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUserWithUnverifiedTwoFactorSecret")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetUserWithUnverifiedTwoFactorSecret called")

	return c.querier.GetUserWithUnverifiedTwoFactorSecret(ctx, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVerifyUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildVerifyUserTwoFactorSecret(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (c *Client) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "VerifyUserTwoFactorSecret")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("VerifyUserTwoFactorSecret called")

	return c.querier.VerifyUserTwoFactorSecret(ctx, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetUserByUsername(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetUserByUsername fetches a user by their username.
func (c *Client) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUserByUsername")
	defer span.End()

	tracing.AttachUsernameToSpan(span, username)
	c.logger.WithValue("username", username).Debug("GetUserByUsername called")

	return c.querier.GetUserByUsername(ctx, username)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllUsersCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetAllUsersCount(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// GetAllUsersCount fetches a count of users from the database that meet a particular filter.
func (c *Client) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllUsersCount")
	defer span.End()

	c.logger.Debug("GetAllUsersCount called")

	return c.querier.GetAllUsersCount(ctx)
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
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetUsers fetches a list of users from the database that meet a particular filter.
func (c *Client) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUsers")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)
	c.logger.WithValue("filter", filter).Debug("GetUsers called")

	return c.querier.GetUsers(ctx, filter)
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
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateUser creates a user.
func (c *Client) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateUser")
	defer span.End()

	tracing.AttachUsernameToSpan(span, input.Username)
	c.logger.WithValue("username", input.Username).Debug("CreateUser called")

	return c.querier.CreateUser(ctx, input)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUpdateUser(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// UpdateUser receives a complete User struct and updates its record in the database.
// NOTE: this function uses the ID provided in the input to make its query.
func (c *Client) UpdateUser(ctx context.Context, updated *v1.User) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateUser")
	defer span.End()

	tracing.AttachUsernameToSpan(span, updated.Username)
	c.logger.WithValue("username", updated.Username).Debug("UpdateUser called")

	return c.querier.UpdateUser(ctx, updated)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUpdateUserPassword(proj)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// UpdateUserPassword updates a user's password hash in the database.
func (c *Client) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateUserPassword")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("UpdateUserPassword called")

	return c.querier.UpdateUserPassword(ctx, userID, newHash)
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
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveUser archives a user.
func (c *Client) ArchiveUser(ctx context.Context, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveUser")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("ArchiveUser called")

	return c.querier.ArchiveUser(ctx, userID)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
