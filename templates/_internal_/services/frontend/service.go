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
			jen.ID("serviceName").String().Equals().Lit("frontend_service"),
		),
		jen.Newline(),
	)

	structFields := []jen.Code{
		jen.ID("useFakeData").ID("bool"),
		jen.ID("templateFuncMap").Qual("html/template", "FuncMap"),
		jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
		jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
		jen.ID("panicker").Qual(proj.InternalPackage("panicking"), "Panicker"),
		jen.ID("localizer").PointerTo().Qual("github.com/nicksnyder/go-i18n/v2/i18n", "Localizer"),
		jen.ID("dataStore").Qual(proj.DatabasePackage(), "DataManager"),
		jen.ID("authService").ID("AuthService"),
		jen.ID("usersService").ID("UsersService"),
		jen.ID("sessionContextDataFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")),
		jen.ID("accountIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()),
		jen.ID("apiClientIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()),
		jen.ID("webhookIDFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()),
	}
	for _, typ := range proj.DataTypes {
		structFields = append(structFields, jen.IDf("%sIDFetcher", typ.Name.UnexportedVarName()).Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.String()))
	}

	code.Add(
		jen.Type().Defs(
			jen.Comment("AuthService is a subset of the larger types.AuthService interface."),
			jen.ID("AuthService").Interface(
				jen.ID("UserAttributionMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("PermissionFilterMiddleware").Params(jen.ID("permissions").Op("...").Qual(proj.InternalAuthorizationPackage(), "Permission")).Params(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))),
				jen.ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.Newline(),
				jen.ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("loginData").PointerTo().Qual(proj.TypesPackage(), "UserLoginInput"),
				).Params(jen.PointerTo().Qual(proj.TypesPackage(), "User"),
					jen.PointerTo().Qual("net/http", "Cookie"),
					jen.ID("error"),
				),
				jen.ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("sessionCtxData").PointerTo().Qual(proj.TypesPackage(), "SessionContextData"),
					jen.ID("req").PointerTo().Qual("net/http", "Request"),
					jen.ID("res").Qual("net/http", "ResponseWriter"),
				).Params(jen.ID("error")),
			),
			jen.Newline(),
			jen.Comment("UsersService is a subset of the larger types.UsersService interface."),
			jen.ID("UsersService").Interface(
				jen.ID("RegisterUser").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("registrationInput").PointerTo().Qual(proj.TypesPackage(), "UserRegistrationInput"),
				).Params(jen.PointerTo().Qual(proj.TypesPackage(), "UserCreationResponse"),
					jen.ID("error"),
				),
				jen.ID("VerifyUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"),
					jen.ID("input").PointerTo().Qual(proj.TypesPackage(), "TOTPSecretVerificationInput"),
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
		jen.ID("useFakeData").MapAssign().ID("cfg").Dot("UseFakeData"),
		jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
		jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("serviceName")),
		jen.ID("panicker").MapAssign().Qual(proj.InternalPackage("panicking"), "NewProductionPanicker").Call(),
		jen.ID("localizer").MapAssign().ID("provideLocalizer").Call(),
		jen.ID("sessionContextDataFetcher").MapAssign().Qual(proj.AuthServicePackage(), "FetchContextFromRequest"),
		jen.ID("authService").MapAssign().ID("authService"),
		jen.ID("usersService").MapAssign().ID("usersService"),
		jen.ID("dataStore").MapAssign().ID("dataStore"),
		jen.ID("apiClientIDFetcher").MapAssign().ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(
			jen.ID("apiClientIDURLParamKey"),
		),
		jen.ID("accountIDFetcher").MapAssign().ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(
			jen.ID("accountIDURLParamKey"),
		),
		jen.ID("webhookIDFetcher").MapAssign().ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(
			jen.ID("webhookIDURLParamKey"),
		),
	}

	for _, typ := range proj.DataTypes {
		tuvn := typ.Name.UnexportedVarName()
		serviceInitFields = append(serviceInitFields,
			jen.IDf("%sIDFetcher", tuvn).MapAssign().ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(
				jen.IDf("%sIDURLParamKey", tuvn),
			),
		)
	}

	serviceInitFields = append(serviceInitFields,
		jen.ID("templateFuncMap").MapAssign().Map(jen.String()).Interface().Valuesln(jen.Lit("relativeTime").MapAssign().ID("relativeTime"),
			jen.Lit("relativeTimeFromPtr").MapAssign().ID("relativeTimeFromPtr"),
		),
	)

	code.Add(
		jen.Comment("ProvideService builds a new Service."),
		jen.Newline(),
		jen.Func().ID("ProvideService").Paramsln(
			jen.ID("cfg").PointerTo().ID("Config"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("authService").ID("AuthService"),
			jen.ID("usersService").ID("UsersService"),
			jen.ID("dataStore").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("routeParamManager").Qual(proj.RoutingPackage(), "RouteParamManager"),
		).Params(jen.ID("Service")).Body(
			jen.ID("svc").Assign().AddressOf().ID("service").Valuesln(serviceInitFields...),
			jen.Newline(),
			jen.ID("svc").Dot("templateFuncMap").Index(jen.Lit("translate")).Equals().ID("svc").Dot("getSimpleLocalizedString"),
			jen.Newline(),
			jen.Return().ID("svc"),
		),
		jen.Newline(),
	)

	return code
}
