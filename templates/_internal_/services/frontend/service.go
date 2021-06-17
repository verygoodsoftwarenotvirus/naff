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
		jen.Const().Defs(
			jen.ID("serviceName").ID("string").Op("=").Lit("frontend_service"),
		),
		jen.Newline(),
	)

	structFields := []jen.Code{
		jen.ID("useFakeData").ID("bool"),
		jen.ID("templateFuncMap").Qual("html/template", "FuncMap"),
		jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
		jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
		jen.ID("panicker").Qual(proj.InternalPackage("panicking"), "Panicker"),
		jen.ID("localizer").Op("*").Qual("github.com/nicksnyder/go-i18n/v2/i18n", "Localizer"),
		jen.ID("dataStore").ID("database").Dot("DataManager"),
		jen.ID("paymentManager").ID("capitalism").Dot("PaymentManager"),
		jen.ID("authService").ID("AuthService"),
		jen.ID("usersService").ID("UsersService"),
		jen.ID("sessionContextDataFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.Op("*").Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")),
		jen.ID("accountIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		jen.ID("apiClientIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
		jen.ID("webhookIDFetcher").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")),
	}
	for _, typ := range proj.DataTypes {
		structFields = append(structFields, jen.IDf("%sIDFetcher", typ.Name.UnexportedVarName()).Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("uint64")))
	}

	code.Add(
		jen.Type().Defs(
			jen.Comment("AuthService is a subset of the larger types.AuthService interface."),
			jen.ID("AuthService").Interface(
				jen.ID("UserAttributionMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("PermissionFilterMiddleware").Params(jen.ID("permissions").Op("...").ID("authorization").Dot("Permission")).Params(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))),
				jen.ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Newline(),
				jen.ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("loginData").Op("*").Qual(proj.TypesPackage(), "UserLoginInput"),
				).Params(jen.Op("*").Qual(proj.TypesPackage(), "User"),
					jen.Op("*").Qual("net/http", "Cookie"),
					jen.ID("error"),
				),
				jen.ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("sessionCtxData").Op("*").Qual(proj.TypesPackage(), "SessionContextData"),
					jen.ID("req").Op("*").Qual("net/http", "Request"),
					jen.ID("res").Qual("net/http", "ResponseWriter"),
				).Params(jen.ID("error")),
			),
			jen.Newline(),
			jen.Comment("UsersService is a subset of the larger types.UsersService interface."),
			jen.ID("UsersService").Interface(
				jen.ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("registrationInput").Op("*").Qual(proj.TypesPackage(), "UserRegistrationInput"),
				).Params(jen.Op("*").Qual(proj.TypesPackage(), "UserCreationResponse"),
					jen.ID("error"),
				),
				jen.ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("input").Op("*").Qual(proj.TypesPackage(), "TOTPSecretVerificationInput"),
				).Params(jen.ID("error")),
			),
			jen.Newline(),
			jen.Comment("Service serves HTML."),
			jen.ID("Service").Interface(jen.ID("SetupRoutes").Params(jen.ID("router").ID("routing").Dot("Router"))),
			jen.Newline(),
			jen.ID("service").Struct(structFields...),
		),
		jen.Newline(),
	)

	serviceInitFields := []jen.Code{
		jen.ID("useFakeData").Op(":").ID("cfg").Dot("UseFakeData"),
		jen.ID("logger").Op(":").Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
		jen.ID("tracer").Op(":").Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("serviceName")),
		jen.ID("panicker").Op(":").Qual(proj.InternalPackage("panicking"), "NewProductionPanicker").Call(),
		jen.ID("localizer").Op(":").ID("provideLocalizer").Call(),
		jen.ID("sessionContextDataFetcher").Op(":").Qual(proj.AuthServicePackage(), "FetchContextFromRequest"),
		jen.ID("authService").Op(":").ID("authService"),
		jen.ID("usersService").Op(":").ID("usersService"),
		jen.ID("paymentManager").Op(":").ID("paymentManager"),
		jen.ID("dataStore").Op(":").ID("dataStore"),
		jen.ID("apiClientIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
			jen.ID("logger"),
			jen.ID("apiClientIDURLParamKey"),
			jen.Lit("API client"),
		),
		jen.ID("accountIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
			jen.ID("logger"),
			jen.ID("accountIDURLParamKey"),
			jen.Lit("account"),
		),
		jen.ID("webhookIDFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
			jen.ID("logger"),
			jen.ID("webhookIDURLParamKey"),
			jen.Lit("webhook"),
		),
	}

	for _, typ := range proj.DataTypes {
		tuvn := typ.Name.UnexportedVarName()
		serviceInitFields = append(serviceInitFields,
			jen.IDf("%sIDFetcher", tuvn).Op(":").ID("routeParamManager").Dot("BuildRouteParamIDFetcher").Call(
				jen.ID("logger"),
				jen.IDf("%sIDURLParamKey", tuvn),
				jen.Lit(typ.Name.SingularCommonName()),
			),
		)
	}

	serviceInitFields = append(serviceInitFields,
		jen.ID("templateFuncMap").Op(":").Map(jen.ID("string")).Interface().Valuesln(jen.Lit("relativeTime").Op(":").ID("relativeTime"),
			jen.Lit("relativeTimeFromPtr").Op(":").ID("relativeTimeFromPtr"),
		),
	)

	code.Add(
		jen.Comment("ProvideService builds a new Service."),
		jen.Newline(),
		jen.Func().ID("ProvideService").Paramsln(
			jen.ID("cfg").Op("*").ID("Config"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("authService").ID("AuthService"),
			jen.ID("usersService").ID("UsersService"),
			jen.ID("dataStore").ID("database").Dot("DataManager"),
			jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager"),
			jen.ID("paymentManager").ID("capitalism").Dot("PaymentManager"),
		).Params(jen.ID("Service")).Body(
			jen.ID("svc").Op(":=").Op("&").ID("service").Valuesln(serviceInitFields...),
			jen.Newline(),
			jen.ID("svc").Dot("templateFuncMap").Index(jen.Lit("translate")).Op("=").ID("svc").Dot("getSimpleLocalizedString"),
			jen.Newline(),
			jen.Return().ID("svc"),
		),
		jen.Newline(),
	)

	return code
}
