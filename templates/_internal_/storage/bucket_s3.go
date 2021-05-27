package storage

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bucketS3DotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("S3Provider").Op("=").Lit("s3"),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("S3Config").Struct(jen.ID("BucketName").ID("string")),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("S3Config")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the S3Config."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("S3Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("BucketName"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
