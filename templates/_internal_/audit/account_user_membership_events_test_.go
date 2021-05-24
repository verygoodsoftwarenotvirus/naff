package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuildUserAddedToAccountEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildUserAddedToAccountEventEntry").Call(jen.ID("exampleAdminUserID"),
					jen.Op("&").Qual(proj.TypesPackage(), "AddUserToAccountInput").Valuesln(
						jen.ID("Reason").Op(":").ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserRemovedFromAccountEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildUserRemovedFromAccountEventEntry").Call(jen.ID("exampleAdminUserID"),
					jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
					jen.Lit("blah blah"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserMarkedAccountAsDefaultEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildUserMarkedAccountAsDefaultEventEntry").Call(jen.ID("exampleAdminUserID"),
					jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildModifyUserPermissionsEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(
				jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildModifyUserPermissionsEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
					jen.ID("exampleAdminUserID"),
					jen.Index().ID("string").Valuesln(
						jen.ID("t").Dot("Name").Call(),
					),
					jen.ID("t").Dot("Name").Call(),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildTransferAccountOwnershipEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(
				jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildTransferAccountOwnershipEventEntry").Call(
					jen.ID("exampleAccountID"), jen.ID("exampleAdminUserID"), jen.ID("fakes").Dot("BuildFakeTransferAccountOwnershipInput").Call(),
				),
			),
		),
		jen.Line(),
	)

	return code
}
