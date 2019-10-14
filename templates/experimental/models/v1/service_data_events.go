package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serviceDataEventsDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(jen.Null().Type().ID("ServiceDataEvent").ID("string"),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Create").ID("ServiceDataEvent").Op("=").Lit("create").Var().ID("Update").ID("ServiceDataEvent").Op("=").Lit("update").Var().ID("Archive").ID("ServiceDataEvent").Op("=").Lit("delete"),

		jen.Line(),
	)
	return ret
}
