package tracing

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func spanAttachersTestDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Func().ID("Test_attachUint8ToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("attachUint8ToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit(1),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("Test_attachUint64ToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("attachUint64ToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("Test_attachStringToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("attachStringToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("Test_attachBooleanToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("attachBooleanToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.ID("false"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("Test_attachSliceToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("attachSliceToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Index().String().Values(jen.ID("t").Dot("Name").Call()),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachFilterToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachFilterToSpan").Call(
						jen.ID("span"),
						jen.Lit(1),
						jen.Lit(2),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachAuditLogEntryIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachAuditLogEntryIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachAuditLogEntryEventTypeToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachAuditLogEntryEventTypeToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachAccountIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachAccountIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachRequestingUserIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachRequestingUserIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachChangeSummarySpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachChangeSummarySpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Index().PointerTo().Qual(proj.TypesPackage(), "FieldChangeSummary").Valuesln(jen.Valuesln(jen.ID("OldValue").Op(":").Lit("blah"), jen.ID("NewValue").Op(":").Lit("butt"))),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachSessionContextDataToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachSessionContextDataToSpan").Call(
						jen.ID("span"),
						jen.AddressOf().Qual(proj.TypesPackage(), "SessionContextData").Valuesln(jen.ID("AccountPermissions").Op(":").ID("nil"), jen.ID("Requester").Op(":").Qual(proj.TypesPackage(), "RequesterInfo").Valuesln(jen.ID("ServicePermissions").Op(":").Qual(proj.InternalPackage("authorization"), "NewServiceRolePermissionChecker").Call(jen.Qual(proj.InternalPackage("authorization"), "ServiceUserRole").Dot("String").Call())), jen.ID("ActiveAccountID").Op(":").Lit(0)),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachAPIClientDatabaseIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachAPIClientDatabaseIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachAPIClientClientIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachAPIClientClientIDToSpan").Call(
						jen.ID("span"),
						jen.Lit("123"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachUserToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachUserToSpan").Call(
						jen.ID("span"),
						jen.ID("exampleUser"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachUserIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachUserIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachUsernameToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachUsernameToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachWebhookIDToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachWebhookIDToSpan").Call(
						jen.ID("span"),
						jen.Lit(123),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachURLToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.List(jen.ID("u"), jen.ID("err")).Assign().Qual("net/url", "ParseRequestURI").Call(jen.Lit("https://todo.verygoodsoftwarenotvirus.ru")),
					jen.Qual(constants.AssertionLibrary, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("AttachURLToSpan").Call(
						jen.ID("span"),
						jen.ID("u"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachRequestURIToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachRequestURIToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachRequestToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.List(jen.ID("req"), jen.ID("err")).Assign().Qual("net/http", "NewRequestWithContext").Call(
						jen.ID("ctx"),
						jen.Qual("net/http", "MethodGet"),
						jen.Lit("/"),
						jen.ID("nil"),
					),
					jen.ID("req").Dot("Header").Dot("Set").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
					jen.Qual(constants.MustAssertPkg, "NoError").Call(
						jen.ID("t"),
						jen.ID("err"),
					),
					jen.Newline(),
					jen.ID("AttachRequestToSpan").Call(
						jen.ID("span"),
						jen.ID("req"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachResponseToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.ID("res").Assign().AddressOf().Qual("net/http", "Response").Valuesln(jen.ID("Header").Op(":").Map(jen.ID("string")).Index().String().Values()),
					jen.ID("res").Dot("Header").Dot("Set").Call(
						jen.ID("t").Dot("Name").Call(),
						jen.Lit("blah"),
					),
					jen.Newline(),
					jen.ID("AttachResponseToSpan").Call(
						jen.ID("span"),
						jen.ID("res"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachErrorToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachErrorToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
						jen.Qual("errors", "New").Call(jen.Lit("blah")),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachDatabaseQueryToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachDatabaseQueryToSpan").Call(
						jen.ID("span"),
						jen.Lit("description"),
						jen.Lit("query"),
						jen.Index().Interface().Values(jen.Lit("blah")),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachQueryFilterToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachQueryFilterToSpan").Call(
						jen.ID("span"),
						jen.Qual(proj.TypesPackage(), "DefaultQueryFilter").Call(),
					),
				),
			),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("with nil"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachQueryFilterToSpan").Call(
						jen.ID("span"),
						jen.ID("nil"),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachSearchQueryToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachSearchQueryToSpan").Call(
						jen.ID("span"),
						jen.ID("t").Dot("Name").Call(),
					),
				),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Func().ID("TestAttachUserAgentDataToSpan").Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
					jen.Newline(),
					jen.ID("AttachUserAgentDataToSpan").Call(
						jen.ID("span"),
						jen.AddressOf().Qual("github.com/mssola/user_agent", "UserAgent").Values(),
					),
				),
			),
		),
		jen.Newline(),
	)

	for _, typ := range proj.DataTypes {
		code.Add(
			jen.Func().IDf("TestAttach%sIDToSpan", typ.Name.Singular()).Params(jen.ID("T").PointerTo().Qual("testing", "T")).Body(
				jen.ID("T").Dot("Parallel").Call(),
				jen.Newline(),
				jen.ID("T").Dot("Run").Call(
					jen.Lit("standard"),
					jen.Func().Params(jen.ID("t").PointerTo().Qual("testing", "T")).Body(
						jen.ID("t").Dot("Parallel").Call(),
						jen.Newline(),
						jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("StartSpan").Call(jen.Qual("context", "Background").Call()),
						jen.Newline(),
						jen.IDf("Attach%sIDToSpan", typ.Name.Singular()).Call(
							jen.ID("span"),
							jen.Lit(123),
						),
					),
				),
			),
			jen.Newline(),
		)
	}

	return code
}
