package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersServiceTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Block(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("expectedUserCount").Assign().Uint64().Call(jen.Lit(123)),
			jen.Line(),
			jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
			jen.ID("mockDB").Dot("UserDataManager").Dot("On").Call(
				jen.Lit("GetAllUserCount"),
				jen.Qual(utils.MockPkg, "Anything"),
			).Dot("Return").Call(jen.ID("expectedUserCount"), jen.Nil()),
			jen.Line(),
			jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
			jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
				jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
				jen.ID("description").String(),
			).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"), jen.Error()).Block(
				jen.Return().List(jen.ID("uc"), jen.Nil()),
			),
			jen.Line(),
			jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
				jen.Qual(proj.InternalConfigV1Package(), "AuthSettings").Values(),
				jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
				jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ucp"), jen.Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "NewNewsman").Call(jen.Nil(), jen.Nil()),
			),
			utils.RequireNoError(jen.Err(), nil),
			jen.Line(),
			utils.AssertExpectationsFor("mockDB", "uc"),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideUsersService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error(),
				).Block(
					jen.Return().List(jen.AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
					jen.Qual(proj.InternalConfigV1Package(), "AuthSettings").Values(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
					jen.AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotNil(jen.ID("service"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with nil userIDFetcher",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error()).Block(
					jen.Return().List(jen.AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
					jen.Qual(proj.InternalConfigV1Package(), "AuthSettings").Values(),
					jen.ID("noop").Dot("ProvideNoopLogger").Call(),
					jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
					jen.AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Nil(),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
					jen.Nil(),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error initializing counter",
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String()).Params(jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error(),
				).Block(
					jen.Return().List(
						jen.AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
						constants.ObligatoryError(),
					),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideUsersService").Callln(
					jen.Qual(proj.InternalConfigV1Package(), "AuthSettings").Values(),
					jen.Qual(utils.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
					jen.AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(), jen.ID("ucp"), jen.Nil(),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
