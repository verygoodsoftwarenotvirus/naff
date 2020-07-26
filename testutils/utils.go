package testutils

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"
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

// BuildOwnershipChain takes a series of names and returns a slice of datatypes with ownership between them.
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
func BuildOwnershipChain(names ...string) (out []models.DataType) {
	for i, name := range names {
		if i == 0 {
			out = append(out,
				models.DataType{
					Name: wordsmith.FromSingularPascalCase(name),
				},
			)
		} else {
			out = append(out,
				models.DataType{
					Name:            wordsmith.FromSingularPascalCase(name),
					BelongsToStruct: wordsmith.FromSingularPascalCase(names[i-1]),
				},
			)
		}
	}

	return
}
