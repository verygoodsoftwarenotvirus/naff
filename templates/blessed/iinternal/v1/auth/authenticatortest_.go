package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authenticatorTestDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("auth_test")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("TestProvideBcryptHashCost").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").ParamPointer().Qual("testing", "T")).Block(
				jen.Qual(pkg.InternalAuthV1Package(), "ProvideBcryptHashCost").Call(),
			)),
		),
		jen.Line(),
	)
	return ret
}
