package iterables

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(typ.Name.PackageName())

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("Config configures the service."),
		jen.Newline(),
		jen.Type().ID("Config").Struct(
			jen.ID("Logging").Qual(proj.InternalLoggingPackage(), "Config").Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("Logging"), false)),
			func() jen.Code {
				if typ.SearchEnabled {
					return jen.ID("SearchIndexPath").String().Tag(utils.BuildStructTag(wordsmith.FromSingularPascalCase("SearchIndexPath"), false))
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Underscore().Qual(constants.ValidationLibrary, "ValidatableWithContext").Equals().Parens(jen.PointerTo().ID("Config")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a Config struct."),
		jen.Newline(),
		jen.Func().Params(jen.ID("cfg").PointerTo().ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual(constants.ValidationLibrary, "ValidateStructWithContext").Callln(
				jen.ID("ctx"),
				jen.ID("cfg"),
				jen.Qual(constants.ValidationLibrary, "Field").Call(
					jen.AddressOf().ID("cfg").Dot("Logging"),
					jen.Qual(constants.ValidationLibrary, "Required"),
				),
				func() jen.Code {
					if typ.SearchEnabled {
						return jen.Qual(constants.ValidationLibrary, "Field").Call(
							jen.AddressOf().ID("cfg").Dot("SearchIndexPath"),
							jen.Qual(constants.ValidationLibrary, "Required"),
						)
					}
					return jen.Null()
				}(),
			)),
		jen.Newline(),
	)

	return code
}
