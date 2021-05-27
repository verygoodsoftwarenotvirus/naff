package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountRoleDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountRole").ID("role"),
			jen.ID("AccountRolePermissionsChecker").Interface(
				jen.ID("HasPermission").Params(jen.ID("Permission")).Params(jen.ID("bool")),
				jen.ID("CanUpdateAccounts").Params().Params(jen.ID("bool")),
				jen.ID("CanDeleteAccounts").Params().Params(jen.ID("bool")),
				jen.ID("CanAddMemberToAccounts").Params().Params(jen.ID("bool")),
				jen.ID("CanRemoveMemberFromAccounts").Params().Params(jen.ID("bool")),
				jen.ID("CanTransferAccountToNewOwner").Params().Params(jen.ID("bool")),
				jen.ID("CanCreateWebhooks").Params().Params(jen.ID("bool")),
				jen.ID("CanSeeWebhooks").Params().Params(jen.ID("bool")),
				jen.ID("CanUpdateWebhooks").Params().Params(jen.ID("bool")),
				jen.ID("CanArchiveWebhooks").Params().Params(jen.ID("bool")),
				jen.ID("CanCreateAPIClients").Params().Params(jen.ID("bool")),
				jen.ID("CanSeeAPIClients").Params().Params(jen.ID("bool")),
				jen.ID("CanDeleteAPIClients").Params().Params(jen.ID("bool")),
				jen.ID("CanSeeAuditLogEntriesForItems").Params().Params(jen.ID("bool")),
				jen.ID("CanSeeAuditLogEntriesForWebhooks").Params().Params(jen.ID("bool")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("AccountMemberRole").ID("AccountRole").Op("=").ID("iota"),
			jen.ID("AccountAdminRole").ID("AccountRole").Op("=").ID("iota"),
			jen.ID("accountAdminRoleName").Op("=").Lit("account_admin"),
			jen.ID("accountMemberRoleName").Op("=").Lit("account_member"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("accountAdmin").Op("=").ID("gorbac").Dot("NewStdRole").Call(jen.ID("accountAdminRoleName")),
			jen.ID("accountMember").Op("=").ID("gorbac").Dot("NewStdRole").Call(jen.ID("accountMemberRoleName")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().ID("accountRoleCollection").Struct(jen.ID("Roles").Index().ID("string")),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("encoding/gob", "Register").Call(jen.ID("accountRoleCollection").Values())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewAccountRolePermissionChecker returns a new checker for a set of Roles."),
		jen.Line(),
		jen.Func().ID("NewAccountRolePermissionChecker").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("AccountRolePermissionsChecker")).Body(
			jen.Return().Op("&").ID("accountRoleCollection").Valuesln(
				jen.ID("Roles").Op(":").ID("roles"))),
		jen.Line(),
	)

	code.Add(
		jen.Func().Params(jen.ID("r").ID("AccountRole")).ID("String").Params().Params(jen.ID("string")).Body(
			jen.Switch(jen.ID("r")).Body(
				jen.Case(jen.ID("AccountMemberRole")).Body(
					jen.Return().ID("accountMemberRoleName")),
				jen.Case(jen.ID("AccountAdminRole")).Body(
					jen.Return().ID("accountAdminRoleName")),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("HasPermission returns whether a user can do something or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("HasPermission").Params(jen.ID("p").ID("Permission")).Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("p"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanUpdateAccounts returns whether a user can update accounts or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanUpdateAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("UpdateAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanDeleteAccounts returns whether a user can delete accounts or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanDeleteAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ArchiveAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanAddMemberToAccounts returns whether a user can add members to accounts or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanAddMemberToAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("AddMemberAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanRemoveMemberFromAccounts returns whether a user can remove members from accounts or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanRemoveMemberFromAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("RemoveMemberAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanTransferAccountToNewOwner returns whether a user can transfer an account to a new owner or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanTransferAccountToNewOwner").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("TransferAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanCreateWebhooks returns whether a user can create webhooks or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanCreateWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("CreateWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanSeeWebhooks returns whether a user can view webhooks or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ReadWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanUpdateWebhooks returns whether a user can update webhooks or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanUpdateWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("UpdateWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanArchiveWebhooks returns whether a user can delete webhooks or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanArchiveWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ArchiveWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanCreateAPIClients returns whether a user can create API clients or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanCreateAPIClients").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("CreateAPIClientsPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanSeeAPIClients returns whether a user can view API clients or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAPIClients").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ReadAPIClientsPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanDeleteAPIClients returns whether a user can delete API clients or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanDeleteAPIClients").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ArchiveAPIClientsPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanSeeAuditLogEntriesForItems returns whether a user can view item audit log entries or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAuditLogEntriesForItems").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ReadItemsAuditLogEntriesPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CanSeeAuditLogEntriesForWebhooks returns whether a user can view item audit log entries or not."),
		jen.Line(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAuditLogEntriesForWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ReadWebhooksAuditLogEntriesPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Line(),
	)

	return code
}
