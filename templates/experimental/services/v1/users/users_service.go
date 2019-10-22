package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func usersServiceDotGo() *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("MiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("user_input"),
			jen.ID("counterName").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "CounterName").Op("=").Lit("users"),
			jen.ID("topicName").Op("=").Lit("users"),
			jen.ID("serviceName").Op("=").Lit("users_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().ID("_").ID("models").Dot("UserDataServer").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.ID("nil")),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.ID("RequestValidator").Interface(
				jen.ID("Validate").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("bool"), jen.ID("error")),
			),
			jen.ID("Service").Struct(
				jen.ID("cookieSecret").Index().ID("byte"), jen.ID("database").ID("database").Dot("Database"),
				jen.ID("authenticator").ID("auth").Dot("Authenticator"),
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("encoderDecoder").ID("encoding").Dot("EncoderDecoder"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("userCounter").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounter"),
				jen.ID("reporter").ID("newsman").Dot("Reporter"),
				jen.ID("userCreationEnabled").ID("bool"),
			),
			jen.ID("UserIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		),

		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUsersService builds a new UsersService"),
		jen.Line(),
		jen.Func().ID("ProvideUsersService").Paramsln(
			jen.ID("ctx").Qual("context", "Context"), jen.ID("authSettings").ID("config").Dot("AuthSettings"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("db").ID("database").Dot("Database"),
			jen.ID("authenticator").ID("auth").Dot("Authenticator"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("encoder").ID("encoding").Dot("EncoderDecoder"),
			jen.ID("counterProvider").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "UnitCounterProvider"),
			jen.ID("reporter").ID("newsman").Dot("Reporter"),
		).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.If(jen.ID("userIDFetcher").Op("==").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("userIDFetcher must be provided"))),
			),
			jen.List(jen.ID("counter"), jen.ID("err")).Op(":=").ID("counterProvider").Call(jen.ID("counterName"), jen.Lit("number of users managed by the users service")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.ID("err"))),
			),
			jen.List(jen.ID("userCount"), jen.ID("err")).Op(":=").ID("db").Dot("GetUserCount").Call(jen.ID("ctx"), jen.ID("nil")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user count: %w"), jen.ID("err"))),
			),
			jen.ID("counter").Dot("IncrementBy").Call(jen.ID("ctx"), jen.ID("userCount")),
			jen.ID("us").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("cookieSecret").Op(":").Index().ID("byte").Call(jen.ID("authSettings").Dot("CookieSecret")),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("database").Op(":").ID("db"),
				jen.ID("authenticator").Op(":").ID("authenticator"),
				jen.ID("userIDFetcher").Op(":").ID("userIDFetcher"),
				jen.ID("encoderDecoder").Op(":").ID("encoder"),
				jen.ID("userCounter").Op(":").ID("counter"),
				jen.ID("reporter").Op(":").ID("reporter"),
				jen.ID("userCreationEnabled").Op(":").ID("authSettings").Dot("EnableUserSignup"),
			),
			jen.Return().List(jen.ID("us"), jen.ID("nil")),
		),
		jen.Line(),
	)
	return ret
}
