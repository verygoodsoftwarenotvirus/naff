package authentication

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func sessionManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("sessionManager").Interface(
				jen.ID("Load").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("token").ID("string")).Params(jen.Qual("context", "Context"), jen.ID("error")),
				jen.ID("RenewToken").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")),
				jen.ID("Get").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("key").ID("string")).Params(jen.Interface()),
				jen.ID("Put").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("key").ID("string"), jen.ID("val").Interface()),
				jen.ID("Commit").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string"), jen.Qual("time", "Time"), jen.ID("error")),
				jen.ID("Destroy").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")),
			),
		),
		jen.Line(),
	)

	return code
}
