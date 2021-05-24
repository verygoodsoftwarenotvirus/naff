package server

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("Build builds a server."),
		jen.Line(),
		jen.Func().ID("Build").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("config").Dot("ServerConfig"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("dbm").ID("database").Dot("DataManager"), jen.ID("db").Op("*").Qual("database/sql", "DB"), jen.ID("authenticator").ID("authentication").Dot("Authenticator")).Params(jen.Op("*").ID("server").Dot("HTTPServer"), jen.ID("error")).Body(
			jen.ID("wire").Dot("Build").Call(
				jen.ID("bleve").Dot("Providers"),
				jen.ID("config").Dot("Providers"),
				jen.ID("database").Dot("Providers"),
				jen.ID("encoding").Dot("Providers"),
				jen.ID("server").Dot("Providers"),
				jen.ID("metrics").Dot("Providers"),
				jen.ID("images").Dot("Providers"),
				jen.ID("uploads").Dot("Providers"),
				jen.ID("observability").Dot("Providers"),
				jen.ID("storage").Dot("Providers"),
				jen.ID("capitalism").Dot("Providers"),
				jen.ID("stripe").Dot("Providers"),
				jen.ID("chi").Dot("Providers"),
				jen.ID("authentication").Dot("Providers"),
				jen.ID("users").Dot("Providers"),
				jen.ID("accounts").Dot("Providers"),
				jen.ID("apiclients").Dot("Providers"),
				jen.ID("webhooks").Dot("Providers"),
				jen.ID("audit").Dot("Providers"),
				jen.ID("admin").Dot("Providers"),
				jen.ID("frontend").Dot("Providers"),
				jen.ID("items").Dot("Providers"),
			),
			jen.Return().List(jen.ID("nil"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
