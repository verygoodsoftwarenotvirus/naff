package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestService_SetupRoutes").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestHelper").Call(jen.ID("t")),
					jen.ID("obligatoryHandler").Op(":=").Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.Qual("net/http", "ResponseWriter"), jen.Op("*").Qual("net/http", "Request")).Body()),
					jen.ID("authService").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/gamut/pkg/types/mock", "AuthService").Valuesln(),
					jen.ID("authService").Dot("On").Call(
						jen.Lit("ServiceAdminMiddleware"),
						jen.ID("mock").Dot("IsType").Call(jen.ID("obligatoryHandler")),
					).Dot("Return").Call(jen.Qual("net/http", "Handler").Call(jen.ID("obligatoryHandler"))),
					jen.ID("authService").Dot("On").Call(
						jen.Lit("PermissionFilterMiddleware"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("authorization").Dot("Permission").Valuesln()),
					).Dot("Return").Call(jen.Func().Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")).Body(
						jen.Return().Qual("net/http", "Handler").Call(jen.ID("obligatoryHandler")))),
					jen.ID("authService").Dot("On").Call(
						jen.Lit("UserAttributionMiddleware"),
						jen.ID("mock").Dot("IsType").Call(jen.ID("obligatoryHandler")),
					).Dot("Return").Call(jen.Qual("net/http", "Handler").Call(jen.ID("obligatoryHandler"))),
					jen.ID("s").Dot("service").Dot("authService").Op("=").ID("authService"),
					jen.ID("router").Op(":=").ID("chi").Dot("NewRouter").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
					jen.ID("s").Dot("service").Dot("SetupRoutes").Call(jen.ID("router")),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
