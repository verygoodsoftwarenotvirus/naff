package client

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_helpersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		out := helpersDotGo(proj)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotPointer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildArgIsNotPointer()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildArgIsNotNil()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArgIsNotPointerOrNil(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildArgIsNotPointerOrNil()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildUnmarshalBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildUnmarshalBody(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateBodyFromStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildCreateBodyFromStruct()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
