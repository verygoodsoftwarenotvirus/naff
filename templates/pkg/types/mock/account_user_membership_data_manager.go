package mock

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipDataManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AccountUserMembershipDataManager").Op("=").Parens(jen.Op("*").ID("AccountUserMembershipDataManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountUserMembershipDataManager").Struct(jen.ID("mock").Dot("Mock")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildSessionContextDataForUser satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("BuildSessionContextDataForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.Op("*").ID("types").Dot("SessionContextData")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetDefaultAccountIDForUser satisfies our interface contract."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("GetDefaultAccountIDForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Get").Call(jen.Lit(0)).Assert(jen.ID("uint64")), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MarkAccountAsUserDefault implements the interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("MarkAccountAsUserDefault").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("changedByUser")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("changedByUser"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserIsMemberOfAccount implements the interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("UserIsMemberOfAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.ID("returnValues").Op(":=").ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
			),
			jen.Return().List(jen.ID("returnValues").Dot("Bool").Call(jen.Lit(0)), jen.ID("returnValues").Dot("Error").Call(jen.Lit(1))),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AddUserToAccount implements the interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("AddUserToAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput"), jen.ID("addedByUser").ID("uint64")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("input"),
				jen.ID("addedByUser"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RemoveUserFromAccount implements the interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("RemoveUserFromAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("removedByUser")).ID("uint64"), jen.ID("reason").ID("string")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("removedByUser"),
				jen.ID("reason"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ModifyUserPermissions implements the interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("ModifyUserPermissions").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("changedByUser")).ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("ModifyUserPermissionsInput")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("changedByUser"),
				jen.ID("input"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TransferAccountOwnership implements the interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AccountUserMembershipDataManager")).ID("TransferAccountOwnership").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("transferredBy")).ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("AccountOwnershipTransferInput")).Params(jen.ID("error")).Body(
			jen.Return().ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("transferredBy"),
				jen.ID("input"),
			).Dot("Error").Call(jen.Lit(0))),
		jen.Line(),
	)

	return code
}
