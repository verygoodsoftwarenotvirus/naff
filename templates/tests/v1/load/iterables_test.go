package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterablesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesDotGo(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"math/rand"
	http1 "net/http"
)

// fetchRandomItem retrieves a random item from the list of available items.
func fetchRandomItem(ctx context.Context, c *http.V1Client) *v1.Item {
	itemsRes, err := c.GetItems(ctx, nil)
	if err != nil || itemsRes == nil || len(itemsRes.Items) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(itemsRes.Items))
	return &itemsRes.Items[randIndex]
}

func buildItemActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateItem": {
			Name: "CreateItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				itemInput := fake.BuildFakeItemCreationInput()

				return c.BuildCreateItemRequest(ctx, itemInput)
			},
			Weight: 100,
		},
		"GetItem": {
			Name: "GetItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				randomItem := fetchRandomItem(ctx, c)
				if randomItem == nil {
					return nil, fmt.Errorf("retrieving random item: %w", ErrUnavailableYet)
				}

				return c.BuildGetItemRequest(ctx, randomItem.ID)
			},
			Weight: 100,
		},
		"GetItems": {
			Name: "GetItems",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				return c.BuildGetItemsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateItem": {
			Name: "UpdateItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				if randomItem := fetchRandomItem(ctx, c); randomItem != nil {
					newItem := fake.BuildFakeItemCreationInput()
					randomItem.Name = newItem.Name
					randomItem.Details = newItem.Details
					return c.BuildUpdateItemRequest(ctx, randomItem)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveItem": {
			Name: "ArchiveItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				randomItem := fetchRandomItem(ctx, c)
				if randomItem == nil {
					return nil, fmt.Errorf("retrieving random item: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveItemRequest(ctx, randomItem.ID)
			},
			Weight: 85,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		call := true
		x := buildParamsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ, call)

		expected := `
package main

import ()

func example(ctx) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFetchRandomSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildFetchRandomSomething(proj, typ)

		expected := `
package example

import (
	"context"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"math/rand"
)

// fetchRandomItem retrieves a random item from the list of available items.
func fetchRandomItem(ctx context.Context, c *http.V1Client) *v1.Item {
	itemsRes, err := c.GetItems(ctx, nil)
	if err != nil || itemsRes == nil || len(itemsRes.Items) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(itemsRes.Items))
	return &itemsRes.Items[randIndex]
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreationArguments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		varPrefix := "example"
		x := buildCreationArguments(proj, varPrefix, typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		ctx,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCreationCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRequisiteCreationCode(proj, typ)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForMethodThatHandlesAnInstanceWithStructs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildParamsForMethodThatHandlesAnInstanceWithStructs(proj, typ)

		expected := `
package main

import ()

func example(ctx, createdItem.ID) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCallArgsForMethodThatHandlesAnInstanceWithRetrievedStructs(proj, typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		ctx,
		randomItem.ID,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandomActionMap(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRandomActionMap(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	http1 "net/http"
)

func buildItemActions(c *http.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateItem": {
			Name: "CreateItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				itemInput := fake.BuildFakeItemCreationInput()

				return c.BuildCreateItemRequest(ctx, itemInput)
			},
			Weight: 100,
		},
		"GetItem": {
			Name: "GetItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				randomItem := fetchRandomItem(ctx, c)
				if randomItem == nil {
					return nil, fmt.Errorf("retrieving random item: %w", ErrUnavailableYet)
				}

				return c.BuildGetItemRequest(ctx, randomItem.ID)
			},
			Weight: 100,
		},
		"GetItems": {
			Name: "GetItems",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				return c.BuildGetItemsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateItem": {
			Name: "UpdateItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				if randomItem := fetchRandomItem(ctx, c); randomItem != nil {
					newItem := fake.BuildFakeItemCreationInput()
					randomItem.Name = newItem.Name
					randomItem.Details = newItem.Details
					return c.BuildUpdateItemRequest(ctx, randomItem)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveItem": {
			Name: "ArchiveItem",
			Action: func() (*http1.Request, error) {
				ctx := context.Background()

				randomItem := fetchRandomItem(ctx, c)
				if randomItem == nil {
					return nil, fmt.Errorf("retrieving random item: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveItemRequest(ctx, randomItem.ID)
			},
			Weight: 85,
		},
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateSomethingBlock(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingBlock(proj, typ)

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx := context.Background()

	itemInput := fake.BuildFakeItemCreationInput()

	return c.BuildCreateItemRequest(ctx, itemInput)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandomDependentIDFetchers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRandomDependentIDFetchers(proj, typ)

		expected := `
package main

import (
	"context"
)

func main() {
	ctx := context.Background()
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetSomethingBlock(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingBlock(proj, typ)

		expected := `
package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	randomItem := fetchRandomItem(ctx, c)
	if randomItem == nil {
		return nil, fmt.Errorf("retrieving random item: %w", ErrUnavailableYet)
	}

	return c.BuildGetItemRequest(ctx, randomItem.ID)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomethingBlock(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingBlock(proj, typ)

		expected := `
package main

import (
	"context"
)

func main() {
	ctx := context.Background()

	return c.BuildGetItemsRequest(ctx, nil)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateChildBlock(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateChildBlock(proj, typ)

		expected := `
package main

import (
	"context"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx := context.Background()

	if randomItem := fetchRandomItem(ctx, c); randomItem != nil {
		newItem := fake.BuildFakeItemCreationInput()
		randomItem.Name = newItem.Name
		randomItem.Details = newItem.Details
		return c.BuildUpdateItemRequest(ctx, randomItem)
	}

	return nil, ErrUnavailableYet
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveSomethingBlock(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingBlock(proj, typ)

		expected := `
package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	randomItem := fetchRandomItem(ctx, c)
	if randomItem == nil {
		return nil, fmt.Errorf("retrieving random item: %w", ErrUnavailableYet)
	}

	return c.BuildArchiveItemRequest(ctx, randomItem.ID)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
