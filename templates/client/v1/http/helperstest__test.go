package client

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_helpersTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		out := helpersTestDotGo(proj)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildHelperTestingType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildHelperTestingType()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestArgIsNotPointerOrNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildTestArgIsNotPointerOrNil()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestArgIsNotPointer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildTestArgIsNotPointer()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestArgIsNotNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildTestArgIsNotNil()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestUnmarshalBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildTestUnmarshalBody(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildHelperTestBreakableStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildHelperTestBreakableStruct()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestCreateBodyFromStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildTestCreateBodyFromStruct()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
