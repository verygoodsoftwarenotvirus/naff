package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(
		jen.Var().ID("UserLoginInputMiddlewareCtxKey").ID("models").Dot(
		"ContextKey",
	).Op("=").Lit("user_login_input").Var().ID("UsernameFormKey").Op("=").Lit("username").Var().ID("PasswordFormKey").Op("=").Lit("password").Var().ID("TOTPTokenFormKey").Op("=").Lit("totp_token"),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("CookieAuthenticationMiddleware checks every request for a user cookie"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CookieAuthenticationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("CookieAuthenticationMiddleware")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot(
				"FetchUserFromRequest",
			).Call(jen.ID("ctx"), jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("error encountered fetching user")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			),
			jen.If(jen.ID("user").Op("!=").ID("nil")).Block(
				jen.ID("req").Op("=").ID("req").Dot(
					"WithContext",
				).Call(jen.Qual("context", "WithValue").Call(jen.Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
					"UserKey",
				), jen.ID("user")), jen.ID("models").Dot(
					"UserIDKey",
				), jen.ID("user").Dot(
					"ID",
				))),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req")),
				jen.Return(),
			),
			jen.Qual("net/http", "Redirect").Call(jen.ID("res"), jen.ID("req"), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthenticationMiddleware authenticates based on either an oauth2 token or a cookie"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("AuthenticationMiddleware").Params(jen.ID("allowValidCookieInLieuOfAValidToken").ID("bool")).Params(jen.Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).Block(
		jen.Return().Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
					"Context",
				).Call(), jen.Lit("AuthenticationMiddleware")),
				jen.Defer().ID("span").Dot(
					"End",
				).Call(),

		jen.Var().ID("user").Op("*").ID("models").Dot(
					"User",
				),
				jen.If(jen.ID("allowValidCookieInLieuOfAValidToken")).Block(
					jen.List(jen.ID("cookieAuth"), jen.ID("err")).Op(":=").ID("s").Dot(
						"DecodeCookieFromRequest",
					).Call(jen.ID("ctx"), jen.ID("req")),
					jen.If(jen.ID("err").Op("==").ID("nil").Op("&&").ID("cookieAuth").Op("!=").ID("nil")).Block(
						jen.List(jen.ID("user"), jen.ID("err")).Op("=").ID("s").Dot(
							"userDB",
						).Dot(
							"GetUser",
						).Call(jen.ID("ctx"), jen.ID("cookieAuth").Dot(
							"UserID",
						)),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
							jen.ID("s").Dot(
								"logger",
							).Dot(
								"Error",
							).Call(jen.ID("err"), jen.Lit("error authenticating request")),
							jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("fetching user"), jen.Qual("net/http", "StatusInternalServerError")),
							jen.Return(),
						),
					),
				),
				jen.If(jen.ID("user").Op("==").ID("nil")).Block(
					jen.List(jen.ID("oauth2Client"), jen.ID("err")).Op(":=").ID("s").Dot(
						"oauth2ClientsService",
					).Dot(
						"ExtractOAuth2ClientFromRequest",
					).Call(jen.ID("ctx"), jen.ID("req")),
					jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("oauth2Client").Op("==").ID("nil")).Block(
						jen.ID("s").Dot(
							"logger",
						).Dot(
							"Error",
						).Call(jen.ID("err"), jen.Lit("fetching oauth2 client")),
						jen.Qual("net/http", "Redirect").Call(jen.ID("res"), jen.ID("req"), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
						jen.Return(),
					),
					jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
						"OAuth2ClientKey",
					), jen.ID("oauth2Client")),
					jen.List(jen.ID("user"), jen.ID("err")).Op("=").ID("s").Dot(
						"userDB",
					).Dot(
						"GetUser",
					).Call(jen.ID("ctx"), jen.ID("oauth2Client").Dot(
						"BelongsTo",
					)),
					jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
						jen.ID("s").Dot(
							"logger",
						).Dot(
							"Error",
						).Call(jen.ID("err"), jen.Lit("error authenticating request")),
						jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("fetching user"), jen.Qual("net/http", "StatusInternalServerError")),
						jen.Return(),
					),
				),
				jen.If(jen.ID("user").Op("==").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Debug",
					).Call(jen.Lit("no user attached to request request")),
					jen.Qual("net/http", "Redirect").Call(jen.ID("res"), jen.ID("req"), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
					"UserKey",
				), jen.ID("user")),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
					"UserIDKey",
				), jen.ID("user").Dot(
					"ID",
				)),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot(
					"UserIsAdminKey",
				), jen.ID("user").Dot(
					"IsAdmin",
				)),
				jen.ID("next").Dot(
					"ServeHTTP",
				).Call(jen.ID("res"), jen.ID("req").Dot(
					"WithContext",
				).Call(jen.ID("ctx"))),
			)),
		),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("AdminMiddleware restricts requests to admin users only"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("AdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("AdminMiddleware")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithRequest",
			).Call(jen.ID("req")),
			jen.List(jen.ID("user"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
				"Value",
			).Call(jen.ID("models").Dot(
				"UserKey",
			)).Assert(jen.Op("*").ID("models").Dot(
				"User",
			)),
			jen.If(jen.Op("!").ID("ok").Op("||").ID("user").Op("==").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Debug",
				).Call(jen.Lit("AdminMiddleware called without user attached to context")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			),
			jen.If(jen.Op("!").ID("user").Dot(
				"IsAdmin",
			)).Block(
				jen.ID("logger").Dot(
					"Debug",
				).Call(jen.Lit("AdminMiddleware called by non-admin user")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			),
			jen.ID("next").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req")),
		)),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Func().Comment("// parseLoginInputFromForm checks a request for a login form, and returns the parsed login data if relevant").ID("parseLoginInputFromForm").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"UserLoginInput",
	)).Block(
		jen.If(jen.ID("err").Op(":=").ID("req").Dot(
			"ParseForm",
		).Call(), jen.ID("err").Op("==").ID("nil")).Block(
			jen.ID("uli").Op(":=").Op("&").ID("models").Dot(
				"UserLoginInput",
			).Valuesln(jen.ID("Username").Op(":").ID("req").Dot(
				"FormValue",
			).Call(jen.ID("UsernameFormKey")), jen.ID("Password").Op(":").ID("req").Dot(
				"FormValue",
			).Call(jen.ID("PasswordFormKey")), jen.ID("TOTPToken").Op(":").ID("req").Dot(
				"FormValue",
			).Call(jen.ID("TOTPTokenFormKey"))),
			jen.If(jen.ID("uli").Dot(
				"Username",
			).Op("!=").Lit("").Op("&&").ID("uli").Dot(
				"Password",
			).Op("!=").Lit("").Op("&&").ID("uli").Dot(
				"TOTPToken",
			).Op("!=").Lit("")).Block(
				jen.Return().ID("uli"),
			),
		),
		jen.Return().ID("nil"),
	),
	jen.Line(),
	)

	ret.Add(
		jen.Comment("UserLoginInputMiddleware fetches user login input from requests"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("UserLoginInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
		jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("UserLoginInputMiddleware")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot(
				"UserLoginInput",
			)),
			jen.If(jen.ID("err").Op(":=").ID("s").Dot(
				"encoderDecoder",
			).Dot(
				"DecodeRequest",
			).Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.If(jen.ID("x").Op("=").ID("parseLoginInputFromForm").Call(jen.ID("req")), jen.ID("x").Op("==").ID("nil")).Block(
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
			),
			jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("x")),
			jen.ID("next").Dot(
				"ServeHTTP",
			).Call(jen.ID("res"), jen.ID("req").Dot(
				"WithContext",
			).Call(jen.ID("ctx"))),
		)),
	),
	jen.Line(),
	)
	return ret
}
