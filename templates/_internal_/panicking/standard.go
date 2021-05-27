package panicking

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func standardDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("NewProductionPanicker produces a production-ready panicker that will actually panic when called."),
		jen.Line(),
		jen.Func().ID("NewProductionPanicker").Params().Params(jen.ID("Panicker")).Body(
			jen.Return().ID("stdLibPanicker").Values(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("stdLibPanicker").Struct(),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("p").ID("stdLibPanicker")).ID("Panic").Params(jen.ID("msg").Interface()).Body(
			jen.ID("panic").Call(jen.ID("msg"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("p").ID("stdLibPanicker")).ID("Panicf").Params(jen.ID("format").ID("string"), jen.ID("args").Op("...").Interface()).Body(
			jen.ID("p").Dot("Panic").Call(jen.Qual("fmt", "Sprintf").Call(
				jen.ID("format"),
				jen.ID("args").Op("..."),
			))),
		jen.Line(),
	)

	return code
}
