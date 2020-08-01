package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockUserDataManagerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mockUserDataManagerDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var _ v1.UserDataManager = (*UserDataManager)(nil)

// UserDataManager is a mocked models.UserDataManager for testing
type UserDataManager struct {
	mock.Mock
}

// GetUser is a mock function.
func (m *UserDataManager) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*v1.User), args.Error(1)
}

// GetUserWithUnverifiedTwoFactorSecret is a mock function.
func (m *UserDataManager) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*v1.User), args.Error(1)
}

// VerifyUserTwoFactorSecret is a mock function.
func (m *UserDataManager) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// GetUserByUsername is a mock function.
func (m *UserDataManager) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*v1.User), args.Error(1)
}

// GetAllUsersCount is a mock function.
func (m *UserDataManager) GetAllUsersCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetUsers is a mock function.
func (m *UserDataManager) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*v1.UserList), args.Error(1)
}

// CreateUser is a mock function.
func (m *UserDataManager) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.User), args.Error(1)
}

// UpdateUser is a mock function.
func (m *UserDataManager) UpdateUser(ctx context.Context, updated *v1.User) error {
	return m.Called(ctx, updated).Error(0)
}

// UpdateUserPassword is a mock function.
func (m *UserDataManager) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	return m.Called(ctx, userID, newHash).Error(0)
}

// ArchiveUser is a mock function.
func (m *UserDataManager) ArchiveUser(ctx context.Context, userID uint64) error {
	return m.Called(ctx, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserDataManager()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// UserDataManager is a mocked models.UserDataManager for testing
type UserDataManager struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildGetUser(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUser is a mock function.
func (m *UserDataManager) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*v1.User), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildGetUserWithUnverifiedTwoFactorSecret(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserWithUnverifiedTwoFactorSecret is a mock function.
func (m *UserDataManager) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*v1.User), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVerifyUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildVerifyUserTwoFactorSecret()

		expected := `
package example

import (
	"context"
)

// VerifyUserTwoFactorSecret is a mock function.
func (m *UserDataManager) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildGetUserByUsername(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserByUsername is a mock function.
func (m *UserDataManager) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*v1.User), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllUsersCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildGetAllUsersCount()

		expected := `
package example

import (
	"context"
)

// GetAllUsersCount is a mock function.
func (m *UserDataManager) GetAllUsersCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildGetUsers(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUsers is a mock function.
func (m *UserDataManager) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*v1.UserList), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildCreateUser(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateUser is a mock function.
func (m *UserDataManager) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.User), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildUpdateUser(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateUser is a mock function.
func (m *UserDataManager) UpdateUser(ctx context.Context, updated *v1.User) error {
	return m.Called(ctx, updated).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUpdateUserPassword()

		expected := `
package example

import (
	"context"
)

// UpdateUserPassword is a mock function.
func (m *UserDataManager) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	return m.Called(ctx, userID, newHash).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildArchiveUser()

		expected := `
package example

import (
	"context"
)

// ArchiveUser is a mock function.
func (m *UserDataManager) ArchiveUser(ctx context.Context, userID uint64) error {
	return m.Called(ctx, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
