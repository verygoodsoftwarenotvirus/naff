package authorization

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func serviceRoleDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(jen.Null(),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("serviceAdminRoleName").Op("=").Lit("service_admin").Var().ID("serviceUserRoleName").Op("=").Lit("service_user"),

		jen.Line(),
	)
	code.Add(jen.Null().Type().ID("ServiceRole").ID("role").Type().ID("ServiceRolePermissionChecker").Interface(
		jen.ID("HasPermission").Params(jen.ID("Permission")).Params(jen.ID("bool")),
		jen.ID("AsAccountRolePermissionChecker").Params().Params(jen.ID("AccountRolePermissionsChecker")),
		jen.ID("IsServiceAdmin").Params().Params(jen.ID("bool")),
		jen.ID("CanCycleCookieSecrets").Params().Params(jen.ID("bool")),
		jen.ID("CanSeeAccountAuditLogEntries").Params().Params(jen.ID("bool")),
		jen.ID("CanSeeAPIClientAuditLogEntries").Params().Params(jen.ID("bool")),
		jen.ID("CanSeeUserAuditLogEntries").Params().Params(jen.ID("bool")),
		jen.ID("CanSeeWebhookAuditLogEntries").Params().Params(jen.ID("bool")),
		jen.ID("CanUpdateUserReputations").Params().Params(jen.ID("bool")),
		jen.ID("CanSeeUserData").Params().Params(jen.ID("bool")),
		jen.ID("CanSearchUsers").Params().Params(jen.ID("bool")),
	),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("invalidServiceRole").ID("ServiceRole").Op("=").ID("iota").Var().ID("ServiceUserRole").ID("ServiceRole").Op("=").ID("iota").Var().ID("ServiceAdminRole").ID("ServiceRole").Op("=").ID("iota"),

		jen.Line(),
	)
	code.Add(jen.Null().Var().ID("serviceUser").Op("=").ID("gorbac").Dot(
		"NewStdRole",
	).Call(jen.ID("serviceUserRoleName")).Var().ID("serviceAdmin").Op("=").ID("gorbac").Dot(
		"NewStdRole",
	).Call(jen.ID("serviceAdminRoleName")),

		jen.Line(),
	)
	code.Add(jen.Func().Params(jen.ID("r").ID("ServiceRole")).ID("String").Params().Params(jen.ID("string")).Body(jen.Switch(jen.ID("r")).Body(
		jen.Case(jen.ID("invalidServiceRole")).Body(jen.Return().Lit("INVALID_SERVICE_ROLE")),
		jen.Case(jen.ID("ServiceUserRole")).Body(jen.Return().ID("serviceUserRoleName")),
		jen.Case(jen.ID("ServiceAdminRole")).Body(jen.Return().ID("serviceAdminRoleName")),
		jen.Default().Body(jen.Return().Lit("")),
	)),

		jen.Line(),
	)
	code.Add(jen.Null().Type().ID("serviceRoleCollection").Struct(jen.ID("Roles").Index().ID("string")),

		jen.Line(),
	)
	code.Add(jen.Func().ID("init").Params().Body(jen.Qual("encoding/gob", "Register").Call(jen.ID("serviceRoleCollection").Valuesln())),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// NewServiceRolePermissionChecker returns a new checker for a set of Roles.").ID("NewServiceRolePermissionChecker").Params(jen.ID("roles").Op("...").ID("string")).Params(jen.ID("ServiceRolePermissionChecker")).Body(jen.Return().Op("&").ID("serviceRoleCollection").Valuesln(jen.ID("Roles").Op(":").ID("roles"))),

		jen.Line(),
	)
	code.Add(jen.Func().Params(jen.ID("r").ID("serviceRoleCollection")).ID("AsAccountRolePermissionChecker").Params().Params(jen.ID("AccountRolePermissionsChecker")).Body(jen.Return().ID("NewAccountRolePermissionChecker").Call(jen.ID("r").Dot(
		"Roles",
	).Op("..."))),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// HasPermission returns whether a user can do something or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("HasPermission").Params(jen.ID("p").ID("Permission")).Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("p"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// IsServiceAdmin returns if a role is an admin.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("IsServiceAdmin").Params().Params(jen.ID("bool")).Body(
		jen.For(jen.List(jen.ID("_"), jen.ID("x")).Op(":=").Range().ID("r").Dot(
			"Roles",
		)).Body(jen.If(jen.ID("x").Op("==").ID("ServiceAdminRole").Dot(
			"String",
		).Call()).Body(jen.Return().ID("true"))),
		jen.Return().ID("false"),
	),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanCycleCookieSecrets returns whether a user can cycle cookie secrets or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanCycleCookieSecrets").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("CycleCookieSecretPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeAccountAuditLogEntries returns whether a user can view account audit log entries or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanSeeAccountAuditLogEntries").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadAccountAuditLogEntriesPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeAPIClientAuditLogEntries returns whether a user can view API client audit log entries or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanSeeAPIClientAuditLogEntries").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadAPIClientAuditLogEntriesPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeUserAuditLogEntries returns whether a user can view user audit log entries or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanSeeUserAuditLogEntries").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadUserAuditLogEntriesPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeWebhookAuditLogEntries returns whether a user can view webhook audit log entries or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanSeeWebhookAuditLogEntries").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadWebhookAuditLogEntriesPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanUpdateUserReputations returns whether a user can update user reputations or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanUpdateUserReputations").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("UpdateUserStatusPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSeeUserData returns whether a user can view users or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanSeeUserData").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("ReadUserPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	code.Add(jen.Func().Comment("// CanSearchUsers returns whether a user can search for users or not.").Params(jen.ID("r").ID("serviceRoleCollection")).ID("CanSearchUsers").Params().Params(jen.ID("bool")).Body(jen.Return().ID("hasPermission").Call(
		jen.ID("SearchUserPermission"),
		jen.ID("r").Dot(
			"Roles",
		).Op("..."),
	)),

		jen.Line(),
	)
	return code
}
