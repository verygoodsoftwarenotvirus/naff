package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_clientTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := clientTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

const (
	defaultLimit = uint8(20)
)

func buildTestClient() (*Client, *v1.MockDatabase) {
	db := v1.BuildMockDatabase()
	c := &Client{
		logger:  noop.ProvideNoopLogger(),
		querier: db,
	}
	return c, db
}

func TestMigrate(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(nil)

		c := &Client{querier: mockDB}
		actual := c.Migrate(ctx)
		assert.NoError(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("bubbles up errors", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(errors.New("blah"))

		c := &Client{querier: mockDB}
		actual := c.Migrate(ctx)
		assert.Error(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestIsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("IsReady", mock.Anything).Return(true)

		c := &Client{querier: mockDB}
		c.IsReady(ctx)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(nil)

		actual, err := ProvideDatabaseClient(ctx, nil, mockDB, true, noop.ProvideNoopLogger())
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error migrating querier", func(t *testing.T) {
		ctx := context.Background()

		expected := errors.New("blah")
		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(expected)

		x, actual := ProvideDatabaseClient(ctx, nil, mockDB, true, noop.ProvideNoopLogger())
		assert.Nil(t, x)
		assert.Error(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildTestClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildTestClient(proj)

		expected := `
package example

import (
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

func buildTestClient() (*Client, *v1.MockDatabase) {
	db := v1.BuildMockDatabase()
	c := &Client{
		logger:  noop.ProvideNoopLogger(),
		querier: db,
	}
	return c, db
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestMigrate(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestMigrate(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestMigrate(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(nil)

		c := &Client{querier: mockDB}
		actual := c.Migrate(ctx)
		assert.NoError(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("bubbles up errors", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(errors.New("blah"))

		c := &Client{querier: mockDB}
		actual := c.Migrate(ctx)
		assert.Error(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestIsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestIsReady(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestIsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("IsReady", mock.Anything).Return(true)

		c := &Client{querier: mockDB}
		c.IsReady(ctx)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestProvideDatabaseClient(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(nil)

		actual, err := ProvideDatabaseClient(ctx, nil, mockDB, true, noop.ProvideNoopLogger())
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error migrating querier", func(t *testing.T) {
		ctx := context.Background()

		expected := errors.New("blah")
		mockDB := v1.BuildMockDatabase()
		mockDB.On("Migrate", mock.Anything).Return(expected)

		x, actual := ProvideDatabaseClient(ctx, nil, mockDB, true, noop.ProvideNoopLogger())
		assert.Nil(t, x)
		assert.Error(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
