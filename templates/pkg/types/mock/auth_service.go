package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authServiceDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AuthService").Op("=").Parens(jen.Op("*").ID("AuthService")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AuthService").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("StatusHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("StatusHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("req"),
				jen.ID("res"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PermissionFilterMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("PermissionFilterMiddleware").Params(jen.ID("perms").Op("...").ID("authorization").Dot("Permission")).Params(jen.Func().Params(jen.Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("perms")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Func().Params(jen.Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BeginSessionHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("BeginSessionHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("req"),
				jen.ID("res"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EndSessionHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("EndSessionHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("req"),
				jen.ID("res"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CycleCookieSecretHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("CycleCookieSecretHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("req"),
				jen.ID("res"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PASETOHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("PASETOHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("req"),
				jen.ID("res"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ChangeActiveAccountHandler satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("ChangeActiveAccountHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("req"),
				jen.ID("res"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CookieRequirementMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("CookieRequirementMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserAttributionMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("UserAttributionMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuthorizationMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("AuthorizationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ServiceAdminMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserLoginInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("UserLoginInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("PASETOCreationInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("PASETOCreationInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ChangeActiveAccountInputMiddleware satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("ChangeActiveAccountInputMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
			jen.Return().ID("m").Dot("Called").Call(jen.ID("next")).Dot("Get").Call(jen.Lit(0)).Assert(jen.Qual("net/http", "Handler"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AuthenticateUser satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("loginData").Op("*").ID("types").Dot("UserLoginInput")).Params(jen.Op("*").ID("types").Dot("User"), jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("loginData"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("User")), jen.ID("returnValues").Dot("Get").Call(jen.Lit(1)).Assert(jen.Op("*").Qual("net/http", "Cookie")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(2))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogoutUser satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuthService")).ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("types").Dot("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("res").Qual("net/http", "ResponseWriter")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("sessionCtxData"),
				jen.ID("req"),
				jen.ID("res"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
