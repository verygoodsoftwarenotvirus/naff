package storage

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bucketFilesystemDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("FilesystemProvider").Op("=").Lit("filesystem"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("FilesystemConfig").Struct(jen.ID("RootDirectory").ID("string")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("FilesystemConfig")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the FilesystemConfig."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("FilesystemConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("RootDirectory"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
