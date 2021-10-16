package workers

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preArchivesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("PreArchivesWorker archives data from the pending archives topic to the database."),
		jen.Newline(),
		jen.Type().ID("PreArchivesWorker").Struct(
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ClientEncoder"),
			jen.ID("postArchivesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("itemsIndexManager").Qual(proj.InternalSearchPackage(), "IndexManager"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvidePreArchivesWorker provides a PreArchivesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvidePreArchivesWorker").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("client").Op("*").Qual("net/http", "Client"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("postArchivesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("searchIndexLocation").Qual(proj.InternalSearchPackage(), "IndexPath"),
			jen.ID("searchIndexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider"),
		).Params(jen.Op("*").ID("PreArchivesWorker"), jen.ID("error")).Body(
			jen.Const().ID("name").Equals().Lit("pre_archives"),
			jen.Newline(),
			jen.List(jen.ID("itemsIndexManager"), jen.ID("err")).Op(":=").ID("searchIndexProvider").Call(
				jen.ID("ctx"),
				jen.ID("logger"),
				jen.ID("client"),
				jen.ID("searchIndexLocation"),
				jen.Lit("items"),
				jen.Lit("name"),
				jen.Lit("description"),
			),
			jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
					jen.Lit("setting up items search index manager: %w"),
					jen.ID("err"),
				))),
			jen.Newline(),
			jen.ID("w").Op(":=").Op("&").ID("PreArchivesWorker").Valuesln(
				jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")).Dot("WithValue").Call(
					jen.Lit("topic"),
					jen.ID("name"),
				), jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().Qual(proj.EncodingPackage(), "ProvideClientEncoder").Call(
					jen.ID("logger"),
					jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
				), jen.ID("postArchivesPublisher").MapAssign().ID("postArchivesPublisher"), jen.ID("dataManager").MapAssign().ID("dataManager"), jen.ID("itemsIndexManager").MapAssign().ID("itemsIndexManager")),
			jen.Newline(),
			jen.Return().List(jen.ID("w"), jen.Nil()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HandleMessage handles a pending archive."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("PreArchivesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Var().ID("msg").Op("*").Qual(proj.TypesPackage(), "PreArchiveMessage"),
			jen.Newline(),
			jen.If(jen.ID("err").Op(":=").ID("w").Dot("encoder").Dot("Unmarshal").Call(
				jen.ID("ctx"),
				jen.ID("message"),
				jen.Op("&").ID("msg"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("w").Dot("logger"),
					jen.ID("span"),
					jen.Lit("unmarshalling message"),
				)),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("msg").Dot("AttributableToUserID"),
			),
			jen.ID("logger").Op(":=").ID("w").Dot("logger").Dot("WithValue").Call(
				jen.Lit("data_type"),
				jen.ID("msg").Dot("DataType"),
			),
			jen.Newline(),
			jen.ID("logger").Dot("Debug").Call(jen.Lit("message read")),
			jen.Newline(),
			jen.Switch(jen.ID("msg").Dot("DataType")).Body(
				jen.Case(jen.Qual(proj.TypesPackage(), "ItemDataType")).Body(
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("dataManager").Dot("ArchiveItem").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("ItemID"),
						jen.ID("msg").Dot("AttributableToAccountID"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
							jen.ID("err"),
							jen.ID("w").Dot("logger"),
							jen.ID("span"),
							jen.Lit("archiving item"),
						)),
					jen.Newline(),
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("itemsIndexManager").Dot("Delete").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("ItemID"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
							jen.ID("err"),
							jen.ID("w").Dot("logger"),
							jen.ID("span"),
							jen.Lit("removing item from index"),
						)),
					jen.Newline(),
					jen.If(jen.ID("w").Dot("postArchivesPublisher").DoesNotEqual().Nil()).Body(
						jen.ID("dcm").Op(":=").Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
							jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
							jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
							jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
						),
						jen.Newline(),
						jen.If(jen.ID("err").Op(":=").ID("w").Dot("postArchivesPublisher").Dot("Publish").Call(
							jen.ID("ctx"),
							jen.ID("dcm"),
						), jen.ID("err").DoesNotEqual().Nil()).Body(
							jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							)),
					)),
				jen.Case(jen.Qual(proj.TypesPackage(), "WebhookDataType")).Body(
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("dataManager").Dot("ArchiveWebhook").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("WebhookID"),
						jen.ID("msg").Dot("AttributableToAccountID"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
							jen.ID("err"),
							jen.ID("w").Dot("logger"),
							jen.ID("span"),
							jen.Lit("creating item"),
						)),
					jen.Newline(),
					jen.If(jen.ID("w").Dot("postArchivesPublisher").DoesNotEqual().Nil()).Body(
						jen.ID("dcm").Op(":=").Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
							jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
							jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
							jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
						),
						jen.Newline(),
						jen.If(jen.ID("err").Op(":=").ID("w").Dot("postArchivesPublisher").Dot("Publish").Call(
							jen.ID("ctx"),
							jen.ID("dcm"),
						), jen.ID("err").DoesNotEqual().Nil()).Body(
							jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							)),
					)),
				jen.Case(jen.Qual(proj.TypesPackage(), "UserMembershipDataType")).Body(
					jen.Break()),
			),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	return code
}
