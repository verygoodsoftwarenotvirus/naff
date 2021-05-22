package querier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := webhooksTestDotGo(proj)

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

func TestClient_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetWebhook", mock.Anything, exampleWebhook.ID, exampleWebhook.BelongsToUser).Return(exampleWebhook, nil)

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, exampleWebhook.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllWebhooksCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expected := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetAllWebhooksCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllWebhooksCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(exampleWebhookList, nil)

		actual, err := c.GetAllWebhooks(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetWebhooks(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()
		filter := v1.DefaultQueryFilter()

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetWebhooks", mock.Anything, exampleUser.ID, filter).Return(exampleWebhookList, nil)

		actual, err := c.GetWebhooks(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()
		filter := (*v1.QueryFilter)(nil)

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetWebhooks", mock.Anything, exampleUser.ID, filter).Return(exampleWebhookList, nil)

		actual, err := c.GetWebhooks(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("CreateWebhook", mock.Anything, exampleInput).Return(exampleWebhook, nil)

		actual, err := c.CreateWebhook(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("UpdateWebhook", mock.Anything, exampleWebhook).Return(expected)

		actual := c.UpdateWebhook(ctx, exampleWebhook)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("ArchiveWebhook", mock.Anything, exampleWebhook.ID, exampleWebhook.BelongsToUser).Return(expected)

		actual := c.ArchiveWebhook(ctx, exampleWebhook.ID, exampleWebhook.BelongsToUser)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetWebhook(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetWebhook", mock.Anything, exampleWebhook.ID, exampleWebhook.BelongsToUser).Return(exampleWebhook, nil)

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, exampleWebhook.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetAllWebhooksCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetAllWebhooksCount(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

func TestClient_GetAllWebhooksCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expected := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetAllWebhooksCount", mock.Anything).Return(expected, nil)

		actual, err := c.GetAllWebhooksCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetAllWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetAllWebhooks(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_GetAllWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(exampleWebhookList, nil)

		actual, err := c.GetAllWebhooks(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetWebhooks(proj)

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

func TestClient_GetWebhooks(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()
		filter := v1.DefaultQueryFilter()

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetWebhooks", mock.Anything, exampleUser.ID, filter).Return(exampleWebhookList, nil)

		actual, err := c.GetWebhooks(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fake.BuildFakeWebhookList()
		filter := (*v1.QueryFilter)(nil)

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("GetWebhooks", mock.Anything, exampleUser.ID, filter).Return(exampleWebhookList, nil)

		actual, err := c.GetWebhooks(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_CreateWebhook(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		exampleInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("CreateWebhook", mock.Anything, exampleInput).Return(exampleWebhook, nil)

		actual, err := c.CreateWebhook(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_UpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_UpdateWebhook(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_UpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("UpdateWebhook", mock.Anything, exampleWebhook).Return(expected)

		actual := c.UpdateWebhook(ctx, exampleWebhook)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_ArchiveWebhook(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestClient_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhook := fake.BuildFakeWebhook()
		var expected error

		c, mockDB := buildTestClient()
		mockDB.WebhookDataManager.On("ArchiveWebhook", mock.Anything, exampleWebhook.ID, exampleWebhook.BelongsToUser).Return(expected)

		actual := c.ArchiveWebhook(ctx, exampleWebhook.ID, exampleWebhook.BelongsToUser)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
