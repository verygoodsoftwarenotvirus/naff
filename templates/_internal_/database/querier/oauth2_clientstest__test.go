package querier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsTestDotGo(proj)

		expected := `
package example

import (
	"context"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(exampleOAuth2Client, nil)

		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := (*v1.OAuth2Client)(nil)

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(expected, errors.New("blah"))

		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientByClientID", mock.Anything, exampleOAuth2Client.ClientID).Return(exampleOAuth2Client, nil)

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientByClientID", mock.Anything, exampleOAuth2Client.ClientID).Return(exampleOAuth2Client, errors.New("blah"))

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetAllOAuth2ClientCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()
		filter := v1.DefaultQueryFilter()

		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientsForUser", mock.Anything, exampleUser.ID, filter).Return(exampleOAuth2ClientList, nil)

		actual, err := c.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()
		filter := (*v1.QueryFilter)(nil)

		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientsForUser", mock.Anything, exampleUser.ID, filter).Return(exampleOAuth2ClientList, nil)

		actual, err := c.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		exampleOAuth2ClientList := (*v1.OAuth2ClientList)(nil)
		filter := v1.DefaultQueryFilter()

		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientsForUser", mock.Anything, exampleUser.ID, filter).Return(exampleOAuth2ClientList, errors.New("blah"))

		actual, err := c.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		mockDB.OAuth2ClientDataManager.On("CreateOAuth2Client", mock.Anything, exampleInput).Return(exampleOAuth2Client, nil)

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()

		expected := (*v1.OAuth2Client)(nil)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		mockDB.OAuth2ClientDataManager.On("CreateOAuth2Client", mock.Anything, exampleInput).Return(expected, errors.New("blah"))

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		var expected error
		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("UpdateOAuth2Client", mock.Anything, exampleOAuth2Client).Return(expected)

		actual := c.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}

func TestClient_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		var expected error
		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("ArchiveOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(expected)

		actual := c.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expected := fmt.Errorf("blah")
		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("ArchiveOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(expected)

		actual := c.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
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

func Test_buildTestClient_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(exampleOAuth2Client, nil)

		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := (*v1.OAuth2Client)(nil)

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(expected, errors.New("blah"))

		actual, err := c.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetOAuth2ClientByClientID(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientByClientID", mock.Anything, exampleOAuth2Client.ClientID).Return(exampleOAuth2Client, nil)

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientByClientID", mock.Anything, exampleOAuth2Client.ClientID).Return(exampleOAuth2Client, errors.New("blah"))

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetOAuth2ClientCount(proj)

		expected := `
package example

import (
	"context"
	"errors"
	v5 "github.com/brianvoe/gofakeit/v5"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"testing"
)

func TestClient_GetOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		filter := v1.DefaultQueryFilter()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientCount", mock.Anything, exampleUserID, filter).Return(expected, nil)

		actual, err := c.GetOAuth2ClientCount(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		exampleUserID := v5.Uint64()
		expected := v5.Uint64()
		c, mockDB := buildTestClient()
		filter := (*v1.QueryFilter)(nil)
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientCount", mock.Anything, exampleUserID, filter).Return(expected, nil)

		actual, err := c.GetOAuth2ClientCount(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleUserID := v5.Uint64()
		expected := v5.Uint64()
		c, mockDB := buildTestClient()
		filter := v1.DefaultQueryFilter()
		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientCount", mock.Anything, exampleUserID, filter).Return(expected, errors.New("blah"))

		actual, err := c.GetOAuth2ClientCount(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetAllOAuth2ClientCount(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"testing"
)

func TestClient_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleCount := uint64(123)

		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("GetAllOAuth2ClientCount", mock.Anything).Return(exampleCount, nil)

		actual, err := c.GetAllOAuth2ClientCount(ctx)
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

func Test_buildTestClient_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_GetOAuth2ClientsForUser(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()
		filter := v1.DefaultQueryFilter()

		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientsForUser", mock.Anything, exampleUser.ID, filter).Return(exampleOAuth2ClientList, nil)

		actual, err := c.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with nil filter", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()
		filter := (*v1.QueryFilter)(nil)

		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientsForUser", mock.Anything, exampleUser.ID, filter).Return(exampleOAuth2ClientList, nil)

		actual, err := c.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()
		exampleOAuth2ClientList := (*v1.OAuth2ClientList)(nil)
		filter := v1.DefaultQueryFilter()

		mockDB.OAuth2ClientDataManager.On("GetOAuth2ClientsForUser", mock.Anything, exampleUser.ID, filter).Return(exampleOAuth2ClientList, errors.New("blah"))

		actual, err := c.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_CreateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"errors"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		mockDB.OAuth2ClientDataManager.On("CreateOAuth2Client", mock.Anything, exampleInput).Return(exampleOAuth2Client, nil)

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		c, mockDB := buildTestClient()

		expected := (*v1.OAuth2Client)(nil)
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		exampleInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		mockDB.OAuth2ClientDataManager.On("CreateOAuth2Client", mock.Anything, exampleInput).Return(expected, errors.New("blah"))

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.Error(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestClient_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_UpdateOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		var expected error
		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("UpdateOAuth2Client", mock.Anything, exampleOAuth2Client).Return(expected)

		actual := c.UpdateOAuth2Client(ctx, exampleOAuth2Client)
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

func Test_buildTestClient_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestClient_ArchiveOAuth2Client(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types/fake"
	"testing"
)

func TestClient_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		var expected error
		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("ArchiveOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(expected)

		actual := c.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, actual)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error returned from querier", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expected := fmt.Errorf("blah")
		c, mockDB := buildTestClient()
		mockDB.OAuth2ClientDataManager.On("ArchiveOAuth2Client", mock.Anything, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser).Return(expected)

		actual := c.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
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
