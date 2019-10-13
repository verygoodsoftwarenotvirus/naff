package v1

import jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func serviceDataEventsDotGo() *jen.File {
	ret := jen.NewFile("models")
	ret.Add(jen.Null().Type().ID("ServiceDataEvent").ID("string"),

		jen.Line(),
	)
	ret.Add(jen.Null().Var().ID("Create").ID("ServiceDataEvent").Op("=").Lit("create").Var().ID("Update").ID("ServiceDataEvent").Op("=").Lit("update").Var().ID("Archive").ID("ServiceDataEvent").Op("=").Lit("delete"),

		jen.Line(),
	)
	return ret
}
