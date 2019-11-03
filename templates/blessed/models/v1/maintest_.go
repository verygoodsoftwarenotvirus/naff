package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainTestDotGo(pkgRoot string, types []models.DataType) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(pkgRoot, types, ret)

	ret.Add(
		jen.Func().ID("TestErrorResponse_Error").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("_").Op("=").Parens(jen.Op("&").ID("ErrorResponse").Values()).Dot("Error").Call(),
			)),
		),
		jen.Line(),
	)
	return ret
}
