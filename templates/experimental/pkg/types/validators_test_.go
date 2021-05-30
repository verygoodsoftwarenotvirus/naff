package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func validatorsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_urlValidator_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("urlValidator").Valuesln(),
					jen.ID("assert").Dot("Nil").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.Lit("https://verygoodsoftwarenotvirus.ru")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("unhappy path"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("urlValidator").Valuesln(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("%s://verygoodsoftwarenotvirus.ru"),
							jen.ID("string").Call(jen.ID("byte").Call(jen.Lit(127))),
						)),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("invalid value"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("urlValidator").Valuesln(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.ID("validating").Dot("F").Call(
							jen.Lit("arbitrary"),
							jen.Lit(123),
						)),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stringDurationValidator_Validate").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("stringDurationValidator").Valuesln(jen.ID("maxDuration").Op(":").Qual("time", "Hour")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.Qual("time", "Minute").Dot("String").Call()),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("invalid value"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("stringDurationValidator").Valuesln(jen.ID("maxDuration").Op(":").Qual("time", "Hour")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.Lit(1234)),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("invalid format"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("stringDurationValidator").Valuesln(jen.ID("maxDuration").Op(":").Qual("time", "Hour")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.Lit("fake lol")),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("too large a max duration"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("x").Op(":=").Op("&").ID("stringDurationValidator").Valuesln(jen.ID("maxDuration").Op(":").Qual("time", "Hour")),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("x").Dot("Validate").Call(jen.Parens(jen.Lit(2400).Op("*").Qual("time", "Hour")).Dot("String").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
