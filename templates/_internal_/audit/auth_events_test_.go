package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestBuildCycleCookieSecretEvent").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildCycleCookieSecretEvent").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildSuccessfulLoginEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildSuccessfulLoginEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBannedUserLoginAttemptEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildBannedUserLoginAttemptEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUnsuccessfulLoginBadPasswordEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUnsuccessfulLoginBadPasswordEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildUnsuccessfulLoginBad2FATokenEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildUnsuccessfulLoginBad2FATokenEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildLogoutEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"), jen.Qual(proj.InternalAuditPackage(), "BuildLogoutEventEntry").Call(jen.ID("exampleUserID"))),
		),
		jen.Line(),
	)

	return code
}
