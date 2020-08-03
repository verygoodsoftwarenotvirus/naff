package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsServiceTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildBuildTestService(proj)...)
	code.Add(buildTestProvideOAuth2ClientsService(proj)...)
	code.Add(buildTest_clientStore_GetByID(proj)...)
	code.Add(buildTestService_HandleAuthorizeRequest()...)
	code.Add(buildTestService_HandleTokenRequest()...)

	return code
}

func buildBuildTestService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildTestService").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().ID("Service")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.ID("manager").Assign().Qual("gopkg.in/oauth2.v3/manage", "NewDefaultManager").Call(),
			jen.List(jen.ID("tokenStore"), jen.Err()).Assign().Qual("gopkg.in/oauth2.v3/store", "NewMemoryTokenStore").Call(),
			utils.RequireNoError(jen.Err(), nil),
			jen.ID("manager").Dot("MustTokenStorage").Call(jen.ID("tokenStore"), jen.Err()),
			jen.ID("server").Assign().Qual("gopkg.in/oauth2.v3/server", "NewDefaultServer").Call(jen.ID("manager")),
			jen.Line(),
			jen.ID("service").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID("database").MapAssign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID(constants.LoggerVarName).MapAssign().Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
				jen.ID("encoderDecoder").MapAssign().AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("authenticator").MapAssign().AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
				jen.ID("urlClientIDExtractor").MapAssign().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
				jen.ID("oauth2ClientCounter").MapAssign().AddressOf().Qual(proj.InternalMetricsV1Package("mock"), "UnitCounter").Values(),
				jen.ID("oauth2Handler").MapAssign().ID("server"),
			),
			jen.Line(),
			jen.Return().ID("service"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestProvideOAuth2ClientsService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestProvideOAuth2ClientsService").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(), jen.Nil()),
				jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(
					jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error(),
				).Body(
					jen.Return(jen.Nil(), jen.Nil()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideOAuth2ClientsService").Callln(
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(), jen.ID("mockDB"),
					jen.AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertNotNil(jen.ID("service"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error providing counter",
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Call(
					jen.Lit("GetAllOAuth2Clients"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Index().PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(), jen.Nil()),
				jen.Line(),
				jen.Var().ID("ucp").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider").Equals().Func().Params(
					jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName"),
					jen.ID("description").String(),
				).Params(
					jen.Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
					jen.Error(),
				).Body(
					jen.Return(jen.Nil(), constants.ObligatoryError()),
				),
				jen.Line(),
				jen.List(jen.ID("service"), jen.Err()).Assign().ID("ProvideOAuth2ClientsService").Callln(
					jen.Qual(constants.NoopLoggingPkg, "ProvideNoopLogger").Call(),
					jen.ID("mockDB"),
					jen.AddressOf().Qual(proj.InternalAuthV1Package("mock"), "Authenticator").Values(),
					jen.Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).SingleLineBlock(jen.Return().Zero()),
					jen.AddressOf().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
					jen.ID("ucp"),
				),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("service"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTest_clientStore_GetByID(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_clientStore_GetByID").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.Line(),
				jen.ID("c").Assign().AddressOf().ID("clientStore").Values(jen.ID("database").MapAssign().ID("mockDB")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("c").Dot("GetByID").Call(
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"), jen.ID("actual").Dot("GetID").Call(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no rows",
				jen.ID(utils.BuildFakeVarName("ID")).Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("ID")),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.Line(),
				jen.ID("c").Assign().AddressOf().ID("clientStore").Values(jen.ID("database").MapAssign().ID("mockDB")),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID("c").Dot("GetByID").Call(jen.ID(utils.BuildFakeVarName("ID"))),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID(utils.BuildFakeVarName("ID")).Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("ID")),
				).Dot("Return").Call(
					jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()),
					utils.BuildError(jen.ID("exampleID")),
				),
				jen.Line(),
				jen.ID("c").Assign().AddressOf().ID("clientStore").Values(jen.ID("database").MapAssign().ID("mockDB")),
				jen.List(jen.Underscore(), jen.Err()).Assign().ID("c").Dot("GetByID").Call(jen.ID(utils.BuildFakeVarName("ID"))),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_HandleAuthorizeRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_HandleAuthorizeRequest").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("moah").Assign().AddressOf().ID("mockOAuth2Handler").Values(),
				jen.ID("moah").Dot("On").Callln(
					jen.Lit("HandleAuthorizeRequest"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("moah"),
				jen.List(jen.ID(constants.RequestVarName), jen.ID(constants.ResponseVarName)).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.Line(),
				utils.AssertNoError(jen.ID("s").Dot("HandleAuthorizeRequest").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)), nil),
				jen.Line(),
				utils.AssertExpectationsFor("moah"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_HandleTokenRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_HandleTokenRequest").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("moah").Assign().AddressOf().ID("mockOAuth2Handler").Values(),
				jen.ID("moah").Dot("On").Callln(
					jen.Lit("HandleTokenRequest"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("moah"),
				jen.List(jen.ID(constants.RequestVarName), jen.ID(constants.ResponseVarName)).Assign().List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot("NewRecorder").Call()),
				jen.Line(),
				utils.AssertNoError(jen.ID("s").Dot("HandleTokenRequest").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)), nil),
				jen.Line(),
				utils.AssertExpectationsFor("moah"),
			),
		),
		jen.Line(),
	}

	return lines
}
