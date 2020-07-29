package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockReadCloserTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := mockReadCloserTestDotGo(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserInterfaceAssurance(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockReadCloserInterfaceAssurance()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockReadCloserDecl()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildNewMockReadCloser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildNewMockReadCloser()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockReadCloserReadHandler()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserClose(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMockReadCloserClose()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
