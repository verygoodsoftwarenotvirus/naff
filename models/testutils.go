package models

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"github.com/stretchr/testify/require"
)

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

// buildOwnershipChain takes a series of names and returns a slice of datatypes with ownership between them.
// So for instance, if you provided `Forum`, `Subforum`, and `Post` as input, the output would be:
// 		[]DataType{
//			{
//				Name: wordsmith.FromSingularPascalCase("Forum"),
//			},
//			{
//				Name:            wordsmith.FromSingularPascalCase("Subforum"),
//				BelongsToStruct: wordsmith.FromSingularPascalCase("Forum"),
//			},
//			{
//				Name:            wordsmith.FromSingularPascalCase("Post"),
//				BelongsToStruct: wordsmith.FromSingularPascalCase("Subforum"),
//			},
//		}
func buildOwnershipChain(names ...string) (out []DataType) {
	for i, name := range names {
		if i == 0 {
			out = append(out,
				DataType{
					Name: wordsmith.FromSingularPascalCase(name),
				},
			)
		} else {
			out = append(out,
				DataType{
					Name:            wordsmith.FromSingularPascalCase(name),
					BelongsToStruct: wordsmith.FromSingularPascalCase(names[i-1]),
				},
			)
		}
	}

	return
}

func renderIndependentStatementToString(t *testing.T, result *jen.Statement) string {
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
