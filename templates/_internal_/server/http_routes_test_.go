package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_buildURLVarChunk").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"), jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
				jen.ID("t").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("expected").Assign().Lit("/{things:stuff}"),
				jen.ID("actual").Assign().ID("buildURLVarChunk").Call(jen.Lit("things"), jen.Lit("stuff")),
				jen.Newline(),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			)),
		),
		jen.Newline(),
	)

	return code
}
