package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func implementationTestDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("apiURLPrefix").Op("=").Lit("/api/v1"),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2InternalErrorHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Qual("errors", "New").Call(jen.Lit("blah")),
				jen.ID("actual").Op(":=").ID("s").Dot(
					"OAuth2InternalErrorHandler",
				).Call(jen.ID("expected")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual").Dot("Error")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2ResponseErrorHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("obligatory"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("exampleInput").Op(":=").Op("&").Qual("gopkg.in/oauth2.v3/errors", "Response").Values(),
				jen.ID("buildTestService").Call(jen.ID("t")).Dot(
					"OAuth2ResponseErrorHandler",
				).Call(jen.ID("exampleInput")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_AuthorizeScopeHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("blah"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("Scopes").Op(":").Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(","))),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("models").Dot(
					"OAuth2ClientKey",
				),
					jen.ID("exampleClient"))),
				jen.ID("req").Dot(
					"URL",
				).Dot(
					"Path",
				).Op("=").Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"AuthorizeScopeHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
					"Code",
				)),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without client attached to request but with client ID attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("blah"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ClientID").Op(":").Lit("blargh"), jen.ID("Scopes").Op(":").Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(","))),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleClient").Dot(
						"ClientID",
					)).Dot(
					"Return",
				).Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot(
					"ClientID",
				))),
				jen.ID("req").Dot(
					"URL",
				).Dot(
					"Path",
				).Op("=").Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"AuthorizeScopeHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
					"Code",
				)),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without client attached to request and now rows found fetching client info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("blah,flarg"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ClientID").Op(":").Lit("blargh"), jen.ID("Scopes").Op(":").Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(","))),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleClient").Dot(
						"ClientID",
					)).Dot(
					"Return",
				).Call(jen.Parens(jen.Op("*").ID("models").Dot(
					"OAuth2Client",
				)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot(
					"ClientID",
				))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"AuthorizeScopeHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusNotFound"), jen.ID("res").Dot(
					"Code",
				)),
				jen.ID("assert").Dot(
					"Empty",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without client attached to request and error fetching client info"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").Lit("blah,flarg"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ClientID").Op(":").Lit("blargh"), jen.ID("Scopes").Op(":").Qual("strings", "Split").Call(jen.ID("expected"), jen.Lit(","))),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleClient").Dot(
						"ClientID",
					)).Dot(
					"Return",
				).Call(jen.Parens(jen.Op("*").ID("models").Dot(
					"OAuth2Client",
				)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot(
					"ClientID",
				))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"AuthorizeScopeHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusInternalServerError"), jen.ID("res").Dot(
					"Code",
				)),
				jen.ID("assert").Dot(
					"Empty",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without client attached to request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"AuthorizeScopeHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot(
					"Code",
				)),
				jen.ID("assert").Dot(
					"Empty",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid scope & client ID but no client"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ClientID").Op(":").Lit("blargh"), jen.ID("Scopes").Op(":").Index().ID("string").Values()),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("exampleClient").Dot(
						"ClientID",
					)).Dot(
					"Return",
				).Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("clientIDKey"), jen.ID("exampleClient").Dot(
					"ClientID",
				))),
				jen.ID("req").Dot(
					"URL",
				).Dot(
					"Path",
				).Op("=").Qual("fmt", "Sprintf").Call(jen.Lit("%s/blah"), jen.ID("apiURLPrefix")),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"AuthorizeScopeHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
					"Code",
				)),
				jen.ID("assert").Dot(
					"Empty",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_UserAuthorizationHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("BelongsTo").Op(":").Lit(1)),
				jen.ID("expected").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("BelongsTo")),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("models").Dot(
					"OAuth2ClientKey",
				),
					jen.ID("exampleClient"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"UserAuthorizationHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("actual"), jen.ID("expected")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without client attached to request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("exampleUser").Op(":=").Op("&").ID("models").Dot(
					"User",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1)),
				jen.ID("expected").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleUser").Dot("ID")),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.ID("models").Dot(
					"UserKey",
				),
					jen.ID("exampleUser"))),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"UserAuthorizationHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("actual"), jen.ID("expected")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with no user info attached"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Op(":=").ID("httptest").Dot(
					"NewRecorder",
				).Call(),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"UserAuthorizationHandler",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
				jen.ID("assert").Dot(
					"Empty",
				).Call(jen.ID("t"), jen.ID("actual")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_ClientAuthorizedHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("true"),
				jen.ID("exampleGrant").Op(":=").ID("oauth2").Dot(
					"AuthorizationCode",
				),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1), jen.ID("ClientID").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").Index().ID("string").Values()),
				jen.ID("stringID").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("stringID")).Dot(
					"Return",
				).Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientAuthorizedHandler",
				).Call(jen.ID("stringID"), jen.ID("exampleGrant")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with password credentials grant"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("false"),
				jen.ID("exampleGrant").Op(":=").ID("oauth2").Dot(
					"PasswordCredentials",
				),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientAuthorizedHandler",
				).Call(jen.Lit("ID"), jen.ID("exampleGrant")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("false"),
				jen.ID("exampleGrant").Op(":=").ID("oauth2").Dot(
					"AuthorizationCode",
				),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1), jen.ID("ClientID").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").Index().ID("string").Values()),
				jen.ID("stringID").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("stringID")).Dot(
					"Return",
				).Call(jen.Parens(jen.Op("*").ID("models").Dot(
					"OAuth2Client",
				)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientAuthorizedHandler",
				).Call(jen.ID("stringID"), jen.ID("exampleGrant")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with disallowed implicit"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("false"),
				jen.ID("exampleGrant").Op(":=").ID("oauth2").Dot(
					"Implicit",
				),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1),
					jen.ID("ClientID").Op(":").Lit("blah"),
					jen.ID("Scopes").Op(":").Index().ID("string").Values(),
				),
				jen.ID("stringID").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("stringID")).Dot(
					"Return",
				).Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientAuthorizedHandler",
				).Call(jen.ID("stringID"), jen.ID("exampleGrant")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_ClientScopeHandler").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("true"),
				jen.ID("exampleScope").Op(":=").Lit("halb"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1), jen.ID("ClientID").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
						jen.ID("exampleScope"))),
				jen.ID("stringID").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("stringID")).Dot(
					"Return",
				).Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientScopeHandler",
				).Call(jen.ID("stringID"), jen.ID("exampleScope")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("false"),
				jen.ID("exampleScope").Op(":=").Lit("halb"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1), jen.ID("ClientID").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
						jen.ID("exampleScope"))),
				jen.ID("stringID").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("stringID")).Dot(
					"Return",
				).Call(jen.Parens(jen.Op("*").ID("models").Dot(
					"OAuth2Client",
				)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientScopeHandler",
				).Call(jen.ID("stringID"), jen.ID("exampleScope")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without valid scope"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
				jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Op(":=").ID("false"),
				jen.ID("exampleScope").Op(":=").Lit("halb"),
				jen.ID("exampleClient").Op(":=").Op("&").ID("models").Dot(
					"OAuth2Client",
				).Valuesln(
					jen.ID("ID").Op(":").Lit(1), jen.ID("ClientID").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").Index().ID("string").Values()),
				jen.ID("stringID").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%d"), jen.ID("exampleClient").Dot("ID")),
				jen.ID("mockDB").Op(":=").ID("database").Dot(
					"BuildMockDatabase",
				).Call(),
				jen.ID("mockDB").Dot(
					"OAuth2ClientDataManager",
				).Dot("On").Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
					"Anything",
				),
					jen.ID("stringID")).Dot(
					"Return",
				).Call(jen.ID("exampleClient"), jen.ID("nil")),
				jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
				jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ClientScopeHandler",
				).Call(jen.ID("stringID"), jen.ID("exampleScope")),
				jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
				jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			)),
		),
		jen.Line(),
	)
	return ret
}
