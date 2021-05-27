package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("exampleAccountID").ID("uint64").Op("=").Lit(123),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildAccountCreationEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildAccountCreationEventEntry").Call(jen.Op("&").Qual(proj.TypesPackage(), "Account").Values(),
					jen.ID("exampleUserID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildAccountUpdateEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildAccountUpdateEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
					jen.ID("exampleUserID"),
					jen.ID("nil"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildAccountArchiveEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildAccountArchiveEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
					jen.ID("exampleUserID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
