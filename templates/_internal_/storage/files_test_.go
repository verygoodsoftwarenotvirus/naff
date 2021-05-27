package storage

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func filesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestUploader_ReadFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleFilename").Op(":=").Lit("hello_world.txt"),
					jen.ID("b").Op(":=").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("b").Dot("WriteAll").Call(
							jen.ID("ctx"),
							jen.ID("exampleFilename"),
							jen.Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
							jen.ID("nil"),
						),
					),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
						jen.ID("bucket").Op(":").ID("b"), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()), jen.ID("filenameFetcher").Op(":").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
							jen.Return().ID("t").Dot("Name").Call())),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("u").Dot("ReadFile").Call(
						jen.ID("ctx"),
						jen.ID("exampleFilename"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid file"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleFilename").Op(":=").Lit("hello_world.txt"),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
						jen.ID("bucket").Op(":").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values()), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()), jen.ID("filenameFetcher").Op(":").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
							jen.Return().ID("t").Dot("Name").Call())),
					jen.List(jen.ID("x"), jen.ID("err")).Op(":=").ID("u").Dot("ReadFile").Call(
						jen.ID("ctx"),
						jen.ID("exampleFilename"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUploader_SaveFile").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("io/ioutil", "TempFile").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
						jen.ID("bucket").Op(":").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values()), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()), jen.ID("filenameFetcher").Op(":").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
							jen.Return().ID("t").Dot("Name").Call())),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("u").Dot("SaveFile").Call(
							jen.ID("ctx"),
							jen.ID("tempFile").Dot("Name").Call(),
							jen.Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestUploader_ServeFiles").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleFilename").Op(":=").Lit("hello_world.txt"),
					jen.ID("b").Op(":=").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("b").Dot("WriteAll").Call(
							jen.ID("ctx"),
							jen.ID("exampleFilename"),
							jen.Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
							jen.ID("nil"),
						),
					),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
						jen.ID("bucket").Op(":").ID("b"), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()), jen.ID("filenameFetcher").Op(":").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
							jen.Return().ID("exampleFilename"))),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("u").Dot("ServeFiles").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nonexistent file"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleFilename").Op(":=").Lit("hello_world.txt"),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
						jen.ID("bucket").Op(":").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values()), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()), jen.ID("filenameFetcher").Op(":").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
							jen.Return().ID("exampleFilename"))),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("u").Dot("ServeFiles").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusNotFound"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing file content"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleFilename").Op(":=").Lit("hello_world.txt"),
					jen.ID("b").Op(":=").ID("memblob").Dot("OpenBucket").Call(jen.Op("&").ID("memblob").Dot("Options").Values()),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("b").Dot("WriteAll").Call(
							jen.ID("ctx"),
							jen.ID("exampleFilename"),
							jen.Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
							jen.ID("nil"),
						),
					),
					jen.ID("u").Op(":=").Op("&").ID("Uploader").Valuesln(
						jen.ID("bucket").Op(":").ID("b"), jen.ID("logger").Op(":").ID("logging").Dot("NewNonOperationalLogger").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("t").Dot("Name").Call()), jen.ID("filenameFetcher").Op(":").Func().Params(jen.Op("*").Qual("net/http", "Request")).Params(jen.ID("string")).Body(
							jen.Return().ID("exampleFilename"))),
					jen.ID("res").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPResponseWriter").Values(),
					jen.ID("res").Dot("On").Call(
						jen.Lit("Write"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("byte").Call(jen.ID("nil"))),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("res").Dot("On").Call(jen.Lit("Header")).Dot("Return").Call(jen.Qual("net/http", "Header").Values()),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("u").Dot("ServeFiles").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
