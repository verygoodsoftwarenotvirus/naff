package config

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("Providers").Op("=").ID("wire").Dot("NewSet").Call(jen.ID("wire").Dot("FieldsOf").Call(
			jen.ID("new").Call(jen.Op("*").ID("ServerConfig")),
			jen.Lit("Auth"),
			jen.Lit("Database"),
			jen.Lit("Observability"),
			jen.Lit("Capitalism"),
			jen.Lit("Meta"),
			jen.Lit("Encoding"),
			jen.Lit("Frontend"),
			jen.Lit("Uploads"),
			jen.Lit("Search"),
			jen.Lit("Server"),
			jen.Lit("Webhooks"),
			jen.Lit("AuditLog"),
		)),
		jen.Line(),
	)

	return code
}
