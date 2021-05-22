package search

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func searchDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code)

	code.Add(buildTypeDefs()...)

	return code
}

func buildTypeDefs() []jen.Code {
	lines := []jen.Code{
		jen.Type().Defs(
			jen.Comment("IndexPath is a type alias for dependency injection's sake"),
			jen.ID("IndexPath").String(),
			jen.Line(),
			jen.Comment("IndexName is a type alias for dependency injection's sake"),
			jen.ID("IndexName").String(),
			jen.Line(),
			jen.Comment("IndexManager is our wrapper interface for a text search index"),
			jen.ID("IndexManager").Interface(
				jen.ID("Index").Params(
					constants.CtxParam(),
					jen.ID("id").Uint64(),
					jen.ID("value").Interface(),
				).Error(),
				jen.ID("Search").Params(
					constants.CtxParam(),
					jen.ID("query").String(),
					constants.UserIDParam(),
				).Params(jen.ID("ids").Index().Uint64(), jen.Err().Error()),
				jen.ID("Delete").Params(
					constants.CtxParam(),
					jen.ID("id").Uint64(),
				).Params(jen.Err().Error()),
			),
			jen.Line(),
			jen.Comment("IndexManagerProvider is a function that provides an IndexManager for a given index."),
			jen.ID("IndexManagerProvider").Func().Params(
				jen.ID("path").ID("IndexPath"),
				jen.ID("name").ID("IndexName"),
				constants.LoggerParam(),
			).Params(jen.ID("IndexManager"), jen.Error()),
		),
		jen.Line(),
	}

	return lines
}
