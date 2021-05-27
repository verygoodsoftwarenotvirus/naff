package metrics

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
			jen.ID("_").ID("UnitCounter").Op("=").Parens(jen.Op("*").ID("unitCounter")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("unitCounter").Struct(jen.ID("counter").ID("metric").Dot("Int64Counter")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("unitCounter")).ID("Increment").Params(jen.ID("ctx").Qual("context", "Context")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("counter").Dot("Add").Call(
				jen.ID("ctx"),
				jen.Lit(1),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("unitCounter")).ID("IncrementBy").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("val").ID("int64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("counter").Dot("Add").Call(
				jen.ID("ctx"),
				jen.ID("val"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("c").Op("*").ID("unitCounter")).ID("Decrement").Params(jen.ID("ctx").Qual("context", "Context")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("c").Dot("counter").Dot("Add").Call(
				jen.ID("ctx"),
				jen.Op("-").Lit(1),
			),
		),
		jen.Line(),
	)

	return code
}
