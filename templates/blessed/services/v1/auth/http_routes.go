package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("auth")

	utils.AddImports(proj, code)

	code.Add(
		jen.Const().Defs(
			jen.Comment("CookieName is the name of the cookie we attach to requests."),
			jen.ID("CookieName").Equals().Lit("todocookie"),
			jen.ID("cookieErrorLogName").Equals().Lit("_COOKIE_CONSTRUCTION_ERROR_"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DecodeCookieFromRequest takes a request object and fetches the cookie data if it is present."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("DecodeCookieFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.ID("ca").PointerTo().Qual(proj.ModelsV1Package(), "CookieAuth"), jen.Err().Error()).Block(
			jen.List(jen.Underscore(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("DecodeCookieFromRequest")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.List(jen.ID("cookie"), jen.Err()).Assign().ID(constants.RequestVarName).Dot("Cookie").Call(jen.ID("CookieName")),
			jen.If(jen.Err().DoesNotEqual().Qual("net/http", "ErrNoCookie").And().ID("cookie").DoesNotEqual().ID("nil")).Block(
				jen.ID("decodeErr").Assign().ID("s").Dot("cookieManager").Dot("Decode").Call(jen.ID("CookieName"), jen.ID("cookie").Dot("Value"), jen.AddressOf().ID("ca")),
				jen.If(jen.ID("decodeErr").DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("decoding request cookie")),
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("decoding request cookie: %w"), jen.ID("decodeErr"))),
				),
				jen.Line(),
				jen.Return().List(jen.ID("ca"), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.Nil(), jen.Qual("net/http", "ErrNoCookie")),
		),
		jen.Line(),
	)

	// if proj.EnableNewsman {
	code.Add(
		jen.Comment("WebsocketAuthFunction is provided to Newsman to determine if a user has access to websockets."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("WebsocketAuthFunction").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.Bool()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("WebsocketAuthFunction")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Comment("First we check to see if there is an OAuth2 token for a valid client attached to the request."),
			jen.Comment("We do this first because it is presumed to be the primary means by which requests are made to the httpServer."),
			jen.List(jen.ID("oauth2Client"), jen.Err()).Assign().ID("s").Dot("oauth2ClientsService").Dot("ExtractOAuth2ClientFromRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
			jen.If(jen.Err().IsEqualTo().ID("nil").And().ID("oauth2Client").DoesNotEqual().ID("nil")).Block(
				jen.Return().True(),
			),
			jen.Line(),
			jen.Comment("In the event there's not a valid OAuth2 token attached to the request, or there is some other OAuth2 issue,"),
			jen.Comment("we next check to see if a valid cookie is attached to the request."),
			jen.List(jen.ID("cookieAuth"), jen.ID("cookieErr")).Assign().ID("s").Dot("DecodeCookieFromRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
			jen.If(jen.ID("cookieErr").IsEqualTo().ID("nil").And().ID("cookieAuth").DoesNotEqual().ID("nil")).Block(
				jen.Return().True(),
			),
			jen.Line(),
			jen.Comment("If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere."),
			jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error authenticated token-authenticated request")),
			jen.Return().False(),
		),
		jen.Line(),
	)
	// }

	code.Add(
		jen.Comment("FetchUserFromRequest takes a request object and fetches the cookie, and then the user for that cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("FetchUserFromRequest").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"), jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("FetchUserFromRequest")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.List(jen.ID("ca"), jen.ID("decodeErr")).Assign().ID("s").Dot("DecodeCookieFromRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
			jen.If(jen.ID("decodeErr").DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching cookie data from request: %w"), jen.ID("decodeErr"))),
			),
			jen.Line(),
			jen.List(jen.ID("user"), jen.ID("userFetchErr")).Assign().ID("s").Dot("userDB").Dot("GetUser").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.ID("ca").Dot("UserID")),
			jen.If(jen.ID("userFetchErr").DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("fetching user from request: %w"), jen.ID("userFetchErr"))),
			),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("ca").Dot("UserID")),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("ca").Dot("Username")),
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.Nil()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LoginHandler is our login route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("LoginHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("LoginHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.List(jen.ID("loginData"), jen.ID("errRes")).Assign().ID("s").Dot("fetchLoginDataFromRequest").Call(jen.ID(constants.RequestVarName)),
				jen.If(jen.ID("errRes").DoesNotEqual().ID("nil").Or().ID("loginData").IsEqualTo().Nil()).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.ID("errRes"), jen.Lit("error encountered fetching login data from request")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("errRes")), jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
					),
					jen.Return(),
				),
				jen.Line(),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("loginData").Dot("user").Dot("ID")),
				jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("loginData").Dot("user").Dot("Username")),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("user"), jen.ID("loginData").Dot("user").Dot("ID")),
				jen.List(jen.ID("loginValid"), jen.Err()).Assign().ID("s").Dot("validateLogin").Call(constants.CtxVar(), jen.PointerTo().ID("loginData")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered validating login")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("valid"), jen.ID("loginValid")),
				jen.Line(),
				jen.If(jen.Not().ID("loginValid")).Block(
					jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("login was invalid")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("login was valid")),
				jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("s").Dot("buildAuthCookie").Call(jen.ID("loginData").Dot("user")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error building cookie")),
					jen.Line(),
					utils.WriteXHeader(constants.ResponseVarName, "StatusInternalServerError"),
					jen.ID("response").Assign().AddressOf().Qual(proj.ModelsV1Package(), "ErrorResponse").Valuesln(
						jen.ID("Code").MapAssign().Qual("net/http", "StatusInternalServerError"),
						jen.ID("Message").MapAssign().Lit("error encountered building cookie"),
					),
					jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(jen.ID(constants.ResponseVarName), jen.ID("response")), jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
					),
					jen.Return(),
				),
				jen.Line(),
				jen.Qual("net/http", "SetCookie").Call(jen.ID(constants.ResponseVarName), jen.ID("cookie")),
				utils.WriteXHeader(constants.ResponseVarName, "StatusNoContent"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogoutHandler is our logout route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("LogoutHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(jen.Underscore(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("LogoutHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.Line(),
				jen.If(jen.List(jen.ID("cookie"), jen.Err()).Assign().ID(constants.RequestVarName).Dot("Cookie").Call(jen.ID("CookieName")), jen.Err().IsEqualTo().ID("nil").And().ID("cookie").DoesNotEqual().ID("nil")).Block(
					jen.ID("c").Assign().ID("s").Dot("buildCookie").Call(jen.Lit("deleted")),
					jen.ID("c").Dot("Expires").Equals().Qual("time", "Time").Values(),
					jen.ID("c").Dot("MaxAge").Equals().Lit(-1),
					jen.Qual("net/http", "SetCookie").Call(jen.ID(constants.ResponseVarName), jen.ID("c")),
				).Else().Block(
					jen.ID(constants.LoggerVarName).Dot("WithError").Call(jen.Err()).Dot("Debug").Call(jen.Lit("logout was called, no cookie was found")),
				),
				jen.Line(),
				utils.WriteXHeader(constants.ResponseVarName, "StatusOK"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CycleSecretHandler rotates the cookie building secret with a new random seccode."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CycleSecretHandler").Params().Params(jen.Qual("net/http", "HandlerFunc")).Block(
			jen.Return().Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(jen.Underscore(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CycleSecretHandler")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.ID(constants.LoggerVarName).Dot("Info").Call(jen.Lit("cycling cookie secret!")),
				jen.Line(),
				jen.ID("s").Dot("cookieManager").Equals().Qual("github.com/gorilla/securecookie", "New").Callln(
					jen.Qual("github.com/gorilla/securecookie", "GenerateRandomKey").Call(jen.Lit(64)),
					jen.Index().Byte().Call(jen.ID("s").Dot("config").Dot("CookieSecret")),
				),
				jen.Line(),
				utils.WriteXHeader(constants.ResponseVarName, "StatusCreated"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("loginData").Struct(jen.ID("loginInput").PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput"), jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("fetchLoginDataFromRequest searches a given HTTP request for parsed login input data, and"),
		jen.Line(),
		jen.Comment("returns a helper struct with the relevant login information."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("fetchLoginDataFromRequest").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().ID("loginData"), jen.PointerTo().Qual(proj.ModelsV1Package(), "ErrorResponse")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("fetchLoginDataFromRequest")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.List(jen.ID("loginInput"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("UserLoginInputMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput")),
			jen.If(jen.Not().ID("ok")).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no UserLoginInput found for /login request")),
				jen.Return().List(jen.Nil(), jen.AddressOf().Qual(proj.ModelsV1Package(), "ErrorResponse").Valuesln(
					jen.ID("Code").MapAssign().Qual("net/http", "StatusUnauthorized")),
				),
			),
			jen.Line(),
			jen.ID("username").Assign().ID("loginInput").Dot("Username"),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUsernameToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("username")),
			jen.Line(),
			jen.Comment("you could ensure there isn't an unsatisfied password reset token"),
			jen.Comment("requested before allowing login here."),
			jen.Line(),
			jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("userDB").Dot("GetUserByUsername").Call(constants.CtxVar(), jen.ID("username")),
			jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("no matching user")),
				jen.Return().List(jen.Nil(), jen.AddressOf().Qual(proj.ModelsV1Package(), "ErrorResponse").Values(jen.ID("Code").MapAssign().Qual("net/http", "StatusBadRequest"))),
			).Else().If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error fetching user")),
				jen.Return().List(jen.Nil(), jen.AddressOf().Qual(proj.ModelsV1Package(), "ErrorResponse").Values(jen.ID("Code").MapAssign().Qual("net/http", "StatusInternalServerError"))),
			),
			jen.Qual(proj.InternalTracingV1Package(), "AttachUserIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("user").Dot("ID")),
			jen.Line(),
			jen.ID("ld").Assign().AddressOf().ID("loginData").Valuesln(
				jen.ID("loginInput").MapAssign().ID("loginInput"),
				jen.ID("user").MapAssign().ID("user"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("ld"), jen.Nil()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("validateLogin takes login information and returns whether or not the login is valid."),
		jen.Line(),
		jen.Comment("In the event that there's an error, this function will return false and the error."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("validateLogin").Params(constants.CtxParam(), jen.ID("loginInfo").ID("loginData")).Params(jen.Bool(), jen.Error()).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(constants.CtxVar(), jen.Lit("validateLogin")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.Comment("alias the relevant data."),
			jen.List(jen.ID("user"), jen.ID("loginInput")).Assign().List(jen.ID("loginInfo").Dot("user"), jen.ID("loginInfo").Dot("loginInput")),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithValue").Call(jen.Lit("username"), jen.ID("user").Dot("Username")),
			jen.Line(),
			jen.Comment("check for login validity."),
			jen.List(jen.ID("loginValid"), jen.Err()).Assign().ID("s").Dot("authenticator").Dot("ValidateLogin").Callln(
				constants.CtxVar(),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("loginInput").Dot("Password"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("loginInput").Dot("TOTPToken"),
				jen.ID("user").Dot("Salt"),
			),
			jen.Line(),
			jen.Comment("if the login is otherwise valid, but the password is too weak, try to rehash it."),
			jen.If(jen.Err().IsEqualTo().Qual(proj.InternalAuthV1Package(), "ErrPasswordHashTooWeak").And().ID("loginValid")).Block(
				jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("hashed password was deemed to weak, updating its hash")),
				jen.Line(),
				jen.Comment("re-hash the password"),
				jen.List(jen.ID("updated"), jen.ID("hashErr")).Assign().ID("s").Dot("authenticator").Dot("HashPassword").Call(constants.CtxVar(), jen.ID("loginInput").Dot("Password")),
				jen.If(jen.ID("hashErr").DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.False(), jen.Qual("fmt", "Errorf").Call(jen.Lit("updating password hash: %w"), jen.ID("hashErr"))),
				),
				jen.Line(),
				jen.Comment("update stored hashed password in the database."),
				jen.ID("user").Dot("HashedPassword").Equals().ID("updated"),
				jen.If(jen.ID("updateErr").Assign().ID("s").Dot("userDB").Dot("UpdateUser").Call(constants.CtxVar(), jen.ID("user")), jen.ID("updateErr").DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.False(), jen.Qual("fmt", "Errorf").Call(jen.Lit("saving updated password hash: %w"), jen.ID("updateErr"))),
				),
				jen.Line(),
				jen.Return().List(jen.ID("loginValid"), jen.Nil()),
			).Else().If(jen.Err().DoesNotEqual().ID("nil").And().Err().DoesNotEqual().Qual(proj.InternalAuthV1Package(), "ErrPasswordHashTooWeak")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("issue validating login")),
				jen.Return().List(jen.False(), jen.Qual("fmt", "Errorf").Call(jen.Lit("validating login: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("loginValid"), jen.Err()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildAuthCookie returns an authentication cookie for a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("buildAuthCookie").Params(jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User")).Params(jen.PointerTo().Qual("net/http", "Cookie"), jen.Error()).Block(
			jen.Comment("NOTE: code here is duplicated into the unit tests for"),
			jen.Comment("DecodeCookieFromRequest any changes made here might need"),
			jen.Comment("to be reflected there."),
			jen.List(jen.ID("encoded"), jen.Err()).Assign().ID("s").Dot("cookieManager").Dot("Encode").Callln(
				jen.ID("CookieName"),
				jen.Qual(proj.ModelsV1Package(), "CookieAuth").Valuesln(
					jen.ID(constants.UserIDVarName).MapAssign().ID("user").Dot("ID"),
					jen.ID("Admin").MapAssign().ID("user").Dot("IsAdmin"),
					jen.ID("Username").MapAssign().ID("user").Dot("Username"),
				),
			),
			jen.Line(),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Comment("NOTE: these errors should be infrequent, and should cause alarm when they do occur"),
				jen.ID("s").Dot(constants.LoggerVarName).Dot("WithName").Call(jen.ID("cookieErrorLogName")).Dot("WithValue").Call(jen.Lit("user_id"), jen.ID("user").Dot("ID")).Dot("Error").Call(jen.Err(), jen.Lit("error encoding cookie")),
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Return().List(jen.ID("s").Dot("buildCookie").Call(jen.ID("encoded")), jen.Nil()),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildCookie provides a consistent way of constructing an HTTP cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("buildCookie").Params(jen.ID("value").String()).Params(jen.PointerTo().Qual("net/http", "Cookie")).Block(
			jen.Comment("https://www.calhoun.io/securing-cookies-in-go/"),
			jen.Return().AddressOf().Qual("net/http", "Cookie").Valuesln(
				jen.ID("Name").MapAssign().ID("CookieName"),
				jen.ID("Value").MapAssign().ID("value"),
				jen.ID("Path").MapAssign().Lit("/"),
				jen.ID("HttpOnly").MapAssign().True(),
				jen.ID("Secure").MapAssign().ID("s").Dot("config").Dot("SecureCookiesOnly"),
				jen.ID("Domain").MapAssign().ID("s").Dot("config").Dot("CookieDomain"),
				jen.ID("Expires").MapAssign().Qual("time", "Now").Call().Dot("Add").Call(jen.ID("s").Dot("config").Dot("CookieLifetime")),
			),
		),
		jen.Line(),
	)

	return code
}
