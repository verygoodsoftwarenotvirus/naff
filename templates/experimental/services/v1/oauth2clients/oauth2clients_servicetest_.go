package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func oauth2ClientsServiceTestDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(ret)

	ret.Add(
		jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Service")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("manager").Op(":=").ID("manage").Dot(
			"NewDefaultManager",
		).Call(),
		jen.List(jen.ID("tokenStore"), jen.ID("err")).Op(":=").Qual("gopkg.in/oauth2.v3/store", "NewMemoryTokenStore").Call(),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.ID("manager").Dot(
			"MustTokenStorage",
		).Call(jen.ID("tokenStore"), jen.ID("err")),
		jen.ID("server").Op(":=").Qual("gopkg.in/oauth2.v3/server", "NewDefaultServer").Call(jen.ID("manager")),
		jen.ID("service").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("database").Op(":").ID("database").Dot(
			"BuildMockDatabase",
		).Call(), jen.ID("logger").Op(":").ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("encoderDecoder").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("authenticator").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(), jen.ID("urlClientIDExtractor").Op(":").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().Lit(0),
	),
	jen.ID("oauth2ClientCounter").Op(":").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(), jen.ID("tokenStore").Op(":").ID("tokenStore"), jen.ID("oauth2Handler").Op(":").ID("server")),
		jen.Return().ID("service"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestProvideOAuth2ClientsService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Index().Op("*").ID("models").Dot(
				"OAuth2Client",
			).Values(), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllOAuth2ClientCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot(
				"Return",
			).Call(),

		jen.Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
	),
	jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
	),
	jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideOAuth2ClientsService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
	),
	jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("ucp")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
		jen.ID("T").Dot("Run").Call(jen.Lit("with error providing counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Index().Op("*").ID("models").Dot(
				"OAuth2Client",
			).Values(), jen.ID("nil")),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllOAuth2ClientCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.ID("nil")),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot(
				"Return",
			).Call(),

		jen.Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
	),
	jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
	),
	jen.ID("error")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideOAuth2ClientsService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
	),
	jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("ucp")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
		jen.ID("T").Dot("Run").Call(jen.Lit("with error fetching oauth2 clients"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllOAuth2Clients"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.Index().Op("*").ID("models").Dot(
				"OAuth2Client",
			).Values(), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetAllOAuth2ClientCount"), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("expected"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "UnitCounter").Values(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("expected")).Dot(
				"Return",
			).Call(),

		jen.Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
	),
	jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
	),
	jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideOAuth2ClientsService").Call(jen.Qual("context", "Background").Call(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "Authenticator").Values(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
	),
	jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "EncoderDecoder").Values(), jen.ID("ucp")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("Test_clientStore_GetByID").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleID").Op(":=").Lit("blah"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.ID("mock").Dot(
				"Anything",
	),
	jen.ID("exampleID")).Dot(
				"Return",
			).Call(jen.Op("&").ID("models").Dot(
				"OAuth2Client",
			).Valuesln(jen.ID("ClientID").Op(":").ID("exampleID")), jen.ID("nil")),
			jen.ID("c").Op(":=").Op("&").ID("clientStore").Valuesln(jen.ID("database").Op(":").ID("mockDB")),
			jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetByID",
			).Call(jen.ID("exampleID")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Equal",
			).Call(jen.ID("t"), jen.ID("exampleID"), jen.ID("actual").Dot(
				"GetID",
			).Call()),
		)),
		jen.ID("T").Dot("Run").Call(jen.Lit("with no rows"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleID").Op(":=").Lit("blah"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.ID("mock").Dot(
				"Anything",
	),
	jen.ID("exampleID")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("database/sql", "ErrNoRows")),
			jen.ID("c").Op(":=").Op("&").ID("clientStore").Valuesln(jen.ID("database").Op(":").ID("mockDB")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetByID",
			).Call(jen.ID("exampleID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
		jen.ID("T").Dot("Run").Call(jen.Lit("with error reading from database"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("exampleID").Op(":=").Lit("blah"),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"OAuth2ClientDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetOAuth2ClientByClientID"), jen.ID("mock").Dot(
				"Anything",
	),
	jen.ID("exampleID")).Dot(
				"Return",
			).Call(jen.Parens(jen.Op("*").ID("models").Dot(
				"OAuth2Client",
			)).Call(jen.ID("nil")), jen.Qual("errors", "New").Call(jen.ID("exampleID"))),
			jen.ID("c").Op(":=").Op("&").ID("clientStore").Valuesln(jen.ID("database").Op(":").ID("mockDB")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("c").Dot(
				"GetByID",
			).Call(jen.ID("exampleID")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_HandleAuthorizeRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("moah").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("moah").Dot(
				"On",
			).Call(jen.Lit("HandleAuthorizeRequest"), jen.ID("mock").Dot(
				"Anything",
	),
	jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("moah"),
			jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot(
				"NewRecorder",
			).Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("s").Dot(
				"HandleAuthorizeRequest",
			).Call(jen.ID("res"), jen.ID("req"))),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().ID("TestService_HandleTokenRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot("Parallel").Call(),
		jen.Line(),
		jen.ID("T").Dot("Run").Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
			jen.ID("moah").Op(":=").Op("&").ID("mockOauth2Handler").Values(),
			jen.ID("moah").Dot(
				"On",
			).Call(jen.Lit("HandleTokenRequest"), jen.ID("mock").Dot(
				"Anything",
	),
	jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("nil")),
			jen.ID("s").Dot(
				"oauth2Handler",
			).Op("=").ID("moah"),
			jen.List(jen.ID("req"), jen.ID("res")).Op(":=").List(jen.ID("buildRequest").Call(jen.ID("t")), jen.ID("httptest").Dot(
				"NewRecorder",
			).Call()),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("s").Dot(
				"HandleTokenRequest",
			).Call(jen.ID("res"), jen.ID("req"))),
		)),
	),
	jen.Line(),
	)
	return ret
}
