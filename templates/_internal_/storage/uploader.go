package storage

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func uploaderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("MemoryProvider").Op("=").Lit("memory"),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("ErrNilConfig").Op("=").Qual("errors", "New").Call(jen.Lit("nil config provided")),
			jen.ID("ErrInvalidConfiguration").Op("=").Qual("errors", "New").Call(jen.Lit("configuration invalid")),
			jen.ID("ErrBucketIsUnavailable").Op("=").Qual("errors", "New").Call(jen.Lit("bucket is unavailable")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Uploader").Struct(
				jen.ID("bucket").Op("*").Qual("gocloud.dev/blob", "Bucket"),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("filenameFetcher").Func().Params(jen.ID("req").Op("*").Qual("net/http", "Request")).Params(jen.ID("string")),
			),
			jen.ID("Config").Struct(
				jen.ID("FilesystemConfig").Op("*").ID("FilesystemConfig"),
				jen.ID("AzureConfig").Op("*").ID("AzureConfig"),
				jen.ID("GCSConfig").Op("*").ID("GCSConfig"),
				jen.ID("S3Config").Op("*").ID("S3Config"),
				jen.ID("BucketName").ID("string"),
				jen.ID("UploadFilenameKey").ID("string"),
				jen.ID("Provider").ID("string"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("Config")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates the Config."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Config")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("c"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("BucketName"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("Provider"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "In").Call(
						jen.ID("AzureProvider"),
						jen.ID("GCSProvider"),
						jen.ID("S3Provider"),
						jen.ID("FilesystemProvider"),
						jen.ID("MemoryProvider"),
					),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("AzureConfig"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("Provider").Op("==").ID("AzureProvider"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Nil")),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("GCSConfig"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("Provider").Op("==").ID("GCSProvider"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Nil")),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("S3Config"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("Provider").Op("==").ID("S3Provider"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Nil")),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("c").Dot("FilesystemConfig"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "When").Call(
						jen.ID("c").Dot("Provider").Op("==").ID("FilesystemProvider"),
						jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
					).Dot("Else").Call(jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Nil")),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewUploadManager provides a new uploads.UploadManager."),
		jen.Line(),
		jen.Func().ID("NewUploadManager").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("cfg").Op("*").ID("Config"), jen.ID("routeParamManager").ID("routing").Dot("RouteParamManager")).Params(jen.Op("*").ID("Uploader"), jen.ID("error")).Body(
			jen.If(jen.ID("cfg").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilConfig"))),
			jen.ID("serviceName").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s_uploader"),
				jen.ID("cfg").Dot("BucketName"),
			),
			jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
				jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")), jen.ID("filenameFetcher").Op(":").ID("routeParamManager").Dot("BuildRouteParamStringIDFetcher").Call(jen.ID("cfg").Dot("UploadFilenameKey"))),
			jen.If(jen.ID("err").Op(":=").ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("upload manager provided invalid config: %w"),
					jen.ID("err"),
				))),
			jen.If(jen.ID("err").Op(":=").ID("u").Dot("selectBucket").Call(
				jen.ID("ctx"),
				jen.ID("cfg"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("initializing bucket: %w"),
					jen.ID("err"),
				))),
			jen.If(jen.List(jen.ID("available"), jen.ID("err")).Op(":=").ID("u").Dot("bucket").Dot("IsAccessible").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("verifying bucket accessibility: %w"),
					jen.ID("err"),
				))).Else().If(jen.Op("!").ID("available")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrBucketIsUnavailable"))),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("u").Op("*").ID("Uploader")).ID("selectBucket").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("cfg").Op("*").ID("Config")).Params(jen.ID("err").ID("error")).Body(
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("cfg").Dot("Provider")))).Body(
				jen.Case(jen.ID("AzureProvider")).Body(
					jen.If(jen.ID("cfg").Dot("AzureConfig").Op("==").ID("nil")).Body(
						jen.Return().ID("ErrNilConfig")), jen.If(jen.List(jen.ID("u").Dot("bucket"), jen.ID("err")).Op("=").ID("provideAzureBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg").Dot("AzureConfig"),
						jen.ID("u").Dot("logger"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().Qual("fmt", "Errorf").Call(
							jen.Lit("initializing azure bucket: %w"),
							jen.ID("err"),
						))),
				jen.Case(jen.ID("GCSProvider")).Body(
					jen.If(jen.ID("cfg").Dot("GCSConfig").Op("==").ID("nil")).Body(
						jen.Return().ID("ErrNilConfig")), jen.If(jen.List(jen.ID("u").Dot("bucket"), jen.ID("err")).Op("=").ID("buildGCSBucket").Call(
						jen.ID("ctx"),
						jen.ID("cfg").Dot("GCSConfig"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().Qual("fmt", "Errorf").Call(
							jen.Lit("initializing gcs bucket: %w"),
							jen.ID("err"),
						))),
				jen.Case(jen.ID("S3Provider")).Body(
					jen.If(jen.ID("cfg").Dot("S3Config").Op("==").ID("nil")).Body(
						jen.Return().ID("ErrNilConfig")), jen.If(jen.List(jen.ID("u").Dot("bucket"), jen.ID("err")).Op("=").ID("s3blob").Dot("OpenBucket").Call(
						jen.ID("ctx"),
						jen.ID("session").Dot("Must").Call(jen.ID("session").Dot("NewSession").Call()),
						jen.ID("cfg").Dot("S3Config").Dot("BucketName"),
						jen.Op("&").ID("s3blob").Dot("Options").Valuesln(
							jen.ID("UseLegacyList").Op(":").ID("false")),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().Qual("fmt", "Errorf").Call(
							jen.Lit("initializing s3 bucket: %w"),
							jen.ID("err"),
						))),
				jen.Case(jen.ID("MemoryProvider")).Body(
					jen.ID("u").Dot("bucket").Op("=").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values())),
				jen.Default().Body(
					jen.If(jen.ID("cfg").Dot("FilesystemConfig").Op("==").ID("nil")).Body(
						jen.Return().ID("ErrNilConfig")), jen.If(jen.List(jen.ID("u").Dot("bucket"), jen.ID("err")).Op("=").ID("fileblob").Dot("OpenBucket").Call(
						jen.ID("cfg").Dot("FilesystemConfig").Dot("RootDirectory"),
						jen.Op("&").ID("fileblob").Dot("Options").Valuesln(
							jen.ID("URLSigner").Op(":").ID("nil"), jen.ID("CreateDir").Op(":").ID("true")),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().Qual("fmt", "Errorf").Call(
							jen.Lit("initializing filesystem bucket: %w"),
							jen.ID("err"),
						))),
			),
			jen.ID("bn").Op(":=").ID("cfg").Dot("BucketName"),
			jen.If(jen.Op("!").Qual("strings", "HasSuffix").Call(
				jen.ID("bn"),
				jen.Lit("_"),
			)).Body(
				jen.ID("bn").Op("=").Qual("fmt", "Sprintf").Call(
					jen.Lit("%s_"),
					jen.ID("cfg").Dot("BucketName"),
				)),
			jen.ID("u").Dot("bucket").Op("=").Qual("gocloud.dev/blob", "PrefixedBucket").Call(
				jen.ID("u").Dot("bucket"),
				jen.ID("bn"),
			),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)

	return code
}
