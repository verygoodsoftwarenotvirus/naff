package client

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_iterablesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := iterablesDotGo(proj, typ)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_attachURIToSpanCall(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(attachURIToSpanCall(proj))

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		firstVar := jen.Nil()

		out := jen.NewFile("main")
		out.Add(buildV1ClientURLBuildingParamsForSingleInstanceOfSomething(proj, firstVar, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForCreatingSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		firstVar := jen.Nil()

		out := jen.NewFile("main")
		out.Add(buildV1ClientURLBuildingParamsForCreatingSomething(proj, firstVar, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForListOfSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		firstVar := jen.Nil()

		out := jen.NewFile("main")
		out.Add(buildV1ClientURLBuildingParamsForListOfSomething(proj, firstVar, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForSearchingSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		firstVar := jen.Nil()

		out := jen.NewFile("main")
		out.Add(buildV1ClientURLBuildingParamsForSearchingSomething(proj, firstVar, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]
		firstVar := jen.Nil()

		out := jen.NewFile("main")
		out.Add(buildV1ClientURLBuildingParamsForMethodThatIncludesItsOwnType(proj, firstVar, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildSomethingExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildSomethingExistsRequest(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildSomethingExists(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildSomethingExists(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildGetSomethingRequestFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildGetSomethingFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildSearchSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildSearchSomethingRequestFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildSearchSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildSearchSomethingFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetListOfSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildGetListOfSomethingRequestFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildGetListOfSomethingFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildCreateSomethingRequestFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildCreateSomethingFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildUpdateSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildUpdateSomethingRequestFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildUpdateSomethingFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveSomethingRequestFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildBuildArchiveSomethingRequestFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		typ := proj.DataTypes[0]

		out := jen.NewFile("main")
		out.Add(buildArchiveSomethingFuncDecl(proj, typ)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
