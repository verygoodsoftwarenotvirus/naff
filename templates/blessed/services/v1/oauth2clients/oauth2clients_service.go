package oauth2clients

import (
	"path/filepath"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func oauth2ClientsServiceDotGo(pkg *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(pkg, ret)

	ret.Add(
		jen.Func().ID("init").Params().Block(
			jen.ID("b").Op(":=").ID("make").Call(jen.Index().ID("byte"), jen.Lit(64)),
			jen.If(jen.List(jen.ID("_"), jen.Err()).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("b")), jen.Err().Op("!=").ID("nil")).Block(
				jen.ID("panic").Call(jen.Err()),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("CreationMiddlewareCtxKey is a string alias for referring to OAuth2 client creation data"),
			jen.ID("CreationMiddlewareCtxKey").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "ContextKey").Op("=").Lit("create_oauth2_client"),
			jen.Line(),
			jen.ID("counterName").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "CounterName").Op("=").Lit("oauth2_clients"),
			jen.ID("counterDescription").ID("string").Op("=").Lit("number of oauth2 clients managed by the oauth2 client service"),
			jen.ID("serviceName").ID("string").Op("=").Lit("oauth2_clients_service"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.ID("_").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "OAuth2ClientDataServer").Op("=").Parens(jen.Op("*").ID("Service")).Call(jen.Nil()),
			jen.ID("_").Qual("gopkg.in/oauth2.v3", "ClientStore").Op("=").Parens(jen.Op("*").ID("clientStore")).Call(jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Type().Defs(
			jen.ID("oauth2Handler").Interface(
				jen.ID("SetAllowGetAccessRequest").Params(jen.ID("bool")),
				jen.ID("SetClientAuthorizedHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientAuthorizedHandler")),
				jen.ID("SetClientScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientScopeHandler")),
				jen.ID("SetClientInfoHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ClientInfoHandler")),
				jen.ID("SetUserAuthorizationHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "UserAuthorizationHandler")),
				jen.ID("SetAuthorizeScopeHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "AuthorizeScopeHandler")),
				jen.ID("SetResponseErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "ResponseErrorHandler")),
				jen.ID("SetInternalErrorHandler").Params(jen.ID("handler").Qual("gopkg.in/oauth2.v3/server", "InternalErrorHandler")),
				jen.ID("ValidationBearerToken").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Qual("gopkg.in/oauth2.v3", "TokenInfo"), jen.ID("error")),
				jen.ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")),
				jen.ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")),
			),
			jen.Line(),
			jen.Comment("ClientIDFetcher is a function for fetching client IDs out of requests"),
			jen.ID("ClientIDFetcher").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			jen.Line(),
			jen.Comment("Service manages our OAuth2 clients via HTTP"),
			jen.ID("Service").Struct(
				jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
				jen.ID("database").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
				jen.ID("authenticator").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "Authenticator"),
				jen.ID("encoderDecoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
				jen.ID("urlClientIDExtractor").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.Line(),
				jen.ID("tokenStore").Qual("gopkg.in/oauth2.v3", "TokenStore"),
				jen.ID("oauth2Handler").ID("oauth2Handler"),
				jen.ID("oauth2ClientCounter").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounter"),
			),

			jen.Line(),
			jen.ID("clientStore").Struct(
				jen.ID("database").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().ID("newClientStore").Params(jen.ID("db").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database")).Params(jen.Op("*").ID("clientStore")).Block(
			jen.ID("cs").Op(":=").Op("&").ID("clientStore").Valuesln(
				jen.ID("database").Op(":").ID("db")),
			jen.Return().ID("cs"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("GetByID implements oauth2.ClientStorage"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("clientStore")).ID("GetByID").Params(jen.ID("id").ID("string")).Params(jen.Qual("gopkg.in/oauth2.v3", "ClientInfo"), jen.ID("error")).Block(
			jen.List(jen.ID("client"), jen.Err()).Op(":=").ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(utils.InlineCtx(), jen.ID("id")),
			jen.Line(),
			jen.If(jen.Err().Op("==").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("invalid client"))),
			).Else().If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("querying for client: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("client"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ProvideOAuth2ClientsService builds a new OAuth2ClientsService"),
		jen.Line(),
		jen.Func().ID("ProvideOAuth2ClientsService").Paramsln(
			utils.CtxParam(),
			jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"),
			jen.ID("db").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Database"),
			jen.ID("authenticator").Qual(filepath.Join(pkg.OutputPath, "internal/v1/auth"), "Authenticator"),
			jen.ID("clientIDFetcher").ID("ClientIDFetcher"), jen.ID("encoderDecoder").Qual(filepath.Join(pkg.OutputPath, "internal/v1/encoding"), "EncoderDecoder"),
			jen.ID("counterProvider").Qual(filepath.Join(pkg.OutputPath, "internal/v1/metrics"), "UnitCounterProvider"),
		).Params(jen.Op("*").ID("Service"), jen.ID("error")).Block(
			jen.List(jen.ID("counter"), jen.Err()).Op(":=").ID("counterProvider").Call(jen.ID("counterName"), jen.ID("counterDescription")),
			jen.If(jen.Err().Op("!=").ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("error initializing counter: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("manager").Op(":=").Qual("gopkg.in/oauth2.v3/manage", "NewDefaultManager").Call(),
			jen.ID("clientStore").Op(":=").ID("newClientStore").Call(jen.ID("db")),
			jen.List(jen.ID("tokenStore"), jen.Err()).Op(":=").Qual("gopkg.in/oauth2.v3/store", "NewMemoryTokenStore").Call(),
			jen.ID("manager").Dot("MapClientStorage").Call(jen.ID("clientStore")),
			jen.ID("manager").Dot("MustTokenStorage").Call(jen.ID("tokenStore"), jen.Err()),
			jen.ID("manager").Dot("SetAuthorizeCodeTokenCfg").Call(jen.Qual("gopkg.in/oauth2.v3/manage", "DefaultAuthorizeCodeTokenCfg")),
			jen.ID("manager").Dot("SetRefreshTokenCfg").Call(jen.Qual("gopkg.in/oauth2.v3/manage", "DefaultRefreshTokenCfg")),
			jen.ID("oHandler").Op(":=").Qual("gopkg.in/oauth2.v3/server", "NewDefaultServer").Call(jen.ID("manager")),
			jen.ID("oHandler").Dot("SetAllowGetAccessRequest").Call(jen.ID("true")),
			jen.Line(),
			jen.ID("s").Op(":=").Op("&").ID("Service").Valuesln(
				jen.ID("database").Op(":").ID("db"),
				jen.ID("logger").Op(":").ID("logger").Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("encoderDecoder").Op(":").ID("encoderDecoder"),
				jen.ID("authenticator").Op(":").ID("authenticator"),
				jen.ID("urlClientIDExtractor").Op(":").ID("clientIDFetcher"),
				jen.ID("oauth2ClientCounter").Op(":").ID("counter"),
				jen.ID("tokenStore").Op(":").ID("tokenStore"),
				jen.ID("oauth2Handler").Op(":").ID("oHandler"),
			),
			jen.Line(),
			jen.ID("initializeOAuth2Handler").Call(jen.ID("s").Dot("oauth2Handler"), jen.ID("s")),
			jen.List(jen.ID("count"), jen.Err()).Op(":=").ID("s").Dot("database").Dot("GetAllOAuth2ClientCount").Call(utils.CtxVar()),
			jen.If(jen.Err().Op("!=").ID("nil").Op("&&").ID("err").Op("!=").Qual("database/sql", "ErrNoRows")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching oauth2 clients: %w"), jen.Err())),
			),
			jen.ID("counter").Dot("IncrementBy").Call(utils.CtxVar(), jen.ID("count")),
			jen.Line(),
			jen.Return().List(jen.ID("s"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("initializeOAuth2Handler"),
		jen.Line(),
		jen.Func().ID("initializeOAuth2Handler").Params(jen.ID("handler").ID("oauth2Handler"), jen.ID("s").Op("*").ID("Service")).Block(
			jen.ID("handler").Dot("SetAllowGetAccessRequest").Call(jen.ID("true")),
			jen.ID("handler").Dot("SetClientAuthorizedHandler").Call(jen.ID("s").Dot("ClientAuthorizedHandler")),
			jen.ID("handler").Dot("SetClientScopeHandler").Call(jen.ID("s").Dot("ClientScopeHandler")),
			jen.ID("handler").Dot("SetClientInfoHandler").Call(jen.Qual("gopkg.in/oauth2.v3/server", "ClientFormHandler")),
			jen.ID("handler").Dot("SetAuthorizeScopeHandler").Call(jen.ID("s").Dot("AuthorizeScopeHandler")), jen.ID("handler").Dot("SetResponseErrorHandler").Call(jen.ID("s").Dot("OAuth2ResponseErrorHandler")),
			jen.ID("handler").Dot("SetInternalErrorHandler").Call(jen.ID("s").Dot("OAuth2InternalErrorHandler")),
			jen.ID("handler").Dot("SetUserAuthorizationHandler").Call(jen.ID("s").Dot("UserAuthorizationHandler")),
			jen.Line(),
			jen.Comment("this sad type cast is here because I have an arbitrary"),
			jen.Comment("test-only interface for OAuth2 interactions."),
			jen.If(jen.List(jen.ID("x"), jen.ID("ok")).Op(":=").ID("handler").Assert(jen.Op("*").Qual("gopkg.in/oauth2.v3/server", "Server")), jen.ID("ok")).Block(
				jen.ID("x").Dot("Config").Dot("AllowedGrantTypes").Op("=").Index().Qual("gopkg.in/oauth2.v3", "GrantType").Valuesln(
					jen.Qual("gopkg.in/oauth2.v3", "ClientCredentials"),
					jen.Comment("oauth2.AuthorizationCode"),
					jen.Comment("oauth2.Refreshing"),
					jen.Comment("oauth2.Implicit"),
				),
			),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HandleAuthorizeRequest is a simple wrapper around the internal server's HandleAuthorizeRequest"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("HandleAuthorizeRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Block(
			jen.Return().ID("s").Dot(
				"oauth2Handler",
			).Dot(
				"HandleAuthorizeRequest",
			).Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("HandleTokenRequest is a simple wrapper around the internal server's HandleTokenRequest"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("HandleTokenRequest").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("error")).Block(
			jen.Return().ID("s").Dot(
				"oauth2Handler",
			).Dot(
				"HandleTokenRequest",
			).Call(jen.ID("res"), jen.ID("req")),
		),
		jen.Line(),
	)
	return ret
}
