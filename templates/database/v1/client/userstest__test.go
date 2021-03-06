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
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)

		actual, err := c.GetUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)

		actual, err := c.GetUserWithUnverifiedTwoFactorSecret(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_VerifyUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("VerifyUserTwoFactorSecret", mock.Anything, exampleUser.ID).Return(nil)

		err := c.VerifyUserTwoFactorSecret(ctx, exampleUser.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUserByUsername", mock.Anything, exampleUser.Username).Return(exampleUser, nil)

		actual, err := c.GetUserByUsername(ctx, exampleUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllUsersCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetAllUsersCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllUsersCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetUsers(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUserList := fake.BuildFakeUserList()
		filter := v1.DefaultQueryFilter()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, filter).Return(exampleUserList, nil)

		actual, err := c.GetUsers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleUserList := fake.BuildFakeUserList()
		filter := (*v1.QueryFilter)(nil)

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, filter).Return(exampleUserList, nil)

		actual, err := c.GetUsers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserDatabaseCreationInputFromUser(exampleUser)

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("CreateUser", mock.Anything, exampleInput).Return(exampleUser, nil)

		actual, err := c.CreateUser(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, exampleUser).Return(expected)

		err := c.UpdateUser(ctx, exampleUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, exampleUser.HashedPassword).Return(expected)

		err := c.UpdateUserPassword(ctx, exampleUser.ID, exampleUser.HashedPassword)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("ArchiveUser", mock.Anything, exampleUser.ID).Return(nil)

		err := c.ArchiveUser(ctx, exampleUser.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetUser(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUser", mock.Anything, exampleUser.ID).Return(exampleUser, nil)

		actual, err := c.GetUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetUserWithUnverifiedTwoFactorSecret(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUserWithUnverifiedTwoFactorSecret", mock.Anything, exampleUser.ID).Return(exampleUser, nil)

		actual, err := c.GetUserWithUnverifiedTwoFactorSecret(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_VerifyUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_VerifyUserTwoFactorSecret(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_VerifyUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("VerifyUserTwoFactorSecret", mock.Anything, exampleUser.ID).Return(nil)

		err := c.VerifyUserTwoFactorSecret(ctx, exampleUser.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetUserByUsername(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUserByUsername", mock.Anything, exampleUser.Username).Return(exampleUser, nil)

		actual, err := c.GetUserByUsername(ctx, exampleUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetAllUsersCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetAllUsersCount(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

func TestClient_GetAllUsersCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetAllUsersCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllUsersCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetUsers(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetUsers(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUserList := fake.BuildFakeUserList()
		filter := v1.DefaultQueryFilter()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, filter).Return(exampleUserList, nil)

		actual, err := c.GetUsers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleUserList := fake.BuildFakeUserList()
		filter := (*v1.QueryFilter)(nil)

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("GetUsers", mock.Anything, filter).Return(exampleUserList, nil)

		actual, err := c.GetUsers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_CreateUser(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleInput := fake.BuildFakeUserDatabaseCreationInputFromUser(exampleUser)

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("CreateUser", mock.Anything, exampleInput).Return(exampleUser, nil)

		actual, err := c.CreateUser(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_UpdateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_UpdateUser(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_UpdateUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("UpdateUser", mock.Anything, exampleUser).Return(expected)

		err := c.UpdateUser(ctx, exampleUser)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_UpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_UpdateUserPassword(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_UpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("UpdateUserPassword", mock.Anything, exampleUser.ID, exampleUser.HashedPassword).Return(expected)

		err := c.UpdateUserPassword(ctx, exampleUser.ID, exampleUser.HashedPassword)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_ArchiveUser(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		c, mockDB := buildTestClient()
		mockDB.UserDataManager.On("ArchiveUser", mock.Anything, exampleUser.ID).Return(nil)

		err := c.ArchiveUser(ctx, exampleUser.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
