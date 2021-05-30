package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("github.com/brianvoe/gofakeit/v5", "Seed").Call(jen.Qual("time", "Now").Call().Dot("UnixNano").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotEmpty").Call(
						jen.ID("t"),
						jen.Parens(jen.Op("&").ID("ErrorResponse").Valuesln()).Dot("Error").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
