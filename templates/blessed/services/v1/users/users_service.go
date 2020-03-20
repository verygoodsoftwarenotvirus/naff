package users

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersServiceDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("MiddlewareCtxKey is the context key we search for when interacting with user-related requests"),
			jen.ID("MiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit("user_input"),
			jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName").Op("=").Lit("users"),
			jen.ID("topicName").Op("=").Lit("users"),
			jen.ID("serviceName").Op("=").Lit("users_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "UserDataServer").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("RequestValidator validates request"),
			jen.ID("RequestValidator").Interface(
				jen.ID("Validate").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("bool"), jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("Service handles our users"),
			jen.ID("Service").Struct(
				jen.ID("cookieSecret").Index().ID("byte"), jen.ID("database").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
				jen.ID("authenticator").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "Authenticator"),
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("encoderDecoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("userCounter").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
				jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
				jen.ID("userCreationEnabled").ID("bool"),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher fetches usernames from requests"),
			jen.ID("UserIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUsersService builds a new UsersService"),
		jen.Line(),
		jen.Func().ID("ProvideUsersService").Paramsln(
			utils.CtxParam(), jen.ID("authSettings").Qual(filepath.Join(pkg.OutputPath, "internal/v1/config"), "AuthSettings"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("db").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
			jen.ID("authenticator").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "Authenticator"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("encoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
			jen.ID("counterProvider").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider"),
			jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
		).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.If(jen.ID("userIDFetcher").Op("==").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("userIDFetcher must be provided"))),
			),
			jen.Line(),
			jen.List(jen.ID("counter"), jen.Err()).Op(":=").ID("counterProvider").Call(jen.ID("counterName"), jen.Lit("number of users managed by the users service")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("userCount"), jen.Err()).Op(":=").ID("db").Dot("GetUserCount").Call(utils.CtxVar(), jen.Nil()),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user count: %w"), jen.Err())),
			),
			jen.ID("counter").Dot("IncrementBy").Call(utils.CtxVar(), jen.ID("userCount")),
			jen.Line(),
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
			jen.Return().List(jen.ID("us"), jen.Nil()),
		),
		jen.Line(),
	)
	return ret
}
