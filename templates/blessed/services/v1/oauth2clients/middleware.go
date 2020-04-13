package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Const().Defs(
			jen.ID("scopesSeparator").Equals().Lit(","),
			jen.ID("apiPathPrefix").Equals().Lit("/api/v1/"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreationInputMiddleware is a middleware for attaching OAuth2 client info to a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "OAuth2ClientCreationInput")),
				jen.Line(),
				jen.Comment("decode value from request"),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("req"), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
					utils.WriteXHeader("res", "StatusBadRequest"),
					jen.Return(),
				),
				jen.Line(),
				utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.ID("CreationMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(utils.CtxVar())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ExtractOAuth2ClientFromRequest extracts OAuth2 client data from a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("ExtractOAuth2ClientFromRequest").Params(utils.CtxParam(), jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client"), jen.Error()).Block(
			jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(utils.CtxVar(), jen.Lit("ExtractOAuth2ClientFromRequest")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("function_name"), jen.Lit("ExtractOAuth2ClientFromRequest")),
			jen.Line(),
			jen.Comment("validate bearer token"),
			jen.List(jen.ID("token"), jen.Err()).Assign().ID("s").Dot("oauth2Handler").Dot("ValidationBearerToken").Call(jen.ID("req")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("validating bearer token: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Comment("fetch client ID"),
			jen.ID("clientID").Assign().ID("token").Dot("GetClientID").Call(),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(jen.Lit("client_id"), jen.ID("clientID")),
			jen.Line(),
			jen.Comment("fetch client by client ID"),
			jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("clientID")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 Client")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Comment("determine the scope"),
			jen.ID("scope").Assign().ID("determineScope").Call(jen.ID("req")),
			jen.ID("hasScope").Assign().ID("c").Dot("HasScope").Call(jen.ID("scope")),
			jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(jen.Lit("scope"), jen.ID("scope")).Dot("WithValue").Call(jen.Lit("scopes"), jen.Qual("strings", "Join").Call(jen.ID("c").Dot("Scopes"), jen.ID("scopesSeparator"))),
			jen.Line(),
			jen.If(jen.Op("!").ID("hasScope")).Block(
				jen.ID("logger").Dot("Info").Call(jen.Lit("rejecting client for invalid scope")),
				jen.Return().List(jen.Nil(), utils.Error("client not authorized for scope")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("c"), jen.Nil()),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("determineScope determines the scope of a request by its URL"),
		jen.Line(),
		jen.Comment("this may be more ideally embedded as a struct field and placed"),
		jen.Line(),
		jen.Comment("in the HTTP server's package instead"),
		jen.Line(),
		jen.Func().ID("determineScope").Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.String()).Block(
			jen.If(jen.Qual("strings", "HasPrefix").Call(jen.ID("req").Dot("URL").Dot("Path"),
				jen.ID("apiPathPrefix"))).Block(
				jen.ID("x").Assign().Qual("strings", "TrimPrefix").Call(jen.ID("req").Dot("URL").Dot("Path"), jen.ID("apiPathPrefix")),
				jen.If(jen.ID("y").Assign().Qual("strings", "Split").Call(jen.ID("x"), jen.Lit("/")), jen.Len(jen.ID("y")).GreaterThan().Zero()).Block(
					jen.ID("x").Equals().ID("y").Index(jen.Zero()),
				),
				jen.Return().ID("x"),
			),
			jen.Line(),
			jen.Return().EmptyString(),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2TokenAuthenticationMiddleware authenticates Oauth tokens"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("OAuth2TokenAuthenticationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("OAuth2TokenAuthenticationMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.List(jen.ID("c"), jen.Err()).Assign().ID("s").Dot("ExtractOAuth2ClientFromRequest").Call(utils.CtxVar(), jen.ID("req")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.Err(), jen.Lit("error authenticated token-authed request")),
					jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("invalid token"), jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.Line(),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("c").Dot("BelongsToUser")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("c").Dot("ID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientIDToSpan").Call(jen.ID("span"), jen.ID("c").Dot("ClientID")),
				jen.Line(),
				jen.Comment("attach both the user ID and the client object to the request. it might seem"),
				jen.Comment("superfluous, but some things should only need to know to look for user IDs"),
				utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("c")),
				utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "UserIDKey"), jen.ID("c").Dot("BelongsToUser")),
				jen.Line(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(utils.CtxVar())),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ClientInfoMiddleware fetches clientOAuth2Client info from requests and attaches it explicitly to a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").PointerTo().Qual("net/http", "Request")).Block(
				jen.List(utils.CtxVar(), jen.ID("span")).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("OAuth2ClientInfoMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.If(jen.ID("v").Assign().ID("req").Dot("URL").Dot("Query").Call().Dot("Get").Call(jen.ID("oauth2ClientIDURIParamKey")), jen.ID("v").DoesNotEqual().EmptyString()).Block(
					jen.ID("logger").Assign().ID("s").Dot("logger").Dot("WithValue").Call(jen.Lit("oauth2_client_id"), jen.ID("v")),
					jen.Line(),
					jen.List(jen.ID("client"), jen.Err()).Assign().ID("s").Dot("database").Dot("GetOAuth2ClientByClientID").Call(utils.CtxVar(), jen.ID("v")),
					jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID("logger").Dot("Error").Call(jen.Err(), jen.Lit("error fetching OAuth2 client")),
						jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("invalid request"), jen.Qual("net/http", "StatusUnauthorized")),
						jen.Return(),
					),
					jen.Line(),
					jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot("ClientID")),
					jen.Qual(proj.InternalTracingV1Package(), "AttachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot("ID")),
					jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot("BelongsToUser")),
					jen.Line(),
					utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("client")),
					utils.CtxVar().Equals().Qual("context", "WithValue").Call(utils.CtxVar(), jen.Qual(proj.ModelsV1Package(), "UserIDKey"), jen.ID("client").Dot("BelongsToUser")),
					jen.Line(),
					jen.ID("req").Equals().ID("req").Dot("WithContext").Call(utils.CtxVar()),
				),
				jen.Line(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("fetchOAuth2ClientFromRequest").Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")).Block(
			jen.List(jen.ID("client"), jen.ID("ok")).Assign().ID("req").Dot("Context").Call().Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "OAuth2Client")),
			jen.Underscore().Equals().ID("ok").Comment("we don't really care, but the linters do"),
			jen.Return().ID("client"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("fetchOAuth2ClientIDFromRequest").Params(jen.ID("req").PointerTo().Qual("net/http", "Request")).Params(jen.String()).Block(
			jen.List(jen.ID("clientID"), jen.ID("ok")).Assign().ID("req").Dot("Context").Call().Dot("Value").Call(jen.ID("clientIDKey")).Assert(jen.String()),
			jen.Underscore().Equals().ID("ok").Comment("we don't really care, but the linters do"),
			jen.Return().ID("clientID"),
		),
		jen.Line(),
	)
	return ret
}
