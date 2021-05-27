package panicking

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func standardTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestNewProductionPanicker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("NewProductionPanicker").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stdLibPanicker_Panic").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("p").Op(":=").ID("NewProductionPanicker").Call(),
					jen.Defer().Func().Params().Body(
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("recover").Call(),
							jen.Lit("expected panic to occur"),
						)).Call(),
					jen.ID("p").Dot("Panic").Call(jen.Lit("blah")),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_stdLibPanicker_Panicf").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("p").Op(":=").ID("NewProductionPanicker").Call(),
					jen.Defer().Func().Params().Body(
						jen.ID("assert").Dot("NotNil").Call(
							jen.ID("t"),
							jen.ID("recover").Call(),
							jen.Lit("expected panic to occur"),
						)).Call(),
					jen.ID("p").Dot("Panicf").Call(jen.Lit("blah")),
				),
			),
		),
		jen.Line(),
	)

	return code
}
