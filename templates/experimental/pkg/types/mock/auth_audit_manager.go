package mock

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func authAuditManagerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("LogCycleCookieSecretEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogCycleCookieSecretEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogSuccessfulLoginEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogSuccessfulLoginEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogBannedUserLoginAttemptEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogBannedUserLoginAttemptEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogUnsuccessfulLoginBadPasswordEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogUnsuccessfulLoginBadPasswordEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogUnsuccessfulLoginBad2FATokenEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogUnsuccessfulLoginBad2FATokenEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogLogoutEvent implements our interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("m").Op("*").ID("AuditLogEntryDataManager")).ID("LogLogoutEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.ID("m").Dot("Called").Call(
				jen.ID("ctx"),
				jen.ID("userID"),
			)),
		jen.Line(),
	)

	return code
}
