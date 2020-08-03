package users

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildUsersServiceConstDefs(proj)...)
	code.Add(buildUsersServiceVarDefs(proj)...)
	code.Add(buildUsersServiceTypeDefs(proj)...)
	code.Add(buildProvideUsersService(proj)...)

	return code
}

func buildUsersServiceConstDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.ID("serviceName").Equals().Lit("users_service"),
			jen.ID("topicName").Equals().Lit("users"),
			jen.ID("counterDescription").Equals().Lit("number of users managed by the users service"),
			jen.ID("counterName").Equals().Qual(proj.InternalMetricsV1Package(), "CounterName").Call(jen.ID("serviceName")),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersServiceVarDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Underscore().Qual(proj.ModelsV1Package(), "UserDataServer").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildUsersServiceTypeDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("RequestValidator validates request."),
			jen.ID("RequestValidator").Interface(
				jen.ID("Validate").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Bool(), jen.Error()),
			),
			jen.Line(),
			jen.ID("secretGenerator").Interface(
				jen.ID("GenerateTwoFactorSecret").Params().Params(jen.String(), jen.Error()),
				jen.ID("GenerateSalt").Params().Params(jen.Index().Byte(), jen.Error()),
			),
			jen.Line(),
			jen.Comment("UserIDFetcher fetches usernames from requests."),
			jen.ID("UserIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
			jen.Line(),
			jen.Comment("Service handles our users."),
			jen.ID("Service").Struct(
				jen.ID("cookieSecret").Index().Byte(),
				jen.ID("userDataManager").Qual(proj.ModelsV1Package(), "UserDataManager"),
				jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
				constants.LoggerParam(),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
				jen.ID("userIDFetcher").ID("UserIDFetcher"),
				jen.ID("userCounter").Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
				jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
				jen.ID("secretGenerator").ID("secretGenerator"),
				jen.ID("userCreationEnabled").Bool(),
			),
			jen.Line(),
		),
		jen.Line(),
	}

	return lines
}

func buildProvideUsersService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideUsersService builds a new UsersService."),
		jen.Line(),
		jen.Func().ID("ProvideUsersService").Paramsln(
			jen.ID("authSettings").Qual(proj.InternalConfigV1Package(), "AuthSettings"),
			constants.LoggerParam(),
			jen.ID("userDataManager").Qual(proj.ModelsV1Package(), "UserDataManager"),
			jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
			jen.ID("userIDFetcher").ID("UserIDFetcher"), jen.ID("encoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
			jen.ID("counterProvider").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider"),
			jen.ID("reporter").Qual("gitlab.com/verygoodsoftwarenotvirus/newsman", "Reporter"),
		).Params(jen.PointerTo().ID("Service"), jen.Error()).Body(
			jen.If(jen.ID("userIDFetcher").IsEqualTo().ID("nil")).Body(
				jen.Return().List(jen.Nil(), utils.Error("userIDFetcher must be provided")),
			),
			jen.Line(),
			jen.List(jen.ID("counter"), jen.Err()).Assign().ID("counterProvider").Call(
				jen.ID("counterName"),
				jen.ID("counterDescription"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID("cookieSecret").MapAssign().Index().Byte().Call(jen.ID("authSettings").Dot("CookieSecret")),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("userDataManager").MapAssign().ID("userDataManager"),
				jen.ID("authenticator").MapAssign().ID("authenticator"),
				jen.ID("userIDFetcher").MapAssign().ID("userIDFetcher"),
				jen.ID("encoderDecoder").MapAssign().ID("encoder"),
				jen.ID("userCounter").MapAssign().ID("counter"),
				jen.ID("reporter").MapAssign().ID("reporter"),
				jen.ID("secretGenerator").MapAssign().AddressOf().ID("standardSecretGenerator").Values(),
				jen.ID("userCreationEnabled").MapAssign().ID("authSettings").Dot("EnableUserSignup"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
