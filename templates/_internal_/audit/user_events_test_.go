package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("exampleAdminUserID").ID("uint64").Op("=").Lit(321),
			jen.ID("exampleUserID").ID("uint64").Op("=").Lit(123),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserCreationEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildUserCreationEventEntry").Call(jen.ID("exampleUserID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserVerifyTwoFactorSecretEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUserVerifyTwoFactorSecretEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserUpdateTwoFactorSecretEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUserUpdateTwoFactorSecretEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserUpdatePasswordEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUserUpdatePasswordEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserUpdateEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUserUpdateEventEntry").Call(jen.ID("exampleUserID"), jen.ID("nil"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUserArchiveEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUserArchiveEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	return code
}
