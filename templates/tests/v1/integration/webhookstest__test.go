package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := webhooksTestDotGo(proj)

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

func checkWebhookEquality(t *testing.T, expected, actual *v1.Webhook) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ContentType, actual.ContentType)
	assert.Equal(t, expected.URL, actual.URL)
	assert.Equal(t, expected.Method, actual.Method)
	assert.NotZero(t, actual.CreatedOn)
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestWebhooks(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, premade)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkWebhookEquality(t, exampleWebhook, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhooks.
			var expected []*v1.Webhook
			for i := 0; i < 5; i++ {
				exampleWebhook := fake.BuildFakeWebhook()
				exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
				createdWebhook, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
				checkValueAndError(t, createdWebhook, err)

				expected = append(expected, createdWebhook)
			}

			// Assert webhook list equality.
			actual, err := todoClient.GetWebhooks(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(t, len(expected) <= len(actual.Webhooks))

			// Clean up.
			for _, webhook := range actual.Webhooks {
				err = todoClient.ArchiveWebhook(ctx, webhook.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch webhook.
			_, err := todoClient.GetWebhook(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Fetch webhook.
			actual, err := todoClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, actual)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhook.ID = nonexistentID

			err := todoClient.UpdateWebhook(ctx, exampleWebhook)
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Change webhook.
			premade.Name = reverse(premade.Name)
			exampleWebhook.Name = premade.Name
			err = todoClient.UpdateWebhook(ctx, premade)
			assert.NoError(t, err)

			// Fetch webhook.
			actual, err := todoClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCheckWebhookEquality(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildCheckWebhookEquality(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"testing"
)

func checkWebhookEquality(t *testing.T, expected, actual *v1.Webhook) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ContentType, actual.ContentType)
	assert.Equal(t, expected.URL, actual.URL)
	assert.Equal(t, expected.Method, actual.Method)
	assert.NotZero(t, actual.CreatedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildReverse(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildReverse()

		expected := `
package example

import ()

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestWebhooks(proj)

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

func TestWebhooks(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, premade)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, premade.ID)
			assert.NoError(t, err)

			actual, err := todoClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)
			checkWebhookEquality(t, exampleWebhook, actual)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhooks.
			var expected []*v1.Webhook
			for i := 0; i < 5; i++ {
				exampleWebhook := fake.BuildFakeWebhook()
				exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
				createdWebhook, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
				checkValueAndError(t, createdWebhook, err)

				expected = append(expected, createdWebhook)
			}

			// Assert webhook list equality.
			actual, err := todoClient.GetWebhooks(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(t, len(expected) <= len(actual.Webhooks))

			// Clean up.
			for _, webhook := range actual.Webhooks {
				err = todoClient.ArchiveWebhook(ctx, webhook.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Fetch webhook.
			_, err := todoClient.GetWebhook(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Fetch webhook.
			actual, err := todoClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, actual)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that doesn't exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhook.ID = nonexistentID

			err := todoClient.UpdateWebhook(ctx, exampleWebhook)
			assert.Error(t, err)
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Change webhook.
			premade.Name = reverse(premade.Name)
			exampleWebhook.Name = premade.Name
			err = todoClient.UpdateWebhook(ctx, premade)
			assert.NoError(t, err)

			// Fetch webhook.
			actual, err := todoClient.GetWebhook(ctx, premade.ID)
			checkValueAndError(t, actual, err)

			// Assert webhook equality.
			checkWebhookEquality(t, exampleWebhook, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, actual.ID)
			assert.NoError(t, err)
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create webhook.
			exampleWebhook := fake.BuildFakeWebhook()
			exampleWebhookInput := fake.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
			premade, err := todoClient.CreateWebhook(ctx, exampleWebhookInput)
			checkValueAndError(t, premade, err)

			// Clean up.
			err = todoClient.ArchiveWebhook(ctx, premade.ID)
			assert.NoError(t, err)
		})
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
