package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverEncoderDecoderDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("ContentTypeHeaderKey").Op("=").Lit("RawHTML-type").Var().ID("contentTypeXML").Op("=").Lit("application/xml").Var().ID("contentTypeJSON").Op("=").Lit("application/json"),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("defaultContentType").Op("=").ID("ContentTypeJSON"),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("ServerEncoderDecoder").Interface(
			jen.ID("RespondWithData").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("val").Interface()),
			jen.ID("EncodeResponseWithStatus").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("val").Interface(), jen.ID("statusCode").ID("int")),
			jen.ID("EncodeErrorResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("msg").ID("string"), jen.ID("statusCode").ID("int")),
			jen.ID("EncodeInvalidInputResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")),
			jen.ID("EncodeNotFoundResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")),
			jen.ID("EncodeUnspecifiedInternalServerErrorResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")),
			jen.ID("EncodeUnauthorizedResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")),
			jen.ID("EncodeInvalidPermissionsResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")),
			jen.ID("DecodeRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("dest").Interface()).Params(jen.ID("error")),
			jen.ID("MustEncode").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("v").Interface()).Params(jen.Index().ID("byte")),
			jen.ID("MustEncodeJSON").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("v").Interface()).Params(jen.Index().ID("byte")),
		).Type().ID("serverEncoderDecoder").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
			jen.ID("panicker").ID("panicking").Dot("Panicker"),
			jen.ID("contentType").ID("ContentType"),
		).Type().ID("encoder").Interface(jen.ID("Encode").Params(jen.Interface()).Params(jen.ID("error"))).Type().ID("decoder").Interface(jen.ID("Decode").Params(jen.ID("v").Interface()).Params(jen.ID("error"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("encodeResponse encodes responses."),
		jen.Line(),
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("encodeResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface(), jen.ID("statusCode").ID("int")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("e").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ResponseStatusKey"),
				jen.ID("statusCode"),
			),
			jen.Var().ID("enc").ID("encoder"),
			jen.Switch(jen.ID("contentTypeFromString").Call(jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("ContentTypeHeaderKey")))).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.ID("contentTypeXML"),
					), jen.ID("enc").Op("=").Qual("encoding/xml", "NewEncoder").Call(jen.ID("res"))),
				jen.Case(jen.ID("ContentTypeJSON")).Body(
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.ID("contentTypeJSON"),
					), jen.Fallthrough()),
				jen.Default().Body(
					jen.ID("enc").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.ID("res"))),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.ID("statusCode")),
			jen.If(jen.ID("err").Op(":=").ID("enc").Dot("Encode").Call(jen.ID("v")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding response"),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeErrorResponse encodes errs to responses."),
		jen.Line(),
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeErrorResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("msg").ID("string"), jen.ID("statusCode").ID("int")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().ID("enc").ID("encoder").Var().ID("logger").Op("=").ID("e").Dot("logger").Dot("WithValue").Call(
				jen.Lit("error_message"),
				jen.ID("msg"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("ResponseStatusKey"),
				jen.ID("statusCode"),
			),
			jen.Switch(jen.ID("contentTypeFromString").Call(jen.ID("res").Dot("Header").Call().Dot("Get").Call(jen.ID("ContentTypeHeaderKey")))).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.ID("contentTypeXML"),
					), jen.ID("enc").Op("=").Qual("encoding/xml", "NewEncoder").Call(jen.ID("res"))),
				jen.Case(jen.ID("ContentTypeJSON")).Body(
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.ID("contentTypeJSON"),
					), jen.Fallthrough()),
				jen.Default().Body(
					jen.ID("enc").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.ID("res"))),
			),
			jen.ID("res").Dot("WriteHeader").Call(jen.ID("statusCode")),
			jen.If(jen.ID("err").Op(":=").ID("enc").Dot("Encode").Call(jen.Op("&").ID("types").Dot("ErrorResponse").Valuesln(jen.ID("Message").Op(":").ID("msg"), jen.ID("Code").Op(":").ID("statusCode"))), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("encoding error response"),
				)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("EncodeInvalidInputResponse encodes a generic 400 error to a response.").Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeInvalidInputResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("invalid input attached to request"),
				jen.Qual("net/http", "StatusBadRequest"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("EncodeNotFoundResponse encodes a generic 404 error to a response.").Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeNotFoundResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("resource not found"),
				jen.Qual("net/http", "StatusNotFound"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("EncodeUnspecifiedInternalServerErrorResponse encodes a generic 500 error to a response.").Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeUnspecifiedInternalServerErrorResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("something has gone awry"),
				jen.Qual("net/http", "StatusInternalServerError"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("EncodeUnauthorizedResponse encodes a generic 401 error to a response.").Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeUnauthorizedResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("invalid credentials provided"),
				jen.Qual("net/http", "StatusUnauthorized"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("EncodeInvalidPermissionsResponse encodes a generic 403 error to a response.").Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeInvalidPermissionsResponse").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("EncodeErrorResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.Lit("invalid permissions"),
				jen.Qual("net/http", "StatusForbidden"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("MustEncodeJSON").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("v").Interface()).Params(jen.Index().ID("byte")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().ID("b").Qual("bytes", "Buffer"),
			jen.If(jen.ID("err").Op(":=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b")).Dot("Encode").Call(jen.ID("v")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("e").Dot("panicker").Dot("Panicf").Call(
					jen.Lit("encoding JSON content: %w"),
					jen.ID("err"),
				)),
			jen.Return().ID("b").Dot("Bytes").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MustEncode encodes data or else."),
		jen.Line(),
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("MustEncode").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("v").Interface()).Params(jen.Index().ID("byte")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().ID("enc").ID("encoder").Var().ID("b").Qual("bytes", "Buffer"),
			jen.Switch(jen.ID("e").Dot("contentType")).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("enc").Op("=").Qual("encoding/xml", "NewEncoder").Call(jen.Op("&").ID("b"))),
				jen.Default().Body(
					jen.ID("enc").Op("=").Qual("encoding/json", "NewEncoder").Call(jen.Op("&").ID("b"))),
			),
			jen.If(jen.ID("err").Op(":=").ID("enc").Dot("Encode").Call(jen.ID("v")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("e").Dot("panicker").Dot("Panicf").Call(
					jen.Lit("encoding %s content: %w"),
					jen.ID("e").Dot("contentType"),
					jen.ID("err"),
				)),
			jen.Return().ID("b").Dot("Bytes").Call(),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RespondWithData encodes successful responses with data."),
		jen.Line(),
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("RespondWithData").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface()).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("encodeResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("v"),
				jen.Qual("net/http", "StatusOK"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("EncodeResponseWithStatus encodes responses and writes the provided status to the response."),
		jen.Line(),
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("EncodeResponseWithStatus").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("v").Interface(), jen.ID("statusCode").ID("int")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("e").Dot("encodeResponse").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("v"),
				jen.ID("statusCode"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("DecodeRequest decodes request bodies into values."),
		jen.Line(),
		jen.Func().Params(jen.ID("e").Op("*").ID("serverEncoderDecoder")).ID("DecodeRequest").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("e").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().ID("d").ID("decoder"),
			jen.Switch(jen.ID("contentTypeFromString").Call(jen.ID("req").Dot("Header").Dot("Get").Call(jen.ID("ContentTypeHeaderKey")))).Body(
				jen.Case(jen.ID("ContentTypeXML")).Body(
					jen.ID("d").Op("=").Qual("encoding/xml", "NewDecoder").Call(jen.ID("req").Dot("Body"))),
				jen.Default().Body(
					jen.ID("d").Op("=").Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body"))),
			),
			jen.Return().ID("d").Dot("Decode").Call(jen.ID("v")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideServerEncoderDecoder provides a ServerEncoderDecoder."),
		jen.Line(),
		jen.Func().ID("ProvideServerEncoderDecoder").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("contentType").ID("ContentType")).Params(jen.ID("ServerEncoderDecoder")).Body(
			jen.Return().Op("&").ID("serverEncoderDecoder").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.Lit("server_encoder_decoder")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("server_encoder_decoder")), jen.ID("panicker").Op(":").ID("panicking").Dot("NewProductionPanicker").Call(), jen.ID("contentType").Op(":").ID("contentType"))),
		jen.Line(),
	)

	return code
}
