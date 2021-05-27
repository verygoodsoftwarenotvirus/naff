package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func apiClientEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("exampleAPIClientDatabaseID").ID("uint64").Op("=").Lit(123),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildAPIClientCreationEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildAPIClientCreationEventEntry").Call(jen.Op("&").Qual(proj.TypesPackage(), "APIClient").Values(), jen.ID("exampleUserID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildAPIClientArchiveEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildAPIClientArchiveEventEntry").Call(
					jen.ID("exampleAccountID"),
					jen.ID("exampleAPIClientDatabaseID"),
					jen.ID("exampleUserID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
