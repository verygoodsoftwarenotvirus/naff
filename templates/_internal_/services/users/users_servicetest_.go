package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildBuildTestService(proj)...)
	code.Add(buildTestProvideUsersService(proj)...)

	return code
}

func buildBuildTestService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("expectedUserCount").Assign().Uint64().Call(jen.Lit(123)),
			jen.Line(),
			jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
			jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
				jen.Lit("GetAllUsersCount"),
				jen.Qual(constants.MockPkg, "Anything"),
			).Dot("Return").Call(jen.ID("expectedUserCount"), jen.Nil()),
			jen.Line(),
			jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(),
			jen.Var().ID("ucp").Qual(proj.InternalMetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
				jen.ID("counterName").Qual(proj.InternalMetricsPackage(), "CounterName"),
				jen.ID("description").String(),
			).Params(jen.Qual(proj.InternalMetricsPackage(), "UnitCounter"), jen.Error()).Body(
				jen.Return().List(jen.ID("uc"), jen.Nil()),
			),
			jen.Line(),
			jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
				jen.Qual(proj.InternalConfigPackage(), "AuthSettings").Values(),
				jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
				jen.Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ucp"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
			),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			utils.AssertExpectationsFor("mockDB", "uc"),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideUsersService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideUsersService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsPackage(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsPackage(), "UnitCounter"),
					jen.Error(),
				).Body(
					jen.Return().List(jen.AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
					jen.Qual(proj.InternalConfigPackage(), "AuthSettings").Values(),
					jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
					jen.Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotNil(jen.ID("service"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil userIDFetcher",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsPackage(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsPackage(), "UnitCounter"),
					jen.Error()).Body(
					jen.Return().List(jen.AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
					jen.Qual(proj.InternalConfigPackage(), "AuthSettings").Values(),
					jen.ID("noop").Dot("NewNonOperationalLogger").Call(),
					jen.Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
					jen.Nil(),
					jen.AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error initializing counter",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsPackage(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsPackage(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsPackage(), "UnitCounter"),
					jen.Error(),
				).Body(
					jen.Return().List(
						jen.AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(),
						constants.ObligatoryError(),
					),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
					jen.Qual(proj.InternalConfigPackage(), "AuthSettings").Values(),
					jen.Qual(proj.InternalLoggingPackage(), "NewNonOperationalLogger").Call(),
					jen.Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
					jen.AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(), jen.ID("ucp"), jen.Nil(),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}
