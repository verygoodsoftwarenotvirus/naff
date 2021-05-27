package authentication

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("errNoUserIDFoundInSession").Op("=").Qual("errors", "New").Call(jen.Lit("no user ID found in session")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("overrideSessionContextDataValuesWithSessionData").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.List(jen.ID("activeAccount"), jen.ID("ok")).Op(":=").ID("s").Dot("sessionManager").Dot("Get").Call(
				jen.ID("ctx"),
				jen.ID("accountIDContextKey"),
			).Assert(jen.ID("uint64")), jen.ID("ok")).Body(
				jen.ID("sessionCtxData").Dot("ActiveAccountID").Op("=").ID("activeAccount")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("getUserIDFromCookie takes a request object and fetches the cookie data if it is present."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("getUserIDFromCookie").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Qual("context", "Context"), jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.If(jen.List(jen.ID("cookie"), jen.ID("cookieErr")).Op(":=").ID("req").Dot("Cookie").Call(jen.ID("s").Dot("config").Dot("Cookies").Dot("Name")), jen.Op("!").Qual("errors", "Is").Call(
				jen.ID("cookieErr"),
				jen.Qual("net/http", "ErrNoCookie"),
			).Op("&&").ID("cookie").Op("!=").ID("nil")).Body(
				jen.Var().Defs(
					jen.ID("token").ID("string"),
					jen.ID("err").ID("error"),
				),
				jen.If(jen.ID("err").Op("=").ID("s").Dot("cookieManager").Dot("Decode").Call(
					jen.ID("s").Dot("config").Dot("Cookies").Dot("Name"),
					jen.ID("cookie").Dot("Value"),
					jen.Op("&").ID("token"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
						jen.Lit("cookie"),
						jen.ID("cookie").Dot("Value"),
					),
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("retrieving session context data"),
					)),
				),
				jen.List(jen.ID("ctx"), jen.ID("err")).Op("=").ID("s").Dot("sessionManager").Dot("Load").Call(
					jen.ID("ctx"),
					jen.ID("token"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("loading session"),
					))),
				jen.If(jen.List(jen.ID("userID"), jen.ID("ok")).Op(":=").ID("s").Dot("sessionManager").Dot("Get").Call(
					jen.ID("ctx"),
					jen.ID("userIDContextKey"),
				).Assert(jen.ID("uint64")), jen.ID("ok")).Body(
					jen.ID("logger").Dot("WithValue").Call(
						jen.ID("keys").Dot("UserIDKey"),
						jen.ID("userID"),
					).Dot("Debug").Call(jen.Lit("determined userID from request cookie")),
					jen.Return().List(jen.ID("ctx"), jen.ID("userID"), jen.ID("nil")),
				),
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("errNoUserIDFoundInSession"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("determining user ID from cookie"),
				)),
			),
			jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Qual("net/http", "ErrNoCookie")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("determineUserFromRequestCookie takes a request object and fetches the cookie, and then the user for that cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("determineUserFromRequestCookie").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")).Dot("WithValue").Call(
				jen.Lit("cookie_count"),
				jen.ID("len").Call(jen.ID("req").Dot("Cookies").Call()),
			),
			jen.List(jen.ID("ctx"), jen.ID("userID"), jen.ID("err")).Op(":=").ID("s").Dot("getUserIDFromCookie").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching cookie data from request"),
				))),
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("s").Dot("userDataManager").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching user from database"),
				))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("user determined from request cookie")),
			jen.Return().List(jen.ID("user"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("validateLogin takes login information and returns whether the login is valid."),
		jen.Line(),
		jen.Comment("In the event of an error, this function will return false and the error."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("validateLogin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("user").Op("*").ID("types").Dot("User"), jen.ID("loginInput").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("user").Dot("Username"),
			),
			jen.List(jen.ID("loginValid"), jen.ID("err")).Op(":=").ID("s").Dot("authenticator").Dot("ValidateLogin").Call(
				jen.ID("ctx"),
				jen.ID("user").Dot("HashedPassword"),
				jen.ID("loginInput").Dot("Password"),
				jen.ID("user").Dot("TwoFactorSecret"),
				jen.ID("loginInput").Dot("TOTPToken"),
			),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.ID("authentication").Dot("ErrInvalidTOTPToken"),
			).Op("||").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.ID("authentication").Dot("ErrPasswordDoesNotMatch"),
			)).Body(
				jen.Return().List(jen.ID("false"), jen.ID("err"))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating login"),
				))),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("login validated")),
			jen.Return().List(jen.ID("loginValid"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildCookie provides a consistent way of constructing an HTTP cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("buildCookie").Params(jen.ID("value").ID("string"), jen.ID("expiry").Qual("time", "Time")).Params(jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Body(
			jen.List(jen.ID("encoded"), jen.ID("err")).Op(":=").ID("s").Dot("cookieManager").Dot("Encode").Call(
				jen.ID("s").Dot("config").Dot("Cookies").Dot("Name"),
				jen.ID("value"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("s").Dot("logger").Dot("WithName").Call(jen.ID("cookieErrorLogName")).Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("error encoding cookie"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.ID("cookie").Op(":=").Op("&").Qual("net/http", "Cookie").Valuesln(jen.ID("Name").Op(":").ID("s").Dot("config").Dot("Cookies").Dot("Name"), jen.ID("Value").Op(":").ID("encoded"), jen.ID("Path").Op(":").Lit("/"), jen.ID("HttpOnly").Op(":").ID("true"), jen.ID("Secure").Op(":").ID("s").Dot("config").Dot("Cookies").Dot("SecureOnly"), jen.ID("Domain").Op(":").ID("s").Dot("config").Dot("Cookies").Dot("Domain"), jen.ID("Expires").Op(":").ID("expiry"), jen.ID("SameSite").Op(":").Qual("net/http", "SameSiteStrictMode"), jen.ID("MaxAge").Op(":").ID("int").Call(jen.Qual("time", "Until").Call(jen.ID("expiry")).Dot("Seconds").Call())),
			jen.Return().List(jen.ID("cookie"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
