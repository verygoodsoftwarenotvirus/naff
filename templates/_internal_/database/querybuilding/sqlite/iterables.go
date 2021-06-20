package sqlite

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()

	code.Add(
		jen.Var().ID("_").Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sSQLQueryBuilder", sn)).Op("=").Parens(jen.Op("*").ID("Sqlite")).Call(jen.ID("nil")),
		jen.Newline(),
	)

	code.Add(buildBuildSomethingExistsQuery(proj, typ)...)
	code.Add(buildBuildGetSomethingQuery(proj, typ)...)
	code.Add(buildBuildGetAllSomethingCountQuery(proj, typ)...)
	code.Add(buildBuildGetBatchOfSomethingQuery(proj, typ)...)
	code.Add(buildBuildGetListOfSomethingQuery(proj, typ)...)
	code.Add(buildBuildGetSomethingWithIDsQuery(proj, typ)...)
	code.Add(buildBuildCreateSomethingQuery(proj, typ)...)
	code.Add(buildBuildUpdateSomethingQuery(proj, typ)...)
	code.Add(buildBuildArchiveSomethingQuery(proj, typ)...)
	code.Add(buildBuildGetAuditLogEntriesForSomethingQuery(proj, typ)...)

	return code
}

func buildDBQuerierExistenceQueryMethodParams(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxParam()}

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
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	return params
}

func buildBuildSomethingExistsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	params := buildDBQuerierExistenceQueryMethodParams(proj, typ)

	lines := []jen.Code{
		jen.Commentf("Build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a user with a given ID exists.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("Build%sExistsQuery", sn).Params(params...).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				)).
					Dotln("Prefix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistencePrefix")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Suffix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistenceSuffix")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().IDf("%sID", uvn),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)),
					).MapAssign().ID("accountID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					).MapAssign().ID("nil"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildDBQuerierRetrievalQueryMethodParams(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxParam()}

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
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	return params
}

func buildBuildGetSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	params := buildDBQuerierRetrievalQueryMethodParams(proj, typ)

	lines := []jen.Code{
		jen.Commentf("BuildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to a user with a given ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildGet%sQuery", sn).Params(params...).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().IDf("%sID", uvn),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)),
					).MapAssign().ID("accountID"),
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					).MapAssign().ID("nil"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllSomethingCountQuery(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Commentf("BuildGetAll%sCountQuery returns a query that fetches the total number of %s in the database.", pn, pcn),
		jen.Newline(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildGetAll%sCountQuery", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("string")).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),

			jen.Return().ID("b").Dot("buildQueryOnly").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.ID("columnCountQueryTemplate"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
				)).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).MapAssign().ID("nil"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetBatchOfSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Commentf("BuildGetBatchOf%sQuery returns a query that fetches every %s in the database within a bucketed range.", pn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildGetBatchOf%sQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("beginID"), jen.ID("endID")).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Gt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().ID("beginID"))).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Lt").Valuesln(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
				).MapAssign().ID("endID"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildDBQuerierListRetrievalQueryBuildingMethodParams(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxParam()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("%sID", pt.Name.UnexportedVarName()))
	}
	if typ.RestrictedToAccountAtSomeLevel(p) {
		lp = append(lp, jen.ID("accountID"))
	}

	if len(lp) > 0 {
		params = append(params, jen.List(lp...).ID("uint64"))
	}

	params = append(params, jen.ID("forAdmin").Bool(), jen.ID("filter").Op("*").Qual(p.TypesPackage(), "QueryFilter"))

	return params
}

func buildBuildGetListOfSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	params := buildDBQuerierListRetrievalQueryBuildingMethodParams(proj, typ)

	lines := []jen.Code{
		jen.Commentf("BuildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given account,", pn, pcn),
		jen.Newline(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildGet%sQuery", pn).Params(params...).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Body(
				jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.ID("string").Call(jen.ID("filter").Dot("SortBy")),
				),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildListQuery").Callln(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
				jen.Nil(),
				jen.Nil(),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)),
				jen.ID("accountID"),
				jen.ID("forAdmin"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetSomethingWithIDsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Commentf("BuildGet%sWithIDsQuery builds a SQL query selecting %s that belong to a given account,", pn, pcn).Newline(),
		jen.Comment("and have IDs that exist within a given set of IDs. Returns both the query and the relevant").Newline(),
		jen.Comment("args to pass to the query executor. This function is primarily intended for use with a search").Newline(),
		jen.Comment("index, which would provide a slice of string IDs to query against. This function accepts a").Newline(),
		jen.Comment("slice of uint64s instead of a slice of strings in order to ensure all the provided strings").Newline(),
		jen.Comment("are valid database IDs, because there's no way in squirrel to escape them in the unnest join,").Newline(),
		jen.Comment("and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.").Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildGet%sWithIDsQuery", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.ID("accountID").ID("uint64")),
			jen.ID("limit").ID("uint8"),
			jen.ID("ids").Index().ID("uint64"),
			jen.ID("includeArchived").ID("bool"),
		).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
			jen.Newline(),
			jen.ID("whenThenStatement").Op(":=").ID("buildWhenThenStatement").Call(jen.ID("ids")),
			jen.ID("where").Op(":=").ID("squirrel").Dot("Eq").Valuesln(jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
				jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
			).MapAssign().ID("ids"),
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
				).MapAssign().ID("nil"),
			),
			jen.Newline(),
			jen.If(jen.Op("!").ID("includeArchived")).Body(
				jen.ID("where").Index(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)),
				)).Op("=").ID("accountID"),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Where").Call(jen.ID("where")).
					Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("CASE %s.%s %s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
					jen.ID("whenThenStatement"),
				)).
					Dotln("Limit").Call(jen.ID("uint64").Call(jen.ID("limit"))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildCreateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("BuildCreate%sQuery takes %s and returns a creation query for that %s and the relevant arguments.", sn, scnwp, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildCreate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").ID("types").Dotf("%sCreationInput", sn)).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),

			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Columns").Callln(
					jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableNameColumn", pn)),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableDetailsColumn", pn)),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)),
				).
					Dotln("Values").Callln(
					jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
					jen.ID("input").Dot("Name"),
					jen.ID("input").Dot("Details"),
					jen.ID("input").Dot("BelongsToAccount"),
				),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("BuildUpdate%sQuery takes %s and returns an update SQL query, with the relevant query parameters.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildUpdate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").ID("types").Dot(sn)).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("BelongsToAccount")),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableNameColumn", pn)),
					jen.ID("input").Dot("Name"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableDetailsColumn", pn)),
					jen.ID("input").Dot("Details"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").MapAssign().ID("input").Dot("ID"),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").MapAssign().ID("nil"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)).MapAssign().ID("input").Dot("BelongsToAccount"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildDBQuerierArchiveQueryMethodParams(typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxParam()}

	lp := []jen.Code{}
	if typ.BelongsToStruct != nil {
		lp = append(lp, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	lp = append(lp, jen.IDf("%sID", typ.Name.UnexportedVarName()))

	if typ.RestrictedToAccountMembers {
		lp = append(lp, jen.ID("accountID"))
	}

	params = append(params, jen.List(lp...).ID("uint64"))

	return params
}

func buildBuildArchiveSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	params := buildDBQuerierArchiveQueryMethodParams(typ)

	lines := []jen.Code{
		jen.Commentf("BuildArchive%sQuery returns a SQL query which marks a given %s belonging to a given account as archived.", sn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildArchive%sQuery", sn).Params(params...).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Set").Call(
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					jen.ID("currentUnixTimeQuery"),
				).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Valuesln(jen.Qual(proj.QuerybuildingPackage(), "IDColumn").MapAssign().IDf("%sID", uvn),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").MapAssign().ID("nil"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)).MapAssign().ID("accountID"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	lines := []jen.Code{
		jen.Commentf("BuildGetAuditLogEntriesFor%sQuery constructs a SQL query for fetching audit log entries relating to %s with a given ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").Op("*").ID("Sqlite")).IDf("BuildGetAuditLogEntriesFor%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).ID("uint64")).Params(jen.ID("query").ID("string"),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.ID("_"), jen.ID("span")).Op(":=").ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			jen.IDf("%sIDKey", typ.Name.UnexportedVarName()).Assign().Qual("fmt", "Sprintf").Call(
				jen.ID("jsonPluckQuery"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
				jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
				jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sAssignmentKey", sn)),
			),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
					Dotln("Where").Call(jen.ID("squirrel").Dot("Eq").Values(jen.IDf("%sIDKey", uvn).MapAssign().IDf("%sID", uvn))).
					Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
					jen.Qual(proj.QuerybuildingPackage(), "CreatedOnColumn"),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}
