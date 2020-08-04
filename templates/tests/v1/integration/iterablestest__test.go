package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
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
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"testing"
)

func checkItemEquality(t *testing.T, expected, actual *v1.Item) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Details, actual.Details, "expected Details for ID %d to be %v, but it was %v ", expected.ID, expected.Details, actual.Details)
	assert.NotZero(t, actual.CreatedOn)
}

func TestItems(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, createdItem)

			// Clean up.
			err = todoClient.ArchiveItem(ctx, createdItem.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)
			checkItemEquality(t, exampleItem, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create items.
			var expected []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItem := fake.BuildFakeItem()
				exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
				createdItem, itemCreationErr := todoClient.CreateItem(ctx, exampleItemInput)
				checkValueAndError(t, createdItem, itemCreationErr)

				expected = append(expected, createdItem)
			}

			// Assert item list equality.
			actual, err := todoClient.GetItems(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Items),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Items),
			)

			// Clean up.
			for _, createdItem := range actual.Items {
				err = todoClient.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Searching", func(T *testing.T) {
		T.Run("should be able to be search for items", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create items.
			exampleItem := fake.BuildFakeItem()
			var expected []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
				exampleItemInput.Name = fmt.Sprintf("%s %d", exampleItemInput.Name, i)
				createdItem, itemCreationErr := todoClient.CreateItem(ctx, exampleItemInput)
				checkValueAndError(t, createdItem, itemCreationErr)

				expected = append(expected, createdItem)
			}

			exampleLimit := uint8(20)

			// Assert item list equality.
			actual, err := todoClient.SearchItems(ctx, exampleItem.Name, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdItem := range expected {
				err = todoClient.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}
		})

		T.Run("should only receive your own items", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// create user and oauth2 client A.
			userA, err := testutil.CreateObligatoryUser(urlToUse, debug)
			require.NoError(t, err)

			ca, err := testutil.CreateObligatoryClient(urlToUse, userA)
			require.NoError(t, err)

			clientA, err := http.NewClient(
				ctx,
				ca.ClientID,
				ca.ClientSecret,
				todoClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				ca.Scopes,
				true,
			)
			checkValueAndError(test, clientA, err)

			// Create items for user A.
			exampleItemA := fake.BuildFakeItem()
			var createdForA []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItemInputA := fake.BuildFakeItemCreationInputFromItem(exampleItemA)
				exampleItemInputA.Name = fmt.Sprintf("%s %d", exampleItemInputA.Name, i)

				createdItem, itemCreationErr := clientA.CreateItem(ctx, exampleItemInputA)
				checkValueAndError(t, createdItem, itemCreationErr)

				createdForA = append(createdForA, createdItem)
			}

			exampleLimit := uint8(20)
			query := exampleItemA.Name

			// create user and oauth2 client B.
			userB, err := testutil.CreateObligatoryUser(urlToUse, debug)
			require.NoError(t, err)

			cb, err := testutil.CreateObligatoryClient(urlToUse, userB)
			require.NoError(t, err)

			clientB, err := http.NewClient(
				ctx,
				cb.ClientID,
				cb.ClientSecret,
				todoClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				cb.Scopes,
				true,
			)
			checkValueAndError(test, clientB, err)

			// Create items for user B.
			exampleItemB := fake.BuildFakeItem()
			exampleItemB.Name = reverse(exampleItemA.Name)
			var createdForB []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItemInputB := fake.BuildFakeItemCreationInputFromItem(exampleItemB)
				exampleItemInputB.Name = fmt.Sprintf("%s %d", exampleItemInputB.Name, i)

				createdItem, itemCreationErr := clientB.CreateItem(ctx, exampleItemInputB)
				checkValueAndError(t, createdItem, itemCreationErr)

				createdForB = append(createdForB, createdItem)
			}

			expected := createdForA

			// Assert item list equality.
			actual, err := clientA.SearchItems(ctx, query, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdItem := range createdForA {
				err = clientA.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}

			for _, createdItem := range createdForB {
				err = clientB.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent item.
			actual, err := todoClient.ItemExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant item exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Fetch item.
			actual, err := todoClient.ItemExists(ctx, createdItem.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent item.
			_, err := todoClient.GetItem(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Fetch item.
			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, actual)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleItem := fake.BuildFakeItem()
			exampleItem.ID = nonexistentID

			assert.Error(t, todoClient.UpdateItem(ctx, exampleItem))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Change item.
			createdItem.Update(exampleItem.ToUpdateInput())
			err = todoClient.UpdateItem(ctx, createdItem)
			assert.NoError(t, err)

			// Fetch item.
			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, todoClient.ArchiveItem(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)

		expected := `
package main

import ()

func example(ctx, createdItem) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

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

func Test_buildListArguments(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		varPrefix := "example"
		x := buildListArguments(proj, varPrefix, typ)

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

func Test_buildCheckSomethingEquality(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCheckSomethingEquality(proj, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func checkItemEquality(t *testing.T, expected, actual *v1.Item) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Details, actual.Details, "expected Details for ID %d to be %v, but it was %v ", expected.ID, expected.Details, actual.Details)
	assert.NotZero(t, actual.CreatedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestSomething(proj, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
	"testing"
)

func TestItems(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, createdItem)

			// Clean up.
			err = todoClient.ArchiveItem(ctx, createdItem.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)
			checkItemEquality(t, exampleItem, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create items.
			var expected []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItem := fake.BuildFakeItem()
				exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
				createdItem, itemCreationErr := todoClient.CreateItem(ctx, exampleItemInput)
				checkValueAndError(t, createdItem, itemCreationErr)

				expected = append(expected, createdItem)
			}

			// Assert item list equality.
			actual, err := todoClient.GetItems(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Items),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Items),
			)

			// Clean up.
			for _, createdItem := range actual.Items {
				err = todoClient.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Searching", func(T *testing.T) {
		T.Run("should be able to be search for items", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create items.
			exampleItem := fake.BuildFakeItem()
			var expected []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
				exampleItemInput.Name = fmt.Sprintf("%s %d", exampleItemInput.Name, i)
				createdItem, itemCreationErr := todoClient.CreateItem(ctx, exampleItemInput)
				checkValueAndError(t, createdItem, itemCreationErr)

				expected = append(expected, createdItem)
			}

			exampleLimit := uint8(20)

			// Assert item list equality.
			actual, err := todoClient.SearchItems(ctx, exampleItem.Name, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdItem := range expected {
				err = todoClient.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}
		})

		T.Run("should only receive your own items", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// create user and oauth2 client A.
			userA, err := testutil.CreateObligatoryUser(urlToUse, debug)
			require.NoError(t, err)

			ca, err := testutil.CreateObligatoryClient(urlToUse, userA)
			require.NoError(t, err)

			clientA, err := http.NewClient(
				ctx,
				ca.ClientID,
				ca.ClientSecret,
				todoClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				ca.Scopes,
				true,
			)
			checkValueAndError(test, clientA, err)

			// Create items for user A.
			exampleItemA := fake.BuildFakeItem()
			var createdForA []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItemInputA := fake.BuildFakeItemCreationInputFromItem(exampleItemA)
				exampleItemInputA.Name = fmt.Sprintf("%s %d", exampleItemInputA.Name, i)

				createdItem, itemCreationErr := clientA.CreateItem(ctx, exampleItemInputA)
				checkValueAndError(t, createdItem, itemCreationErr)

				createdForA = append(createdForA, createdItem)
			}

			exampleLimit := uint8(20)
			query := exampleItemA.Name

			// create user and oauth2 client B.
			userB, err := testutil.CreateObligatoryUser(urlToUse, debug)
			require.NoError(t, err)

			cb, err := testutil.CreateObligatoryClient(urlToUse, userB)
			require.NoError(t, err)

			clientB, err := http.NewClient(
				ctx,
				cb.ClientID,
				cb.ClientSecret,
				todoClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				cb.Scopes,
				true,
			)
			checkValueAndError(test, clientB, err)

			// Create items for user B.
			exampleItemB := fake.BuildFakeItem()
			exampleItemB.Name = reverse(exampleItemA.Name)
			var createdForB []*v1.Item
			for i := 0; i < 5; i++ {
				// Create item.
				exampleItemInputB := fake.BuildFakeItemCreationInputFromItem(exampleItemB)
				exampleItemInputB.Name = fmt.Sprintf("%s %d", exampleItemInputB.Name, i)

				createdItem, itemCreationErr := clientB.CreateItem(ctx, exampleItemInputB)
				checkValueAndError(t, createdItem, itemCreationErr)

				createdForB = append(createdForB, createdItem)
			}

			expected := createdForA

			// Assert item list equality.
			actual, err := clientA.SearchItems(ctx, query, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdItem := range createdForA {
				err = clientA.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}

			for _, createdItem := range createdForB {
				err = clientB.ArchiveItem(ctx, createdItem.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent item.
			actual, err := todoClient.ItemExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant item exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Fetch item.
			actual, err := todoClient.ItemExists(ctx, createdItem.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent item.
			_, err := todoClient.GetItem(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Fetch item.
			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, actual)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleItem := fake.BuildFakeItem()
			exampleItem.ID = nonexistentID

			assert.Error(t, todoClient.UpdateItem(ctx, exampleItem))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Change item.
			createdItem.Update(exampleItem.ToUpdateInput())
			err = todoClient.UpdateItem(ctx, createdItem)
			assert.NoError(t, err)

			// Fetch item.
			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, todoClient.ArchiveItem(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

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

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

// Create item. exampleItem :=  fake.BuildFakeItem  () exampleItemInput :=  fake.BuildFakeItemCreationInputFromItem  (exampleItem) createdItem,err := todoClient . CreateItem (ctx,exampleItemInput) checkValueAndError (t,createdItem,err)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCreationCodeWithoutType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRequisiteCreationCodeWithoutType(proj, typ)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCleanupCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		includeSelf := true
		x := buildRequisiteCleanupCode(proj, typ, includeSelf)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
)

// Clean up item.  assert.NoError  (t,todoClient . ArchiveItem (ctx,createdItem . ID))
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildParamsForMethodThatHandlesAnInstanceWithStructsButIDsOnly(proj, typ)

		expected := `
package main

import ()

func example(ctx, createdItem.ID) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEqualityCheckLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildEqualityCheckLines(typ)

		expected := `
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Details, actual.Details, "expected Details for ID %d to be %v, but it was %v ", expected.ID, expected.Details, actual.Details)
	assert.NotZero(t, actual.CreatedOn)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCreationCodeFor404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		indexToStop := 1
		x := buildRequisiteCreationCodeFor404Tests(proj, typ, indexToStop)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCleanupCodeFor404s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		indexToStop := 1
		x := buildRequisiteCleanupCodeFor404s(proj, typ, indexToStop)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreationArgumentsFor404s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		varPrefix := "example"
		indexToStop := 1
		x := buildCreationArgumentsFor404s(proj, varPrefix, typ, indexToStop)

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

func Test_buildSubtestsForCreation404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSubtestsForCreation404Tests(proj, typ)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestCreating(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestCreating(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func main() {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, createdItem)

			// Clean up.
			err = todoClient.ArchiveItem(ctx, createdItem.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)
			checkItemEquality(t, exampleItem, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})
}
`
		actual := testutils.RenderIndependentStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestListing(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestListing(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Create items.
	var expected []*v1.Item
	for i := 0; i < 5; i++ {
		// Create item.
		exampleItem := fake.BuildFakeItem()
		exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
		createdItem, itemCreationErr := todoClient.CreateItem(ctx, exampleItemInput)
		checkValueAndError(t, createdItem, itemCreationErr)

		expected = append(expected, createdItem)
	}

	// Assert item list equality.
	actual, err := todoClient.GetItems(ctx, nil)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual.Items),
		"expected %d to be <= %d",
		len(expected),
		len(actual.Items),
	)

	// Clean up.
	for _, createdItem := range actual.Items {
		err = todoClient.ArchiveItem(ctx, createdItem.ID)
		assert.NoError(t, err)
	}
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestSearching(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestSearching(proj, typ)

		expected := `
package main

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Create items.
	exampleItem := fake.BuildFakeItem()
	var expected []*v1.Item
	for i := 0; i < 5; i++ {
		// Create item.
		exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
		exampleItemInput.Name = fmt.Sprintf("%s %d", exampleItemInput.Name, i)
		createdItem, itemCreationErr := todoClient.CreateItem(ctx, exampleItemInput)
		checkValueAndError(t, createdItem, itemCreationErr)

		expected = append(expected, createdItem)
	}

	exampleLimit := uint8(20)

	// Assert item list equality.
	actual, err := todoClient.SearchItems(ctx, exampleItem.Name, exampleLimit)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual),
		"expected results length %d to be <= %d",
		len(expected),
		len(actual),
	)

	// Clean up.
	for _, createdItem := range expected {
		err = todoClient.ArchiveItem(ctx, createdItem.ID)
		assert.NoError(t, err)
	}
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestSearchingForOnlyYourOwnItems(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestSearchingForOnlyYourOwnItems(proj, typ)

		expected := `
package main

import (
	"context"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	noop "gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	testutil "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/tests/v1/testutil"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// create user and oauth2 client A.
	userA, err := testutil.CreateObligatoryUser(urlToUse, debug)
	require.NoError(t, err)

	ca, err := testutil.CreateObligatoryClient(urlToUse, userA)
	require.NoError(t, err)

	clientA, err := http.NewClient(
		ctx,
		ca.ClientID,
		ca.ClientSecret,
		todoClient.URL,
		noop.ProvideNoopLogger(),
		buildHTTPClient(),
		ca.Scopes,
		true,
	)
	checkValueAndError(test, clientA, err)

	// Create items for user A.
	exampleItemA := fake.BuildFakeItem()
	var createdForA []*v1.Item
	for i := 0; i < 5; i++ {
		// Create item.
		exampleItemInputA := fake.BuildFakeItemCreationInputFromItem(exampleItemA)
		exampleItemInputA.Name = fmt.Sprintf("%s %d", exampleItemInputA.Name, i)

		createdItem, itemCreationErr := clientA.CreateItem(ctx, exampleItemInputA)
		checkValueAndError(t, createdItem, itemCreationErr)

		createdForA = append(createdForA, createdItem)
	}

	exampleLimit := uint8(20)
	query := exampleItemA.Name

	// create user and oauth2 client B.
	userB, err := testutil.CreateObligatoryUser(urlToUse, debug)
	require.NoError(t, err)

	cb, err := testutil.CreateObligatoryClient(urlToUse, userB)
	require.NoError(t, err)

	clientB, err := http.NewClient(
		ctx,
		cb.ClientID,
		cb.ClientSecret,
		todoClient.URL,
		noop.ProvideNoopLogger(),
		buildHTTPClient(),
		cb.Scopes,
		true,
	)
	checkValueAndError(test, clientB, err)

	// Create items for user B.
	exampleItemB := fake.BuildFakeItem()
	exampleItemB.Name = reverse(exampleItemA.Name)
	var createdForB []*v1.Item
	for i := 0; i < 5; i++ {
		// Create item.
		exampleItemInputB := fake.BuildFakeItemCreationInputFromItem(exampleItemB)
		exampleItemInputB.Name = fmt.Sprintf("%s %d", exampleItemInputB.Name, i)

		createdItem, itemCreationErr := clientB.CreateItem(ctx, exampleItemInputB)
		checkValueAndError(t, createdItem, itemCreationErr)

		createdForB = append(createdForB, createdItem)
	}

	expected := createdForA

	// Assert item list equality.
	actual, err := clientA.SearchItems(ctx, query, exampleLimit)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual),
		"expected results length %d to be <= %d",
		len(expected),
		len(actual),
	)

	// Clean up.
	for _, createdItem := range createdForA {
		err = clientA.ArchiveItem(ctx, createdItem.ID)
		assert.NoError(t, err)
	}

	for _, createdItem := range createdForB {
		err = clientB.ArchiveItem(ctx, createdItem.ID)
		assert.NoError(t, err)
	}
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestExistenceCheckingShouldFailWhenTryingToReadSomethingThatDoesNotExist(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestExistenceCheckingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Attempt to fetch nonexistent item.
	actual, err := todoClient.ItemExists(ctx, nonexistentID)
	assert.NoError(t, err)
	assert.False(t, actual)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestExistenceCheckingShouldBeReadable(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestExistenceCheckingShouldBeReadable(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Create item.
	exampleItem := fake.BuildFakeItem()
	exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
	createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
	checkValueAndError(t, createdItem, err)

	// Fetch item.
	actual, err := todoClient.ItemExists(ctx, createdItem.ID)
	assert.NoError(t, err)
	assert.True(t, actual)

	// Clean up item.
	assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestReadingShouldFailWhenTryingToReadSomethingThatDoesNotExist(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Attempt to fetch nonexistent item.
	_, err := todoClient.GetItem(ctx, nonexistentID)
	assert.Error(t, err)
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestReadingShouldBeReadable(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestReadingShouldBeReadable(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Create item.
	exampleItem := fake.BuildFakeItem()
	exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
	createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
	checkValueAndError(t, createdItem, err)

	// Fetch item.
	actual, err := todoClient.GetItem(ctx, createdItem.ID)
	checkValueAndError(t, actual, err)

	// Assert item equality.
	checkItemEquality(t, exampleItem, actual)

	// Clean up item.
	assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildParamsForCheckingATypeThatDoesNotExistButIncludesPredecessorID(proj, typ)

		expected := `
package main

import ()

func example(ctx, nonexistentID) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForCheckingATypeThatDoesNotExistAndIncludesItsOwnerVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildParamsForCheckingATypeThatDoesNotExistAndIncludesItsOwnerVar(proj, typ)

		expected := `
package main

import ()

func example(ctx, exampleItem) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCreationCodeForUpdate404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		nonexistentArgIndex := 1
		x := buildRequisiteCreationCodeForUpdate404Tests(proj, typ, nonexistentArgIndex)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

// Create item. exampleItem :=  fake.BuildFakeItem  () exampleItemInput :=  fake.BuildFakeItemCreationInputFromItem  (exampleItem) createdItem,err := todoClient . CreateItem (ctx,exampleItemInput) checkValueAndError (t,createdItem,err)
// Change item. createdItem . Update (exampleItem . ToUpdateInput ()) err = todoClient . UpdateItem (ctx,createdItem)  assert.Error  (t,err)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCleanupCodeForUpdate404s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildRequisiteCleanupCodeForUpdate404s(proj, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
)

// Clean up item.  assert.NoError  (t,todoClient . ArchiveItem (ctx,createdItem . ID))
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSubtestsForUpdate404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSubtestsForUpdate404Tests(proj, typ)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestUpdating(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestUpdating(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func main() {
	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleItem := fake.BuildFakeItem()
			exampleItem.ID = nonexistentID

			assert.Error(t, todoClient.UpdateItem(ctx, exampleItem))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Change item.
			createdItem.Update(exampleItem.ToUpdateInput())
			err = todoClient.UpdateItem(ctx, createdItem)
			assert.NoError(t, err)

			// Fetch item.
			actual, err := todoClient.GetItem(ctx, createdItem.ID)
			checkValueAndError(t, actual, err)

			// Assert item equality.
			checkItemEquality(t, exampleItem, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})
}
`
		actual := testutils.RenderFunctionBodyToString(t, []jen.Code{x})

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestUpdatingShouldBeUpdatable(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestUpdatingShouldBeUpdatable(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Create item.
	exampleItem := fake.BuildFakeItem()
	exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
	createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
	checkValueAndError(t, createdItem, err)

	// Change item.
	createdItem.Update(exampleItem.ToUpdateInput())
	err = todoClient.UpdateItem(ctx, createdItem)
	assert.NoError(t, err)

	// Fetch item.
	actual, err := todoClient.GetItem(ctx, createdItem.ID)
	checkValueAndError(t, actual, err)

	// Assert item equality.
	checkItemEquality(t, exampleItem, actual)
	assert.NotNil(t, actual.LastUpdatedOn)

	// Clean up item.
	assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestUpdatingShouldFailWhenTryingToChangeSomethingThatDoesNotExist(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	exampleItem := fake.BuildFakeItem()
	exampleItem.ID = nonexistentID

	assert.Error(t, todoClient.UpdateItem(ctx, exampleItem))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCreationCodeFor404DeletionTests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		nonexistentArgIndex := 1
		x := buildRequisiteCreationCodeFor404DeletionTests(proj, typ, nonexistentArgIndex)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCleanupCodeFor404DeletionTests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		indexToStop := 1
		x := buildRequisiteCleanupCodeFor404DeletionTests(proj, typ, indexToStop)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
)

// Clean up item.  assert.NoError  (t,todoClient . ArchiveItem (ctx,createdItem . ID))
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForDeletionWithNonexistentID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		indexToNotExist := 1
		x := buildParamsForDeletionWithNonexistentID(proj, typ, indexToNotExist)

		expected := `
package main

import ()

func example(ctx, createdItem.ID) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSubtestsForDeletion404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSubtestsForDeletion404Tests(proj, typ)

		expected := `
package example

import ()
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDeleting(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDeleting(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func main() {
	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, todoClient.ArchiveItem(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create item.
			exampleItem := fake.BuildFakeItem()
			exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
			createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
			checkValueAndError(t, createdItem, err)

			// Clean up item.
			assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
		})
	})
}
`
		actual := testutils.RenderIndependentStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDeletingShouldBeAbleToBeDeleted(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDeletingShouldBeAbleToBeDeleted(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	// Create item.
	exampleItem := fake.BuildFakeItem()
	exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
	createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
	checkValueAndError(t, createdItem, err)

	// Clean up item.
	assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDeletingShouldFailForNonexistent(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDeletingShouldFailForNonexistent(proj, typ)

		expected := `
package main

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
)

func main() {
	ctx, span := tracing.StartSpan(context.Background(), t.Name())
	defer span.End()

	assert.Error(t, todoClient.ArchiveItem(ctx, nonexistentID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
