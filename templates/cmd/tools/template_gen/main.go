package template_gen

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func mainDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("defaultTemplateFuncMap").Op("=").Map(jen.ID("string")).Interface().Values(),
		jen.Line(),
	)

	code.Add(buildWriteFile()...)
	code.Add(buildMainFunc()...)

	return code
}

func buildWriteFile() []jen.Code {
	return []jen.Code{
		jen.Func().ID("writeFile").Params(jen.List(jen.ID("path"), jen.ID("out")).ID("string")).Params(jen.ID("error")).Body(
			jen.ID("containingDir").Op(":=").Qual("path/filepath", "Dir").Call(jen.ID("path")),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").Qual("os", "MkdirAll").Call(jen.ID("containingDir"), jen.Octal(777)), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("error writing to filepath %q: %w"), jen.ID("path"), jen.ID("err")),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").Qual("io/ioutil", "WriteFile").Call(jen.ID("path"), jen.Index().ID("byte").Call(jen.ID("out")), jen.Octal(644)), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(jen.Lit("error writing to filepath %q: %w"), jen.ID("path"), jen.ID("err")),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}
}

func buildMainFunc() []jen.Code {
	return []jen.Code{
		jen.Func().ID("main").Params().Body(
			jen.For(jen.List(jen.ID("path"), jen.ID("cfg")).Op(":=").Range().ID("editorConfigs")).Body(
				jen.If(jen.ID("err").Op(":=").ID("writeFile").Call(jen.ID("path"), jen.ID("buildBasicEditorTemplate").Call(jen.ID("cfg"))), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Qual("log", "Fatal").Call(jen.ID("err")),
				),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("path"), jen.ID("cfg")).Op(":=").Range().ID("tableConfigs")).Body(
				jen.If(jen.ID("err").Op(":=").ID("writeFile").Call(jen.ID("path"), jen.ID("buildBasicTableTemplate").Call(jen.ID("cfg"))), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Qual("log", "Fatal").Call(jen.ID("err")),
				),
			),
			jen.Line(),
			jen.For(jen.List(jen.ID("path"), jen.ID("cfg")).Op(":=").Range().ID("creatorConfigs")).Body(
				jen.If(jen.ID("err").Op(":=").ID("writeFile").Call(jen.ID("path"), jen.ID("buildBasicCreatorTemplate").Call(jen.ID("cfg"))), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Qual("log", "Fatal").Call(jen.ID("err")),
				),
			),
		),
		jen.Line(),
	}
}
