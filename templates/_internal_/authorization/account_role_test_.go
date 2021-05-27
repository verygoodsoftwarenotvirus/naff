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
		jen.Func().ID("TestNewAccountRolePermissionChecker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("NewAccountRolePermissionChecker").Call(jen.ID("AccountMemberRole").Dot("String").Call()),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanUpdateAccounts").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanDeleteAccounts").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanAddMemberToAccounts").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanCreateWebhooks").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeWebhooks").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanUpdateWebhooks").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanArchiveWebhooks").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanCreateAPIClients").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeAPIClients").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanDeleteAPIClients").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeAuditLogEntriesForItems").Call(),
					),
					jen.ID("assert").Dot("False").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeAuditLogEntriesForWebhooks").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account admin"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("r").Op(":=").ID("NewAccountRolePermissionChecker").Call(jen.ID("AccountAdminRole").Dot("String").Call()),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanUpdateAccounts").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanDeleteAccounts").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanAddMemberToAccounts").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanCreateWebhooks").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeWebhooks").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanUpdateWebhooks").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanArchiveWebhooks").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanCreateAPIClients").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeAPIClients").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanDeleteAPIClients").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeAuditLogEntriesForItems").Call(),
					),
					jen.ID("assert").Dot("True").Call(
						jen.ID("t"),
						jen.ID("r").Dot("CanSeeAuditLogEntriesForWebhooks").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
