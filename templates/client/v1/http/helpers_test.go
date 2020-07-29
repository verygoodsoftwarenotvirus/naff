package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_helpersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := helpersDotGo(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotPointer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildArgIsNotPointer()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildArgIsNotNil()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotPointerOrNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildArgIsNotPointerOrNil()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildUnmarshalBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildUnmarshalBody(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateBodyFromStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildCreateBodyFromStruct()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
