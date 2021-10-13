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
		jen.Func().ID("dummyIDFetcher").Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()).Body(
			jen.Return().EmptyString()),
		jen.Newline(),
	)

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.ID("cfg").Assign().AddressOf().ID("Config").Values(),
		jen.ID("logger").Assign().Qual(proj.InternalLoggingPackage(), "NewNoopLogger").Call(),
		jen.ID("authService").Assign().AddressOf().Qual(proj.TypesPackage("mock"), "AuthService").Values(),
		jen.ID("usersService").Assign().AddressOf().Qual(proj.TypesPackage("mock"), "UsersService").Values(),
		jen.ID("dataManager").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
		jen.Newline(),
		jen.ID("rpm").Assign().Qual(proj.RoutingPackage("mock"), "NewRouteParamManager").Call(),
		jen.ID("rpm").Dot("On").Call(
			jen.Lit("BuildRouteParamStringIDFetcher"),
			jen.ID("apiClientIDURLParamKey"),
		).Dot("Return").Call(jen.ID("dummyIDFetcher")),
		jen.ID("rpm").Dot("On").Call(
			jen.Lit("BuildRouteParamStringIDFetcher"),
			jen.ID("accountIDURLParamKey"),
		).Dot("Return").Call(jen.ID("dummyIDFetcher")),
		jen.ID("rpm").Dot("On").Call(
			jen.Lit("BuildRouteParamStringIDFetcher"),
			jen.ID("webhookIDURLParamKey"),
		).Dot("Return").Call(jen.ID("dummyIDFetcher")),
	}

	for _, typ := range proj.DataTypes {
		bodyLines = append(bodyLines,
			jen.ID("rpm").Dot("On").Call(
				jen.Lit("BuildRouteParamStringIDFetcher"),
				jen.IDf("%sIDURLParamKey", typ.Name.UnexportedVarName()),
			).Dot("Return").Call(jen.ID("dummyIDFetcher")),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("s").Assign().ID("ProvideService").Callln(
			jen.ID("cfg"),
			jen.ID("logger"),
			jen.ID("authService"),
			jen.ID("usersService"),
			jen.ID("dataManager"),
			jen.ID("rpm"),
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
		jen.Func().ID("TestProvideService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			bodyLines...,
		),
		jen.Newline(),
	)

	return code
}
