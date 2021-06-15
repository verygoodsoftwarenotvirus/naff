package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
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
				jen.ID("PermissionFilterMiddleware").Params(jen.ID("permissions").Op("...").ID("authorization").Dot("Permission")).Params(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))),
				jen.ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("loginData").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Op("*").ID("types").Dot("User"), jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")),
				jen.ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("res").Qual("net/http", "ResponseWriter")).Params(jen.ID("error")),
			),
			jen.ID("UsersService").Interface(
				jen.ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("registrationInput").Op("*").ID("types").Dot("UserRegistrationInput")).Params(jen.Op("*").ID("types").Dot("UserCreationResponse"), jen.ID("error")),
				jen.ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("TOTPSecretVerificationInput")).Params(jen.ID("error")),
			),
			jen.ID("Service").Interface(jen.ID("SetupRoutes").Params(jen.ID("router").ID("routing").Dot("Router"))),
			jen.ID("service").Struct(
				jen.ID("useFakeData").ID("bool"),
				jen.ID("templateFuncMap").Qual("html/template", "FuncMap"),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("panicker").ID("panicking").Dot("Panicker"),
				jen.ID("localizer").Op("*").ID("i18n").Dot("Localizer"),
				jen.ID("dataStore").ID("database").Dot("DataManager"),
				jen.ID("paymentManager").ID("capitalism").Dot("PaymentManager"),
				jen.ID("authService").ID("AuthService"),
				jen.ID("usersService").ID("UsersService"),
				jen.ID("sessionContextDataFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")),
				jen.ID("webhookIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.ID("apiClientIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.ID("accountIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
				jen.ID("itemIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvideService builds a new Service."),
		jen.Newline(),
		jen.Func().ID("ProvideService").Params(jen.ID("cfg").Op("*").ID("Config"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("authService").ID("AuthService"), jen.ID("usersService").ID("UsersService"), jen.ID("dataStore").ID("database").Dot("DataManager"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager"), jen.ID("paymentManager").ID("capitalism").Dot("PaymentManager")).Params(jen.ID("Service")).Body(
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(jen.ID("useFakeData").Op(":").ID("cfg").Dot("UseFakeData"), jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")), jen.ID("panicker").Op(":").ID("panicking").Dot("NewProductionPanicker").Call(), jen.ID("localizer").Op(":").ID("provideLocalizer").Call(), jen.ID("sessionContextDataFetcher").Op(":").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/services/authentication", "FetchContextFromRequest"), jen.ID("authService").Op(":").ID("authService"), jen.ID("usersService").Op(":").ID("usersService"), jen.ID("paymentManager").Op(":").ID("paymentManager"), jen.ID("dataStore").Op(":").ID("dataStore"), jen.ID("apiClientIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("apiClientIDURLParamKey"),
				jen.Lit("API client"),
			), jen.ID("accountIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("accountIDURLParamKey"),
				jen.Lit("account"),
			), jen.ID("webhookIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("webhookIDURLParamKey"),
				jen.Lit("webhook"),
			), jen.ID("itemIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.ID("itemIDURLParamKey"),
				jen.Lit("item"),
			), jen.ID("templateFuncMap").Op(":").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("relativeTime").Op(":").ID("relativeTime"), jen.Lit("relativeTimeFromPtr").Op(":").ID("relativeTimeFromPtr"))),
			jen.ID("svc").Dot("templateFuncMap").Index(jen.Lit("translate")).Op("=").ID("svc").Dot("getSimpleLocalizedString"),
			jen.Return().ID("svc"),
		),
		jen.Newline(),
	)

	return code
}
