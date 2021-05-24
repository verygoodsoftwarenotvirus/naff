package panicking

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mockDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("MockPanicker").Struct(jen.ID("mock").Dot("Mock")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewMockPanicker produces a production-ready panicker that will actually panic when called."),
		jen.Line(),
		jen.Func().ID("NewMockPanicker").Params().Params(jen.Op("*").ID("MockPanicker")).Body(
			jen.Return().Op("&").ID("MockPanicker").Values()),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Panic satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("MockPanicker")).ID("Panic").Params(jen.ID("msg").Interface()).Body(
			jen.ID("p").Dot("Called").Call(jen.ID("msg"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Panicf satisfies our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("MockPanicker")).ID("Panicf").Params(jen.ID("format").ID("string"), jen.ID("args").Op("...").Interface()).Body(
			jen.ID("p").Dot("Called").Call(
				jen.ID("format"),
				jen.ID("args"),
			)),
		jen.Line(),
	)

	return code
}
