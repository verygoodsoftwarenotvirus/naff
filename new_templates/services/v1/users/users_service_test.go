package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersServiceTestDotGo() *jen.File {
	ret := jen.NewFile("users")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Func().ID("buildTestService").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Service")).Block(
		jen.ID("t").Dot(
			"Helper",
		).Call(),
		jen.ID("expectedUserCount").Op(":=").ID("uint64").Call(jen.Lit(123)),
		jen.ID("mockDB").Op(":=").ID("database").Dot(
			"BuildMockDatabase",
		).Call(),
		jen.ID("mockDB").Dot(
			"UserDataManager",
		).Dot(
			"On",
		).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
			"Anything",
		), jen.Parens(jen.Op("*").ID("models").Dot(
			"QueryFilter",
		)).Call(jen.ID("nil"))).Dot(
			"Return",
		).Call(jen.ID("expectedUserCount"), jen.ID("nil")),
		jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
		jen.ID("uc").Dot(
			"On",
		).Call(jen.Lit("IncrementBy"), jen.ID("mock").Dot(
			"Anything",
		)),
		jen.Null().Var().ID("ucp").ID("metrics").Dot(
			"UnitCounterProvider",
		).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
			"CounterName",
		), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
			"UnitCounter",
		), jen.ID("error")).Block(
			jen.Return().List(jen.ID("uc"), jen.ID("nil")),
		),
		jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Call(jen.Qual("context", "Background").Call(), jen.ID("config").Dot(
			"AuthSettings",
		).Valuesln(), jen.ID("noop").Dot(
			"ProvideNoopLogger",
		).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
			jen.Return().Lit(0),
		), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("newsman").Dot(
			"NewNewsman",
		).Call(jen.ID("nil"), jen.ID("nil"))),
		jen.ID("require").Dot(
			"NoError",
		).Call(jen.ID("t"), jen.ID("err")),
		jen.Return().ID("service"),
	),
	)
	ret.Add(jen.Func().ID("TestProvideUsersService").Params(jen.ID("T").Op("*").Qual("testing", "T")).Block(
		jen.ID("T").Dot(
			"Parallel",
		).Call(),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("happy path"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("mockUserCount"), jen.ID("nil")),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("mockUserCount")).Dot(
				"Return",
			).Call(),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Call(jen.Qual("context", "Background").Call(), jen.ID("config").Dot(
				"AuthSettings",
			).Valuesln(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("assert").Dot(
				"NoError",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"NotNil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with nil userIDFetcher"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("mockUserCount"), jen.ID("nil")),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("mockUserCount")).Dot(
				"Return",
			).Call(),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Call(jen.Qual("context", "Background").Call(), jen.ID("config").Dot(
				"AuthSettings",
			).Valuesln(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(), jen.ID("nil"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error initializing counter"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("mockUserCount"), jen.ID("nil")),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.ID("uc").Dot(
				"On",
			).Call(jen.Lit("IncrementBy"), jen.ID("mockUserCount")).Dot(
				"Return",
			).Call(),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Call(jen.Qual("context", "Background").Call(), jen.ID("config").Dot(
				"AuthSettings",
			).Valuesln(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
		jen.ID("T").Dot(
			"Run",
		).Call(jen.Lit("with error getting user count"), jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Block(
			jen.ID("mockUserCount").Op(":=").ID("uint64").Call(jen.Lit(0)),
			jen.ID("mockDB").Op(":=").ID("database").Dot(
				"BuildMockDatabase",
			).Call(),
			jen.ID("mockDB").Dot(
				"UserDataManager",
			).Dot(
				"On",
			).Call(jen.Lit("GetUserCount"), jen.ID("mock").Dot(
				"Anything",
			), jen.ID("mock").Dot(
				"Anything",
			)).Dot(
				"Return",
			).Call(jen.ID("mockUserCount"), jen.Qual("errors", "New").Call(jen.Lit("blah"))),
			jen.ID("uc").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/metrics/v1/mock", "UnitCounter").Valuesln(),
			jen.Null().Var().ID("ucp").ID("metrics").Dot(
				"UnitCounterProvider",
			).Op("=").Func().Params(jen.ID("counterName").ID("metrics").Dot(
				"CounterName",
			), jen.ID("description").ID("string")).Params(jen.ID("metrics").Dot(
				"UnitCounter",
			), jen.ID("error")).Block(
				jen.Return().List(jen.ID("uc"), jen.ID("nil")),
			),
			jen.List(jen.ID("service"), jen.ID("err")).Op(":=").ID("ProvideUsersService").Call(jen.Qual("context", "Background").Call(), jen.ID("config").Dot(
				"AuthSettings",
			).Valuesln(), jen.ID("noop").Dot(
				"ProvideNoopLogger",
			).Call(), jen.ID("mockDB"), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/auth/v1/mock", "Authenticator").Valuesln(), jen.Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")).Block(
				jen.Return().Lit(0),
			), jen.Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/encoding/v1/mock", "EncoderDecoder").Valuesln(), jen.ID("ucp"), jen.ID("nil")),
			jen.ID("assert").Dot(
				"Error",
			).Call(jen.ID("t"), jen.ID("err")),
			jen.ID("assert").Dot(
				"Nil",
			).Call(jen.ID("t"), jen.ID("service")),
		)),
	),
	)
	return ret
}
