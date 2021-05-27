package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("errorFromResponse returns library errors according to a response's status code."),
		jen.Line(),
		jen.Func().ID("errorFromResponse").Params(jen.ID("res").Op("*").Qual("net/http", "Response")).Params(jen.ID("error")).Body(
			jen.If(jen.ID("res").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilResponse")),
			jen.Switch(jen.ID("res").Dot("StatusCode")).Body(
				jen.Case(jen.Qual("net/http", "StatusNotFound")).Body(
					jen.Return().ID("ErrNotFound")),
				jen.Case(jen.Qual("net/http", "StatusBadRequest")).Body(
					jen.Return().ID("ErrInvalidRequestInput")),
				jen.Case(jen.Qual("net/http", "StatusUnauthorized"), jen.Qual("net/http", "StatusForbidden")).Body(
					jen.Return().ID("ErrUnauthorized")),
				jen.Case(jen.Qual("net/http", "StatusInternalServerError")).Body(
					jen.Return().ID("ErrInternalServerError")),
				jen.Default().Body(
					jen.Return().ID("nil")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("argIsNotPointer checks an argument and returns whether it is a pointer."),
		jen.Line(),
		jen.Func().ID("argIsNotPointer").Params(jen.ID("i").Interface()).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.If(jen.ID("i").Op("==").ID("nil").Op("||").Qual("reflect", "TypeOf").Call(jen.ID("i")).Dot("Kind").Call().Op("!=").Qual("reflect", "Ptr")).Body(
				jen.Return().List(jen.ID("true"), jen.ID("ErrArgumentIsNotPointer"))),
			jen.Return().List(jen.ID("false"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("argIsNotNil checks an argument and returns whether it is nil."),
		jen.Line(),
		jen.Func().ID("argIsNotNil").Params(jen.ID("i").Interface()).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.If(jen.ID("i").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("true"), jen.ID("ErrNilInputProvided"))),
			jen.Return().List(jen.ID("false"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("argIsNotPointerOrNil does what it says on the tin. This function is primarily useful for detecting"),
		jen.Line(),
		jen.Comment("if a destination value is valid before decoding an HTTP response, for instance."),
		jen.Line(),
		jen.Func().ID("argIsNotPointerOrNil").Params(jen.ID("i").Interface()).Params(jen.ID("error")).Body(
			jen.If(jen.List(jen.ID("nn"), jen.ID("err")).Op(":=").ID("argIsNotNil").Call(jen.ID("i")), jen.ID("nn").Op("||").ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err")),
			jen.If(jen.List(jen.ID("np"), jen.ID("err")).Op(":=").ID("argIsNotPointer").Call(jen.ID("i")), jen.ID("np").Op("||").ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("err")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("unmarshalBody takes an HTTP response and JSON decodes its body into a destination value. The error returned here"),
		jen.Line(),
		jen.Comment("should only ever be received in testing, and should never be encountered by an end-user."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("unmarshalBody").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Op("*").Qual("net/http", "Response"), jen.ID("dest").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithResponse").Call(jen.ID("res")),
			jen.If(jen.ID("err").Op(":=").ID("argIsNotPointerOrNil").Call(jen.ID("dest")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("nil marshal target"),
				)),
			jen.List(jen.ID("bodyBytes"), jen.ID("err")).Op(":=").Qual("io", "ReadAll").Call(jen.ID("res").Dot("Body")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling error response"),
				)),
			jen.If(jen.ID("res").Dot("StatusCode").Op(">=").Qual("net/http", "StatusBadRequest")).Body(
				jen.ID("apiErr").Op(":=").Op("&").ID("types").Dot("ErrorResponse").Valuesln(jen.ID("Code").Op(":").ID("res").Dot("StatusCode")),
				jen.If(jen.ID("err").Op("=").ID("c").Dot("encoder").Dot("Unmarshal").Call(
					jen.ID("ctx"),
					jen.ID("bodyBytes"),
					jen.Op("&").ID("apiErr"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("unmarshalling error response"),
					),
					jen.ID("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("unmarshalling error response"),
					),
					jen.ID("tracing").Dot("AttachErrorToSpan").Call(
						jen.ID("span"),
						jen.Lit(""),
						jen.ID("err"),
					),
				),
				jen.Return().ID("apiErr"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("encoder").Dot("Unmarshal").Call(
				jen.ID("ctx"),
				jen.ID("bodyBytes"),
				jen.Op("&").ID("dest"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling response body"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
