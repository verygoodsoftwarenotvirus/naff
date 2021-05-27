package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func helpersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("exampleInvalidForm").Qual("io", "Reader").Op("=").Qual("strings", "NewReader").Call(jen.Lit("a=|%%%=%%%%%%")),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_buildRedirectURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("buildRedirectURL").Call(
						jen.Lit("/from"),
						jen.Lit("/to"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_pluckRedirectURL").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.ID("expected").Op(":=").Lit(""),
					jen.ID("actual").Op(":=").ID("pluckRedirectURL").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_htmxRedirectTo").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("htmxRedirectTo").Call(
						jen.ID("res"),
						jen.Lit("/example"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_parseListOfTemplates").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleTemplateA").Op(":=").Lit(`<div> hi </div>`),
					jen.ID("exampleTemplateB").Op(":=").Lit(`<div> bye </div>`),
					jen.ID("actual").Op(":=").ID("parseListOfTemplates").Call(
						jen.ID("nil"),
						jen.ID("exampleTemplateA"),
						jen.ID("exampleTemplateB"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_renderStringToResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("thing").Op(":=").ID("t").Dot("Name").Call(),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("renderStringToResponse").Call(
						jen.ID("thing"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_renderBytesToResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("thing").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("renderBytesToResponse").Call(
						jen.ID("thing"),
						jen.ID("res"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error writing response"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("thing").Op(":=").Index().ID("byte").Call(jen.ID("t").Dot("Name").Call()),
					jen.ID("res").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockHTTPResponseWriter").Values(),
					jen.ID("res").Dot("On").Call(
						jen.Lit("Write"),
						jen.ID("mock").Dot("Anything"),
					).Dot("Return").Call(
						jen.Op("-").Lit(1),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("s").Dot("renderBytesToResponse").Call(
						jen.ID("thing"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_mergeFuncMaps").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("inputA").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
						jen.Lit("things").Op(":").Func().Params().Body()),
					jen.ID("inputB").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
						jen.Lit("stuff").Op(":").Func().Params().Body()),
					jen.ID("expected").Op(":=").Qual("html/template", "FuncMap").Valuesln(
						jen.Lit("things").Op(":").Func().Params().Body(), jen.Lit("stuff").Op(":").Func().Params().Body()),
					jen.ID("actual").Op(":=").ID("mergeFuncMaps").Call(
						jen.ID("inputA"),
						jen.ID("inputB"),
					),
					jen.For(jen.ID("key").Op(":=").Range().ID("expected")).Body(
						jen.ID("assert").Dot("Contains").Call(
							jen.ID("t"),
							jen.ID("actual"),
							jen.ID("key"),
						)),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_extractFormFromRequest").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("expected").Op(":=").Qual("net/url", "Values").Valuesln(
						jen.Lit("things").Op(":").Index().ID("string").Valuesln(
							jen.Lit("stuff"))),
					jen.ID("exampleReq").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/things"),
						jen.Qual("strings", "NewReader").Call(jen.ID("expected").Dot("Encode").Call()),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
						jen.ID("ctx"),
						jen.ID("exampleReq"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil request body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleBody").Op(":=").Op("&").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "MockReadCloser").Values(),
					jen.ID("exampleBody").Dot("On").Call(
						jen.Lit("Read"),
						jen.ID("mock").Dot("Anything"),
					).Dot("Return").Call(
						jen.Lit(0),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("exampleReq").Op(":=").Op("&").Qual("net/http", "Request").Valuesln(
						jen.ID("Body").Op(":").ID("exampleBody")),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
						jen.ID("ctx"),
						jen.ID("exampleReq"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid body"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleReq").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/things"),
						jen.ID("exampleInvalidForm"),
					),
					jen.List(jen.ID("actual"), jen.ID("err")).Op(":=").ID("s").Dot("extractFormFromRequest").Call(
						jen.ID("ctx"),
						jen.ID("exampleReq"),
					),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_renderTemplateToResponse").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleTemplate").Op(":=").Lit(`<div> hi </div>`),
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleTemplate"),
						jen.ID("nil"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.ID("nil"),
						jen.ID("res"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleTemplate").Op(":=").Lit(`<div> {{ .Something }} </div>`),
					jen.ID("tmpl").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleTemplate"),
						jen.ID("nil"),
					),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("s").Dot("renderTemplateToResponse").Call(
						jen.ID("ctx"),
						jen.ID("tmpl"),
						jen.Struct(jen.ID("Thing").ID("string")).Values(),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_renderTemplateIntoBaseTemplate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("s").Dot("renderTemplateIntoBaseTemplate").Call(
							jen.Lit("<div>hi</div>"),
							jen.ID("nil"),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestService_parseTemplate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("exampleTemplate").Op(":=").Lit(`<div> hi </div>`),
					jen.ID("actual").Op(":=").ID("s").Dot("parseTemplate").Call(
						jen.ID("ctx"),
						jen.Lit(""),
						jen.ID("exampleTemplate"),
						jen.ID("nil"),
					),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
