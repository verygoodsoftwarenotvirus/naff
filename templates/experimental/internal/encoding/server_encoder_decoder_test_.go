package encoding

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serverEncoderDecoderTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Null().Type().ID("example").Struct(jen.ID("Name").ID("string")),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("broken").Struct(jen.ID("Name").Qual("encoding/json", "Number")),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_encodeResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("ex").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.List(jen.ID("encoderDecoder"), jen.ID("ok")).Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					).Assert(jen.Op("*").ID("serverEncoderDecoder")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Dot("encodeResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("ex"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Body").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("{%q:%q}\n"),
							jen.Lit("name"),
							jen.ID("ex").Dot("Name"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("as XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("ex").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.List(jen.ID("encoderDecoder"), jen.ID("ok")).Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					).Assert(jen.Op("*").ID("serverEncoderDecoder")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.Lit("application/xml"),
					),
					jen.ID("encoderDecoder").Dot("encodeResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("ex"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("<example><name>%s</name></example>"),
							jen.ID("expectation"),
						),
						jen.ID("res").Dot("Body").Dot("String").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with broken structure"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("ex").Op(":=").Op("&").ID("broken").Valuesln(jen.ID("Name").Op(":").Qual("encoding/json", "Number").Call(jen.ID("expectation"))),
					jen.List(jen.ID("encoderDecoder"), jen.ID("ok")).Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					).Assert(jen.Op("*").ID("serverEncoderDecoder")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Dot("encodeResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("ex"),
						jen.Qual("net/http", "StatusOK"),
					),
					jen.ID("assert").Dot("Empty").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Body").Dot("String").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeErrorResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleMessage").Op(":=").Lit("something went awry"),
					jen.ID("exampleCode").Op(":=").Qual("net/http", "StatusBadRequest"),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Dot("EncodeErrorResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("exampleMessage"),
						jen.ID("exampleCode"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Body").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("{\"message\":%q,\"code\":%d}\n"),
							jen.ID("exampleMessage"),
							jen.ID("exampleCode"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected status code to match"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("as XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleMessage").Op(":=").Lit("something went awry"),
					jen.ID("exampleCode").Op(":=").Qual("net/http", "StatusBadRequest"),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.Lit("application/xml"),
					),
					jen.ID("encoderDecoder").Dot("EncodeErrorResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("exampleMessage"),
						jen.ID("exampleCode"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("<ErrorResponse><Message>%s</Message><Code>%d</Code></ErrorResponse>"),
							jen.ID("exampleMessage"),
							jen.ID("exampleCode"),
						),
						jen.ID("res").Dot("Body").Dot("String").Call(),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("exampleCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected status code to match"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeInvalidInputResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("encoderDecoder").Dot("EncodeInvalidInputResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.ID("expectedCode").Op(":=").Qual("net/http", "StatusBadRequest"),
					jen.ID("assert").Dot("EqualValues").Call(
						jen.ID("t"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected code to be %d, got %d instead"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeNotFoundResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.ID("expectedCode").Op(":=").Qual("net/http", "StatusNotFound"),
					jen.ID("assert").Dot("EqualValues").Call(
						jen.ID("t"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected code to be %d, got %d instead"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeUnspecifiedInternalServerErrorResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.ID("expectedCode").Op(":=").Qual("net/http", "StatusInternalServerError"),
					jen.ID("assert").Dot("EqualValues").Call(
						jen.ID("t"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected code to be %d, got %d instead"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeUnauthorizedResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("encoderDecoder").Dot("EncodeUnauthorizedResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.ID("expectedCode").Op(":=").Qual("net/http", "StatusUnauthorized"),
					jen.ID("assert").Dot("EqualValues").Call(
						jen.ID("t"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected code to be %d, got %d instead"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeInvalidPermissionsResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("encoderDecoder").Dot("EncodeInvalidPermissionsResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
					jen.ID("expectedCode").Op(":=").Qual("net/http", "StatusForbidden"),
					jen.ID("assert").Dot("EqualValues").Call(
						jen.ID("t"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected code to be %d, got %d instead"),
						jen.ID("expectedCode"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_MustEncodeJSON").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("expected").Op(":=").Lit(`{"name":"TestServerEncoderDecoder_MustEncodeJSON/standard"}
`),
					jen.ID("actual").Op(":=").ID("string").Call(jen.ID("encoderDecoder").Dot("MustEncodeJSON").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					)),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with panic"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.Defer().Func().Params().Body(
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("recover").Call(),
						)).Call(),
					jen.ID("encoderDecoder").Dot("MustEncodeJSON").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("broken").Valuesln(jen.ID("Name").Op(":").Qual("encoding/json", "Number").Call(jen.ID("t").Dot("Name").Call())),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_MustEncode").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with JSON"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("expected").Op(":=").Lit(`{"name":"TestServerEncoderDecoder_MustEncode/with_JSON"}
`),
					jen.ID("actual").Op(":=").ID("string").Call(jen.ID("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					)),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeXML"),
					),
					jen.ID("expected").Op(":=").Lit(`<example><name>TestServerEncoderDecoder_MustEncode/with_XML</name></example>`),
					jen.ID("actual").Op(":=").ID("string").Call(jen.ID("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("t").Dot("Name").Call()),
					)),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with broken struct"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.List(jen.ID("encoderDecoder"), jen.ID("ok")).Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					).Assert(jen.Op("*").ID("serverEncoderDecoder")),
					jen.ID("require").Dot("True").Call(
						jen.ID("t"),
						jen.ID("ok"),
					),
					jen.Defer().Func().Params().Body(
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("recover").Call(),
						)).Call(),
					jen.ID("encoderDecoder").Dot("MustEncode").Call(
						jen.ID("ctx"),
						jen.Op("&").ID("broken").Valuesln(jen.ID("Name").Op(":").Qual("encoding/json", "Number").Call(jen.ID("t").Dot("Name").Call())),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_RespondWithData").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("ex").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("encoderDecoder").Dot("RespondWithData").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("ex"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Body").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("{%q:%q}\n"),
							jen.Lit("name"),
							jen.ID("ex").Dot("Name"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("as XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("ex").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("res").Dot("Header").Call().Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.Lit("application/xml"),
					),
					jen.ID("encoderDecoder").Dot("RespondWithData").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("ex"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("<example><name>%s</name></example>"),
							jen.ID("expectation"),
						),
						jen.ID("res").Dot("Body").Dot("String").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_EncodeResponseWithStatus").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("ex").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("expected").Op(":=").Lit(666),
					jen.ID("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
						jen.ID("ctx"),
						jen.ID("res"),
						jen.ID("ex"),
						jen.ID("expected"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("res").Dot("Code"),
						jen.Lit("expected code to be %d, but got %d"),
						jen.ID("expected"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("res").Dot("Body").Dot("String").Call(),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("{%q:%q}\n"),
							jen.Lit("name"),
							jen.ID("ex").Dot("Name"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestServerEncoderDecoder_DecodeRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("e").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("e")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("bs")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.ID("contentTypeJSON"),
					),
					jen.Var().ID("x").ID("example"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder").Dot("DecodeRequest").Call(
							jen.ID("ctx"),
							jen.ID("req"),
							jen.Op("&").ID("x"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Name"),
						jen.ID("e").Dot("Name"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("as XML"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("expectation").Op(":=").Lit("name"),
					jen.ID("e").Op(":=").Op("&").ID("example").Valuesln(jen.ID("Name").Op(":").ID("expectation")),
					jen.ID("encoderDecoder").Op(":=").ID("ProvideServerEncoderDecoder").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("ContentTypeJSON"),
					),
					jen.List(jen.ID("bs"), jen.ID("err")).Op(":=").Qual("encoding/xml", "Marshal").Call(jen.ID("e")),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("https://todo.verygoodsoftwarenotvirus.ru"),
						jen.Qual("bytes", "NewReader").Call(jen.ID("bs")),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.ID("ContentTypeHeaderKey"),
						jen.ID("contentTypeXML"),
					),
					jen.Var().ID("x").ID("example"),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("encoderDecoder").Dot("DecodeRequest").Call(
							jen.ID("ctx"),
							jen.ID("req"),
							jen.Op("&").ID("x"),
						),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Name"),
						jen.ID("e").Dot("Name"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
