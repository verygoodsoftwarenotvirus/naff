package types

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("PaidAccountBillingStatus").ID("AccountBillingStatus").Op("=").Lit("paid"),
			jen.ID("UnpaidAccountBillingStatus").ID("AccountBillingStatus").Op("=").Lit("unpaid"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("AccountBillingStatus").ID("string"),
			jen.ID("Account").Struct(
				jen.ID("ArchivedOn").Op("*").ID("uint64"),
				jen.ID("SubscriptionPlanID").Op("*").ID("uint64"),
				jen.ID("LastUpdatedOn").Op("*").ID("uint64"),
				jen.ID("Name").ID("string"),
				jen.ID("BillingStatus").ID("AccountBillingStatus"),
				jen.ID("ContactEmail").ID("string"),
				jen.ID("ContactPhone").ID("string"),
				jen.ID("PaymentProcessorCustomerID").ID("string"),
				jen.ID("ExternalID").ID("string"),
				jen.ID("Members").Index().Op("*").ID("AccountUserMembership"),
				jen.ID("CreatedOn").ID("uint64"),
				jen.ID("ID").ID("uint64"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
			jen.ID("AccountList").Struct(
				jen.ID("Accounts").Index().Op("*").ID("Account"),
				jen.ID("Pagination"),
			),
			jen.ID("AccountCreationInput").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("ContactEmail").ID("string"),
				jen.ID("ContactPhone").ID("string"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
			jen.ID("AccountUpdateInput").Struct(
				jen.ID("Name").ID("string"),
				jen.ID("ContactEmail").ID("string"),
				jen.ID("ContactPhone").ID("string"),
				jen.ID("BelongsToUser").ID("uint64"),
			),
			jen.ID("AccountDataManager").Interface(
				jen.ID("GetAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("Account"), jen.ID("error")),
				jen.ID("GetAllAccountsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")),
				jen.ID("GetAllAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("resultChannel").Chan().Index().Op("*").ID("Account"), jen.ID("bucketSize").ID("uint16")).Params(jen.ID("error")),
				jen.ID("GetAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("AccountList"), jen.ID("error")),
				jen.ID("GetAccountsForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("QueryFilter")).Params(jen.Op("*").ID("AccountList"), jen.ID("error")),
				jen.ID("CreateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("AccountCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("Account"), jen.ID("error")),
				jen.ID("UpdateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("Account"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("FieldChangeSummary")).Params(jen.ID("error")),
				jen.ID("ArchiveAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID"), jen.ID("archivedByUser")).ID("uint64")).Params(jen.ID("error")),
				jen.ID("GetAuditLogEntriesForAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.Index().Op("*").ID("AuditLogEntry"), jen.ID("error")),
			),
			jen.ID("AccountDataService").Interface(
				jen.ID("ListHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("CreateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ReadHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("UpdateHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ArchiveHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("AddMemberHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("RemoveMemberHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("MarkAsDefaultAccountHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ModifyMemberPermissionsHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("TransferAccountOwnershipHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("AuditEntryHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Update merges an AccountUpdateInput with an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("Account")).ID("Update").Params(jen.ID("input").Op("*").ID("AccountUpdateInput")).Params(jen.Index().Op("*").ID("FieldChangeSummary")).Body(
			jen.Var().Defs(
				jen.ID("out").Index().Op("*").ID("FieldChangeSummary"),
			),
			jen.If(jen.ID("input").Dot("Name").Op("!=").Lit("").Op("&&").ID("input").Dot("Name").Op("!=").ID("x").Dot("Name")).Body(
				jen.ID("out").Op("=").ID("append").Call(
					jen.ID("out"),
					jen.Op("&").ID("FieldChangeSummary").Valuesln(jen.ID("FieldName").Op(":").Lit("Name"), jen.ID("OldValue").Op(":").ID("x").Dot("Name"), jen.ID("NewValue").Op(":").ID("input").Dot("Name")),
				),
				jen.ID("x").Dot("Name").Op("=").ID("input").Dot("Name"),
			),
			jen.Return().ID("out"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("AccountCreationInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a AccountUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("AccountUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("AccountUpdateInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext validates a AccountUpdateInput."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("AccountUpdateInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("x"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("x").Dot("Name"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AccountCreationInputForNewUser creates a new AccountInputCreation struct for a given user."),
		jen.Line(),
		jen.Func().ID("AccountCreationInputForNewUser").Params(jen.ID("u").Op("*").ID("User")).Params(jen.Op("*").ID("AccountCreationInput")).Body(
			jen.Return().Op("&").ID("AccountCreationInput").Valuesln(jen.ID("Name").Op(":").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s_default"),
				jen.ID("u").Dot("Username"),
			), jen.ID("BelongsToUser").Op(":").ID("u").Dot("ID"))),
		jen.Line(),
	)

	return code
}
