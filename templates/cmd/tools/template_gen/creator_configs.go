package template_gen

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func creatorConfigsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	code.Add(buildBuildBasicCreatorTemplate()...)
	code.Add(basicCreatorTemplateConfig()...)
	code.Add(buildCreatorConfigs(proj.DataTypes)...)

	return code
}

func buildBuildBasicCreatorTemplate() []jen.Code {
	return []jen.Code{
		jen.Comment("//go:embed templates/creator.gotpl"),
		jen.Newline(),
		jen.Var().ID("basicCreatorTemplateSrc").String(),
		jen.Newline(),
		jen.Func().ID("buildBasicCreatorTemplate").Params(jen.ID("cfg").PointerTo().ID("basicCreatorTemplateConfig")).Params(jen.String()).Body(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("parseTemplate").Call(jen.Lit(""), jen.ID("basicCreatorTemplateSrc"), jen.ID("nil")).Dot("Execute").Call(jen.AddressOf().ID("b"), jen.ID("cfg")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("String").Call(),
		),
	}
}

func basicCreatorTemplateConfig() []jen.Code {
	return []jen.Code{
		jen.Type().ID("basicCreatorTemplateConfig").Struct(
			jen.ID("Title").String(),
			jen.ID("SubmissionURL").String(),
			jen.ID("Fields").Slice().ID("formField"),
		),
		jen.Newline(),
	}
}

func determineFormType(t string) string {
	switch t {
	case "string":
		return "text"
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		return "number"
	default:
		return "text"
	}
}

func buildCreatorConfigs(types []models.DataType) []jen.Code {
	iterableCreatorConfigs := []jen.Code{
		jen.Lit("internal/services/frontend/templates/partials/generated/creators/account_creator.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("New Account"),
			jen.ID("SubmissionURL").Op(":").Lit("/accounts/new/submit"),
			jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("name"),
					jen.ID("FormName").Op(":").Lit("name"),
					jen.ID("StructFieldName").Op(":").Lit("Name"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/creators/api_client_creator.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("New API Client"),
			jen.ID("SubmissionURL").Op(":").Lit("/api_clients/new/submit"),
			jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("name"),
					jen.ID("FormName").Op(":").Lit("name"),
					jen.ID("StructFieldName").Op(":").Lit("Name"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("client_id"),
					jen.ID("FormName").Op(":").Lit("client_id"),
					jen.ID("StructFieldName").Op(":").Lit("ClientID"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/creators/webhook_creator.gotpl").Op(":").Valuesln(
			jen.ID("Title").Op(":").Lit("New Webhook"),
			jen.ID("SubmissionURL").Op(":").Lit("/webhooks/new/submit"),
			jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("name"),
					jen.ID("StructFieldName").Op(":").Lit("Name"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("Method"),
					jen.ID("StructFieldName").Op(":").Lit("Method"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("ContentType"),
					jen.ID("StructFieldName").Op(":").Lit("ContentType"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("URL"),
					jen.ID("StructFieldName").Op(":").Lit("URL"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true"),
				),
			),
		),
	}
	for _, typ := range types {
		var formFieldValues []jen.Code
		for _, field := range typ.Fields {
			if field.ValidForCreationInput {
				fn := field.Name
				formFieldValues = append(formFieldValues, jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit(fn.UnexportedVarName()),
					jen.ID("FormName").Op(":").Lit(fn.UnexportedVarName()),
					jen.ID("StructFieldName").Op(":").Lit(fn.Singular()),
					jen.ID("InputType").Op(":").Lit(determineFormType(field.Type)),
					jen.ID("Required").Op(":").ID("true"),
				))
			}
		}

		n := typ.Name

		iterableCreatorConfigs = append(iterableCreatorConfigs, jen.Litf("internal/services/frontend/templates/partials/generated/creators/%s_creator.gotpl", n.RouteName()).Op(":").Valuesln(
			jen.ID("Title").Op(":").Litf("New %s", n.Singular()),
			jen.ID("SubmissionURL").Op(":").Litf("/%s/new/submit", n.PluralRouteName()),
			jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(
				formFieldValues...,
			),
		),
		)
	}

	return []jen.Code{
		jen.Var().ID("creatorConfigs").Equals().Map(jen.String()).PointerTo().ID("basicCreatorTemplateConfig").Valuesln(
			iterableCreatorConfigs...,
		),
		jen.Newline(),
	}
}
