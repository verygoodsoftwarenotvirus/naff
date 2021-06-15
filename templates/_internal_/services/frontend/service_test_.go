package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("dummyIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Body(
			jen.Return().Lit(0)),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("cfg").Op(":=").Op("&").ID("Config").Values(),
			jen.ID("logger").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("authService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "AuthService").Values(),
			jen.ID("usersService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/pkg/types/mock", "UsersService").Values(),
			jen.ID("dataManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
			jen.ID("rpm").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "NewRouteParamManager").Call(),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("logger"),
				jen.ID("apiClientIDURLParamKey"),
				jen.Lit("API client"),
			).Dot("Return").Call(jen.ID("dummyIDFetcher")),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("logger"),
				jen.ID("accountIDURLParamKey"),
				jen.Lit("account"),
			).Dot("Return").Call(jen.ID("dummyIDFetcher")),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("logger"),
				jen.ID("webhookIDURLParamKey"),
				jen.Lit("webhook"),
			).Dot("Return").Call(jen.ID("dummyIDFetcher")),
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.ID("logger"),
				jen.ID("itemIDURLParamKey"),
				jen.Lit("item"),
			).Dot("Return").Call(jen.ID("dummyIDFetcher")),
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
		jen.Newline(),
	)

	return code
}
