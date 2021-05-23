package random

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func randDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("arbitrarySize").ID("uint16").Op("=").Lit(128),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("_").ID("Generator").Op("=").Parens(jen.Op("*").ID("standardGenerator")).Call(jen.ID("nil")),
		jen.ID("defaultGenerator").Op("=").ID("NewGenerator").Call(jen.ID("logging").Dot("NewNonOperationalLogger").Call()),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").Qual("crypto/rand", "Read").Call(jen.ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("arbitrarySize"),
			)), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err")))),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("Generator").Interface(
			jen.ID("GenerateBase32EncodedString").Params(jen.Qual("context", "Context"), jen.ID("int")).Params(jen.ID("string"), jen.ID("error")),
			jen.ID("GenerateBase64EncodedString").Params(jen.Qual("context", "Context"), jen.ID("int")).Params(jen.ID("string"), jen.ID("error")),
			jen.ID("GenerateRawBytes").Params(jen.Qual("context", "Context"), jen.ID("int")).Params(jen.Index().ID("byte"), jen.ID("error")),
		).Type().ID("standardGenerator").Struct(
			jen.ID("logger").ID("logging").Dot("Logger"),
			jen.ID("tracer").ID("tracing").Dot("Tracer"),
			jen.ID("randReader").Qual("io", "Reader"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewGenerator builds a new Generator."),
		jen.Line(),
		jen.Func().ID("NewGenerator").Params(jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("Generator")).Body(
			jen.Return().Op("&").ID("standardGenerator").Valuesln(
				jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.Lit("random_string_generator")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("secret_generator")), jen.ID("randReader").Op(":").Qual("crypto/rand", "Reader"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateBase32EncodedString generates a base64-encoded string of a securely random byte array of a given length."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").Op("*").ID("standardGenerator")).ID("GenerateBase32EncodedString").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("g").Dot("logger").Dot("WithValue").Call(
				jen.Lit("requested_length"),
				jen.ID("length"),
			),
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("length"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("g").Dot("randReader").Dot("Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading from secure random source"),
				))),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateBase64EncodedString generates a base64-encoded string of a securely random byte array of a given length."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").Op("*").ID("standardGenerator")).ID("GenerateBase64EncodedString").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("g").Dot("logger").Dot("WithValue").Call(
				jen.Lit("requested_length"),
				jen.ID("length"),
			),
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("length"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("g").Dot("randReader").Dot("Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading from secure random source"),
				))),
			jen.Return().List(jen.Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateRawBytes generates a securely random byte array."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").Op("*").ID("standardGenerator")).ID("GenerateRawBytes").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.Index().ID("byte"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("g").Dot("logger").Dot("WithValue").Call(
				jen.Lit("requested_length"),
				jen.ID("length"),
			),
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("length"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("g").Dot("randReader").Dot("Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading from secure random source"),
				))),
			jen.Return().List(jen.ID("b"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateBase32EncodedString generates a base64-encoded string of a securely random byte array of a given length."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").Op("*").ID("standardGenerator")).ID("GenerateBase32EncodedString").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("g").Dot("logger").Dot("WithValue").Call(
				jen.Lit("requested_length"),
				jen.ID("length"),
			),
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("length"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("g").Dot("randReader").Dot("Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading from secure random source"),
				))),
			jen.Return().List(jen.Qual("encoding/base32", "StdEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateBase64EncodedString generates a base64-encoded string of a securely random byte array of a given length."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").Op("*").ID("standardGenerator")).ID("GenerateBase64EncodedString").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("g").Dot("logger").Dot("WithValue").Call(
				jen.Lit("requested_length"),
				jen.ID("length"),
			),
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("length"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("g").Dot("randReader").Dot("Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading from secure random source"),
				))),
			jen.Return().List(jen.Qual("encoding/base64", "RawURLEncoding").Dot("EncodeToString").Call(jen.ID("b")), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GenerateRawBytes generates a securely random byte array."),
		jen.Line(),
		jen.Func().Params(jen.ID("g").Op("*").ID("standardGenerator")).ID("GenerateRawBytes").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("length").ID("int")).Params(jen.Index().ID("byte"), jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("tracing").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("g").Dot("logger").Dot("WithValue").Call(
				jen.Lit("requested_length"),
				jen.ID("length"),
			),
			jen.ID("b").Op(":=").ID("make").Call(
				jen.Index().ID("byte"),
				jen.ID("length"),
			),
			jen.If(jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID("g").Dot("randReader").Dot("Read").Call(jen.ID("b")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("reading from secure random source"),
				))),
			jen.Return().List(jen.ID("b"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
