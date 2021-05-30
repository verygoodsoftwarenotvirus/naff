package fakes

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("BuildFakeAccountUserMembership builds a faked item."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUserMembership").Params().Params(jen.Op("*").ID("types").Dot("AccountUserMembership")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountUserMembership").Valuesln(jen.ID("ID").Op(":").ID("uint64").Call(jen.Qual("github.com/brianvoe/gofakeit/v5", "Uint32").Call()), jen.ID("BelongsToUser").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("BelongsToAccount").Op(":").Qual("github.com/brianvoe/gofakeit/v5", "Uint64").Call(), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()), jen.ID("CreatedOn").Op(":").Lit(0), jen.ID("ArchivedOn").Op(":").ID("nil"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountUserMembershipList builds a faked AccountUserMembershipList."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUserMembershipList").Params().Params(jen.Op("*").ID("types").Dot("AccountUserMembershipList")).Body(
			jen.Var().Defs(
				jen.ID("examples").Index().Op("*").ID("types").Dot("AccountUserMembership"),
			),
			jen.For(jen.ID("i").Op(":=").Lit(0), jen.ID("i").Op("<").ID("exampleQuantity"), jen.ID("i").Op("++")).Body(
				jen.ID("examples").Op("=").ID("append").Call(
					jen.ID("examples"),
					jen.ID("BuildFakeAccountUserMembership").Call(),
				)),
			jen.Return().Op("&").ID("types").Dot("AccountUserMembershipList").Valuesln(jen.ID("Pagination").Op(":").ID("types").Dot("Pagination").Valuesln(jen.ID("Page").Op(":").Lit(1), jen.ID("Limit").Op(":").Lit(20), jen.ID("FilteredCount").Op(":").ID("exampleQuantity").Op("/").Lit(2), jen.ID("TotalCount").Op(":").ID("exampleQuantity")), jen.ID("AccountUserMemberships").Op(":").ID("examples")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountUserMembershipUpdateInputFromAccountUserMembership builds a faked AccountUserMembershipUpdateInput from an item."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUserMembershipUpdateInputFromAccountUserMembership").Params(jen.ID("accountUserMembership").Op("*").ID("types").Dot("AccountUserMembership")).Params(jen.Op("*").ID("types").Dot("AccountUserMembershipUpdateInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountUserMembershipUpdateInput").Valuesln(jen.ID("BelongsToUser").Op(":").ID("accountUserMembership").Dot("BelongsToUser"), jen.ID("BelongsToAccount").Op(":").ID("accountUserMembership").Dot("BelongsToAccount"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountUserMembershipCreationInput builds a faked AccountUserMembershipCreationInput."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUserMembershipCreationInput").Params().Params(jen.Op("*").ID("types").Dot("AccountUserMembershipCreationInput")).Body(
			jen.Return().ID("BuildFakeAccountUserMembershipCreationInputFromAccountUserMembership").Call(jen.ID("BuildFakeAccountUserMembership").Call())),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildFakeAccountUserMembershipCreationInputFromAccountUserMembership builds a faked AccountUserMembershipCreationInput from an item."),
		jen.Line(),
		jen.Func().ID("BuildFakeAccountUserMembershipCreationInputFromAccountUserMembership").Params(jen.ID("accountUserMembership").Op("*").ID("types").Dot("AccountUserMembership")).Params(jen.Op("*").ID("types").Dot("AccountUserMembershipCreationInput")).Body(
			jen.Return().Op("&").ID("types").Dot("AccountUserMembershipCreationInput").Valuesln(jen.ID("BelongsToUser").Op(":").ID("accountUserMembership").Dot("BelongsToUser"), jen.ID("BelongsToAccount").Op(":").ID("accountUserMembership").Dot("BelongsToAccount"))),
		jen.Line(),
	)

	return code
}
