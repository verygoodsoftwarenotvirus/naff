package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	code.PackageComment("// +build wireinject\n")

	utils.AddImports(proj, code, false)

	imports := []jen.Code{}
	for _, typ := range proj.DataTypes {
		imports = append(imports, jen.Qual(proj.ServicePackage(typ.Name.PackageName()), "Providers"))
	}

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
				jen.Qual(proj.HTTPServerPackage(), "Providers"),
				jen.Qual(proj.MetricsPackage(), "Providers"),
				jen.Qual(proj.InternalImagesPackage(), "Providers"),
				jen.Qual(proj.UploadsPackage(), "Providers"),
				jen.Qual(proj.ObservabilityPackage(), "Providers"),
				jen.Qual(proj.StoragePackage(), "Providers"),
				jen.Qual(proj.CapitalismPackage(), "Providers"),
				jen.Qual(proj.CapitalismPackage("stripe"), "Providers"),
				jen.Qual(proj.RoutingPackage("chi"), "Providers"),
				jen.Qual(proj.InternalAuthenticationPackage(), "Providers"),
				jen.Qual(proj.AuthServicePackage(), "Providers"),
				jen.Qual(proj.UsersServicePackage(), "Providers"),
				jen.Qual(proj.AccountsServicePackage(), "Providers"),
				jen.Qual(proj.APIClientsServicePackage(), "Providers"),
				jen.Qual(proj.WebhooksServicePackage(), "Providers"),
				jen.Qual(proj.AuditServicePackage(), "Providers"),
				jen.Qual(proj.AdminServicePackage(), "Providers"),
				jen.Qual(proj.FrontendServicePackage(), "Providers"),
				jen.Null().Add(imports...),
			),
			jen.Line(),
			jen.Return().List(jen.ID("nil"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
