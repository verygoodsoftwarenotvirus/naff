package images

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func imagesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("newAvatarUploadRequest").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("filename").ID("string"), jen.ID("avatar").Qual("io", "Reader")).Params(jen.Op("*").Qual("net/http", "Request")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
			jen.ID("body").Op(":=").Op("&").Qual("bytes", "Buffer").Valuesln(),
			jen.ID("writer").Op(":=").Qual("mime/multipart", "NewWriter").Call(jen.ID("body")),
			jen.List(jen.ID("part"), jen.ID("err")).Op(":=").ID("writer").Dot("CreateFormFile").Call(
				jen.Lit("avatar"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("avatar.%s"),
					jen.Qual("path/filepath", "Ext").Call(jen.ID("filename")),
				),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.List(jen.ID("_"), jen.ID("err")).Op("=").Qual("io", "Copy").Call(
				jen.ID("part"),
				jen.ID("avatar"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("writer").Dot("Close").Call(),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodPost"),
				jen.Lit("https://tests.verygoodsoftwarenotvirus.ru"),
				jen.ID("body"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.ID("headerContentType"),
				jen.ID("writer").Dot("FormDataContentType").Call(),
			),
			jen.Return().ID("req"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildPNGBytes").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("bytes", "Buffer")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("b").Op(":=").ID("new").Call(jen.Qual("bytes", "Buffer")),
			jen.ID("exampleImage").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildArbitraryImage").Call(jen.Lit(256)),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.Qual("image/png", "Encode").Call(
					jen.ID("b"),
					jen.ID("exampleImage"),
				),
			),
			jen.ID("expected").Op(":=").ID("b").Dot("Bytes").Call(),
			jen.Return().Qual("bytes", "NewBuffer").Call(jen.ID("expected")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildJPEGBytes").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("bytes", "Buffer")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("b").Op(":=").ID("new").Call(jen.Qual("bytes", "Buffer")),
			jen.ID("exampleImage").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildArbitraryImage").Call(jen.Lit(256)),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.Qual("image/jpeg", "Encode").Call(
					jen.ID("b"),
					jen.ID("exampleImage"),
					jen.Op("&").Qual("image/jpeg", "Options").Valuesln(jen.ID("Quality").Op(":").Qual("image/jpeg", "DefaultQuality")),
				),
			),
			jen.ID("expected").Op(":=").ID("b").Dot("Bytes").Call(),
			jen.Return().Qual("bytes", "NewBuffer").Call(jen.ID("expected")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildGIFBytes").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").Qual("bytes", "Buffer")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("b").Op(":=").ID("new").Call(jen.Qual("bytes", "Buffer")),
			jen.ID("exampleImage").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildArbitraryImage").Call(jen.Lit(256)),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.Qual("image/gif", "Encode").Call(
					jen.ID("b"),
					jen.ID("exampleImage"),
					jen.Op("&").Qual("image/gif", "Options").Valuesln(jen.ID("NumColors").Op(":").Lit(256)),
				),
			),
			jen.ID("expected").Op(":=").ID("b").Dot("Bytes").Call(),
			jen.Return().Qual("bytes", "NewBuffer").Call(jen.ID("expected")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestImage_DataURI").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").ID("t").Dot("Name").Call(), jen.ID("ContentType").Op(":").Lit("things/stuff"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(12345)),
					jen.ID("expected").Op(":=").Lit("data:things/stuff;base64,VGVzdEltYWdlX0RhdGFVUkkvc3RhbmRhcmQ="),
					jen.ID("actual").Op(":=").ID("i").Dot("DataURI").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestImage_Write").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").ID("t").Dot("Name").Call(), jen.ID("ContentType").Op(":").Lit("things/stuff"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(12345)),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("i").Dot("Write").Call(jen.ID("res")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with write error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").ID("t").Dot("Name").Call(), jen.ID("ContentType").Op(":").Lit("things/stuff"), jen.ID("Data").Op(":").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()), jen.ID("Size").Op(":").Lit(12345)),
					jen.ID("res").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPResponseWriter").Valuesln(),
					jen.ID("res").Dot("On").Call(jen.Lit("Header")).Dot("Return").Call(jen.Qual("net/http", "Header").Valuesln()),
					jen.ID("res").Dot("On").Call(
						jen.Lit("Write"),
						jen.ID("mock").Dot("IsType").Call(jen.Index().ID("byte").Call(jen.ID("nil"))),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("i").Dot("Write").Call(jen.ID("res")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestImage_Thumbnail").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("imgBytes").Op(":=").ID("buildPNGBytes").Call(jen.ID("t")).Dot("Bytes").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("Filename").Op(":").ID("t").Dot("Name").Call(), jen.ID("ContentType").Op(":").ID("imagePNG"), jen.ID("Data").Op(":").ID("imgBytes"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("imgBytes"))),
					jen.List(jen.ID("tempFile"), jen.ID("err")).Op(":=").Qual("os", "CreateTemp").Call(
						jen.Lit(""),
						jen.Lit(""),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("i").Dot("Thumbnail").Call(
						jen.Lit(123),
						jen.Lit(123),
						jen.ID("tempFile").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.Qual("os", "Remove").Call(jen.ID("tempFile").Dot("Name").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content type"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(jen.ID("ContentType").Op(":").ID("t").Dot("Name").Call()),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("i").Dot("Thumbnail").Call(
						jen.Lit(123),
						jen.Lit(123),
						jen.ID("t").Dot("Name").Call(),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestLimitFileSize").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("imgBytes").Op(":=").ID("buildPNGBytes").Call(jen.ID("t")),
					jen.ID("req").Op(":=").ID("newAvatarUploadRequest").Call(
						jen.ID("t"),
						jen.Lit("avatar.png"),
						jen.ID("imgBytes"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("LimitFileSize").Call(
						jen.Lit(0),
						jen.ID("res"),
						jen.ID("req"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_uploadProcessor_Process").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("p").Op(":=").ID("NewImageUploadProcessor").Call(jen.ID("nil")),
					jen.ID("expectedFieldName").Op(":=").Lit("avatar"),
					jen.ID("imgBytes").Op(":=").ID("buildPNGBytes").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").ID("imgBytes").Dot("Bytes").Call(),
					jen.ID("req").Op(":=").ID("newAvatarUploadRequest").Call(
						jen.ID("t"),
						jen.Lit("avatar.png"),
						jen.ID("imgBytes"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot("Process").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("expectedFieldName"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual").Dot("Data"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing form file"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("p").Op(":=").ID("NewImageUploadProcessor").Call(jen.ID("nil")),
					jen.ID("expectedFieldName").Op(":=").Lit("avatar"),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("https://tests.verygoodsoftwarenotvirus.ru"),
						jen.ID("nil"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot("Process").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("expectedFieldName"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid content type"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("p").Op(":=").ID("NewImageUploadProcessor").Call(jen.ID("nil")),
					jen.ID("expectedFieldName").Op(":=").Lit("avatar"),
					jen.ID("imgBytes").Op(":=").ID("buildPNGBytes").Call(jen.ID("t")),
					jen.ID("req").Op(":=").ID("newAvatarUploadRequest").Call(
						jen.ID("t"),
						jen.Lit("avatar.pizza"),
						jen.ID("imgBytes"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot("Process").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("expectedFieldName"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error decoding image"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("p").Op(":=").ID("NewImageUploadProcessor").Call(jen.ID("nil")),
					jen.ID("expectedFieldName").Op(":=").Lit("avatar"),
					jen.ID("req").Op(":=").ID("newAvatarUploadRequest").Call(
						jen.ID("t"),
						jen.Lit("avatar.png"),
						jen.Qual("bytes", "NewBufferString").Call(jen.Lit("")),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("p").Dot("Process").Call(
						jen.ID("ctx"),
						jen.ID("req"),
						jen.ID("expectedFieldName"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
