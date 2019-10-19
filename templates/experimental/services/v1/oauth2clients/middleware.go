package oauth2clients

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("oauth2clients")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("scopesSeparator").Op("=").Lit(",").Var().ID("apiPathPrefix").Op("=").Lit("/api/v1/"),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CreationInputMiddleware is a middleware for attaching OAuth2 client info to a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("CreationInputMiddleware")),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(
					"OAuth2ClientCreationInput",
				)),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"DecodeRequest",
				).Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Error",
					).Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
					jen.ID("res").Dot(
						"WriteHeader",
					).Call(jen.Qual("net/http", "StatusBadRequest")),
					jen.Return(),
				),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("CreationMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req").Dot(
					"WithContext",
				).Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("ExtractOAuth2ClientFromRequest extracts OAuth2 client data from a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("ExtractOAuth2ClientFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
	),
	jen.ID("error")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("ExtractOAuth2ClientFromRequest")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValue",
			).Call(jen.Lit("function_name"), jen.Lit("ExtractOAuth2ClientFromRequest")),
			jen.List(jen.ID("token"), jen.ID("err")).Op(":=").ID("s").Dot(
				"oauth2Handler",
			).Dot(
				"ValidationBearerToken",
			).Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("validating bearer token: %w"), jen.ID("err"))),
			),
			jen.ID("clientID").Op(":=").ID("token").Dot(
				"GetClientID",
			).Call(),
			jen.ID("logger").Op("=").ID("logger").Dot(
				"WithValue",
			).Call(jen.Lit("client_id"), jen.ID("clientID")),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
				"database",
			).Dot(
				"GetOAuth2ClientByClientID",
			).Call(jen.ID("ctx"), jen.ID("clientID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("error fetching OAuth2 Client")),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("scope").Op(":=").ID("determineScope").Call(jen.ID("req")),
			jen.ID("hasScope").Op(":=").ID("c").Dot(
				"HasScope",
			).Call(jen.ID("scope")),
			jen.ID("logger").Op("=").ID("logger").Dot(
				"WithValue",
			).Call(jen.Lit("scope"), jen.ID("scope")).Dot(
				"WithValue",
			).Call(jen.Lit("scopes"), jen.Qual("strings", "Join").Call(jen.ID("c").Dot(
				"Scopes",
	),
	jen.ID("scopesSeparator"))),
			jen.If(jen.Op("!").ID("hasScope")).Block(
				jen.ID("logger").Dot(
					"Info",
				).Call(jen.Lit("rejecting client for invalid scope")),
				jen.Return().List(jen.ID("nil"), jen.Qual("errors", "New").Call(jen.Lit("client not authorized for scope"))),
			),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("determineScope determines the scope of a request by its URL"),
		jen.Line(),
		jen.Comment("// this may be more ideally embedded as a struct field and placed"),
		jen.Line(),
		jen.Comment("// in the HTTP server's package instead"),
		jen.Line(),
		jen.Func().ID("determineScope").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Block(
			jen.If(jen.Qual("strings", "HasPrefix").Call(jen.ID("req").Dot(
				"URL",
			).Dot(
				"Path",
	),
	jen.ID("apiPathPrefix"))).Block(
				jen.ID("x").Op(":=").Qual("strings", "TrimPrefix").Call(jen.ID("req").Dot(
					"URL",
				).Dot(
					"Path",
	),
	jen.ID("apiPathPrefix")),
				jen.If(jen.ID("y").Op(":=").Qual("strings", "Split").Call(jen.ID("x"), jen.Lit("/")), jen.ID("len").Call(jen.ID("y")).Op(">").Lit(0)).Block(
					jen.ID("x").Op("=").ID("y").Index(jen.Lit(0)),
				),
				jen.Return().ID("x"),
			),
			jen.Return().Lit(""),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2TokenAuthenticationMiddleware authenticates Oauth tokens"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("OAuth2TokenAuthenticationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("OAuth2TokenAuthenticationMiddleware")),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("s").Dot(
					"ExtractOAuth2ClientFromRequest",
				).Call(jen.ID("ctx"), jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Error",
					).Call(jen.ID("err"), jen.Lit("error authenticated token-authed request")),
					jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("invalid token"), jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("c").Dot(
					"BelongsTo",
				)),
				jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("c").Dot(
					"ID",
				)),
				jen.ID("attachOAuth2ClientIDToSpan").Call(jen.ID("span"), jen.ID("c").Dot(
					"ClientID",
				)),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
					"OAuth2ClientKey",
	),
	jen.ID("c")),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
					"UserIDKey",
	),
	jen.ID("c").Dot(
					"BelongsTo",
				)),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req").Dot(
					"WithContext",
				).Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("OAuth2ClientInfoMiddleware fetches clientOAuth2Client info from requests and attaches it explicitly to a request"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("OAuth2ClientInfoMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("OAuth2ClientInfoMiddleware")),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),
				jen.If(jen.ID("v").Op(":=").ID("req").Dot(
					"URL",
				).Dot(
					"Query",
				).Call().Dot(
					"Get",
				).Call(jen.ID("oauth2ClientIDURIParamKey")), jen.ID("v").Op("!=").Lit("")).Block(
					jen.ID("logger").Op(":=").ID("s").Dot(
						"logger",
					).Dot(
						"WithValue",
					).Call(jen.Lit("oauth2_client_id"), jen.ID("v")),
					jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("s").Dot(
						"database",
					).Dot(
						"GetOAuth2ClientByClientID",
					).Call(jen.ID("ctx"), jen.ID("v")),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.ID("logger").Dot(
							"Error",
						).Call(jen.ID("err"), jen.Lit("error fetching OAuth2 client")),
						jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("invalid request"), jen.Qual("net/http", "StatusUnauthorized")),
						jen.Return(),
					),
					jen.ID("attachOAuth2ClientIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot(
						"ClientID",
					)),
					jen.ID("attachOAuth2ClientDatabaseIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot(
						"ID",
					)),
					jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("client").Dot(
						"BelongsTo",
					)),
					jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
						"OAuth2ClientKey",
	),
	jen.ID("client")),
					jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
						"UserIDKey",
	),
	jen.ID("client").Dot(
						"BelongsTo",
					)),
					jen.ID("req").Op("=").ID("req").Dot(
						"WithContext",
					).Call(jen.ID("ctx")),
				),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("fetchOAuth2ClientFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"OAuth2Client",
	)).Block(
		jen.List(jen.ID("client"), jen.ID("ok")).Op(":=").ID("req").Dot(
			"Context",
		).Call().Dot(
			"Value",
		).Call(jen.ID("models").Dot(
			"OAuth2ClientKey",
		)).Assert(jen.Op("*").ID("models").Dot(
			"OAuth2Client",
		)),
		jen.ID("_").Op("=").ID("ok"),
		jen.Return().ID("client"),
	),
		jen.Line(),
	)

	ret.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("fetchOAuth2ClientIDFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Block(
		jen.List(jen.ID("clientID"), jen.ID("ok")).Op(":=").ID("req").Dot(
			"Context",
		).Call().Dot(
			"Value",
		).Call(jen.ID("clientIDKey")).Assert(jen.ID("string")),
		jen.ID("_").Op("=").ID("ok"),
		jen.Return().ID("clientID"),
	),
		jen.Line(),
	)
	return ret
}
