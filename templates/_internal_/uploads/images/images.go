package images

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func imagesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("headerContentType").Op("=").Lit("Content-Type"),
			jen.ID("imagePNG").Op("=").Lit("image/png"),
			jen.ID("imageJPEG").Op("=").Lit("image/jpeg"),
			jen.ID("imageGIF").Op("=").Lit("image/gif"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("ErrInvalidContentType").Op("=").Qual("errors", "New").Call(jen.Lit("invalid content type")),
			jen.ID("ErrInvalidImageContentType").Op("=").Qual("errors", "New").Call(jen.Lit("invalid image content type")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("Image").Struct(
				jen.ID("Filename").ID("string"),
				jen.ID("ContentType").ID("string"),
				jen.ID("Data").Index().ID("byte"),
				jen.ID("Size").ID("int"),
			),
			jen.ID("ImageUploadProcessor").Interface(jen.ID("Process").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error"))),
			jen.ID("uploadProcessor").Struct(
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("logger").ID("logging").Dot("Logger"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DataURI converts image to base64 data URI."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("Image")).ID("DataURI").Params().Params(jen.ID("string")).Body(
			jen.Return().Qual("fmt", "Sprintf").Call(
				jen.Lit("data:%s;base64,%s"),
				jen.ID("i").Dot("ContentType"),
				jen.Qual("encoding/base64", "StdEncoding").Dot("EncodeToString").Call(jen.ID("i").Dot("Data")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Write image to HTTP response."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("Image")).ID("Write").Params(jen.ID("w").Qual("net/http", "ResponseWriter")).Params(jen.ID("error")).Body(
			jen.ID("w").Dot("Header").Call().Dot("Set").Call(
				jen.ID("headerContentType"),
				jen.ID("i").Dot("ContentType"),
			),
			jen.ID("w").Dot("Header").Call().Dot("Set").Call(
				jen.Lit("RawHTML-Length"),
				jen.Qual("strconv", "Itoa").Call(jen.ID("i").Dot("Size")),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("w").Dot("Write").Call(jen.ID("i").Dot("Data")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("writing image to HTTP response: %w"),
					jen.ID("err"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Thumbnail creates a thumbnail from an image."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("Image")).ID("Thumbnail").Params(jen.List(jen.ID("width"), jen.ID("height")).ID("uint"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error")).Body(
			jen.List(jen.ID("t"), jen.ID("err")).Op(":=").ID("newThumbnailer").Call(jen.ID("i").Dot("ContentType")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.Return().ID("t").Dot("Thumbnail").Call(
				jen.ID("i"),
				jen.ID("width"),
				jen.ID("height"),
				jen.ID("filename"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewImageUploadProcessor provides a new ImageUploadProcessor."),
		jen.Line(),
		jen.Func().ID("NewImageUploadProcessor").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("ImageUploadProcessor")).Body(
			jen.Return().Op("&").ID("uploadProcessor").Valuesln(
				jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.Lit("image_upload_processor")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("image_upload_processor")))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LimitFileSize limits the size of uploaded files, for use before Process."),
		jen.Line(),
		jen.Func().ID("LimitFileSize").Params(jen.ID("maxSize").ID("uint16"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.If(jen.ID("maxSize").Op("==").Lit(0)).Body(
				jen.ID("maxSize").Op("=").Lit(4096)),
			jen.ID("req").Dot("Body").Op("=").Qual("net/http", "MaxBytesReader").Call(
				jen.ID("res"),
				jen.ID("req").Dot("Body"),
				jen.ID("int64").Call(jen.ID("maxSize")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("contentTypeFromFilename").Params(jen.ID("filename").ID("string")).Params(jen.ID("string")).Body(
			jen.Return().Qual("mime", "TypeByExtension").Call(jen.Qual("path/filepath", "Ext").Call(jen.ID("filename")))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("validateContentType").Params(jen.ID("filename").ID("string")).Params(jen.ID("error")).Body(
			jen.ID("contentType").Op(":=").ID("contentTypeFromFilename").Call(jen.ID("filename")),
			jen.Switch(jen.Qual("strings", "TrimSpace").Call(jen.Qual("strings", "ToLower").Call(jen.ID("contentType")))).Body(
				jen.Case(jen.ID("imagePNG"), jen.ID("imageJPEG"), jen.ID("imageGIF")).Body(
					jen.Return().ID("nil")),
				jen.Default().Body(
					jen.Return().Qual("fmt", "Errorf").Call(
						jen.Lit("%w: %s"),
						jen.ID("ErrInvalidContentType"),
						jen.ID("contentType"),
					)),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Process extracts an image from an *http.Request."),
		jen.Line(),
		jen.Func().Params(jen.ID("p").Op("*").ID("uploadProcessor")).ID("Process").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("filename").ID("string")).Params(jen.Op("*").ID("Image"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("p").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("p").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.List(jen.ID("file"), jen.ID("info"), jen.ID("err")).Op(":=").ID("req").Dot("FormFile").Call(jen.ID("filename")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("parsing file from request"),
				))),
			jen.If(jen.ID("contentTypeErr").Op(":=").ID("validateContentType").Call(jen.ID("info").Dot("Filename")), jen.ID("contentTypeErr").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("contentTypeErr"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating the content type"),
				))),
			jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("io", "ReadAll").Call(jen.ID("file")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading file from request"),
				))),
			jen.If(jen.List(jen.ID("_"), jen.ID("_"), jen.ID("err")).Op("=").Qual("image", "Decode").Call(jen.Qual("bytes", "NewReader").Call(jen.ID("bs"))), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding the image data"),
				))),
			jen.ID("i").Op(":=").Op("&").ID("Image").Valuesln(
				jen.ID("Filename").Op(":").ID("info").Dot("Filename"), jen.ID("ContentType").Op(":").ID("contentTypeFromFilename").Call(jen.ID("filename")), jen.ID("Data").Op(":").ID("bs"), jen.ID("Size").Op(":").ID("len").Call(jen.ID("bs"))),
			jen.Return().List(jen.ID("i"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
