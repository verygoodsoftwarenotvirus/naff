package apiclients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTest_randString()...)
	code.Add(buildBuildRequest()...)
	code.Add(buildTest_fetchUserID(proj)...)
	code.Add(buildTestService_ListHandler(proj)...)
	code.Add(buildTestService_CreateHandler(proj)...)
	code.Add(buildTestService_ReadHandler(proj)...)
	code.Add(buildTestService_ArchiveHandler(proj)...)

	return code
}

func buildTest_randString() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_randString").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("actual").Assign().ID("randString").Call(),
				utils.AssertNotEmpty(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildRequest() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("buildRequest").Params(jen.ID("t").PointerTo().Qual("testing", "T")).Params(jen.PointerTo().Qual("net/http", "Request")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.Line(),
			jen.List(jen.ID(constants.RequestVarName), jen.Err()).Assign().Qual("net/http", "NewRequest").Callln(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit("https://verygoodsoftwarenotvirus.ru"),
				jen.Nil(),
			),
			jen.Line(),
			utils.RequireNotNil(jen.ID(constants.RequestVarName), nil),
			utils.AssertNoError(jen.Err(), nil),
			jen.Return().ID(constants.RequestVarName),
		),
		jen.Line(),
	}

	return lines
}

func buildTest_fetchUserID(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("Test_fetchUserID").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "User"),
				jen.Line(),
				jen.Comment("for the service.fetchUserID() call"),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("fetchUserID").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("User")).Dot("ID"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without context value present",
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("expected").Assign().Uint64().Call(jen.Zero()),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("fetchUserID").Call(jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_ListHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_ListHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID(utils.BuildFakeVarName("User")).Assign().Qual(proj.FakeModelsPackage(), "BuildFakeUser").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2ClientList"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.QueryFilter")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2ClientList")),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.Comment("for the service.fetchUserID() call"),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "ed"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no rows returned",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.QueryFilter")),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientList")).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2ClientList")),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "ed"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.QueryFilter")),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2ClientList")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2ClientList"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientsForUser"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("ID"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.QueryFilter")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2ClientList")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2ClientList")),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ListHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "ed"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_CreateHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_CreateHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Callln(
					jen.Lit("GetUserByUsername"), jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Username")).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")),
					jen.Nil(),
				),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("a").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.ID("a").Dot("On").Callln(
					jen.Lit("ValidateLogin"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Password"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("TOTPToken"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
				).Dot("Return").Call(jen.True(), jen.Nil()),
				jen.ID("s").Dot("authenticator").Equals().ID("a"),
				jen.Line(),
				jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("Increment"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("oauth2ClientCounter").Equals().ID("uc"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2Client")),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.ID("creationMiddlewareCtxKey"),
						jen.ID(utils.BuildFakeVarName("Input")),
					),
				),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusCreated"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "a", "uc", "ed"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with missing input",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error getting user",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Callln(
					jen.Lit("GetUserByUsername"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Username"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), "User")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.ID("creationMiddlewareCtxKey"),
						jen.ID(utils.BuildFakeVarName("Input")),
					),
				),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid credentials",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Callln(
					jen.Lit("GetUserByUsername"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Username"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input"))).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("a").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.ID("a").Dot("On").Callln(
					jen.Lit("ValidateLogin"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Password"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("TOTPToken"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
				).Dot("Return").Call(jen.False(), jen.Nil()),
				jen.ID("s").Dot("authenticator").Equals().ID("a"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.ID("creationMiddlewareCtxKey"),
						jen.ID(utils.BuildFakeVarName("Input")),
					),
				),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "a"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error validating password",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Callln(
					jen.Lit("GetUserByUsername"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Username"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("a").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.ID("a").Dot("On").Callln(
					jen.Lit("ValidateLogin"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Password"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("TOTPToken"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
				).Dot("Return").Call(jen.True(), constants.ObligatoryError()),
				jen.ID("s").Dot("authenticator").Equals().ID("a"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.ID("creationMiddlewareCtxKey"),
						jen.ID(utils.BuildFakeVarName("Input")),
					),
				),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "a"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error creating oauth2 client",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Callln(
					jen.Lit("GetUserByUsername"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Username"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("a").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.ID("a").Dot("On").Callln(
					jen.Lit("ValidateLogin"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Password"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("TOTPToken"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
				).Dot("Return").Call(jen.True(), jen.Nil()),
				jen.ID("s").Dot("authenticator").Equals().ID("a"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.ID("creationMiddlewareCtxKey"),
						jen.ID(utils.BuildFakeVarName("Input")),
					),
				),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "a"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				utils.BuildFakeVarWithCustomName(proj, "exampleInput", "OAuth2ClientCreationInputFromClient", jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				jen.Line(),
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("UserDataManager").Dot("On").Callln(
					jen.Lit("GetUserByUsername"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Username"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("User")), jen.Nil()),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("CreateOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("Input")),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("a").Assign().AddressOf().Qual(proj.InternalAuthPackage("mock"), "Authenticator").Values(),
				jen.ID("a").Dot("On").Callln(
					jen.Lit("ValidateLogin"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("HashedPassword"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("Password"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("TwoFactorSecret"),
					jen.ID(utils.BuildFakeVarName("Input")).Dot("TOTPToken"),
					jen.ID(utils.BuildFakeVarName("User")).Dot("Salt"),
				).Dot("Return").Call(jen.True(), jen.Nil()),
				jen.ID("s").Dot("authenticator").Equals().ID("a"),
				jen.Line(),
				jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("Increment"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("oauth2ClientCounter").Equals().ID("uc"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2Client")),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.ID("creationMiddlewareCtxKey"),
						jen.ID(utils.BuildFakeVarName("Input")),
					),
				),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("CreateHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusCreated"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "a", "uc", "ed"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_ReadHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_ReadHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2Client")),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "ed"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no rows found",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching client from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.TypesPackage(), "OAuth2Client")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error encoding response",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("ed").Assign().AddressOf().Qual(proj.InternalEncodingPackage("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Call(
					jen.Lit("EncodeResponse"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.Qual(constants.MockPkg, "AnythingOfType").Call(jen.Lit("*models.OAuth2Client")),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ReadHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "ed"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_ArchiveHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_ArchiveHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildFakeVar(proj, "User"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("uc").Assign().AddressOf().Qual(proj.InternalMetricsPackage("mock"), "UnitCounter").Values(),
				jen.ID("uc").Dot("On").Call(jen.Lit("Decrement"), jen.Qual(constants.MockPkg, "Anything")).Dot("Return").Call(),
				jen.ID("s").Dot("oauth2ClientCounter").Equals().ID("uc"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusNoContent"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB", "uc"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no rows found",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error deleting record",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("BelongsToUser").Equals().ID("exampleUser").Dot("ID"),
				jen.Line(),
				jen.ID("s").Dot("urlClientIDExtractor").Equals().Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()).Body(
					jen.Return().ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabasePackage(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("ArchiveOAuth2Client"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName),
				).Dot("Return").Call(constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.TypesPackage(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("s").Dot("ArchiveHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}
