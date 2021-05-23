package audit

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func userEventsDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile("audit")

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("UserAssignmentKey").Op("=").Lit("user_id"),
			jen.ID("UserCreationEvent").Op("=").Lit("user_created"),
			jen.ID("UserVerifyTwoFactorSecretEvent").Op("=").Lit("user_two_factor_secret_verified"),
			jen.ID("UserUpdateTwoFactorSecretEvent").Op("=").Lit("user_two_factor_secret_changed"),
			jen.ID("UserUpdateEvent").Op("=").Lit("user_updated"),
			jen.ID("UserUpdatePasswordEvent").Op("=").Lit("user_password_updated"),
			jen.ID("UserArchiveEvent").Op("=").Lit("user_archived"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserCreationEventEntry builds an entry creation input for when a user is created."),
		jen.Line(),
		jen.Func().ID("BuildUserCreationEventEntry").Params(jen.ID("userID").ID("uint64")).Params(jen.Op("*").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput")).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(), "AuditLogEntryCreationInput").Valuesln(
				jen.ID("EventType").Op(":").ID("UserCreationEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("UserAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserVerifyTwoFactorSecretEventEntry builds an entry creation input for when a user verifies their two factor secret."),
		jen.Line(),
		jen.Func().ID("BuildUserVerifyTwoFactorSecretEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			).Valuesln(
				jen.ID("EventType").Op(":").ID("UserVerifyTwoFactorSecretEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserUpdateTwoFactorSecretEventEntry builds an entry creation input for when a user updates their two factor secret."),
		jen.Line(),
		jen.Func().ID("BuildUserUpdateTwoFactorSecretEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			).Valuesln(
				jen.ID("EventType").Op(":").ID("UserUpdateTwoFactorSecretEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserUpdatePasswordEventEntry builds an entry creation input for when a user updates their passwords."),
		jen.Line(),
		jen.Func().ID("BuildUserUpdatePasswordEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			).Valuesln(
				jen.ID("EventType").Op(":").ID("UserUpdatePasswordEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserUpdateEventEntry builds an entry creation input for when a user is updated."),
		jen.Line(),
		jen.Func().ID("BuildUserUpdateEventEntry").Params(
			jen.ID("userID").ID("uint64"),
			jen.ID("changes").Index().Op("*").Qual(proj.TypesPackage(),
				"FieldChangeSummary",
			)).Params(
			jen.Op("*").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			).Valuesln(
				jen.ID("EventType").Op(":").ID("UserUpdateEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
					jen.ID("ChangesAssignmentKey").Op(":").ID("changes"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("BuildUserArchiveEventEntry builds an entry creation input for when a user is archived."),
		jen.Line(),
		jen.Func().ID("BuildUserArchiveEventEntry").Params(
			jen.ID("userID").ID("uint64")).Params(
			jen.Op("*").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			)).Body(
			jen.Return().Op("&").Qual(proj.TypesPackage(),
				"AuditLogEntryCreationInput",
			).Valuesln(
				jen.ID("EventType").Op(":").ID("UserArchiveEvent"),
				jen.ID("Context").Op(":").Map(jen.ID("string")).Interface().Valuesln(
					jen.ID("UserAssignmentKey").Op(":").ID("userID"),
					jen.ID("ActorAssignmentKey").Op(":").ID("userID"),
				),
			),
		),
		jen.Line(),
	)

	return code
}
