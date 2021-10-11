package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountRoleTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("TestNewAccountRolePermissionChecker").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account user"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("r").Assign().ID("NewAccountRolePermissionChecker").Call(jen.ID("AccountMemberRole").Dot("String").Call()),
					jen.Newline(),
					utils.AssertFalse(jen.ID("r").Dot("CanUpdateAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("UpdateAccountPermission")), jen.ID("r").Dot("CanUpdateAccounts").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanDeleteAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ArchiveAccountPermission")), jen.ID("r").Dot("CanDeleteAccounts").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanAddMemberToAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("AddMemberAccountPermission")), jen.ID("r").Dot("CanAddMemberToAccounts").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("RemoveMemberAccountPermission")), jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("TransferAccountPermission")), jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanCreateWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("CreateWebhooksPermission")), jen.ID("r").Dot("CanCreateWebhooks").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanSeeWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ReadWebhooksPermission")), jen.ID("r").Dot("CanSeeWebhooks").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanUpdateWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("UpdateWebhooksPermission")), jen.ID("r").Dot("CanUpdateWebhooks").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanArchiveWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ArchiveWebhooksPermission")), jen.ID("r").Dot("CanArchiveWebhooks").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanCreateAPIClients").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("CreateAPIClientsPermission")), jen.ID("r").Dot("CanCreateAPIClients").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanSeeAPIClients").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ReadAPIClientsPermission")), jen.ID("r").Dot("CanSeeAPIClients").Call(), nil),
					utils.AssertFalse(jen.ID("r").Dot("CanDeleteAPIClients").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ArchiveAPIClientsPermission")), jen.ID("r").Dot("CanDeleteAPIClients").Call(), nil),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account admin"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("r").Assign().ID("NewAccountRolePermissionChecker").Call(jen.ID("AccountAdminRole").Dot("String").Call()),
					jen.Newline(),
					utils.AssertTrue(jen.ID("r").Dot("CanUpdateAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("UpdateAccountPermission")), jen.ID("r").Dot("CanUpdateAccounts").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanDeleteAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ArchiveAccountPermission")), jen.ID("r").Dot("CanDeleteAccounts").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanAddMemberToAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("AddMemberAccountPermission")), jen.ID("r").Dot("CanAddMemberToAccounts").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("RemoveMemberAccountPermission")), jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("TransferAccountPermission")), jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanCreateWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("CreateWebhooksPermission")), jen.ID("r").Dot("CanCreateWebhooks").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanSeeWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ReadWebhooksPermission")), jen.ID("r").Dot("CanSeeWebhooks").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanUpdateWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("UpdateWebhooksPermission")), jen.ID("r").Dot("CanUpdateWebhooks").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanArchiveWebhooks").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ArchiveWebhooksPermission")), jen.ID("r").Dot("CanArchiveWebhooks").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanCreateAPIClients").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("CreateAPIClientsPermission")), jen.ID("r").Dot("CanCreateAPIClients").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanSeeAPIClients").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ReadAPIClientsPermission")), jen.ID("r").Dot("CanSeeAPIClients").Call(), nil),
					utils.AssertTrue(jen.ID("r").Dot("CanDeleteAPIClients").Call(), nil),
					utils.AssertEqual(jen.ID("r").Dot("HasPermission").Call(jen.ID("ArchiveAPIClientsPermission")), jen.ID("r").Dot("CanDeleteAPIClients").Call(), nil),
				),
			),
		),
		jen.Newline(),
	)

	return code
}
