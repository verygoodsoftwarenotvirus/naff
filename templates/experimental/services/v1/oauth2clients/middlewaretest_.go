package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareTestDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("_").Qual("net/http", "Handler").Op("=").Parens(jen.Op("*").ID("mockHTTPHandler")).Call(jen.ID("nil")),
	jen.Line(),
	)

	ret.Add(
		jen.Type().ID("mockHTTPHandler").Struct(jen.Qual("github.com/stretchr/testify/mock",
		"Mock",
	)),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("m").Op("*").ID("mockHTTPHandler")).ID("ServeHTTP").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
		jen.ID("m").Dot(
			"Called",
		).Call(jen.ID("res"), jen.ID("req")),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_CreationInputMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request")), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
			)),
			jen.ID("h").Op(":=").ID("s").Dot(
				"CreationInputMiddleware",
			).Call(jen.ID("mh")),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("expected").Op(":=").ID("models").Dot(
				"OAuth2ClientCreationInput",
			).Valuesln(
	jen.ID("RedirectURI").Op(":").Lit("https://blah.com")),
			jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("expected")),
			jen.ID("require").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("req").Dot(
				"Body",
			).Op("=").Qual("io/ioutil", "NopCloser").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error decoding request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("ed").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(),
			jen.ID("ed").Dot(
				"On",
			).Call(jen.Lit("DecodeRequest"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request")), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"encoderDecoder",
			).Op("=").ID("ed"),
			jen.ID("mh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
			)),
			jen.ID("h").Op(":=").ID("s").Dot(
				"CreationInputMiddleware",
			).Call(jen.ID("mh")),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("h").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusBadRequest"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_RequestIsAuthenticated").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("THIS IS A FAKE CLIENT ID"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("things"))),
			jen.ID("mh").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ValidationBearerToken"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(jen.Op("&").Qual("gopkg.in/oauth2.v3/models", "Token").Valuesln(
	jen.ID("ClientID").Op(":").ID("expected").Dot(
				"ClientID",
			)), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("mh"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.ID("expected").Dot(
				"ClientID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/api/v1/things"),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"ExtractOAuth2ClientFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot("NoError").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error validating token"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mh").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ValidationBearerToken"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").Qual("gopkg.in/oauth2.v3/models", "Token")).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("mh"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"ExtractOAuth2ClientFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("THIS IS A FAKE CLIENT ID")),
			jen.ID("mh").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ValidationBearerToken"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(jen.Op("&").Qual("gopkg.in/oauth2.v3/models", "Token").Valuesln(
	jen.ID("ClientID").Op(":").ID("expected").Dot(
				"ClientID",
			)), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("mh"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.ID("expected").Dot(
				"ClientID",
			)).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"ExtractOAuth2ClientFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with invalid scope"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("THIS IS A FAKE CLIENT ID"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("things"))),
			jen.ID("mh").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ValidationBearerToken"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(jen.Op("&").Qual("gopkg.in/oauth2.v3/models", "Token").Valuesln(
	jen.ID("ClientID").Op(":").ID("expected").Dot(
				"ClientID",
			)), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("mh"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.ID("expected").Dot(
				"ClientID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/api/v1/stuff"),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot(
				"ExtractOAuth2ClientFromRequest",
			).Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.ID("req")),
			jen.ID("assert").Dot("Error").Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("actual")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2TokenAuthenticationMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("THIS IS A FAKE CLIENT ID"), jen.ID("Scopes").Op(":").Index().ID("string").Valuesln(
	jen.Lit("things"))),
			jen.ID("mh").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ValidationBearerToken"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(jen.Op("&").Qual("gopkg.in/oauth2.v3/models", "Token").Valuesln(
	jen.ID("ClientID").Op(":").ID("expected").Dot(
				"ClientID",
			)), jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("mh"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.ID("expected").Dot(
				"ClientID",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
			).Op("=").Lit("/api/v1/things"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("mhh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
			jen.ID("mhh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"OAuth2TokenAuthenticationMiddleware",
			).Call(jen.ID("mhh")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error authenticating request"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("mh").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("mh").Dot(
				"On",
			).Call(jen.Lit("ValidationBearerToken"), jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").Qual("gopkg.in/oauth2.v3/models", "Token")).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("mh"),
			jen.ID("res").Op(":=").ID("httptest").Dot(
				"NewRecorder",
			).Call(),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")),
			jen.ID("mhh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
			jen.ID("mhh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(),
			jen.ID("s").Dot(
				"OAuth2TokenAuthenticationMiddleware",
			).Call(jen.ID("mhh")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_OAuth2ClientInfoMiddleware").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Lit("blah"),
			jen.ID("mhh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
			jen.ID("mhh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("q").Op(":=").Qual("net/url", "Values").Values(),
			jen.ID("q").Dot("Set").Call(jen.ID("oauth2ClientIDURIParamKey"), jen.ID("expected")),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"RawQuery",
			).Op("=").ID("q").Dot(
				"Encode",
			).Call(),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.ID("expected")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Values(), jen.ID("nil")),
			jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
			jen.ID("s").Dot(
				"OAuth2ClientInfoMiddleware",
			).Call(jen.ID("mhh")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusOK"), jen.ID("res").Dot(
				"Code",
			)),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Lit("blah"),
			jen.ID("mhh").Op(":=").Op("&").ID("mockHTTPHandler").Values(),
			jen.ID("mhh").Dot(
				"On",
			).Call(jen.Lit("ServeHTTP"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.Qual("github.com/stretchr/testify/mock",
				"AnythingOfType",
			).Call(jen.Lit("*http.Request"))).Dot(
				"Return",
			).Call(),
			jen.List(jen.ID("res"), jen.ID("req")).Op(":=").List(jen.ID("httptest").Dot(
				"NewRecorder",
			).Call(), jen.ID("buildRequest").Call(jen.ID("t"))),
			jen.ID("q").Op(":=").Qual("net/url", "Values").Values(),
			jen.ID("q").Dot("Set").Call(jen.ID("oauth2ClientIDURIParamKey"), jen.ID("expected")),
			jen.ID("req").Dot(
				"URL",
			).Dot(
				"RawQuery",
			).Op("=").ID("q").Dot(
				"Encode",
			).Call(),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.Qual("github.com/stretchr/testify/mock",
				"Anything",
	),
	jen.ID("expected")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("s").Dot("Database").Op("=").ID("mockDB"),
			jen.ID("s").Dot(
				"OAuth2ClientInfoMiddleware",
			).Call(jen.ID("mhh")).Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.Qual("net/http", "StatusUnauthorized"), jen.ID("res").Dot(
				"Code",
			)),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_fetchOAuth2ClientFromRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("THIS IS A FAKE CLIENT ID")),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")).Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.Qual("context", "Background").Call(), jen.ID("models").Dot(
				"OAuth2ClientKey",
	),
	jen.ID("expected"))),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"fetchOAuth2ClientFromRequest",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected"), jen.ID("actual")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without value present"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("s").Dot(
				"fetchOAuth2ClientFromRequest",
			).Call(jen.ID("buildRequest").Call(jen.ID("t")))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_fetchOAuth2ClientIDFromRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("expected").Op(":=").Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(
	jen.ID("ClientID").Op(":").Lit("THIS IS A FAKE CLIENT ID")),
			jen.ID("req").Op(":=").ID("buildRequest").Call(jen.ID("t")).Dot(
				"WithContext",
			).Call(jen.Qual("context", "WithValue").Call(jen.Qual("context", "Background").Call(), jen.ID("clientIDKey"), jen.ID("expected").Dot(
				"ClientID",
			))),
			jen.ID("actual").Op(":=").ID("s").Dot(
				"fetchOAuth2ClientIDFromRequest",
			).Call(jen.ID("req")),
			jen.ID("assert").Dot("Equal").Call(jen.ID("t"), jen.ID("expected").Dot(
				"ClientID",
	),
	jen.ID("actual")),
		)),
			jen.Line(),
			jen.ID("T").Dot("Run").Call(jen.Lit("without value present"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("assert").Dot(
				"Empty",
			).Call(jen.ID("t"), jen.ID("s").Dot(
				"fetchOAuth2ClientIDFromRequest",
			).Call(jen.ID("buildRequest").Call(jen.ID("t")))),
		)),
	),
	jen.Line(),
	)
	return ret
}
