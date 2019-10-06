package models

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func mainTestDotGo() *jen.File {
	ret := jen.NewFile("$1")
	utils.AddImports(ret)


	ret.Add(jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		utils.ParallelTest(nil),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("_").Op("=").Parens(jen.Op("&").ID("ErrorResponse").Valuesln()).Dot(
				"Error",
			).Call(),
		)),
	),
	)
	return ret
}
