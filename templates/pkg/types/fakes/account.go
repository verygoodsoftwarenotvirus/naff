package fakes

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeAccount builds a faked account."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccount").Params().Params(jen.Op("*").ID("types").Dot("Account")).Body(
			jen.Return().Op("&").ID("types").Dot("Account").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("Name").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Word").Call(), jen.ID("BillingStatus").Op(":").ID("types").Dot("PaidAccountBillingStatus"), jen.ID("ContactEmail").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Email").Call(), jen.ID("ContactPhone").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "PhoneFormatted").Call(), jen.ID("PaymentProcessorCustomerID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("BelongsToUser").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("Members").Op(":").ID("BuildFakeAccountUserMembershipList").Call().Dot("AccountUserMemberships"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountForUser builds a faked account."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountForUser").Params(jen.ID("u").Op("*").ID("types").Dot("User")).Params(jen.Op("*").ID("types").Dot("Account")).Body(
			jen.Return().Op("&").ID("types").Dot("Account").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("ExternalID").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "UUID").Call(), jen.ID("Name").Op(":").ID("u").Dot("Username"), jen.ID("CreatedOn").Op(":").ID("uint64").Call(jen.ID("uint32").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Date").Call().Dot("Unix").Call())), jen.ID("BelongsToUser").Op(":").ID("u").Dot("ID"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountList builds a faked AccountList."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountList").Params().Params(jen.Op("*").ID("types").Dot("AccountList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("Account"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeAccount").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("AccountList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("FilteredCount").Op(":").ID("exampleQuantity").Op("/").Lit(2), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("Accounts").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountUpdateInput builds a faked AccountUpdateInput from an account."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUpdateInput").Params().Params(jen.Op("*").ID("types").Dot("AccountUpdateInput")).Body(
			jen.ID("account").Op(":=").ID("BuildFakeAccount").Call(),
			jen.Return().Op("&").ID("types").Dot("AccountUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("account").Dot("Name"), jen.ID("BelongsToUser").Op(":").ID("account").Dot("BelongsToUser")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountUpdateInputFromAccount builds a faked AccountUpdateInput from an account."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUpdateInputFromAccount").Params(jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.Op("*").ID("types").Dot("AccountUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountUpdateInput").Valuesln(jen.ID("Name").Op(":").ID("account").Dot("Name"), jen.ID("BelongsToUser").Op(":").ID("account").Dot("BelongsToUser"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountCreationInput builds a faked AccountCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountCreationInput").Params().Params(jen.Op("*").ID("types").Dot("AccountCreationInput")).Body(
			jen.ID("account").Op(":=").ID("BuildFakeAccount").Call(),
			jen.Return().ID("BuildFakeAccountCreationInputFromAccount").Call(jen.ID("account")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountCreationInputFromAccount builds a faked AccountCreationInput from an account."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountCreationInputFromAccount").Params(jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.Op("*").ID("types").Dot("AccountCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountCreationInput").Valuesln(jen.ID("Name").Op(":").ID("account").Dot("Name"), jen.ID("ContactEmail").Op(":").ID("account").Dot("ContactEmail"), jen.ID("ContactPhone").Op(":").ID("account").Dot("ContactPhone"), jen.ID("BelongsToUser").Op(":").ID("account").Dot("BelongsToUser"))),
		jen.Line(),
	)

	return code
}
