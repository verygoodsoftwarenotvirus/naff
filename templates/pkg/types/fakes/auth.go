package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeSessionContextData builds a faked SessionContextData."),
		jen.Line(),
		jen.Func().ID("BuildFakeSessionContextData").Params().Params(jen.Op("*").ID("types").Dot("SessionContextData")).Body(
			jen.ID("fakeAccountID").Op(":=").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(),
			jen.Return().Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("fakeAccountID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())), jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("Reputation").Op(":").ID("types").Dot("GoodStandingAccountStatus"), jen.ID("ReputationExplanation").Op(":").Lit(""), jen.ID("UserID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call())), jen.ID("ActiveAccountID").Op(":").ID("fakeAccountID")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeSessionContextDataForAccount builds a faked SessionContextData."),
		jen.Line(),
		jen.Func().ID("BuildFakeSessionContextDataForAccount").Params(jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.Op("*").ID("types").Dot("SessionContextData")).Body(
			jen.ID("fakeAccountID").Op(":=").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(),
			jen.Return().Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("AccountPermissions").Op(":").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Valuesln(jen.ID("account").Dot("ID").Op(":").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call())), jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("Reputation").Op(":").ID("types").Dot("GoodStandingAccountStatus"), jen.ID("ReputationExplanation").Op(":").Lit(""), jen.ID("UserID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call())), jen.ID("ActiveAccountID").Op(":").ID("fakeAccountID")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAddUserToAccountInput builds a faked AddUserToAccountInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeAddUserToAccountInput").Params().Params(jen.Op("*").ID("types").Dot("AddUserToAccountInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("Reason").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Sentence").Call(jen.Lit(10)), jen.ID("UserID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("AccountID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeUserPermissionModificationInput builds a faked ModifyUserPermissionsInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeUserPermissionModificationInput").Params().Params(jen.Op("*").ID("types").Dot("ModifyUserPermissionsInput")).Body(
			jen.Return().Op("&").ID("types").Dot("ModifyUserPermissionsInput").Valuesln(jen.ID("Reason").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Sentence").Call(jen.Lit(10)), jen.ID("NewRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeTransferAccountOwnershipInput builds a faked AccountOwnershipTransferInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeTransferAccountOwnershipInput").Params().Params(jen.Op("*").ID("types").Dot("AccountOwnershipTransferInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountOwnershipTransferInput").Valuesln(jen.ID("Reason").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Sentence").Call(jen.Lit(10)), jen.ID("CurrentOwner").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("NewOwner").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeChangeActiveAccountInput builds a faked ChangeActiveAccountInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeChangeActiveAccountInput").Params().Params(jen.Op("*").ID("types").Dot("ChangeActiveAccountInput")).Body(
			jen.Return().Op("&").ID("types").Dot("ChangeActiveAccountInput").Valuesln(jen.ID("AccountID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakePASETOCreationInput builds a faked PASETOCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakePASETOCreationInput").Params().Params(jen.Op("*").ID("types").Dot("PASETOCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("PASETOCreationInput").Valuesln(jen.ID("ClientID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("RequestTime").Op(":").Qual("time", "Now").Call().Dot("Unix").Call())),
		jen.Line(),
	)

	return code
}
