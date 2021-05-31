package template_gen

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func tableConfigsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	code.Add(buildBasicTableTemplateConfig()...)

	code.Add(
		jen.Comment("//go:embed templates/table.gotpl"),
		jen.Line(),
		jen.Var().ID("basicTableTemplateSrc").ID("string"),
		jen.Line(),
	)

	code.Add(buildBuildBasicTableTemplate()...)
	code.Add(buildTableConfigs(proj.DataTypes)...)

	return code
}

func buildBasicTableTemplateConfig() []jen.Code {
	return []jen.Code{
		jen.Type().ID("basicTableTemplateConfig").Struct(jen.ID("SearchURL").ID("string"),
			jen.ID("CreatorPageURL").ID("string"),
			jen.ID("RowDataFieldName").ID("string"),
			jen.ID("Title").ID("string"),
			jen.ID("CreatorPagePushURL").ID("string"),
			jen.ID("CellFields").Index().ID("string"),
			jen.ID("Columns").Index().ID("string"),
			jen.ID("EnableSearch").ID("bool"),
			jen.ID("ExcludeIDRow").ID("bool"),
			jen.ID("ExcludeLink").ID("bool"),
			jen.ID("IncludeLastUpdatedOn").ID("bool"),
			jen.ID("IncludeCreatedOn").ID("bool"),
			jen.ID("IncludeDeleteRow").ID("bool"),
		),
		jen.Line(),
	}
}

func buildBuildBasicTableTemplate() []jen.Code {
	return []jen.Code{
		jen.Func().ID("buildBasicTableTemplate").Params(jen.ID("cfg").Op("*").ID("basicTableTemplateConfig")).Params(jen.ID("string")).Body(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("parseTemplate").Call(jen.Lit(""), jen.ID("basicTableTemplateSrc"), jen.ID("nil")).Dot("Execute").Call(jen.Op("&").ID("b"), jen.ID("cfg")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("b").Dot("String").Call(),
		),
		jen.Line(),
	}
}

func buildTableConfigs(types []models.DataType) []jen.Code {
	tableConfigs := []jen.Code{
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/api_clients_table.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("API Clients"),
			jen.ID("CreatorPagePushURL").Op(":").Lit("/api_clients/new"),
			jen.ID("CreatorPageURL").Op(":").Lit("/dashboard_pages/api_clients/new"),
			jen.ID("Columns").Op(":").Index().ID("string").Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("External ID"),
				jen.Lit("Client ID"),
				jen.Lit("Belongs To User"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").Op(":").Index().ID("string").Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("ExternalID"),
				jen.Lit("ClientID"),
				jen.Lit("BelongsToAccount"),
				jen.Lit("CreatedOn"),
			),
			jen.ID("RowDataFieldName").Op(":").Lit("Clients"),
			jen.ID("IncludeLastUpdatedOn").Op(":").ID("false"),
			jen.ID("IncludeCreatedOn").Op(":").ID("true"),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/accounts_table.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("Accounts"),
			jen.ID("CreatorPagePushURL").Op(":").Lit("/accounts/new"),
			jen.ID("CreatorPageURL").Op(":").Lit("/dashboard_pages/accounts/new"),
			jen.ID("Columns").Op(":").Index().ID("string").Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("External ID"),
				jen.Lit("Belongs To User"),
				jen.Lit("Last Updated On"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").Op(":").Index().ID("string").Valuesln(
				jen.Lit("Name"),
				jen.Lit("ExternalID"),
				jen.Lit("BelongsToAccount"),
			),
			jen.ID("RowDataFieldName").Op(":").Lit("Accounts"),
			jen.ID("IncludeLastUpdatedOn").Op(":").ID("true"),
			jen.ID("IncludeCreatedOn").Op(":").ID("true"),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/users_table.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("Users"),
			jen.ID("Columns").Op(":").Index().ID("string").Valuesln(
				jen.Lit("ID"),
				jen.Lit("Username"),
				jen.Lit("Last Updated On"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").Op(":").Index().ID("string").Valuesln(
				jen.Lit("Username"),
			),
			jen.ID("EnableSearch").Op(":").ID("true"),
			jen.ID("RowDataFieldName").Op(":").Lit("Users"),
			jen.ID("IncludeLastUpdatedOn").Op(":").ID("true"),
			jen.ID("IncludeCreatedOn").Op(":").ID("true"),
			jen.ID("IncludeDeleteRow").Op(":").ID("false"),
			jen.ID("ExcludeLink").Op(":").ID("true"),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/webhooks_table.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("Webhooks"),
			jen.ID("CreatorPagePushURL").Op(":").Lit("/accounts/webhooks/new"),
			jen.ID("CreatorPageURL").Op(":").Lit("/dashboard_pages/accounts/webhooks/new"),
			jen.ID("Columns").Op(":").Index().ID("string").Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("Method"),
				jen.Lit("URL"),
				jen.Lit("Content Type"),
				jen.Lit("Belongs To Account"),
				jen.Lit("Last Updated On"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").Op(":").Index().ID("string").Valuesln(
				jen.Lit("Name"),
				jen.Lit("Method"),
				jen.Lit("URL"),
				jen.Lit("ContentType"),
				jen.Lit("BelongsToAccount"),
			),
			jen.ID("RowDataFieldName").Op(":").Lit("Webhooks"),
			jen.ID("IncludeLastUpdatedOn").Op(":").ID("true"),
			jen.ID("IncludeCreatedOn").Op(":").ID("true"),
		),
	}

	for _, typ := range types {
		var cellFields []jen.Code
		columns := []jen.Code{
			jen.Lit("ID"),
		}
		for _, field := range typ.Fields {
			cellFields = append(cellFields, jen.Lit(field.Name.Singular()))
			columns = append(columns, jen.Lit(field.Name.Singular()))
		}
		columns = append(columns, jen.Lit("Last Updated On"), jen.Lit("Created On"))

		tn := typ.Name
		tableConfigs = append(tableConfigs,
			jen.Litf("internal/services/frontend/templates/partials/generated/tables/%s_table.gotpl", tn.PluralRouteName()).Op(":").Valuesln(
				jen.ID("Title").Op(":").Lit(tn.Plural()),
				jen.ID("CreatorPagePushURL").Op(":").Litf("/%s/new", tn.PluralRouteName()),
				jen.ID("CreatorPageURL").Op(":").Litf("/dashboard_pages/%s/new", tn.PluralRouteName()),
				jen.ID("Columns").Op(":").Index().ID("string").Valuesln(columns...),
				jen.ID("CellFields").Op(":").Index().ID("string").Valuesln(cellFields...),
				jen.ID("RowDataFieldName").Op(":").Lit(tn.Plural()),
				jen.ID("IncludeLastUpdatedOn").Op(":").ID("true"),
				jen.ID("IncludeCreatedOn").Op(":").ID("true"),
				jen.ID("IncludeDeleteRow").Op(":").ID("true"),
			),
		)
	}

	return []jen.Code{
		jen.Var().ID("tableConfigs").Op("=").Map(jen.ID("string")).Op("*").ID("basicTableTemplateConfig").Valuesln(tableConfigs...),
		jen.Line(),
	}
}
