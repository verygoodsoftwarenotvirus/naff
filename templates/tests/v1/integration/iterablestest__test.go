package integration

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
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

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildParamsForMethodThatIncludesItsOwnTypeInItsParamsAndHasFullStructs(proj, typ)

		expected := `
package main

import ()

func example(ctx, createdThing.ID, createdYetAnotherThing) {}
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

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		varPrefix := "example"
		x := buildCreationArguments(proj, varPrefix, typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		ctx,
		exampleThing.ID,
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

	T.Run("without search enabled", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.SearchEnabled = false

		x := buildTestSomething(proj, typ)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	tracing "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/tracing"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
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
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildRequisiteCreationCode(proj, typ)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

	// Create yet another thing.
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
	exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
	createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
	checkValueAndError(t, createdYetAnotherThing, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCreationCodeWithoutType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildRequisiteCreationCodeWithoutType(proj, typ)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

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
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {

	// Clean up item.
	assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		includeSelf := true
		x := buildRequisiteCleanupCode(proj, typ, includeSelf)

		expected := `
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {

	// Clean up another thing.
	assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

	// Clean up thing.
	assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))

	// Clean up yet another thing.
	assert.NoError(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

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

	T.Run("with gamut", func(t *testing.T) {
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.LastDataType()

		x := buildEqualityCheckLines(typ)

		expected := `
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.String, actual.String, "expected String for ID %d to be %v, but it was %v ", expected.ID, expected.String, actual.String)
	assert.Equal(t, *expected.PointerToString, *actual.PointerToString, "expected PointerToString to be %v, but it was %v ", expected.PointerToString, actual.PointerToString)
	assert.Equal(t, expected.Bool, actual.Bool, "expected Bool for ID %d to be %v, but it was %v ", expected.ID, expected.Bool, actual.Bool)
	assert.Equal(t, *expected.PointerToBool, *actual.PointerToBool, "expected PointerToBool to be %v, but it was %v ", expected.PointerToBool, actual.PointerToBool)
	assert.Equal(t, expected.Int, actual.Int, "expected Int for ID %d to be %v, but it was %v ", expected.ID, expected.Int, actual.Int)
	assert.Equal(t, *expected.PointerToInt, *actual.PointerToInt, "expected PointerToInt to be %v, but it was %v ", expected.PointerToInt, actual.PointerToInt)
	assert.Equal(t, expected.Int8, actual.Int8, "expected Int8 for ID %d to be %v, but it was %v ", expected.ID, expected.Int8, actual.Int8)
	assert.Equal(t, *expected.PointerToInt8, *actual.PointerToInt8, "expected PointerToInt8 to be %v, but it was %v ", expected.PointerToInt8, actual.PointerToInt8)
	assert.Equal(t, expected.Int16, actual.Int16, "expected Int16 for ID %d to be %v, but it was %v ", expected.ID, expected.Int16, actual.Int16)
	assert.Equal(t, *expected.PointerToInt16, *actual.PointerToInt16, "expected PointerToInt16 to be %v, but it was %v ", expected.PointerToInt16, actual.PointerToInt16)
	assert.Equal(t, expected.Int32, actual.Int32, "expected Int32 for ID %d to be %v, but it was %v ", expected.ID, expected.Int32, actual.Int32)
	assert.Equal(t, *expected.PointerToInt32, *actual.PointerToInt32, "expected PointerToInt32 to be %v, but it was %v ", expected.PointerToInt32, actual.PointerToInt32)
	assert.Equal(t, expected.Int64, actual.Int64, "expected Int64 for ID %d to be %v, but it was %v ", expected.ID, expected.Int64, actual.Int64)
	assert.Equal(t, *expected.PointerToInt64, *actual.PointerToInt64, "expected PointerToInt64 to be %v, but it was %v ", expected.PointerToInt64, actual.PointerToInt64)
	assert.Equal(t, expected.Uint, actual.Uint, "expected Uint for ID %d to be %v, but it was %v ", expected.ID, expected.Uint, actual.Uint)
	assert.Equal(t, *expected.PointerToUint, *actual.PointerToUint, "expected PointerToUint to be %v, but it was %v ", expected.PointerToUint, actual.PointerToUint)
	assert.Equal(t, expected.Uint8, actual.Uint8, "expected Uint8 for ID %d to be %v, but it was %v ", expected.ID, expected.Uint8, actual.Uint8)
	assert.Equal(t, *expected.PointerToUint8, *actual.PointerToUint8, "expected PointerToUint8 to be %v, but it was %v ", expected.PointerToUint8, actual.PointerToUint8)
	assert.Equal(t, expected.Uint16, actual.Uint16, "expected Uint16 for ID %d to be %v, but it was %v ", expected.ID, expected.Uint16, actual.Uint16)
	assert.Equal(t, *expected.PointerToUint16, *actual.PointerToUint16, "expected PointerToUint16 to be %v, but it was %v ", expected.PointerToUint16, actual.PointerToUint16)
	assert.Equal(t, expected.Uint32, actual.Uint32, "expected Uint32 for ID %d to be %v, but it was %v ", expected.ID, expected.Uint32, actual.Uint32)
	assert.Equal(t, *expected.PointerToUint32, *actual.PointerToUint32, "expected PointerToUint32 to be %v, but it was %v ", expected.PointerToUint32, actual.PointerToUint32)
	assert.Equal(t, expected.Uint64, actual.Uint64, "expected Uint64 for ID %d to be %v, but it was %v ", expected.ID, expected.Uint64, actual.Uint64)
	assert.Equal(t, *expected.PointerToUint64, *actual.PointerToUint64, "expected PointerToUint64 to be %v, but it was %v ", expected.PointerToUint64, actual.PointerToUint64)
	assert.Equal(t, expected.Float32, actual.Float32, "expected Float32 for ID %d to be %v, but it was %v ", expected.ID, expected.Float32, actual.Float32)
	assert.Equal(t, *expected.PointerToFloat32, *actual.PointerToFloat32, "expected PointerToFloat32 to be %v, but it was %v ", expected.PointerToFloat32, actual.PointerToFloat32)
	assert.Equal(t, expected.Float64, actual.Float64, "expected Float64 for ID %d to be %v, but it was %v ", expected.ID, expected.Float64, actual.Float64)
	assert.Equal(t, *expected.PointerToFloat64, *actual.PointerToFloat64, "expected PointerToFloat64 to be %v, but it was %v ", expected.PointerToFloat64, actual.PointerToFloat64)
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
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		indexToStop := len(proj.DataTypes)
		x := buildRequisiteCreationCodeFor404Tests(proj, typ, indexToStop)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		indexToStop := len(proj.DataTypes) / 2
		x := buildRequisiteCreationCodeFor404Tests(proj, typ, indexToStop)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCleanupCodeFor404s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		indexToStop := len(proj.DataTypes)
		x := buildRequisiteCleanupCodeFor404s(proj, typ, indexToStop)

		expected := `
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {

	// Clean up another thing.
	assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

	// Clean up thing.
	assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreationArgumentsFor404s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		varPrefix := "example"
		indexToStop := 1
		x := buildCreationArgumentsFor404s(proj, varPrefix, typ, indexToStop)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		ctx,
		exampleThing.ID,
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
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildSubtestsForCreation404Tests(proj, typ)

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

	T.Run("should fail to create for nonexistent thing", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = nonexistentID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, nonexistentID, exampleYetAnotherThingInput)

		assert.Nil(t, createdYetAnotherThing)
		assert.Error(t, err)
	})

	T.Run("should fail to create for nonexistent another thing", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// Create thing.
		exampleThing := fake.BuildFakeThing()
		exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
		createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
		checkValueAndError(t, createdThing, err)

		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = nonexistentID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)

		assert.Nil(t, createdYetAnotherThing)
		assert.Error(t, err)

		// Clean up thing.
		assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
	})
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

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

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

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

	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

	// Create yet another things.
	var expected []*v1.YetAnotherThing
	for i := 0; i < 5; i++ {
		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, yetAnotherThingCreationErr := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
		checkValueAndError(t, createdYetAnotherThing, yetAnotherThingCreationErr)

		expected = append(expected, createdYetAnotherThing)
	}

	// Assert yet another thing list equality.
	actual, err := todoClient.GetYetAnotherThings(ctx, createdThing.ID, createdAnotherThing.ID, nil)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual.YetAnotherThings),
		"expected %d to be <= %d",
		len(expected),
		len(actual.YetAnotherThings),
	)

	// Clean up.
	for _, createdYetAnotherThing := range actual.YetAnotherThings {
		err = todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID)
		assert.NoError(t, err)
	}

	// Clean up another thing.
	assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

	// Clean up thing.
	assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
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

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		for i := range proj.DataTypes {
			proj.DataTypes[i].SearchEnabled = true
		}

		typ := proj.LastDataType()
		typ.Fields = testprojects.BuildTodoApp().DataTypes[0].Fields

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

	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

	// Create yet another things.
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	var expected []*v1.YetAnotherThing
	for i := 0; i < 5; i++ {
		// Create yet another thing.
		exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		exampleYetAnotherThingInput.Name = fmt.Sprintf("%s %d", exampleYetAnotherThingInput.Name, i)
		createdYetAnotherThing, yetAnotherThingCreationErr := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
		checkValueAndError(t, createdYetAnotherThing, yetAnotherThingCreationErr)

		expected = append(expected, createdYetAnotherThing)
	}

	exampleLimit := uint8(20)

	// Assert yet another thing list equality.
	actual, err := todoClient.SearchYetAnotherThings(ctx, exampleYetAnotherThing.Name, exampleLimit)
	checkValueAndError(t, actual, err)
	assert.True(
		t,
		len(expected) <= len(actual),
		"expected results length %d to be <= %d",
		len(expected),
		len(actual),
	)

	// Clean up.
	for _, createdYetAnotherThing := range expected {
		err = todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID)
		assert.NoError(t, err)
	}

	// Clean up another thing.
	assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

	// Clean up thing.
	assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
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

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildParamsForCheckingATypeThatDoesNotExistAndIncludesItsOwnerVar(proj, typ)

		expected := `
package main

import ()

func example(ctx, createdThing.ID, exampleYetAnotherThing) {}
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
package main

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create item.
	exampleItem := fake.BuildFakeItem()
	exampleItemInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)
	createdItem, err := todoClient.CreateItem(ctx, exampleItemInput)
	checkValueAndError(t, createdItem, err)

	// Change item.
	createdItem.Update(exampleItem.ToUpdateInput())
	err = todoClient.UpdateItem(ctx, createdItem)
	assert.Error(t, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		nonexistentArgIndex := len(proj.DataTypes) - 1
		x := buildRequisiteCreationCodeForUpdate404Tests(proj, typ, nonexistentArgIndex)

		expected := `
package main

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

	// Create yet another thing.
	exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
	exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
	exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
	createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
	checkValueAndError(t, createdYetAnotherThing, err)

	// Change yet another thing.
	createdYetAnotherThing.Update(exampleYetAnotherThing.ToUpdateInput())
	err = todoClient.UpdateYetAnotherThing(ctx, createdThing.ID, createdYetAnotherThing)
	assert.Error(t, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRequisiteCleanupCodeForUpdate404s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildRequisiteCleanupCodeForUpdate404s(proj, typ)

		expected := `
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {

	// Clean up yet another thing.
	assert.NoError(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID))

	// Clean up another thing.
	assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

	// Clean up thing.
	assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSubtestsForUpdate404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildSubtestsForUpdate404Tests(proj, typ)

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

	T.Run("it should return an error when trying to update something that belongs to a thing that does not exist", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// Create thing.
		exampleThing := fake.BuildFakeThing()
		exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
		createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
		checkValueAndError(t, createdThing, err)

		// Create another thing.
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = createdThing.ID
		exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
		createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
		checkValueAndError(t, createdAnotherThing, err)

		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
		checkValueAndError(t, createdYetAnotherThing, err)

		// Change yet another thing.
		createdYetAnotherThing.Update(exampleYetAnotherThing.ToUpdateInput())
		err = todoClient.UpdateYetAnotherThing(ctx, nonexistentID, createdYetAnotherThing)
		assert.Error(t, err)

		// Clean up yet another thing.
		assert.NoError(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID))

		// Clean up another thing.
		assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

		// Clean up thing.
		assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
	})

	T.Run("it should return an error when trying to update something that belongs to an another thing that does not exist", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// Create thing.
		exampleThing := fake.BuildFakeThing()
		exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
		createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
		checkValueAndError(t, createdThing, err)

		// Create another thing.
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = createdThing.ID
		exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
		createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
		checkValueAndError(t, createdAnotherThing, err)

		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
		checkValueAndError(t, createdYetAnotherThing, err)

		// Change yet another thing.
		createdYetAnotherThing.Update(exampleYetAnotherThing.ToUpdateInput())
		createdYetAnotherThing.BelongsToAnotherThing = nonexistentID
		err = todoClient.UpdateYetAnotherThing(ctx, createdThing.ID, createdYetAnotherThing)
		assert.Error(t, err)

		// Clean up yet another thing.
		assert.NoError(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID))

		// Clean up another thing.
		assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

		// Clean up thing.
		assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
	})
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

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
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		nonexistentArgIndex := 1
		x := buildRequisiteCreationCodeFor404DeletionTests(proj, typ, nonexistentArgIndex)

		expected := `
package main

import (
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
)

func main() {
	// Create thing.
	exampleThing := fake.BuildFakeThing()
	exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
	createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
	checkValueAndError(t, createdThing, err)

	// Create another thing.
	exampleAnotherThing := fake.BuildFakeAnotherThing()
	exampleAnotherThing.BelongsToThing = createdThing.ID
	exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
	createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
	checkValueAndError(t, createdAnotherThing, err)

}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

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
package main

import (
	assert "github.com/stretchr/testify/assert"
)

func main() {

	// Clean up item.
	assert.NoError(t, todoClient.ArchiveItem(ctx, createdItem.ID))
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParamsForDeletionWithNonexistentID(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		indexToNotExist := 1
		x := buildParamsForDeletionWithNonexistentID(proj, typ, indexToNotExist)

		expected := `
package main

import ()

func example(ctx, createdThing.ID, nonexistentID, createdYetAnotherThing.ID) {}
`
		actual := testutils.RenderFunctionParamsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSubtestsForDeletion404Tests(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()

		x := buildSubtestsForDeletion404Tests(proj, typ)

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

	T.Run("returns error when trying to archive post belonging to nonexistent thing", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// Create thing.
		exampleThing := fake.BuildFakeThing()
		exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
		createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
		checkValueAndError(t, createdThing, err)

		// Create another thing.
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = createdThing.ID
		exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
		createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
		checkValueAndError(t, createdAnotherThing, err)

		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
		checkValueAndError(t, createdYetAnotherThing, err)

		assert.Error(t, todoClient.ArchiveYetAnotherThing(ctx, nonexistentID, createdAnotherThing.ID, createdYetAnotherThing.ID))

		// Clean up yet another thing.
		assert.NoError(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID))

		// Clean up another thing.
		assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

		// Clean up thing.
		assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
	})

	T.Run("returns error when trying to archive post belonging to nonexistent another thing", func(t *testing.T) {
		ctx, span := tracing.StartSpan(context.Background(), t.Name())
		defer span.End()

		// Create thing.
		exampleThing := fake.BuildFakeThing()
		exampleThingInput := fake.BuildFakeThingCreationInputFromThing(exampleThing)
		createdThing, err := todoClient.CreateThing(ctx, exampleThingInput)
		checkValueAndError(t, createdThing, err)

		// Create another thing.
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = createdThing.ID
		exampleAnotherThingInput := fake.BuildFakeAnotherThingCreationInputFromAnotherThing(exampleAnotherThing)
		createdAnotherThing, err := todoClient.CreateAnotherThing(ctx, exampleAnotherThingInput)
		checkValueAndError(t, createdAnotherThing, err)

		// Create yet another thing.
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = createdAnotherThing.ID
		exampleYetAnotherThingInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)
		createdYetAnotherThing, err := todoClient.CreateYetAnotherThing(ctx, createdThing.ID, exampleYetAnotherThingInput)
		checkValueAndError(t, createdYetAnotherThing, err)

		assert.Error(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, nonexistentID, createdYetAnotherThing.ID))

		// Clean up yet another thing.
		assert.NoError(t, todoClient.ArchiveYetAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID, createdYetAnotherThing.ID))

		// Clean up another thing.
		assert.NoError(t, todoClient.ArchiveAnotherThing(ctx, createdThing.ID, createdAnotherThing.ID))

		// Clean up thing.
		assert.NoError(t, todoClient.ArchiveThing(ctx, createdThing.ID))
	})
}
`
		actual := testutils.RenderFunctionBodyToString(t, x)

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
