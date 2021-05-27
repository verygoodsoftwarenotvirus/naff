package querier

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountUserMembershipsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("AccountUserMembershipDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("accountMemberRolesSeparator").Op("=").Lit(","),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAccountUserMembership takes a database Scanner (i.e. *sql.Row) and scans the result into an AccountUserMembership struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAccountUserMembership").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner")).Params(jen.ID("x").Op("*").ID("types").Dot("AccountUserMembership"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("x").Op("=").Op("&").ID("types").Dot("AccountUserMembership").Values(),
			jen.Var().Defs(
				jen.ID("rawAccountRoles").ID("string"),
			),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("BelongsToUser"), jen.Op("&").ID("x").Dot("BelongsToAccount"), jen.Op("&").ID("rawAccountRoles"), jen.Op("&").ID("x").Dot("DefaultAccount"), jen.Op("&").ID("x").Dot("CreatedOn"), jen.Op("&").ID("x").Dot("LastUpdatedOn"), jen.Op("&").ID("x").Dot("ArchivedOn")),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("q").Dot("logger"),
					jen.ID("span"),
					jen.Lit("scanning account user memberships"),
				))),
			jen.If(jen.ID("roles").Op(":=").Qual("strings", "Split").Call(
				jen.ID("rawAccountRoles"),
				jen.ID("accountMemberRolesSeparator"),
			), jen.ID("len").Call(jen.ID("roles")).Op(">").Lit(0)).Body(
				jen.ID("x").Dot("AccountRoles").Op("=").ID("roles")).Else().Body(
				jen.ID("x").Dot("AccountRoles").Op("=").Index().ID("string").Values()),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanAccountUserMemberships takes some database rows and turns them into a slice of memberships."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanAccountUserMemberships").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator")).Params(jen.ID("defaultAccount").ID("uint64"), jen.ID("accountRolesMap").Map(jen.ID("uint64")).Index().ID("string"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("accountRolesMap").Op("=").Map(jen.ID("uint64")).Index().ID("string").Values(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("x"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanAccountUserMembership").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.Lit(0), jen.ID("nil"), jen.ID("scanErr"))),
				jen.If(jen.ID("x").Dot("DefaultAccount").Op("&&").ID("defaultAccount").Op("==").Lit(0)).Body(
					jen.ID("defaultAccount").Op("=").ID("x").Dot("BelongsToAccount")),
				jen.ID("accountRolesMap").Index(jen.ID("x").Dot("BelongsToAccount")).Op("=").ID("x").Dot("AccountRoles"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("handling rows"),
				))),
			jen.Return().List(jen.ID("defaultAccount"), jen.ID("accountRolesMap"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildSessionContextDataForUser does ."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("BuildSessionContextDataForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("SessionContextData"), jen.ID("error")).Body(
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
			jen.List(jen.ID("user"), jen.ID("err")).Op(":=").ID("q").Dot("GetUser").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching user from database"),
				))),
			jen.List(jen.ID("getAccountMembershipsQuery"), jen.ID("getAccountMembershipsArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAccountMembershipsForUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("membershipRows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("account memberships for user"),
				jen.ID("getAccountMembershipsQuery"),
				jen.ID("getAccountMembershipsArgs").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching user's memberships from database"),
				))),
			jen.List(jen.ID("defaultAccountID"), jen.ID("accountRolesMap"), jen.ID("err")).Op(":=").ID("q").Dot("scanAccountUserMemberships").Call(
				jen.ID("ctx"),
				jen.ID("membershipRows"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning user's memberships from database"),
				))),
			jen.ID("actualAccountRolesMap").Op(":=").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker").Values(),
			jen.For(jen.List(jen.ID("accountID"), jen.ID("roles")).Op(":=").Range().ID("accountRolesMap")).Body(
				jen.ID("actualAccountRolesMap").Index(jen.ID("accountID")).Op("=").ID("authorization").Dot("NewAccountRolePermissionChecker").Call(jen.ID("roles").Op("..."))),
			jen.ID("sessionCtxData").Op(":=").Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("UserID").Op(":").ID("user").Dot("ID"), jen.ID("Reputation").Op(":").ID("user").Dot("ServiceAccountStatus"), jen.ID("ReputationExplanation").Op(":").ID("user").Dot("ReputationExplanation"), jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("user").Dot("ServiceRoles").Op("..."))), jen.ID("AccountPermissions").Op(":").ID("actualAccountRolesMap"), jen.ID("ActiveAccountID").Op(":").ID("defaultAccountID")),
			jen.Return().List(jen.ID("sessionCtxData"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetDefaultAccountIDForUser does ."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetDefaultAccountIDForUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.Lit(0), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetDefaultAccountIDForUserQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			),
			jen.Var().Defs(
				jen.ID("id").ID("uint64"),
			),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("default account ID query"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			).Dot("Scan").Call(jen.Op("&").ID("id")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("executing id query"),
				))),
			jen.Return().List(jen.ID("id"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MarkAccountAsUserDefault does a thing."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("MarkAccountAsUserDefault").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("changedByUser")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0).Op("||").ID("accountID").Op("==").Lit(0).Op("||").ID("changedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("changedByUser"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("changedByUser"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildMarkAccountAsUserDefaultQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user default account assignment"),
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
					jen.Lit("assigning user default account"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserMarkedAccountAsDefaultEventEntry").Call(
					jen.ID("userID"),
					jen.ID("accountID"),
					jen.ID("changedByUser"),
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
					jen.Lit("account not found for user"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("account marked as default")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserIsMemberOfAccount does a thing."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UserIsMemberOfAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID")).ID("uint64")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0).Op("||").ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUserIsMemberOfAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
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
					jen.Lit("performing user membership check query"),
				))),
			jen.Return().List(jen.ID("result"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ModifyUserPermissions does a thing."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ModifyUserPermissions").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("changedByUser")).ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("ModifyUserPermissionsInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0).Op("||").ID("userID").Op("==").Lit(0).Op("||").ID("changedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"), jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("changedByUser"), jen.Lit("new_roles").Op(":").ID("input").Dot("NewRoles"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("changedByUser"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildModifyUserPermissionsQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
				jen.ID("input").Dot("NewRoles"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user account permissions modification"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil").Op("&&").Op("!").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("database/sql", "ErrNoRows"),
			)).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("modifying user account permissions"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildModifyUserPermissionsEventEntry").Call(
					jen.ID("userID"),
					jen.ID("accountID"),
					jen.ID("changedByUser"),
					jen.ID("input").Dot("NewRoles"),
					jen.ID("input").Dot("Reason"),
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
					jen.Lit("writing user account membership permission modification audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user permissions modified")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TransferAccountOwnership does a thing."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("TransferAccountOwnership").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("transferredBy")).ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("AccountOwnershipTransferInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0).Op("||").ID("transferredBy").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("transferredBy"), jen.Lit("current_owner").Op(":").ID("input").Dot("CurrentOwner"), jen.Lit("new_owner").Op(":").ID("input").Dot("NewOwner"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("NewOwner"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("transferredBy"),
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
			jen.List(jen.ID("transferAccountOwnershipQuery"), jen.ID("transferAccountOwnershipArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildTransferAccountOwnershipQuery").Call(
				jen.ID("ctx"),
				jen.ID("input").Dot("CurrentOwner"),
				jen.ID("input").Dot("NewOwner"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user ownership transfer"),
				jen.ID("transferAccountOwnershipQuery"),
				jen.ID("transferAccountOwnershipArgs"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("transferring account to new owner"),
				),
			),
			jen.List(jen.ID("transferAccountMembershipQuery"), jen.ID("transferAccountMembershipArgs")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildTransferAccountMembershipsQuery").Call(
				jen.ID("ctx"),
				jen.ID("input").Dot("CurrentOwner"),
				jen.ID("input").Dot("NewOwner"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user memberships transfer"),
				jen.ID("transferAccountMembershipQuery"),
				jen.ID("transferAccountMembershipArgs"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("transferring account memberships"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildTransferAccountOwnershipEventEntry").Call(
					jen.ID("accountID"),
					jen.ID("transferredBy"),
					jen.ID("input"),
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
					jen.Lit("writing account ownership transfer audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("account transferred to new owner")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AddUserToAccount does a thing."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("AddUserToAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput"), jen.ID("addedByUser").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("addedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("addedByUser"), jen.ID("keys").Dot("UserIDKey").Op(":").ID("input").Dot("UserID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("input").Dot("AccountID"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("UserID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("AccountID"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("addedByUser"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildAddUserToAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user account membership creation"),
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
					jen.Lit("creating user account membership"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserAddedToAccountEventEntry").Call(
					jen.ID("addedByUser"),
					jen.ID("input"),
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
			jen.ID("logger").Dot("Info").Call(jen.Lit("user added to account")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RemoveUserFromAccount removes a user's membership to an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("RemoveUserFromAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("userID"), jen.ID("accountID"), jen.ID("removedByUser")).ID("uint64"), jen.ID("reason").ID("string")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("userID").Op("==").Lit(0).Op("||").ID("accountID").Op("==").Lit(0).Op("||").ID("removedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("reason").Op("==").Lit("")).Body(
				jen.Return().ID("ErrEmptyInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("UserIDKey").Op(":").ID("userID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("ReasonKey").Op(":").ID("reason"), jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("removedByUser"))),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("removedByUser"),
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
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildRemoveUserFromAccountQuery").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("user membership removal"),
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
					jen.Lit("removing user from account"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserRemovedFromAccountEventEntry").Call(
					jen.ID("userID"),
					jen.ID("accountID"),
					jen.ID("removedByUser"),
					jen.ID("reason"),
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
					jen.Lit("writing remove user from account audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user removed from account")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
