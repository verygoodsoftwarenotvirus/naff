package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preArchivesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Type().Defs(
			jen.ID("PreArchivesWorker").Struct(
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
				jen.ID("encoder").ID("encoding").Dot("ClientEncoder"),
				jen.ID("postArchivesPublisher").ID("publishers").Dot("Publisher"),
				jen.ID("dataManager").ID("database").Dot("DataManager"),
				jen.ID("itemsIndexManager").ID("search").Dot("IndexManager"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvidePreArchivesWorker provides a PreArchivesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvidePreArchivesWorker").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("logger").ID("logging").Dot("Logger"), jen.ID("client").Op("*").Qual("net/http", "Client"), jen.ID("dataManager").ID("database").Dot("DataManager"), jen.ID("postArchivesPublisher").ID("publishers").Dot("Publisher"), jen.ID("searchIndexLocation").ID("search").Dot("IndexPath"), jen.ID("searchIndexProvider").ID("search").Dot("IndexManagerProvider")).Params(jen.Op("*").ID("PreArchivesWorker"), jen.ID("error")).Body(
			jen.Var().Defs(
				jen.ID("name").Op("=").Lit("pre_archives"),
			),
			jen.List(jen.ID("itemsIndexManager"), jen.ID("err")).Op(":=").ID("searchIndexProvider").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("searchIndexLocation"),
				jen.Lit("items"),
				jen.Lit("name"),
				jen.Lit("description"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("setting up items search index manager: %w"),
					jen.ID("err"),
				))),
			jen.ID("w").Op(":=").Op("&").ID("PreArchivesWorker").Valuesln(jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")).Dot("WithValue").Call(
				jen.Lit("topic"),
				jen.ID("name"),
			), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("name")), jen.ID("encoder").Op(":").ID("encoding").Dot("ProvideClientEncoder").Call(
				jen.ID("logger"),
				jen.ID("encoding").Dot("ContentTypeJSON"),
			), jen.ID("postArchivesPublisher").Op(":").ID("postArchivesPublisher"), jen.ID("dataManager").Op(":").ID("dataManager"), jen.ID("itemsIndexManager").Op(":").ID("itemsIndexManager")),
			jen.Return().List(jen.ID("w"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HandleMessage handles a pending archive."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("PreArchivesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Var().Defs(
				jen.ID("msg").Op("*").ID("types").Dot("PreArchiveMessage"),
			),
			jen.If(jen.ID("err").Op(":=").ID("w").Dot("encoder").Dot("Unmarshal").Call(
				jen.ID("ctx"),
				jen.ID("message"),
				jen.Op("&").ID("msg"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("w").Dot("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling message"),
				)),
			jen.ID("tracing").Dot("AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("msg").Dot("AttributableToUserID"),
			),
			jen.ID("logger").Op(":=").ID("w").Dot("logger").Dot("WithValue").Call(
				jen.Lit("data_type"),
				jen.ID("msg").Dot("DataType"),
			),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("message read")),
			jen.Switch(jen.ID("msg").Dot("DataType")).Body(
				jen.Case(jen.ID("types").Dot("ItemDataType")).Body(
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("dataManager").Dot("ArchiveItem").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("RelevantID"),
						jen.ID("msg").Dot("AttributableToAccountID"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("w").Dot("logger"),
							jen.ID("span"),
							jen.Lit("archiving item"),
						)), jen.If(jen.ID("err").Op(":=").ID("w").Dot("itemsIndexManager").Dot("Delete").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("RelevantID"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("w").Dot("logger"),
							jen.ID("span"),
							jen.Lit("removing item from index"),
						)), jen.If(jen.ID("w").Dot("postArchivesPublisher").Op("!=").ID("nil")).Body(
						jen.ID("dcm").Op(":=").Op("&").ID("types").Dot("DataChangeMessage").Valuesln(jen.ID("DataType").Op(":").ID("msg").Dot("DataType"), jen.ID("AttributableToUserID").Op(":").ID("msg").Dot("AttributableToUserID"), jen.ID("AttributableToAccountID").Op(":").ID("msg").Dot("AttributableToAccountID")),
						jen.If(jen.ID("err").Op(":=").ID("w").Dot("postArchivesPublisher").Dot("Publish").Call(
							jen.ID("ctx"),
							jen.ID("dcm"),
						), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().ID("observability").Dot("PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							)),
					)),
				jen.Case(jen.ID("types").Dot("WebhookDataType")).Body(
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("dataManager").Dot("ArchiveWebhook").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("RelevantID"),
						jen.ID("msg").Dot("AttributableToAccountID"),
					), jen.ID("err").Op("!=").ID("nil")).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("w").Dot("logger"),
							jen.ID("span"),
							jen.Lit("creating item"),
						)), jen.If(jen.ID("w").Dot("postArchivesPublisher").Op("!=").ID("nil")).Body(
						jen.ID("dcm").Op(":=").Op("&").ID("types").Dot("DataChangeMessage").Valuesln(jen.ID("DataType").Op(":").ID("msg").Dot("DataType"), jen.ID("AttributableToUserID").Op(":").ID("msg").Dot("AttributableToUserID"), jen.ID("AttributableToAccountID").Op(":").ID("msg").Dot("AttributableToAccountID")),
						jen.If(jen.ID("err").Op(":=").ID("w").Dot("postArchivesPublisher").Dot("Publish").Call(
							jen.ID("ctx"),
							jen.ID("dcm"),
						), jen.ID("err").Op("!=").ID("nil")).Body(
							jen.Return().ID("observability").Dot("PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							)),
					)),
				jen.Case(jen.ID("types").Dot("UserMembershipDataType")).Body(
					jen.Break()),
			),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	return code
}
