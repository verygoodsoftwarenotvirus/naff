package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountEventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.Comment("AccountAssignmentKey is the key we use to indicate that an audit log entry is associated with an account."),
			jen.ID("AccountAssignmentKey").Op("=").Lit("account_id"),
			jen.Comment("AccountCreationEvent events indicate a user created an account."),
			jen.ID("AccountCreationEvent").Op("=").Lit("account_created"),
			jen.Comment("AccountUpdateEvent events indicate a user updated an account."),
			jen.ID("AccountUpdateEvent").Op("=").Lit("account_updated"),
			jen.Comment("AccountArchiveEvent events indicate a user deleted an account."),
			jen.ID("AccountArchiveEvent").Op("=").Lit("account_archived"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAccountCreationEventEntry builds an entry creation input for when an account is created."),
		jen.Line(),
		jen.Func().ID("BuildAccountCreationEventEntry").Params(
			jen.ID("account").Op("*").Qual(proj.TypesPackage(), "Account"),
			jen.ID("createdByUser").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("AccountCreationEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("AccountAssignmentKey").Op(":").ID("account").Dot("ID"),
					jen.ID("UserAssignmentKey").Op(":").ID("account").Dot("BelongsToUser"),
					jen.ID("ActorAssignmentKey").Op(":").ID("createdByUser"),
					jen.ID("CreationAssignmentKey").Op(":").ID("account"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAccountUpdateEventEntry builds an entry creation input for when an account is updated."),
		jen.Line(),
		jen.Func().ID("BuildAccountUpdateEventEntry").Params(
			jen.List(jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("changedByUser")).ID("uint64"),
			jen.ID("changes").Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("AccountUpdateEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("UserAssignmentKey").Op(":").ID("userID"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.ID("ChangesAssignmentKey").Op(":").ID("changes"),
					jen.ID("ActorAssignmentKey").Op(":").ID("changedByUser"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildAccountArchiveEventEntry builds an entry creation input for when an account is archived."),
		jen.Line(),
		jen.Func().ID("BuildAccountArchiveEventEntry").Params(
			jen.List(jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("archivedByUser")).ID("uint64"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("AccountArchiveEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("UserAssignmentKey").Op(":").ID("userID"),
					jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
					jen.ID("ActorAssignmentKey").Op(":").ID("archivedByUser"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
