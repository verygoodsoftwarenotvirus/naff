package requests

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func testHelpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("exampleURI").Op("=").Lit("https://todo.verygoodsoftwarenotvirus.ru").Var().ID("asciiControlChar").Op("=").ID("string").Call(jen.ID("byte").Call(jen.Lit(127))),
		jen.Line(),
	)

	code.Add(
		jen.Var().ID("parsedExampleURL").Op("*").Qual("net/url", "URL").Var().ID("invalidParsedURL").Op("*").Qual("net/url", "URL"),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("parsedExampleURL"), jen.ID("err")).Op("=").Qual("net/url", "Parse").Call(jen.ID("exampleURI")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.ID("u").Op(":=").ID("mustParseURL").Call(jen.Lit("https://verygoodsoftwarenotvirus.ru")),
			jen.ID("u").Dot("Scheme").Op("=").Qual("fmt", "Sprintf").Call(
				jen.Lit(`%s://`),
				jen.ID("asciiControlChar"),
			),
			jen.ID("invalidParsedURL").Op("=").ID("u"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("mustParseURL parses a url or otherwise panics."),
		jen.Line(),
		jen.Func().ID("mustParseURL").Params(jen.ID("raw").ID("string")).Params(jen.Op("*").Qual("net/url", "URL")).Body(
			jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "Parse").Call(jen.ID("raw")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("u"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("valuer").Map(jen.ID("string")).Index().ID("string"),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("v").ID("valuer")).ID("ToValues").Params().Params(jen.Qual("net/url", "Values")).Body(
			jen.Return().Qual("net/url", "Values").Call(jen.ID("v"))),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("requestSpec").Struct(
			jen.ID("path").ID("string"),
			jen.ID("method").ID("string"),
			jen.ID("query").ID("string"),
			jen.ID("pathArgs").Index().Interface(),
			jen.ID("bodyShouldBeEmpty").ID("bool"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("newRequestSpec").Params(jen.ID("bodyShouldBeEmpty").ID("bool"), jen.List(jen.ID("method"), jen.ID("query"), jen.ID("path")).ID("string"), jen.ID("pathArgs").Op("...").Interface()).Params(jen.Op("*").ID("requestSpec")).Body(
			jen.Return().Op("&").ID("requestSpec").Valuesln(jen.ID("path").Op(":").ID("path"), jen.ID("pathArgs").Op(":").ID("pathArgs"), jen.ID("method").Op(":").ID("method"), jen.ID("query").Op(":").ID("query"), jen.ID("bodyShouldBeEmpty").Op(":").ID("bodyShouldBeEmpty"))),
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
			jen.ID("require").Dot("NotEmpty").Call(
				jen.ID("t"),
				jen.ID("bodyBytes"),
			),
			jen.ID("require").Dot("NoError").Call(
				jen.ID("t"),
				jen.ID("err"),
			),
			jen.If(jen.ID("spec").Dot("bodyShouldBeEmpty")).Body(
				jen.ID("bodyLines").Op(":=").Qual("strings", "Split").Call(
					jen.ID("string").Call(jen.ID("bodyBytes")),
					jen.Lit("\n"),
				),
				jen.ID("require").Dot("NotEmpty").Call(
					jen.ID("t"),
					jen.ID("bodyLines"),
				),
				jen.ID("assert").Dot("Empty").Call(
					jen.ID("t"),
					jen.ID("bodyLines").Index(jen.ID("len").Call(jen.ID("bodyLines")).Op("-").Lit(1)),
					jen.Lit("body was expected to be empty, and was not empty"),
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
		jen.Func().ID("buildTestRequestBuilder").Params().Params(jen.Op("*").ID("Builder")).Body(
			jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
			jen.Return().Op("&").ID("Builder").Valuesln(jen.ID("url").Op(":").ID("parsedExampleURL"), jen.ID("logger").Op(":").ID("l"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("test")), jen.ID("encoder").Op(":").ID("encoding").Dot("ProvideClientEncoder").Call(
				jen.ID("l"),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			)),
		),
		jen.Line(),
	)

	code.Add(
		jen.Null().Type().ID("testHelper").Struct(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("builder").Op("*").ID("Builder"),
			jen.ID("exampleUser").Op("*").ID("types").Dot("User"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestHelper").Params().Params(jen.Op("*").ID("testHelper")).Body(
			jen.ID("helper").Op(":=").Op("&").ID("testHelper").Valuesln(jen.ID("ctx").Op(":").Qual("context", "Background").Call(), jen.ID("builder").Op(":").ID("buildTestRequestBuilder").Call(), jen.ID("exampleUser").Op(":").ID("fakes").Dot("BuildFakeUser").Call()),
			jen.ID("helper").Dot("exampleUser").Dot("HashedPassword").Op("=").Lit(""),
			jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecret").Op("=").Lit(""),
			jen.ID("helper").Dot("exampleUser").Dot("TwoFactorSecretVerifiedOn").Op("=").ID("nil"),
			jen.Return().ID("helper"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("buildTestRequestBuilderWithInvalidURL").Params().Params(jen.Op("*").ID("Builder")).Body(
			jen.ID("l").Op(":=").ID("logging").Dot("NewNonOperationalLogger").Call(),
			jen.Return().Op("&").ID("Builder").Valuesln(jen.ID("url").Op(":").ID("invalidParsedURL"), jen.ID("logger").Op(":").ID("l"), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.Lit("test")), jen.ID("encoder").Op(":").ID("encoding").Dot("ProvideClientEncoder").Call(
				jen.ID("l"),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			)),
		),
		jen.Line(),
	)

	return code
}
