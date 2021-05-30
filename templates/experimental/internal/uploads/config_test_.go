package uploads

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func configTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestConfig_ValidateWithContext").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(jen.ID("Storage").Op(":").ID("storage").Dot("Config").Valuesln(jen.ID("FilesystemConfig").Op(":").Op("&").ID("storage").Dot("FilesystemConfig").Valuesln(jen.ID("RootDirectory").Op(":").Lit("/blah")), jen.ID("AzureConfig").Op(":").Op("&").ID("storage").Dot("AzureConfig").Valuesln(jen.ID("BucketName").Op(":").Lit("blahs"), jen.ID("Retrying").Op(":").Op("&").ID("storage").Dot("AzureRetryConfig").Valuesln()), jen.ID("GCSConfig").Op(":").Op("&").ID("storage").Dot("GCSConfig").Valuesln(jen.ID("ServiceAccountKeyFilepath").Op(":").Lit("/blah/blah"), jen.ID("BucketName").Op(":").Lit("blah"), jen.ID("Scopes").Op(":").ID("nil")), jen.ID("S3Config").Op(":").Op("&").ID("storage").Dot("S3Config").Valuesln(jen.ID("BucketName").Op(":").Lit("blahs")), jen.ID("BucketName").Op(":").Lit("blahs"), jen.ID("UploadFilenameKey").Op(":").Lit("blahs"), jen.ID("Provider").Op(":").Lit("blahs")), jen.ID("Debug").Op(":").ID("false")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
