package querier

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildScanSomethingRow(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("scan%s takes a database Scanner (i.e. *sql.Row) and scans the result into %s struct.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("scan%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").Qual(proj.DatabasePackage(), "Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("x").Op("*").Qual(proj.TypesPackage(), sn), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.Line(),
			jen.ID("x").Op("=").Op("&").Qual(proj.TypesPackage(), sn).Values(),
			jen.Line(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("ExternalID"), jen.Op("&").ID("x").Dot("Name"), jen.Op("&").ID("x").Dot("Details"), jen.Op("&").ID("x").Dot("CreatedOn"), jen.Op("&").ID("x").Dot("LastUpdatedOn"), jen.Op("&").ID("x").Dot("ArchivedOn"), jen.Op("&").ID("x").Dot("BelongsToAccount")),
			jen.Line(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Op("=").ID("append").Call(
					jen.ID("targetVars"),
					jen.Op("&").ID("filteredCount"),
					jen.Op("&").ID("totalCount"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit(""),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildScanListOfSomethingRows(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("scan%s takes some database rows and turns them into a slice of %s.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("scan%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").Qual(proj.DatabasePackage(), "ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID(puvn).Index().Op("*").Qual(proj.TypesPackage(), sn), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).ID("uint64"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("x"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Assign().ID("q").Dotf("scan%s", sn).Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").Op("!=").ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("scanErr")),
				),
				jen.Line(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("filteredCount").Op("==").Lit(0)).Body(
						jen.ID("filteredCount").Op("=").ID("fc"),
					),
					jen.Line(),
					jen.If(jen.ID("totalCount").Op("==").Lit(0)).Body(
						jen.ID("totalCount").Op("=").ID("tc"),
					)),
				jen.Line(),
				jen.ID(puvn).Op("=").ID("append").Call(
					jen.ID(puvn),
					jen.ID("x"),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("handling rows"),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID(puvn), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("%sExists fetches whether %s exists from the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("%sExists", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("accountID")).ID("uint64")).Params(jen.ID("exists").ID("bool"), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("false"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("Build%sExistsQuery", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("accountID"),
			),
			jen.Line(),
			jen.List(jen.ID("result"), jen.ID("err")).Assign().ID("q").Dot("performBooleanQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("false"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("performing %s existence check", scn),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID("result"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Get%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("accountID")).ID("uint64")).Params(jen.Op("*").Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("keys").Dotf("%sIDKey", sn).Op(":").IDf("%sID", uvn),
				jen.ID("keys").Dot("UserIDKey").Op(":").ID("accountID"),
			),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGet%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("accountID"),
			),
			jen.ID("row").Assign().ID("q").Dot("getOneRow").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit(uvn),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.Line(),
			jen.List(jen.ID(uvn), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", sn).Call(
				jen.ID("ctx"),
				jen.ID("row"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("scanning %s", scn),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID(uvn), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildGetAllSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sCount fetches the count of %s from the database that meet a particular filter.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("uint64"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger"),
			jen.Line(),
			jen.List(jen.ID("count"), jen.ID("err")).Assign().ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dotf("BuildGetAll%sCountQuery", pn).Call(jen.ID("ctx")),
				jen.Litf("fetching count of %s", pcn),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("querying for count of %s", pcn),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildGetAllSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%s fetches a list of all %s in the database.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("GetAll%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().Op("*").Qual(proj.TypesPackage(), sn), jen.ID("batchSize").ID("uint16")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.ID("results").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided"),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("batch_size"),
				jen.ID("batchSize"),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.ID("err")).Assign().ID("q").Dotf("GetAll%sCount", pn).Call(jen.ID("ctx")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("fetching count of %s", pcn),
				),
			),
			jen.Line(),
			jen.For(jen.ID("beginID").Assign().ID("uint64").Call(jen.Lit(1)), jen.ID("beginID").Op("<=").ID("count"), jen.ID("beginID").Op("+=").ID("uint64").Call(jen.ID("batchSize"))).Body(
				jen.ID("endID").Assign().ID("beginID").Op("+").ID("uint64").Call(jen.ID("batchSize")),
				jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).ID("uint64")).Body(
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGetBatchOf%sQuery", pn).Call(
						jen.ID("ctx"),
						jen.ID("begin"),
						jen.ID("end"),
					),
					jen.ID("logger").Op("=").ID("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(jen.Lit("query").Op(":").ID("query"), jen.Lit("begin").Op(":").ID("begin"), jen.Lit("end").Op(":").ID("end"))),
					jen.Line(),
					jen.List(jen.ID("rows"), jen.ID("queryErr")).Assign().ID("q").Dot("db").Dot("Query").Call(
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
					jen.Line(),
					jen.List(jen.ID(puvn), jen.ID("_"), jen.ID("_"), jen.ID("scanErr")).Assign().ID("q").Dotf("scan%s", pn).Call(
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
					jen.Line(),
					jen.ID("results").ReceiveFromChannel().ID(puvn),
				).Call(
					jen.ID("beginID"),
					jen.ID("endID"),
				),
			),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}
}

func buildGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Get%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("filter").Op("*").Qual(proj.TypesPackage(), "QueryFilter")).Params(jen.ID("x").Op("*").ID("types").Dotf("%sList", sn), jen.ID("err").ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.ID("x").Op("=").Op("&").ID("types").Dotf("%sList", sn).Values(),
			jen.ID("logger").Assign().ID("filter").Dot("AttachToLogger").Call(jen.ID("q").Dot("logger")).Dot("WithValue").Call(
				jen.ID("keys").Dot("AccountIDKey"),
				jen.ID("accountID"),
			),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachQueryFilterToSpan").Call(
				jen.ID("span"),
				jen.ID("filter"),
			),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Op("=").List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGet%sQuery", pn).Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("false"),
				jen.ID("filter"),
			),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Assign().ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Lit(puvn),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("executing %s list retrieval query", pcn),
				)),
			),
			jen.Line(),
			jen.If(jen.List(jen.ID("x").Dot(pn), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Op("=").ID("q").Dotf("scan%s", pn).Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("true"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("scanning %s", pcn),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildGetListOfSomethingWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("Get%sWithIDs fetches %s from the database within a given set of IDs.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Get%sWithIDs", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("accountID").ID("uint64"), jen.ID("limit").ID("uint8"), jen.ID("ids").Index().ID("uint64")).Params(jen.Index().Op("*").Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Line(),
			jen.If(jen.ID("limit").Op("==").Lit(0)).Body(
				jen.ID("limit").Op("=").ID("uint8").Call(jen.Qual(proj.TypesPackage(), "DefaultLimit")),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("keys").Dot("UserIDKey").Op(":").ID("accountID"),
				jen.Lit("limit").Op(":").ID("limit"),
				jen.Lit("id_count").Op(":").ID("len").Call(jen.ID("ids")),
			)),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGet%sWithIDsQuery", pn).Call(
				jen.ID("ctx"),
				jen.ID("accountID"),
				jen.ID("limit"),
				jen.ID("ids"),
				jen.ID("false"),
			),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Assign().ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Litf("%s with IDs", pcn),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("fetching %s from database", pcn),
				)),
			),
			jen.Line(),
			jen.List(jen.ID(puvn), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", pn).Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("scanning %s", pcn),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID(puvn), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Create%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").ID("types").Dotf("%sCreationInput", sn), jen.ID("createdByUser").ID("uint64")).Params(jen.Op("*").Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.ID("input").Op("==").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided")),
			),
			jen.Line(),
			jen.If(jen.ID("createdByUser").Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dot("RequesterIDKey"),
				jen.ID("createdByUser"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("createdByUser"),
			),
			jen.Line(),
			jen.List(jen.ID("tx"), jen.ID("err")).Assign().ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildCreate%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Line(),
			jen.Commentf("create the %s.", scn),
			jen.List(jen.ID("id"), jen.ID("err")).Assign().ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("false"),
				jen.Litf("%s creation", scn),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("creating %s", scn),
				)),
			),
			jen.Line(),
			jen.ID("x").Assign().Op("&").Qual(proj.TypesPackage(), sn).Valuesln(jen.ID("ID").Op(":").ID("id"), jen.ID("Name").Op(":").ID("input").Dot("Name"), jen.ID("Details").Op(":").ID("input").Dot("Details"), jen.ID("BelongsToAccount").Op(":").ID("input").Dot("BelongsToAccount"), jen.ID("CreatedOn").Op(":").ID("q").Dot("currentTime").Call()),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dotf("Build%sCreationEventEntry", sn).Call(
					jen.ID("x"),
					jen.ID("createdByUser"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("writing %s creation audit log entry", scn),
				)),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.ID("x").Dot("ID"),
			),
			jen.ID("logger").Dot("Info").Call(jen.Litf("%s created", scn)),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the provided input to have a valid ID.", sn, scn, sn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("updated").Op("*").Qual(proj.TypesPackage(), sn), jen.ID("changedByUser").ID("uint64"), jen.ID("changes").Index().Op("*").Qual(proj.TypesPackage(), "FieldChangeSummary")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.ID("updated").Op("==").ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided"),
			),
			jen.Line(),
			jen.If(jen.ID("changedByUser").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.ID("updated").Dot("ID"),
			),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.ID("updated").Dot("ID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("BelongsToAccount"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("changedByUser"),
			),
			jen.Line(),
			jen.List(jen.ID("tx"), jen.ID("err")).Assign().ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildUpdate%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Litf("%s update", scn),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("updating %s", scn),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dotf("Build%sUpdateEventEntry", sn).Call(
					jen.ID("changedByUser"),
					jen.ID("updated").Dot("ID"),
					jen.ID("updated").Dot("BelongsToAccount"),
					jen.ID("changes"),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("writing %s update audit log entry", scn),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				),
			),
			jen.Line(),
			jen.ID("logger").Dot("Info").Call(jen.Litf("%s updated", scn)),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}
}

func buildArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("Archive%s archives %s from the database by its ID.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Archive%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("accountID"), jen.ID("archivedBy")).ID("uint64")).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Line(),
			jen.If(jen.ID("accountID").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Line(),
			jen.If(jen.ID("archivedBy").Op("==").Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Line(),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("accountID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("archivedBy"),
			),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValues").Call(jen.Map(jen.ID("string")).Interface().Valuesln(
				jen.ID("keys").Dotf("%sIDKey", sn).Op(":").IDf("%sID", uvn),
				jen.ID("keys").Dot("UserIDKey").Op(":").ID("archivedBy"),
				jen.ID("keys").Dot("AccountIDKey").Op(":").ID("accountID"),
			)),
			jen.Line(),
			jen.List(jen.ID("tx"), jen.ID("err")).Assign().ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildArchive%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
				jen.ID("accountID"),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Litf("%s archive", scn),
				jen.ID("query"),
				jen.ID("args"),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("updating %s", scn),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("audit").Dotf("Build%sArchiveEventEntry", sn).Call(
					jen.ID("archivedBy"),
					jen.ID("accountID"),
					jen.IDf("%sID", uvn),
				),
			), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("writing %s archive audit log entry", scn),
				),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op("=").ID("tx").Dot("Commit").Call(), jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				),
			),
			jen.Line(),
			jen.ID("logger").Dot("Info").Call(jen.Litf("%s archived", scn)),
			jen.Line(),
			jen.Return().ID("nil"),
		),
		jen.Line(),
	}
}

func buildGetAuditLogEntriesForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Commentf("GetAuditLogEntriesFor%s fetches a list of audit log entries from the database that relate to a given %s.", sn, scn),
		jen.Line(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("GetAuditLogEntriesFor%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.IDf("%sID", uvn).ID("uint64")).Params(jen.Index().Op("*").Qual(proj.TypesPackage(), "AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Line(),
			jen.If(jen.IDf("%sID", uvn).Op("==").Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Line(),
			jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGetAuditLogEntriesFor%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Assign().ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Litf("audit log entries for %s", scn),
				jen.ID("query"),
				jen.ID("args").Op("..."),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("querying database for audit log entries"),
				)),
			),
			jen.Line(),
			jen.List(jen.ID("auditLogEntries"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dot("scanAuditLogEntries").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Lit("scanning audit log entries"),
				)),
			),
			jen.Line(),
			jen.Return().List(jen.ID("auditLogEntries"), jen.ID("nil")),
		),
		jen.Line(),
	}
}

func newIterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	sn := typ.Name.Singular()

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.ID("_").ID("types").Dotf("%sDataManager", sn).Op("=").Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Line(),
	)

	code.Add(buildScanSomethingRow(proj, typ)...)
	code.Add(buildScanListOfSomethingRows(proj, typ)...)
	code.Add(buildSomethingExists(proj, typ)...)
	code.Add(buildGetSomething(proj, typ)...)
	code.Add(buildGetAllSomethingCount(proj, typ)...)
	code.Add(buildGetAllSomething(proj, typ)...)
	code.Add(buildGetListOfSomething(proj, typ)...)
	code.Add(buildGetListOfSomethingWithIDs(proj, typ)...)
	code.Add(buildCreateSomething(proj, typ)...)
	code.Add(buildUpdateSomething(proj, typ)...)
	code.Add(buildArchiveSomething(proj, typ)...)
	code.Add(buildGetAuditLogEntriesForSomething(proj, typ)...)

	return code
}
