package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func middlewareDotGo() *jen.File {
	ret := jen.NewFile("auth")

	utils.AddImports(ret)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("UserLoginInputMiddlewareCtxKey is the context key for login input"),
			jen.ID("UserLoginInputMiddlewareCtxKey").ID("models").Dot("ContextKey").Op("=").Lit("user_login_input"),
			jen.Line(),
			jen.Comment("UsernameFormKey is the string we look for in request forms for username information"),
			jen.ID("UsernameFormKey").Op("=").Lit("username"),
			jen.Comment("PasswordFormKey is the string we look for in request forms for password information"),
			jen.ID("PasswordFormKey").Op("=").Lit("password"),
			jen.Comment("TOTPTokenFormKey is the string we look for in request forms for TOTP token information"),
			jen.ID("TOTPTokenFormKey").Op("=").Lit("totp_token"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("CookieAuthenticationMiddleware checks every request for a user cookie"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("CookieAuthenticationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("CookieAuthenticationMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.Comment("fetch the user from the request"),
				jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("FetchUserFromRequest").Call(jen.ID("ctx"), jen.ID("req")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered fetching user")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.Line(),
				jen.If(jen.ID("user").Op("!=").ID("nil")).Block(
					jen.ID("req").Op("=").ID("req").Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot("UserKey"), jen.ID("user")),
							jen.ID("models").Dot("UserIDKey"),
							jen.ID("user").Dot("ID"),
						),
					),
					jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("if no error was attached to the request, tell them to login first"),
				jen.Qual("net/http", "Redirect").Call(jen.ID("res"), jen.ID("req"), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("AuthenticationMiddleware authenticates based on either an oauth2 token or a cookie"),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("Service")).ID("AuthenticationMiddleware").Params(jen.ID("allowValidCookieInLieuOfAValidToken").ID("bool")).Params(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).Block(
			jen.Return().Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
				jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("AuthenticationMiddleware")),
					jen.Defer().ID("span").Dot("End").Call(),
					jen.Line(),
					jen.Comment("let's figure out who the user is"),
					jen.Var().ID("user").Op("*").ID("models").Dot("User"),
					jen.Line(),
					jen.Comment("check for a cookie first if we can"),
					jen.If(jen.ID("allowValidCookieInLieuOfAValidToken")).Block(
						jen.List(jen.ID("cookieAuth"), jen.ID("err")).Op(":=").ID("s").Dot("DecodeCookieFromRequest").Call(jen.ID("ctx"), jen.ID("req")),
						jen.Line(),
						jen.If(jen.ID("err").Op("==").ID("nil").Op("&&").ID("cookieAuth").Op("!=").ID("nil")).Block(
							jen.List(jen.ID("user"), jen.ID("err")).Op("=").ID("s").Dot("userDB").Dot("GetUser").Call(jen.ID("ctx"), jen.ID("cookieAuth").Dot("UserID")),
							jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
								jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error authenticating request")),
								jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("fetching user"), jen.Qual("net/http", "StatusInternalServerError")),
								jen.Comment("if we get here, then we just don't have a valid cookie, and we need to move on"),
								jen.Return(),
							),
						),
					),
					jen.Line(),
					jen.Comment("if the cookie wasn't present, or didn't indicate who the user is"),
					jen.If(jen.ID("user").Op("==").ID("nil")).Block(
						jen.Comment("check to see if there is an OAuth2 token for a valid client attached to the request."),
						jen.Comment("We do this first because it is presumed to be the primary means by which requests are made to the httpServer."),
						jen.List(jen.ID("oauth2Client"), jen.ID("err")).Op(":=").ID("s").Dot("oauth2ClientsService").Dot("ExtractOAuth2ClientFromRequest").Call(jen.ID("ctx"), jen.ID("req")),
						jen.If(jen.ID("err").Op("!=").ID("nil").Op("||").ID("oauth2Client").Op("==").ID("nil")).Block(
							jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("fetching oauth2 client")),
							jen.Qual("net/http", "Redirect").Call(jen.ID("res"), jen.ID("req"), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
							jen.Return(),
						),
						jen.Line(),
						jen.Comment("attach the oauth2 client and user's info to the request"),
						jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot("OAuth2ClientKey"), jen.ID("oauth2Client")),
						jen.List(jen.ID("user"), jen.ID("err")).Op("=").ID("s").Dot("userDB").Dot("GetUser").Call(jen.ID("ctx"), jen.ID("oauth2Client").Dot("BelongsTo")),
						jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
							jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error authenticating request")),
							jen.Qual("net/http", "Error").Call(jen.ID("res"), jen.Lit("fetching user"), jen.Qual("net/http", "StatusInternalServerError")),
							jen.Return(),
						),
					),
					jen.Line(),
					jen.Comment("If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere"),
					jen.If(jen.ID("user").Op("==").ID("nil")).Block(
						jen.ID("s").Dot("logger").Dot("Debug").Call(jen.Lit("no user attached to request request")),
						jen.Qual("net/http", "Redirect").Call(jen.ID("res"), jen.ID("req"), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
						jen.Return(),
					),
					jen.Line(),
					jen.Comment("elsewise, load the request with extra context"),
					jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot("UserKey"), jen.ID("user")),
					jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot("UserIDKey"), jen.ID("user").Dot("ID")),
					jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("models").Dot("UserIsAdminKey"), jen.ID("user").Dot("IsAdmin")),
					jen.Line(),
					jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(jen.ID("ctx"))),
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
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("AdminMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
				jen.List(jen.ID("user"), jen.ID("ok")).Op(":=").ID("ctx").Dot("Value").Call(jen.ID("models").Dot("UserKey")).Assert(jen.Op("*").ID("models").Dot("User")),
				jen.Line(),
				jen.If(jen.Op("!").ID("ok").Op("||").ID("user").Op("==").ID("nil")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("AdminMiddleware called without user attached to context")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.Line(),
				jen.If(jen.Op("!").ID("user").Dot("IsAdmin")).Block(
					jen.ID("logger").Dot("Debug").Call(jen.Lit("AdminMiddleware called by non-admin user")),
					jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
					jen.Return(),
				),
				jen.Line(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req")),
			)),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Comment("parseLoginInputFromForm checks a request for a login form, and returns the parsed login data if relevant"),
		jen.Line(),
		jen.Func().ID("parseLoginInputFromForm").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot("UserLoginInput")).Block(
			jen.If(jen.ID("err").Op(":=").ID("req").Dot("ParseForm").Call(), jen.ID("err").Op("==").ID("nil")).Block(
				jen.ID("uli").Op(":=").Op("&").ID("models").Dot("UserLoginInput").Valuesln(
					jen.ID("Username").Op(":").ID("req").Dot("FormValue").Call(jen.ID("UsernameFormKey")),
					jen.ID("Password").Op(":").ID("req").Dot("FormValue").Call(jen.ID("PasswordFormKey")),
					jen.ID("TOTPToken").Op(":").ID("req").Dot("FormValue").Call(jen.ID("TOTPTokenFormKey")),
				),
				jen.Line(),
				jen.If(jen.ID("uli").Dot("Username").Op("!=").Lit("").Op("&&").ID("uli").Dot("Password").Op("!=").Lit("").Op("&&").ID("uli").Dot("TOTPToken").Op("!=").Lit("")).Block(
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
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot("Context").Call(), jen.Lit("UserLoginInputMiddleware")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Line(),
				jen.ID("x").Op(":=").ID("new").Call(jen.ID("models").Dot("UserLoginInput")),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID("req"), jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.If(jen.ID("x").Op("=").ID("parseLoginInputFromForm").Call(jen.ID("req")), jen.ID("x").Op("==").ID("nil")).Block(
						jen.ID("s").Dot("logger").Dot("Error").Call(jen.ID("err"), jen.Lit("error encountered decoding request body")),
						jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusBadRequest")),
						jen.Return(),
					),
				),
				jen.Line(),
				jen.ID("ctx").Op("=").Qual("context", "WithValue").Call(jen.ID("ctx"), jen.ID("UserLoginInputMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID("res"), jen.ID("req").Dot("WithContext").Call(jen.ID("ctx"))),
			)),
		),
		jen.Line(),
	)
	return ret
}
