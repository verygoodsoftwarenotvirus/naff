package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func callerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("runtimeFrameBuffer").Op("=").Lit(3).Var().ID("counterBuffer").Op("=").Lit(2).Var().ID("this").Op("=").Lit("gitlab.com/verygoodsoftwarenotvirus/todo/"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("GetCallerName is largely (and respectfully) inspired by/copied from https://stackoverflow.com/a/35213181").ID("GetCallerName").Params().Params(jen.ID("string")).Body(
			jen.ID("programCounters").Op(":=").ID("make").Call(
				jen.Index().ID("uintptr"),
				jen.ID("runtimeFrameBuffer").Op("+").ID("counterBuffer"),
			),
			jen.ID("n").Op(":=").Qual("runtime", "Callers").Call(
				jen.Lit(0),
				jen.ID("programCounters"),
			),
			jen.ID("frame").Op(":=").Qual("runtime", "Frame").Valuesln(jen.ID("Function").Op(":").Lit("unknown")),
			jen.If(jen.ID("n").Op(">").Lit(0)).Body(
				jen.ID("frames").Op(":=").Qual("runtime", "CallersFrames").Call(jen.ID("programCounters").Index(jen.Empty(), jen.ID("n"))),
				jen.For(jen.List(jen.ID("more"), jen.ID("frameIndex")).Op(":=").List(jen.ID("true"), jen.Lit(0)), jen.ID("more").Op("&&").ID("frameIndex").Op("<=").ID("runtimeFrameBuffer"), jen.ID("frameIndex").Op("++")).Body(
					jen.If(jen.ID("frameIndex").Op("==").ID("runtimeFrameBuffer")).Body(
						jen.List(jen.ID("frame"), jen.ID("more")).Op("=").ID("frames").Dot("Next").Call()).Else().Body(
						jen.List(jen.ID("_"), jen.ID("more")).Op("=").ID("frames").Dot("Next").Call())),
			),
			jen.Return().Qual("strings", "TrimPrefix").Call(
				jen.ID("frame").Dot("Function"),
				jen.ID("this"),
			),
		),
		jen.Line(),
	)

	return code
}
