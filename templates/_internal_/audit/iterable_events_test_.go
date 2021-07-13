package audit

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterableEventsTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)
	n := typ.Name

	code.Add(
		jen.Const().Defs(
			jen.IDf("example%sID", n.Singular()).Uint64().Equals().Lit(123),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuild%sCreationEventEntry", n.Singular()).Params(
			jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Qual(constants.AssertionLibrary, "NotNil").Call(
				jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sCreationEventEntry", n.Singular())).Call(
					jen.AddressOf().Qual(proj.TypesPackage(), n.Singular()).Values(),
					jen.ID("exampleAccountID"),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuild%sUpdateEventEntry", n.Singular()).Params(
			jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Qual(constants.AssertionLibrary, "NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sUpdateEventEntry", n.Singular())).Call(
					jen.ID("exampleUserID"),
					jen.IDf("example%sID", n.Singular()),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID("exampleAccountID")
						}
						return jen.Null()
					}(),
					jen.ID("nil"),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().IDf("TestBuild%sArchiveEventEntry", n.Singular()).Params(
			jen.ID("t").PointerTo().Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.Newline(),
			jen.Qual(constants.AssertionLibrary, "NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sArchiveEventEntry", n.Singular())).Call(jen.ID("exampleUserID"),
					jen.IDf("example%sID", n.Singular()),
					func() jen.Code {
						if typ.BelongsToAccount {
							return jen.ID("exampleAccountID")
						}
						return jen.Null()
					}(),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
