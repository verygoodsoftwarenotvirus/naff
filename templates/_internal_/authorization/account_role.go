package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountRoleDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.Comment("AccountRole describes a role a user has for an Account context."),
			jen.ID("AccountRole").ID("role"),
			jen.Newline(),
			jen.Comment("AccountRolePermissionsChecker checks permissions for one or more account Roles."),
			jen.ID("AccountRolePermissionsChecker").Interface(
				jen.ID("HasPermission").Params(jen.ID("Permission")).Params(jen.ID("bool")),
				jen.Newline(),
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
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Const().Defs(
			jen.Comment("AccountMemberRole is a role for a plain account participant."),
			jen.ID("AccountMemberRole").ID("AccountRole").Equals().ID("iota"),
			jen.Comment("AccountAdminRole is a role for someone who can manipulate the specifics of an account."),
			jen.ID("AccountAdminRole").ID("AccountRole").Equals().ID("iota"),
			jen.Newline(),
			jen.ID("accountAdminRoleName").Equals().Lit("account_admin"),
			jen.ID("accountMemberRoleName").Equals().Lit("account_member"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("accountAdmin").Equals().Qual(constants.RBACLibrary, "NewStdRole").Call(jen.ID("accountAdminRoleName")),
			jen.ID("accountMember").Equals().Qual(constants.RBACLibrary, "NewStdRole").Call(jen.ID("accountMemberRoleName")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Type().ID("accountRoleCollection").Struct(jen.ID("Roles").Index().String()),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("encoding/gob", "Register").Call(jen.ID("accountRoleCollection").Values())),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("NewAccountRolePermissionChecker returns a new checker for a set of Roles."),
		jen.Newline(),
		jen.Func().ID("NewAccountRolePermissionChecker").Params(jen.ID("roles").Op("...").String()).Params(jen.ID("AccountRolePermissionsChecker")).Body(
			jen.Return().AddressOf().ID("accountRoleCollection").Valuesln(
				jen.ID("Roles").MapAssign().ID("roles"))),
		jen.Newline(),
	)

	code.Add(
		jen.Func().Params(jen.ID("r").ID("AccountRole")).ID("String").Params().Params(jen.String()).Body(
			jen.Switch(jen.ID("r")).Body(
				jen.Case(jen.ID("AccountMemberRole")).Body(
					jen.Return().ID("accountMemberRoleName")),
				jen.Case(jen.ID("AccountAdminRole")).Body(
					jen.Return().ID("accountAdminRoleName")),
				jen.Default().Body(
					jen.Return().Lit("")),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HasPermission returns whether a user can do something or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("HasPermission").Params(jen.ID("p").ID("Permission")).Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("p"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanUpdateAccounts returns whether a user can update accounts or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanUpdateAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("UpdateAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanDeleteAccounts returns whether a user can delete accounts or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanDeleteAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ArchiveAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanAddMemberToAccounts returns whether a user can add members to accounts or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanAddMemberToAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("AddMemberAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanRemoveMemberFromAccounts returns whether a user can remove members from accounts or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanRemoveMemberFromAccounts").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("RemoveMemberAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanTransferAccountToNewOwner returns whether a user can transfer an account to a new owner or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanTransferAccountToNewOwner").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("TransferAccountPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanCreateWebhooks returns whether a user can create webhooks or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanCreateWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("CreateWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanSeeWebhooks returns whether a user can view webhooks or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ReadWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanUpdateWebhooks returns whether a user can update webhooks or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanUpdateWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("UpdateWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanArchiveWebhooks returns whether a user can delete webhooks or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanArchiveWebhooks").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ArchiveWebhooksPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanCreateAPIClients returns whether a user can create API clients or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanCreateAPIClients").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("CreateAPIClientsPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanSeeAPIClients returns whether a user can view API clients or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanSeeAPIClients").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ReadAPIClientsPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CanDeleteAPIClients returns whether a user can delete API clients or not."),
		jen.Newline(),
		jen.Func().Params(jen.ID("r").ID("accountRoleCollection")).ID("CanDeleteAPIClients").Params().Params(jen.ID("bool")).Body(
			jen.Return().ID("hasPermission").Call(
				jen.ID("ArchiveAPIClientsPermission"),
				jen.ID("r").Dot("Roles").Op("..."),
			)),
		jen.Newline(),
	)

	return code
}
