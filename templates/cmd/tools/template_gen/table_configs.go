package template_gen

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func tableConfigsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	code.Add(buildBasicTableTemplateConfig()...)

	code.Add(
		jen.Comment("//go:embed templates/table.gotpl"),
		jen.Newline(),
		jen.Var().ID("basicTableTemplateSrc").String(),
		jen.Newline(),
	)

	code.Add(buildBuildBasicTableTemplate()...)
	code.Add(buildTableConfigs(proj.DataTypes)...)

	return code
}

func buildBasicTableTemplateConfig() []jen.Code {
	return []jen.Code{
		jen.Type().ID("basicTableTemplateConfig").Struct(jen.ID("SearchURL").String(),
			jen.ID("CreatorPageURL").String(),
			jen.ID("RowDataFieldName").String(),
			jen.ID("Title").String(),
			jen.ID("CreatorPagePushURL").String(),
			jen.ID("CellFields").Index().String(),
			jen.ID("Columns").Index().String(),
			jen.ID("EnableSearch").ID("bool"),
			jen.ID("ExcludeIDRow").ID("bool"),
			jen.ID("ExcludeLink").ID("bool"),
			jen.ID("IncludeLastUpdatedOn").ID("bool"),
			jen.ID("IncludeCreatedOn").ID("bool"),
			jen.ID("IncludeDeleteRow").ID("bool"),
		),
		jen.Newline(),
	}
}

func buildBuildBasicTableTemplate() []jen.Code {
	return []jen.Code{
		jen.Func().ID("buildBasicTableTemplate").Params(jen.ID("cfg").PointerTo().ID("basicTableTemplateConfig")).Params(jen.String()).Body(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("parseTemplate").Call(jen.Lit(""), jen.ID("basicTableTemplateSrc"), jen.Nil()).Dot("Execute").Call(jen.AddressOf().ID("b"), jen.ID("cfg")), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("String").Call(),
		),
		jen.Newline(),
	}
}

func buildTableConfigs(types []models.DataType) []jen.Code {
	tableConfigs := []jen.Code{
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/api_clients_table.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("API Clients"),
			jen.ID("CreatorPagePushURL").MapAssign().Lit("/api_clients/new"),
			jen.ID("CreatorPageURL").MapAssign().Lit("/dashboard_pages/api_clients/new"),
			jen.ID("Columns").MapAssign().Index().String().Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("Client ID"),
				jen.Lit("Belongs To User"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").MapAssign().Index().String().Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("ClientID"),
				jen.Lit("BelongsToUser"),
				jen.Lit("CreatedOn"),
			),
			jen.ID("RowDataFieldName").MapAssign().Lit("Clients"),
			jen.ID("IncludeLastUpdatedOn").MapAssign().False(),
			jen.ID("IncludeCreatedOn").MapAssign().True(),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/accounts_table.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("Accounts"),
			jen.ID("CreatorPagePushURL").MapAssign().Lit("/accounts/new"),
			jen.ID("CreatorPageURL").MapAssign().Lit("/dashboard_pages/accounts/new"),
			jen.ID("Columns").MapAssign().Index().String().Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("Belongs To User"),
				jen.Lit("Last Updated On"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").MapAssign().Index().String().Valuesln(
				jen.Lit("Name"),
				jen.Lit("BelongsToUser"),
			),
			jen.ID("RowDataFieldName").MapAssign().Lit("Accounts"),
			jen.ID("IncludeLastUpdatedOn").MapAssign().True(),
			jen.ID("IncludeCreatedOn").MapAssign().True(),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/users_table.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("Users"),
			jen.ID("Columns").MapAssign().Index().String().Valuesln(
				jen.Lit("ID"),
				jen.Lit("Username"),
				jen.Lit("Last Updated On"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").MapAssign().Index().String().Valuesln(
				jen.Lit("Username"),
			),
			jen.ID("EnableSearch").MapAssign().True(),
			jen.ID("RowDataFieldName").MapAssign().Lit("Users"),
			jen.ID("IncludeLastUpdatedOn").MapAssign().True(),
			jen.ID("IncludeCreatedOn").MapAssign().True(),
			jen.ID("IncludeDeleteRow").MapAssign().False(),
			jen.ID("ExcludeLink").MapAssign().True(),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/tables/webhooks_table.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("Webhooks"),
			jen.ID("CreatorPagePushURL").MapAssign().Lit("/accounts/webhooks/new"),
			jen.ID("CreatorPageURL").MapAssign().Lit("/dashboard_pages/accounts/webhooks/new"),
			jen.ID("Columns").MapAssign().Index().String().Valuesln(
				jen.Lit("ID"),
				jen.Lit("Name"),
				jen.Lit("Method"),
				jen.Lit("URL"),
				jen.Lit("Content Type"),
				jen.Lit("Belongs To Account"),
				jen.Lit("Last Updated On"),
				jen.Lit("Created On"),
			),
			jen.ID("CellFields").MapAssign().Index().String().Valuesln(
				jen.Lit("Name"),
				jen.Lit("Method"),
				jen.Lit("URL"),
				jen.Lit("ContentType"),
				jen.Lit("BelongsToAccount"),
			),
			jen.ID("RowDataFieldName").MapAssign().Lit("Webhooks"),
			jen.ID("IncludeLastUpdatedOn").MapAssign().True(),
			jen.ID("IncludeCreatedOn").MapAssign().True(),
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
			jen.Litf("internal/services/frontend/templates/partials/generated/tables/%s_table.gotpl", tn.PluralRouteName()).MapAssign().Valuesln(
				jen.ID("Title").MapAssign().Lit(tn.Plural()),
				jen.ID("CreatorPagePushURL").MapAssign().Litf("/%s/new", tn.PluralRouteName()),
				jen.ID("CreatorPageURL").MapAssign().Litf("/dashboard_pages/%s/new", tn.PluralRouteName()),
				jen.ID("Columns").MapAssign().Index().String().Valuesln(columns...),
				jen.ID("CellFields").MapAssign().Index().String().Valuesln(cellFields...),
				jen.ID("RowDataFieldName").MapAssign().Lit(tn.Plural()),
				jen.ID("IncludeLastUpdatedOn").MapAssign().True(),
				jen.ID("IncludeCreatedOn").MapAssign().True(),
				jen.ID("IncludeDeleteRow").MapAssign().True(),
			),
		)
	}

	return []jen.Code{
		jen.Var().ID("tableConfigs").Equals().Map(jen.String()).PointerTo().ID("basicTableTemplateConfig").Valuesln(tableConfigs...),
		jen.Newline(),
	}
}
