package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serviceDataEventsDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(
		jen.Type().ID("ServiceDataEvent").ID("string"),
		jen.Line(),
	)

	ret.Add(
		jen.Const().Defs(
			jen.ID("Create").ID("ServiceDataEvent").Op("=").Lit("create"),
			jen.ID("Update").ID("ServiceDataEvent").Op("=").Lit("update"),
			jen.ID("Archive").ID("ServiceDataEvent").Op("=").Lit("delete"),
		),
		jen.Line(),
	)
	return ret
}
