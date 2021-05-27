package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func billingTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_renderPrice").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("inputsAndExpectations").Op(":=").Map(jen.ID("uint32")).ID("string").Valuesln(
						jen.Lit(12345).Op(":").Lit("$123.45"), jen.Lit(42069).Op(":").Lit("$420.69"), jen.Lit(666).Op(":").Lit("$6.66")),
					jen.For(jen.List(jen.ID("input"), jen.ID("expectation")).Op(":=").Range().ID("inputsAndExpectations")).Body(
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("expectation"),
							jen.ID("renderPrice").Call(jen.ID("input")),
						)),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with too large a number"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("fakePanicker").Op(":=").ID("panicking").Dot("NewMockPanicker").Call(),
					jen.ID("fakePanicker").Dot("On").Call(
						jen.Lit("Panic"),
						jen.ID("mock").Dot("IsType").Call(jen.Lit("")),
					).Dot("Return").Call(),
					jen.ID("pricePanicker").Op("=").ID("fakePanicker"),
					jen.ID("renderPrice").Call(jen.ID("arbitraryPriceMax").Op("+").Lit(1)),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("fakePanicker"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_service_handleCheckoutSessionStart").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("examplePlanID").Op(":=").Lit("example_plan"),
					jen.ID("exampleSessionID").Op(":=").Lit("example_session"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/whatever?plan=%s"),
							jen.ID("examplePlanID"),
						),
						jen.ID("nil"),
					),
					jen.ID("mpm").Op(":=").Op("&").ID("capitalism").Dot("MockPaymentManager").Values(),
					jen.ID("mpm").Dot("On").Call(
						jen.Lit("CreateCheckoutSession"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("examplePlanID"),
					).Dot("Return").Call(
						jen.ID("exampleSessionID"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("paymentManager").Op("=").ID("mpm"),
					jen.ID("s").Dot("handleCheckoutSessionStart").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusOK"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mpm"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with missing plan ID"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleCheckoutSessionStart").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusBadRequest"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with error creating checkout session"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("examplePlanID").Op(":=").Lit("example_plan"),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("/whatever?plan=%s"),
							jen.ID("examplePlanID"),
						),
						jen.ID("nil"),
					),
					jen.ID("mpm").Op(":=").Op("&").ID("capitalism").Dot("MockPaymentManager").Values(),
					jen.ID("mpm").Dot("On").Call(
						jen.Lit("CreateCheckoutSession"),
						jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "ContextMatcher"),
						jen.ID("examplePlanID"),
					).Dot("Return").Call(
						jen.Lit(""),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
					jen.ID("s").Dot("paymentManager").Op("=").ID("mpm"),
					jen.ID("s").Dot("handleCheckoutSessionStart").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusInternalServerError"),
						jen.ID("res").Dot("Code"),
					),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("mpm"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_service_handleCheckoutSuccess").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleCheckoutSuccess").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusTooEarly"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_service_handleCheckoutCancel").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleCheckoutCancel").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusTooEarly"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_service_handleCheckoutFailure").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("s").Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("res").Op(":=").ID("httptest").Dot("NewRecorder").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodPost"),
						jen.Lit("/whatever"),
						jen.ID("nil"),
					),
					jen.ID("s").Dot("handleCheckoutFailure").Call(
						jen.ID("res"),
						jen.ID("req"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.Qual("net/http", "StatusTooEarly"),
						jen.ID("res").Dot("Code"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
