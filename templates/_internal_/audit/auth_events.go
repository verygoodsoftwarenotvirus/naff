package audit

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authEventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("ActorAssignmentKey").Op("=").Lit("performed_by"),
			jen.ID("ChangesAssignmentKey").Op("=").Lit("changes"),
			jen.ID("CreationAssignmentKey").Op("=").Lit("created"),
			jen.ID("AccountRolesKey").Op("=").Lit("account_roles"),
			jen.ID("PermissionsKey").Op("=").Lit("permissions"),
			jen.ID("ReasonKey").Op("=").Lit("reason"),
			jen.ID("UserBannedEvent").Op("=").Lit("user_banned"),
			jen.ID("AccountTerminatedEvent").Op("=").Lit("account_terminated"),
			jen.ID("CycleCookieSecretEvent").Op("=").Lit("cookie_secret_cycled"),
			jen.ID("SuccessfulLoginEvent").Op("=").Lit("user_logged_in"),
			jen.ID("LogoutEvent").Op("=").Lit("user_logged_out"),
			jen.ID("BannedUserLoginAttemptEvent").Op("=").Lit("banned_user_login_attempt"),
			jen.ID("UnsuccessfulLoginBadPasswordEvent").Op("=").Lit("user_login_failed_bad_password"),
			jen.ID("UnsuccessfulLoginBad2FATokenEvent").Op("=").Lit("user_login_failed_bad_2FA_token"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildCycleCookieSecretEvent builds an entry creation input for when a cookie secret is cycled."),
		jen.Line(),
		jen.Func().ID("BuildCycleCookieSecretEvent").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("CycleCookieSecretEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildSuccessfulLoginEventEntry builds an entry creation input for when a user successfully logs in."),
		jen.Line(),
		jen.Func().ID("BuildSuccessfulLoginEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("SuccessfulLoginEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildBannedUserLoginAttemptEventEntry builds an entry creation input for when a user successfully logs in."),
		jen.Line(),
		jen.Func().ID("BuildBannedUserLoginAttemptEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("BannedUserLoginAttemptEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUnsuccessfulLoginBadPasswordEventEntry builds an entry creation input for when a user fails to log in because of a bad passwords."),
		jen.Line(),
		jen.Func().ID("BuildUnsuccessfulLoginBadPasswordEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UnsuccessfulLoginBadPasswordEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUnsuccessfulLoginBad2FATokenEventEntry builds an entry creation input for when a user fails to log in because of a bad two factor token."),
		jen.Line(),
		jen.Func().ID("BuildUnsuccessfulLoginBad2FATokenEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UnsuccessfulLoginBad2FATokenEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildLogoutEventEntry builds an entry creation input for when a user logs out."),
		jen.Line(),
		jen.Func().ID("BuildLogoutEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("LogoutEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
