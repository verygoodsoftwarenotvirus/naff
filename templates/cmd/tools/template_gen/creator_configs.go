package template_gen

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
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
			jen.If(jen.ID("err").Assign().ID("parseTemplate").Call(jen.Lit(""), jen.ID("basicCreatorTemplateSrc"), jen.Nil()).Dot("Execute").Call(jen.AddressOf().ID("b"), jen.ID("cfg")), jen.ID("err").DoesNotEqual().Nil()).Body(
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
		jen.Lit("internal/services/frontend/templates/partials/generated/creators/account_creator.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("New Account"),
			jen.ID("SubmissionURL").MapAssign().Lit("/accounts/new/submit"),
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("name"),
					jen.ID("FormName").MapAssign().Lit("name"),
					jen.ID("StructFieldName").MapAssign().Lit("Name"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
				),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/creators/api_client_creator.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("New API Client"),
			jen.ID("SubmissionURL").MapAssign().Lit("/api_clients/new/submit"),
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("name"),
					jen.ID("FormName").MapAssign().Lit("name"),
					jen.ID("StructFieldName").MapAssign().Lit("Name"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
				),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("client_id"),
					jen.ID("FormName").MapAssign().Lit("client_id"),
					jen.ID("StructFieldName").MapAssign().Lit("ClientID"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
				),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/creators/webhook_creator.gotpl").MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Lit("New Webhook"),
			jen.ID("SubmissionURL").MapAssign().Lit("/webhooks/new/submit"),
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("name"),
					jen.ID("StructFieldName").MapAssign().Lit("Name"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
				),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("Method"),
					jen.ID("StructFieldName").MapAssign().Lit("Method"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
				),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("ContentType"),
					jen.ID("StructFieldName").MapAssign().Lit("ContentType"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
				),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("URL"),
					jen.ID("StructFieldName").MapAssign().Lit("URL"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().True(),
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
					jen.ID("LabelName").MapAssign().Lit(fn.UnexportedVarName()),
					jen.ID("FormName").MapAssign().Lit(fn.UnexportedVarName()),
					jen.ID("StructFieldName").MapAssign().Lit(fn.Singular()),
					jen.ID("InputType").MapAssign().Lit(determineFormType(field.Type)),
					jen.ID("Required").MapAssign().True(),
				))
			}
		}

		n := typ.Name

		iterableCreatorConfigs = append(iterableCreatorConfigs, jen.Litf("internal/services/frontend/templates/partials/generated/creators/%s_creator.gotpl", n.RouteName()).MapAssign().Valuesln(
			jen.ID("Title").MapAssign().Litf("New %s", n.Singular()),
			jen.ID("SubmissionURL").MapAssign().Litf("/%s/new/submit", n.PluralRouteName()),
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
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
