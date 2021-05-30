package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	code.HeaderComment("+build wireinject")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("Build builds a server."),
		jen.Line(),
		jen.Func().ID("Build").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("cfg").Op("*").Qual(proj.ConfigPackage(), "InstanceConfig"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
		).Params(jen.Op("*").ID("server").Dot("HTTPServer"), jen.ID("error")).Body(
			jen.Qual("github.com/google/wire", "Build").Callln(
				jen.Qual(proj.InternalSearchPackage("bleve"), "Providers"),
				jen.Qual(proj.ConfigPackage(), "Providers"),
				jen.Qual(proj.DatabasePackage(), "Providers"),
				jen.Qual(proj.DatabasePackage("config"), "Providers"),
				jen.Qual(proj.EncodingPackage(), "Providers"),
				jen.Qual(proj.InternalPackage("server"), "Providers"),
				jen.Qual(proj.MetricsPackage(), "Providers"),
				jen.Qual(proj.InternalPackage("uploads", "images"), "Providers"),
				jen.Qual(proj.InternalPackage("uploads"), "Providers"),
				jen.Qual(proj.ObservabilityPackage(), "Providers"),
				jen.Qual(proj.InternalPackage("storage"), "Providers"),
				jen.Qual(proj.InternalPackage("capitalism"), "Providers"),
				jen.Qual(proj.InternalPackage("capitalism", "stripe"), "Providers"),
				jen.Qual(proj.InternalPackage("routing", "chi"), "Providers"),
				jen.Qual(proj.InternalPackage("authentication"), "Providers"),
				jen.Qual(proj.ServicePackage("authentication"), "Providers"),
				jen.Qual(proj.ServicePackage("users"), "Providers"),
				jen.Qual(proj.ServicePackage("accounts"), "Providers"),
				jen.Qual(proj.ServicePackage("apiclients"), "Providers"),
				jen.Qual(proj.ServicePackage("webhooks"), "Providers"),
				jen.Qual(proj.ServicePackage("audit"), "Providers"),
				jen.Qual(proj.ServicePackage("admin"), "Providers"),
				jen.Qual(proj.ServicePackage("frontend"), "Providers"),
				jen.Qual(proj.ServicePackage("items"), "Providers"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("nil"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
