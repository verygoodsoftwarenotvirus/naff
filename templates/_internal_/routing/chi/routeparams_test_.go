package chi

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func routeparamsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestNewRouteParamManager").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("NewRouteParamManager").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_BuildRouteParamIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").Op("&").ID("chiRouteParamManager").Values(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleKey").Op(":=").Lit("blah"),
					jen.ID("fn").Op(":=").ID("r").Dot("BuildRouteParamIDFetcher").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("exampleKey"),
						jen.Lit("thing"),
					),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(123)),
					jen.ID("req").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildTestRequest").Call(jen.ID("t")).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("ctx"),
						jen.ID("chi").Dot("RouteCtxKey"),
						jen.Op("&").ID("chi").Dot("Context").Valuesln(
							jen.ID("URLParams").Op(":").ID("chi").Dot("RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Valuesln(
									jen.ID("exampleKey")), jen.ID("Values").Op(":").Index().ID("string").Valuesln(
									jen.Qual("strconv", "FormatUint").Call(
										jen.ID("expected"),
										jen.Lit(10),
									)))),
					)),
					jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expected"),
						jen.ID("actual"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with invalid value somehow"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").Op("&").ID("chiRouteParamManager").Values(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleKey").Op(":=").Lit("blah"),
					jen.ID("fn").Op(":=").ID("r").Dot("BuildRouteParamIDFetcher").Call(
						jen.ID("logging").Dot("NewNonOperationalLogger").Call(),
						jen.ID("exampleKey"),
						jen.Lit("thing"),
					),
					jen.ID("expected").Op(":=").ID("uint64").Call(jen.Lit(0)),
					jen.ID("req").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildTestRequest").Call(jen.ID("t")),
					jen.ID("req").Op("=").ID("req").Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("ctx"),
						jen.ID("chi").Dot("RouteCtxKey"),
						jen.Op("&").ID("chi").Dot("Context").Valuesln(
							jen.ID("URLParams").Op(":").ID("chi").Dot("RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Valuesln(
									jen.ID("exampleKey")), jen.ID("Values").Op(":").Index().ID("string").Valuesln(
									jen.Lit("expected")))),
					)),
					jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
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
		jen.Func().ID("Test_BuildRouteParamStringIDFetcher").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").Op("&").ID("chiRouteParamManager").Values(),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.ID("exampleKey").Op(":=").Lit("blah"),
					jen.ID("fn").Op(":=").ID("r").Dot("BuildRouteParamStringIDFetcher").Call(jen.ID("exampleKey")),
					jen.ID("expectedInt").Op(":=").ID("uint64").Call(jen.Lit(123)),
					jen.ID("expected").Op(":=").Qual("strconv", "FormatUint").Call(
						jen.ID("expectedInt"),
						jen.Lit(10),
					),
					jen.ID("req").Op(":=").Qual("gitlab.com/verygoodsoftwarenotvirus/todo/tests/utils", "BuildTestRequest").Call(jen.ID("t")).Dot("WithContext").Call(jen.Qual("context", "WithValue").Call(
						jen.ID("ctx"),
						jen.ID("chi").Dot("RouteCtxKey"),
						jen.Op("&").ID("chi").Dot("Context").Valuesln(
							jen.ID("URLParams").Op(":").ID("chi").Dot("RouteParams").Valuesln(
								jen.ID("Keys").Op(":").Index().ID("string").Valuesln(
									jen.ID("exampleKey")), jen.ID("Values").Op(":").Index().ID("string").Valuesln(
									jen.Qual("strconv", "FormatUint").Call(
										jen.ID("expectedInt"),
										jen.Lit(10),
									)))),
					)),
					jen.ID("actual").Op(":=").ID("fn").Call(jen.ID("req")),
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

	return code
}
