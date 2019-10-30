package auth

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func authenticatorTestDotGo(rootPkg string) *jen.File {
	ret := jen.NewFile("auth_test")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("TestProvideBcryptHashCost").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.Qual(filepath.Join(rootPkg, "internal/v1/auth"), "ProvideBcryptHashCost").Call(),
			)),
		),
		jen.Line(),
	)
	return ret
}
