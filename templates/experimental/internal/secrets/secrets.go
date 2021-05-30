package secrets

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func secretsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("tracerName").Op("=").Lit("secret_manager"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidKeeper").Op("=").Qual("errors", "New").Call(jen.Lit("invalid keeper")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("SecretManager").Interface(
				jen.ID("Encrypt").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")),
				jen.ID("Decrypt").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("content").ID("string"), jen.ID("v").Interface()).Params(jen.ID("error")),
			),
			jen.ID("secretManager").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("keeper").Op("*").Qual("gocloud.dev/secrets", "Keeper"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ProvideSecretManager builds a new SecretManager."),
		jen.Line(),
		jen.Func().ID("ProvideSecretManager").Params(jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("keeper").Op("*").Qual("gocloud.dev/secrets", "Keeper")).Params(jen.ID("SecretManager"), jen.ID("error")).Body(
			jen.If(jen.ID("keeper").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("errInvalidKeeper"))),
			jen.ID("sm").Op(":=").Op("&").ID("secretManager").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("tracerName")), jen.ID("keeper").Op(":").ID("keeper")),
			jen.Return().List(jen.ID("sm"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Encrypt does the following:"),
		jen.Line(),
		jen.Func().Comment("//\t\t1. JSON encodes a given value").Comment("//\t\t2. encrypts that encoded data").Comment("//\t\t3. base64 URL encodes that encrypted data").Params(jen.ID("sm").Op("*").ID("secretManager")).ID("Encrypt").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("value").Interface()).Params(jen.ID("string"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("jsonBytes"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("value")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encoding value to JSON: %w"),
					jen.ID("err"),
				))),
			jen.List(jen.ID("encrypted"), jen.ID("err")).Op(":=").ID("sm").Dot("keeper").Dot("Encrypt").Call(
				jen.ID("ctx"),
				jen.ID("jsonBytes"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(""), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("encrypting JSON encoded bytes: %w"),
					jen.ID("err"),
				))),
			jen.ID("encoded").Op(":=").Qual("encoding/base64", "URLEncoding").Dot("EncodeToString").Call(jen.ID("encrypted")),
			jen.Return().List(jen.ID("encoded"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Decrypt does the following:"),
		jen.Line(),
		jen.Func().Comment("//\t\t1. base64 URL decodes the provided data").Comment("//\t\t2. decrypts that encoded data").Comment("//\t\t3. JSON decodes that decrypted data into the target variable.").Params(jen.ID("sm").Op("*").ID("secretManager")).ID("Decrypt").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("content").ID("string"), jen.ID("v").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.List(jen.ID("decoded"), jen.ID("err")).Op(":=").Qual("encoding/base64", "URLEncoding").Dot("DecodeString").Call(jen.ID("content")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("decoding base64 encoded content: %w"),
					jen.ID("err"),
				)),
			jen.List(jen.ID("jsonBytes"), jen.ID("err")).Op(":=").ID("sm").Dot("keeper").Dot("Decrypt").Call(
				jen.ID("ctx"),
				jen.ID("decoded"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("decrypting decoded bytes into JSON: %w"),
					jen.ID("err"),
				)),
			jen.ID("err").Op("=").Qual("encoding/json", "Unmarshal").Call(
				jen.ID("jsonBytes"),
				jen.Op("&").ID("v"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual("fmt", "Errorf").Call(
					jen.Lit("decoded JSON bytes into value: %w"),
					jen.ID("err"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
