package uploads

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func uploaderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().ID("RootUploadDirectory").ID("string").Type().ID("UploadManager").Interface(
			jen.ID("SaveFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("path").ID("string"), jen.ID("content").Index().ID("byte")).Params(jen.ID("error")),
			jen.ID("ReadFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("path").ID("string")).Params(jen.Index().ID("byte"), jen.ID("error")),
			jen.ID("ServeFiles").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
		),
		jen.Line(),
	)

	return code
}
