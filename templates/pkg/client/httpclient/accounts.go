package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func accountsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("SwitchActiveAccount will switch the account on whose behalf requests are made."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("SwitchActiveAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("c").Dot("authMethod").Op("==").ID("cookieAuthMethod")).Body(
				jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildSwitchActiveAccountRequest").Call(
					jen.ID("ctx"),
					jen.ID("accountID"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("building account switch request"),
					)),
				jen.If(jen.ID("err").Op("=").ID("c").Dot("executeAndUnmarshal").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.ID("c").Dot("authedClient"),
					jen.ID("nil"),
				), jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("executing account switch request"),
					)),
			),
			jen.ID("c").Dot("accountID").Op("=").ID("accountID"),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccount retrieves an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Account"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetAccountRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building account retrieval request"),
				))),
			jen.Var().Defs(
				jen.ID("account").Op("*").ID("types").Dot("Account"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("account"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving account"),
				))),
			jen.Return().List(jen.ID("account"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAccounts retrieves a list of accounts."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAccounts").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("AccountList"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("c").Dot("loggerWithFilter").Call(jen.ID("filter")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetAccountsRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building account list request"),
				))),
			jen.Var().Defs(
				jen.ID("accounts").Op("*").ID("types").Dot("AccountList"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("accounts"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving accounts"),
				))),
			jen.Return().List(jen.ID("accounts"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateAccount creates an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AccountCreationInput")).Params(jen.Op("*").ID("types").Dot("Account"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.Lit("account_name"),
				jen.ID("input").Dot("Name"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				))),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildCreateAccountRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building account creation request"),
				))),
			jen.Var().Defs(
				jen.ID("account").Op("*").ID("types").Dot("Account"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("account"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating account"),
				))),
			jen.Return().List(jen.ID("account"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateAccount updates an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("account").Op("*").ID("types").Dot("Account")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("account").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("account").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("account").Dot("ID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildUpdateAccountRequest").Call(
				jen.ID("ctx"),
				jen.ID("account"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building account update request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("account"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating account"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveAccount archives an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildArchiveAccountRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building account archive request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving account"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AddUserToAccount adds a user to an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("AddUserToAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("AddUserToAccountInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("input").Dot("AccountID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("input").Dot("UserID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("AccountID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("UserID"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				)),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildAddUserRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building add user to account request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("adding user to account"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("MarkAsDefault marks a given account as the default for a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("MarkAsDefault").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildMarkAsDefaultRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building mark account as default request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("marking account as default"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("RemoveUserFromAccount removes a user from an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("RemoveUserFromAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64"), jen.ID("reason").ID("string")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("reason").Op("==").Lit("")).Body(
				jen.Return().ID("ErrEmptyInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildRemoveUserRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
				jen.ID("reason"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building remove user from account request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("removing user from account"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ModifyMemberPermissions modifies a given user's permissions for a given account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ModifyMemberPermissions").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("accountID"), jen.ID("userID")).ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("ModifyUserPermissionsInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("userID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				)),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildModifyMemberPermissionsRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("userID"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building modify account member permissions request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("modifying user account permissions"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("TransferAccountOwnership transfers ownership of an account to a given user."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("TransferAccountOwnership").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("AccountOwnershipTransferInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			).Dot("WithValue").Call(
				jen.Lit("old_owner"),
				jen.ID("input").Dot("CurrentOwner"),
			).Dot("WithValue").Call(
				jen.Lit("new_owner"),
				jen.ID("input").Dot("NewOwner"),
			),
			jen.ID("tracing").Dot("AttachToSpan").Call(
				jen.ID("span"),
				jen.Lit("old_owner"),
				jen.ID("input").Dot("CurrentOwner"),
			),
			jen.ID("tracing").Dot("AttachToSpan").Call(
				jen.ID("span"),
				jen.Lit("new_owner"),
				jen.ID("input").Dot("NewOwner"),
			),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("validating input"),
				)),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildTransferAccountOwnershipRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building transfer account ownership request"),
				)),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("transferring account to user"),
				)),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogForAccount retrieves a list of audit log entries pertaining to an account."),
		jen.Line(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAuditLogForAccount").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.List(jen.ID("req"), jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetAuditLogForAccountRequest").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building fetch audit log entries for account request"),
				))),
			jen.Var().Defs(
				jen.ID("entries").Index().Op("*").ID("types").Dot("AuditLogEntry"),
			),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("entries"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching audit log entries for account"),
				))),
			jen.Return().List(jen.ID("entries"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
