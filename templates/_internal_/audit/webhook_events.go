package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhookEventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("WebhookAssignmentKey").Op("=").Lit("webhook_id"),
			jen.ID("WebhookCreationEvent").Op("=").Lit("webhook_created"),
			jen.ID("WebhookUpdateEvent").Op("=").Lit("webhook_updated"),
			jen.ID("WebhookArchiveEvent").Op("=").Lit("webhook_archived"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildWebhookCreationEventEntry builds an entry creation input for when a webhook is created."),
		jen.Line(),
		jen.Func().ID("BuildWebhookCreationEventEntry").Params(
			jen.ID("webhook").Op("*").Qual(proj.TypesPackage(), "Webhook"),
			jen.ID("createdByUser").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("WebhookCreationEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("createdByUser"),
					jen.ID("CreationAssignmentKey").Op(":").ID("webhook"),
					jen.ID("WebhookAssignmentKey").Op(":").ID("webhook").Dot("ID"),
					jen.ID("AccountAssignmentKey").Op(":").ID("webhook").Dot("BelongsToAccount"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildWebhookUpdateEventEntry builds an entry creation input for when a webhook is updated."),
		jen.Line(),
		jen.Func().ID("BuildWebhookUpdateEventEntry").Params(
			jen.List(jen.ID("changedByUser"),
				jen.ID("accountID"),
				jen.ID("webhookID")).ID("uint64"),
			jen.ID("changes").Index().Op("*").Qual(proj.TypesPackage(),
				"FieldChangeSummary",
			)).Params(
			jen.Op("*").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			).Valuesln(
				jen.ID("EventType").Op(":").ID("WebhookUpdateEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("changedByUser"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.ID("WebhookAssignmentKey").Op(":").ID("webhookID"),
					jen.ID("ChangesAssignmentKey").Op(":").ID("changes"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildWebhookArchiveEventEntry builds an entry creation input for when a webhook is archived."),
		jen.Line(),
		jen.Func().ID("BuildWebhookArchiveEventEntry").Params(
			jen.List(jen.ID("archivedByUser"), jen.ID("accountID"), jen.ID("webhookID")).ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("WebhookArchiveEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("archivedByUser"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.ID("WebhookAssignmentKey").Op(":").ID("webhookID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
