package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AccountDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAccount takes a database Scanner (i.e. *sql.Row) and scans the result into an Account struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("account").Op("*").ID("types").Dot("Account"), jen.ID("membership").Op("*").ID("types").Dot("AccountUserMembership"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.ID("account").Op("=").Op("&").ID("types").Dot("Account").Valuesln(jen.ID("Members").Op(":").Index().Op("*").ID("types").Dot("AccountUserMembership").Valuesln()),
			jen.ID("membership").Op("=").Op("&").ID("types").Dot("AccountUserMembership").Valuesln(),
			jen.Var().Defs(
				jen.ID("rawRoles").ID("string"),
			),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("account").Dot("ID"), jen.Op("&").ID("account").Dot("ExternalID"), jen.Op("&").ID("account").Dot("Name"), jen.Op("&").ID("account").Dot("BillingStatus"), jen.Op("&").ID("account").Dot("ContactEmail"), jen.Op("&").ID("account").Dot("ContactPhone"), jen.Op("&").ID("account").Dot("PaymentProcessorCustomerID"), jen.Op("&").ID("account").Dot("SubscriptionPlanID"), jen.Op("&").ID("account").Dot("CreatedOn"), jen.Op("&").ID("account").Dot("LastUpdatedOn"), jen.Op("&").ID("account").Dot("ArchivedOn"), jen.Op("&").ID("account").Dot("BelongsToUser"), jen.Op("&").ID("membership").Dot("ID"), jen.Op("&").ID("membership").Dot("BelongsToUser"), jen.Op("&").ID("membership").Dot("BelongsToAccount"), jen.Op("&").ID("rawRoles"), jen.Op("&").ID("membership").Dot("DefaultAccount"), jen.Op("&").ID("membership").Dot("CreatedOn"), jen.Op("&").ID("membership").Dot("LastUpdatedOn"), jen.Op("&").ID("membership").Dot("ArchivedOn")),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Op("=").ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("filteredCount"),
					jen.Op("&").ID("totalCount"),
				)),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching memberships from database"),
				))),
			jen.ID("membership").Dot("AccountRoles").Op("=").Qual("strings", "Split").Call(
				jen.ID("rawRoles"),
				jen.ID("accountMemberRolesSeparator"),
			),
			jen.Return().List(jen.ID("account"), jen.ID("membership"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAccounts takes some database rows and turns them into a slice of accounts."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("accounts").Index().Op("*").ID("types").Dot("Account"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.ID("accounts").Op("=").Index().Op("*").ID("types").Dot("Account").Valuesln(),
			jen.Var().Defs(
				jen.ID("currentAccount").Op("*").ID("types").Dot("Account"),
			),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("account"), jen.ID("membership"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAccount").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("scanErr"))),
				jen.If(jen.ID("currentAccount").Op("==").ID("nil")).Body(
					jen.ID("currentAccount").Op("=").ID("account")),
				jen.If(jen.ID("currentAccount").Dot("ID").Op("!=").ID("account").Dot("ID")).Body(
					jen.ID("accounts").Op("=").ID("append").Call(
						jen.ID("accounts"),
						jen.ID("currentAccount"),
					),
					jen.ID("currentAccount").Op("=").ID("account"),
				),
				jen.ID("currentAccount").Dot("Members").Op("=").ID("append").Call(
					jen.ID("currentAccount").Dot("Members"),
					jen.ID("membership"),
				),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("filteredCount").Op("==").Lit(0)).Body(
						jen.ID("filteredCount").Op("=").ID("fc")),
					jen.If(jen.ID("totalCount").Op("==").Lit(0)).Body(
						jen.ID("totalCount").Op("=").ID("tc")),
				),
			),
			jen.If(jen.ID("currentAccount").Op("!=").ID("nil")).Body(
				jen.ID("accounts").Op("=").ID("append").Call(
					jen.ID("accounts"),
					jen.ID("currentAccount"),
				)),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("handling rows"),
				))),
			jen.Return().List(jen.ID("accounts"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccount fetches an account from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64")).Params(jen.Op("*").ID("types").Dot("Account"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0).Op("||").ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("account"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing accounts list retrieval query"),
				))),
			jen.List(jen.ID("accounts"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAccounts").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				))),
			jen.Var().Defs(
				jen.ID("account").Op("*").ID("types").Dot("Account"),
			),
			jen.If(jen.ID("len").Call(jen.ID("accounts")).Op(">").Lit(0)).Body(
				jen.ID("account").Op("=").ID("accounts").Index(jen.Lit(0))),
			jen.If(jen.ID("account").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("database/sql", "ErrNoRows"))),
			jen.Return().List(jen.ID("account"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAccountsCount fetches the count of accounts from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllAccountsCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAllAccountsCountQuery").Call(jen.ID("ctx")),
				jen.Lit("fetching count of all accounts"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for count of accounts"),
				))),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllAccounts fetches a list of all accounts in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").ID("types").Dot("Account"), jen.ID("batchSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("results").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.If(jen.ID("batchSize").Op("==").Lit(0)).Body(
				jen.ID("batchSize").Op("=").ID("defaultBatchSize")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("batch_size"),
				jen.ID("batchSize"),
			),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("GetAllAccountsCount").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching count of accounts"),
				)),
			jen.For(jen.ID("beginID").Op(":=").ID("uint64").Call(jen.Lit(1)), jen.ID("beginID").Op("<=").ID("count"), jen.ID("beginID").Op("+=").ID("uint64").Call(jen.ID("batchSize"))).Body(
				jen.ID("endID").Op(":=").ID("beginID").Op("+").ID("uint64").Call(jen.ID("batchSize")),
				jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).ID("uint64")).Body(
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetBatchOfAccountsQuery").Call(
						jen.ID("ctx"),
						jen.ID("begin"),
						jen.ID("end"),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("query").Op(":").ID("query"), jen.Lit("begin").Op(":").ID("begin"), jen.Lit("end").Op(":").ID("end"))),
					jen.List(jen.ID("rows"), jen.ID("queryErr")).Op(":=").ID("q").Dot("db").Dot("Query").Call(
						jen.ID("query"),
						jen.ID("args").Op("..."),
					),
					jen.If(jen.Qual("errors", "Is").Call(
						jen.ID("queryErr"),
						jen.Qual("database/sql", "ErrNoRows"),
					)).Body(
						jen.Return()).Else().If(jen.ID("queryErr").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("queryErr"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("querying for database rows"),
						),
						jen.Return(),
					),
					jen.List(jen.ID("accounts"), jen.ID("_"), jen.ID("_"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAccounts").Call(
						jen.ID("ctx"),
						jen.ID("rows"),
						jen.ID("false"),
					),
					jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
						jen.ID("observability").Dot("AcknowledgeError").Call(
							jen.ID("scanErr"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("scanning database rows"),
						),
						jen.Return(),
					),
					jen.ID("results").ReceiveFromChannel().ID("accounts"),
				).Call(
					jen.ID("beginID"),
					jen.ID("endID"),
				),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccounts fetches a list of accounts from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dot("AccountList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("q").Dot("logger")).Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("x").Op("=").Op("&").ID("types").Dot("AccountList").Valuesln(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAccountsQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("accounts"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing accounts list retrieval query"),
				))),
			jen.If(jen.List(jen.ID("x").Dot("Accounts"), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dot("scanAccounts").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning accounts from database"),
				))),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccountsForAdmin fetches a list of accounts from the database that meet a particular filter for all users."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAccountsForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dot("AccountList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("q").Dot("logger")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("x").Op("=").Op("&").ID("types").Dot("AccountList").Valuesln(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAccountsQuery").Call(
				jen.ID("ctx"),
				jen.Lit(0),
				jen.ID("true"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("accounts for admin"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying database for accounts"),
				))),
			jen.If(jen.List(jen.ID("x").Dot("Accounts"), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dot("scanAccounts").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning accounts"),
				))),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateAccount creates an account in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("CreateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AccountCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Account"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("createdByUser").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("createdByUser"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("input").Dot("BelongsToUser"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("createdByUser"),
			),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				))),
			jen.List(jen.ID("accountCreationQuery"), jen.ID("accountCreationArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildAccountCreationQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.List(jen.ID("id"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("false"),
				jen.Lit("account creation"),
				jen.ID("accountCreationQuery"),
				jen.ID("accountCreationArgs"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating account"),
				)),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("id"),
			),
			jen.ID("account").Op(":=").Op("&").ID("types").Dot("Account").Valuesln(jen.ID("ID").Op(":").ID("id"), jen.ID("Name").Op(":").ID("input").Dot("Name"), jen.ID("BelongsToUser").Op(":").ID("input").Dot("BelongsToUser"), jen.ID("BillingStatus").Op(":").ID("types").Dot("UnpaidAccountBillingStatus"), jen.ID("ContactEmail").Op(":").ID("input").Dot("ContactEmail"), jen.ID("ContactPhone").Op(":").ID("input").Dot("ContactPhone"), jen.ID("CreatedOn").Op(":").ID("q").Dot("currentTime").Call()),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildAccountCreationEventEntry").Call(
					jen.ID("account"),
					jen.ID("createdByUser"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing account creation audit log event entry"),
				)),
			),
			jen.ID("addInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("UserID").Op(":").ID("input").Dot("BelongsToUser"), jen.ID("AccountID").Op(":").ID("account").Dot("ID"), jen.ID("Reason").Op(":").Lit("account creation"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountAdminRole").Dot("String").Call())),
			jen.List(jen.ID("addUserToAccountQuery"), jen.ID("addUserToAccountArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildAddUserToAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("addInput"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("account user membership creation"),
				jen.ID("addUserToAccountQuery"),
				jen.ID("addUserToAccountArgs"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating account membership"),
				)),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildUserAddedToAccountEventEntry").Call(
					jen.ID("createdByUser"),
					jen.ID("addInput"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing account membership creation audit log event entry"),
				)),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				))),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("account").Dot("ID"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("account created")),
			jen.Return().List(jen.ID("account"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateAccount updates a particular account. Note that UpdateAccount expects the provided input to have a valid ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Account"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("changedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("updated").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("updated").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("changedByUser"),
			),
			jen.ID("tracing").Dot("AttachChangeSummarySpan").Call(
				jen.ID("span"),
				jen.Lit("account"),
				jen.ID("changes"),
			),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUpdateAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("account update"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating account"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildAccountUpdateEventEntry").Call(
					jen.ID("updated").Dot("BelongsToUser"),
					jen.ID("updated").Dot("ID"),
					jen.ID("changedByUser"),
					jen.ID("changes"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing account update audit log event entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("account updated")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveAccount archives an account from the database by its ID."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ArchiveAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID"), jen.ID("archivedByUser")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0).Op("||").ID("userID").Op("==").Lit(0).Op("||").ID("archivedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("archivedByUser"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"))),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildArchiveAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("account archive"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving account"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildAccountArchiveEventEntry").Call(
					jen.ID("userID"),
					jen.ID("accountID"),
					jen.ID("archivedByUser"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing account archive audit log event entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("account archived")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForAccount fetches a list of audit log entries from the database that relate to a given account."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAuditLogEntriesForAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAuditLogEntriesForAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("audit log entries for account"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying database for audit log entries"),
				))),
			jen.List(jen.ID("auditLogEntries"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAuditLogEntries").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning audit log entries"),
				))),
			jen.Return().List(jen.ID("auditLogEntries"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
