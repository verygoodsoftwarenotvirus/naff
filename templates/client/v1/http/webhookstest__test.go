package client

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

		proj := testprojects.TodoApp
		x := webhooksTestDotGo(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildGetWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientBuildGetWebhookRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientGetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientGetWebhook(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildGetWebhooksRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientBuildGetWebhooksRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientGetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientGetWebhooks(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildCreateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientBuildCreateWebhookRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientCreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientCreateWebhook(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildUpdateWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientBuildUpdateWebhookRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientUpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientUpdateWebhook(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientBuildArchiveWebhookRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientBuildArchiveWebhookRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1ClientArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestV1ClientArchiveWebhook(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
