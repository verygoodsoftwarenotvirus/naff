package audit

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableEventsTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)
	n := typ.Name

	code.Add(
		jen.Const().Defs(
			jen.IDf("example%sID", n.Singular()).ID("uint64").Op("=").Lit(123),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().IDf("TestBuild%sCreationEventEntry", n.Singular()).Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sCreationEventEntry", n.Singular())).Call(jen.Op("&").Qual(proj.TypesPackage(), n.Singular()).Values(),
					jen.ID("exampleAccountID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().IDf("TestBuild%sUpdateEventEntry", n.Singular()).Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sUpdateEventEntry", n.Singular())).Call(jen.ID("exampleUserID"),
					jen.IDf("example%sID", n.Singular()),
					jen.ID("exampleAccountID"),
					jen.ID("nil"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().IDf("TestBuild%sArchiveEventEntry", n.Singular()).Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Line(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sArchiveEventEntry", n.Singular())).Call(jen.ID("exampleUserID"),
					jen.IDf("example%sID", n.Singular()),
					jen.ID("exampleAccountID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
