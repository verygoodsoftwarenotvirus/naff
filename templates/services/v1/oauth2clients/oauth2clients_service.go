package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildServiceInit()...)
	code.Add(buildServiceConstDefs(proj)...)
	code.Add(buildServiceVarDefs(proj)...)
	code.Add(buildServiceTypeDefs(proj)...)
	code.Add(buildServiceNewClientStore(proj)...)
	code.Add(buildServiceGetByID()...)
	code.Add(buildServiceProvideOAuth2ClientsService(proj)...)
	code.Add(buildServiceInitializeOAuth2Handler()...)
	code.Add(buildServiceHandleAuthorizeRequest()...)
	code.Add(buildServiceHandleTokenRequest()...)

	return code
}

func buildServiceInit() []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("init").Params().Body(
			jen.ID("b").Assign().ID("make").Call(jen.Index().Byte(), jen.Lit(64)),
			jen.If(jen.List(jen.Underscore(), jen.Err()).Assign().Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceConstDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("creationMiddlewareCtxKey is a string alias for referring to OAuth2 client creation data."),
			jen.ID("creationMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("create_oauth2_client"),
			jen.Line(),
			jen.ID("counterName").Qual(proj.InternalMetricsV1Package(), "CounterName").Equals().Lit("oauth2_clients"),
			jen.ID("counterDescription").String().Equals().Lit("number of oauth2 clients managed by the oauth2 client service"),
			jen.ID("serviceName").String().Equals().Lit("oauth2_clients_service"),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceVarDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Var().Defs(
			jen.Underscore().Qual(proj.ModelsV1Package(), "OAuth2ClientDataServer").Equals().Parens(jen.PointerTo().ID("Service")).Call(jen.Nil()),
			jen.Underscore().Qual("gopkg.in/oauth2.v3", "ClientStore").Equals().Parens(jen.PointerTo().ID("clientStore")).Call(jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceTypeDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.ID("oauth2Handler").Interface(
				jen.ID("SetAllowGetAccessRequest").Params(jen.Bool()),
				jen.ID("SetClientAuthorizedHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler")),
				jen.ID("SetClientScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler")),
				jen.ID("SetClientInfoHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientInfoHandler")),
				jen.ID("SetUserAuthorizationHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler")),
				jen.ID("SetAuthorizeScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler")),
				jen.ID("SetResponseErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler")),
				jen.ID("SetInternalErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler")),
				jen.ID("ValidationBearerToken").Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.Qual("gopkg.in/oauth2.v3", "TokenInfo"), jen.Error()),
				jen.ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()),
				jen.ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()),
			),
			jen.Line(),
			jen.Comment("ClientIDFetcher is a function for fetching client IDs out of requests."),
			jen.ID("ClientIDFetcher").Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
			jen.Line(),
			jen.Comment("Service manages our OAuth2 clients via HTTP."),
			jen.ID("Service").Struct(
				constants.LoggerParam(),
				jen.ID("database").Qual(proj.DatabaseV1Package(), "DataManager"),
				jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
				jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
				jen.ID("urlClientIDExtractor").Func().Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Uint64()),
				jen.ID("oauth2Handler").ID("oauth2Handler"),
				jen.ID("oauth2ClientCounter").Qual(proj.InternalMetricsV1Package(), "UnitCounter"),
			),

			jen.Line(),
			jen.ID("clientStore").Struct(
				jen.ID("database").Qual(proj.DatabaseV1Package(), "DataManager"),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceNewClientStore(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Func().ID("newClientStore").Params(jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager")).Params(jen.PointerTo().ID("clientStore")).Body(
			jen.ID("cs").Assign().AddressOf().ID("clientStore").Valuesln(
				jen.ID("database").MapAssign().ID("db")),
			jen.Return().ID("cs"),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceGetByID() []jen.Code {
	lines := []jen.Code{
		jen.Comment("GetByID implements oauth2.ClientStorage"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("clientStore")).ID("GetByID").Params(jen.ID("id").String()).Params(jen.Qual("gopkg.in/oauth2.v3", "ClientInfo"), jen.Error()).Body(
			jen.List(jen.ID("client"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(constants.InlineCtx(), jen.ID("id")),
			jen.Line(),
			jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Body(
				jen.Return().List(jen.Nil(), utils.Error("invalid client")),
			).Else().If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for client: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceProvideOAuth2ClientsService(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("ProvideOAuth2ClientsService builds a new OAuth2ClientsService."),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientsService").Paramsln(
			constants.LoggerParam(),
			jen.ID("db").Qual(proj.DatabaseV1Package(), "DataManager"),
			jen.ID("authenticator").Qual(proj.InternalAuthV1Package(), "Authenticator"),
			jen.ID("clientIDFetcher").ID("ClientIDFetcher"), jen.ID("encoderDecoder").Qual(proj.InternalEncodingV1Package(), "EncoderDecoder"),
			jen.ID("counterProvider").Qual(proj.InternalMetricsV1Package(), "UnitCounterProvider"),
		).Params(jen.PointerTo().ID("Service"), jen.Error()).Body(
			jen.ID("manager").Assign().Qual("gopkg.in/oauth2.v3/manage", "NewDefaultManager").Call(),
			jen.ID("clientStore").Assign().ID("newClientStore").Call(jen.ID("db")),
			jen.ID("manager").Dot("MapClientStorage").Call(jen.ID("clientStore")),
			jen.List(jen.ID("tokenStore"), jen.ID("tokenStoreErr")).Assign().Qual("gopkg.in/oauth2.v3/store", "NewMemoryTokenStore").Call(),
			jen.ID("manager").Dot("MustTokenStorage").Call(jen.ID("tokenStore"), jen.ID("tokenStoreErr")),
			jen.ID("manager").Dot("SetAuthorizeCodeTokenCfg").Call(jen.Qual("gopkg.in/oauth2.v3/manage", "DefaultAuthorizeCodeTokenCfg")),
			jen.ID("manager").Dot("SetRefreshTokenCfg").Call(jen.Qual("gopkg.in/oauth2.v3/manage", "DefaultRefreshTokenCfg")),
			jen.ID("oHandler").Assign().Qual("gopkg.in/oauth2.v3/server", "NewDefaultServer").Call(jen.ID("manager")),
			jen.ID("oHandler").Dot("SetAllowGetAccessRequest").Call(jen.True()),
			jen.Line(),
			jen.ID("svc").Assign().AddressOf().ID("Service").Valuesln(
				jen.ID("database").MapAssign().ID("db"),
				jen.ID(constants.LoggerVarName).MapAssign().ID(constants.LoggerVarName).Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("encoderDecoder").MapAssign().ID("encoderDecoder"),
				jen.ID("authenticator").MapAssign().ID("authenticator"),
				jen.ID("urlClientIDExtractor").MapAssign().ID("clientIDFetcher"),
				jen.ID("oauth2Handler").MapAssign().ID("oHandler"),
			),
			jen.ID("initializeOAuth2Handler").Call(jen.ID("svc")),
			jen.Line(),
			jen.Var().Err().Error(),
			jen.If(jen.List(jen.ID("svc").Dot("oauth2ClientCounter"), jen.Err()).Equals().ID("counterProvider").Call(jen.ID("counterName"), jen.ID("counterDescription")), jen.Err().DoesNotEqual().Nil()).Body(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("svc"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceInitializeOAuth2Handler() []jen.Code {
	lines := []jen.Code{
		jen.Comment("initializeOAuth2Handler."),
		jen.Line(),
		jen.Func().ID("initializeOAuth2Handler").Params(
			jen.ID("svc").PointerTo().ID("Service"),
		).Body(
			jen.ID("svc").Dot("oauth2Handler").Dot("SetAllowGetAccessRequest").Call(jen.True()),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetClientAuthorizedHandler").Call(jen.ID("svc").Dot("ClientAuthorizedHandler")),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetClientScopeHandler").Call(jen.ID("svc").Dot("ClientScopeHandler")),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetClientInfoHandler").Call(jen.Qual("gopkg.in/oauth2.v3/server", "ClientFormHandler")),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetAuthorizeScopeHandler").Call(jen.ID("svc").Dot("AuthorizeScopeHandler")),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetResponseErrorHandler").Call(jen.ID("svc").Dot("OAuth2ResponseErrorHandler")),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetInternalErrorHandler").Call(jen.ID("svc").Dot("OAuth2InternalErrorHandler")),
			jen.ID("svc").Dot("oauth2Handler").Dot("SetUserAuthorizationHandler").Call(jen.ID("svc").Dot("UserAuthorizationHandler")),
			jen.Line(),
			jen.Comment("this sad type cast is here because I have an arbitrary."),
			jen.Comment("test-only interface for OAuth2 interactions."),
			jen.If(jen.List(jen.ID("x"), jen.ID("ok")).Assign().ID("svc").Dot("oauth2Handler").Assert(jen.PointerTo().Qual("gopkg.in/oauth2.v3/server", "Server")), jen.ID("ok")).Body(
				jen.ID("x").Dot("Config").Dot("AllowedGrantTypes").Equals().Index().Qual("gopkg.in/oauth2.v3", "GrantType").Valuesln(
					jen.Qual("gopkg.in/oauth2.v3", "ClientCredentials"),
					jen.Comment("oauth2.AuthorizationCode"),
					jen.Comment("oauth2.Refreshing"),
					jen.Comment("oauth2.Implicit"),
				),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceHandleAuthorizeRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("HandleAuthorizeRequest is a simple wrapper around the internal server's HandleAuthorizeRequest."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("HandleAuthorizeRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Body(
			jen.Return().ID("s").Dot(
				"oauth2Handler",
			).Dot(
				"HandleAuthorizeRequest",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}

func buildServiceHandleTokenRequest() []jen.Code {
	lines := []jen.Code{
		jen.Comment("HandleTokenRequest is a simple wrapper around the internal server's HandleTokenRequest."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("HandleTokenRequest").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Error()).Body(
			jen.Return().ID("s").Dot(
				"oauth2Handler",
			).Dot(
				"HandleTokenRequest",
			).Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
		),
		jen.Line(),
	}

	return lines
}
