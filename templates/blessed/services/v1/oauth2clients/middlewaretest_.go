package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareTestDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Var().Underscore().Qual("net/http", "Handler").Equals().Parens(jen.PointerTo().ID("mockHTTPHandler")).Call(jen.Nil()),
		jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockHTTPHandler").Struct(jen.Qual(utils.MockPkg, "Mock")),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").PointerTo().ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").ParamPointer().Qual("net/http", "Request")).Block(
			jen.ID("m").Dot("Called").Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Callln(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")), jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Nil()),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ServeHTTP"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("expected").Assign().Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput").Valuesln(
					jen.ID("RedirectURI").MapAssign().Lit("https://blah.com"),
				),
				jen.List(jen.ID("bs"), jen.Err()).Assign().Qual("encoding/json", "Marshal").Call(jen.ID("expected")),
				utils.RequireNoError(jen.Err(), nil),
				jen.ID("req").Dot("Body").Equals().Qual("io/ioutil", "NopCloser").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))),
				jen.Line(),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error decoding request",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("ed").Assign().VarPointer().Qual(proj.InternalEncodingV1Package("mock"), "EncoderDecoder").Values(),
				jen.ID("ed").Dot("On").Callln(
					jen.Lit("DecodeRequest"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
					jen.Qual(utils.MockPkg, "Anything"),
				).Dot("Return").Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("encoderDecoder").Equals().ID("ed"),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ServeHTTP"), jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "Anything"),
				),
				jen.Line(),
				jen.ID("h").Assign().ID("s").Dot("CreationInputMiddleware").Call(jen.ID("mh")),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("h").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertEqual(jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Func().ID("TestService_RequestIsAuthenticated").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("THIS IS A FAKE CLIENT ID"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("things")),
				),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ValidationBearerToken"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(jen.VarPointer().Qual("gopkg.in/oauth2.v3/models", "Token").Values(jen.ID("ClientID").MapAssign().ID("expected").Dot("ClientID")), jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("mh"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("expected").Dot("ClientID"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Lit("/api/v1/things"),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ExtractOAuth2ClientFromRequest").Call(jen.ID("req").Dot("Context").Call(), jen.ID("req")),
				jen.Line(),
				utils.AssertNoError(jen.Err(), nil),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error validating token",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ValidationBearerToken"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(jen.Parens(jen.ParamPointer().Qual("gopkg.in/oauth2.v3/models", "Token")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("mh"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ExtractOAuth2ClientFromRequest").Call(jen.ID("req").Dot("Context").Call(), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error fetching from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("THIS IS A FAKE CLIENT ID"),
				),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ValidationBearerToken"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(jen.VarPointer().Qual("gopkg.in/oauth2.v3/models", "Token").Values(jen.ID("ClientID").MapAssign().ID("expected").Dot("ClientID")), jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("mh"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("expected").Dot("ClientID"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ExtractOAuth2ClientFromRequest").Call(jen.ID("req").Dot("Context").Call(), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with invalid scope",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("THIS IS A FAKE CLIENT ID"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("things")),
				),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ValidationBearerToken"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(jen.VarPointer().Qual("gopkg.in/oauth2.v3/models", "Token").Values(jen.ID("ClientID").MapAssign().ID("expected").Dot("ClientID")), jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("mh"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("expected").Dot("ClientID"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Lit("/api/v1/stuff"),
				jen.List(jen.ID("actual"), jen.Err()).Assign().ID("s").Dot("ExtractOAuth2ClientFromRequest").Call(jen.ID("req").Dot("Context").Call(), jen.ID("req")),
				jen.Line(),
				utils.AssertError(jen.Err(), nil),
				utils.AssertNil(jen.ID("actual"), nil),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2TokenAuthenticationMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			jen.Comment("These tests have a lot of overlap to those of ExtractOAuth2ClientFromRequest, which is deliberate"),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("THIS IS A FAKE CLIENT ID"),
					jen.ID("Scopes").MapAssign().Index().String().Values(jen.Lit("things")),
				),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ValidationBearerToken"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(jen.VarPointer().Qual("gopkg.in/oauth2.v3/models", "Token").Values(jen.ID("ClientID").MapAssign().ID("expected").Dot("ClientID")), jen.Nil()),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("mh"),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("expected").Dot("ClientID"),
				).Dot("Return").Call(jen.ID("expected"), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.ID("req").Dot("URL").Dot("Path").Equals().Lit("/api/v1/things"),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.Line(),
				jen.ID("mhh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mhh").Dot("On").Callln(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("s").Dot("OAuth2TokenAuthenticationMiddleware").Call(jen.ID("mhh")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error authenticating request",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mh").Assign().VarPointer().ID("mockOauth2Handler").Values(),
				jen.ID("mh").Dot("On").Callln(
					jen.Lit("ValidationBearerToken"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(jen.Parens(jen.ParamPointer().Qual("gopkg.in/oauth2.v3/models", "Token")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("oauth2Handler").Equals().ID("mh"),
				jen.Line(),
				jen.ID("res").Assign().ID("httptest").Dot("NewRecorder").Call(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")),
				jen.Line(),
				jen.ID("mhh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mhh").Dot("On").Callln(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(),
				jen.Line(),
				jen.ID("s").Dot("OAuth2TokenAuthenticationMiddleware").Call(jen.ID("mhh")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2ClientInfoMiddleware").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mhh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mhh").Dot("On").Callln(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("res"), jen.ID("req")).Assign().List(jen.ID("httptest").Dot("NewRecorder").Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
				jen.ID("q").Assign().Qual("net/url", "Values").Values(),
				jen.ID("q").Dot("Set").Call(jen.ID("oauth2ClientIDURIParamKey"), jen.ID("expected")),
				jen.ID("req").Dot("URL").Dot("RawQuery").Equals().ID("q").Dot("Encode").Call(),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("expected"),
				).Dot("Return").Call(jen.VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Values(), jen.Nil()),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("s").Dot("OAuth2ClientInfoMiddleware").Call(jen.ID("mhh")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertEqual(jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"with error reading from database",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().Lit("blah"),
				jen.Line(),
				jen.ID("mhh").Assign().VarPointer().ID("mockHTTPHandler").Values(),
				jen.ID("mhh").Dot("On").Callln(
					jen.Lit("ServeHTTP"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.Qual(utils.MockPkg, "AnythingOfType").Call(jen.Lit("*http.Request")),
				).Dot("Return").Call(),
				jen.Line(),
				jen.List(jen.ID("res"), jen.ID("req")).Assign().List(jen.ID("httptest").Dot("NewRecorder").Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
				jen.ID("q").Assign().Qual("net/url", "Values").Values(),
				jen.ID("q").Dot("Set").Call(jen.ID("oauth2ClientIDURIParamKey"), jen.ID("expected")),
				jen.ID("req").Dot("URL").Dot("RawQuery").Equals().ID("q").Dot("Encode").Call(),
				jen.Line(),
				jen.ID("mockDB").Assign().Qual(proj.DatabaseV1Package(), "BuildMockDatabase").Call(),
				jen.ID("mockDB").Dot("OAuth2ClientDataManager").Dot("On").Callln(
					jen.Lit("GetOAuth2ClientByClientID"),
					jen.Qual(utils.MockPkg, "Anything"),
					jen.ID("expected"),
				).Dot("Return").Call(jen.Parens(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Call(jen.Nil()), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
				jen.ID("s").Dot("database").Equals().ID("mockDB"),
				jen.Line(),
				jen.ID("s").Dot("OAuth2ClientInfoMiddleware").Call(jen.ID("mhh")).Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
				utils.AssertEqual(jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot("Code"), nil),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Func().ID("TestService_fetchOAuth2ClientFromRequest").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("THIS IS A FAKE CLIENT ID"),
				),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						utils.CtxVar(),
						jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"),
						jen.ID("expected"),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("fetchOAuth2ClientFromRequest").Call(jen.ID("req")),
				utils.AssertEqual(jen.ID("expected"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without value present",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				utils.AssertNil(jen.ID("s").Dot("fetchOAuth2ClientFromRequest").Call(jen.ID("buildRequest").Call(jen.ID("t"))), nil),
			),
			jen.Line(),
		),
	)

	ret.Add(
		jen.Func().ID("TestService_fetchOAuth2ClientIDFromRequest").Params(jen.ID("T").ParamPointer().Qual("testing", "T")).Block(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Line(),
			utils.BuildSubTest(
				"happy path",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.ID("expected").Assign().VarPointer().Qual(proj.ModelsV1Package(), "OAuth2Client").Valuesln(
					jen.ID("ClientID").MapAssign().Lit("THIS IS A FAKE CLIENT ID"),
				),
				jen.Line(),
				jen.ID("req").Assign().ID("buildRequest").Call(jen.ID("t")).Dot("WithContext").Callln(
					jen.Qual("context", "WithValue").Callln(
						utils.CtxVar(),
						jen.ID("clientIDKey"),
						jen.ID("expected").Dot("ClientID"),
					),
				),
				jen.Line(),
				jen.ID("actual").Assign().ID("s").Dot("fetchOAuth2ClientIDFromRequest").Call(jen.ID("req")),
				utils.AssertEqual(jen.ID("expected").Dot("ClientID"), jen.ID("actual"), nil),
			),
			jen.Line(),
			utils.BuildSubTestWithoutContext(
				"without value present",
				jen.ID("s").Assign().ID("buildTestService").Call(jen.ID("t")),
				jen.Qual("github.com/stretchr/testify/assert",
					"Empty",
				).Call(jen.ID("t"), jen.ID("s").Dot(
					"fetchOAuth2ClientIDFromRequest",
				).Call(jen.ID("buildRequest").Call(jen.ID("t")))),
			),
			jen.Line(),
		),
	)
	return ret
}
