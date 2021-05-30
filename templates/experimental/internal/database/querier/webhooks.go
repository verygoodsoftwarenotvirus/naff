package querier

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func webhooksDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dot("WebhookDataManager").Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanWebhook is a consistent way to turn a *sql.Row into a webhook struct."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").ID("database").Dot("Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("webhook").Op("*").ID("types").Dot("Webhook"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.ID("webhook").Op("=").Op("&").ID("types").Dot("Webhook").Valuesln(),
			jen.Var().Defs(
				jen.List(jen.ID("eventsStr"), jen.ID("dataTypesStr"), jen.ID("topicsStr")).ID("string"),
			),
			jen.ID("targetVars").Op(":=").Index().Interface().Valuesln(jen.Op("&").ID("webhook").Dot("ID"), jen.Op("&").ID("webhook").Dot("ExternalID"), jen.Op("&").ID("webhook").Dot("Name"), jen.Op("&").ID("webhook").Dot("ContentType"), jen.Op("&").ID("webhook").Dot("URL"), jen.Op("&").ID("webhook").Dot("Method"), jen.Op("&").ID("eventsStr"), jen.Op("&").ID("dataTypesStr"), jen.Op("&").ID("topicsStr"), jen.Op("&").ID("webhook").Dot("CreatedOn"), jen.Op("&").ID("webhook").Dot("LastUpdatedOn"), jen.Op("&").ID("webhook").Dot("ArchivedOn"), jen.Op("&").ID("webhook").Dot("BelongsToAccount")),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Op("=").ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("filteredCount"),
					jen.Op("&").ID("totalCount"),
				)),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning webhook"),
				))),
			jen.If(jen.ID("events").Op(":=").Qual("strings", "Split").Call(
				jen.ID("eventsStr"),
				jen.ID("querybuilding").Dot("WebhooksTableEventsSeparator"),
			), jen.ID("len").Call(jen.ID("events")).Op(">=").Lit(1).Op("&&").ID("events").Index(jen.Lit(0)).Op("!=").Lit("")).Body(
				jen.ID("webhook").Dot("Events").Op("=").ID("events")),
			jen.If(jen.ID("dataTypes").Op(":=").Qual("strings", "Split").Call(
				jen.ID("dataTypesStr"),
				jen.ID("querybuilding").Dot("WebhooksTableDataTypesSeparator"),
			), jen.ID("len").Call(jen.ID("dataTypes")).Op(">=").Lit(1).Op("&&").ID("dataTypes").Index(jen.Lit(0)).Op("!=").Lit("")).Body(
				jen.ID("webhook").Dot("DataTypes").Op("=").ID("dataTypes")),
			jen.If(jen.ID("topics").Op(":=").Qual("strings", "Split").Call(
				jen.ID("topicsStr"),
				jen.ID("querybuilding").Dot("WebhooksTableTopicsSeparator"),
			), jen.ID("len").Call(jen.ID("topics")).Op(">=").Lit(1).Op("&&").ID("topics").Index(jen.Lit(0)).Op("!=").Lit("")).Body(
				jen.ID("webhook").Dot("Topics").Op("=").ID("topics")),
			jen.Return().List(jen.ID("webhook"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("scanWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").ID("database").Dot("ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("webhooks").Index().Op("*").ID("types").Dot("Webhook"), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("webhook"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanWebhook").Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("scanErr"))),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("filteredCount").Op("==").Lit(0)).Body(
						jen.ID("filteredCount").Op("=").ID("fc")),
					jen.If(jen.ID("totalCount").Op("==").Lit(0)).Body(
						jen.ID("totalCount").Op("=").ID("tc")),
				),
				jen.ID("webhooks").Op("=").ID("append").Call(
					jen.ID("webhooks"),
					jen.ID("webhook"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("rows").Dot("Err").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching webhook from database"),
				))),
			jen.If(jen.ID("err").Op("=").ID("rows").Dot("Close").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching webhook from database"),
				))),
			jen.Return().List(jen.ID("webhooks"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetWebhook fetches a webhook from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID")).ID("uint64")).Params(jen.Op("*").ID("types").Dot("Webhook"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("webhookID").Op("==").Lit(0).Op("||").ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("WebhookIDKey").Op(":").ID("webhookID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetWebhookQuery").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
				jen.ID("accountID"),
			),
			jen.ID("row").Op(":=").ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("webhook"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.List(jen.ID("webhook"), jen.ID("_"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanWebhook").Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning webhook"),
				))),
			jen.Return().List(jen.ID("webhook"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllWebhooksCount fetches the count of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllWebhooksCount").Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("q").Dot("logger"),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAllWebhooksCountQuery").Call(jen.ID("ctx")),
				jen.Lit("fetching count of webhooks"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying for count of webhooks"),
				))),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetWebhooks fetches a list of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("WebhookList"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.ID("x").Op(":=").Op("&").ID("types").Dot("WebhookList").Valuesln(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit"))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetWebhooksQuery").Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("filter"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("webhooks"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching webhook from database"),
				))),
			jen.If(jen.List(jen.ID("x").Dot("Webhooks"), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dot("scanWebhooks").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning database response"),
				))),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAllWebhooks fetches a list of webhooks from the database that meet a particular filter."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAllWebhooks").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("resultChannel").Chan().Index().Op("*").ID("types").Dot("Webhook"), jen.ID("batchSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("batchSize").Op("==").Lit(0)).Body(
				jen.ID("batchSize").Op("=").ID("defaultBatchSize")),
			jen.If(jen.ID("resultChannel").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("batch_size"),
				jen.ID("batchSize"),
			),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID("q").Dot("GetAllWebhooksCount").Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("fetching count of webhooks"),
				)),
			jen.ID("increment").Op(":=").ID("uint64").Call(jen.ID("batchSize")),
			jen.For(jen.ID("beginID").Op(":=").ID("uint64").Call(jen.Lit(1)), jen.ID("beginID").Op("<=").ID("count"), jen.ID("beginID").Op("+=").ID("increment")).Body(
				jen.ID("endID").Op(":=").ID("beginID").Op("+").ID("increment"),
				jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).ID("uint64")).Body(
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetBatchOfWebhooksQuery").Call(
						jen.ID("ctx"),
						jen.ID("begin"),
						jen.ID("end"),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("query").Op(":").ID("query"), jen.Lit("begin").Op(":").ID("begin"), jen.Lit("end").Op(":").ID("end"))),
					jen.List(jen.ID("rows"), jen.ID("queryErr")).Op(":=").ID("q").Dot("db").Dot("Query").Call(
						jen.ID("query"),
						jen.ID("args").Op("..."),
					),
					jen.If(jen.Qual("errors", "Is").Call(
						jen.ID("queryErr"),
						jen.Qual("database/sql", "ErrNoRows"),
					)).Body(
						jen.Return()).Else().If(jen.ID("queryErr").Op("!=").ID("nil")).Body(
						jen.ID("logger").Dot("Error").Call(
							jen.ID("queryErr"),
							jen.Lit("querying for database rows"),
						),
						jen.Return(),
					),
					jen.List(jen.ID("webhooks"), jen.ID("_"), jen.ID("_"), jen.ID("scanErr")).Op(":=").ID("q").Dot("scanWebhooks").Call(
						jen.ID("ctx"),
						jen.ID("rows"),
						jen.ID("false"),
					),
					jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
						jen.ID("logger").Dot("Error").Call(
							jen.ID("scanErr"),
							jen.Lit("scanning database rows"),
						),
						jen.Return(),
					),
					jen.ID("resultChannel").ReceiveFromChannel().ID("webhooks"),
				).Call(
					jen.ID("beginID"),
					jen.ID("endID"),
				),
			),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("CreateWebhook creates a webhook in a database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("CreateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("WebhookCreationInput"), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Webhook"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided"))),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("createdByUser"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("input").Dot("BelongsToAccount"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("input").Dot("BelongsToAccount"),
			),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				))),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildCreateWebhookQuery").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.List(jen.ID("id"), jen.ID("err")).Op(":=").ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("false"),
				jen.Lit("webhook creation"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("creating webhook"),
				)),
			),
			jen.ID("x").Op(":=").Op("&").ID("types").Dot("Webhook").Valuesln(jen.ID("ID").Op(":").ID("id"), jen.ID("Name").Op(":").ID("input").Dot("Name"), jen.ID("ContentType").Op(":").ID("input").Dot("ContentType"), jen.ID("URL").Op(":").ID("input").Dot("URL"), jen.ID("Method").Op(":").ID("input").Dot("Method"), jen.ID("Events").Op(":").ID("input").Dot("Events"), jen.ID("DataTypes").Op(":").ID("input").Dot("DataTypes"), jen.ID("Topics").Op(":").ID("input").Dot("Topics"), jen.ID("BelongsToAccount").Op(":").ID("input").Dot("BelongsToAccount"), jen.ID("CreatedOn").Op(":").ID("q").Dot("currentTime").Call()),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildWebhookCreationEventEntry").Call(
					jen.ID("x"),
					jen.ID("createdByUser"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing webhook creation audit log entry"),
				)),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				))),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("x").Dot("ID"),
			),
			jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("x").Dot("ID"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Lit("webhook created")),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("UpdateWebhook updates a particular webhook."),
		jen.Line(),
		jen.Func().Comment("NOTE: this function expects the provided input to have a non-zero ID.").Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("UpdateWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").ID("types").Dot("Webhook"), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").ID("types").Dot("FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("changedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.If(jen.ID("updated").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided")),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("changedByUser"),
			),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("ID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("BelongsToAccount"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("updated").Dot("ID"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("changedByUser"),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("updated").Dot("BelongsToAccount"),
			),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildUpdateWebhookQuery").Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("webhook update"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("updating webhook"),
				),
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating webhook"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildWebhookUpdateEventEntry").Call(
					jen.ID("changedByUser"),
					jen.ID("updated").Dot("BelongsToAccount"),
					jen.ID("updated").Dot("ID"),
					jen.ID("changes"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("logger").Dot("Error").Call(
					jen.ID("err"),
					jen.Lit("writing webhook update audit log entry"),
				),
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing webhook update audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("webhook updated")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("ArchiveWebhook archives a webhook from the database."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("ArchiveWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.ID("webhookID"), jen.ID("accountID"), jen.ID("archivedByUserID")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("webhookID").Op("==").Lit(0).Op("||").ID("accountID").Op("==").Lit(0).Op("||").ID("archivedByUserID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided")),
			jen.ID("tracing").Dot("AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("archivedByUserID"),
			),
			jen.ID("tracing").Dot("AttachWebhookIDToSpan").Call(
				jen.ID("span"),
				jen.ID("webhookID"),
			),
			jen.ID("tracing").Dot("AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.ID("keys").Dot("WebhookIDKey").Op(":").ID("webhookID"), jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"), jen.ID("keys").Dot("RequesterIDKey").Op(":").ID("archivedByUserID"))),
			jen.List(jen.ID("tx"), jen.ID("err")).Op(":=").ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildArchiveWebhookQuery").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
				jen.ID("accountID"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Lit("webhook archive"),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving webhook"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dot("BuildWebhookArchiveEventEntry").Call(
					jen.ID("archivedByUserID"),
					jen.ID("accountID"),
					jen.ID("webhookID"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("writing webhook archive audit log entry"),
				),
			),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			jen.ID("logger").Dot("Info").Call(jen.Lit("webhook archived")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("GetAuditLogEntriesForWebhook fetches a list of audit log entries from the database that relate to a given webhook."),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).ID("GetAuditLogEntriesForWebhook").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("webhookID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.If(jen.ID("webhookID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided"))),
			jen.ID("logger").Op(":=").ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("WebhookIDKey"),
				jen.ID("webhookID"),
			),
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID("q").Dot("sqlQueryBuilder").Dot("BuildGetAuditLogEntriesForWebhookQuery").Call(
				jen.ID("ctx"),
				jen.ID("webhookID"),
			),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit("audit log entries for webhook"),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying database for audit log entries"),
				))),
			jen.List(jen.ID("auditLogEntries"), jen.ID("_"), jen.ID("err")).Op(":=").ID("q").Dot("scanAuditLogEntries").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning response from database"),
				))),
			jen.Return().List(jen.ID("auditLogEntries"), jen.ID("nil")),
		),
		jen.Line(),
	)

	return code
}
