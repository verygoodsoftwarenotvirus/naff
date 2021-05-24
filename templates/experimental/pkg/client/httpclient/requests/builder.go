package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func builderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("clientName").Op("=").Lit("todo_client_v1"),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("Builder").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
			jen.ID("url").Op("*").Qual("net/url", "URL"),
			jen.ID("encoder").ID("encoding").Dot("ClientEncoder"),
			jen.ID("panicker").ID("panicking").Dot("Panicker"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("id").Params(jen.ID("x").ID("uint64")).Params(jen.ID("string")).Body(
			jen.Return().Qual("strconv", "FormatUint").Call(
				jen.ID("x"),
				jen.Lit(10),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewBuilder builds a new API client for us."),
		jen.Line(),
		jen.Func().ID("NewBuilder").Params(jen.ID("u").Op("*").Qual("net/url", "URL"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("encoder").ID("encoding").Dot("ClientEncoder")).Params(jen.Op("*").ID("Builder"), jen.ID("error")).Body(
			jen.ID("l").Op(":=").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")),
			jen.If(jen.ID("u").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNoURLProvided"))),
			jen.If(jen.ID("encoder").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilEncoderProvided"))),
			jen.ID("c").Op(":=").Op("&").ID("Builder").Valuesln(jen.ID("url").Op(":").ID("u"), jen.ID("logger").Op(":").ID("l"), jen.ID("encoder").Op(":").ID("encoder"), jen.ID("panicker").Op(":").ID("panicking").Dot("NewProductionPanicker").Call(), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("clientName"))),
			jen.Return().List(jen.ID("c"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("URL provides the client's URL.").Params(jen.ID("b").Op("*").ID("Builder")).ID("URL").Params().Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.Return().ID("b").Dot("url")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("SetURL provides the client's URL.").Params(jen.ID("b").Op("*").ID("Builder")).ID("SetURL").Params(jen.ID("u").Op("*").Qual("net/url", "URL")).Params(jen.ID("error")).Body(
			jen.If(jen.ID("u").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNoURLProvided")),
			jen.ID("b").Dot("url").Op("=").ID("u"),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildURL builds standard service URLs."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("qp").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("u").Op(":=").ID("b").Dot("buildAPIV1URL").Call(
				jen.ID("ctx"),
				jen.ID("qp"),
				jen.ID("parts").Op("..."),
			), jen.ID("u").Op("!=").ID("nil")).Body(
				jen.Return().ID("u").Dot("String").Call()),
			jen.Return().Lit(""),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Must requires that a given request be built without error."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("Must").Params(jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("err").ID("error")).Params(jen.Op("*").Qual("net/http", "Request")).Body(
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("b").Dot("panicker").Dot("Panic").Call(jen.ID("err"))),
			jen.Return().ID("req"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildRawURL").Params(jen.ID("u").Op("*").Qual("net/url", "URL"), jen.ID("qp").Qual("net/url", "Values"), jen.ID("includeVersionPrefix").ID("bool"), jen.ID("parts").Op("...").ID("string")).Params(jen.Op("*").Qual("net/url", "URL"), jen.ID("error")).Body(
			jen.ID("tu").Op(":=").Op("*").ID("u"),
			jen.If(jen.ID("includeVersionPrefix")).Body(
				jen.ID("parts").Op("=").ID("append").Call(
					jen.Index().ID("string").Valuesln(jen.Lit("api"), jen.Lit("v1")),
					jen.ID("parts").Op("..."),
				)),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Qual("path", "Join").Call(jen.ID("parts").Op("..."))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("err"))),
			jen.If(jen.ID("qp").Op("!=").ID("nil")).Body(
				jen.ID("u").Dot("RawQuery").Op("=").ID("qp").Dot("Encode").Call()),
			jen.Return().List(jen.ID("tu").Dot("ResolveReference").Call(jen.ID("u")), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildRawURL takes a given set of query parameters and url parts, and returns a parsed url object from them."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("buildAPIV1URL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("queryParams").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tu").Op(":=").Op("*").ID("b").Dot("url"),
			jen.ID("parts").Op("=").ID("append").Call(
				jen.Index().ID("string").Valuesln(jen.Lit("api"), jen.Lit("v1")),
				jen.ID("parts").Op("..."),
			),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.Qual("path", "Join").Call(jen.ID("parts").Op("..."))),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("b").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("building url"),
				),
				jen.Return().ID("nil"),
			),
			jen.If(jen.ID("queryParams").Op("!=").ID("nil")).Body(
				jen.ID("u").Dot("RawQuery").Op("=").ID("queryParams").Dot("Encode").Call()),
			jen.ID("out").Op(":=").ID("tu").Dot("ResolveReference").Call(jen.ID("u")),
			jen.ID("tracing").Dot("AttachURLToSpan").Call(
				jen.ID("span"),
				jen.ID("out"),
			),
			jen.Return().ID("out"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("buildUnversionedURL builds a url without the v1 API prefix. It should otherwise be identical to buildRawURL.").Params(jen.ID("b").Op("*").ID("Builder")).ID("buildUnversionedURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("qp").Qual("net/url", "Values"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").ID("buildRawURL").Call(
				jen.ID("b").Dot("url"),
				jen.ID("qp"),
				jen.ID("false"),
				jen.ID("parts").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("b").Dot("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("building unversioned url"),
				),
				jen.Return().Lit(""),
			),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildWebsocketURL builds a standard url and then converts its scheme to the websocket protocol."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildWebsocketURL").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("parts").Op("...").ID("string")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("u").Op(":=").ID("b").Dot("buildAPIV1URL").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
				jen.ID("parts").Op("..."),
			),
			jen.ID("u").Dot("Scheme").Op("=").Lit("ws"),
			jen.Return().ID("u").Dot("String").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildHealthCheckRequest builds a health check HTTP request."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("BuildHealthCheckRequest").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("u").Op(":=").Op("*").ID("b").Dot("url"),
			jen.ID("uri").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s://%s/_meta_/ready"),
				jen.ID("u").Dot("Scheme"),
				jen.ID("u").Dot("Host"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.Qual("net/http", "MethodGet"),
				jen.ID("uri"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("b").Dot("logger"),
					jen.ID("span"),
					jen.Lit("building user status request"),
				))),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("buildDataRequest builds an HTTP request for a given method, url, and body data."),
		jen.Line(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Builder")).ID("buildDataRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("method"), jen.ID("uri")).ID("string"), jen.ID("in").Interface()).Params(jen.Op("*").Qual("net/http", "Request"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("b").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequestMethodKey"),
				jen.ID("method"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("URLKey"),
				jen.ID("uri"),
			),
			jen.List(jen.ID("body"), jen.ID("err")).Op(":=").ID("b").Dot("encoder").Dot("EncodeReader").Call(
				jen.ID("ctx"),
				jen.ID("in"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding request"),
				))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
				jen.ID("ctx"),
				jen.ID("method"),
				jen.ID("uri"),
				jen.ID("body"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building request"),
				))),
			jen.ID("req").Dot("Header").Dot("Set").Call(
				jen.Lit("RawHTML-type"),
				jen.ID("b").Dot("encoder").Dot("ContentType").Call(),
			),
			jen.ID("tracing").Dot("AttachURLToSpan").Call(
				jen.ID("span"),
				jen.ID("req").Dot("URL"),
			),
			jen.Return().List(jen.ID("req"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
