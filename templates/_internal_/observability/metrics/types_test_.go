package metrics

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func typesTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestEnsureUnitCounter").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("ucp").Op(":=").Func().Params(jen.ID("string"), jen.ID("string")).Params(jen.ID("UnitCounter")).Body(
						jen.Return().Op("&").ID("noopUnitCounter").Values()),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("EnsureUnitCounter").Call(
							jen.ID("ucp"),
							jen.ID("logging").Dot("NewNoopLogger").Call(),
							jen.Lit(""),
							jen.Lit(""),
						),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil UnitCounterProvider"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("assert").Dot("NotNil").Call(
						jen.ID("t"),
						jen.ID("EnsureUnitCounter").Call(
							jen.ID("nil"),
							jen.ID("logging").Dot("NewNoopLogger").Call(),
							jen.Lit(""),
							jen.Lit(""),
						),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
