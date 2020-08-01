package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockWebhookDataManagerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mockWebhookDataManagerDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var _ v1.WebhookDataManager = (*WebhookDataManager)(nil)

// WebhookDataManager is a mocked models.WebhookDataManager for testing
type WebhookDataManager struct {
	mock.Mock
}

// GetWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	args := m.Called(ctx, webhookID, userID)
	return args.Get(0).(*v1.Webhook), args.Error(1)
}

// GetAllWebhooksCount satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetAllWebhooksCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetWebhooks satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*v1.WebhookList), args.Error(1)
}

// GetAllWebhooks satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	args := m.Called(ctx)
	return args.Get(0).(*v1.WebhookList), args.Error(1)
}

// CreateWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.Webhook), args.Error(1)
}

// UpdateWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) UpdateWebhook(ctx context.Context, updated *v1.Webhook) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	return m.Called(ctx, webhookID, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockWebhookDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockWebhookDataManager()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// WebhookDataManager is a mocked models.WebhookDataManager for testing
type WebhookDataManager struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockGetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMockGetWebhook(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	args := m.Called(ctx, webhookID, userID)
	return args.Get(0).(*v1.Webhook), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockGetAllWebhooksCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockGetAllWebhooksCount()

		expected := `
package example

import (
	"context"
)

// GetAllWebhooksCount satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetAllWebhooksCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockGetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMockGetWebhooks(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhooks satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*v1.WebhookList), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockGetAllWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMockGetAllWebhooks(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllWebhooks satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	args := m.Called(ctx)
	return args.Get(0).(*v1.WebhookList), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockCreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMockCreateWebhook(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*v1.Webhook), args.Error(1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockUpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMockUpdateWebhook(proj)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) UpdateWebhook(ctx context.Context, updated *v1.Webhook) error {
	return m.Called(ctx, updated).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMockArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockArchiveWebhook()

		expected := `
package example

import (
	"context"
)

// ArchiveWebhook satisfies our WebhookDataManager interface.
func (m *WebhookDataManager) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	return m.Called(ctx, webhookID, userID).Error(0)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
