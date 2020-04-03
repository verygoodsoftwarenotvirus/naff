package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func implementationTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("apiURLPrefix").Equals().Lit("/api/v1"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2InternalErrorHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
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
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2ResponseErrorHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"obligatory",
				jen.ID("exampleInput").Assign().VarPointer().Qual("gopkg.in/oauth2.v3/errors", "Response").Values(),
				jen.ID("buildTestService").Call(jen.ID("t")).Dot("OAuth2ResponseErrorHandler").Call(jen.ID("exampleInput")),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AuthorizeScopeHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("blah"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("Scopes").MapAssign().Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(",")),
				),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("exampleClient")),
				),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request but with client ID attached",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("blah"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("blargh"), jen.ID("Scopes").MapAssign().Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(",")),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleClient").Dot("ClientID"),
				).Dot("Return").Call(jen.ID("exampleClient"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot("ClientID")),
				),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request and now rows found fetching client info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("blah,flarg"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("blargh"), jen.ID("Scopes").MapAssign().Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(",")),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleClient").Dot("ClientID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot("ClientID")),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request and error fetching client info",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("blah,flarg"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("blargh"), jen.ID("Scopes").MapAssign().Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(",")),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleClient").Dot("ClientID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot("ClientID")),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid scope & client ID but no client",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("blargh"), jen.ID("Scopes").MapAssign().Index().String().Values(),
				),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("exampleClient").Dot("ClientID"),
				).Dot("Return").Call(jen.ID("exampleClient"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot("ClientID")),
				),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("AuthorizeScopeHandler").Call(jen.ID("res"), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserAuthorizationHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(jen.ID("BelongsToUser").MapAssign().Add(utils.FakeUint64Func())),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("BelongsToUser")),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("exampleClient")),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("UserAuthorizationHandler").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("actual"), jen.ID("expected"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without client attached to request",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Assign().VarPointer().Qual(proj.ModelsV1Package(), "User").Values(jen.ID("ID").MapAssign().Add(utils.FakeUint64Func())),
				jen.ID("expected").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleUser").Dot("ID")),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Equals().ID("req").Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Call(jen.ID("req").Dot("Context").Call(), jen.Qual(proj.ModelsV1Package(), "UserKey"), jen.ID("exampleUser")),
				),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("UserAuthorizationHandler").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("actual"), jen.ID("expected"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with no user info attached",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("UserAuthorizationHandler").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertError(jen.Err(), nil),
				utils.AssertEmpty(jen.ID("actual"), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_ClientAuthorizedHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("true"),
				jen.Line(),
				jen.ID("exampleGrant").Assign().Qual("gopkg.in/oauth2.v3", "AuthorizationCode"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("ClientID").MapAssign().Lit("blah"),
					jen.ID("Scopes").MapAssign().Index().String().Values(),
				),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock", "Anything"), jen.ID("stringID"),
				).Dot("Return").Call(jen.ID("exampleClient"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.ID("stringID"), jen.ID("exampleGrant")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with password credentials grant",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("false"),
				jen.ID("exampleGrant").Assign().Qual("gopkg.in/oauth2.v3", "PasswordCredentials"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.Lit("ID"), jen.ID("exampleGrant")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("false"),
				jen.ID("exampleGrant").Assign().Qual("gopkg.in/oauth2.v3", "AuthorizationCode"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()), jen.ID("ClientID").MapAssign().Lit("blah"), jen.ID("Scopes").MapAssign().Index().String().Values(),
				),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.ID("stringID"), jen.ID("exampleGrant")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with disallowed implicit",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("false"),
				jen.Line(),
				jen.ID("exampleGrant").Assign().Qual("gopkg.in/oauth2.v3", "Implicit"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("ClientID").MapAssign().Lit("blah"),
					jen.ID("Scopes").MapAssign().Index().String().Values(),
				),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.ID("exampleClient"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientAuthorizedHandler").Call(jen.ID("stringID"), jen.ID("exampleGrant")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_ClientScopeHandler").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("true"),
				jen.Line(),
				jen.ID("exampleScope").Assign().Lit("halb"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("ClientID").MapAssign().Lit("blah"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.ID("exampleScope")),
				),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.ID("exampleClient"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientScopeHandler").Call(jen.ID("stringID"), jen.ID("exampleScope")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertNoError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("false"),
				jen.Line(),
				jen.ID("exampleScope").Assign().Lit("halb"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("ClientID").MapAssign().Lit("blah"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.ID("exampleScope")),
				),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientScopeHandler").Call(jen.ID("stringID"), jen.ID("exampleScope")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without valid scope",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().ID("false"),
				jen.Line(),
				jen.ID("exampleScope").Assign().Lit("halb"),
				jen.ID("exampleClient").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ID").MapAssign().Add(utils.FakeUint64Func()),
					jen.ID("ClientID").MapAssign().Lit("blah"),
					jen.ID("Scopes").MapAssign().Index().String().Values(),
				),
				jen.ID("stringID").Assign().Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual("github.com/stretchr/testify/mock", "Anything"),
					jen.ID("stringID"),
				).Dot("Return").Call(jen.ID("exampleClient"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ClientScopeHandler").Call(jen.ID("stringID"), jen.ID("exampleScope")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
				utils.AssertError(jen.Err(), nil),
			),
		),
		jen.Line(),
	)
	return ret
}
