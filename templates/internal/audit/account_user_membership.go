package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("UserAddedToAccountEvent").Op("=").Lit("user_added_to_account"),
			jen.ID("UserAccountPermissionsModifiedEvent").Op("=").Lit("user_account_permissions_modified"),
			jen.ID("UserRemovedFromAccountEvent").Op("=").Lit("user_removed_from_account"),
			jen.ID("AccountMarkedAsDefaultEvent").Op("=").Lit("account_marked_as_default"),
			jen.ID("AccountTransferredEvent").Op("=").Lit("account_transferred"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("// BuildUserAddedToAccountEventEntry builds an entry creation input for when a membership is created."),
		jen.Line(),
		jen.Func().ID("BuildUserAddedToAccountEventEntry").Params(
			jen.ID("addedBy").ID("uint64"),
			jen.ID("input").Op("*").Qual(proj.TypesPackage(), "AddUserToAccountInput")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.ID("contextMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("ActorAssignmentKey").Op(":").ID("addedBy"),
				jen.ID("AccountAssignmentKey").Op(":").ID("input").Dot("AccountID"),
				jen.ID("UserAssignmentKey").Op(":").ID("input").Dot("UserID"),
			),
			jen.If(jen.ID("input").Dot("Reason").Op("!=").Lit("")).Body(
				jen.ID("contextMap").Index(jen.ID("ReasonKey")).Op("=").ID("input").Dot("Reason")),
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UserAddedToAccountEvent"),
				jen.ID("Context").Op(":").ID("contextMap"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("// BuildUserRemovedFromAccountEventEntry builds an entry creation input for when a membership is archived."),
		jen.Line(),
		jen.Func().ID("BuildUserRemovedFromAccountEventEntry").Params(
			jen.List(jen.ID("removedBy"), jen.ID("removed"), jen.ID("accountID")).ID("uint64"),
			jen.ID("reason").ID("string")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.ID("contextMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("ActorAssignmentKey").Op(":").ID("removedBy"),
				jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
				jen.ID("UserAssignmentKey").Op(":").ID("removed"),
			),
			jen.If(jen.ID("reason").Op("!=").Lit("")).Body(
				jen.ID("contextMap").Index(jen.ID("ReasonKey")).Op("=").ID("reason")),
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UserRemovedFromAccountEvent"),
				jen.ID("Context").Op(":").ID("contextMap"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("// BuildUserMarkedAccountAsDefaultEventEntry builds an entry creation input for when a membership is created."),
		jen.Line(),
		jen.Func().ID("BuildUserMarkedAccountAsDefaultEventEntry").Params(
			jen.List(jen.ID("performedBy"), jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.ID("contextMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("ActorAssignmentKey").Op(":").ID("performedBy"),
				jen.ID("UserAssignmentKey").Op(":").ID("userID"),
				jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
			),
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("AccountMarkedAsDefaultEvent"),
				jen.ID("Context").Op(":").ID("contextMap"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("// BuildModifyUserPermissionsEventEntry builds an entry creation input for when a membership is created."),
		jen.Line(),
		jen.Func().ID("BuildModifyUserPermissionsEventEntry").Params(
			jen.List(jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("modifiedBy")).ID("uint64"),
			jen.ID("newRoles").Index().ID("string"),
			jen.ID("reason").ID("string"),
		).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.ID("contextMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("ActorAssignmentKey").Op(":").ID("modifiedBy"),
				jen.ID("AccountAssignmentKey").Op(":").ID("accountID"),
				jen.ID("UserAssignmentKey").Op(":").ID("userID"),
				jen.ID("AccountRolesKey").Op(":").ID("newRoles"),
			),
			jen.If(jen.ID("reason").Op("!=").Lit("")).Body(
				jen.ID("contextMap").Index(jen.ID("ReasonKey")).Op("=").ID("reason")),
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UserAccountPermissionsModifiedEvent"),
				jen.ID("Context").Op(":").ID("contextMap"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("// BuildTransferAccountOwnershipEventEntry builds an entry creation input for when a membership is created."),
		jen.Line(),
		jen.Func().ID("BuildTransferAccountOwnershipEventEntry").Params(
			jen.List(jen.ID("accountID"),
				jen.ID("changedBy")).ID("uint64"),
			jen.ID("input").Op("*").Qual(proj.TypesPackage(), "AccountOwnershipTransferInput")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.ID("contextMap").Op(":=").Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("ActorAssignmentKey").Op(":").ID("changedBy"),
				jen.Lit("old_owner").Op(":").ID("input").Dot("CurrentOwner"),
				jen.Lit("new_owner").Op(":").ID("input").Dot("NewOwner"),
				jen.ID("AccountAssignmentKey").Op(":").ID("accountID")),
			jen.If(jen.ID("input").Dot("Reason").Op("!=").Lit("")).Body(
				jen.ID("contextMap").Index(jen.ID("ReasonKey")).Op("=").ID("input").Dot("Reason")),
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("AccountTransferredEvent"),
				jen.ID("Context").Op(":").ID("contextMap"),
			),
		),
		jen.Line(),
	)

	return code
}
