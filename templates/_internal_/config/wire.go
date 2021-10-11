package config

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func wireDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	serviceConfigLines := []jen.Code{
		jen.New(jen.PointerTo().ID("ServicesConfigurations")),
		jen.Lit("Auth"),
		jen.Lit("Frontend"),
		jen.Lit("Webhooks"),
		jen.Lit("Websockets"),
		jen.Lit("Accounts"),
	}
	for _, typ := range proj.DataTypes {
		serviceConfigLines = append(serviceConfigLines, jen.Lit(typ.Name.Plural()))
	}

	code.Add(
		jen.Var().Defs(
			jen.Comment("Providers represents this package's offering to the dependency injector."),
			jen.ID("Providers").Equals().Qual(constants.DependencyInjectionPkg, "NewSet").Callln(
				jen.ID("ProvideDatabaseClient"),
				jen.Qual(constants.DependencyInjectionPkg, "FieldsOf").Callln(
					jen.New(jen.PointerTo().ID("InstanceConfig")),
					jen.Lit("Database"),
					jen.Lit("Observability"),
					jen.Lit("Meta"),
					jen.Lit("Encoding"),
					jen.Lit("Uploads"),
					jen.Lit("Search"),
					jen.Lit("Events"),
					jen.Lit("Server"),
					jen.Lit("Services"),
				),
				jen.Qual(constants.DependencyInjectionPkg, "FieldsOf").Callln(
					serviceConfigLines...,
				),
			),
		),
		jen.Newline(),
	)

	return code
}
