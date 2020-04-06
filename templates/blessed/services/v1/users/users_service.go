package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersServiceDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("users")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("serviceName").String().Equals().Lit("users_service"),
			jen.ID("topicName").String().Equals().Lit("users"),
			jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName").Equals().Lit("users"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.ModelsV1Package(), "UserDataServer").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.Comment("RequestValidator validates request"),
			jen.ID("RequestValidator").Interface(
				jen.ID("Validate").Params(jen.ID("req").ParamPointer().Qual("net/http", "Request")).Params(jen.Bool(), jen.Error()),
			),
			jen.Line(),
			jen.Comment("Service handles our users"),
			jen.ID("Service").Struct(
				jen.ID("cookieSecret").Index().Byte(), jen.ID("database").Qual(proj.DatabaseV1Package(), "Database"),
				jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("userCounter").Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
				jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
				jen.ID("userCreationEnabled").Bool(),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher fetches usernames from requests"),
			jen.ID("UserIDFetcher").Func().Params(jen.ParamPointer().Qual("net/http", "Request")).Params(jen.Uint64()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideUsersService builds a new UsersService"),
		jen.Line(),
		jen.Func().ID("ProvideUsersService").Paramsln(
			utils.CtxParam(), jen.ID("authSettings").Qual(proj.InternalConfigV1Package(), "AuthSettings"),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("db").Qual(proj.DatabaseV1Package(), "Database"),
			jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
			jen.ID("counterProvider").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider"),
			jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
		).Params(jen.PointerTo().ID("Service"), jen.Error()).Block(
			jen.If(jen.ID("userIDFetcher").Op("==").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("userIDFetcher must be provided"))),
			),
			jen.Line(),
			jen.List(jen.ID("counter"), jen.Err()).Assign().ID("counterProvider").Call(jen.ID("counterName"), jen.Lit("number of users managed by the users service")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("userCount"), jen.Err()).Assign().ID("db").Dot("GetAllUserCount").Call(utils.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user count: %w"), jen.Err())),
			),
			jen.ID("counter").Dot("IncrementBy").Call(utils.CtxVar(), jen.ID("userCount")),
			jen.Line(),
			jen.ID("us").Assign().VarPointer().ID("Service").Valuesln(
				jen.ID("cookieSecret").MapAssign().Index().Byte().Call(jen.ID("authSettings").Dot("CookieSecret")),
				jen.ID("logger").MapAssign().ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("database").MapAssign().ID("db"),
				jen.ID("authenticator").MapAssign().ID("authenticator"),
				jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"),
				jen.ID("encoderDecoder").MapAssign().ID("encoder"),
				jen.ID("userCounter").MapAssign().ID("counter"),
				jen.ID("reporter").MapAssign().ID("reporter"),
				jen.ID("userCreationEnabled").MapAssign().ID("authSettings").Dot("EnableUserSignup"),
			),
			jen.Return().List(jen.ID("us"), jen.Nil()),
		),
		jen.Line(),
	)
	return ret
}
