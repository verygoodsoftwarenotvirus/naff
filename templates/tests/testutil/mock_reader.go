package testutil

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockReaderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("MockReadCloser").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Read implements the io.ReadCloser interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockReadCloser")).ID("Read").Params(jen.ID("p").Index().ID("byte")).Params(jen.ID("n").ID("int"), jen.ID("err").ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(jen.ID("p")),
			jen.Return().List(jen.ID("returnValues").Dot("Int").Call(jen.Lit(0)), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Close implements the io.ReadCloser interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("MockReadCloser")).ID("Close").Params().Params(jen.ID("err").ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(),
			jen.Return().ID("returnValues").Dot("Error").Call(jen.Lit(0)),
		),
		jen.Line(),
	)

	return code
}
