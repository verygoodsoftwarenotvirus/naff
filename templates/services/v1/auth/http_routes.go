package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildHTTPRoutesConstantDefs()...)
	code.Add(buildDecodeCookieFromRequest(proj)...)

	// if proj.EnableNewsman {
	code.Add(buildWebsocketAuthFunction(proj)...)
	// }

	code.Add(buildFetchUserFromCookie(proj)...)
	code.Add(buildLoginHandler(proj)...)
	code.Add(buildLogoutHandler(proj)...)
	code.Add(buildStatusHandler(proj)...)
	code.Add(buildCycleSecretHandler(proj)...)

	code.Add(
		jen.Type().ID("loginData").Struct(jen.ID("loginInput").PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput"), jen.ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User")),
		jen.Line(),
	)

	code.Add(buildFetchLoginDataFromRequest(proj)...)
	code.Add(buildValidateLogin(proj)...)
	code.Add(buildBuildCookie()...)

	return code
}

func buildHTTPRoutesConstantDefs() []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("CookieName is the name of the cookie we attach to requests."),
			jen.ID("CookieName").Equals().Lit("todocookie"),
			jen.ID("cookieErrorLogName").Equals().Lit("_COOKIE_CONSTRUCTION_ERROR_"),
			jen.Line(),
			jen.ID("sessionInfoKey").Equals().Lit("session_info"),
		),
		jen.Line(),
	}

	return lines
}

func buildDecodeCookieFromRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("DecodeCookieFromRequest takes a request object and fetches the cookie data if it is present."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("DecodeCookieFromRequest").Params(
			constants.CtxParam(),
			jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request"),
		).Params(
			jen.ID("ca").PointerTo().Qual(proj.ModelsV1Package(), "SessionInfo"),
			jen.Err().Error(),
		).Block(
			utils.StartSpan(proj, true, "DecodeCookieFromRequest"),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.List(jen.ID("cookie"), jen.Err()).Assign().ID(constants.RequestVarName).Dot("Cookie").Call(jen.ID("CookieName")),
			jen.If(jen.Err().DoesNotEqual().Qual("net/http", "ErrNoCookie").And().ID("cookie").DoesNotEqual().ID("nil")).Block(
				jen.Var().ID("token").String(),
				jen.ID("decodeErr").Assign().ID("s").Dot("cookieManager").Dot("Decode").Call(jen.ID("CookieName"), jen.ID("cookie").Dot("Value"), jen.AddressOf().ID("token")),
				jen.If(jen.ID("decodeErr").DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("decoding request cookie")),
					jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("decoding request cookie: %w"), jen.ID("decodeErr"))),
				),
				jen.Line(),
				jen.Var().ID("sessionErr").Error(),
				jen.List(jen.ID(constants.ContextVarName), jen.ID("sessionErr")).Equals().ID("s").Dot("sessionManager").Dot("Load").Call(
					jen.ID(constants.ContextVarName),
					jen.ID("token"),
				),
				jen.If(jen.ID("sessionErr").DoesNotEqual().Nil()).Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(
						jen.ID("sessionErr"),
						jen.Lit("error loading token"),
					),
					jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("error loading token"))),
				),
				jen.Line(),
				jen.List(jen.ID("si"), jen.ID("ok")).Assign().ID("s").Dot("sessionManager").Dot("Get").Call(
					jen.ID(constants.ContextVarName),
					jen.ID("sessionInfoKey"),
				).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "SessionInfo")),
				jen.If(jen.Not().ID("ok")).Block(
					jen.ID("errToReturn").Assign().Qual("errors", "New").Call(
						jen.Lit("no session info attached to context"),
					),
					jen.ID(constants.LoggerVarName).Dot("Error").Call(
						jen.ID("errToReturn"),
						jen.Lit("fetching session data"),
					),
					jen.Return(jen.Nil(), jen.ID("errToReturn")),
				),
				jen.Line(),
				jen.Return().List(jen.ID("si"), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.Nil(), jen.Qual("net/http", "ErrNoCookie")),
		),
		jen.Line(),
	}

	return lines
}

func buildWebsocketAuthFunction(proj *models.Project) []jen.Code {
	lines := []jen.Code{
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
	}

	return lines
}

func buildFetchUserFromCookie(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("fetchUserFromCookie takes a request object and fetches the cookie, and then the user for that cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("fetchUserFromCookie").Params(constants.CtxParam(), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "User"), jen.Error()).Block(
			utils.StartSpan(proj, true, "fetchUserFromCookie"),
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
			jen.Line(),
			jen.Return().List(jen.ID("user"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}

func buildLoginHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("LoginHandler is our login route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("LoginHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
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
			jen.Var().ID("sessionErr").Error(),
			jen.List(jen.ID(constants.ContextVarName), jen.ID("sessionErr")).Equals().ID("s").Dot("sessionManager").Dot("Load").Call(
				jen.ID(constants.ContextVarName),
				jen.EmptyString(),
			),
			jen.If(jen.ID("sessionErr").DoesNotEqual().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.ID("sessionErr"),
					jen.Lit("error loading token"),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.Line(),
			jen.If(
				jen.ID("renewTokenErr").Assign().ID("s").Dot("sessionManager").Dot("RenewToken").Call(jen.ID(constants.ContextVarName)),
				jen.ID("renewTokenErr").DoesNotEqual().Nil(),
			).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("error encountered renewing token"),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.ID("s").Dot("sessionManager").Dot("Put").Call(
				constants.CtxVar(),
				jen.ID("sessionInfoKey"),
				jen.ID("loginData").Dot("user").Dot("ToSessionInfo").Call(),
			),
			jen.Line(),
			jen.List(
				jen.ID("token"),
				jen.ID("expiry"),
				jen.Err(),
			).Assign().ID("s").Dot("sessionManager").Dot("Commit").Call(
				constants.CtxVar(),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("error encountered writing to session store"),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.Line(),
			jen.List(jen.ID("cookie"), jen.Err()).Assign().ID("s").Dot("buildCookie").Call(jen.ID("token"), jen.ID("expiry")),
			jen.If(jen.Err().DoesNotEqual().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("error encountered building cookie"),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.Line(),
			jen.Qual("net/http", "SetCookie").Call(jen.ID(constants.ResponseVarName), jen.ID("cookie")),
			utils.WriteXHeader(constants.ResponseVarName, "StatusNoContent"),
		),
		jen.Line(),
	}

	return lines
}

func buildLogoutHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("LogoutHandler is our logout route."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("LogoutHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("LogoutHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.List(constants.CtxVar(), jen.ID("sessionErr")).Assign().ID("s").Dot("sessionManager").Dot("Load").Call(
				constants.CtxVar(),
				jen.EmptyString(),
			),
			jen.If(jen.ID("sessionErr").DoesNotEqual().Nil()).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.ID("sessionErr"),
					jen.Lit("error loading token"),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.Line(),
			jen.If(
				jen.Err().Assign().ID("s").Dot("sessionManager").Dot("Clear").Call(constants.CtxVar()),
				jen.Err().DoesNotEqual().Nil(),
			).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(
					jen.Err(),
					jen.Lit("clearing user session"),
				),
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
					jen.Qual("net/http", "StatusInternalServerError"),
				),
				jen.Return(),
			),
			jen.Line(),
			jen.If(
				jen.List(jen.ID("cookie"), jen.ID("cookieRetrievalErr")).Assign().ID(constants.RequestVarName).Dot("Cookie").Call(jen.ID("CookieName")),
				jen.ID("cookieRetrievalErr").IsEqualTo().Nil().And().ID("cookie").DoesNotEqual().Nil(),
			).Block(
				jen.If(
					jen.List(jen.ID("c"), jen.ID("cookieBuildingErr")).Assign().ID("s").Dot("buildCookie").Call(
						jen.Lit("deleted"),
						jen.Qual("time", "Time").Values(),
					),
					jen.ID("cookieBuildingErr").IsEqualTo().Nil().And().ID("c").DoesNotEqual().Nil(),
				).Block(
					jen.ID("c").Dot("MaxAge").Equals().Lit(-1),
					jen.Qual("net/http", "SetCookie").Call(jen.ID(constants.ResponseVarName), jen.ID("c")),
				).Else().Block(
					jen.ID(constants.LoggerVarName).Dot("Error").Call(
						jen.ID("cookieBuildingErr"),
						jen.Lit("error encountered building cookie"),
					),
					jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(
						jen.Qual("net/http", "StatusInternalServerError"),
					),
					jen.Return(),
				),
			).Else().Block(
				jen.ID(constants.LoggerVarName).Dot("WithError").Call(jen.ID("cookieRetrievalErr")).Dot("Debug").Call(jen.Lit("logout was called, no cookie was found")),
			),
			jen.Line(),
			utils.WriteXHeader(constants.ResponseVarName, "StatusOK"),
		),
		jen.Line(),
	}

	return lines
}

func buildStatusHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("StatusHandler returns the user info for the user making the request."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("StatusHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
			jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("StatusHandler")),
			jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
			jen.Line(),
			jen.ID(constants.LoggerVarName).Assign().ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
			jen.Line(),
			jen.Var().ID("sr").PointerTo().Qual(proj.ModelsV1Package(), "StatusResponse"),
			jen.List(jen.ID("userInfo"), jen.Err()).Assign().ID("s").Dot("fetchUserFromCookie").Call(
				constants.CtxVar(),
				jen.ID(constants.RequestVarName),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Block(
				jen.ID(constants.ResponseVarName).Dot("WriteHeader").Call(jen.Qual("net/http", "StatusUnauthorized")),
				jen.ID("sr").Equals().AddressOf().Qual(proj.ModelsV1Package(), "StatusResponse").Valuesln(
					jen.ID("Authenticated").MapAssign().False(),
					jen.ID("IsAdmin").MapAssign().False(),
				),
			).Else().Block(
				jen.ID("sr").Equals().AddressOf().Qual(proj.ModelsV1Package(), "StatusResponse").Valuesln(
					jen.ID("Authenticated").MapAssign().True(),
					jen.ID("IsAdmin").MapAssign().ID("userInfo").Dot("IsAdmin"),
				),
			),
			jen.Line(),
			jen.If(
				jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("EncodeResponse").Call(
					jen.ID(constants.ResponseVarName),
					jen.ID("sr"),
				),
				jen.Err().DoesNotEqual().Nil(),
			).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("encoding response")),
			),
		),
	}

	return lines
}

func buildCycleSecretHandler(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CycleSecretHandler rotates the cookie building secret with a new random secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CycleSecretHandler").Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
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
		jen.Line(),
	}

	return lines
}

func buildFetchLoginDataFromRequest(proj *models.Project) []jen.Code {
	lines := []jen.Code{
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
			jen.List(jen.ID("loginInput"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.ID("userLoginInputMiddlewareCtxKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput")),
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
	}

	return lines
}

func buildValidateLogin(proj *models.Project) []jen.Code {
	lines := []jen.Code{
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
			jen.If(jen.Err().IsEqualTo().Qual(proj.InternalAuthV1Package(), "ErrCostTooLow").And().ID("loginValid")).Block(
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
			).Else().If(jen.Err().DoesNotEqual().ID("nil").And().Err().DoesNotEqual().Qual(proj.InternalAuthV1Package(), "ErrCostTooLow")).Block(
				jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("issue validating login")),
				jen.Return().List(jen.False(), jen.Qual("fmt", "Errorf").Call(jen.Lit("validating login: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("loginValid"), jen.Err()),
		),
		jen.Line(),
	}

	return lines
}

func buildBuildCookie() []jen.Code {
	lines := []jen.Code{
		jen.Comment("buildCookie provides a consistent way of constructing an HTTP cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("buildCookie").Params(
			jen.ID("value").String(),
			jen.ID("expiry").Qual("time", "Time"),
		).Params(
			jen.PointerTo().Qual("net/http", "Cookie"),
			jen.Error(),
		).Block(
			jen.List(jen.ID("encoded"), jen.Err()).Assign().ID("s").Dot("cookieManager").Dot("Encode").Call(
				jen.ID("CookieName"),
				jen.ID("value"),
			),
			jen.If(jen.Err().DoesNotEqual().Nil()).Block(
				jen.Comment("NOTE: these errors should be infrequent, and should cause alarm when they do occur"),
				jen.ID("s").Dot("logger").Dot("WithName").Call(jen.ID("cookieErrorLogName")).Dot("Error").Call(
					jen.Err(),
					jen.Lit("error encoding cookie"),
				),
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.Comment("https://www.calhoun.io/securing-cookies-in-go/"),
			jen.ID("cookie").Assign().AddressOf().Qual("net/http", "Cookie").Valuesln(
				jen.ID("Name").MapAssign().ID("CookieName"),
				jen.ID("Value").MapAssign().ID("encoded"),
				jen.ID("Path").MapAssign().Lit("/"),
				jen.ID("HttpOnly").MapAssign().True(),
				jen.ID("Secure").MapAssign().ID("s").Dot("config").Dot("SecureCookiesOnly"),
				jen.ID("Domain").MapAssign().ID("s").Dot("config").Dot("CookieDomain"),
				jen.ID("Expires").MapAssign().ID("expiry"),
			),
			jen.Line(),
			jen.Return(jen.ID("cookie"), jen.Nil()),
		),
		jen.Line(),
	}

	return lines
}
