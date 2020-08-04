package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterablesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_BuildItemExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleItem := fake.BuildFakeItem()
		actual, err := c.BuildItemExistsRequest(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleItem.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ItemExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleItem.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.ItemExists(ctx, exampleItem.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ItemExists(ctx, exampleItem.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleItem := fake.BuildFakeItem()

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetItemRequest(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleItem.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleItem.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleItem))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItem(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItem, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetItem(ctx, exampleItem.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleItem.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItem(ctx, exampleItem.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetItemsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetItemsRequest(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetItems(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/items"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)

		exampleItemList := fake.BuildFakeItemList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleItemList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItems(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItemList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetItems(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItems(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildSearchItemsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildSearchItemsRequest(ctx, exampleQuery, limit)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_SearchItems(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/items/search"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		exampleItemList := fake.BuildFakeItemList().Items

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.URL.Query().Get(v1.SearchQueryKey), exampleQuery, "expected and actual search query param do not match")
					assert.Equal(t, req.URL.Query().Get(v1.LimitQueryKey), strconv.FormatUint(uint64(limit), 10), "expected and actual limit query param do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleItemList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.SearchItems(ctx, exampleQuery, limit)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItemList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchItems(ctx, exampleQuery, limit)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.URL.Query().Get(v1.SearchQueryKey), exampleQuery, "expected and actual search query param do not match")
					assert.Equal(t, req.URL.Query().Get(v1.LimitQueryKey), strconv.FormatUint(uint64(limit), 10), "expected and actual limit query param do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.SearchItems(ctx, exampleQuery, limit)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateItemRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateItem(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/items"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.ItemCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToUser = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleItem))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateItem(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItem, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateItem(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateItemRequest(ctx, exampleItem)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleItem))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateItem(ctx, exampleItem)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		err := buildTestClientWithInvalidURL(t).UpdateItem(ctx, exampleItem)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleItem := fake.BuildFakeItem()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveItemRequest(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleItem.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveItem(ctx, exampleItem.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		err := buildTestClientWithInvalidURL(t).ArchiveItem(ctx, exampleItem.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildSomethingExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_BuildSomethingExistsRequest(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildItemExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleItem := fake.BuildFakeItem()
		actual, err := c.BuildItemExistsRequest(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleItem.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_SomethingExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_SomethingExists(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_ItemExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleItem.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.ItemExists(ctx, exampleItem.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ItemExists(ctx, exampleItem.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildGetSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_BuildGetSomethingRequest(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildGetItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleItem := fake.BuildFakeItem()

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetItemRequest(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleItem.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_GetSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_GetSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestV1Client_GetItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleItem.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleItem))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItem(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItem, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetItem(ctx, exampleItem.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleItem.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItem(ctx, exampleItem.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildSearchSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_BuildSearchSomethingRequest(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildSearchItemsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildSearchItemsRequest(ctx, exampleQuery, limit)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_SearchSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_SearchSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestV1Client_SearchItems(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/items/search"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		exampleItemList := fake.BuildFakeItemList().Items

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.URL.Query().Get(v1.SearchQueryKey), exampleQuery, "expected and actual search query param do not match")
					assert.Equal(t, req.URL.Query().Get(v1.LimitQueryKey), strconv.FormatUint(uint64(limit), 10), "expected and actual limit query param do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleItemList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.SearchItems(ctx, exampleQuery, limit)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItemList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchItems(ctx, exampleQuery, limit)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		limit := v1.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.URL.Query().Get(v1.SearchQueryKey), exampleQuery, "expected and actual search query param do not match")
					assert.Equal(t, req.URL.Query().Get(v1.LimitQueryKey), strconv.FormatUint(uint64(limit), 10), "expected and actual limit query param do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.SearchItems(ctx, exampleQuery, limit)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildGetListOfSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_BuildGetListOfSomethingRequest(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildGetItemsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetItemsRequest(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_GetListOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_GetListOfSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_GetItems(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/items"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)

		exampleItemList := fake.BuildFakeItemList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleItemList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItems(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItemList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetItems(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		filter := (*v1.QueryFilter)(nil)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetItems(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildTestV1Client_GetListOfSomething(proj, proj.LastDataType())

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_GetYetAnotherThings(T *testing.T) {
	T.Parallel()

	const expectedPath = fmt.Sprintf("/api/v1/things/%d/another_things/%d/yet_another_things", exampleThing.ID, exampleAnotherThing.ID)

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		filter := (*v1.QueryFilter)(nil)

		exampleYetAnotherThingList := fake.BuildFakeYetAnotherThingList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleYetAnotherThingList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetYetAnotherThings(ctx, exampleThing.ID, exampleAnotherThing.ID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleYetAnotherThingList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		filter := (*v1.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetYetAnotherThings(ctx, exampleThing.ID, exampleAnotherThing.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		filter := (*v1.QueryFilter)(nil)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetYetAnotherThings(ctx, exampleThing.ID, exampleAnotherThing.ID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildCreateSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildTestV1Client_BuildCreateSomethingRequest(proj, proj.LastDataType())

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildCreateYetAnotherThingRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = exampleThing.ID
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateYetAnotherThingRequest(ctx, exampleThing.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_CreateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_CreateSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_CreateItem(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/items"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.ItemCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToUser = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleItem))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateItem(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleItem, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateItem(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildTestV1Client_CreateSomething(proj, proj.LastDataType())

		expected := `
package example

import (
	"context"
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_CreateYetAnotherThing(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/things/%d/another_things/%d/yet_another_things"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *v1.YetAnotherThingCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToAnotherThing = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleYetAnotherThing))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateYetAnotherThing(ctx, exampleThing.ID, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleYetAnotherThing, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateYetAnotherThing(ctx, exampleThing.ID, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildUpdateSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_BuildUpdateSomethingRequest(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_BuildUpdateItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateItemRequest(ctx, exampleItem)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_UpdateSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_UpdateSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"encoding/json"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_UpdateItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleItem))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateItem(ctx, exampleItem)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		err := buildTestClientWithInvalidURL(t).UpdateItem(ctx, exampleItem)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildArchiveSomethingRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_BuildArchiveSomethingRequest(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestV1Client_BuildArchiveItemRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleItem := fake.BuildFakeItem()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveItemRequest(ctx, exampleItem.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleItem.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_ArchiveSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestV1Client_ArchiveSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1Client_ArchiveItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/items/%d", exampleItem.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveItem(ctx, exampleItem.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		err := buildTestClientWithInvalidURL(t).ArchiveItem(ctx, exampleItem.ID)
		assert.Error(t, err, "error should be returned")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
