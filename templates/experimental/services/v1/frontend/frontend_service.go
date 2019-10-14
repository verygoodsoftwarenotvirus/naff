package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func frontendServiceDotGo() *jen.File {
	ret := jen.NewFile("frontend")

	utils.AddImports(ret)

	ret.Add(jen.Null(),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("serviceName").Op("=").Lit("frontend_service"),

		jen.Line(),
	)
	ret.Add(jen.Null().Type().ID("Service").Struct(jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("config").ID("config").Dot(
		"FrontendSettings",
	)),

		jen.Line(),
	)
	ret.Add(jen.Func().Comment("// ProvideFrontendService provides the frontend service to dependency injection").ID("ProvideFrontendService").Params(jen.ID("logger").ID("logging").Dot(
		"Logger",
	), jen.ID("cfg").ID("config").Dot(
		"FrontendSettings",
	)).Params(jen.Op("*").ID("Service")).Block(
		jen.ID("svc").Op(":=").Op("&").ID("Service").Valuesln(jen.ID("config").Op(":").ID("cfg"), jen.ID("logger").Op(":").ID("logger").Dot(
			"WithName",
		).Call(jen.ID("serviceName"))),
		jen.Return().ID("svc"),
	),

		jen.Line(),
	)
	return ret
}
