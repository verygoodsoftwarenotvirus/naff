package types

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("SessionContextDataKey").ID("ContextKey").Op("=").Lit("session_context_data"),
			jen.ID("UserIDContextKey").ID("ContextKey").Op("=").Lit("user_id"),
			jen.ID("AccountIDContextKey").ID("ContextKey").Op("=").Lit("account_id"),
			jen.ID("UserLoginInputContextKey").ID("ContextKey").Op("=").Lit("user_login_input"),
			jen.ID("UserRegistrationInputContextKey").ID("ContextKey").Op("=").Lit("user_registration_input"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("init").Params().Body(
			jen.Qual("encoding/gob", "Register").Call(jen.Op("&").ID("SessionContextData").Valuesln())),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("UserAccountMembershipInfo").Struct(
				jen.ID("AccountName").ID("string"),
				jen.ID("AccountRoles").Index().ID("string"),
				jen.ID("AccountID").ID("uint64"),
			),
			jen.ID("SessionContextData").Struct(
				jen.ID("AccountPermissions").Map(jen.ID("uint64")).ID("authorization").Dot("AccountRolePermissionsChecker"),
				jen.ID("Requester").ID("RequesterInfo"),
				jen.ID("ActiveAccountID").ID("uint64"),
			),
			jen.ID("RequesterInfo").Struct(
				jen.ID("ServicePermissions").ID("authorization").Dot("ServiceRolePermissionChecker"),
				jen.ID("Reputation").ID("accountStatus"),
				jen.ID("ReputationExplanation").ID("string"),
				jen.ID("UserID").ID("uint64"),
			),
			jen.ID("UserStatusResponse").Struct(
				jen.ID("UserReputation").ID("accountStatus"),
				jen.ID("UserReputationExplanation").ID("string"),
				jen.ID("ActiveAccount").ID("uint64"),
				jen.ID("UserIsAuthenticated").ID("bool"),
			),
			jen.ID("ChangeActiveAccountInput").Struct(jen.ID("AccountID").ID("uint64")),
			jen.ID("PASETOCreationInput").Struct(
				jen.ID("ClientID").ID("string"),
				jen.ID("AccountID").ID("uint64"),
				jen.ID("RequestTime").ID("int64"),
				jen.ID("RequestedLifetime").ID("uint64"),
			),
			jen.ID("PASETOResponse").Struct(
				jen.ID("Token").ID("string"),
				jen.ID("ExpiresAt").ID("string"),
			),
			jen.ID("AuthService").Interface(
				jen.ID("StatusHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("BeginSessionHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("EndSessionHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("CycleCookieSecretHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("PASETOHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("ChangeActiveAccountHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")),
				jen.ID("PermissionFilterMiddleware").Params(jen.ID("permissions").Op("...").ID("authorization").Dot("Permission")).Params(jen.Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler"))),
				jen.ID("CookieRequirementMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("UserAttributionMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("AuthorizationMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("ServiceAdminMiddleware").Params(jen.ID("next").Qual("net/http", "Handler")).Params(jen.Qual("net/http", "Handler")),
				jen.ID("AuthenticateUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("loginData").Op("*").ID("UserLoginInput")).Params(jen.Op("*").ID("User"), jen.Op("*").Qual("net/http", "Cookie"), jen.ID("error")),
				jen.ID("LogoutUser").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("sessionCtxData").Op("*").ID("SessionContextData"), jen.ID("req").Op("*").Qual("net/http", "Request"), jen.ID("res").Qual("net/http", "ResponseWriter")).Params(jen.ID("error")),
			),
			jen.ID("AuthAuditManager").Interface(
				jen.ID("LogCycleCookieSecretEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")),
				jen.ID("LogSuccessfulLoginEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")),
				jen.ID("LogBannedUserLoginAttemptEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")),
				jen.ID("LogUnsuccessfulLoginBadPasswordEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")),
				jen.ID("LogUnsuccessfulLoginBad2FATokenEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")),
				jen.ID("LogLogoutEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("ChangeActiveAccountInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our  provided UserLoginInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("PASETOCreationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("ClientID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("RequestTime"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidatableWithContext").Op("=").Parens(jen.Op("*").ID("PASETOCreationInput")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ValidateWithContext ensures our  provided UserLoginInput meets expectations."),
		jen.Line(),
		jen.Func().Params(jen.ID("i").Op("*").ID("PASETOCreationInput")).ID("ValidateWithContext").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("error")).Body(
			jen.Return().Qual("github.com/go-ozzo/ozzo-validation/v4", "ValidateStructWithContext").Call(
				jen.ID("ctx"),
				jen.ID("i"),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("ClientID"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
				jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Field").Call(
					jen.Op("&").ID("i").Dot("RequestTime"),
					jen.Qual("github.com/go-ozzo/ozzo-validation/v4", "Required"),
				),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("AccountRolePermissionsChecker returns the relevant AccountRolePermissionsChecker."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("SessionContextData")).ID("AccountRolePermissionsChecker").Params().Params(jen.ID("authorization").Dot("AccountRolePermissionsChecker")).Body(
			jen.Return().ID("x").Dot("AccountPermissions").Index(jen.ID("x").Dot("ActiveAccountID"))),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ServiceRolePermissionChecker returns the relevant ServiceRolePermissionChecker."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("SessionContextData")).ID("ServiceRolePermissionChecker").Params().Params(jen.ID("authorization").Dot("ServiceRolePermissionChecker")).Body(
			jen.Return().ID("x").Dot("Requester").Dot("ServicePermissions")),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ToBytes returns the gob encoded session info."),
		jen.Line(),
		jen.Func().Params(jen.ID("x").Op("*").ID("SessionContextData")).ID("ToBytes").Params().Params(jen.Index().ID("byte")).Body(
			jen.Var().Defs(
				jen.ID("b").Qual("bytes", "Buffer"),
			),
			jen.If(jen.ID("err").Op(":=").Qual("encoding/gob", "NewEncoder").Call(jen.Op("&").ID("b")).Dot("Encode").Call(jen.ID("x")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("panic").Call(jen.ID("err"))),
			jen.Return().ID("b").Dot("Bytes").Call(),
		),
		jen.Line(),
	)

	return code
}
