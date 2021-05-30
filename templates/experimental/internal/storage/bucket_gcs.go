package storage

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bucketGcsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("GCSProvider").Op("=").Lit("gcs"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("GCSBlobConfig").Struct(jen.ID("GoogleAccessID").ID("string")),
			jen.ID("GCSConfig").Struct(
				jen.ID("BlobSettings").ID("GCSBlobConfig"),
				jen.ID("ServiceAccountKeyFilepath").ID("string"),
				jen.ID("BucketName").ID("string"),
				jen.ID("Scopes").Index().ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildGCSBucket").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("GCSConfig")).Params(jen.Op("*").Qual("gocloud.dev/blob", "Bucket"), jen.ID("error")).Body(
			jen.Var().Defs(
				jen.ID("creds").Op("*").ID("google").Dot("Credentials"),
				jen.ID("bucket").Op("*").Qual("gocloud.dev/blob", "Bucket"),
			),
			jen.If(jen.ID("cfg").Dot("ServiceAccountKeyFilepath").Op("!=").Lit("")).Body(
				jen.List(jen.ID("serviceAccountKeyBytes"), jen.ID("err")).Op(":=").Qual("os", "ReadFile").Call(jen.ID("cfg").Dot("ServiceAccountKeyFilepath")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("reading service account key file: %w"),
						jen.ID("err"),
					))),
				jen.If(jen.List(jen.ID("creds"), jen.ID("err")).Op("=").ID("google").Dot("CredentialsFromJSON").Call(
					jen.ID("ctx"),
					jen.ID("serviceAccountKeyBytes"),
					jen.ID("cfg").Dot("Scopes").Op("..."),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("using service account key credentials: %w"),
						jen.ID("err"),
					))),
			).Else().Body(
				jen.Var().Defs(
					jen.ID("err").ID("error"),
				),
				jen.If(jen.List(jen.ID("creds"), jen.ID("err")).Op("=").Qual("gocloud.dev/gcp", "DefaultCredentials").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
						jen.Lit("constructing GCPKMS credentials: %w"),
						jen.ID("err"),
					))),
			),
			jen.List(jen.ID("gcsClient"), jen.ID("gcsClientErr")).Op(":=").Qual("gocloud.dev/gcp", "NewHTTPClient").Call(
				jen.ID("nil"),
				jen.Qual("gocloud.dev/gcp", "CredentialsTokenSource").Call(jen.ID("creds")),
			),
			jen.If(jen.ID("gcsClientErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("constructing GCPKMS client: %w"),
					jen.ID("gcsClientErr"),
				))),
			jen.ID("blobOpts").Op(":=").Op("&").ID("gcsblob").Dot("Options").Valuesln(jen.ID("GoogleAccessID").Op(":").ID("cfg").Dot("BlobSettings").Dot("GoogleAccessID")),
			jen.List(jen.ID("bucket"), jen.ID("err")).Op(":=").ID("gcsblob").Dot("OpenBucket").Call(
				jen.ID("ctx"),
				jen.ID("gcsClient"),
				jen.ID("cfg").Dot("BucketName"),
				jen.ID("blobOpts"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("initializing filesystem bucket: %w"),
					jen.ID("err"),
				))),
			jen.Return().List(jen.ID("bucket"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("GCSConfig")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the GCSConfig."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("GCSConfig")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
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
