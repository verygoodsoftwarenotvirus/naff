package authorization

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountRoleTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	trueAssertions := []jen.Code{}
	falseAssertions := []jen.Code{}
	for _, typ := range proj.DataTypes {
		trueAssertions = append(trueAssertions, utils.AssertTrue(jen.ID("r").Dotf("CanSeeAuditLogEntriesFor%s", typ.Name.Plural()).Call(), nil))
		falseAssertions = append(falseAssertions, utils.AssertFalse(jen.ID("r").Dotf("CanSeeAuditLogEntriesFor%s", typ.Name.Plural()).Call(), nil))
	}

	code.Add(
		jen.Func().ID("TestNewAccountRolePermissionChecker").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account user"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("r").Op(":=").ID("NewAccountRolePermissionChecker").Call(jen.ID("AccountMemberRole").Dot("String").Call()),
						jen.Newline(),
						utils.AssertFalse(jen.ID("r").Dot("CanUpdateAccounts").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanDeleteAccounts").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanAddMemberToAccounts").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanCreateWebhooks").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanSeeWebhooks").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanUpdateWebhooks").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanArchiveWebhooks").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanCreateAPIClients").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanSeeAPIClients").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanDeleteAPIClients").Call(), nil),
						utils.AssertFalse(jen.ID("r").Dot("CanSeeAuditLogEntriesForWebhooks").Call(), nil),
					}, falseAssertions...)...,
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("account admin"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					append([]jen.Code{
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.ID("r").Op(":=").ID("NewAccountRolePermissionChecker").Call(jen.ID("AccountAdminRole").Dot("String").Call()),
						jen.Newline(),
						utils.AssertTrue(jen.ID("r").Dot("CanUpdateAccounts").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanDeleteAccounts").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanAddMemberToAccounts").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanRemoveMemberFromAccounts").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanTransferAccountToNewOwner").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanCreateWebhooks").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanSeeWebhooks").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanUpdateWebhooks").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanArchiveWebhooks").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanCreateAPIClients").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanSeeAPIClients").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanDeleteAPIClients").Call(), nil),
						utils.AssertTrue(jen.ID("r").Dot("CanSeeAuditLogEntriesForWebhooks").Call(), nil),
					}, trueAssertions...)...,
				),
			),
		),
		jen.Newline(),
	)

	return code
}
