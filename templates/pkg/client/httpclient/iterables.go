package httpclient

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Comment("ItemExists retrieves whether an item exists."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ItemExists").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.ID("bool"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("false"),
					jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildItemExistsRequest").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("building item existence request"),
					),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("exists"),
				jen.ID("err")).Op(":=").ID("c").Dot("responseIsOK").Call(
				jen.ID("ctx"),
				jen.ID("req"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("checking existence for item #%d"),
						jen.ID("itemID"),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("exists"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("GetItem gets an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("building get item request"),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID("item").Op("*").ID("types").Dot("Item"),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("item"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("retrieving item"),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("item"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	if typ.SearchEnabled {
		code.Add(
			jen.Comment("SearchItems searches through a list of items."),
			jen.Newline(),
			jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("SearchItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("limit").ID("uint8")).Params(jen.Index().Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
				jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
				jen.Defer().ID("span").Dot("End").Call(),
				jen.Newline(),
				jen.If(jen.ID("query").Op("==").Lit("")).Body(
					jen.Return().List(jen.ID("nil"),
						jen.ID("ErrEmptyQueryProvided")),
				),
				jen.Newline(),
				jen.If(jen.ID("limit").Op("==").Lit(0)).Body(
					jen.ID("limit").Op("=").Lit(20)),
				jen.Newline(),
				jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("SearchQueryKey"),
					jen.ID("query"),
				).Dot("WithValue").Call(
					jen.ID("keys").Dot("FilterLimitKey"),
					jen.ID("limit"),
				),
				jen.Newline(),
				jen.List(jen.ID("req"),
					jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildSearchItemsRequest").Call(
					jen.ID("ctx"),
					jen.ID("query"),
					jen.ID("limit"),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"),
						jen.ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("building search for items request"),
						),
					),
				),
				jen.Newline(),
				jen.Var().ID("items").Index().Op("*").ID("types").Dot("Item"),
				jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
					jen.ID("ctx"),
					jen.ID("req"),
					jen.Op("&").ID("items"),
				),
					jen.ID("err").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"),
						jen.ID("observability").Dot("PrepareError").Call(
							jen.ID("err"),
							jen.ID("logger"),
							jen.ID("span"),
							jen.Lit("retrieving items"),
						),
					),
				),
				jen.Newline(),
				jen.Return().List(jen.ID("items"), jen.ID("nil")),
			),
			jen.Newline(),
		)
	}

	code.Add(
		jen.Comment("GetItems retrieves a list of items."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetItems").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("filter").Op("*").ID("types").Dot("QueryFilter")).Params(jen.Op("*").ID("types").Dot("ItemList"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("loggerWithFilter").Call(jen.ID("filter")),
			jen.ID("tracing").Dot("AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetItemsRequest").Call(
				jen.ID("ctx"),
				jen.ID("filter"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("building items list request"),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID("items").Op("*").ID("types").Dot("ItemList"),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("items"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("retrieving items"),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("items"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("CreateItem creates an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("CreateItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dot("ItemCreationInput")).Params(jen.Op("*").ID("types").Dot("Item"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("ErrNilInputProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger"),
			jen.Newline(),
			jen.If(jen.ID("err").Op(":=").ID("input").Dot("ValidateWithContext").Call(jen.ID("ctx")),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("validating input"),
					),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildCreateItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("building create item request"),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID("item").Op("*").ID("types").Dot("Item"),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("item"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("creating item"),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("item"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("UpdateItem updates an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("UpdateItem").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("item").Op("*").ID("types").Dot("Item")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("item").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided"),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("item").Dot("ID"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildUpdateItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("item"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building update item request"),
				)),
			jen.Newline(),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("item"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("updating item #%d"),
					jen.ID("item").Dot("ID"),
				)),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("ArchiveItem archives an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("ArchiveItem").Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("itemID").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildArchiveItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("building archive item request"),
				)),
			jen.Newline(),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.ID("nil"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("archiving item #%d"),
					jen.ID("itemID"),
				)),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("GetAuditLogForItem retrieves a list of audit log entries pertaining to an item."),
		jen.Newline(),
		jen.Func().Params(jen.ID("c").Op("*").ID("Client")).ID("GetAuditLogForItem").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("itemID").ID("uint64")).Params(jen.Index().Op("*").ID("types").Dot("AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Op(":=").ID("c").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("itemID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			jen.ID("logger").Op(":=").ID("c").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("ItemIDKey"),
				jen.ID("itemID"),
			),
			jen.Newline(),
			jen.List(jen.ID("req"),
				jen.ID("err")).Op(":=").ID("c").Dot("requestBuilder").Dot("BuildGetAuditLogForItemRequest").Call(
				jen.ID("ctx"),
				jen.ID("itemID"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("building get audit log entries for item request"),
					),
				),
			),
			jen.Newline(),
			jen.Var().ID("entries").Index().Op("*").ID("types").Dot("AuditLogEntry"),
			jen.If(jen.ID("err").Op("=").ID("c").Dot("fetchAndUnmarshal").Call(
				jen.ID("ctx"),
				jen.ID("req"),
				jen.Op("&").ID("entries"),
			),
				jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"),
					jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("err"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("retrieving plan"),
					),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("entries"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	return code
}
