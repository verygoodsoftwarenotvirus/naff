package template_gen

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func editorConfigsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	code.Add(
		jen.Comment("//go:embed templates/editor.gotpl"),
		jen.Newline(),
		jen.Var().ID("basicEditorTemplateSrc").String(),
		jen.Newline(),
	)

	code.Add(buildBuildBasicEditorTemplate()...)
	code.Add(buildBasicEditorTemplateConfig()...)
	code.Add(buildEditorConfigs(proj.DataTypes)...)

	return code
}

func buildBuildBasicEditorTemplate() []jen.Code {
	return []jen.Code{
		jen.Func().ID("buildBasicEditorTemplate").Params(jen.ID("cfg").PointerTo().ID("basicEditorTemplateConfig")).Params(jen.String()).Body(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("parseTemplate").Call(jen.Lit(""), jen.ID("basicEditorTemplateSrc"), jen.Nil()).Dot("Execute").Call(jen.AddressOf().ID("b"), jen.ID("cfg")), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.ID("panic").Call(jen.ID("err")),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("String").Call(),
		),
		jen.Newline(),
	}
}

func buildBasicEditorTemplateConfig() []jen.Code {
	return []jen.Code{
		jen.Type().ID("basicEditorTemplateConfig").Struct(jen.ID("SubmissionURL").String(),
			jen.ID("Fields").Index().ID("formField"),
		),
		jen.Newline(),
	}
}

func buildEditorConfigs(types []models.DataType) []jen.Code {
	editorTemplateConfigs := []jen.Code{
		jen.Lit("internal/services/frontend/templates/partials/generated/editors/account_editor.gotpl").MapAssign().Valuesln(
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("name"),
					jen.ID("FormName").MapAssign().Lit("name"),
					jen.ID("StructFieldName").MapAssign().Lit("Name"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true"),
				),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/editors/api_client_editor.gotpl").MapAssign().Valuesln(
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("name"),
					jen.ID("FormName").MapAssign().Lit("name"),
					jen.ID("StructFieldName").MapAssign().Lit("Name"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("client_id"),
					jen.ID("FormName").MapAssign().Lit("client_id"),
					jen.ID("StructFieldName").MapAssign().Lit("ClientID"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true")),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/editors/webhook_editor.gotpl").MapAssign().Valuesln(
			jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("name"),
					jen.ID("StructFieldName").MapAssign().Lit("Name"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("Method"),
					jen.ID("StructFieldName").MapAssign().Lit("Method"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("ContentType"),
					jen.ID("StructFieldName").MapAssign().Lit("ContentType"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit("URL"),
					jen.ID("StructFieldName").MapAssign().Lit("URL"),
					jen.ID("InputType").MapAssign().Lit("text"),
					jen.ID("Required").MapAssign().ID("true"),
				),
			),
		),
	}
	for _, typ := range types {
		var formFields []jen.Code
		for _, field := range typ.Fields {
			if field.ValidForCreationInput {
				fn := field.Name
				formFields = append(formFields, jen.Valuesln(
					jen.ID("LabelName").MapAssign().Lit(fn.UnexportedVarName()),
					jen.ID("FormName").MapAssign().Lit(fn.UnexportedVarName()),
					jen.ID("StructFieldName").MapAssign().Lit(fn.Singular()),
					jen.ID("InputType").MapAssign().Lit(determineFormType(field.Type)),
					jen.ID("Required").MapAssign().ID("true"),
				))
			}
		}

		editorTemplateConfigs = append(
			editorTemplateConfigs,
			jen.Litf("internal/services/frontend/templates/partials/generated/editors/%s_editor.gotpl", typ.Name.RouteName()).MapAssign().Valuesln(
				jen.ID("Fields").MapAssign().Index().ID("formField").Valuesln(formFields...),
			),
		)
	}

	return []jen.Code{
		jen.Var().ID("editorConfigs").Equals().Map(jen.String()).PointerTo().ID("basicEditorTemplateConfig").Valuesln(
			editorTemplateConfigs...,
		),
		jen.Newline(),
	}
}
