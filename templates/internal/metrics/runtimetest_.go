package metrics

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func runtimeTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestRecordRuntimeStats()...)
	code.Add(buildTestRegisterDefaultViews()...)

	return code
}

func buildTestRecordRuntimeStats() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestRecordRuntimeStats").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("this is sort of an obligatory test for coverage's sake."),
			jen.Line(),
			jen.ID("d").Assign().Qual("time", "Second"),
			jen.ID("sf").Assign().ID("RecordRuntimeStats").Call(jen.ID("d").Op("/").Lit(5)),
			jen.Qual("time", "Sleep").Call(jen.ID("d")),
			jen.ID("sf").Call(),
		),
		jen.Line(),
	}

	return lines
}

func buildTestRegisterDefaultViews() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestRegisterDefaultViews").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("obligatory"),
			utils.RequireNoError(jen.ID("RegisterDefaultViews").Call(), nil),
		),
		jen.Line(),
	}

	return lines
}
