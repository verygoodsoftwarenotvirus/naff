package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
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
		jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
		jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
		jen.ID("panicker").Qual(proj.InternalPackage("panicking"), "Panicker"),
		jen.ID("authService").ID("AuthService"),
		jen.ID("config").PointerTo().ID("Config"),
		jen.ID("sessionContextDataFetcher").Func().Params(jen.PointerTo().Qual("net/http", "Request")).Params(jen.PointerTo().Qual(proj.TypesPackage(), "SessionContextData"), jen.ID("error")),
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
			jen.ID("Service").Interface(jen.ID("SetupRoutes").Params(jen.ID("router").Qual(proj.RoutingPackage(), "Router"))),
			jen.Newline(),
			jen.ID("service").Struct(structFields...),
		),
		jen.Newline(),
	)

	serviceInitFields := []jen.Code{
		jen.ID("config").MapAssign().ID("cfg"),
		jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
		jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("serviceName")),
		jen.ID("panicker").MapAssign().Qual(proj.InternalPackage("panicking"), "NewProductionPanicker").Call(),
		jen.ID("sessionContextDataFetcher").MapAssign().Qual(proj.AuthServicePackage(), "FetchContextFromRequest"),
		jen.ID("authService").MapAssign().ID("authService"),
	}

	code.Add(
		jen.Comment("ProvideService builds a new Service."),
		jen.Newline(),
		jen.Func().ID("ProvideService").Paramsln(
			jen.ID("cfg").PointerTo().ID("Config"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("authService").ID("AuthService"),
		).Params(jen.ID("Service")).Body(
			jen.ID("svc").Assign().AddressOf().ID("service").Valuesln(serviceInitFields...),
			jen.Newline(),
			jen.If(jen.ID("cfg").Dot("Debug")).Body(
				jen.ID("svc").Dot(constants.LoggerVarName).Dot("SetLevel").Call(jen.Qual(proj.InternalLoggingPackage(), "DebugLevel")),
			),
			jen.Newline(),
			jen.Return().ID("svc"),
		),
		jen.Newline(),
	)

	return code
}
