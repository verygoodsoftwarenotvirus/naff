package testutils

import (
	"bytes"
	"strings"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/require"
)

func RemoveImportBlock(input string) string {
	var (
		importBlockStartLine,
		importBlockEndLine int
	)
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "import (") {
			importBlockStartLine = i - 1
		}
		if importBlockStartLine != 0 && line == ")" {
			importBlockEndLine = i + 1
			break
		}
	}

	return strings.Join(append(lines[:importBlockStartLine], lines[importBlockEndLine:]...), "\n")
}

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

func RenderOuterStatementToString(t *testing.T, code ...jen.Code) string {
	t.Helper()

	out := jen.NewFile("example")

	out.Add(code...)

	b := bytes.NewBufferString("\n")
	require.NoError(t, out.Render(b))

	return b.String()
}

func RenderFileToString(t *testing.T, f *jen.File) string {
	t.Helper()

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
