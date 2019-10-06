package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func usersServiceDotGo() *jen.File {
	ret := jen.NewFile("users")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("MiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("user_input").Var().ID("counterName").ID("metrics").Dot(
		"CounterName",
	).Op("=").Lit("users").Var().ID("topicName").Op("=").Lit("users").Var().ID("serviceName").Op("=").Lit("users_service"),
	)
	ret.Add(jen.Null().Var().ID("_").ID("models").Dot(
		"UserDataServer",
	).Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
	)
	ret.Add(jen.Null().Type().ID("RequestValidator").Interface(jen.ID("Validate").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("bool"), jen.ID("error"))).Type().ID("Service").Struct(
		jen.ID("cookieSecret").Index().ID("byte"),
		jen.ID("database").ID("database").Dot(
			"Database",
		),
		jen.ID("authenticator").ID("auth").Dot(
			"Authenticator",
		),
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("encoderDecoder").ID("encoding").Dot(
			"EncoderDecoder",
		),
		jen.ID("userIDFetcher").ID("UserIDFetcher"),
		jen.ID("userCounter").ID("metrics").Dot(
			"UnitCounter",
		),
		jen.ID("reporter").ID("newsman").Dot(
			"Reporter",
		),
		jen.ID("userCreationEnabled").ID("bool"),
	).Type().ID("UserIDFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
