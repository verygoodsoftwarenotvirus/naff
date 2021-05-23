package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountRoleDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Null(),

		jen.Line(),
	)
	code.Add(jen.Null().Type().ID("AccountRole").ID("role").Type().ID("AccountRolePermissionsChecker").Interface(
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

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("AccountMemberRole").ID("AccountRole").Op("=").ID("iota").Var().ID("AccountAdminRole").ID("AccountRole").Op("=").ID("iota").Var().ID("accountAdminRoleName").Op("=").Lit("account_admin").Var().ID("accountMemberRoleName").Op("=").Lit("account_member"),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("accountAdmin").Op("=").ID("gorbac").Dot(
		"NewStdRole",
	).Call(jen.ID("accountAdminRoleName")).Var().ID("accountMember").Op("=").ID("gorbac").Dot(
		"NewStdRole",
	).Call(jen.ID("accountMemberRoleName")),

		jen.Line(),
	)
	code.Add(jen.Null().Type().ID("accountRoleCollection").Struct(jen.ID("Roles").Index().ID("string")),

		jen.Line(),
	)
	code.Add(jen.Func().ID("init").Params().Body(jen.Qual("encoding/gob", "Register").Call(jen.ID("accountRoleCollection").Valuesln())),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// NewAccountRolePermissionChecker returns a new checker for a set of Roles.").ID("NewAccountRolePermissionChecker").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("AccountRolePermissionsChecker")).Body(jen.Return().Op("&").ID("accountRoleCollection").Valuesln(jen.ID("Roles").Op(":").ID("roles"))),

		jen.Line(),
	)
	code.Add(jen.Func().Params(jen.ID("r").ID("AccountRole")).ID("String").Params().Params(jen.ID("string")).Body(jen.Switch(jen.ID("r")).Body(
		jen.Case(jen.ID("AccountMemberRole")).Body(jen.Return().ID("accountMemberRoleName")),
		jen.Case(jen.ID("AccountAdminRole")).Body(jen.Return().ID("accountAdminRoleName")),
		jen.Default().Body(jen.Return().Lit("")),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// HasPermission returns whether a user can do something or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("HasPermission").Params(jen.ID("p").ID("Permission")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("p"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanUpdateAccounts returns whether a user can update accounts or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanUpdateAccounts").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("UpdateAccountPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanDeleteAccounts returns whether a user can delete accounts or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanDeleteAccounts").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ArchiveAccountPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanAddMemberToAccounts returns whether a user can add members to accounts or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanAddMemberToAccounts").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("AddMemberAccountPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanRemoveMemberFromAccounts returns whether a user can remove members from accounts or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanRemoveMemberFromAccounts").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("RemoveMemberAccountPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanTransferAccountToNewOwner returns whether a user can transfer an account to a new owner or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanTransferAccountToNewOwner").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("TransferAccountPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanCreateWebhooks returns whether a user can create webhooks or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanCreateWebhooks").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("CreateWebhooksPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeWebhooks returns whether a user can view webhooks or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeWebhooks").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadWebhooksPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanUpdateWebhooks returns whether a user can update webhooks or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanUpdateWebhooks").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("UpdateWebhooksPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanArchiveWebhooks returns whether a user can delete webhooks or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanArchiveWebhooks").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ArchiveWebhooksPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanCreateAPIClients returns whether a user can create API clients or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanCreateAPIClients").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("CreateAPIClientsPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeAPIClients returns whether a user can view API clients or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAPIClients").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadAPIClientsPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanDeleteAPIClients returns whether a user can delete API clients or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanDeleteAPIClients").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ArchiveAPIClientsPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeAuditLogEntriesForItems returns whether a user can view item audit log entries or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAuditLogEntriesForItems").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadItemsAuditLogEntriesPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeAuditLogEntriesForWebhooks returns whether a user can view item audit log entries or not.").Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAuditLogEntriesForWebhooks").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadWebhooksAuditLogEntriesPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	return code
}
