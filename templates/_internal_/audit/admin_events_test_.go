package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuildUserBanEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildUserBanEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleUserID"),
					jen.Lit("reason"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildAccountTerminationEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildAccountTerminationEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleUserID"),
					jen.Lit("reason"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
