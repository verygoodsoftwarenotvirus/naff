package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func usersDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("UserDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("serviceRolesSeparator").Op("=").Lit(","),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanUser provides a consistent way to scan something like a *sql.Row into a Requester struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("user").Op("*").ID("types").Dot("User"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.ID("user").Op("=").Op("&").ID("types").Dot("User").Values(),
			jen.Var().Defs(
				jen.ID("rawRoles").ID("string"),
			),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("user").Dot("ID"), jen.Op("&").ID("user").Dot("ExternalID"), jen.Op("&").ID("user").Dot("Username"), jen.Op("&").ID("user").Dot("AvatarSrc"), jen.Op("&").ID("user").Dot("HashedPassword"), jen.Op("&").ID("user").Dot("RequiresPasswordChange"), jen.Op("&").ID("user").Dot("PasswordLastChangedOn"), jen.Op("&").ID("user").Dot("TwoFactorSecret"), jen.Op("&").ID("user").Dot("TwoFactorSecretVerifiedOn"), jen.Op("&").ID("rawRoles"), jen.Op("&").ID("user").Dot("ServiceAccountStatus"), jen.Op("&").ID("user").Dot("ReputationExplanation"), jen.Op("&").ID("user").Dot("CreatedOn"), jen.Op("&").ID("user").Dot("LastUpdatedOn"), jen.Op("&").ID("user").Dot("ArchivedOn")),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Op("=").ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("filteredCount"),
					jen.Op("&").ID("totalCount"),
				)),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				))),
			jen.If(jen.ID("roles").Op(":=").Qual("strings", "Split").Call(
				jen.ID("rawRoles"),
				jen.ID("serviceRolesSeparator"),
			), jen.ID("len").Call(jen.ID("roles")).Op(">").Lit(0)).Body(
				jen.ID("user").Dot("ServiceRoles").Op("=").ID("roles")).Else().Body(
				jen.ID("user").Dot("ServiceRoles").Op("=").Index().ID("string").Values()),
			jen.Return().List(jen.ID("user"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanUsers takes database rows and loads them into a slice of Requester structs."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("users").Index().Op("*").ID("types").Dot("User"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("user"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanUser").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("scanErr"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("scanning user result"),
					))),
				jen.If(jen.ID("includeCounts").Op("&&").ID("filteredCount").Op("==").Lit(0)).Body(
					jen.ID("filteredCount").Op("=").ID("fc")),
				jen.If(jen.ID("includeCounts").Op("&&").ID("totalCount").Op("==").Lit(0)).Body(
					jen.ID("totalCount").Op("=").ID("tc")),
				jen.ID("users").Op("=").ID("append").Call(
					jen.ID("users"),
					jen.ID("user"),
				),
			),
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
			jen.Return().List(jen.ID("users"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("getUser fetches a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("getUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("withVerifiedTOTPSecret").ID("bool")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Var().Defs(
				jen.ID("query").ID("string"),
				jen.ID("args").Index().Interface(),
			),
			jen.If(jen.ID("withVerifiedTOTPSecret")).Body(
				jen.List(jen.ID("query"), jen.ID("args")).Op("=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetUserQuery").Call(
					jen.ID("ctx"),
					jen.ID("userID"),
				)).Else().Body(
				jen.List(jen.ID("query"), jen.ID("args")).Op("=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetUserWithUnverifiedTwoFactorSecretQuery").Call(
					jen.ID("ctx"),
					jen.ID("userID"),
				)),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("user"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.List(jen.ID("u"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanUser").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				))),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("createUser creates a user. The `user` and `account` parameters are meant to be filled out."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("createUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("user").Op("*").ID("types").Dot("User"), jen.ID("account").Op("*").ID("types").Dot("Account"), jen.ID("userCreationQuery").ID("string"), jen.ID("userCreationArgs").Index().Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("username"),
				jen.ID("user").Dot("Username"),
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
			jen.List(jen.ID("userID"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("false"),
				jen.Lit("user creation"),
				jen.ID("userCreationQuery"),
				jen.ID("userCreationArgs"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating user"),
				),
			),
			jen.ID("user").Dot("ID").Op("=").ID("userID"),
			jen.ID("account").Dot("BelongsToUser").Op("=").ID("user").Dot("ID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserCreationEventEntry").Call(jen.ID("user").Dot("ID")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing user creation audit log entry"),
				),
			),
			jen.ID("accountCreationInput").Op(":=").ID("types").Dot("AccountCreationInputForNewUser").Call(jen.ID("user")),
			jen.List(jen.ID("accountCreationQuery"), jen.ID("accountCreationArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildAccountCreationQuery").Call(
				jen.ID("ctx"),
				jen.ID("accountCreationInput"),
			),
			jen.List(jen.ID("accountID"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
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
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("create account"),
				),
			),
			jen.ID("account").Dot("ID").Op("=").ID("accountID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildAccountCreationEventEntry").Call(
					jen.ID("account"),
					jen.ID("user").Dot("ID"),
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
					jen.Lit("writing account creation audit log entry"),
				),
			),
			jen.List(jen.ID("addUserToAccountQuery"), jen.ID("addUserToAccountArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildCreateMembershipForNewUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
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
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing account user membership creation audit log entry"),
				),
			),
			jen.ID("addToAccountInput").Op(":=").Op("&").ID("types").Dot("AddUserToAccountInput").Valuesln(jen.ID("UserID").Op(":").ID("user").Dot("ID"), jen.ID("AccountID").Op(":").ID("account").Dot("ID"), jen.ID("AccountRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("AccountMemberRole").Dot("String").Call()), jen.ID("Reason").Op(":").Lit("account creation")),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserAddedToAccountEventEntry").Call(
					jen.ID("userID"),
					jen.ID("addToAccountInput"),
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
					jen.Lit("writing user added to account audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("user").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("account").Dot("ID"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user and account created")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserHasStatus fetches whether an item exists from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UserHasStatus").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("statuses").Op("...").ID("string")).Params(jen.ID("banned").ID("bool"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrInvalidIDProvided"))),
			jen.If(jen.ID("len").Call(jen.ID("statuses")).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("true"), jen.ID("nil"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			).Dot("WithValue").Call(
				jen.Lit("statuses"),
				jen.ID("statuses"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUserHasStatusQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("statuses").Op("..."),
			),
			jen.List(jen.ID("result"), jen.ID("err")).Op(":=").ID("q").Dot("performBooleanQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("performing user status check"),
				))),
			jen.Return().List(jen.ID("result"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUser fetches a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Return().ID("q").Dot("getUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("true"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified 2FA secret."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetUserWithUnverifiedTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.Return().ID("q").Dot("getUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("false"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUserByUsername fetches a user by their username."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetUserByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("username").ID("string")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("username").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrEmptyInputProvided"))),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("username"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("username"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetUserByUsernameQuery").Call(
				jen.ID("ctx"),
				jen.ID("username"),
			),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("user"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.List(jen.ID("u"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanUser").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("err"))),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				)),
			),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SearchForUsersByUsername fetches a list of users whose usernames begin with a given query."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("SearchForUsersByUsername").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("usernameQuery").ID("string")).Params(jen.Index().Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("usernameQuery").Op("==").Lit("")).Body(
				jen.Return().List(jen.Index().Op("*").ID("types").Dot("User").Values(), jen.ID("ErrEmptyInputProvided"))),
			jen.ID("tracing").Dot("AttachSearchQueryToSpan").Call(
				jen.ID("span"),
				jen.ID("usernameQuery"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("SearchQueryKey"),
				jen.ID("usernameQuery"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildSearchForUserByUsernameQuery").Call(
				jen.ID("ctx"),
				jen.ID("usernameQuery"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("user search by username"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("err"))),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying database for users"),
				)),
			),
			jen.List(jen.ID("u"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanUsers").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				))),
			jen.Return().List(jen.ID("u"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllUsersCount fetches a count of users from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllUsersCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAllUsersCountQuery").Call(jen.ID("ctx")),
				jen.Lit("fetching count of users"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for count of users"),
				))),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetUsers fetches a list of users from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetUsers").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dot("UserList"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("x").Op("=").Op("&").ID("types").Dot("UserList").Values(),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("logger").Op(":=").ID("filter").Dot("AttachToLogger").Call(jen.ID("q").Dot("logger")),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetUsersQuery").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("users"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user"),
				))),
			jen.If(jen.List(jen.ID("x").Dot("Users"), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dot("scanUsers").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("loading response from database"),
				))),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateUser creates a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("CreateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("UserDataStoreCreationInput")).Params(jen.Op("*").ID("types").Dot("User"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("Username"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("input").Dot("Username"),
			),
			jen.List(jen.ID("userCreationQuery"), jen.ID("userCreationArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildCreateUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.ID("user").Op(":=").Op("&").ID("types").Dot("User").Valuesln(jen.ID("Username").Op(":").ID("input").Dot("Username"), jen.ID("HashedPassword").Op(":").ID("input").Dot("HashedPassword"), jen.ID("TwoFactorSecret").Op(":").ID("input").Dot("TwoFactorSecret"), jen.ID("ServiceRoles").Op(":").Index().ID("string").Valuesln(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call()), jen.ID("CreatedOn").Op(":").ID("q").Dot("currentTime").Call()),
			jen.ID("account").Op(":=").Op("&").ID("types").Dot("Account").Valuesln(jen.ID("Name").Op(":").ID("input").Dot("Username"), jen.ID("SubscriptionPlanID").Op(":").ID("nil"), jen.ID("CreatedOn").Op(":").ID("q").Dot("currentTime").Call()),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("createUser").Call(
				jen.ID("ctx"),
				jen.ID("user"),
				jen.ID("account"),
				jen.ID("userCreationQuery"),
				jen.ID("userCreationArgs"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating user"),
				))),
			jen.Return().List(jen.ID("user"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUser receives a complete Requester struct and updates its record in the database."),
		jen.Line(),
		jen.Comment("NOTE: this function uses the ID provided in the input to make its query."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("User"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("updated").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("tracing").Dot("AttachUsernameToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("Username"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UsernameKey"),
				jen.ID("updated").Dot("Username"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUpdateUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user update"),
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
					jen.Lit("updating user"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserUpdateEventEntry").Call(
					jen.ID("updated").Dot("ID"),
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
					jen.Lit("writing user update audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user updated")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUserPassword updates a user's passwords hash in the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateUserPassword").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newHash").ID("string")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("newHash").Op("==").Lit("")).Body(
				jen.Return().ID("ErrEmptyInputProvided")),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUpdateUserPasswordQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("newHash"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user passwords update"),
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
					jen.Lit("updating user password"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserUpdatePasswordEventEntry").Call(jen.ID("userID")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing user password update audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user password updated")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateUserTwoFactorSecret marks a user's two factor secret as validated."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateUserTwoFactorSecret").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("newSecret").ID("string")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("newSecret").Op("==").Lit("")).Body(
				jen.Return().ID("ErrEmptyInputProvided")),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUpdateUserTwoFactorSecretQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("newSecret"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user 2FA secret update"),
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
					jen.Lit("updating user 2FA secret"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserUpdateTwoFactorSecretEventEntry").Call(jen.ID("userID")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing update 2FA secret audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user two factor secret updated")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MarkUserTwoFactorSecretAsVerified marks a user's two factor secret as validated."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("MarkUserTwoFactorSecretAsVerified").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildVerifyUserTwoFactorSecretQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user two factor secret verification"),
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
					jen.Lit("writing verified two factor status to database"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserVerifyTwoFactorSecretEventEntry").Call(jen.ID("userID")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing user 2FA secret verification audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user two factor secret verified")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveUser archives a user."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ArchiveUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
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
			jen.List(jen.ID("archiveUserQuery"), jen.ID("archiveUserArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildArchiveUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user archive"),
				jen.ID("archiveUserQuery"),
				jen.ID("archiveUserArgs"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving user"),
				),
			),
			jen.List(jen.ID("archiveMembershipsQuery"), jen.ID("archiveMembershipsArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildArchiveAccountMembershipsForUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user memberships archive"),
				jen.ID("archiveMembershipsQuery"),
				jen.ID("archiveMembershipsArgs"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving user account memberships"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserArchiveEventEntry").Call(jen.ID("userID")),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing user archive audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user archived")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForUser fetches a list of audit log entries from the database that relate to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAuditLogEntriesForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAuditLogEntriesForUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("audit log entries for user"),
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
					jen.Lit("scanning response from database"),
				))),
			jen.Return().List(jen.ID("auditLogEntries"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
