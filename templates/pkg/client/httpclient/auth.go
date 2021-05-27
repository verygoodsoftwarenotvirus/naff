package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("UserStatus fetches a user's status."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UserStatus").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").ID("types").Dot("UserStatusResponse"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildUserStatusRequest").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Var().Defs(
				jen.ID("output").Op("*").ID("types").Dot("UserStatusResponse"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("output"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving user status"),
				))),
			jen.Return().List(jen.ID("output"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BeginSession fetches a login cookie."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("BeginSession").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("input").Dot("Username"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildLoginRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building login request"),
				))),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("unauthenticatedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing login request"),
				))),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.ID("cookies").Op(":=").ID("res").Dot("Cookies").Call(),
			jen.If(jen.ID("len").Call(jen.ID("cookies")).Op(">").Lit(0)).Body(
				jen.Return().List(jen.ID("cookies").Index(jen.Lit(0)), jen.ID("nil"))),
			jen.Return().List(jen.ID("nil"), jen.ID("ErrNoCookiesReturned")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EndSession logs a user out."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("EndSession").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildLogoutRequest").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building logout request"),
				)),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("authedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing logout request"),
				)),
			jen.ID("c").Dot("authedClient").Dot("Transport").Op("=").ID("newDefaultRoundTripper").Call(jen.ID("c").Dot("authedClient").Dot("Timeout")),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ChangePassword changes a user's password."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ChangePassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.ID("input").Op("*").ID("types").Dot("PasswordUpdateInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("cookie").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrCookieRequired")),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildChangePasswordRequest").Call(
				jen.ID("ctx"),
				jen.ID("cookie"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building change password request"),
				)),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("unauthenticatedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("changing password"),
				)),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.If(jen.ID("res").Dot("StatusCode").Op("!=").Qual("net/http", "StatusOK")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("errInvalidResponseCode"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("invalid response code: %d"),
					jen.ID("res").Dot("StatusCode"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CycleTwoFactorSecret cycles a user's 2FA secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CycleTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cookie").Op("*").Qual("net/http", "Cookie"), jen.ID("input").Op("*").ID("types").Dot("TOTPSecretRefreshInput")).Params(jen.Op("*").ID("types").Dot("TOTPSecretRefreshResponse"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("cookie").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrCookieRequired"))),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildCycleTwoFactorSecretRequest").Call(
				jen.ID("ctx"),
				jen.ID("cookie"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building cycle two factor secret request"),
				))),
			jen.Var().Defs(
				jen.ID("output").Op("*").ID("types").Dot("TOTPSecretRefreshResponse"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("output"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("cycling two factor secret"),
				))),
			jen.Return().List(jen.ID("output"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("VerifyTOTPSecret verifies a 2FA secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("VerifyTOTPSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("token").ID("string")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("strconv", "ParseUint").Call(
				jen.ID("token"),
				jen.Lit(10),
				jen.Lit(64),
			), jen.ID("token").Op("==").Lit("").Op("||").ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("invalid token provided: %q"),
					jen.ID("token"),
				)),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildVerifyTOTPSecretRequest").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("token"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building verify two factor secret request"),
				)),
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID("c").Dot("fetchResponseToRequest").Call(
				jen.ID("ctx"),
				jen.ID("c").Dot("unauthenticatedClient"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("verifying two factor secret"),
				)),
			jen.ID("c").Dot("closeResponseBody").Call(
				jen.ID("ctx"),
				jen.ID("res"),
			),
			jen.If(jen.ID("res").Dot("StatusCode").Op("==").Qual("net/http", "StatusBadRequest")).Body(
				jen.Return().ID("ErrInvalidTOTPToken")).Else().If(jen.ID("res").Dot("StatusCode").Op("!=").Qual("net/http", "StatusAccepted")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("errInvalidResponseCode"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("erroneous response code when validating TOTP secret: %d"),
					jen.ID("res").Dot("StatusCode"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
