package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvideService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("authService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/pkg/types/mock", "AuthService").Valuesln(),
			jen.ID("s").Op(":=").ID("ProvideService").Call(
				jen.ID("cfg"),
				jen.ID("logger"),
				jen.ID("authService"),
			),
			jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
				jen.ID("t"),
				jen.ID("authService"),
			),
			jen.ID("assert").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("s"),
			),
		),
		jen.Newline(),
	)

	return code
}
