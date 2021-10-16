package template_gen

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func templatesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, true)

	code.Add(buildParseTemplate()...)
	code.Add(buildMergeFuncMaps()...)
	code.Add(buildFormFieldDecl()...)

	return code
}

func buildParseTemplate() []jen.Code {
	return []jen.Code{
		jen.Func().ID("parseTemplate").Params(jen.List(jen.ID("name"), jen.ID("source")).String(), jen.ID("funcMap").Qual("text/template", "FuncMap")).Params(jen.PointerTo().Qual("text/template", "Template")).Body(
			jen.Return().Qual("text/template", "Must").Call(jen.Qual("text/template", "New").Call(jen.ID("name")).Dot("Funcs").Call(jen.ID("mergeFuncMaps").Call(jen.ID("defaultTemplateFuncMap"), jen.ID("funcMap"))).Dot("Parse").Call(jen.ID("source"))),
		),
		jen.Newline(),
	}
}

func buildMergeFuncMaps() []jen.Code {
	return []jen.Code{
		jen.Func().ID("mergeFuncMaps").Params(jen.List(jen.ID("a"),
			jen.ID("b")).Qual("text/template", "FuncMap")).Params(jen.Qual("text/template", "FuncMap")).Body(
			jen.ID("out").Assign().Map(jen.String()).Interface().Values(),
			jen.Newline(),
			jen.For(jen.List(jen.ID("k"),
				jen.ID("v")).Assign().Range().ID("a")).Body(
				jen.ID("out").Index(jen.ID("k")).Equals().ID("v"),
			),
			jen.Newline(),
			jen.For(jen.List(jen.ID("k"),
				jen.ID("v")).Assign().Range().ID("b")).Body(
				jen.ID("out").Index(jen.ID("k")).Equals().ID("v"),
			),
			jen.Newline(),
			jen.Return().ID("out"),
		),
		jen.Newline(),
	}
}

func buildFormFieldDecl() []jen.Code {
	return []jen.Code{
		jen.Type().ID("formField").Struct(
			jen.ID("LabelName").String(),
			jen.ID("StructFieldName").String(),
			jen.ID("FormName").String(),
			jen.ID("TagID").String(),
			jen.ID("InputType").String(),
			jen.ID("InputPlaceholder").String(),
			jen.ID("Required").ID("bool"),
		),
		jen.Newline(),
	}
}
