package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminEventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildUserBanEventEntry builds an entry creation input for when a user is banned."),
		jen.Line(),
		jen.Func().ID("BuildUserBanEventEntry").Params(
			jen.List(jen.ID("banGiver"), jen.ID("banRecipient")).ID("uint64"),
			jen.ID("reason").ID("string"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UserBannedEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("banGiver"),
					jen.ID("UserAssignmentKey").Op(":").ID("banRecipient"),
					jen.ID("ReasonKey").Op(":").ID("reason"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAccountTerminationEventEntry builds an entry creation input for when an account is terminated."),
		jen.Line(),
		jen.Func().ID("BuildAccountTerminationEventEntry").Params(
			jen.List(jen.ID("terminator"), jen.ID("terminee")).ID("uint64"),
			jen.ID("reason").ID("string")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("AccountTerminatedEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("terminator"),
					jen.ID("UserAssignmentKey").Op(":").ID("terminee"),
					jen.ID("ReasonKey").Op(":").ID("reason"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
