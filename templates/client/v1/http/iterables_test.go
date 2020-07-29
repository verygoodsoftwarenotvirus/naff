package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterablesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := iterablesDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"net/url"
	"strconv"
)

const (
	itemsBasePath = "items"
)

// BuildItemExistsRequest builds an HTTP request for checking the existence of an item.
func (c *V1Client) BuildItemExistsRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildItemExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(itemID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ItemExists retrieves whether or not an item exists.
func (c *V1Client) ItemExists(ctx context.Context, itemID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ItemExists")
	defer span.End()

	req, err := c.BuildItemExistsRequest(ctx, itemID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetItemRequest builds an HTTP request for fetching an item.
func (c *V1Client) BuildGetItemRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(itemID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetItem retrieves an item.
func (c *V1Client) GetItem(ctx context.Context, itemID uint64) (item *v1.Item, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetItem")
	defer span.End()

	req, err := c.BuildGetItemRequest(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &item); retrieveErr != nil {
		return nil, retrieveErr
	}

	return item, nil
}

// BuildSearchItemsRequest builds an HTTP request for querying items.
func (c *V1Client) BuildSearchItemsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildSearchItemsRequest")
	defer span.End()

	params := url.Values{}
	params.Set(v1.SearchQueryKey, query)
	params.Set(v1.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := c.BuildURL(
		params,
		itemsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// SearchItems searches for a list of items.
func (c *V1Client) SearchItems(ctx context.Context, query string, limit uint8) (items []v1.Item, err error) {
	ctx, span := tracing.StartSpan(ctx, "SearchItems")
	defer span.End()

	req, err := c.BuildSearchItemsRequest(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &items); retrieveErr != nil {
		return nil, retrieveErr
	}

	return items, nil
}

// BuildGetItemsRequest builds an HTTP request for fetching items.
func (c *V1Client) BuildGetItemsRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetItemsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		itemsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetItems retrieves a list of items.
func (c *V1Client) GetItems(ctx context.Context, filter *v1.QueryFilter) (items *v1.ItemList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetItems")
	defer span.End()

	req, err := c.BuildGetItemsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &items); retrieveErr != nil {
		return nil, retrieveErr
	}

	return items, nil
}

// BuildCreateItemRequest builds an HTTP request for creating an item.
func (c *V1Client) BuildCreateItemRequest(ctx context.Context, input *v1.ItemCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateItem creates an item.
func (c *V1Client) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (item *v1.Item, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateItem")
	defer span.End()

	req, err := c.BuildCreateItemRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &item)
	return item, err
}

// BuildUpdateItemRequest builds an HTTP request for updating an item.
func (c *V1Client) BuildUpdateItemRequest(ctx context.Context, item *v1.Item) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(item.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, item)
}

// UpdateItem updates an item.
func (c *V1Client) UpdateItem(ctx context.Context, item *v1.Item) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateItem")
	defer span.End()

	req, err := c.BuildUpdateItemRequest(ctx, item)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &item)
}

// BuildArchiveItemRequest builds an HTTP request for updating an item.
func (c *V1Client) BuildArchiveItemRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(itemID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveItem archives an item.
func (c *V1Client) ArchiveItem(ctx context.Context, itemID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveItem")
	defer span.End()

	req, err := c.BuildArchiveItemRequest(ctx, itemID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_attachURIToSpanCall(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := attachURIToSpanCall(proj)

		expected := `
package main

import (
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

func main() {
	tracing.AttachRequestURIToSpan(span, uri)
}
`
		actual := testutils.RenderIndependentStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj, typ)

		expected := `
package main

import (
	"strconv"
)

func main() {
	exampleFunction(nil, itemsBasePath, strconv.FormatUint(itemID, 10))
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForCreatingSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildV1ClientURLBuildingParamsForCreatingSomething(proj, typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(nil, itemsBasePath)
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForListOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildV1ClientURLBuildingParamsForListOfSomething(proj, typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(filter.ToValues(), itemsBasePath)
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForSearchingSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildV1ClientURLBuildingParamsForSearchingSomething(typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(params, itemsBasePath, "search")
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj, typ)

		expected := `
package main

import (
	"strconv"
)

func main() {
	exampleFunction(nil, itemsBasePath, strconv.FormatUint(item.ID, 10))
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildSomethingExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildSomethingExistsRequest(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildItemExistsRequest builds an HTTP request for checking the existence of an item.
func (c *V1Client) BuildItemExistsRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildItemExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(itemID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildSomethingExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildSomethingExists(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ItemExists retrieves whether or not an item exists.
func (c *V1Client) ItemExists(ctx context.Context, itemID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ItemExists")
	defer span.End()

	req, err := c.BuildItemExistsRequest(ctx, itemID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildGetSomethingRequestFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildGetItemRequest builds an HTTP request for fetching an item.
func (c *V1Client) BuildGetItemRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(itemID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildGetSomethingFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetItem retrieves an item.
func (c *V1Client) GetItem(ctx context.Context, itemID uint64) (item *v1.Item, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetItem")
	defer span.End()

	req, err := c.BuildGetItemRequest(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &item); retrieveErr != nil {
		return nil, retrieveErr
	}

	return item, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildSearchSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildSearchSomethingRequestFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"net/url"
	"strconv"
)

// BuildSearchItemsRequest builds an HTTP request for querying items.
func (c *V1Client) BuildSearchItemsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildSearchItemsRequest")
	defer span.End()

	params := url.Values{}
	params.Set(v1.SearchQueryKey, query)
	params.Set(v1.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := c.BuildURL(
		params,
		itemsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildSearchSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildSearchSomethingFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// SearchItems searches for a list of items.
func (c *V1Client) SearchItems(ctx context.Context, query string, limit uint8) (items []v1.Item, err error) {
	ctx, span := tracing.StartSpan(ctx, "SearchItems")
	defer span.End()

	req, err := c.BuildSearchItemsRequest(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &items); retrieveErr != nil {
		return nil, retrieveErr
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetListOfSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildGetListOfSomethingRequestFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildGetItemsRequest builds an HTTP request for fetching items.
func (c *V1Client) BuildGetItemsRequest(ctx context.Context, filter *v1.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetItemsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		itemsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetItems retrieves a list of items.
func (c *V1Client) GetItems(ctx context.Context, filter *v1.QueryFilter) (items *v1.ItemList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetItems")
	defer span.End()

	req, err := c.BuildGetItemsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &items); retrieveErr != nil {
		return nil, retrieveErr
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildCreateSomethingRequestFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

// BuildCreateItemRequest builds an HTTP request for creating an item.
func (c *V1Client) BuildCreateItemRequest(ctx context.Context, input *v1.ItemCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildCreateSomethingFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateItem creates an item.
func (c *V1Client) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (item *v1.Item, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateItem")
	defer span.End()

	req, err := c.BuildCreateItemRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &item)
	return item, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildUpdateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildUpdateSomethingRequestFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
	"strconv"
)

// BuildUpdateItemRequest builds an HTTP request for updating an item.
func (c *V1Client) BuildUpdateItemRequest(ctx context.Context, item *v1.Item) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(item.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, item)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateItem updates an item.
func (c *V1Client) UpdateItem(ctx context.Context, item *v1.Item) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateItem")
	defer span.End()

	req, err := c.BuildUpdateItemRequest(ctx, item)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &item)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildBuildArchiveSomethingRequestFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	"net/http"
	"strconv"
)

// BuildArchiveItemRequest builds an HTTP request for updating an item.
func (c *V1Client) BuildArchiveItemRequest(ctx context.Context, itemID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveItemRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		itemsBasePath,
		strconv.FormatUint(itemID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingFuncDecl(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

// ArchiveItem archives an item.
func (c *V1Client) ArchiveItem(ctx context.Context, itemID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveItem")
	defer span.End()

	req, err := c.BuildArchiveItemRequest(ctx, itemID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
