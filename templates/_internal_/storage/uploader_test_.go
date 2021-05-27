package storage

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func uploaderTestDotGo(proj *models.Project) *jen.File {
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
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("Provider").Op(":").ID("FilesystemProvider"), jen.ID("FilesystemConfig").Op(":").Op("&").ID("FilesystemConfig").Valuesln(
							jen.ID("RootDirectory").Op(":").ID("t").Dot("Name").Call())),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("cfg").Dot("ValidateWithContext").Call(jen.ID("ctx")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestNewUploadManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call(), jen.ID("Provider").Op(":").ID("MemoryProvider")),
					jen.ID("rpm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "RouteParamManager").Values(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamStringIDFetcher"),
						jen.ID("cfg").Dot("UploadFilenameKey"),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
						jen.Return().ID("t").Dot("Name").Call())),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("NewUploadManager").Call(
						jen.ID("ctx"),
						jen.ID("l"),
						jen.ID("cfg"),
						jen.ID("rpm"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("rpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("rpm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "RouteParamManager").Values(),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("NewUploadManager").Call(
						jen.ID("ctx"),
						jen.ID("l"),
						jen.ID("nil"),
						jen.ID("rpm"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("rpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Values(),
					jen.ID("rpm").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/routing/mock", "RouteParamManager").Values(),
					jen.ID("rpm").Dot("On").Call(
						jen.Lit("BuildRouteParamStringIDFetcher"),
						jen.ID("cfg").Dot("UploadFilenameKey"),
					).Dot("Return").Call(jen.Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
						jen.Return().ID("t").Dot("Name").Call())),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("NewUploadManager").Call(
						jen.ID("ctx"),
						jen.ID("l"),
						jen.ID("cfg"),
						jen.ID("rpm"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("rpm"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUploader_selectBucket").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("azure happy path"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("AzureProvider"), jen.ID("AzureConfig").Op(":").Op("&").ID("AzureConfig").Valuesln(
							jen.ID("AccountName").Op(":").Lit("blah"), jen.ID("BucketName").Op(":").Lit("blahs"), jen.ID("Retrying").Op(":").Op("&").ID("AzureRetryConfig").Values())),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("azure with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("AzureProvider"), jen.ID("AzureConfig").Op(":").ID("nil")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("gcs with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("GCSProvider"), jen.ID("GCSConfig").Op(":").ID("nil")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("s3 happy path"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("S3Provider"), jen.ID("S3Config").Op(":").Op("&").ID("S3Config").Valuesln(
							jen.ID("BucketName").Op(":").ID("t").Dot("Name").Call())),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("s3 with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("S3Provider"), jen.ID("S3Config").Op(":").ID("nil")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("memory provider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("MemoryProvider")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("filesystem happy path"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("tempDir").Op(":=").Qual("os", "TempDir").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("FilesystemProvider"), jen.ID("FilesystemConfig").Op(":").Op("&").ID("FilesystemConfig").Valuesln(
							jen.ID("RootDirectory").Op(":").ID("tempDir"))),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("filesystem with nil config"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Values(),
					jen.ID("cfg").Op(":=").Op("&").ID("Config").Valuesln(
						jen.ID("Provider").Op(":").ID("FilesystemProvider"), jen.ID("FilesystemConfig").Op(":").ID("nil")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("u").Dot("selectBucket").Call(
							jen.ID("ctx"),
							jen.ID("cfg"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
