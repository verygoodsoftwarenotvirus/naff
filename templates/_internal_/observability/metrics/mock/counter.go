package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func counterDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("metrics").Dot("UnitCounter").Op("=").Parens(jen.Op("*").ID("UnitCounter")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("UnitCounter").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Increment implements our UnitCounter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")).Body(
			jen.ID("m").Dot("Called").Call(jen.ID("ctx"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("IncrementBy implements our UnitCounter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("int64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("val"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Decrement implements our UnitCounter interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("UnitCounter")).ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context")).Body(
			jen.ID("m").Dot("Called").Call(jen.ID("ctx"))),
		jen.Line(),
	)

	return code
}
