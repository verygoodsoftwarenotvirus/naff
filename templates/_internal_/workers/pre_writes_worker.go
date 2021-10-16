package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preWritesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("PreWritesWorker writes data from the pending writes topic to the database."),
		jen.Newline(),
		jen.Type().ID("PreWritesWorker").Struct(
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ClientEncoder"),
			jen.ID("postWritesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("itemsIndexManager").Qual(proj.InternalSearchPackage(), "IndexManager"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvidePreWritesWorker provides a PreWritesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvidePreWritesWorker").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("client").Op("*").Qual("net/http", "Client"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("postWritesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("searchIndexLocation").Qual(proj.InternalSearchPackage(), "IndexPath"),
			jen.ID("searchIndexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider"),
		).Params(jen.Op("*").ID("PreWritesWorker"), jen.ID("error")).Body(
			jen.Const().ID("name").Equals().Lit("pre_writes"),
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
			jen.ID("w").Op(":=").Op("&").ID("PreWritesWorker").Valuesln(
				jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")).Dot("WithValue").Call(
					jen.Lit("topic"),
					jen.ID("name"),
				), jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().Qual(proj.EncodingPackage(), "ProvideClientEncoder").Call(
					jen.ID("logger"),
					jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
				), jen.ID("postWritesPublisher").MapAssign().ID("postWritesPublisher"), jen.ID("dataManager").MapAssign().ID("dataManager"), jen.ID("itemsIndexManager").MapAssign().ID("itemsIndexManager")),
			jen.Newline(),
			jen.Return().List(jen.ID("w"), jen.ID("err")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HandleMessage handles a pending write."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("PreWritesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Var().ID("msg").Op("*").Qual(proj.TypesPackage(), "PreWriteMessage"),
			jen.Newline(),
			jen.If(jen.ID("err").Op(":=").ID("w").Dot("encoder").Dot("Unmarshal").Call(
				jen.ID("ctx"),
				jen.ID("message"),
				jen.Op("&").ID("msg"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
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
					jen.List(jen.ID("item"), jen.ID("err")).Op(":=").ID("w").Dot("dataManager").Dot("CreateItem").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("Item"),
					), jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("creating item"),
						)),
					jen.Newline(),
					jen.If(jen.ID("err").Equals().ID("w").Dot("itemsIndexManager").Dot("Index").Call(
						jen.ID("ctx"),
						jen.ID("item").Dot("ID"),
						jen.ID("item"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("indexing the item"),
						)),
					jen.Newline(),
					jen.If(jen.ID("w").Dot("postWritesPublisher").DoesNotEqual().Nil()).Body(
						jen.ID("dcm").Op(":=").Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
							jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"), jen.ID("Item").MapAssign().ID("item"),
							jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
							jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
						),
						jen.Newline(),
						jen.If(jen.ID("err").Equals().ID("w").Dot("postWritesPublisher").Dot("Publish").Call(jen.ID("ctx"), jen.ID("dcm")), jen.ID("err").DoesNotEqual().Nil()).Body(
							jen.Return().ID("observability").Dot("PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing to post-writes topic"),
							),
						),
					),
				),
				jen.Case(jen.Qual(proj.TypesPackage(), "WebhookDataType")).Body(
					jen.List(jen.ID("webhook"), jen.ID("err")).Op(":=").ID("w").Dot("dataManager").Dot("CreateWebhook").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("Webhook"),
					),
					jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("creating webhook"),
						)),
					jen.Newline(),
					jen.If(jen.ID("w").Dot("postWritesPublisher").DoesNotEqual().Nil()).Body(
						jen.ID("dcm").Op(":=").Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
							jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
							jen.ID("Webhook").MapAssign().ID("webhook"),
							jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
							jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
						),
						jen.Newline(),
						jen.If(jen.ID("err").Equals().ID("w").Dot("postWritesPublisher").Dot("Publish").Call(
							jen.ID("ctx"),
							jen.ID("dcm"),
						), jen.ID("err").DoesNotEqual().Nil()).Body(
							jen.Return().ID("observability").Dot("PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							)),
					)),
				jen.Case(jen.Qual(proj.TypesPackage(), "UserMembershipDataType")).Body(
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("dataManager").Dot("AddUserToAccount").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("UserMembership"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("creating webhook"),
						)),
					jen.Newline(),
					jen.If(jen.ID("w").Dot("postWritesPublisher").DoesNotEqual().Nil()).Body(
						jen.ID("dcm").Op(":=").Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
							jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"), jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"), jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID")),
						jen.If(jen.ID("err").Op(":=").ID("w").Dot("postWritesPublisher").Dot("Publish").Call(
							jen.ID("ctx"),
							jen.ID("dcm"),
						), jen.ID("err").DoesNotEqual().Nil()).Body(
							jen.Return().ID("observability").Dot("PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							)),
					)),
			),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	return code
}
