package main

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func frontendServiceDotGo() *jen.File {
	ret := jen.NewFile("frontend")
	ret.Add(jen.Null(),
	)
	ret.Add(jen.Null().Var().ID("serviceName").Op("=").Lit("frontend_service"),
	)
	ret.Add(jen.Null().Type().ID("Service").Struct(
		jen.ID("logger").ID("logging").Dot(
			"Logger",
		),
		jen.ID("config").ID("config").Dot(
			"FrontendSettings",
		),
	),
	)
	ret.Add(jen.Func(),
	)
	return ret
}
