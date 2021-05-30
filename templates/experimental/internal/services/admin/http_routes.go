package admin

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func httpRoutesDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("UserIDURIParamKey").Op("=").Lit("userID"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UserReputationChangeHandler changes a user's status."),
		jen.Line(),
		jen.Func().Params(jen.ID("s").Op("*").ID("service")).ID("UserReputationChangeHandler").Params(jen.ID("res").Qual("net/http", "ResponseWriter"), jen.ID("req").Op("*").Qual("net/http", "Request")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("s").Dot("tracer").Dot("StartSpan").Call(jen.ID("req").Dot("Context").Call()),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("s").Dot("logger").Dot("WithRequest").Call(jen.ID("req")),
			jen.ID("tracing").Dot("AttachRequestToSpan").Call(
				jen.ID("span"),
				jen.ID("req"),
			),
			jen.List(jen.ID("sessionCtxData"), jen.ID("err")).Op(":=").ID("s").Dot("sessionContextDataFetcher").Call(jen.ID("req")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("retrieving session context data"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("input").Op(":=").ID("new").Call(jen.ID("types").Dot("UserReputationUpdateInput")),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("encoderDecoder").Dot("DecodeRequest").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("observability").Dot("AcknowledgeError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("decoding request body"),
				),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("invalid request content"),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.If(jen.ID("err").Op("=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("ValidationErrorKey"),
					jen.ID("err"),
				).Dot("Debug").Call(jen.Lit("invalid input attached to request")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.ID("err").Dot("Error").Call(),
					jen.Qual("net/http", "StatusBadRequest"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("new_status"),
				jen.ID("input").Dot("NewReputation"),
			),
			jen.ID("tracing").Dot("AttachSessionContextDataToSpan").Call(
				jen.ID("span"),
				jen.ID("sessionCtxData"),
			),
			jen.If(jen.Op("!").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("CanUpdateUserReputations").Call()).Body(
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeErrorResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
					jen.Lit("inadequate permissions for route"),
					jen.Qual("net/http", "StatusForbidden"),
				),
				jen.Return(),
			),
			jen.ID("requester").Op(":=").ID("sessionCtxData").Dot("Requester").Dot("UserID"),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("ban_giver"),
				jen.ID("requester"),
			),
			jen.Var().Defs(
				jen.ID("allowed").ID("bool"),
			),
			jen.Switch(jen.ID("input").Dot("NewReputation")).Body(
				jen.Case(jen.ID("types").Dot("BannedUserAccountStatus"), jen.ID("types").Dot("TerminatedUserReputation")).Body(
					jen.ID("allowed").Op("=").ID("sessionCtxData").Dot("Requester").Dot("ServicePermissions").Dot("CanUpdateUserReputations").Call()),
				jen.Case(jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("types").Dot("UnverifiedAccountStatus")).Body(
					jen.ID("allowed").Op("=").ID("true")),
			),
			jen.If(jen.Op("!").ID("allowed")).Body(
				jen.ID("logger").Dot("Info").Call(jen.Lit("ban attempt made by admin without appropriate permissions")),
				jen.ID("s").Dot("encoderDecoder").Dot("EncodeInvalidPermissionsResponse").Call(
					jen.ID("ctx"),
					jen.ID("res"),
				),
				jen.Return(),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.Lit("status_change_recipient"),
				jen.ID("input").Dot("TargetUserID"),
			),
			jen.If(jen.ID("err").Op("=").ID("s").Dot("userDB").Dot("UpdateUserReputation").Call(
				jen.ID("ctx"),
				jen.ID("input").Dot("TargetUserID"),
				jen.ID("input"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("err"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.ID("s").Dot("encoderDecoder").Dot("EncodeNotFoundResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					)).Else().Body(
					jen.ID("observability").Dot("AcknowledgeError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("retrieving session context data"),
					),
					jen.ID("s").Dot("encoderDecoder").Dot("EncodeUnspecifiedInternalServerErrorResponse").Call(
						jen.ID("ctx"),
						jen.ID("res"),
					),
				),
				jen.Return(),
			),
			jen.Switch(jen.ID("input").Dot("NewReputation")).Body(
				jen.Case(jen.ID("types").Dot("BannedUserAccountStatus")).Body(
					jen.ID("s").Dot("auditLog").Dot("LogUserBanEvent").Call(
						jen.ID("ctx"),
						jen.ID("requester"),
						jen.ID("input").Dot("TargetUserID"),
						jen.ID("input").Dot("Reason"),
					)),
				jen.Case(jen.ID("types").Dot("TerminatedUserReputation")).Body(
					jen.ID("s").Dot("auditLog").Dot("LogAccountTerminationEvent").Call(
						jen.ID("ctx"),
						jen.ID("requester"),
						jen.ID("input").Dot("TargetUserID"),
						jen.ID("input").Dot("Reason"),
					)),
				jen.Case(jen.ID("types").Dot("GoodStandingAccountStatus"), jen.ID("types").Dot("UnverifiedAccountStatus")).Body(),
			),
			jen.ID("s").Dot("encoderDecoder").Dot("EncodeResponseWithStatus").Call(
				jen.ID("ctx"),
				jen.ID("res"),
				jen.ID("nil"),
				jen.Qual("net/http", "StatusAccepted"),
			),
		),
		jen.Line(),
	)

	return code
}
