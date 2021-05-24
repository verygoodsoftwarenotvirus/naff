package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func adminDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().ID("_").ID("types").Dot("AdminUserDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("UpdateUserReputation updates a user's account status.").Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateUserReputation").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64"), jen.ID("input").Op("*").ID("types").Dot("UserReputationUpdateInput")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("UserIDKey"),
				jen.ID("userID"),
			),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildSetUserStatusQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op(":=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("user status update query"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("user status update"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("user reputation updated")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogUserBanEvent saves a UserBannedEvent in the audit log table."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogUserBanEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("banGiver"), jen.ID("banRecipient")).ID("uint64"), jen.ID("reason").ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("banRecipient"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUserBanEventEntry").Call(
					jen.ID("banGiver"),
					jen.ID("banRecipient"),
					jen.ID("reason"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogAccountTerminationEvent saves a UserBannedEvent in the audit log table."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogAccountTerminationEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("terminator"), jen.ID("terminee")).ID("uint64"), jen.ID("reason").ID("string")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("terminee"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildAccountTerminationEventEntry").Call(
					jen.ID("terminator"),
					jen.ID("terminee"),
					jen.ID("reason"),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogCycleCookieSecretEvent implements our AuditLogEntryDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogCycleCookieSecretEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildCycleCookieSecretEvent").Call(jen.ID("userID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogSuccessfulLoginEvent implements our AuditLogEntryDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogSuccessfulLoginEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildSuccessfulLoginEventEntry").Call(jen.ID("userID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogBannedUserLoginAttemptEvent implements our AuditLogEntryDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogBannedUserLoginAttemptEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildBannedUserLoginAttemptEventEntry").Call(jen.ID("userID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogUnsuccessfulLoginBadPasswordEvent implements our AuditLogEntryDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogUnsuccessfulLoginBadPasswordEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUnsuccessfulLoginBadPasswordEventEntry").Call(jen.ID("userID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().Comment("LogUnsuccessfulLoginBad2FATokenEvent implements our AuditLogEntryDataManager interface.").Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogUnsuccessfulLoginBad2FATokenEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildUnsuccessfulLoginBad2FATokenEventEntry").Call(jen.ID("userID")),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("LogLogoutEvent implements our AuditLogEntryDataManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("LogLogoutEvent").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("userID").ID("uint64")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("userID"),
			),
			jen.ID("q").Dot("createAuditLogEntry").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Qual("gitlab.com/verygoodsoftwarenotvirus/todo/internal/audit", "BuildLogoutEventEntry").Call(jen.ID("userID")),
			),
		),
		jen.Line(),
	)

	return code
}
