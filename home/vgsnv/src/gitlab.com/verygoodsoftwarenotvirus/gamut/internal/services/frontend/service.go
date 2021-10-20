package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("serviceName").ID("string").Op("=").Lit("frontend_service"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AuthService").Interface(
				jen.ID("UserAttributionMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("PermissionFilterMiddleware").Params(jen.ID("permissions").Op("...").ID("authorization").Dot("Permission")).Params(jen.Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))),
				jen.ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("loginData").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Op("*").ID("types").Dot("User"), jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")),
				jen.ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("res").Qual("net/http", "ResponseWriter")).Params(jen.ID("error")),
			),
			jen.ID("UsersService").Interface(
				jen.ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("registrationInput").Op("*").ID("types").Dot("UserRegistrationInput")).Params(jen.Op("*").ID("types").Dot("UserCreationResponse"), jen.ID("error")),
				jen.ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Params(jen.ID("error")),
			),
			jen.ID("Service").Interface(jen.ID("SetupRoutes").Params(jen.ID("router").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/internal/routing", "Router"))),
			jen.ID("service").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("panicker").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/internal/panicking", "Panicker"),
				jen.ID("authService").ID("AuthService"),
				jen.ID("config").Op("*").ID("Config"),
				jen.ID("sessionContextDataFetcher").Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new Service."),
		jen.Newline(),
		jen.Func().ID("ProvideService").Params(jen.ID("cfg").Op("*").ID("Config"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("authService").ID("AuthService")).Params(jen.ID("Service")).Body(
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(jen.ID("config").Op(":").ID("cfg"), jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")), jen.ID("panicker").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/internal/panicking", "NewProductionPanicker").Call(), jen.ID("sessionContextDataFetcher").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/internal/services/authentication", "FetchContextFromRequest"), jen.ID("authService").Op(":").ID("authService")),
			jen.If(jen.ID("cfg").Dot("Debug")).Body(
				jen.ID("svc").Dot("logger").Dot("SetLevel").Call(jen.ID("logging").Dot("DebugLevel"))),
			jen.Return().ID("svc"),
		),
		jen.Newline(),
	)

	return code
}
