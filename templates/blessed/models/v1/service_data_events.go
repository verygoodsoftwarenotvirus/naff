package v1

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceDataEventsDotGo(proj *models.Project) *jen.File {
	ret := jen.NewFile("models")

	utils.AddImports(proj, ret)

	ret.Add(
		jen.Comment("ServiceDataEvent is a simple string alias."),
		jen.Line(),
		jen.Type().ID("ServiceDataEvent").String(),
		jen.Line(),
	)

	ret.Add(
		jen.Const().Defs(
			jen.Comment("Create represents a create event."),
			jen.ID("Create").ID("ServiceDataEvent").Equals().Lit("create"),
			jen.Comment("Update represents an update event."),
			jen.ID("Update").ID("ServiceDataEvent").Equals().Lit("update"),
			jen.Comment("Archive represents an archive event."),
			jen.ID("Archive").ID("ServiceDataEvent").Equals().Lit("archive"),
		),
		jen.Line(),
	)

	return ret
}
