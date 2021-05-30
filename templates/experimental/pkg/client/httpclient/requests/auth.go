package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("authBasePath").Op("=").Lit("auth"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserStatusRequest builds an HTTP request that fetches a user's status."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildUserStatusRequest").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("b").Dot("logger"),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("authBasePath"),
				jen.Lit("status"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildLoginRequest builds an HTTP request that fetches a login cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildLoginRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("Username"),
			),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("login"),
			),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildLogoutRequest builds an HTTP request that clears the user's session."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildLogoutRequest").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("logout"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("b").Dot("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildChangePasswordRequest builds a request to change a user's password."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildChangePasswordRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.ID("input").Op("*").ID("types").Dot("PasswordUpdateInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("password"),
				jen.Lit("new"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPut"),
				jen.ID("uri"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCycleTwoFactorSecretRequest builds a request to change a user's 2FA secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildCycleTwoFactorSecretRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.ID("input").Op("*").ID("types").Dot("TOTPSecretRefreshInput")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("cookie").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrCookieRequired"))),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger"),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("totp_secret"),
				jen.Lit("new"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.ID("req").Dot("AddCookie").Call(jen.ID("cookie")),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildVerifyTOTPSecretRequest builds a request to validate a 2FA secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildVerifyTOTPSecretRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("token").ID("string")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("token"),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("token").Op("==").Lit("").Op("||").ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("invalid token provided"),
				))),
			jen.ID("uri").Op(":=").ID("b").Dot("buildUnversionedURL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("usersBasePath"),
				jen.Lit("totp_secret"),
				jen.Lit("verify"),
			),
			jen.Return().ID("b").Dot("buildDataRequest").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.ID("uri"),
				jen.Op("&").ID("types").Dot("TOTPSecretVerificationInput").Valuesln(jen.ID("TOTPToken").Op(":").ID("token"), jen.ID("UserID").Op(":").ID("userID")),
			),
		),
		jen.Line(),
	)

	return code
}
