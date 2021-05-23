package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookEventsTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit_test")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("exampleWebhookID").ID("uint64").Op("=").Lit(123),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildWebhookCreationEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildWebhookCreationEventEntry").Call(jen.Op("&").Qual(proj.TypesPackage(), "Webhook").Values(),
					jen.ID("exampleUserID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildWebhookUpdateEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildWebhookUpdateEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
					jen.ID("exampleWebhookID"),
					jen.ID("nil"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestBuildWebhookArchiveEventEntry").Params(
			jen.ID("t").Op("*").Qual("testing", "T")).Body(
			jen.ID("t").Dot("Parallel").Call(),
			jen.ID("assert").Dot("NotNil").Call(jen.ID("t"),
				jen.Qual(proj.InternalAuditPackage(), "BuildWebhookArchiveEventEntry").Call(jen.ID("exampleUserID"),
					jen.ID("exampleAccountID"),
					jen.ID("exampleWebhookID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
