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
		jen.Func().ID("buildBasicEditorTemplate").Params(jen.ID("cfg").PointerTo().ID("basicEditorTemplateConfig")).Params(jen.ID("string")).Body(
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("parseTemplate").Call(jen.Lit(""), jen.ID("basicEditorTemplateSrc"), jen.ID("nil")).Dot("Execute").Call(jen.AddressOf().ID("b"), jen.ID("cfg")), jen.ID("err").Op("!=").ID("nil")).Body(
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
		jen.Lit("internal/services/frontend/templates/partials/generated/editors/account_editor.gotpl").Op(":").Valuesln(
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
		jen.Lit("internal/services/frontend/templates/partials/generated/editors/api_client_editor.gotpl").Op(":").Valuesln(
			jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("name"),
					jen.ID("FormName").Op(":").Lit("name"),
					jen.ID("StructFieldName").Op(":").Lit("Name"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("client_id"),
					jen.ID("FormName").Op(":").Lit("client_id"),
					jen.ID("StructFieldName").Op(":").Lit("ClientID"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true")),
			),
		),
		jen.Lit("internal/services/frontend/templates/partials/generated/editors/webhook_editor.gotpl").Op(":").Valuesln(
			jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("name"),
					jen.ID("StructFieldName").Op(":").Lit("Name"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("Method"),
					jen.ID("StructFieldName").Op(":").Lit("Method"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true")),
				jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit("ContentType"),
					jen.ID("StructFieldName").Op(":").Lit("ContentType"),
					jen.ID("InputType").Op(":").Lit("text"),
					jen.ID("Required").Op(":").ID("true")),
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
		var formFields []jen.Code
		for _, field := range typ.Fields {
			if field.ValidForCreationInput {
				fn := field.Name
				formFields = append(formFields, jen.Valuesln(
					jen.ID("LabelName").Op(":").Lit(fn.UnexportedVarName()),
					jen.ID("FormName").Op(":").Lit(fn.UnexportedVarName()),
					jen.ID("StructFieldName").Op(":").Lit(fn.Singular()),
					jen.ID("InputType").Op(":").Lit(determineFormType(field.Type)),
					jen.ID("Required").Op(":").ID("true"),
				))
			}
		}

		editorTemplateConfigs = append(
			editorTemplateConfigs,
			jen.Litf("internal/services/frontend/templates/partials/generated/editors/%s_editor.gotpl", typ.Name.RouteName()).Op(":").Valuesln(
				jen.ID("Fields").Op(":").Index().ID("formField").Valuesln(formFields...),
			),
		)
	}

	return []jen.Code{
		jen.Var().ID("editorConfigs").Equals().Map(jen.ID("string")).PointerTo().ID("basicEditorTemplateConfig").Valuesln(
			editorTemplateConfigs...,
		),
		jen.Newline(),
	}
}
