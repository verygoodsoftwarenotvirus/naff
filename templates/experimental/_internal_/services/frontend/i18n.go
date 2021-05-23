package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func i18NDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("translationsDir").Qual("embed", "FS"),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("provideLocalizer").Params().Params(jen.Op("*").ID("i18n").Dot("Localizer")).Body(
			jen.ID("bundle").Op(":=").ID("i18n").Dot("NewBundle").Call(jen.ID("language").Dot("English")),
			jen.ID("bundle").Dot("RegisterUnmarshalFunc").Call(
				jen.Lit("toml"),
				jen.ID("toml").Dot("Unmarshal"),
			),
			jen.List(jen.ID("translationFolderContents"), jen.ID("folderReadErr")).Op(":=").Qual("io/fs", "ReadDir").Call(
				jen.ID("translationsDir"),
				jen.Lit("translations"),
			),
			jen.If(jen.ID("folderReadErr").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(
					jen.Lit("error reading translations folder: %w"),
					jen.ID("folderReadErr"),
				))),
			jen.For(jen.List(jen.ID("_"), jen.ID("f")).Op(":=").Range().ID("translationFolderContents")).Body(
				jen.ID("translationFilename").Op(":=").Qual("path/filepath", "Join").Call(
					jen.Lit("translations"),
					jen.ID("f").Dot("Name").Call(),
				),
				jen.List(jen.ID("translationFileBytes"), jen.ID("fileReadErr")).Op(":=").Qual("io/fs", "ReadFile").Call(
					jen.ID("translationsDir"),
					jen.ID("translationFilename"),
				),
				jen.If(jen.ID("fileReadErr").Op("!=").ID("nil")).Body(
					jen.ID("panic").Call(jen.Qual("fmt", "Errorf").Call(
						jen.Lit("error reading translation file %q: %w"),
						jen.ID("translationFilename"),
						jen.ID("fileReadErr"),
					))),
				jen.ID("bundle").Dot("MustParseMessageFileBytes").Call(
					jen.ID("translationFileBytes"),
					jen.ID("f").Dot("Name").Call(),
				),
			),
			jen.Return().ID("i18n").Dot("NewLocalizer").Call(
				jen.ID("bundle"),
				jen.Lit("en"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("getSimpleLocalizedString").Params(jen.ID("messageID").ID("string")).Params(jen.ID("string")).Body(
			jen.Return().ID("s").Dot("localizer").Dot("MustLocalize").Call(jen.Op("&").ID("i18n").Dot("LocalizeConfig").Valuesln(
				jen.ID("MessageID").Op(":").ID("messageID"), jen.ID("DefaultMessage").Op(":").ID("nil"), jen.ID("TemplateData").Op(":").ID("nil"), jen.ID("Funcs").Op(":").ID("nil")))),
		jen.Line(),
	)

	return code
}
