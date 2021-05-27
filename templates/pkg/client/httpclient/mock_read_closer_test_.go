package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockReadCloserTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("io", "ReadCloser").Op("=").Parens(jen.Op("*").ID("mockReadCloser")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("mockReadCloser").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("newMockReadCloser returns a new mock io.ReadCloser."),
		jen.Line(),
		jen.Func().ID("newMockReadCloser").Params().Params(jen.Op("*").ID("mockReadCloser")).Body(
			jen.Return().Op("&").ID("mockReadCloser").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadHandler implements the ReadHandler part of our mockReadCloser."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("mockReadCloser")).ID("Read").Params(jen.ID("b").Index().ID("byte")).Params(jen.ID("int"), jen.ID("error")).Body(
			jen.ID("retVals").Op(":=").ID("m").Dot("Called").Call(jen.ID("b")),
			jen.Return().List(jen.ID("retVals").Dot("Int").Call(jen.Lit(0)), jen.ID("retVals").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Close implements the Closer part of our mockReadCloser."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("mockReadCloser")).ID("Close").Params().Params(jen.ID("err").ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call().Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
