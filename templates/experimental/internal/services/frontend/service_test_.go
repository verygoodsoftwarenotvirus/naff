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
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("authService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuthService").Valuesln(),
			jen.ID("usersService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UsersService").Valuesln(),
			jen.ID("dataManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
			jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
			jen.ID("s").Op(":=").ID("ProvideService").Call(
				jen.ID("cfg"),
				jen.ID("logger"),
				jen.ID("authService"),
				jen.ID("usersService"),
				jen.ID("dataManager"),
				jen.ID("rpm"),
				jen.ID("capitalism").Dot("NewMockPaymentManager").Call(),
			),
			jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
				jen.ID("t"),
				jen.ID("authService"),
				jen.ID("usersService"),
				jen.ID("dataManager"),
				jen.ID("rpm"),
			),
			jen.Return().ID("s").Assert(jen.Op("*").ID("service")),
		),
		jen.Line(),
	)

	return code
}
