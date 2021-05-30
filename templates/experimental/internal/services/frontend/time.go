package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func timeDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("gomentPanicker").Op("=").ID("panicking").Dot("NewProductionPanicker").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("mustGoment").Params(jen.ID("ts").ID("uint64")).Params(jen.Op("*").ID("goment").Dot("Goment")).Body(
			jen.List(jen.ID("g"), jen.ID("err")).Op(":=").ID("goment").Dot("Unix").Call(jen.ID("int64").Call(jen.ID("ts"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("gomentPanicker").Dot("Panic").Call(jen.ID("err"))),
			jen.Return().ID("g"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("relativeTime").Params(jen.ID("ts").ID("uint64")).Params(jen.ID("string")).Body(
			jen.Return().ID("mustGoment").Call(jen.ID("ts")).Dot("FromNow").Call()),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("relativeTimeFromPtr").Params(jen.ID("ts").Op("*").ID("uint64")).Params(jen.ID("string")).Body(
			jen.If(jen.ID("ts").Op("==").ID("nil")).Body(
				jen.Return().Lit("never")),
			jen.Return().ID("relativeTime").Call(jen.Op("*").ID("ts")),
		),
		jen.Line(),
	)

	return code
}
