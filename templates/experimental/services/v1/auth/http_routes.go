package auth

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func httpRoutesDotGo() *jen.File {
	ret := jen.NewFile("auth")
	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("CookieName").Op("=").Lit("todocookie").Var().ID("cookieErrorLogName").Op("=").Lit("_COOKIE_CONSTRUCTION_ERROR_"),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachUserIDToSpan provides a consistent way to attach a userID to a given span").ID("attachUserIDToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("userID").ID("uint64")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("user_id"), jen.Qual("strconv", "FormatUint").Call(jen.ID("userID"), jen.Lit(10)))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// attachUsernameToSpan provides a consistent way to attach a username to a given span").ID("attachUsernameToSpan").Params(jen.ID("span").Op("*").Qual("go.opencensus.io/trace", "Span"), jen.ID("username").ID("string")).Block(
		jen.If(jen.ID("span").Op("!=").ID("nil")).Block(
			jen.ID("span").Dot(
				"AddAttributes",
			).Call(jen.Qual("go.opencensus.io/trace", "StringAttribute").Call(jen.Lit("username"), jen.ID("username"))),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// DecodeCookieFromRequest takes a request object and fetches the cookie data if it is present").Params(jen.ID("s").Op("*").ID("Service")).ID("DecodeCookieFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("ca").Op("*").ID("models").Dot(
		"CookieAuth",
	), jen.ID("err").ID("error")).Block(
		jen.List(jen.ID("_"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("DecodeCookieFromRequest")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("req").Dot(
			"Cookie",
		).Call(jen.ID("CookieName")),
		jen.If(jen.ID("err").Op("!=").Qual("net/http", "ErrNoCookie").Op("&&").ID("cookie").Op("!=").ID("nil")).Block(
			jen.ID("decodeErr").Op(":=").ID("s").Dot(
				"cookieManager",
			).Dot(
				"Decode",
			).Call(jen.ID("CookieName"), jen.ID("cookie").Dot(
				"Value",
			), jen.Op("&").ID("ca")),
			jen.If(jen.ID("decodeErr").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("decoding request cookie")),
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("decoding request cookie: %w"), jen.ID("decodeErr"))),
			),
			jen.Return().List(jen.ID("ca"), jen.ID("nil")),
		),
		jen.Return().List(jen.ID("nil"), jen.Qual("net/http", "ErrNoCookie")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// WebsocketAuthFunction is provided to Newsman to determine if a user has access to websockets").Params(jen.ID("s").Op("*").ID("Service")).ID("WebsocketAuthFunction").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("bool")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
			"Context",
		).Call(), jen.Lit("WebsocketAuthFunction")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.List(jen.ID("oauth2Client"), jen.ID("err")).Op(":=").ID("s").Dot(
			"oauth2ClientsService",
		).Dot(
			"ExtractOAuth2ClientFromRequest",
		).Call(jen.ID("ctx"), jen.ID("req")),
		jen.If(jen.ID("err").Op("==").ID("nil").Op("&&").ID("oauth2Client").Op("!=").ID("nil")).Block(
			jen.Return().ID("true"),
		),
		jen.List(jen.ID("cookieAuth"), jen.ID("cookieErr")).Op(":=").ID("s").Dot(
			"DecodeCookieFromRequest",
		).Call(jen.ID("ctx"), jen.ID("req")),
		jen.If(jen.ID("cookieErr").Op("==").ID("nil").Op("&&").ID("cookieAuth").Op("!=").ID("nil")).Block(
			jen.Return().ID("true"),
		),
		jen.ID("s").Dot(
			"logger",
		).Dot(
			"Error",
		).Call(jen.ID("err"), jen.Lit("error authenticated token-authenticated request")),
		jen.Return().ID("false"),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// FetchUserFromRequest takes a request object and fetches the cookie, and then the user for that cookie").Params(jen.ID("s").Op("*").ID("Service")).ID("FetchUserFromRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("models").Dot(
		"User",
	), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("FetchUserFromRequest")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.List(jen.ID("ca"), jen.ID("decodeErr")).Op(":=").ID("s").Dot(
			"DecodeCookieFromRequest",
		).Call(jen.ID("ctx"), jen.ID("req")),
		jen.If(jen.ID("decodeErr").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching cookie data from request: %w"), jen.ID("decodeErr"))),
		),
		jen.List(jen.ID("user"), jen.ID("userFetchErr")).Op(":=").ID("s").Dot(
			"userDB",
		).Dot(
			"GetUser",
		).Call(jen.ID("req").Dot(
			"Context",
		).Call(), jen.ID("ca").Dot(
			"UserID",
		)),
		jen.If(jen.ID("userFetchErr").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user from request: %w"), jen.ID("userFetchErr"))),
		),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("ca").Dot(
			"UserID",
		)),
		jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("ca").Dot(
			"Username",
		)),
		jen.Return().List(jen.ID("user"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// LoginHandler is our login route").Params(jen.ID("s").Op("*").ID("Service")).ID("LoginHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("LoginHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.List(jen.ID("loginData"), jen.ID("errRes")).Op(":=").ID("s").Dot(
				"fetchLoginDataFromRequest",
			).Call(jen.ID("req")),
			jen.If(jen.ID("errRes").Op("!=").ID("nil")).Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"Error",
				).Call(jen.ID("errRes"), jen.Lit("error encountered fetching login data from request")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"EncodeResponse",
				).Call(jen.ID("res"), jen.ID("errRes")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Error",
					).Call(jen.ID("err"), jen.Lit("encoding response")),
				),
				jen.Return(),
			).Else().If(jen.ID("loginData").Op("==").ID("nil")).Block(
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			),
			jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("loginData").Dot(
				"user",
			).Dot(
				"ID",
			)),
			jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("loginData").Dot(
				"user",
			).Dot(
				"Username",
			)),
			jen.ID("logger").Op(":=").ID("s").Dot(
				"logger",
			).Dot(
				"WithValue",
			).Call(jen.Lit("user"), jen.ID("loginData").Dot(
				"user",
			).Dot(
				"ID",
			)),
			jen.List(jen.ID("loginValid"), jen.ID("err")).Op(":=").ID("s").Dot(
				"validateLogin",
			).Call(jen.ID("ctx"), jen.Op("*").ID("loginData")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("error encountered validating login")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot(
				"WithValue",
			).Call(jen.Lit("valid"), jen.ID("loginValid")),
			jen.If(jen.Op("!").ID("loginValid")).Block(
				jen.ID("logger").Dot(
					"Debug",
				).Call(jen.Lit("login was invalid")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.Return(),
			),
			jen.ID("logger").Dot(
				"Debug",
			).Call(jen.Lit("login was valid")),
			jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("s").Dot(
				"buildAuthCookie",
			).Call(jen.ID("loginData").Dot(
				"user",
			)),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot(
					"Error",
				).Call(jen.ID("err"), jen.Lit("error building cookie")),
				jen.ID("res").Dot(
					"WriteHeader",
				).Call(jen.Qual("net/http", "StatusInternalServerError")),
				jen.ID("response").Op(":=").Op("&").ID("models").Dot(
					"ErrorResponse",
				).Valuesln(jen.ID("Code").Op(":").Qual("net/http", "StatusInternalServerError"), jen.ID("Message").Op(":").Lit("error encountered building cookie")),
				jen.If(jen.ID("err").Op(":=").ID("s").Dot(
					"encoderDecoder",
				).Dot(
					"EncodeResponse",
				).Call(jen.ID("res"), jen.ID("response")), jen.ID("err").Op("!=").ID("nil")).Block(
					jen.ID("s").Dot(
						"logger",
					).Dot(
						"Error",
					).Call(jen.ID("err"), jen.Lit("encoding response")),
				),
				jen.Return(),
			),
			jen.Qual("net/http", "SetCookie").Call(jen.ID("res"), jen.ID("cookie")),
			jen.ID("res").Dot(
				"WriteHeader",
			).Call(jen.Qual("net/http", "StatusNoContent")),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// LogoutHandler is our logout route").Params(jen.ID("s").Op("*").ID("Service")).ID("LogoutHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("LogoutHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.If(jen.List(jen.ID("cookie"), jen.ID("err")).Op(":=").ID("req").Dot(
				"Cookie",
			).Call(jen.ID("CookieName")), jen.ID("err").Op("==").ID("nil").Op("&&").ID("cookie").Op("!=").ID("nil")).Block(
				jen.ID("c").Op(":=").ID("s").Dot(
					"buildCookie",
				).Call(jen.Lit("deleted")),
				jen.ID("c").Dot(
					"Expires",
				).Op("=").Qual("time", "Time").Valuesln(),
				jen.ID("c").Dot(
					"MaxAge",
				).Op("=").Op("-").Lit(1),
				jen.Qual("net/http", "SetCookie").Call(jen.ID("res"), jen.ID("c")),
			).Else().Block(
				jen.ID("s").Dot(
					"logger",
				).Dot(
					"WithError",
				).Call(jen.ID("err")).Dot(
					"Debug",
				).Call(jen.Lit("logout was called, no cookie was found")),
			),
			jen.ID("res").Dot(
				"WriteHeader",
			).Call(jen.Qual("net/http", "StatusOK")),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// CycleSecretHandler rotates the cookie building secret with a new random secret").Params(jen.ID("s").Op("*").ID("Service")).ID("CycleSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
		jen.Return().Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"Info",
			).Call(jen.Lit("cycling cookie secret!")),
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
				"Context",
			).Call(), jen.Lit("CycleSecretHandler")),
			jen.Defer().ID("span").Dot(
				"End",
			).Call(),
			jen.ID("s").Dot(
				"cookieManager",
			).Op("=").ID("securecookie").Dot(
				"New",
			).Call(jen.ID("securecookie").Dot(
				"GenerateRandomKey",
			).Call(jen.Lit(64)), jen.Index().ID("byte").Call(jen.ID("s").Dot(
				"config",
			).Dot(
				"CookieSecret",
			))),
			jen.ID("res").Dot(
				"WriteHeader",
			).Call(jen.Qual("net/http", "StatusCreated")),
		),
	),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("loginData").Struct(jen.ID("loginInput").Op("*").ID("models").Dot(
		"UserLoginInput",
	), jen.ID("user").Op("*").ID("models").Dot(
		"User",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// fetchLoginDataFromRequest searches a given HTTP request for parsed login input data, and").Comment("// returns a helper struct with the relevant login information").Params(jen.ID("s").Op("*").ID("Service")).ID("fetchLoginDataFromRequest").Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("loginData"), jen.Op("*").ID("models").Dot(
		"ErrorResponse",
	)).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("req").Dot(
			"Context",
		).Call(), jen.Lit("fetchLoginDataFromRequest")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.List(jen.ID("loginInput"), jen.ID("ok")).Op(":=").ID("ctx").Dot(
			"Value",
		).Call(jen.ID("UserLoginInputMiddlewareCtxKey")).Assert(jen.Op("*").ID("models").Dot(
			"UserLoginInput",
		)),
		jen.If(jen.Op("!").ID("ok")).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"Debug",
			).Call(jen.Lit("no UserLoginInput found for /login request")),
			jen.Return().List(jen.ID("nil"), jen.Op("&").ID("models").Dot(
				"ErrorResponse",
			).Valuesln(jen.ID("Code").Op(":").Qual("net/http", "StatusUnauthorized"))),
		),
		jen.ID("username").Op(":=").ID("loginInput").Dot(
			"Username",
		),
		jen.ID("attachUsernameToSpan").Call(jen.ID("span"), jen.ID("username")),
		jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot(
			"userDB",
		).Dot(
			"GetUserByUsername",
		).Call(jen.ID("ctx"), jen.ID("username")),
		jen.If(jen.ID("err").Op("==").Qual("database/sql", "ErrNoRows")).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("no matching user")),
			jen.Return().List(jen.ID("nil"), jen.Op("&").ID("models").Dot(
				"ErrorResponse",
			).Valuesln(jen.ID("Code").Op(":").Qual("net/http", "StatusBadRequest"))),
		).Else().If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("error fetching user")),
			jen.Return().List(jen.ID("nil"), jen.Op("&").ID("models").Dot(
				"ErrorResponse",
			).Valuesln(jen.ID("Code").Op(":").Qual("net/http", "StatusInternalServerError"))),
		),
		jen.ID("attachUserIDToSpan").Call(jen.ID("span"), jen.ID("user").Dot(
			"ID",
		)),
		jen.ID("ld").Op(":=").Op("&").ID("loginData").Valuesln(jen.ID("loginInput").Op(":").ID("loginInput"), jen.ID("user").Op(":").ID("user")),
		jen.Return().List(jen.ID("ld"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// validateLogin takes login information and returns whether or not the login is valid.").Comment("// In the event that there's an error, this function will return false and the error.").Params(jen.ID("s").Op("*").ID("Service")).ID("validateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("loginInfo").ID("loginData")).Params(jen.ID("bool"), jen.ID("error")).Block(
		jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").Qual("go.opencensus.io/trace", "StartSpan").Call(jen.ID("ctx"), jen.Lit("validateLogin")),
		jen.Defer().ID("span").Dot(
			"End",
		).Call(),
		jen.ID("user").Op(":=").ID("loginInfo").Dot(
			"user",
		),
		jen.ID("loginInput").Op(":=").ID("loginInfo").Dot(
			"loginInput",
		),
		jen.ID("logger").Op(":=").ID("s").Dot(
			"logger",
		).Dot(
			"WithValue",
		).Call(jen.Lit("username"), jen.ID("user").Dot(
			"Username",
		)),
		jen.List(jen.ID("loginValid"), jen.ID("err")).Op(":=").ID("s").Dot(
			"authenticator",
		).Dot(
			"ValidateLogin",
		).Call(jen.ID("ctx"), jen.ID("user").Dot(
			"HashedPassword",
		), jen.ID("loginInput").Dot(
			"Password",
		), jen.ID("user").Dot(
			"TwoFactorSecret",
		), jen.ID("loginInput").Dot(
			"TOTPToken",
		), jen.ID("user").Dot(
			"Salt",
		)),
		jen.If(jen.ID("err").Op("==").ID("auth").Dot(
			"ErrPasswordHashTooWeak",
		).Op("&&").ID("loginValid")).Block(
			jen.ID("logger").Dot(
				"Debug",
			).Call(jen.Lit("hashed password was deemed to weak, updating its hash")),
			jen.List(jen.ID("updated"), jen.ID("hashErr")).Op(":=").ID("s").Dot(
				"authenticator",
			).Dot(
				"HashPassword",
			).Call(jen.ID("ctx"), jen.ID("loginInput").Dot(
				"Password",
			)),
			jen.If(jen.ID("hashErr").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("false"), jen.Qual("fmt", "Errorf").Call(jen.Lit("updating password hash: %w"), jen.ID("hashErr"))),
			),
			jen.ID("user").Dot(
				"HashedPassword",
			).Op("=").ID("updated"),
			jen.If(jen.ID("updateErr").Op(":=").ID("s").Dot(
				"userDB",
			).Dot(
				"UpdateUser",
			).Call(jen.ID("ctx"), jen.ID("user")), jen.ID("updateErr").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("false"), jen.Qual("fmt", "Errorf").Call(jen.Lit("saving updated password hash: %w"), jen.ID("updateErr"))),
			),
		).Else().If(jen.ID("err").Op("!=").ID("nil").Op("&&").ID("err").Op("!=").ID("auth").Dot(
			"ErrPasswordHashTooWeak",
		)).Block(
			jen.ID("logger").Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("issue validating login")),
			jen.Return().List(jen.ID("false"), jen.Qual("fmt", "Errorf").Call(jen.Lit("validating login: %w"), jen.ID("err"))),
		),
		jen.Return().List(jen.ID("loginValid"), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildAuthCookie returns an authentication cookie for a given user").Params(jen.ID("s").Op("*").ID("Service")).ID("buildAuthCookie").Params(jen.ID("user").Op("*").ID("models").Dot(
		"User",
	)).Params(jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Block(
		jen.List(jen.ID("encoded"), jen.ID("err")).Op(":=").ID("s").Dot(
			"cookieManager",
		).Dot(
			"Encode",
		).Call(jen.ID("CookieName"), jen.ID("models").Dot(
			"CookieAuth",
		).Valuesln(jen.ID("UserID").Op(":").ID("user").Dot(
			"ID",
		), jen.ID("Admin").Op(":").ID("user").Dot(
			"IsAdmin",
		), jen.ID("Username").Op(":").ID("user").Dot(
			"Username",
		))),
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.ID("s").Dot(
				"logger",
			).Dot(
				"WithName",
			).Call(jen.ID("cookieErrorLogName")).Dot(
				"WithValue",
			).Call(jen.Lit("user_id"), jen.ID("user").Dot(
				"ID",
			)).Dot(
				"Error",
			).Call(jen.ID("err"), jen.Lit("error encoding cookie")),
			jen.Return().List(jen.ID("nil"), jen.ID("err")),
		),
		jen.Return().List(jen.ID("s").Dot(
			"buildCookie",
		).Call(jen.ID("encoded")), jen.ID("nil")),
	),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// buildCookie provides a consistent way of constructing an HTTP cookie").Params(jen.ID("s").Op("*").ID("Service")).ID("buildCookie").Params(jen.ID("value").ID("string")).Params(jen.Op("*").Qual("net/http", "Cookie")).Block(
		jen.Return().Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("CookieName"), jen.ID("Value").Op(":").ID("value"), jen.ID("Path").Op(":").Lit("/"), jen.ID("HttpOnly").Op(":").ID("true"), jen.ID("Secure").Op(":").ID("s").Dot(
			"config",
		).Dot(
			"SecureCookiesOnly",
		), jen.ID("Domain").Op(":").ID("s").Dot(
			"config",
		).Dot(
			"CookieDomain",
		), jen.ID("Expires").Op(":").Qual("time", "Now").Call().Dot(
			"Add",
		).Call(jen.ID("s").Dot(
			"config",
		).Dot(
			"CookieLifetime",
		))),
	),

		jen.Line(),
	)
	return ret
}
