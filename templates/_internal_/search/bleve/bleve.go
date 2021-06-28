package bleve

import (
	"fmt"
	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	models "gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func bleveDotGo(proj *models.Project) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Const().Defs(
			jen.ID("base").Equals().Lit(10),
			jen.ID("bitSize").Equals().Lit(64),
			jen.Newline(),
			jen.Comment("testingSearchIndexName is an index name that is only valid for testing's sake."),
			jen.ID("testingSearchIndexName").ID("search").Dot("IndexName").Equals().Lit("example_index_name"),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().Defs(
			jen.ID("errInvalidIndexName").Equals().Qual("errors", "New").Call(jen.Lit("invalid index name")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Var().ID("_").ID("search").Dot("IndexManager").Equals().Parens(jen.PointerTo().ID("bleveIndexManager")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(
		jen.Type().Defs(
			jen.ID("bleveIndexManager").Struct(
				jen.ID("index").Qual(constants.SearchLibrary, "Index"),
				jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger"),
				jen.ID("tracer").Qual(proj.InternalTracingPackage(), "Tracer"),
			),
		),
		jen.Newline(),
	)

	typeCases := []jen.Code{
		jen.Case(jen.ID("testingSearchIndexName")).Body(
			jen.List(jen.ID("index"), jen.ID("err")).Equals().Qual(constants.SearchLibrary, "New").Call(
				jen.String().Call(jen.ID("path")),
				jen.Qual(constants.SearchLibrary, "NewIndexMapping").Call(),
			),
		),
	}
	for _, typ := range proj.DataTypes {
		if typ.SearchEnabled {
			typeCases = append(typeCases,
				jen.Case(jen.Qual(proj.TypesPackage(), fmt.Sprintf("%sSearchIndexName", typ.Name.Plural()))).Body(
					jen.List(jen.ID("index"), jen.ID("err")).Equals().Qual(constants.SearchLibrary, "New").Call(
						jen.String().Call(jen.ID("path")),
						jen.IDf("build%sMapping", typ.Name.Singular()).Call(),
					),
				),
			)
		}
	}
	typeCases = append(typeCases,
		jen.Default().Body(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("opening %s index: %w"),
				jen.ID("name"),
				jen.ID("errInvalidIndexName"),
			)),
		),
	)

	code.Add(
		jen.Comment("NewBleveIndexManager instantiates a bleve index."),
		jen.Newline(),
		jen.Func().ID("NewBleveIndexManager").Params(jen.ID("path").ID("search").Dot("IndexPath"), jen.ID("name").ID("search").Dot("IndexName"), jen.ID("logger").Qual(proj.InternalLoggingPackage(), "Logger")).Params(jen.ID("search").Dot("IndexManager"), jen.ID("error")).Body(
			jen.Var().ID("index").Qual(constants.SearchLibrary, "Index"),
			jen.Newline(),
			jen.List(jen.ID("preexistingIndex"), jen.ID("err")).Assign().Qual(constants.SearchLibrary, "Open").Call(jen.String().Call(jen.ID("path"))),
			jen.If(jen.ID("err").IsEqualTo().ID("nil")).Body(
				jen.ID("index").Equals().ID("preexistingIndex"),
			),
			jen.Newline(),
			jen.If(jen.Qual("errors", "Is").Call(jen.ID("err"), jen.Qual(constants.SearchLibrary, "ErrorIndexPathDoesNotExist")).Op("||").Qual("errors", "Is").Call(jen.ID("err"), jen.Qual(constants.SearchLibrary, "ErrorIndexMetaMissing"))).Body(
				jen.ID("logger").Dot("WithValue").Call(
					jen.Lit("path"),
					jen.ID("path"),
				).Dot("Debug").Call(jen.Lit("tried to open existing index, but didn't find it")),
				jen.Newline(),
				jen.Switch(jen.ID("name")).Body(typeCases...),
				jen.Newline(),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
					jen.ID("logger").Dot("Error").Call(
						jen.ID("err"),
						jen.Lit("failed to create new index"),
					),
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
			),
			jen.Newline(),
			jen.ID("serviceName").Assign().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s_search"),
				jen.ID("name"),
			),
			jen.Newline(),
			jen.ID("im").Assign().AddressOf().ID("bleveIndexManager").Valuesln(
				jen.ID("index").Op(":").ID("index"),
				jen.ID("logger").Op(":").Qual(proj.InternalLoggingPackage(), "EnsureLogger").Call(jen.ID("logger")).Dot("WithName").Call(jen.ID("serviceName")),
				jen.ID("tracer").Op(":").ID("tracing").Dot("NewTracer").Call(jen.ID("serviceName")),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("im"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("Index implements our IndexManager interface."),
		jen.Newline(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("Index").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").Uint64(), jen.ID("value").Interface()).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			).Dot("Debug").Call(jen.Lit("adding to index")),
			jen.Newline(),
			jen.Return().ID("sm").Dot("index").Dot("Index").Call(
				jen.Qual("strconv", "FormatUint").Call(
					jen.ID("id"),
					jen.ID("base"),
				),
				jen.ID("value"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("search executes search queries."),
		jen.Newline(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("search").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").String(), jen.ID("accountID").Uint64(), jen.ID("forServiceAdmin").ID("bool")).Params(jen.ID("ids").Index().Uint64(), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dot("AttachSearchQueryToSpan").Call(
				jen.ID("span"),
				jen.ID("query"),
			),
			jen.ID("logger").Assign().ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("SearchQueryKey"),
				jen.ID("query"),
			),
			jen.Newline(),
			jen.If(jen.ID("query").IsEqualTo().Lit("")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("search").Dot("ErrEmptyQueryProvided")),
			),
			jen.Newline(),
			jen.If(jen.Op("!").ID("forServiceAdmin").Op("&&").ID("accountID").Op("!=").Lit(0)).Body(
				jen.ID("logger").Equals().ID("logger").Dot("WithValue").Call(
					jen.ID("keys").Dot("AccountIDKey"),
					jen.ID("accountID"),
				),
			),
			jen.Newline(),
			jen.ID("q").Assign().Qual(constants.SearchLibrary, "NewFuzzyQuery").Call(jen.ID("query")),
			jen.ID("q").Dot("SetFuzziness").Call(jen.Qual("github.com/blevesearch/bleve/v2/search/searcher", "MaxFuzziness")),
			jen.Newline(),
			jen.List(jen.ID("searchResults"), jen.ID("err")).Assign().ID("sm").Dot("index").Dot("SearchInContext").Call(
				jen.ID("ctx"),
				jen.Qual(constants.SearchLibrary, "NewSearchRequest").Call(jen.ID("q")),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("performing search query"),
				),
				),
			),
			jen.Newline(),
			jen.For(jen.List(jen.ID("_"), jen.ID("result")).Assign().Range().ID("searchResults").Dot("Hits")).Body(
				jen.List(jen.ID("x"), jen.ID("parseErr")).Assign().Qual("strconv", "ParseUint").Call(
					jen.ID("result").Dot("ID"),
					jen.ID("base"),
					jen.ID("bitSize"),
				),
				jen.If(jen.ID("parseErr").Op("!=").ID("nil")).Body(
					jen.Comment("this should literally never happen"),
					jen.Return().List(jen.ID("nil"), jen.ID("observability").Dot("PrepareError").Call(
						jen.ID("parseErr"),
						jen.ID("logger"),
						jen.ID("span"),
						jen.Lit("parsing integer stored in search index for #%s"),
						jen.ID("result").Dot("ID"),
					),
					),
				),
				jen.Newline(),
				jen.ID("ids").Equals().ID("append").Call(
					jen.ID("ids"),
					jen.ID("x"),
				),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("ids"), jen.ID("nil")),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("Search implements our IndexManager interface."),
		jen.Newline(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("Search").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").String(), jen.ID("accountID").ID("uint64")).Params(jen.ID("ids").Index().Uint64(), jen.ID("err").ID("error")).Body(
			jen.Return().ID("sm").Dot("search").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.ID("accountID"),
				jen.ID("false"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("SearchForAdmin implements our IndexManager interface."),
		jen.Newline(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("SearchForAdmin").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("query").ID("string")).Params(jen.ID("ids").Index().Uint64(), jen.ID("err").ID("error")).Body(
			jen.Return().ID("sm").Dot("search").Call(
				jen.ID("ctx"),
				jen.ID("query"),
				jen.Lit(0),
				jen.ID("true"),
			),
		),
		jen.Newline(),
	)

	code.Add(
		jen.Comment("Delete implements our IndexManager interface."),
		jen.Newline(),
		jen.Func().Params(jen.ID("sm").PointerTo().ID("bleveIndexManager")).ID("Delete").Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("id").ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("sm").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("logger").Assign().ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("sm").Dot("index").Dot("Delete").Call(jen.Qual("strconv", "FormatUint").Call(
				jen.ID("id"),
				jen.ID("base"),
			),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().ID("observability").Dot("PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("removing from index"),
				),
			),
			jen.Newline(),
			jen.ID("sm").Dot("logger").Dot("WithValue").Call(
				jen.Lit("id"),
				jen.ID("id"),
			).Dot("Debug").Call(jen.Lit("removed from index")),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	)

	return code
}
