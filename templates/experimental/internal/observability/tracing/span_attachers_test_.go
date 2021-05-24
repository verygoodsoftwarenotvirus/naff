package tracing

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanAttachersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_attachUint8ToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("attachUint8ToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit(1),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_attachUint64ToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("attachUint64ToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_attachStringToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("attachStringToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_attachBooleanToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("attachBooleanToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.ID("false"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("Test_attachSliceToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Index().ID("string").Valuesln(jen.ID("t").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachFilterToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachFilterToSpan").Call(
						jen.ID("span"),
						jen.Lit(1),
						jen.Lit(2),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachAuditLogEntryIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachAuditLogEntryIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachAuditLogEntryEventTypeToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachAuditLogEntryEventTypeToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachAccountIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachAccountIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachRequestingUserIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachRequestingUserIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachChangeSummarySpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachChangeSummarySpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Index().Op("*").ID("types").Dot("FieldChangeSummary").Valuesln(jen.Valuesln(jen.ID("OldValue").Op(":").Lit("blah"), jen.ID("NewValue").Op(":").Lit("butt"))),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachSessionContextDataToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachSessionContextDataToSpan").Call(
						jen.ID("span"),
						jen.Op("&").ID("types").Dot("SessionContextData").Valuesln(jen.ID("AccountPermissions").Op(":").ID("nil"), jen.ID("Requester").Op(":").ID("types").Dot("RequesterInfo").Valuesln(jen.ID("ServicePermissions").Op(":").ID("authorization").Dot("NewServiceRolePermissionChecker").Call(jen.ID("authorization").Dot("ServiceUserRole").Dot("String").Call())), jen.ID("ActiveAccountID").Op(":").Lit(0)),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachAPIClientDatabaseIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachAPIClientDatabaseIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachAPIClientClientIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachAPIClientClientIDToSpan").Call(
						jen.ID("span"),
						jen.Lit("123"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachUserToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachUserToSpan").Call(
						jen.ID("span"),
						jen.ID("exampleUser"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachUserIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachUserIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachUsernameToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachUsernameToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachWebhookIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachWebhookIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachURLToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.List(jen.ID("u"), jen.ID("err")).Op(":=").Qual("net/url", "ParseRequestURI").Call(jen.Lit("https://todo.verygoodsoftwarenotvirus.ru")),
					jen.ID("assert").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("AttachURLToSpan").Call(
						jen.ID("span"),
						jen.ID("u"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachRequestURIToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachRequestURIToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachRequestToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.List(jen.ID("req"), jen.ID("err")).Op(":=").Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
					jen.ID("require").Dot("NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.ID("AttachRequestToSpan").Call(
						jen.ID("span"),
						jen.ID("req"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachResponseToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("res").Op(":=").Op("&").Qual("net/http", "Response").Valuesln(jen.ID("Header").Op(":").Map(jen.ID("string")).Index().ID("string").Valuesln()),
					jen.ID("res").Dot("Header").Dot("Set").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
					jen.ID("AttachResponseToSpan").Call(
						jen.ID("span"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachErrorToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachErrorToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachDatabaseQueryToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachDatabaseQueryToSpan").Call(
						jen.ID("span"),
						jen.Lit("description"),
						jen.Lit("query"),
						jen.Index().Interface().Valuesln(jen.Lit("blah")),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachQueryFilterToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachQueryFilterToSpan").Call(
						jen.ID("span"),
						jen.ID("types").Dot("DefaultQueryFilter").Call(),
					),
				),
			),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachQueryFilterToSpan").Call(
						jen.ID("span"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachSearchQueryToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachSearchQueryToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachUserAgentDataToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachUserAgentDataToSpan").Call(
						jen.ID("span"),
						jen.Op("&").Qual("github.com/mssola/user_agent", "UserAgent").Valuesln(),
					),
				),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Func().ID("TestAttachItemIDToSpan").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("AttachItemIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Line(),
	)

	return code
}
