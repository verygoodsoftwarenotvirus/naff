package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockOauth2ClientDataManagerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mockOauth2ClientDataManagerDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

var _ v1.OAuth2ClientDataManager = (*OAuth2ClientDataManager)(nil)

// OAuth2ClientDataManager is a mocked models.OAuth2ClientDataManager for testing
type OAuth2ClientDataManager struct {
	mock.Mock
}

// GetOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, clientID, userID)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}

// GetOAuth2ClientByClientID is a mock function.
func (m *OAuth2ClientDataManager) GetOAuth2ClientByClientID(ctx context.Context, identifier string) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, identifier)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}

// GetAllOAuth2ClientCount is a mock function.
func (m *OAuth2ClientDataManager) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllOAuth2Clients is a mock function.
func (m *OAuth2ClientDataManager) GetAllOAuth2Clients(ctx context.Context) ([]*v1.OAuth2Client, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*v1.OAuth2Client), args.Error(1)
}

// GetOAuth2ClientsForUser is a mock function.
func (m *OAuth2ClientDataManager) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*v1.OAuth2ClientList), args.Error(1)
}

// CreateOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}

// UpdateOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) UpdateOAuth2Client(ctx context.Context, updated *v1.OAuth2Client) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	return m.Called(ctx, clientID, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildOAuth2ClientDataManager()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// OAuth2ClientDataManager is a mocked models.OAuth2ClientDataManager for testing
type OAuth2ClientDataManager struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, clientID, userID)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientByClientID(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetOAuth2ClientByClientID is a mock function.
func (m *OAuth2ClientDataManager) GetOAuth2ClientByClientID(ctx context.Context, identifier string) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, identifier)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildGetAllOAuth2ClientCount()

		expected := `
package example

import (
	"context"
)

// GetAllOAuth2ClientCount is a mock function.
func (m *OAuth2ClientDataManager) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2Clients(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetAllOAuth2Clients is a mock function.
func (m *OAuth2ClientDataManager) GetAllOAuth2Clients(ctx context.Context) ([]*v1.OAuth2Client, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*v1.OAuth2Client), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientsForUser(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetOAuth2ClientsForUser is a mock function.
func (m *OAuth2ClientDataManager) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*v1.OAuth2ClientList), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildCreateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.OAuth2Client), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUpdateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// UpdateOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) UpdateOAuth2Client(ctx context.Context, updated *v1.OAuth2Client) error {
	return m.Called(ctx, updated).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildArchiveOAuth2Client()

		expected := `
package example

import (
	"context"
)

// ArchiveOAuth2Client is a mock function.
func (m *OAuth2ClientDataManager) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	return m.Called(ctx, clientID, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
