package encoding

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func clientEncoderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("ClientEncoder").Interface(
				jen.ID("ContentType").Params().Params(jen.ID("string")),
				jen.ID("Unmarshal").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Index().ID("byte"), jen.ID("v").Interface()).Params(jen.ID("error")),
				jen.ID("Encode").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("dest").Qual("io", "Writer"), jen.ID("v").Interface()).Params(jen.ID("error")),
				jen.ID("EncodeReader").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Interface()).Params(jen.Qual("io", "Reader"), jen.ID("error")),
			),
			jen.ID("clientEncoder").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("contentType").Op("*").ID("contentType"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").Op("*").ID("clientEncoder")).ID("Unmarshal").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Index().ID("byte"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("e").Dot("logger").Dot("WithValue").Call(
				jen.Lit("data_length"),
				jen.ID("len").Call(jen.ID("data")),
			),
			jen.Var().Defs(
				jen.ID("unmarshalFunc").Func().Params(jen.ID("data").Index().ID("byte"), jen.ID("v").Interface()).Params(jen.ID("error")),
			),
			jen.Switch(jen.ID("e").Dot("contentType")).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("unmarshalFunc").Op("=").Qual("encoding/xml", "Unmarshal")),
				jen.Default().Body(
					jen.ID("unmarshalFunc").Op("=").Qual("encoding/json", "Unmarshal")),
			),
			jen.If(jen.ID("err").Op(":=").ID("unmarshalFunc").Call(
				jen.ID("data"),
				jen.ID("v"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling JSON content"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").Op("*").ID("clientEncoder")).ID("Encode").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("dest").Qual("io", "Writer"), jen.ID("data").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("e").Dot("logger"),
			jen.Var().Defs(
				jen.ID("err").ID("error"),
			),
			jen.Switch(jen.ID("e").Dot("contentType")).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("err").Op("=").Qual("encoding/xml", "NewEncoder").Call(jen.ID("dest")).Dot("Encode").Call(jen.ID("data"))),
				jen.Default().Body(
					jen.ID("err").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.ID("dest")).Dot("Encode").Call(jen.ID("data"))),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding JSON content"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").Op("*").ID("clientEncoder")).ID("EncodeReader").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("data").Interface()).Params(jen.Qual("io", "Reader"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().Defs(
				jen.ID("marshalFunc").Func().Params(jen.ID("v").Interface()).Params(jen.Index().ID("byte"), jen.ID("error")),
			),
			jen.Switch(jen.ID("e").Dot("contentType")).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("marshalFunc").Op("=").Qual("encoding/xml", "Marshal")),
				jen.Default().Body(
					jen.ID("marshalFunc").Op("=").Qual("encoding/json", "Marshal")),
			),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").ID("marshalFunc").Call(jen.ID("data")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("e").Dot("logger"),
					jen.ID("span"),
					jen.Lit("marshaling to XML"),
				))),
			jen.Return().List(jen.Qual("bytes", "NewReader").Call(jen.ID("out")), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideClientEncoder provides a ClientEncoder."),
		jen.Line(),
		jen.Func().ID("ProvideClientEncoder").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("encoding").Op("*").ID("contentType")).Params(jen.ID("ClientEncoder")).Body(
			jen.Return().Op("&").ID("clientEncoder").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.Lit("client_encoder")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("client_encoder")), jen.ID("contentType").Op(":").ID("encoding"))),
		jen.Line(),
	)

	return code
}
