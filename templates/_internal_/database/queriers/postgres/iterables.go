package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTopConstBlock(proj, typ, dbvendor)...)
	code.Add(buildTopVarBlock(proj, typ, dbvendor)...)
	code.Add(buildScanSomething(proj, typ, dbvendor)...)
	code.Add(buildScanMultipleSomethings(proj, typ, dbvendor)...)
	code.Add(buildSomethingExists(proj, typ, dbvendor)...)
	code.Add(buildGetSomething(proj, typ, dbvendor)...)
	code.Add(buildGetTotalSomethingCount(proj, typ, dbvendor)...)
	code.Add(buildGetSomethingsList(proj, typ, dbvendor)...)
	code.Add(buildGetSomethingWithIDsQuery(proj, typ, dbvendor)...)
	code.Add(buildGetSomethingWithIDs(proj, typ, dbvendor)...)
	code.Add(buildCreateSomething(proj, typ, dbvendor)...)
	code.Add(buildUpdateSomething(proj, typ, dbvendor)...)
	code.Add(buildArchiveSomething(proj, typ, dbvendor)...)

	return code
}

func convertArgsToCode(args []interface{}) (code []jen.Code) {
	for _, arg := range args {
		if c, ok := arg.(models.Coder); ok {
			code = append(code, c.Code())
		}
	}
	return
}

func determineSelectColumns(typ models.DataType) []string {
	tableName := typ.Name.PluralRouteName()

	selectColumns := []string{
		fmt.Sprintf("%s.id", tableName),
	}

	for _, field := range typ.Fields {
		selectColumns = append(selectColumns, fmt.Sprintf("%s.%s", tableName, field.Name.RouteName()))
	}
	selectColumns = append(selectColumns,
		fmt.Sprintf("%s.created_on", tableName),
		fmt.Sprintf("%s.last_updated_on", tableName),
		fmt.Sprintf("%s.archived_on", tableName),
	)

	return selectColumns
}

func buildTopConstBlock(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	lines := []jen.Code{
		jen.Const().Defs(
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					prn := typ.Name.PluralRouteName()
					puvn := typ.Name.Plural()
					btspn := typ.BelongsToStruct.PluralUnexportedVarName()
					btsrn := typ.BelongsToStruct.RouteName()
					btsprn := typ.BelongsToStruct.PluralRouteName()

					return jen.IDf("%sOn%sJoinClause", btspn, puvn).Equals().Litf("%s ON %s.belongs_to_%s=%s.id", btsprn, prn, btsrn, btsprn)
				}
				return jen.Null()
			}(),
		),
	}

	return lines
}

func buildTopVarBlock(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	tableName := typ.Name.PluralRouteName()

	columns := []jen.Code{
		jen.Litf("%s.id", tableName),
	}

	for _, field := range typ.Fields {
		columns = append(columns, jen.Litf("%s.%s", tableName, field.Name.RouteName()))
	}

	columns = append(columns,
		jen.Litf("%s.created_on", tableName),
		jen.Litf("%s.last_updated_on", tableName),
		jen.Litf("%s.archived_on", tableName),
	)

	if typ.BelongsToStruct != nil {
		columns = append(columns, jen.Litf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName()))
	}

	if typ.BelongsToAccount {
		columns = append(columns, jen.Litf("%s.belongs_to_account", tableName))
	}

	lines := []jen.Code{
		jen.Var().Defs(
			jen.ID("_").ID("types").Dotf("%sDataManager", sn).Equals().Parens(jen.Op("*").ID("SQLQuerier")).Call(jen.Nil()),
			jen.Newline(),
			jen.Commentf("%sTableColumns are the columns for the %s table.", puvn, tableName),
			jen.IDf("%sTableColumns", puvn).Equals().Index().String().Valuesln(
				columns...,
			),
			jen.Newline(),
			utils.ConditionalCode(len(proj.FindOwnerTypeChain(typ)) > 0, jen.IDf("get%sJoins", pn).Equals().Index().String().Valuesln(
				buildJoinsListForListQuery(proj, typ)...,
			)),
		),
		jen.Newline(),
	}

	return lines
}

func buildIDBoilerplate(proj *models.Project, typ models.DataType, includeType bool, returnVal jen.Code) []jen.Code {
	lines := []jen.Code{}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines,
			jen.If(jen.IDf("%sID", dep.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", dep.Name.Singular())), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", dep.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", dep.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	if includeType {
		lines = append(lines,
			jen.If(jen.IDf("%sID", typ.Name.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return().List(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
		)
	}

	if typ.RestrictedToAccountAtSomeLevel(proj) {
		lines = append(lines,
			jen.If(jen.ID("accountID").IsEqualTo().EmptyString()).Body(
				jen.Return().List(returnVal, jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), "AccountIDKey"), jen.ID("accountID")),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("accountID")),
			jen.Newline(),
		)
	}

	return lines
}

// end helpers

func buildScanSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	columns := []jen.Code{
		jen.Op("&").ID("x").Dot("ID"),
	}

	for _, field := range typ.Fields {
		columns = append(columns, jen.Op("&").ID("x").Dot(field.Name.Singular()))
	}

	columns = append(columns,
		jen.Op("&").ID("x").Dot("CreatedOn"),
		jen.Op("&").ID("x").Dot("LastUpdatedOn"),
		jen.Op("&").ID("x").Dot("ArchivedOn"),
	)

	if typ.BelongsToStruct != nil {
		columns = append(columns, jen.Op("&").ID("x").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToAccount {
		columns = append(columns, jen.Op("&").ID("x").Dot("BelongsToAccount"))
	}

	bodyLines := []jen.Code{
		jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
			jen.Lit("include_counts"),
			jen.ID("includeCounts"),
		),
		jen.Newline(),
		jen.ID("x").Equals().Op("&").ID("types").Dotf(sn).Values(),
		jen.Newline(),
		jen.ID("targetVars").Assign().Index().Interface().Valuesln(columns...),
		jen.Newline(),
		jen.If(jen.ID("includeCounts")).Body(
			jen.ID("targetVars").Equals().ID("append").Call(
				jen.ID("targetVars"),
				jen.Op("&").ID("filteredCount"),
				jen.Op("&").ID("totalCount"),
			)),
		jen.Newline(),
		jen.If(jen.ID("err").Equals().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Op("...")), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Zero(), jen.Zero(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit(""),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID("x"), jen.ID("filteredCount"), jen.ID("totalCount"), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("scan%s takes a database Scanner (i.e. *sql.Row) and scans the result into %s struct.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("scan%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("scan").Qual(proj.DatabasePackage(), "Scanner"), jen.ID("includeCounts").ID("bool")).Params(jen.ID("x").Op("*").ID("types").Dot(sn), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).Uint64(), jen.ID("err").ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildScanMultipleSomethings(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("_"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
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
			jen.If(jen.ID("scanErr").DoesNotEqual().Nil()).Body(
				jen.Return().List(jen.Nil(), jen.Zero(), jen.Zero(), jen.ID("scanErr"))),
			jen.Newline(),
			jen.If(jen.ID("includeCounts")).Body(
				jen.If(jen.ID("filteredCount").IsEqualTo().Zero()).Body(
					jen.ID("filteredCount").Equals().ID("fc")),
				jen.Newline(),
				jen.If(jen.ID("totalCount").IsEqualTo().Zero()).Body(
					jen.ID("totalCount").Equals().ID("tc")),
			),
			jen.Newline(),
			jen.ID(puvn).Equals().ID("append").Call(
				jen.ID(puvn),
				jen.ID("x"),
			),
		),
		jen.Newline(),
		jen.If(jen.ID("err").Equals().ID("q").Dot("checkRowsForErrorAndClose").Call(
			jen.ID("ctx"),
			jen.ID("rows"),
		), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Zero(), jen.Zero(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Lit("handling rows"),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID(puvn), jen.ID("filteredCount"), jen.ID("totalCount"), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Commentf("scan%s takes some database rows and turns them into a slice of %s.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("scan%s", pn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("rows").Qual(proj.DatabasePackage(), "ResultIterator"), jen.ID("includeCounts").ID("bool")).Params(jen.ID(puvn).Index().Op("*").ID("types").Dot(sn), jen.List(jen.ID("filteredCount"), jen.ID("totalCount")).Uint64(), jen.ID("err").ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildSomethingExists(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	tableName := typ.Name.PluralRouteName()

	eqArgs := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): whatever,
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		eqArgs[fmt.Sprintf("%s.belongs_to_account", tableName)] = whatever
	}
	if typ.BelongsToStruct != nil {
		eqArgs[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whatever
	}

	whereValues := typ.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(proj)
	qb := queryBuilderForDatabase(dbvendor).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	query, args, err := qb.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger"),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true, jen.False())...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.List(jen.ID("result"), jen.ID("err")).Assign().ID("q").Dot("performBooleanQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.IDf("%sExistenceQuery", uvn),
			jen.ID("args"),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.False(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("performing %s existence check", scn),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID("result"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Const().IDf("%sExistenceQuery", uvn).Equals().Lit(query),
		jen.Newline(),
		jen.Commentf("%sExists fetches whether %s exists from the database.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("%sExists", sn).Params(typ.BuildDBClientExistenceMethodParams(proj)...).Params(jen.ID("exists").ID("bool"), jen.ID("err").ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildPrefixedStringColumns(typ models.DataType) []string {
	tableName := typ.Name.PluralRouteName()
	out := []string{fmt.Sprintf("%s.id", tableName)}

	for _, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s.%s", tableName, field.Name.RouteName()))
	}

	out = append(out,
		fmt.Sprintf("%s.created_on", tableName),
		fmt.Sprintf("%s.last_updated_on", tableName),
		fmt.Sprintf("%s.archived_on", tableName),
	)

	if typ.BelongsToStruct != nil {
		out = append(out, fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName()))
	}
	if typ.BelongsToAccount {
		out = append(out, fmt.Sprintf("%s.belongs_to_account", tableName))
	}

	return out
}

func buildGetSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	tableName := typ.Name.PluralRouteName()

	whereValues := typ.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(proj)

	cols := buildPrefixedStringColumns(typ)
	qb := queryBuilderForDatabase(dbvendor).Select(cols...).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)
	qb = qb.Where(whereValues)

	query, args, err := qb.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger"),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, true, jen.Nil())...)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.ID("row").Assign().ID("q").Dot("getOneRow").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Lit(uvn),
			jen.IDf("get%sQuery", sn),
			jen.ID("args"),
		),
		jen.Newline(),
		jen.List(jen.ID(uvn), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", sn).Call(
			jen.ID("ctx"),
			jen.ID("row"),
			jen.False(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("scanning %s", uvn),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID(uvn), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Const().IDf("get%sQuery", sn).Equals().Lit(query),
		jen.Newline(),
		jen.Commentf("Get%s fetches %s from the database.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Get%s", sn).Params(typ.BuildDBClientRetrievalMethodParams(proj)...).Params(jen.Op("*").ID("types").Dot(sn), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildJoinsListForListQuery(p *models.Project, typ models.DataType) []jen.Code {
	out := []jen.Code{}

	if typ.BelongsToStruct != nil {
		out = append(out, jen.IDf("%sOn%sJoinClause", typ.BelongsToStruct.PluralUnexportedVarName(), typ.Name.Plural()))
	}

	owners := p.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			out = append(out, jen.IDf("%sOn%sJoinClause", pt.BelongsToStruct.PluralUnexportedVarName(), pt.Name.Plural()))

		}
	}

	return out
}

func buildGetSomethingsList(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	prn := typ.Name.PluralRouteName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger"),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildIDBoilerplate(proj, typ, false, jen.Nil())...)

	bodyLines = append(bodyLines,
		jen.ID("x").Equals().Op("&").ID("types").Dotf("%sList", sn).Values(),
		jen.ID("logger").Equals().ID("filter").Dot("AttachToLogger").Call(jen.ID("logger")),
		jen.Qual(proj.InternalTracingPackage(), "AttachQueryFilterToSpan").Call(
			jen.ID("span"),
			jen.ID("filter"),
		),
		jen.Newline(),
		jen.If(jen.ID("filter").DoesNotEqual().Nil()).Body(
			jen.List(jen.ID("x").Dot("Page"), jen.ID("x").Dot("Limit")).Equals().List(jen.ID("filter").Dot("Page"), jen.ID("filter").Dot("Limit")),
		),
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dot("buildListQuery").Callln(
			jen.ID("ctx"),
			jen.Lit(prn),
			func() jen.Code {
				if len(proj.FindOwnerTypeChain(typ)) > 0 {
					return jen.IDf("get%sJoins", pn)
				}
				return jen.Nil()
			}(),
			jen.Nil(),
			jen.ID("accountOwnershipColumn"),
			jen.IDf("%sTableColumns", puvn),
			func() jen.Code {
				if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
					return jen.ID("accountID")
				}
				return jen.EmptyString()
			}(),
			jen.False(),
			jen.ID("filter"),
		),
		jen.Newline(),
		jen.List(jen.ID("rows"), jen.ID("err")).Assign().ID("q").Dot("performReadQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Lit(puvn),
			jen.ID("query"),
			jen.ID("args"),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("executing %s list retrieval query", pcn),
			))),
		jen.Newline(),
		jen.If(jen.List(jen.ID("x").Dot(pn), jen.ID("x").Dot("FilteredCount"), jen.ID("x").Dot("TotalCount"), jen.ID("err")).Equals().ID("q").Dotf("scan%s", pn).Call(
			jen.ID("ctx"),
			jen.ID("rows"),
			jen.True(),
		), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("scanning %s", pcn),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID("x"), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Get%s", pn).Params(typ.BuildDBClientListRetrievalMethodParams(proj)...).Params(jen.ID("x").Op("*").ID("types").Dotf("%sList", sn), jen.ID("err").ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildGetTotalSomethingCount(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)

	query, _, err := sqlBuilder.Select(fmt.Sprintf("COUNT(%s.id)", tableName)).
		From(tableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", tableName): nil,
		}).ToSql()
	if err != nil {
		panic(err)
	}

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger"),
		jen.Newline(),
		jen.List(jen.ID("count"), jen.ID("err")).Assign().ID("q").Dot("performCountQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.IDf("getTotal%sCountQuery", pn),
			jen.Litf("fetching count of %s", pcn),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Zero(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("querying for count of %s", pcn),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID("count"), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Const().IDf("getTotal%sCountQuery", pn).Equals().Lit(query),
		jen.Newline(),
		jen.Commentf("GetTotal%sCount fetches the count of %s from the database that meet a particular filter.", sn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("GetTotal%sCount", sn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.Uint64(), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildGetSomethingWithIDsQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()

	tableName := typ.Name.PluralRouteName()

	bodyLines := []jen.Code{
		jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("withIDsWhere").Assign().ID("squirrel").Dot("Eq").Valuesln(
			jen.Litf("%s.id", tableName).Op(":").ID("ids"),
			jen.Litf("%s.archived_on", tableName).Op(":").Nil(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.Litf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName()).Op(":").IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
		utils.ConditionalCode(typ.BelongsToAccount,
			jen.If(jen.ID("accountID").DoesNotEqual().EmptyString()).Body(
				jen.ID("withIDsWhere").Index(jen.Litf("%s.belongs_to_account", tableName)).Equals().ID("accountID"),
			),
		),
		jen.Newline(),
		jen.ID("subqueryBuilder").Assign().ID("q").Dot("sqlBuilder").Dot("Select").Call(jen.IDf("%sTableColumns", puvn).Op("...")).
			Dotln("From").Call(jen.Lit(puvn)).
			Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("unnest('{%s}'::text[])"), jen.ID("joinIDs").Call(jen.ID("ids")))).
			Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d"), jen.ID("limit"))),
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Assign().ID("q").Dot("sqlBuilder").Dot("Select").Call(jen.IDf("%sTableColumns", puvn).Op("...")).
			Dotln("FromSelect").Call(jen.ID("subqueryBuilder"), jen.Lit(puvn)).
			Dotln("Where").Call(jen.ID("withIDsWhere")).Dot("ToSql").Call(),
		jen.Newline(),
		jen.ID("q").Dot("logQueryBuildingError").Call(
			jen.ID("span"),
			jen.ID("err"),
		),
		jen.Newline(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	}

	lines := []jen.Code{
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("buildGet%sWithIDsQuery", pn).Params(typ.BuildGetListOfSomethingFromIDsParams(proj)...).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildGetSomethingWithIDs(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger"),
		jen.Newline(),
	}

	if typ.BelongsToStruct != nil {
		bodyLines = append(bodyLines,
			jen.If(jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).IsEqualTo().EmptyString()).Body(
				jen.Return(jen.Nil(), jen.ID("ErrInvalidIDProvided")),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.BelongsToStruct.Singular())), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.BelongsToStruct.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())),
			jen.Newline(),
		)
	}

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.If(jen.ID("ids").IsEqualTo().Nil()).Body(
			jen.Return(jen.Nil(), jen.ID("ErrNilInputProvided")),
		),
		jen.Newline(),
		jen.If(jen.ID("limit").IsEqualTo().Zero()).Body(
			jen.ID("limit").Equals().ID("uint8").Call(jen.Qual(proj.TypesPackage(), "DefaultLimit")),
		),
		jen.Newline(),
		jen.ID("logger").Equals().ID("logger").Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
			jen.Lit("limit").MapAssign().ID("limit"), jen.Lit("id_count").MapAssign().ID("len").Call(jen.ID("ids")))),
		jen.Newline(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID("q").Dotf("buildGet%sWithIDsQuery", pn).Call(
			constants.CtxVar(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("accountID")),
			jen.ID("limit"),
			jen.ID("ids"),
		),
		jen.Newline(),
		jen.List(jen.ID("rows"), jen.ID("err")).Assign().ID("q").Dot("performReadQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Litf("%s with IDs", pcn),
			jen.ID("query"),
			jen.ID("args"),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("fetching %s from database", pcn),
			))),
		jen.Newline(),
		jen.List(jen.ID(puvn), jen.ID("_"), jen.ID("_"), jen.ID("err")).Assign().ID("q").Dotf("scan%s", pn).Call(
			jen.ID("ctx"),
			jen.ID("rows"),
			jen.False(),
		),
		jen.If(jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("scanning %s", pcn),
			))),
		jen.Newline(),
		jen.Return().List(jen.ID(puvn), jen.Nil()),
	)

	lines := []jen.Code{
		jen.Commentf("Get%sWithIDs fetches %s from the database within a given set of IDs.", pn, pcn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Get%sWithIDs", pn).Params(typ.BuildGetListOfSomethingFromIDsParams(proj)...).Params(jen.Index().Op("*").ID("types").Dot(sn), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildCreateSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	uvn := typ.Name.UnexportedVarName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)

	creationColumns := []string{
		"id",
	}
	args := []interface{}{models.NewCodeWrapper(jen.ID("input").Dot("ID"))}
	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, field.Name.RouteName())
		args = append(args, models.NewCodeWrapper(jen.ID("input").Dot(field.Name.Singular())))
	}

	if typ.BelongsToStruct != nil {
		creationColumns = append(creationColumns, fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
		args = append(args, models.NewCodeWrapper(jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular())))
	}

	if typ.BelongsToAccount {
		creationColumns = append(creationColumns, "belongs_to_account")
		args = append(args, models.NewCodeWrapper(jen.ID("input").Dot("BelongsToAccount")))
	}

	query, args, err := sqlBuilder.Insert(tableName).
		Columns(creationColumns...).
		Values(args...).
		ToSql()

	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	fieldValues := []jen.Code{jen.ID("ID").MapAssign().ID("input").Dot("ID")}
	for _, field := range typ.Fields {
		fieldValues = append(fieldValues, jen.ID(field.Name.Singular()).MapAssign().ID("input").Dot(field.Name.Singular()))
	}

	if typ.BelongsToStruct != nil {
		fieldValues = append(fieldValues, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToAccount {
		fieldValues = append(fieldValues, jen.ID("BelongsToAccount").MapAssign().ID("input").Dot("BelongsToAccount"))
	}

	fieldValues = append(fieldValues, jen.ID("CreatedOn").MapAssign().ID("q").Dot("currentTime").Call())

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.If(jen.ID("input").IsEqualTo().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.ID("ErrNilInputProvided")),
		),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
			jen.ID("keys").Dotf("%sIDKey", sn),
			jen.ID("input").Dot("ID"),
		),
		jen.Newline(),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.Commentf("create the %s.", scn),
		jen.If(jen.ID("err").Assign().ID("q").Dot("performWriteQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Litf("%s creation", scn),
			jen.IDf("%sCreationQuery", uvn),
			jen.ID("args"),
		), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().List(jen.Nil(), jen.Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("creating %s", scn),
			))),
		jen.Newline(),
		jen.ID("x").Assign().Op("&").ID("types").Dot(sn).Valuesln(fieldValues...),
		jen.Newline(),
		jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
			jen.ID("span"),
			jen.ID("x").Dot("ID"),
		),
		jen.ID("logger").Dot("Info").Call(jen.Litf("%s created", scn)),
		jen.Newline(),
		jen.Return().List(jen.ID("x"), jen.Nil()),
	}

	lines := []jen.Code{
		jen.Const().IDf("%sCreationQuery", uvn).Equals().Lit(query),
		jen.Newline(),
		jen.Commentf("Create%s creates %s in the database.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Create%s", sn).Params(typ.BuildDBClientCreationMethodParams(proj)...).Params(jen.Op("*").ID("types").Dot(sn), jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildUpdateSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)

	updateWhere := squirrel.Eq{
		"id":          models.NewCodeWrapper(jen.ID("updated").Dot("ID")),
		"archived_on": nil,
	}

	if typ.BelongsToStruct != nil {
		updateWhere[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.ID("updated").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		updateWhere["belongs_to_account"] = models.NewCodeWrapper(jen.ID("updated").Dot("BelongsToAccount"))
	}

	updateBuilder := sqlBuilder.Update(tableName)

	for _, field := range typ.Fields {
		updateBuilder = updateBuilder.Set(field.Name.RouteName(), models.NewCodeWrapper(jen.ID("updated").Dot(field.Name.Singular())))
	}

	updateBuilder = updateBuilder.Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).Where(updateWhere)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	bodyLines := []jen.Code{
		jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		jen.If(jen.ID("updated").IsEqualTo().Nil()).Body(
			jen.Return().ID("ErrNilInputProvided"),
		),
		jen.Newline(),
		jen.ID("logger").Assign().ID("q").Dot("logger").Dot("WithValue").Call(
			jen.ID("keys").Dotf("%sIDKey", sn),
			jen.ID("updated").Dot("ID"),
		),
		jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
			jen.ID("span"),
			jen.ID("updated").Dot("ID"),
		),
		utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(
			jen.ID("span"),
			jen.ID("updated").Dot("BelongsToAccount"),
		)),
		jen.Newline(),
		jen.ID("args").Assign().Index().Interface().Valuesln(dbCallArgs...),
		jen.Newline(),
		jen.If(jen.ID("err").Assign().ID("q").Dot("performWriteQuery").Call(
			jen.ID("ctx"),
			jen.ID("q").Dot("db"),
			jen.Litf("%s update", scn),
			jen.IDf("update%sQuery", sn),
			jen.ID("args"),
		), jen.ID("err").DoesNotEqual().Nil()).Body(
			jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
				jen.ID("err"),
				jen.ID("logger"),
				jen.ID("span"),
				jen.Litf("updating %s", scn),
			)),
		jen.Newline(),
		jen.ID("logger").Dot("Info").Call(jen.Litf("%s updated", scn)),
		jen.Newline(),
		jen.Return().Nil(),
	}

	lines := []jen.Code{
		jen.Const().IDf("update%sQuery", sn).Equals().Lit(query),
		jen.Newline(),
		jen.Commentf("Update%s updates a particular %s.", sn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").Op("*").ID("SQLQuerier")).IDf("Update%s", sn).Params(typ.BuildDBClientUpdateMethodParams(proj, "updated")...).Params(jen.ID("error")).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildArchiveSomething(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	tableName := typ.Name.PluralRouteName()
	sqlBuilder := queryBuilderForDatabase(dbvendor)
	archiveWhere := squirrel.Eq{
		"id":          models.NewCodeWrapper(jen.IDf("%sID", uvn)),
		"archived_on": nil,
	}

	if typ.BelongsToStruct != nil {
		archiveWhere[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	if typ.BelongsToAccount {
		archiveWhere["belongs_to_account"] = models.NewCodeWrapper(jen.ID("accountID"))
	}

	query, args, err := sqlBuilder.Update(tableName).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(dbvendor))).
		Where(archiveWhere).ToSql()

	if err != nil {
		panic(err)
	}

	dbCallArgs := convertArgsToCode(args)

	return []jen.Code{
		jen.Const().IDf("archive%sQuery", sn).Equals().Lit(query),
		jen.Newline(),
		jen.Commentf("Archive%s archives %s from the database by its ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("q").PointerTo().ID("SQLQuerier")).IDf("Archive%s", sn).Params(typ.BuildDBClientArchiveMethodParams()...).Params(jen.ID("error")).Body(
			jen.List(jen.ID("ctx"), jen.ID("span")).Assign().ID("q").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID(constants.LoggerVarName).Assign().ID("q").Dot(constants.LoggerVarName),
			jen.Newline(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.If(jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).IsEqualTo().EmptyString()).Body(jen.Return().ID("ErrInvalidIDProvided"))
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qualf(proj.ConstantKeysPackage(), "%sIDKey", typ.BelongsToStruct.Singular()), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.Qualf(proj.InternalTracingPackage(), "Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
				}
				return jen.Null()
			}(),
			jen.Newline(),
			jen.If(jen.IDf("%sID", uvn).IsEqualTo().EmptyString()).Body(
				jen.Return().ID("ErrInvalidIDProvided"),
			),
			jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), fmt.Sprintf("%sIDKey", typ.Name.Singular())), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Qual(proj.InternalTracingPackage(), fmt.Sprintf("Attach%sIDToSpan", typ.Name.Singular())).Call(jen.ID(constants.SpanVarName), jen.IDf("%sID", typ.Name.UnexportedVarName())),
			jen.Newline(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.If(jen.ID("accountID").IsEqualTo().EmptyString()).Body(jen.Return().ID("ErrInvalidIDProvided"))),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID(constants.LoggerVarName).Equals().ID(constants.LoggerVarName).Dot("WithValue").Call(jen.Qual(proj.ConstantKeysPackage(), "AccountIDKey"), jen.ID("accountID"))),
			utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID(constants.SpanVarName), jen.ID("accountID"))),
			jen.Newline(),
			jen.ID("args").Assign().Index().Interface().Valuesln(
				dbCallArgs...,
			//func() jen.Code {
			//	if typ.BelongsToStruct != nil {
			//		return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
			//	}
			//	return jen.Null()
			//}(),
			//utils.ConditionalCode(typ.BelongsToAccount, jen.ID("accountID")),
			//jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.If(jen.ID("err").Assign().ID("q").Dot("performWriteQuery").Call(
				jen.ID("ctx"),
				jen.ID("q").Dot("db"),
				jen.Litf("%s archive", scn),
				jen.IDf("archive%sQuery", sn),
				jen.ID("args"),
			), jen.ID("err").DoesNotEqual().Nil()).Body(
				jen.Return().Qual(proj.ObservabilityPackage(), "PrepareError").Call(
					jen.ID("err"),
					jen.ID("logger"),
					jen.ID("span"),
					jen.Litf("updating %s", scn),
				)),
			jen.Newline(),
			constants.LoggerVar().Dot("Info").Call(jen.Litf("%s archived", scn)),
			jen.Newline(),
			jen.Return().Nil(),
		),
		jen.Newline(),
	}
}
