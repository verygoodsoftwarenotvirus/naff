package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountUserMembership").Struct(
				jen.ID("ArchivedOn").Op("*").ID("uint64"),
				jen.ID("LastUpdatedOn").Op("*").ID("uint64"),
				jen.ID("AccountRoles").Index().ID("string"),
				jen.ID("BelongsToUser").ID("uint64"),
				jen.ID("BelongsToAccount").ID("uint64"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("ID").ID("uint64"),
				jen.ID("DefaultAccount").ID("bool"),
			),
			jen.ID("AccountUserMembershipList").Struct(
				jen.ID("AccountUserMemberships").Index().Op("*").ID("AccountUserMembership"),
				jen.ID("Pagination"),
			),
			jen.ID("AccountUserMembershipCreationInput").Struct(
				jen.ID("BelongsToUser").ID("uint64"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("AccountUserMembershipUpdateInput").Struct(
				jen.ID("BelongsToUser").ID("uint64"),
				jen.ID("BelongsToAccount").ID("uint64"),
			),
			jen.ID("AddUserToAccountInput").Struct(
				jen.ID("Reason").ID("string"),
				jen.ID("AccountRoles").Index().ID("string"),
				jen.ID("UserID").ID("uint64"),
				jen.ID("AccountID").ID("uint64"),
			),
			jen.ID("AccountOwnershipTransferInput").Struct(
				jen.ID("Reason").ID("string"),
				jen.ID("CurrentOwner").ID("uint64"),
				jen.ID("NewOwner").ID("uint64"),
			),
			jen.ID("ModifyUserPermissionsInput").Struct(
				jen.ID("Reason").ID("string"),
				jen.ID("NewRoles").Index().ID("string"),
			),
			jen.ID("AccountUserMembershipDataManager").Interface(
				jen.ID("BuildSessionContextDataForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("SessionContextData"), jen.ID("error")),
				jen.ID("GetDefaultAccountIDForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("MarkAccountAsUserDefault").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("changedByUser")).ID("uint64")).Params(jen.ID("error")),
				jen.ID("UserIsMemberOfAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("bool"), jen.ID("error")),
				jen.ID("ModifyUserPermissions").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID"), jen.ID("changedByUser")).ID("uint64"), jen.ID("input").Op("*").ID("ModifyUserPermissionsInput")).Params(jen.ID("error")),
				jen.ID("TransferAccountOwnership").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("transferredBy").ID("uint64"), jen.ID("input").Op("*").ID("AccountOwnershipTransferInput")).Params(jen.ID("error")),
				jen.ID("AddUserToAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("AddUserToAccountInput"), jen.ID("addedByUser").ID("uint64")).Params(jen.ID("error")),
				jen.ID("RemoveUserFromAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("removedByUser")).ID("uint64"), jen.ID("reason").ID("string")).Params(jen.ID("error")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("AddUserToAccountInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a ModifyUserPermissionsInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("ModifyUserPermissionsInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("NewRoles"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Reason"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("AccountOwnershipTransferInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a ModifyUserPermissionsInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("ModifyUserPermissionsInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("NewRoles"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Reason"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("ModifyUserPermissionsInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a ModifyUserPermissionsInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("ModifyUserPermissionsInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("NewRoles"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Reason"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	return code
}
