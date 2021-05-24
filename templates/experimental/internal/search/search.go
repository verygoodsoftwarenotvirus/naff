package search

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func searchDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Null().Type().ID("IndexPath").ID("string").Type().ID("IndexName").ID("string").Type().ID("IndexManager").Interface(
			jen.ID("Index").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64"), jen.ID("value").Interface()).Params(jen.ID("error")),
			jen.ID("Search").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("accountID").ID("uint64")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")),
			jen.ID("SearchForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")),
			jen.ID("Delete").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64")).Params(jen.ID("err").ID("error")),
		).Type().ID("IndexManagerProvider").Params(jen.ID("path").ID("IndexPath"), jen.ID("name").ID("IndexName"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("IndexManager"), jen.ID("error")),
		jen.Line(),
	)

	return code
}
