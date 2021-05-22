package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_databaseMockDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := databaseMockDotGo(proj)

		expected := `
package example

import (
	"context"
	mock1 "github.com/stretchr/testify/mock"
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
)

var _ DataManager = (*MockDatabase)(nil)

// BuildMockDatabase builds a mock database.
func BuildMockDatabase() *MockDatabase {
	return &MockDatabase{
		ItemDataManager:         &mock.ItemDataManager{},
		UserDataManager:         &mock.UserDataManager{},
		OAuth2ClientDataManager: &mock.OAuth2ClientDataManager{},
		WebhookDataManager:      &mock.WebhookDataManager{},
	}
}

// MockDatabase is our mock database structure.
type MockDatabase struct {
	mock1.Mock

	*mock.ItemDataManager
	*mock.UserDataManager
	*mock.OAuth2ClientDataManager
	*mock.WebhookDataManager
}

// Migrate satisfies the DataManager interface.
func (m *MockDatabase) Migrate(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// IsReady satisfies the DataManager interface.
func (m *MockDatabase) IsReady(ctx context.Context) (ready bool) {
	return m.Called(ctx).Bool(0)
}

var _ ResultIterator = (*MockResultIterator)(nil)

// MockResultIterator is our mock sql.Rows structure.
type MockResultIterator struct {
	mock1.Mock
}

// Scan satisfies the ResultIterator interface.
func (m *MockResultIterator) Scan(dest ...interface{}) error {
	return m.Called(dest...).Error(0)
}

// Next satisfies the ResultIterator interface.
func (m *MockResultIterator) Next() bool {
	return m.Called().Bool(0)
}

// Err satisfies the ResultIterator interface.
func (m *MockResultIterator) Err() error {
	return m.Called().Error(0)
}

// Close satisfies the ResultIterator interface.
func (m *MockResultIterator) Close() error {
	return m.Called().Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildMockDatabase(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildBuildMockDatabase(proj)

		expected := `
package example

import (
	mock "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
)

// BuildMockDatabase builds a mock database.
func BuildMockDatabase() *MockDatabase {
	return &MockDatabase{
		ItemDataManager:         &mock.ItemDataManager{},
		UserDataManager:         &mock.UserDataManager{},
		OAuth2ClientDataManager: &mock.OAuth2ClientDataManager{},
		WebhookDataManager:      &mock.WebhookDataManager{},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockDatabase(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildMockDatabase(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	mock1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/mock"
)

// MockDatabase is our mock database structure.
type MockDatabase struct {
	mock.Mock

	*mock1.ItemDataManager
	*mock1.UserDataManager
	*mock1.OAuth2ClientDataManager
	*mock1.WebhookDataManager
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMigrate(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMigrate()

		expected := `
package example

import (
	"context"
)

// Migrate satisfies the DataManager interface.
func (m *MockDatabase) Migrate(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIsReady(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildIsReady()

		expected := `
package example

import (
	"context"
)

// IsReady satisfies the DataManager interface.
func (m *MockDatabase) IsReady(ctx context.Context) (ready bool) {
	return m.Called(ctx).Bool(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildResultIterator(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildResultIterator()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

var _ ResultIterator = (*MockResultIterator)(nil)

// MockResultIterator is our mock sql.Rows structure.
type MockResultIterator struct {
	mock.Mock
}

// Scan satisfies the ResultIterator interface.
func (m *MockResultIterator) Scan(dest ...interface{}) error {
	return m.Called(dest...).Error(0)
}

// Next satisfies the ResultIterator interface.
func (m *MockResultIterator) Next() bool {
	return m.Called().Bool(0)
}

// Err satisfies the ResultIterator interface.
func (m *MockResultIterator) Err() error {
	return m.Called().Error(0)
}

// Close satisfies the ResultIterator interface.
func (m *MockResultIterator) Close() error {
	return m.Called().Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
