package models

import (
	"bytes"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/assert"
)

func buildObligatorySuperPalabra() wordsmith.SuperPalabra {
	return wordsmith.FromSingularPascalCase("Thing")
}

func renderFunctionParamsToString(t *testing.T, params []jen.Code) string {
	t.Helper()

	f := jen.NewFile("main")
	f.Add(jen.Func().ID("example").Params(params...).Block())
	b := bytes.NewBufferString("\n")
	require.NoError(t, f.Render(b))

	return b.String()
}

func renderCallArgsToString(t *testing.T, args []jen.Code) string {
	t.Helper()

	b := bytes.NewBufferString("\n")
	f := jen.NewFile("main")
	f.Add(
		jen.Func().ID("main").Params().Body(
			jen.ID("exampleFunction").Call(args...),
		),
	)
	require.NoError(t, f.Render(b))

	return b.String()
}

func Test_buildFakeVarName(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		exampleInput := "fart"

		expected := "exampleFart"
		actual := buildFakeVarName(exampleInput)

		assert.Equal(t, expected, actual)
	})
}

func Test_ctxParam(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(
			jen.Func().ID("test").Params(
				ctxParam(),
			).Block(),
		)

		var b bytes.Buffer
		assert.NoError(t, out.Render(&b))

		expected := `
package main

import (
	"context"
)

func test(ctx context.Context) {}
`
		actual := "\n" + b.String()

		assert.Equal(t, expected, actual)
	})
}

func Test_ctxVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		out := jen.NewFile("main")
		out.Add(
			jen.Func().ID("test").Params().Block(
				ctxVar(),
			),
		)

		var b bytes.Buffer
		assert.NoError(t, out.Render(&b))

		expected := `
package main

import ()

func test() {
	ctx
}
`
		actual := "\n" + b.String()

		assert.Equal(t, expected, actual)
	})
}

// DataType tests

func TestDataType_OwnedByAUserAtSomeLevel(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToUser: true,
		}
		p := &Project{
			DataTypes: []DataType{
				dt,
			},
		}

		assert.True(t, dt.OwnedByAUserAtSomeLevel(p))
	})

	T.Run("with multi-level ownership", func(t *testing.T) {
		t.Parallel()

		dtA := DataType{
			Name:          wordsmith.FromSingularPascalCase("A"),
			BelongsToUser: true,
		}
		dtB := DataType{
			Name:            wordsmith.FromSingularPascalCase("B"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
		}
		dtC := DataType{
			Name:            wordsmith.FromSingularPascalCase("C"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
		}
		p := &Project{
			DataTypes: []DataType{
				dtA,
				dtB,
				dtC,
			},
		}

		assert.True(t, dtC.OwnedByAUserAtSomeLevel(p))
	})
}

func TestDataType_RestrictedToUserAtSomeLevel(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			BelongsToUser:    true,
			RestrictedToUser: true,
		}
		p := &Project{
			DataTypes: []DataType{
				dt,
			},
		}

		assert.True(t, dt.RestrictedToUserAtSomeLevel(p))
	})

	T.Run("with multi-level ownership", func(t *testing.T) {
		t.Parallel()

		dtA := DataType{
			Name:             wordsmith.FromSingularPascalCase("A"),
			BelongsToUser:    true,
			RestrictedToUser: true,
		}
		dtB := DataType{
			Name:            wordsmith.FromSingularPascalCase("B"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
		}
		dtC := DataType{
			Name:            wordsmith.FromSingularPascalCase("C"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
		}
		p := &Project{
			DataTypes: []DataType{
				dtA,
				dtB,
				dtC,
			},
		}

		assert.True(t, dtC.RestrictedToUserAtSomeLevel(p))
	})
}

func TestDataType_MultipleOwnersBelongingToUser(T *testing.T) {
	T.Parallel()

	T.Run("with multi-level ownership", func(t *testing.T) {
		t.Parallel()

		dtA := DataType{
			Name:          wordsmith.FromSingularPascalCase("A"),
			BelongsToUser: true,
		}
		dtB := DataType{
			Name:            wordsmith.FromSingularPascalCase("B"),
			BelongsToUser:   true,
			BelongsToStruct: wordsmith.FromSingularPascalCase("A"),
		}
		dtC := DataType{
			Name:            wordsmith.FromSingularPascalCase("C"),
			BelongsToUser:   true,
			BelongsToStruct: wordsmith.FromSingularPascalCase("B"),
		}
		p := &Project{
			DataTypes: []DataType{
				dtA,
				dtB,
				dtC,
			},
		}

		assert.True(t, dtC.MultipleOwnersBelongingToUser(p))
	})
}

func TestDataType_buildGetSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("simple", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.buildGetSomethingParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildArchiveSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.buildArchiveSomethingParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionExistenceMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionExistenceMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionArchiveMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionArchiveMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientArchiveMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientArchiveMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientExistenceMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientExistenceMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context, thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierArchiveMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveQueryMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func example(thingID uint64) {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierArchiveQueryMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

// tests pass up to here

func TestDataType_BuildDBQuerierRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierRetrievalQueryMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierExistenceQueryMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_ModifyQueryBuildingStatementWithJoinClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.ModifyQueryBuildingStatementWithJoinClauses(p, jen.Null())

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildJoinClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		expected := ""
		actual := dt.buildJoinClause("table1", "table2", "table3")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_ModifyQueryBuilderWithJoinClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		s := squirrel.Select("fart")

		result := dt.ModifyQueryBuilderWithJoinClauses(p, s)

		expected := ""
		actual, _, _ := result.ToSql()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildDBQuerierSingleInstanceQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildDBQuerierSingleInstanceQueryMethodConditionalClauses(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierExistenceQueryMethodConditionalClauses(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildDBQuerierSingleInstanceQueryMethodQueryBuildingClauses(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierRetrievalQueryMethodConditionalClauses(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalQueryMethodConditionalClauses(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierListRetrievalQueryMethodConditionalClauses(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierExistenceMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildGetSomethingArgs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildArchiveSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.buildArchiveSomethingArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientExistenceMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildDBClientExistenceMethodCallArgs(p)))
	})
}

func TestDataType_BuildDBClientRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildDBClientRetrievalMethodCallArgs(p)))
	})
}

func TestDataType_BuildDBClientArchiveMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBClientArchiveMethodCallArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierExistenceQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierExistenceQueryBuildingArgs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierRetrievalQueryBuildingArgs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBQuerierArchiveQueryBuildingArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionExistenceMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildInterfaceDefinitionExistenceMethodCallArgs(p)))
	})
}

func TestDataType_BuildInterfaceDefinitionRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildInterfaceDefinitionRetrievalMethodCallArgs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionArchiveMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx, thingID)
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildInterfaceDefinitionArchiveMethodCallArgs()))
	})
}

func TestDataType_buildGetSomethingArgsWithExampleVariables(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildGetSomethingArgsWithExampleVariables(p, true)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientRetrievalTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildHTTPClientRetrievalTestCallArgs(p)))

	})
}

func TestDataType_buildSingleInstanceQueryTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.buildSingleInstanceQueryTestCallArgs(p)))
	})
}

func TestDataType_buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildArgsForMethodThatHandlesAnInstanceWithStructsAndUser(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForDBQuerierExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForDBQuerierRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForServiceRouteExistenceCheck(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForServiceRouteExistenceCheck(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildDBQuerierBuildSomethingExistsQueryTestCallArgs(p)))
	})
}

func TestDataType_BuildDBQuerierRetrievalQueryTestCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildDBQuerierRetrievalQueryTestCallArgs(p)))
	})
}

func TestDataType_BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierSomethingExistsQueryBuilderTestPreQueryLines(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierGetSomethingQueryBuilderTestPreQueryLines(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierGetListOfSomethingQueryBuilderTestPreQueryLines(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierCreateSomethingQueryBuilderTestPreQueryLines(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierUpdateSomethingQueryBuilderTestPreQueryLines(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierUpdateSomethingTestPrerequisiteVariables(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierArchiveSomethingTestPrerequisiteVariables(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierArchiveSomethingQueryBuilderTestPreQueryLines(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildGetSomethingLogValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildGetSomethingLogValues(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildGetListOfSomethingLogValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildGetListOfSomethingLogValues(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetListOfSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("while being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildGetListOfSomethingParams(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("without being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildGetListOfSomethingParams(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildMockDataManagerListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierListRetrievalMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierListRetrievalQueryBuildingMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierListRetrievalQueryBuildingMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildCreateSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("while being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildCreateSomethingParams(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("without being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildCreateSomethingParams(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockInterfaceDefinitionCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildMockInterfaceDefinitionCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildInterfaceDefinitionCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBClientCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import (
	"context"
)

func example() {}
`
		actual := renderFunctionParamsToString(t, dt.BuildDBQuerierCreationMethodParams(p))

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationQueryBuildingMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("while being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierCreationQueryBuildingMethodParams(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("without being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierCreationQueryBuildingMethodParams(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildCreateSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.buildCreateSomethingArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockInterfaceDefinitionCreationMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildMockInterfaceDefinitionCreationMethodCallArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationMethodQueryBuildingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBQuerierCreationMethodQueryBuildingArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForDBQuerierTestOfListRetrievalQueryBuilder(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfUpdateQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildArgsForDBQuerierTestOfUpdateQueryBuilder()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfArchiveQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildArgsForDBQuerierTestOfArchiveQueryBuilder()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForDBQuerierTestOfUpdateMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildArgsForDBQuerierTestOfUpdateMethod()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierCreationMethodArgsToUseFromMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBQuerierCreationMethodArgsToUseFromMethodTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsToUseForDBQuerierCreationQueryBuildingTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildArgsToUseForDBQuerierCreationQueryBuildingTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientCreationMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBClientCreationMethodCallArgs()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildUpdateSomethingParams(T *testing.T) {
	T.Parallel()

	T.Run("while being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildUpdateSomethingParams(p, "updated", true)

		assert.Equal(t, expected, actual)
	})

	T.Run("while being models package", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildUpdateSomethingParams(p, "updated", false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBClientUpdateMethodParams(p, "updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierUpdateMethodParams(p, "updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateQueryBuildingMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierUpdateQueryBuildingMethodParams(p, "updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildInterfaceDefinitionUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildInterfaceDefinitionUpdateMethodParams(p, "updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerUpdateMethodParams(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildMockDataManagerUpdateMethodParams(p, "updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildUpdateSomethingArgsWithExampleVars(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildUpdateSomethingArgsWithExampleVars(p, "updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildUpdateSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.buildUpdateSomethingArgs("updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientUpdateMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBClientUpdateMethodCallArgs("updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBQuerierUpdateMethodArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildDBQuerierUpdateMethodArgs("updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerUpdateMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildMockDataManagerUpdateMethodCallArgs("updated")

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildGetListOfSomethingArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildGetListOfSomethingArgs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDBClientListRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildDBClientListRetrievalMethodCallArgs(p)))
	})
}

func TestDataType_BuildDBQuerierListRetrievalMethodArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDBQuerierListRetrievalMethodArgs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildMockDataManagerListRetrievalMethodCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.BuildMockDataManagerListRetrievalMethodCallArgs(p)))

	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWithOwnerStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildVarDeclarationsOfDependentStructsWithOwnerStruct(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildVarDeclarationsOfDependentStructsWithoutUsingOwnerStruct(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientBuildCreationRequestMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForDBQueriersExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForDBQueriersExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForDBQueriersCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForDBQueriersCreationMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildVarDeclarationsOfDependentStructsWhereEachStructIsImportant(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildVarDeclarationsOfDependentStructsWhereOnlySomeStructsAreImportant(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientRetrievalMethodTestDependentObjects(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildHTTPClientRetrievalMethodTestDependentObjects(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientBuildArchiveRequestMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientBuildExistenceRequestMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientBuildRetrievalRequestMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientArchiveMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildDependentObjectsForHTTPClientListRetrievalTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildDependentObjectsForHTTPClientListRetrievalTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientListRetrievalTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientListRetrievalTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientBuildListRetrievalRequestMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildVarDeclarationsOfDependentStructsForUpdateFunction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildVarDeclarationsOfDependentStructsForUpdateFunction(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientUpdateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientBuildUpdateRequestMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildDependentObjectsForHTTPClientCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildDependentObjectsForHTTPClientCreationMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientUpdateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientArchiveMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientListMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientListMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientSearchMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientSearchMethodTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatStringForHTTPClientCreateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatStringForHTTPClientCreateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatCallArgsForHTTPClientRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatCallArgsForHTTPClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientListMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatCallArgsForHTTPClientListMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatCallArgsForHTTPClientCreationMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildFormatCallArgsForHTTPClientUpdateTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildFormatCallArgsForHTTPClientUpdateTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientExistenceRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientExistenceRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientExistenceRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientExistenceRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientExistenceMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientExistenceMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientCreateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientCreateRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientRetrievalRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientRetrievalRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientRetrievalRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientRetrievalRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientRetrievalMethod(T *testing.T) {
	T.Parallel()

	T.Run("as call", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientRetrievalMethod(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("not as call", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientRetrievalMethod(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientCreateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientCreateRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientCreateMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientCreateMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientUpdateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientUpdateRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientUpdateRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientUpdateRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientUpdateMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientUpdateMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientArchiveRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientArchiveRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveRequestBuildingMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientArchiveRequestBuildingMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientArchiveMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientArchiveMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildParamsForMethodThatHandlesAnInstanceWithStructs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildParamsForMethodThatHandlesAnInstanceWithStructs(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientExistenceRequestBuildingMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientRetrievalRequestBuilderMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientArchiveRequestBuildingMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientArchiveMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientArchiveMethodTestURLFormatCall(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientCreationRequestBuildingMethodArgsForTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildHTTPClientCreationRequestBuildingMethodArgsForTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildHTTPClientCreationMethodArgsForTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildHTTPClientCreationMethodArgsForTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildArgsForHTTPClientListRequestMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildArgsForHTTPClientListRequestMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientListRequestMethod(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientListRequestMethod(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildParamsForHTTPClientMethodThatFetchesAList(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildParamsForHTTPClientMethodThatFetchesAList(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildCallArgsForHTTPClientListRetrievalRequestBuildingMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildCallArgsForHTTPClientListRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildCallArgsForHTTPClientUpdateRequestBuildingMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForHTTPClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildCallArgsForHTTPClientUpdateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarDecs(T *testing.T) {
	T.Parallel()

	T.Run("creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarDecs(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("without creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarDecs(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarDecForModifierFuncs(T *testing.T) {
	T.Parallel()

	T.Run("creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarDecForModifierFuncs(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("without creating context", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarDecForModifierFuncs(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBClientRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientCreateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBClientCreateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBClientArchiveMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarDecsForDBQuerierRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarDecsForListFunction(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarDecsForListFunction(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBClientListRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("with filter", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(p, true)

		assert.Equal(t, expected, actual)
	})

	T.Run("without filter", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarsForDBQuerierListRetrievalMethodTest(p, false)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarCallArgsForCreation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarCallArgsForCreation(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_buildRequisiteFakeVarCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		expected := `
package main

import ()

func main() {
	exampleFunction()
}
`
		assert.Equal(t, expected, renderCallArgsToString(t, dt.buildRequisiteFakeVarCallArgs(p)))
	})
}

func TestDataType_buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.buildRequisiteFakeVarCallArgsForServicesThatUseExampleUser(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForServiceExistenceHandlerTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForServiceReadHandlerTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceCreateHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForServiceCreateHandlerTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForServiceUpdateHandlerTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForServiceArchiveHandlerTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForDBClientExistenceMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForDBClientRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForDBClientArchiveMethodTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildExpectedQueryArgsForDBQueriersListRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForDBQueriersListRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildRequisiteFakeVarCallArgsForDBQueriersArchiveMethodTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForDBClientCreationMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildCallArgsForDBClientCreationMethodTest()

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForDBClientListRetrievalMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildCallArgsForDBClientListRetrievalMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildRequisiteVarsForDBClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}
		p := buildExampleTodoListProject()

		var expected []jen.Code
		actual := dt.BuildRequisiteVarsForDBClientUpdateMethodTest(p)

		assert.Equal(t, expected, actual)
	})
}

func TestDataType_BuildCallArgsForDBClientUpdateMethodTest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		dt := DataType{
			Name: buildObligatorySuperPalabra(),
		}

		var expected []jen.Code
		actual := dt.BuildCallArgsForDBClientUpdateMethodTest()

		assert.Equal(t, expected, actual)
	})
}
