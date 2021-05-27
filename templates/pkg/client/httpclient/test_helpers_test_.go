package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func testHelpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("exampleURI").Op("=").Lit("https://todo.verygoodsoftwarenotvirus.ru"),
			jen.ID("asciiControlChar").Op("=").ID("string").Call(jen.ID("byte").Call(jen.Lit(127))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("mustParseURL parses a URL string or otherwise panics."),
		jen.Line(),
		jen.Func().ID("mustParseURL").Params(jen.ID("raw").ID("string")).Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "ParseRequestURI").Call(jen.ID("raw")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("u"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("argleBargle").Struct(jen.ID("Name").ID("string")),
			jen.ID("valuer").Map(jen.ID("string")).Index().ID("string"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("v").ID("valuer")).ID("ToValues").Params().Params(jen.Qual("net/url", "Values")).Body(
			jen.Return().Qual("net/url", "Values").Call(jen.ID("v"))),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("requestSpec").Struct(
				jen.ID("path").ID("string"),
				jen.ID("method").ID("string"),
				jen.ID("query").ID("string"),
				jen.ID("pathArgs").Index().Interface(),
				jen.ID("bodyShouldBeEmpty").ID("bool"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("newRequestSpec").Params(jen.ID("bodyShouldBeEmpty").ID("bool"), jen.List(jen.ID("method"), jen.ID("query"), jen.ID("path")).ID("string"), jen.ID("pathArgs").Op("...").Interface()).Params(jen.Op("*").ID("requestSpec")).Body(
			jen.Return().Op("&").ID("requestSpec").Valuesln(jen.ID("path").Op(":").ID("path"), jen.ID("pathArgs").Op(":").ID("pathArgs"), jen.ID("method").Op(":").ID("method"), jen.ID("query").Op(":").ID("query"), jen.ID("bodyShouldBeEmpty").Op(":").ID("bodyShouldBeEmpty"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("assertErrorMatches").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.List(jen.ID("err1"), jen.ID("err2")).ID("error")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("assert").Dot("True").Call(
				jen.ID("t"),
				jen.Qual("errors", "Is").Call(
					jen.ID("err1"),
					jen.ID("err2"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("assertRequestQuality").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("spec").Op("*").ID("requestSpec")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("expectedPath").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.ID("spec").Dot("path"),
				jen.ID("spec").Dot("pathArgs").Op("..."),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("req"),
				jen.Lit("provided req must not be nil"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("spec"),
				jen.Lit("provided spec must not be nil"),
			),
			jen.List(jen.ID("bodyBytes"), jen.ID("err")).Op(":=").ID("httputil").Dot("DumpRequest").Call(
				jen.ID("req"),
				jen.ID("true"),
			),
			jen.ID("assert").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("bodyBytes"),
			),
			jen.ID("assert").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.If(jen.ID("spec").Dot("bodyShouldBeEmpty")).Body(
				jen.ID("bodyLines").Op(":=").Qual("strings", "Split").Call(
					jen.ID("string").Call(jen.ID("bodyBytes")),
					jen.Lit("\n"),
				),
				jen.ID("assert").Dot("Empty").Call(
					jen.ID("t"),
					jen.ID("bodyLines").Index(jen.ID("len").Call(jen.ID("bodyLines")).Op("-").Lit(1)),
				),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("spec").Dot("query"),
				jen.ID("req").Dot("URL").Dot("Query").Call().Dot("Encode").Call(),
				jen.Lit("expected query to be %q, but was %q instead"),
				jen.ID("spec").Dot("query"),
				jen.ID("req").Dot("URL").Dot("Query").Call().Dot("Encode").Call(),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("expectedPath"),
				jen.ID("req").Dot("URL").Dot("Path"),
				jen.Lit("expected path to be %q, but was %q instead"),
				jen.ID("expectedPath"),
				jen.ID("req").Dot("URL").Dot("Path"),
			),
			jen.ID("assert").Dot("Equal").Call(
				jen.ID("t"),
				jen.ID("spec").Dot("method"),
				jen.ID("req").Dot("Method"),
				jen.Lit("expected method to be %q, but was %q instead"),
				jen.ID("spec").Dot("method"),
				jen.ID("req").Dot("Method"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("createBodyFromStruct takes any value in and returns an io.Reader for placement within http.NewRequest's last argument."),
		jen.Line(),
		jen.Func().ID("createBodyFromStruct").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("in").Interface()).Params(jen.Qual("io", "Reader")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.List(jen.ID("out"), jen.ID("err")).Op(":=").Qual("encoding/json", "Marshal").Call(jen.ID("in")),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().Qual("bytes", "NewReader").Call(jen.ID("out")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClient").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("ts").Op("*").ID("httptest").Dot("Server")).Params(jen.Op("*").ID("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("ts"),
			),
			jen.List(jen.ID("client"), jen.ID("err")).Op(":=").ID("NewClient").Call(
				jen.ID("mustParseURL").Call(jen.Lit("https://todo.verygoodsoftwarenotvirus.ru")),
				jen.ID("UsingLogger").Call(jen.ID("logging").Dot("NewNoopLogger").Call()),
				jen.ID("UsingJSON").Call(),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("client"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("client").Dot("requestBuilder").Dot("SetURL").Call(jen.ID("mustParseURL").Call(jen.ID("ts").Dot("URL"))),
			),
			jen.ID("client").Dot("unauthenticatedClient").Op("=").ID("ts").Dot("Client").Call(),
			jen.ID("client").Dot("authedClient").Op("=").ID("ts").Dot("Client").Call(),
			jen.Return().ID("client"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildSimpleTestClient").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Client"), jen.Op("*").ID("httptest").Dot("Server")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.ID("nil")),
			jen.Return().List(jen.ID("buildTestClient").Call(
				jen.ID("t"),
				jen.ID("ts"),
			), jen.ID("ts")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClientWithInvalidURL").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("l").Op(":=").ID("logging").Dot("NewNoopLogger").Call(),
			jen.ID("u").Op(":=").ID("mustParseURL").Call(jen.Lit("https://todo.verygoodsoftwarenotvirus.ru")),
			jen.ID("u").Dot("Scheme").Op("=").Qual("fmt", "Sprintf").Call(
				jen.Lit(`%s://`),
				jen.ID("asciiControlChar"),
			),
			jen.List(jen.ID("c"), jen.ID("err")).Op(":=").ID("NewClient").Call(
				jen.ID("u"),
				jen.ID("UsingLogger").Call(jen.ID("l")),
				jen.ID("UsingDebug").Call(jen.ID("true")),
			),
			jen.ID("require").Dot("NotNil").Call(
				jen.ID("t"),
				jen.ID("c"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.Return().ID("c"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClientWithStatusCodeResponse").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("spec").Op("*").ID("requestSpec"), jen.ID("code").ID("int")).Params(jen.Op("*").ID("Client"), jen.Op("*").ID("httptest").Dot("Server")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.ID("t").Dot("Helper").Call(),
				jen.ID("assertRequestQuality").Call(
					jen.ID("t"),
					jen.ID("req"),
					jen.ID("spec"),
				),
				jen.ID("res").Dot("WriteHeader").Call(jen.ID("code")),
			))),
			jen.Return().List(jen.ID("buildTestClient").Call(
				jen.ID("t"),
				jen.ID("ts"),
			), jen.ID("ts")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClientWithInvalidResponse").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("spec").Op("*").ID("requestSpec")).Params(jen.Op("*").ID("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.ID("t").Dot("Helper").Call(),
				jen.ID("assertRequestQuality").Call(
					jen.ID("t"),
					jen.ID("req"),
					jen.ID("spec"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.Lit("BLAH")),
				),
			))),
			jen.Return().ID("buildTestClient").Call(
				jen.ID("t"),
				jen.ID("ts"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClientWithJSONResponse").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("spec").Op("*").ID("requestSpec"), jen.ID("outputBody").Interface()).Params(jen.Op("*").ID("Client"), jen.Op("*").ID("httptest").Dot("Server")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.ID("t").Dot("Helper").Call(),
				jen.ID("assertRequestQuality").Call(
					jen.ID("t"),
					jen.ID("req"),
					jen.ID("spec"),
				),
				jen.ID("assert").Dot("NoError").Call(
					jen.ID("t"),
					jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("outputBody")),
				),
			))),
			jen.Return().List(jen.ID("buildTestClient").Call(
				jen.ID("t"),
				jen.ID("ts"),
			), jen.ID("ts")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClientWithRequestBodyValidation").Params(jen.ID("t").Op("*").Qual("testing", "T"), jen.ID("spec").Op("*").ID("requestSpec"), jen.List(jen.ID("inputBody"), jen.ID("expectedInput"), jen.ID("outputBody")).Interface()).Params(jen.Op("*").ID("Client")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.ID("t").Dot("Helper").Call(),
				jen.ID("assertRequestQuality").Call(
					jen.ID("t"),
					jen.ID("req"),
					jen.ID("spec"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.Qual("encoding/json", "NewDecoder").Call(jen.ID("req").Dot("Body")).Dot("Decode").Call(jen.Op("&").ID("inputBody")),
				),
				jen.ID("assert").Dot("Equal").Call(
					jen.ID("t"),
					jen.ID("expectedInput"),
					jen.ID("inputBody"),
				),
				jen.ID("require").Dot("NoError").Call(
					jen.ID("t"),
					jen.Qual("encoding/json", "NewEncoder").Call(jen.ID("res")).Dot("Encode").Call(jen.ID("outputBody")),
				),
			))),
			jen.Return().ID("buildTestClient").Call(
				jen.ID("t"),
				jen.ID("ts"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestClientThatWaitsTooLong").Params(jen.ID("t").Op("*").Qual("testing", "T")).Params(jen.Op("*").ID("Client"), jen.Op("*").ID("httptest").Dot("Server")).Body(
			jen.ID("t").Dot("Helper").Call(),
			jen.ID("ts").Op(":=").ID("httptest").Dot("NewTLSServer").Call(jen.Qual("net/http", "HandlerFunc").Call(jen.Func().Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
				jen.Qual("time", "Sleep").Call(jen.Lit(24).Op("*").Qual("time", "Hour"))))),
			jen.ID("c").Op(":=").ID("buildTestClient").Call(
				jen.ID("t"),
				jen.ID("ts"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("c").Dot("SetOptions").Call(jen.ID("UsingTimeout").Call(jen.Qual("time", "Millisecond"))),
			),
			jen.Return().List(jen.ID("c"), jen.ID("ts")),
		),
		jen.Line(),
	)

	return code
}
