package bleve

import (
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bleveDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("base").Op("=").Lit(10),
			jen.ID("bitSize").Op("=").Lit(64),
			jen.ID("testingSearchIndexName").ID("search").Dot("IndexName").Op("=").Lit("example_index_name"),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidIndexName").Op("=").Qual("errors", "New").Call(jen.Lit("invalid index name")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("search").Dot("IndexManager").Op("=").Parens(jen.Op("*").ID("bleveIndexManager")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("bleveIndexManager").Struct(
				jen.ID("index").Qual("github.com/blevesearch/bleve/v2", "Index"),
				jen.ID("logger").ID("logging").Dot("Logger"),
				jen.ID("tracer").ID("tracing").Dot("Tracer"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("NewBleveIndexManager instantiates a bleve index."),
		jen.Line(),
		jen.Func().ID("NewBleveIndexManager").Params(jen.ID("path").ID("search").Dot("IndexPath"), jen.ID("name").ID("search").Dot("IndexName"), jen.ID("logger").ID("logging").Dot("Logger")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
			jen.Var().Defs(
				jen.ID("index").Qual("github.com/blevesearch/bleve/v2", "Index"),
			),
			jen.List(jen.ID("preexistingIndex"), jen.ID("err")).Op(":=").Qual("github.com/blevesearch/bleve/v2", "Open").Call(jen.ID("string").Call(jen.ID("path"))),
			jen.If(jen.ID("err").Op("==").ID("nil")).Body(
				jen.ID("index").Op("=").ID("preexistingIndex")),
			jen.If(jen.Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("github.com/blevesearch/bleve/v2", "ErrorIndexPathDoesNotExist"),
			).Op("||").Qual("errors", "Is").Call(
				jen.ID("err"),
				jen.Qual("github.com/blevesearch/bleve/v2", "ErrorIndexMetaMissing"),
			)).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("path"),
					jen.ID("path"),
				).Dot("Debug").Call(jen.Lit("tried to open existing index, but didn't find it")),
				jen.Switch(jen.ID("name")).Body(
					jen.Case(jen.ID("testingSearchIndexName")).Body(
						jen.List(jen.ID("index"), jen.ID("err")).Op("=").Qual("github.com/blevesearch/bleve/v2", "New").Call(
							jen.ID("string").Call(jen.ID("path")),
							jen.Qual("github.com/blevesearch/bleve/v2", "NewIndexMapping").Call(),
						)),
					jen.Case(jen.ID("types").Dot("ItemsSearchIndexName")).Body(
						jen.List(jen.ID("index"), jen.ID("err")).Op("=").Qual("github.com/blevesearch/bleve/v2", "New").Call(
							jen.ID("string").Call(jen.ID("path")),
							jen.ID("buildItemMapping").Call(),
						)),
					jen.Default().Body(
						jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit("opening %s index: %w"),
							jen.ID("name"),
							jen.ID("errInvalidIndexName"),
						))),
				),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("failed to create new index"),
					),
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
			),
			jen.ID("serviceName").Op(":=").Qual("fmt", "Sprintf").Call(
				jen.Lit("%s_search"),
				jen.ID("name"),
			),
			jen.ID("im").Op(":=").Op("&").ID("bleveIndexManager").Valuesln(jen.ID("index").Op(":").ID("index"), jen.ID("logger").Op(":").ID("logging").Dot("EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")), jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName"))),
			jen.Return().List(jen.ID("im"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Index implements our IndexManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").Op("*").ID("bleveIndexManager")).ID("Index").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64"), jen.ID("value").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			).Dot("Debug").Call(jen.Lit("adding to index")),
			jen.Return().ID("sm").Dot("index").Dot("Index").Call(
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.ID("base"),
				),
				jen.ID("value"),
			),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("search executes search queries."),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").Op("*").ID("bleveIndexManager")).ID("search").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("accountID").ID("uint64"), jen.ID("forServiceAdmin").ID("bool")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("tracing").Dot("AttachSearchQueryToSpan").Call(
				jen.ID("span"),
				jen.ID("query"),
			),
			jen.ID("logger").Op(":=").ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("SearchQueryKey"),
				jen.ID("query"),
			),
			jen.If(jen.ID("query").Op("==").Lit("")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("search").Dot("ErrEmptyQueryProvided"))),
			jen.If(jen.Op("!").ID("forServiceAdmin").Op("&&").ID("accountID").Op("!=").Lit(0)).Body(
				jen.ID("logger").Op("=").ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("AccountIDKey"),
					jen.ID("accountID"),
				)),
			jen.ID("q").Op(":=").Qual("github.com/blevesearch/bleve/v2", "NewFuzzyQuery").Call(jen.ID("query")),
			jen.ID("q").Dot("SetFuzziness").Call(jen.ID("searcher").Dot("MaxFuzziness")),
			jen.List(jen.ID("searchResults"), jen.ID("err")).Op(":=").ID("sm").Dot("index").Dot("SearchInContext").Call(
				jen.ID("ctx"),
				jen.Qual("github.com/blevesearch/bleve/v2", "NewSearchRequest").Call(jen.ID("q")),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("performing search query"),
				))),
			jen.For(jen.List(jen.ID("_"), jen.ID("result")).Op(":=").Range().ID("searchResults").Dot("Hits")).Body(
				jen.List(jen.ID("x"), jen.ID("parseErr")).Op(":=").Qual("strconv", "ParseUint").Call(
					jen.ID("result").Dot("ID"),
					jen.ID("base"),
					jen.ID("bitSize"),
				),
				jen.If(jen.ID("parseErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("parseErr"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("parsing integer stored in search index for #%s"),
						jen.ID("result").Dot("ID"),
					))),
				jen.ID("ids").Op("=").ID("append").Call(
					jen.ID("ids"),
					jen.ID("x"),
				),
			),
			jen.Return().List(jen.ID("ids"), jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Search implements our IndexManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").Op("*").ID("bleveIndexManager")).ID("Search").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string"), jen.ID("accountID").ID("uint64")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.Return().ID("sm").Dot("search").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("accountID"),
				jen.ID("false"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("SearchForAdmin implements our IndexManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").Op("*").ID("bleveIndexManager")).ID("SearchForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string")).Params(jen.ID("ids").Index().ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.Return().ID("sm").Dot("search").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.Lit(0),
				jen.ID("true"),
			)),
		jen.Line(),
	)

	code.Add(
		jen.Comment("Delete implements our IndexManager interface."),
		jen.Line(),
		jen.Func().Params(jen.ID("sm").Op("*").ID("bleveIndexManager")).ID("Delete").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.ID("logger").Op(":=").ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			),
			jen.If(jen.ID("err").Op(":=").ID("sm").Dot("index").Dot("Delete").Call(jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.ID("base"),
			)), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("removing from index"),
				)),
			jen.ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			).Dot("Debug").Call(jen.Lit("removed from index")),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	)

	return code
}
