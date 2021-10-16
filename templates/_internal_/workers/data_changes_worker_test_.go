package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func dataChangesWorkerTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestProvideDataChangesWorker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("actual").Assign().ID("ProvideDataChangesWorker").Call(jen.Qual(proj.InternalLoggingPackage(), "NewZerologLogger").Call()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestDataChangesWorker_HandleMessage").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("actual").Assign().ID("ProvideDataChangesWorker").Call(jen.Qual(proj.InternalLoggingPackage(), "NewZerologLogger").Call()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("actual").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.Index().ID("byte").Call(jen.Lit("{}")),
						),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("invalid input"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("actual").Assign().ID("ProvideDataChangesWorker").Call(jen.Qual(proj.InternalLoggingPackage(), "NewZerologLogger").Call()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("actual"),
					),
					jen.Newline(),
					jen.ID("ctx").Assign().Qual("context", "Background").Call(),
					jen.ID("assert").Dot("Error").Call(
						jen.ID("t"),
						jen.ID("actual").Dot("HandleMessage").Call(
							jen.ID("ctx"),
							jen.Index().ID("byte").Call(jen.Lit("} bad JSON lol")),
						),
					),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
