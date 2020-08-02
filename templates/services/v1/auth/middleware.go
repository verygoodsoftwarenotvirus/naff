package auth

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func middlewareDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildMiddlewareConstantDefs(proj)...)
	code.Add(buildCookieAuthenticationMiddleware(proj)...)
	code.Add(buildAuthenticationMiddleware(proj)...)
	code.Add(buildAdminMiddleware(proj)...)
	code.Add(buildparseLoginInputFromForm(proj)...)
	code.Add(buildUserLoginInputMiddleware(proj)...)

	return code
}

func buildMiddlewareConstantDefs(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			jen.Comment("userLoginInputMiddlewareCtxKey is the context key for login input."),
			jen.ID("userLoginInputMiddlewareCtxKey").Qual(proj.ModelsV1Package(), "ContextKey").Equals().Lit("user_login_input"),
			jen.Line(),
			jen.Comment("usernameFormKey is the string we look for in request forms for username information."),
			jen.ID("usernameFormKey").Equals().Lit("username"),
			jen.Comment("passwordFormKey is the string we look for in request forms for password information."),
			jen.ID("passwordFormKey").Equals().Lit("password"),
			jen.Comment("totpTokenFormKey is the string we look for in request forms for TOTP token information."),
			jen.ID("totpTokenFormKey").Equals().Lit("totpToken"),
		),
		jen.Line(),
	}

	return lines
}

func buildCookieAuthenticationMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("CookieAuthenticationMiddleware checks every request for a user cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("CookieAuthenticationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("CookieAuthenticationMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.Comment("fetch the user from the request."),
				jen.List(jen.ID("user"), jen.Err()).Assign().ID("s").Dot("fetchUserFromCookie").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered fetching user")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.If(jen.ID("user").DoesNotEqual().ID("nil")).Block(
					jen.ID(constants.RequestVarName).Equals().ID(constants.RequestVarName).Dot("WithContext").Callln(
						jen.Qual("context", "WithValue").Callln(
							jen.ID("ctx"),
							jen.Qual(proj.ModelsV1Package(), "SessionInfoKey"),
							jen.ID("user").Dot("ToSessionInfo").Call(),
						),
					),
					jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
					jen.Return(),
				),
				jen.Line(),
				jen.Comment("if no error was attached to the request, tell them to login first."),
				jen.Qual("net/http", "Redirect").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildAuthenticationMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("AuthenticationMiddleware authenticates based on either an oauth2 token or a cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("AuthenticationMiddleware").Params(jen.ID("allowValidCookieInLieuOfAValidToken").Bool()).Params(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).Block(
			jen.Return().Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
				jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
					jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("AuthenticationMiddleware")),
					jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
					jen.Line(),
					jen.Comment("let's figure out who the user is."),
					jen.Var().ID("user").PointerTo().Qual(proj.ModelsV1Package(), "User"),
					jen.Line(),
					jen.Comment("check for a cookie first if we can."),
					jen.If(jen.ID("allowValidCookieInLieuOfAValidToken")).Block(
						jen.List(jen.ID("cookieAuth"), jen.Err()).Assign().ID("s").Dot("DecodeCookieFromRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
						jen.Line(),
						jen.If(jen.Err().IsEqualTo().ID("nil").And().ID("cookieAuth").DoesNotEqual().ID("nil")).Block(
							jen.List(jen.ID("user"), jen.Err()).Equals().ID("s").Dot("userDB").Dot("GetUser").Call(constants.CtxVar(), jen.ID("cookieAuth").Dot("UserID")),
							jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
								jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error authenticating request")),
								jen.Qual("net/http", "Error").Call(jen.ID(constants.ResponseVarName), jen.Lit("fetching user"), jen.Qual("net/http", "StatusInternalServerError")),
								jen.Comment("if we get here, then we just don't have a valid cookie, and we need to move on."),
								jen.Return(),
							),
						),
					),
					jen.Line(),
					jen.Comment("if the cookie wasn't present, or didn't indicate who the user is."),
					jen.If(jen.ID("user").IsEqualTo().ID("nil")).Block(
						jen.Comment("check to see if there is an OAuth2 token for a valid client attached to the request."),
						jen.Comment("We do this first because it is presumed to be the primary means by which requests are made to the httpServer."),
						jen.List(jen.ID("oauth2Client"), jen.Err()).Assign().ID("s").Dot("oauth2ClientsService").Dot("ExtractOAuth2ClientFromRequest").Call(constants.CtxVar(), jen.ID(constants.RequestVarName)),
						jen.If(jen.Err().DoesNotEqual().ID("nil").Or().ID("oauth2Client").IsEqualTo().ID("nil")).Block(
							jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("fetching oauth2 client")),
							jen.Qual("net/http", "Redirect").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
							jen.Return(),
						),
						jen.Line(),
						jen.Comment("attach the oauth2 client and user's info to the request."),
						constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.Qual(proj.ModelsV1Package(), "OAuth2ClientKey"), jen.ID("oauth2Client")),
						jen.List(jen.ID("user"), jen.Err()).Equals().ID("s").Dot("userDB").Dot("GetUser").Call(constants.CtxVar(), jen.ID("oauth2Client").Dot(constants.UserOwnershipFieldName)),
						jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
							jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error authenticating request")),
							jen.Qual("net/http", "Error").Call(jen.ID(constants.ResponseVarName), jen.Lit("fetching user"), jen.Qual("net/http", "StatusInternalServerError")),
							jen.Return(),
						),
					),
					jen.Line(),
					jen.Comment("If your request gets here, you're likely either trying to get here, or desperately trying to get anywhere."),
					jen.If(jen.ID("user").IsEqualTo().ID("nil")).Block(
						jen.ID("s").Dot(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("no user attached to request request")),
						jen.Qual("net/http", "Redirect").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName), jen.Lit("/login"), jen.Qual("net/http", "StatusUnauthorized")),
						jen.Return(),
					),
					jen.Line(),
					jen.Comment("elsewise, load the request with extra context."),
					constants.CtxVar().Equals().Qual("context", "WithValue").Call(
						constants.CtxVar(),
						jen.Qual(proj.ModelsV1Package(), "SessionInfoKey"),
						jen.ID("user").Dot("ToSessionInfo").Call(),
					),
					jen.Line(),
					jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
				)),
			),
		),
		jen.Line(),
	}

	return lines
}

func buildAdminMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("AdminMiddleware restricts requests to admin users only."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("AdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("AdminMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID(constants.LoggerVarName).Assign().ID("s").Dot(constants.LoggerVarName).Dot("WithRequest").Call(jen.ID(constants.RequestVarName)),
				jen.List(jen.ID("si"), jen.ID("ok")).Assign().ID(constants.ContextVarName).Dot("Value").Call(jen.Qual(proj.ModelsV1Package(), "SessionInfoKey")).Assert(jen.PointerTo().Qual(proj.ModelsV1Package(), "SessionInfo")),
				jen.Line(),
				jen.If(jen.Not().ID("ok").Or().ID("si").IsEqualTo().ID("nil")).Block(
					jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("AdminMiddleware called without user attached to context")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.If(jen.Not().ID("si").Dot("UserIsAdmin")).Block(
					jen.ID(constants.LoggerVarName).Dot("Debug").Call(jen.Lit("AdminMiddleware called by non-admin user")),
					utils.WriteXHeader(constants.ResponseVarName, "StatusUnauthorized"),
					jen.Return(),
				),
				jen.Line(),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName)),
			)),
		),
		jen.Line(),
	}

	return lines
}

func buildparseLoginInputFromForm(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("parseLoginInputFromForm checks a request for a login form, and returns the parsed login data if relevant."),
		jen.Line(),
		jen.Func().ID("parseLoginInputFromForm").Params(jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), "UserLoginInput")).Block(
			jen.If(jen.Err().Assign().ID(constants.RequestVarName).Dot("ParseForm").Call(), jen.Err().IsEqualTo().ID("nil")).Block(
				jen.ID("uli").Assign().AddressOf().Qual(proj.ModelsV1Package(), "UserLoginInput").Valuesln(
					jen.ID("Username").MapAssign().ID(constants.RequestVarName).Dot("FormValue").Call(jen.ID("usernameFormKey")),
					jen.ID("Password").MapAssign().ID(constants.RequestVarName).Dot("FormValue").Call(jen.ID("passwordFormKey")),
					jen.ID("TOTPToken").MapAssign().ID(constants.RequestVarName).Dot("FormValue").Call(jen.ID("totpTokenFormKey")),
				),
				jen.Line(),
				jen.If(jen.ID("uli").Dot("Username").DoesNotEqual().EmptyString().And().ID("uli").Dot("Password").DoesNotEqual().EmptyString().And().ID("uli").Dot("TOTPToken").DoesNotEqual().EmptyString()).Block(
					jen.Return().ID("uli"),
				),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}

	return lines
}

func buildUserLoginInputMiddleware(proj *models.Project) []jen.Code {
	lines := []jen.Code{
		jen.Comment("UserLoginInputMiddleware fetches user login input from requests."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").PointerTo().ID("Service")).ID("UserLoginInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Block(
			jen.Return().Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID(constants.ResponseVarName).Qual("net/http", "ResponseWriter"), jen.ID(constants.RequestVarName).PointerTo().Qual("net/http", "Request")).Block(
				jen.List(constants.CtxVar(), jen.ID(constants.SpanVarName)).Assign().Qual(proj.InternalTracingV1Package(), "StartSpan").Call(jen.ID(constants.RequestVarName).Dot("Context").Call(), jen.Lit("UserLoginInputMiddleware")),
				jen.Defer().ID(constants.SpanVarName).Dot("End").Call(),
				jen.Line(),
				jen.ID("x").Assign().ID("new").Call(jen.Qual(proj.ModelsV1Package(), "UserLoginInput")),
				jen.If(jen.Err().Assign().ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(jen.ID(constants.RequestVarName), jen.ID("x")), jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.If(jen.ID("x").Equals().ID("parseLoginInputFromForm").Call(jen.ID(constants.RequestVarName)), jen.ID("x").IsEqualTo().ID("nil")).Block(
						jen.ID("s").Dot(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("error encountered decoding request body")),
						utils.WriteXHeader(constants.ResponseVarName, "StatusBadRequest"),
						jen.Return(),
					),
				),
				jen.Line(),
				constants.CtxVar().Equals().Qual("context", "WithValue").Call(constants.CtxVar(), jen.ID("userLoginInputMiddlewareCtxKey"), jen.ID("x")),
				jen.ID("next").Dot("ServeHTTP").Call(jen.ID(constants.ResponseVarName), jen.ID(constants.RequestVarName).Dot("WithContext").Call(constants.CtxVar())),
			)),
		),
		jen.Line(),
	}

	return lines
}
