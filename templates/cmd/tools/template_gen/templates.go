package template_gen

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
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
		jen.Func().ID("parseTemplate").Params(jen.List(jen.ID("name"), jen.ID("source")).ID("string"), jen.ID("funcMap").Qual("text/template", "FuncMap")).Params(jen.Op("*").Qual("text/template", "Template")).Body(
			jen.Return().Qual("text/template", "Must").Call(jen.Qual("text/template", "New").Call(jen.ID("name")).Dot("Funcs").Call(jen.ID("mergeFuncMaps").Call(jen.ID("defaultTemplateFuncMap"), jen.ID("funcMap"))).Dot("Parse").Call(jen.ID("source"))),
		),
		jen.Line(),
	}
}

func buildMergeFuncMaps() []jen.Code {
	return []jen.Code{
		jen.Func().ID("mergeFuncMaps").Params(jen.List(jen.ID("a"),
			jen.ID("b")).Qual("text/template", "FuncMap")).Params(jen.Qual("text/template", "FuncMap")).Body(
			jen.ID("out").Op(":=").Map(jen.ID("string")).Interface().Values(),
			jen.Line(),
			jen.For(jen.List(jen.ID("k"),
				jen.ID("v")).Op(":=").Range().ID("a")).Body(
				jen.ID("out").Index(jen.ID("k")).Op("=").ID("v"),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("k"),
				jen.ID("v")).Op(":=").Range().ID("b")).Body(
				jen.ID("out").Index(jen.ID("k")).Op("=").ID("v"),
			),
			jen.Line(),
			jen.Return().ID("out"),
		),
		jen.Line(),
	}
}

func buildFormFieldDecl() []jen.Code {
	return []jen.Code{
		jen.Type().ID("formField").Struct(
			jen.ID("LabelName").ID("string"),
			jen.ID("StructFieldName").ID("string"),
			jen.ID("FormName").ID("string"),
			jen.ID("TagID").ID("string"),
			jen.ID("InputType").ID("string"),
			jen.ID("InputPlaceholder").ID("string"),
			jen.ID("Required").ID("bool"),
		),
		jen.Line(),
	}
}
