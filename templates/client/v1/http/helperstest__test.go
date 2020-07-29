package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_helpersTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := helpersTestDotGo(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildHelperTestingType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildHelperTestingType()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestArgIsNotPointerOrNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestArgIsNotPointerOrNil()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestArgIsNotPointer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestArgIsNotPointer()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestArgIsNotNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestArgIsNotNil()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestUnmarshalBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildTestUnmarshalBody(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildHelperTestBreakableStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildHelperTestBreakableStruct()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestCreateBodyFromStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestCreateBodyFromStruct()

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
