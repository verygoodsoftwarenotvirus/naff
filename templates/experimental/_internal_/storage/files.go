package storage

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func filesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("SaveFile saves a file to the blob."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").Op("*").ID("Uploader")).ID("SaveFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("path").ID("string"), jen.ID("content").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("u").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("err").Op(":=").ID("u").Dot("bucket").Dot("WriteAll").Call(
				jen.ID("ctx"),
				jen.ID("path"),
				jen.ID("content"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("writing file content: %w"),
					jen.ID("err"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ReadFile reads a file from the blob."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").Op("*").ID("Uploader")).ID("ReadFile").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("path").ID("string")).Params(jen.Index().ID("byte"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("u").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("r"), jen.ID("err")).Op(":=").ID("u").Dot("bucket").Dot("NewReader").Call(
				jen.ID("ctx"),
				jen.ID("path"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("fetching file: %w"),
					jen.ID("err"),
				))),
			jen.Defer().Func().Params().Body(
				jen.If(jen.ID("closeErr").Op(":=").ID("r").Dot("Close").Call(), jen.ID("closeErr").Op("!=").ID("nil")).Body(
					jen.ID("u").Dot("logger").Dot("Error").Call(
						jen.ID("closeErr"),
						jen.Lit("error closing file reader"),
					))).Call(),
			jen.List(jen.ID("fileBytes"), jen.ID("err")).Op(":=").Qual("io", "ReadAll").Call(jen.ID("r")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("reading file: %w"),
					jen.ID("err"),
				))),
			jen.Return().List(jen.ID("fileBytes"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ServeFiles saves a file to the blob."),
		jen.Line(),
		jen.Func().Params(jen.ID("u").Op("*").ID("Uploader")).ID("ServeFiles").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("u").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("fileName").Op(":=").ID("u").Dot("filenameFetcher").Call(jen.ID("req")),
			jen.List(jen.ID("fileBytes"), jen.ID("err")).Op(":=").ID("u").Dot("ReadFile").Call(
				jen.ID("ctx"),
				jen.ID("fileName"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("u").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("trying to read uploaded file"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.Qual("net/http", "StatusNotFound")),
				jen.Return(),
			),
			jen.If(jen.List(jen.ID("attrs"), jen.ID("err")).Op(":=").ID("u").Dot("bucket").Dot("Attributes").Call(
				jen.ID("ctx"),
				jen.ID("fileName"),
			), jen.ID("attrs").Op("!=").ID("nil").Op("&&").ID("err").Op("==").ID("nil")).Body(
				jen.ID("res").Dot("Header").Call().Dot("Set").Call(
					jen.ID("encoding").Dot("ContentTypeHeaderKey"),
					jen.ID("attrs").Dot("ContentType"),
				)),
			jen.If(jen.List(jen.ID("_"), jen.ID("copyErr")).Op(":=").Qual("io", "Copy").Call(
				jen.ID("res"),
				jen.Qual("bytes", "NewReader").Call(jen.ID("fileBytes")),
			), jen.ID("copyErr").Op("!=").ID("nil")).Body(
				jen.ID("u").Dot("logger").Dot("Error").Call(
					jen.ID("copyErr"),
					jen.Lit("copying file bytes to response"),
				)),
		),
		jen.Line(),
	)

	return code
}
