package testutils

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/require"
)

func RenderFunctionParamsToString(t *testing.T, params []jen.Code) string {
	t.Helper()

	f := jen.NewFile("main")
	f.Add(jen.Func().ID("example").Params(params...).Block())
	b := bytes.NewBufferString("\n")
	require.NoError(t, f.Render(b))

	return b.String()
}

func RenderCallArgsToString(t *testing.T, args []jen.Code) string {
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

//// buildOwnershipChain takes a series of names and returns a slice of datatypes with ownership between them.
//// So for instance, if you provided `Forum`, `Subforum`, and `Post` as input, the output would be:
//// 		[]DataType{
////			{
////				Name: wordsmith.FromSingularPascalCase("Forum"),
////			},
////			{
////				Name:            wordsmith.FromSingularPascalCase("Subforum"),
////				BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
////			},
////			{
////				Name:            wordsmith.FromSingularPascalCase("Post"),
////				BelongsToStruct: wordsmith.FromSingularPascalCase("Subforum"),
////			},
////		}
//func BuildOwnershipChain(names ...string) (out []models.DataType) {
//	for i, name := range names {
//		if i == 0 {
//			out = append(out,
//				models.DataType{
//					Name: wordsmith.FromSingularPascalCase(name),
//				},
//			)
//		} else {
//			out = append(out,
//				models.DataType{
//					Name:            wordsmith.FromSingularPascalCase(name),
//					BelongsToStruct: wordsmith.FromSingularPascalCase(names[i-1]),
//				},
//			)
//		}
//	}
//
//	return
//}

func RenderIndependentStatementToString(t *testing.T, result jen.Code) string {
	t.Helper()

	f := jen.NewFile("main")
	f.Add(
		jen.Func().ID("main").Params().Body(
			result,
		),
	)
	b := bytes.NewBufferString("\n")
	require.NoError(t, f.Render(b))

	return b.String()
}

func RenderMapEntriesWithStringKeysToString(t *testing.T, values []jen.Code) string {
	t.Helper()

	f := jen.NewFile("main")
	f.Add(
		jen.Func().ID("main").Params().Body(
			jen.ID("exampleMap").Assign().Map(jen.String()).Interface().Valuesln(
				values...,
			),
		),
	)
	b := bytes.NewBufferString("\n")
	require.NoError(t, f.Render(b))

	return b.String()
}

func RenderVariableDeclarationsToString(t *testing.T, vars []jen.Code) string {
	t.Helper()

	f := jen.NewFile("main")
	f.Add(
		jen.Func().ID("main").Params().Body(
			vars...,
		),
	)
	b := bytes.NewBufferString("\n")
	require.NoError(t, f.Render(b))

	return b.String()
}
