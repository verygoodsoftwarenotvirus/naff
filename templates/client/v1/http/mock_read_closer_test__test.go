package client

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mockReadCloserTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		out := mockReadCloserTestDotGo(proj)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserInterfaceAssurance(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildMockReadCloserInterfaceAssurance()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildMockReadCloserDecl()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildNewMockReadCloser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildNewMockReadCloser()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildMockReadCloserReadHandler()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMockReadCloserClose(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(buildMockReadCloserClose()...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
