package server

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
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
		jen.Newline(),
		jen.Func().ID("Build").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			constants.LoggerVar().Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("cfg").PointerTo().Qual(proj.InternalConfigPackage(), "InstanceConfig"),
		).Params(jen.PointerTo().ID("server").Dot("HTTPServer"), jen.ID("error")).Body(
			jen.Qual(constants.DependencyInjectionPkg, "Build").Callln(
				append([]jen.Code{
					jen.Qual(proj.InternalSearchPackage("elasticsearch"), "Providers"),
					jen.Qual(proj.InternalConfigPackage(), "Providers"),
					jen.Qual(proj.DatabasePackage(), "Providers"),
					jen.Qual(proj.DatabasePackage("config"), "Providers"),
					jen.Qual(proj.EncodingPackage(), "Providers"),
					jen.Qual(proj.InternalMessageQueueConfigPackage(), "Providers"),
					jen.Qual(proj.HTTPServerPackage(), "Providers"),
					jen.Qual(proj.MetricsPackage(), "Providers"),
					jen.Qual(proj.InternalImagesPackage(), "Providers"),
					jen.Qual(proj.UploadsPackage(), "Providers"),
					jen.Qual(proj.ObservabilityPackage(), "Providers"),
					jen.Qual(proj.StoragePackage(), "Providers"),
					jen.Qual(proj.RoutingPackage("chi"), "Providers"),
					jen.Qual(proj.InternalAuthenticationPackage(), "Providers"),
					jen.Qual(proj.AuthServicePackage(), "Providers"),
					jen.Qual(proj.UsersServicePackage(), "Providers"),
					jen.Qual(proj.AccountsServicePackage(), "Providers"),
					jen.Qual(proj.APIClientsServicePackage(), "Providers"),
					jen.Qual(proj.WebhooksServicePackage(), "Providers"),
					jen.Qual(proj.WebsocketsServicePackage(), "Providers"),
					jen.Qual(proj.AdminServicePackage(), "Providers"),
					jen.Qual(proj.FrontendServicePackage(), "Providers"),
				},
					imports...,
				)...,
			),
			jen.Newline(),
			jen.Return().List(jen.Nil(), jen.Nil()),
		),
		jen.Newline(),
	)

	return code
}
