package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
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

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("cfg").Op(":=").Op("&").ID("Config").Values(),
		jen.ID("logger").Op(":=").Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
		jen.ID("authService").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "AuthService").Values(),
		jen.ID("usersService").Op(":=").Op("&").Qual(proj.TypesPackage("mock"), "UsersService").Values(),
		jen.ID("dataManager").Op(":=").ID("database").Dot("BuildMockDatabase").Call(),
		jen.Newline(),
		jen.ID("rpm").Op(":=").Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
		jen.ID("rpm").Dot("On").Call(
			jen.Lit("BuildRouteParamIDFetcher"),
			jen.Qual(constants.MockPkg, "IsType").Call(jen.ID("logger")),
			jen.ID("apiClientIDURLParamKey"),
			jen.Lit("API client"),
		).Dot("Return").Call(jen.ID("dummyIDFetcher")),
		jen.ID("rpm").Dot("On").Call(
			jen.Lit("BuildRouteParamIDFetcher"),
			jen.Qual(constants.MockPkg, "IsType").Call(jen.ID("logger")),
			jen.ID("accountIDURLParamKey"),
			jen.Lit("account"),
		).Dot("Return").Call(jen.ID("dummyIDFetcher")),
		jen.ID("rpm").Dot("On").Call(
			jen.Lit("BuildRouteParamIDFetcher"),
			jen.Qual(constants.MockPkg, "IsType").Call(jen.ID("logger")),
			jen.ID("webhookIDURLParamKey"),
			jen.Lit("webhook"),
		).Dot("Return").Call(jen.ID("dummyIDFetcher")),
	}

	for _, typ := range proj.DataTypes {
		bodyLines = append(bodyLines,
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamIDFetcher"),
				jen.Qual(constants.MockPkg, "IsType").Call(jen.ID("logger")),
				jen.IDf("%sIDURLParamKey", typ.Name.UnexportedVarName()),
				jen.Lit(typ.Name.SingularCommonName()),
			).Dot("Return").Call(jen.ID("dummyIDFetcher")),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("s").Op(":=").ID("ProvideService").Callln(
			jen.ID("cfg"),
			jen.ID("logger"),
			jen.ID("authService"),
			jen.ID("usersService"),
			jen.ID("dataManager"),
			jen.ID("rpm"),
			jen.Qual(proj.CapitalismPackage(), "NewMockPaymentManager").Call(),
		),
		jen.Newline(),
		jen.Qual(constants.MockPkg, "AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("authService"),
			jen.ID("usersService"),
			jen.ID("dataManager"),
			jen.ID("rpm"),
		),
		jen.Qual(constants.AssertionLibrary, "NotNil").Call(jen.ID("t"), jen.ID("s")),
	)

	code.Add(
		jen.Func().ID("TestProvideService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	)

	return code
}
