package frontend

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func languagesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_detailsForLanguage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.For(jen.List(jen.ID("_"), jen.ID("lang")).Op(":=").Range().ID("supportedLanguages")).Body(
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("detailsForLanguage").Call(jen.ID("lang")),
						)),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns default for nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("englishDetails"),
						jen.ID("detailsForLanguage").Call(jen.ID("nil")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_determineLanguage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.For(jen.List(jen.ID("expected"), jen.ID("deets")).Op(":=").Range().ID("languageDetails")).Body(
						jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
							jen.Qual("net/http", "MethodGet"),
							jen.Lit("/things"),
							jen.ID("nil"),
						),
						jen.ID("req").Dot("Header").Dot("Set").Call(
							jen.Lit("Accept-Language"),
							jen.ID("deets").Dot("Abbreviation"),
						),
						jen.ID("actual").Op(":=").ID("determineLanguage").Call(jen.ID("req")),
						jen.ID("assert").Dot("Equal").Call(
							jen.ID("t"),
							jen.ID("expected"),
							jen.ID("actual"),
							jen.Lit("expected result to be interpreted as English"),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.Lit("Accept-Language"),
						jen.Lit("en-US"),
					),
					jen.ID("actual").Op(":=").ID("determineLanguage").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("english"),
						jen.ID("actual"),
						jen.Lit("expected result to be interpreted as English"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns default for nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("actual").Op(":=").ID("determineLanguage").Call(jen.ID("nil")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("defaultLanguage"),
						jen.ID("actual"),
						jen.Lit("expected result to be interpreted as English"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns default for invalid language header"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.Lit("Accept-Language"),
						jen.Lit(""),
					),
					jen.ID("actual").Op(":=").ID("determineLanguage").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("defaultLanguage"),
						jen.ID("actual"),
						jen.Lit("expected result to be interpreted as English"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns default for language header that yields no results but does not error"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.Lit("Accept-Language"),
						jen.Lit("fleeb-FLORP"),
					),
					jen.ID("actual").Op(":=").ID("determineLanguage").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("defaultLanguage"),
						jen.ID("actual"),
						jen.Lit("expected result to be interpreted as English"),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("returns default for language not found"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("req").Op(":=").ID("httptest").Dot("NewRequest").Call(
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/things"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.Lit("Accept-Language"),
						jen.Lit("zu-HM"),
					),
					jen.ID("actual").Op(":=").ID("determineLanguage").Call(jen.ID("req")),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("defaultLanguage"),
						jen.ID("actual"),
						jen.Lit("expected result to be interpreted as English"),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
