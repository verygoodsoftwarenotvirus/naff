package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
)

func serviceDataEventsDotGo() *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(ret)

	ret.Add(
		jen.Comment("ServiceDataEvent is a simple string alias"),
		jen.Line(),
		jen.Type().ID("ServiceDataEvent").ID("string"),
		jen.Line(),
	)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("Create represents a create event"),
			jen.ID("Create").ID("ServiceDataEvent").Op("=").Lit("create"),
			jen.Comment("Update represents an update event"),
			jen.ID("Update").ID("ServiceDataEvent").Op("=").Lit("update"),
			jen.Comment("Archive represents an archive event"),
			jen.ID("Archive").ID("ServiceDataEvent").Op("=").Lit("archive"),
		),
		jen.Line(),
	)
	return ret
}