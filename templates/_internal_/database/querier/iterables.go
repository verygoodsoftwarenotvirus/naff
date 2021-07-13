package querier

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func buildScanFields(typ models.DataType) (scanFields []jen.Code) {
	scanFields = []jen.Code{
		jen.AddressOf().ID("x").Dot("ID"),
		jen.AddressOf().ID("x").Dot("ExternalID"),
	}

	for _, field := range typ.Fields {
		scanFields = append(scanFields, jen.AddressOf().ID("x").Dot(field.Name.Singular()))
	}

	scanFields = append(scanFields,
		jen.AddressOf().ID("x").Dot("CreatedOn"),
		jen.AddressOf().ID("x").Dot("LastUpdatedOn"),
		jen.AddressOf().ID("x").Dot("ArchivedOn"),
	)

	if typ.BelongsToStruct != nil {
		scanFields = append(scanFields, jen.AddressOf().ID("x").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		scanFields = append(scanFields, jen.AddressOf().ID("x").Dot(constants.AccountOwnershipFieldName))
	}

	return scanFields
}

func buildScanSomethingRow(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("scan%s takes a database Scanner (i.e. *sql.Row) and scans the result into %s struct.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("scan%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").Qual(proj.DatabasePackage(), "Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("x").PointerTo().Qual(proj.TypesPackage(), sn), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).Uint64(), jen.Err().ID("error")).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.Newline(),
			jen.ID("x").Equals().AddressOf().Qual(proj.TypesPackage(), sn).Values(),
			jen.Newline(),
			jen.ID("targetVars").Assign().Index().Interface().Valuesln(
				buildScanFields(typ)...,
			),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.ID("targetVars").Equals().ID("append").Call(
					jen.ID("targetVars"),
					jen.AddressOf().ID("filteredCount"),
					jen.AddressOf().ID("totalCount"),
				),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit(""),
				)),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("x"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Newline(),
	}
}

func buildScanListOfSomethingRows(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("scan%s takes some database rows and turns them into a slice of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("scan%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").Qual(proj.DatabasePackage(), "ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID(puvn).Index().PointerTo().Qual(proj.TypesPackage(), sn), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).Uint64(), jen.Err().ID("error")).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Lit("include_counts"),
				jen.ID("includeCounts"),
			),
			jen.Newline(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Body(
				jen.List(jen.ID("x"), jen.ID("fc"), jen.ID("tc"), jen.ID("scanErr")).Assign().ID("q").Dotf("scan%s", sn).Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("includeCounts"),
				),
				jen.If(jen.ID("scanErr").DoesNotEqual().ID("nil")).Body(
					jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.ID("scanErr")),
				),
				jen.Newline(),
				jen.If(jen.ID("includeCounts")).Body(
					jen.If(jen.ID("filteredCount").IsEqualTo().Lit(0)).Body(
						jen.ID("filteredCount").Equals().ID("fc"),
					),
					jen.Newline(),
					jen.If(jen.ID("totalCount").IsEqualTo().Lit(0)).Body(
						jen.ID("totalCount").Equals().ID("tc"),
					)),
				jen.Newline(),
				jen.ID(puvn).Equals().ID("append").Call(
					jen.ID(puvn),
					jen.ID("x"),
				),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("q").Dot("checkRowsForErrorAndClose").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Lit(0), jen.Lit(0), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("handling rows"),
				)),
			),
			jen.Newline(),
			jen.Return().List(jen.ID(puvn), jen.ID("filteredCount"), jen.ID("totalCount"), jen.ID("nil")),
		),
		jen.Newline(),
	}
}

func buildIDBoilerplate(proj *models.Project, typ models.DataType, includeType bool, returnVal jen.Code) []jen.Code {
	lines := []jen.Code{}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().Zero()).Body(
				jen.Return(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	if includeType {
		lines = append(lines,
			jen.If(jen.IDf("%sID", typ.Name.UnexportedVarName()).IsEqualTo().Lit(0)).Body(
				jen.Return().List(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		lines = append(lines,
			jen.If(jen.ID("accountID").IsEqualTo().Lit(0)).Body(
				jen.Return().List(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), "AccountIDKey"), jen.ID("accountID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("accountID")),
			jen.Newline(),
		)
	}

	return lines
}

func buildDBRetrievalQueryBuilderArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.RestrictedToAccountAtSomeLevel(p) {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...))
	}

	return params
}

func buildSomethingExists(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true, jen.False())...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("Build%sExistsQuery", sn).Call(
			buildDBRetrievalQueryBuilderArgs(proj, typ)...,
		),
		jen.Newline(),
		jen.List(jen.ID("result"), jen.Err()).Assign().ID("q").Dot("performBooleanQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.ID("query"),
			jen.ID("args"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("false"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("performing %s existence check", scn),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("result"), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("%sExists fetches whether %s exists from the database.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("%sExists", sn).Params(typ.BuildDBClientExistenceMethodParams(proj)...).Params(jen.ID("exists").ID("bool"), jen.Err().ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildGetSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true, jen.Nil())...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGet%sQuery", sn).Call(
			buildDBRetrievalQueryBuilderArgs(proj, typ)...,
		),
		jen.ID("row").Assign().ID("q").Dot("getOneRow").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Lit(uvn),
			jen.ID("query"),
			jen.ID("args").Spread(),
		),
		jen.Newline(),
		jen.List(jen.ID(uvn), jen.Underscore(), jen.Underscore(), jen.Err()).Assign().ID("q").Dotf("scan%s", sn).Call(
			jen.ID("ctx"),
			jen.ID("row"),
			jen.ID("false"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("scanning %s", scn),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID(uvn), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the database.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Get%s", sn).Params(typ.BuildDBClientRetrievalMethodParams(proj)...).Params(jen.PointerTo().Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildGetAllSomethingCount(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("GetAll%sCount fetches the count of %s from the database that meet a particular filter.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Uint64(), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
			jen.Newline(),
			jen.List(jen.ID("count"), jen.Err()).Assign().ID("q").Dot("performCountQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.ID("q").Dot("sqlQueryBuilder").Dotf("BuildGetAll%sCountQuery", pn).Call(jen.ID("ctx")),
				jen.Litf("fetching count of %s", pcn),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.Lit(0), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("querying for count of %s", pcn),
				)),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("count"), jen.ID("nil")),
		),
		jen.Newline(),
	}
}

func buildGetAllSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.If(jen.ID("results").IsEqualTo().ID("nil")).Body(
			jen.Return().ID("ErrNilInputProvided"),
		),
		jen.Newline(),
		constants.LoggerVar().Assign().ID("q").Dot("logger").Dot("WithValue").Call(
			jen.Lit("batch_size"),
			jen.ID("batchSize"),
		),
		jen.Newline(),
		jen.List(jen.ID("count"), jen.Err()).Assign().ID("q").Dotf("GetAll%sCount", pn).Call(jen.ID("ctx")),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("fetching count of %s", pcn),
			),
		),
		jen.Newline(),
		jen.For(jen.ID("beginID").Assign().Uint64().Call(jen.Lit(1)), jen.ID("beginID").Op("<=").ID("count"), jen.ID("beginID").Op("+=").Uint64().Call(jen.ID("batchSize"))).Body(
			jen.ID("endID").Assign().ID("beginID").Plus().Uint64().Call(jen.ID("batchSize")),
			jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).Uint64()).Body(
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGetBatchOf%sQuery", pn).Call(
					jen.ID("ctx"),
					jen.ID("begin"),
					jen.ID("end"),
				),
				constants.LoggerVar().Equals().ID("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(jen.Lit("query").MapAssign().ID("query"), jen.Lit("begin").MapAssign().ID("begin"), jen.Lit("end").MapAssign().ID("end"))),
				jen.Newline(),
				jen.List(jen.ID("rows"), jen.ID("queryErr")).Assign().ID("q").Dot("db").Dot("Query").Call(
					jen.ID("query"),
					jen.ID("args").Spread(),
				),
				jen.If(jen.Qual("errors", "Is").Call(
					jen.ID("queryErr"),
					jen.Qual("database/sql", "ErrNoRows"),
				)).Body(
					jen.Return()).Else().If(jen.ID("queryErr").DoesNotEqual().ID("nil")).Body(
					constants.LoggerVar().Dot("Error").Call(
						jen.ID("queryErr"),
						jen.Lit("querying for database rows"),
					),
					jen.Return(),
				),
				jen.Newline(),
				jen.List(jen.ID(puvn), jen.Underscore(), jen.Underscore(), jen.ID("scanErr")).Assign().ID("q").Dotf("scan%s", pn).Call(
					jen.ID("ctx"),
					jen.ID("rows"),
					jen.ID("false"),
				),
				jen.If(jen.ID("scanErr").DoesNotEqual().ID("nil")).Body(
					constants.LoggerVar().Dot("Error").Call(
						jen.ID("scanErr"),
						jen.Lit("scanning database rows"),
					),
					jen.Return(),
				),
				jen.Newline(),
				jen.ID("results").ReceiveFromChannel().ID(puvn),
			).Call(
				jen.ID("beginID"),
				jen.ID("endID"),
			),
		),
		jen.Newline(),
		jen.Return().ID("nil"),
	}

	return []jen.Code{
		jen.Commentf("GetAll%s fetches a list of all %s in the database.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("GetAll%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("results").Chan().Index().PointerTo().Qual(proj.TypesPackage(), sn), jen.ID("batchSize").ID("uint16")).Params(jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildListRetrievalQueryBuildingMethodArgs(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToAccountAtSomeLevel(p) {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...))
	}

	params = append(params, jen.False(), jen.ID("filter"))

	return params
}

func buildGetListOfSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, false, jen.Nil())...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("x").Equals().AddressOf().ID("types").Dotf("%sList", sn).Values(),
		constants.LoggerVar().Equals().ID("filter").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Qual(proj.InternalTracingPackage(), "AttachQueryFilterToSpan").Call(
			jen.ID("span"),
			jen.ID("filter"),
		),
		jen.Newline(),
		jen.If(jen.ID("filter").DoesNotEqual().ID("nil")).Body(
			jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Equals().List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit")),
		),
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGet%sQuery", pn).Call(
			buildListRetrievalQueryBuildingMethodArgs(proj, typ)...,
		),
		jen.Newline(),
		jen.List(jen.ID("rows"), jen.Err()).Assign().ID("q").Dot("performReadQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Lit(puvn),
			jen.ID("query"),
			jen.ID("args").Spread(),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("executing %s list retrieval query", pcn),
			)),
		),
		jen.Newline(),
		jen.If(jen.List(jen.ID("x").Dot(pn), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.Err()).Equals().ID("q").Dotf("scan%s", pn).Call(
			jen.ID("ctx"),
			jen.ID("rows"),
			jen.ID("true"),
		), jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("scanning %s", pcn),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("x"), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Get%s", pn).Params(typ.BuildDBClientListRetrievalMethodParams(proj)...).Params(jen.ID("x").PointerTo().ID("types").Dotf("%sList", sn), jen.Err().ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildIDBoilerplateForIDsQuery(proj *models.Project, typ models.DataType, returnVal jen.Code) []jen.Code {
	lines := []jen.Code{}

	if typ.BelongsToStruct != nil {
		lines = append(lines,
			jen.If(jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).IsEqualTo().Zero()).Body(
				jen.Return(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.BelongsToStruct.Singular())), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Newline(),
		)
	}

	return lines
}

func buildGetListOfSomethingWithIDs(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	pcn := typ.Name.PluralCommonName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplateForIDsQuery(proj, typ, jen.Nil())...)

	bodyLines = append(bodyLines,
		utils.ConditionalCode(typ.BelongsToAccount, jen.If(jen.ID("accountID").IsEqualTo().Lit(0)).Body(
			jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
		)),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), "AccountIDKey"), jen.ID("accountID"))),
		utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
			jen.ID("span"),
			jen.ID("accountID"),
		)),
		jen.Newline(),
		jen.If(jen.ID("limit").IsEqualTo().Lit(0)).Body(
			jen.ID("limit").Equals().ID("uint8").Call(jen.Qual(proj.TypesPackage(), "DefaultLimit")),
		),
		jen.Newline(),
		constants.LoggerVar().Equals().ID("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
			jen.Lit("limit").MapAssign().ID("limit"),
			jen.Lit("id_count").MapAssign().ID("len").Call(jen.ID("ids")),
		)),
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGet%sWithIDsQuery", pn).Call(
			jen.ID("ctx"),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("accountID")),
			jen.ID("limit"),
			jen.ID("ids"),
			func() jen.Code {
				if typ.BelongsToAccount {
					if typ.RestrictedToAccountAtSomeLevel(proj) {
						return jen.True()
					} else {
						return jen.False()
					}
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
		jen.List(jen.ID("rows"), jen.Err()).Assign().ID("q").Dot("performReadQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Litf("%s with IDs", pcn),
			jen.ID("query"),
			jen.ID("args").Spread(),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("fetching %s from database", pcn),
			)),
		),
		jen.Newline(),
		jen.List(jen.ID(puvn), jen.Underscore(), jen.Underscore(), jen.Err()).Assign().ID("q").Dotf("scan%s", pn).Call(
			jen.ID("ctx"),
			jen.ID("rows"),
			jen.ID("false"),
		),
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
			jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.Err(),
				constants.LoggerVar(),
				jen.ID("span"),
				jen.Litf("scanning %s", pcn),
			)),
		),
		jen.Newline(),
		jen.Return().List(jen.ID(puvn), jen.ID("nil")),
	)

	return []jen.Code{
		jen.Commentf("Get%sWithIDs fetches %s from the database within a given set of IDs.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Get%sWithIDs", pn).Params(typ.BuildGetListOfSomethingFromIDsParams(proj)...).Params(jen.Index().PointerTo().Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}
}

func buildCreateInitFields(typ models.DataType) []jen.Code {
	createInitColumns := []jen.Code{
		jen.ID("ID").MapAssign().ID("id"),
	}

	queryBuildingArgs := typ.BuildDBQuerierCreationMethodQueryBuildingArgs()
	queryBuildingArgs = queryBuildingArgs[:len(queryBuildingArgs)-1]
	queryBuildingArgs = append(queryBuildingArgs, jen.ID("x"))

	for _, field := range typ.Fields {
		fn := field.Name.Singular()
		createInitColumns = append(createInitColumns, jen.ID(fn).MapAssign().ID("input").Dot(fn))
	}

	if typ.BelongsToStruct != nil {
		createInitColumns = append(createInitColumns, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		createInitColumns = append(createInitColumns, jen.ID(constants.AccountOwnershipFieldName).MapAssign().ID("input").Dot(constants.AccountOwnershipFieldName))
	}

	createInitColumns = append(createInitColumns, jen.ID("CreatedOn").MapAssign().ID("q").Dot("currentTime").Call())

	return createInitColumns
}

func buildCreateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Create%s", sn).Params(typ.BuildDBClientCreationMethodParams(proj)...).Params(jen.PointerTo().Qual(proj.TypesPackage(), sn), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("input").IsEqualTo().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrNilInputProvided")),
			),
			jen.Newline(),
			jen.If(jen.ID("createdByUser").IsEqualTo().Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.Qual(proj.ConstantKeysPackage(), "RequesterIDKey"),
				jen.ID("createdByUser"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("createdByUser"),
			),
			jen.Newline(),
			jen.List(jen.ID("tx"), jen.Err()).Assign().ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				)),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildCreate%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.ID("input"),
			),
			jen.Newline(),
			jen.Commentf("create the %s.", scn),
			jen.List(jen.ID("id"), jen.Err()).Assign().ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.ID("false"),
				jen.Litf("%s creation", scn),
				jen.ID("query"),
				jen.ID("args"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("creating %s", scn),
				)),
			),
			jen.Newline(),
			jen.ID("x").Assign().AddressOf().Qual(proj.TypesPackage(), sn).Valuesln(
				buildCreateInitFields(typ)...,
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sCreationEventEntry", sn)).Call(
					jen.ID("x"),
					jen.ID("createdByUser"),
				),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("writing %s creation audit log entry", scn),
				)),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("tx").Dot("Commit").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				)),
			),
			jen.Newline(),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.ID("x").Dot("ID"),
			),
			constants.LoggerVar().Dot("Info").Call(jen.Litf("%s created", scn)),
			jen.Newline(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Newline(),
	}
}

func buildUpdateSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the provided input to have a valid ID.", sn, scn, sn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Update%s", sn).Params(typ.BuildDBClientUpdateMethodParams(proj, "updated")...).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("updated").IsEqualTo().ID("nil")).Body(
				jen.Return().ID("ErrNilInputProvided"),
			),
			jen.Newline(),
			jen.If(jen.ID("changedByUser").IsEqualTo().Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.Newline(),
			constants.LoggerVar().Assign().ID("q").Dot("logger").Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.ID("updated").Dot("ID"),
			),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.ID("updated").Dot("ID"),
			),
			utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
				jen.ID("span"),
				jen.ID("updated").Dot("BelongsToAccount"),
			)),
			jen.Qual(proj.InternalTracingPackage(), "AttachRequestingUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("changedByUser"),
			),
			jen.Newline(),
			jen.List(jen.ID("tx"), jen.Err()).Assign().ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildUpdate%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.ID("updated"),
			),
			jen.If(jen.Err().Equals().ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Litf("%s update", scn),
				jen.ID("query"),
				jen.ID("args"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("updating %s", scn),
				),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sUpdateEventEntry", sn)).Call(
					jen.ID("changedByUser"),
					jen.ID("updated").Dot("ID"),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("updated").Dot("BelongsToAccount")),
					jen.ID("changes"),
				),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("writing %s update audit log entry", scn),
				),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("tx").Dot("Commit").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				),
			),
			jen.Newline(),
			constants.LoggerVar().Dot("Info").Call(jen.Litf("%s updated", scn)),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	}
}

func buildArchiveSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("Archive%s archives %s from the database by its ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Archive%s", sn).Params(typ.BuildDBClientArchiveMethodParams()...).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).IsEqualTo().Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.If(jen.ID("accountID").IsEqualTo().Lit(0)).Body(jen.Return().ID("ErrInvalidIDProvided"))),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), "AccountIDKey"), jen.ID("accountID"))),
			utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("accountID"))),
			jen.Newline(),
			jen.If(jen.ID("archivedBy").IsEqualTo().Lit(0)).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), "RequesterIDKey"), jen.ID("archivedBy")),
			jen.Qual(proj.InternalTracingPackage(), "AttachUserIDToSpan").Call(
				jen.ID("span"),
				jen.ID("archivedBy"),
			),
			jen.Newline(),
			jen.List(jen.ID("tx"), jen.Err()).Assign().ID("q").Dot("db").Dot("BeginTx").Call(
				jen.ID("ctx"),
				jen.ID("nil"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("beginning transaction"),
				),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildArchive%sQuery", sn).Call(
				jen.ID("ctx"),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
					}
					return jen.Null()
				}(),
				jen.IDf("%sID", uvn),
				utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.ID("accountID")),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("q").Dot("performWriteQueryIgnoringReturn").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Litf("%s archive", scn),
				jen.ID("query"),
				jen.ID("args"),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("updating %s", scn),
				),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("q").Dot("createAuditLogEntryInTransaction").Call(
				jen.ID("ctx"),
				jen.ID("tx"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("Build%sArchiveEventEntry", sn)).Call(
					jen.ID("archivedBy"),
					utils.ConditionalCode(typ.BelongsToAccount, jen.ID("accountID")),
					jen.IDf("%sID", uvn),
				),
			), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.ID("q").Dot("rollbackTransaction").Call(
					jen.ID("ctx"),
					jen.ID("tx"),
				),
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Litf("writing %s archive audit log entry", scn),
				),
			),
			jen.Newline(),
			jen.If(jen.Err().Equals().ID("tx").Dot("Commit").Call(), jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("committing transaction"),
				),
			),
			jen.Newline(),
			constants.LoggerVar().Dot("Info").Call(jen.Litf("%s archived", scn)),
			jen.Newline(),
			jen.Return().ID("nil"),
		),
		jen.Newline(),
	}
}

func buildGetAuditLogEntriesForSomething(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	return []jen.Code{
		jen.Commentf("GetAuditLogEntriesFor%s fetches a list of audit log entries from the database that relate to a given %s.", sn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("GetAuditLogEntriesFor%s", sn).Params(constants.CtxParam(), jen.IDf("%sID", uvn).Uint64()).Params(jen.Index().PointerTo().Qual(proj.TypesPackage(), "AuditLogEntry"), jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).IsEqualTo().Lit(0)).Body(
				jen.Return().List(jen.ID("nil"), jen.ID("ErrInvalidIDProvided")),
			),
			constants.LoggerVar().Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(
				jen.ID("keys").Dotf("%sIDKey", sn),
				jen.IDf("%sID", uvn),
			),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", sn)).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("sqlQueryBuilder").Dotf("BuildGetAuditLogEntriesFor%sQuery", sn).Call(
				jen.ID("ctx"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID("q").Dot("performReadQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Litf("audit log entries for %s", scn),
				jen.ID("query"),
				jen.ID("args").Spread(),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("querying database for audit log entries"),
				)),
			),
			jen.Newline(),
			jen.List(jen.ID("auditLogEntries"), jen.Underscore(), jen.Err()).Assign().ID("q").Dot("scanAuditLogEntries").Call(
				jen.ID("ctx"),
				jen.ID("rows"),
				jen.ID("false"),
			),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Body(
				jen.Return().List(jen.ID("nil"), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.Err(),
					constants.LoggerVar(),
					jen.ID("span"),
					jen.Lit("scanning audit log entries"),
				)),
			),
			jen.Newline(),
			jen.Return().List(jen.ID("auditLogEntries"), jen.ID("nil")),
		),
		jen.Newline(),
	}
}

func newIterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	sn := typ.Name.Singular()

	utils.AddImports(proj, code, false)

	code.Add(
		jen.Var().Defs(
			jen.Underscore().Qual(proj.TypesPackage(), fmt.Sprintf("%sDataManager", sn)).Equals().Parens(jen.PointerTo().ID("SQLQuerier")).Call(jen.ID("nil")),
		),
		jen.Newline(),
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
