package workers

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func preUpdatesWorkerDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("PreUpdatesWorker updates data from the pending updates topic to the database."),
		jen.Newline(),
		jen.Type().ID("PreUpdatesWorker").Struct(
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
			jen.ID("encoder").Qual(proj.EncodingPackage(), "ClientEncoder"),
			jen.ID("postUpdatesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("itemsIndexManager").Qual(proj.InternalSearchPackage(), "IndexManager"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ProvidePreUpdatesWorker provides a PreUpdatesWorker."),
		jen.Newline(),
		jen.Func().ID("ProvidePreUpdatesWorker").Paramsln(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
			jen.ID("client").Op("*").Qual("net/http", "Client"),
			jen.ID("dataManager").Qual(proj.DatabasePackage(), "DataManager"),
			jen.ID("postUpdatesPublisher").Qual(proj.InternalMessageQueuePublishersPackage(), "Publisher"),
			jen.ID("searchIndexLocation").Qual(proj.InternalSearchPackage(), "IndexPath"),
			jen.ID("searchIndexProvider").Qual(proj.InternalSearchPackage(), "IndexManagerProvider"),
		).Params(jen.Op("*").ID("PreUpdatesWorker"), jen.ID("error")).Body(
			jen.Const().ID("name").Equals().Lit("pre_updates"),
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
			jen.ID("w").Op(":=").Op("&").ID("PreUpdatesWorker").Valuesln(
				jen.ID("logger").MapAssign().Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("name")).Dot("WithValue").Call(
					jen.Lit("topic"),
					jen.ID("name"),
				), jen.ID("tracer").MapAssign().Qual(proj.InternalTracingPackage(), "NewTracer").Call(jen.ID("name")), jen.ID("encoder").MapAssign().Qual(proj.EncodingPackage(), "ProvideClientEncoder").Call(
					jen.ID("logger"),
					jen.Qual(proj.EncodingPackage(), "ContentTypeJSON"),
				), jen.ID("postUpdatesPublisher").MapAssign().ID("postUpdatesPublisher"), jen.ID("dataManager").MapAssign().ID("dataManager"), jen.ID("itemsIndexManager").MapAssign().ID("itemsIndexManager")),
			jen.Newline(),
			jen.Return().List(jen.ID("w"), jen.Nil()),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("HandleMessage handles a pending update."),
		jen.Newline(),
		jen.Func().Params(jen.ID("w").Op("*").ID("PreUpdatesWorker")).ID("HandleMessage").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("message").Index().ID("byte")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("w").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Var().ID("msg").Op("*").Qual(proj.TypesPackage(), "PreUpdateMessage"),
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
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("dataManager").Dot("UpdateItem").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("Item"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("creating item"),
						)),
					jen.Newline(),
					jen.If(jen.ID("err").Op(":=").ID("w").Dot("itemsIndexManager").Dot("Index").Call(
						jen.ID("ctx"),
						jen.ID("msg").Dot("Item").Dot("ID"),
						jen.ID("msg").Dot("Item"),
					), jen.ID("err").DoesNotEqual().Nil()).Body(
						jen.Return().ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("indexing the item"),
						)),
					jen.Newline(),
					jen.If(jen.ID("w").Dot("postUpdatesPublisher").DoesNotEqual().Nil()).Body(
						jen.ID("dcm").Op(":=").Op("&").Qual(proj.TypesPackage(), "DataChangeMessage").Valuesln(
							jen.ID("DataType").MapAssign().ID("msg").Dot("DataType"),
							jen.ID("Item").MapAssign().ID("msg").Dot("Item"),
							jen.ID("AttributableToUserID").MapAssign().ID("msg").Dot("AttributableToUserID"),
							jen.ID("AttributableToAccountID").MapAssign().ID("msg").Dot("AttributableToAccountID"),
						),
						jen.Newline(),
						jen.If(
							jen.ID("err").Op(":=").ID("w").Dot("postUpdatesPublisher").Dot("Publish").Call(jen.ID("ctx"), jen.ID("dcm")),
							jen.ID("err").DoesNotEqual().Nil(),
						).Body(
							jen.Return().ID("observability").Dot("PrepareError").Call(
								jen.ID("err"),
								jen.ID("logger"),
								jen.ID("span"),
								jen.Lit("publishing data change message"),
							),
						),
					),
				),
				jen.Case(jen.Qual(proj.TypesPackage(), "UserMembershipDataType"), jen.Qual(proj.TypesPackage(), "WebhookDataType")).Body(
					jen.Break()),
			),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	)

	return code
}
