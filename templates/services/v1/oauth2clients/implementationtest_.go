package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func implementationTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildImplementationTestConstantDefs()...)
	code.Add(buildTestService_OAuth2InternalErrorHandler()...)
	code.Add(buildTestService_OAuth2ResponseErrorHandler()...)
	code.Add(buildTestService_AuthorizeScopeHandler(proj)...)
	code.Add(buildTestService_UserAuthorizationHandler(proj)...)
	code.Add(buildTestService_ClientAuthorizedHandler(proj)...)
	code.Add(buildTestService_ClientScopeHandler(proj)...)

	return code
}

func buildImplementationTestConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("apiURLPrefix").Equals().Lit("/api/v1"),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_OAuth2InternalErrorHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_OAuth2InternalErrorHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Qual("errors", "New").Call(jen.Lit("blah")),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("OAuth2InternalErrorHandler").Call(jen.ID("expected")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual").Dot("Error"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_OAuth2ResponseErrorHandler() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_OAuth2ResponseErrorHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID(utils.BuildFakeVarName("Input")).Assign().AddressOf().Qual("gopkg.in/oauth2.v3/errors", "Response").Values(),
				jen.ID("buildTestService").Call(jen.ID("t")).Dot("OAuth2ResponseErrorHandler").Call(jen.ID(utils.BuildFakeVarName("Input"))),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_AuthorizeScopeHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_AuthorizeScopeHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"),
						jen.ID(utils.BuildFakeVarName("OAuth2Client")),
					),
				),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Qual("fmt", "Sprintf").Call(
					jen.Lit("%s/%s"),
					jen.ID("apiURLPrefix"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes").Index(jen.Zero()),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes").Index(jen.Zero()), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request but with client ID attached",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID")),
				),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Qual("fmt", "Sprintf").Call(
					jen.Lit("%s/%s"),
					jen.ID("apiURLPrefix"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes").Index(jen.Zero()),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertEqual(jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes").Index(jen.Zero()), jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request and now rows found fetching client info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID")),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request and error fetching client info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID")),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid scope & client ID but no client",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Dot("URL").Dot("Path").Equals().Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ClientID")),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID(constants.ResponseVarName).Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_UserAuthorizationHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_UserAuthorizationHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot(constants.UserOwnershipFieldName)),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID(utils.BuildFakeVarName("OAuth2Client"))),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("UserAuthorizationHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("actual"), jen.ID("expected"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "User"),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("User")).Dot("ID")),
				jen.Line(),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(
						jen.ID(constants.RequestVarName).Dot("Context").Call(),
						jen.Qual(proj.ModelsV1Package(), "SessionInfoKey"),
						jen.ID(utils.BuildFakeVarName("User")).Dot("ToSessionInfo").Call(),
					),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("UserAuthorizationHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("actual"), jen.ID("expected"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no user info attached",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(constants.RequestVarName).Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID(constants.ResponseVarName).Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("UserAuthorizationHandler").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_ClientAuthorizedHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_ClientAuthorizedHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Grant")).Assign().Qual("gopkg.in/oauth2.v3", "AuthorizationCode"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"), jen.Qual(constants.MockPkg, "Anything"), jen.ID("stringID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.ID("stringID"), jen.ID(utils.BuildFakeVarName("Grant"))),
				utils.AssertTrue(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with password credentials grant",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(utils.BuildFakeVarName("Grant")).Assign().Qual("gopkg.in/oauth2.v3", "PasswordCredentials"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.Lit("ID"), jen.ID(utils.BuildFakeVarName("Grant"))),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID(utils.BuildFakeVarName("Grant")).Assign().Qual("gopkg.in/oauth2.v3", "AuthorizationCode"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.ID("stringID"), jen.ID(utils.BuildFakeVarName("Grant"))),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with disallowed implicit",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Grant")).Assign().Qual("gopkg.in/oauth2.v3", "Implicit"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.ID("stringID"), jen.ID(utils.BuildFakeVarName("Grant"))),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildTestService_ClientScopeHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("TestService_ClientScopeHandler").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientScopeHandler").Call(
					jen.ID("stringID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes").Index(jen.Zero()),
				),
				utils.AssertTrue(jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), constants.ObligatoryError()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientScopeHandler").Call(
					jen.ID("stringID"),
					jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("Scopes").Index(jen.Zero()),
				),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without valid scope",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID(utils.BuildFakeVarName("Scope")).Assign().Lit("halb"),
				utils.BuildFakeVar(proj, "OAuth2Client"),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID(utils.BuildFakeVarName("OAuth2Client")).Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(constants.MockPkg, "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.ID(utils.BuildFakeVarName("OAuth2Client")), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientScopeHandler").Call(jen.ID("stringID"), jen.ID(utils.BuildFakeVarName("Scope"))),
				utils.AssertFalse(jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
				jen.Line(),
				utils.AssertExpectationsFor("mockDB"),
			),
		),
		jen.Line(),
	}

	return lines
}
