package querybuilders

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func iterablesDotGo(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) *jen.File {
	code := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(proj, code, false)

	sn := typ.Name.Singular()

	code.Add(
		jen.Var().Underscore().Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sSQLQueryBuilder", sn)).Equals().Parens(jen.PointerTo().ID(dbvendor.Singular())).Call(jen.Nil()),
		jen.Newline(),
	)

	code.Add(buildBuildSomethingExistsQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildGetSomethingQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildGetAllSomethingCountQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildGetBatchOfSomethingQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildGetListOfSomethingQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildGetSomethingWithIDsQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildCreateSomethingQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildUpdateSomethingQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildArchiveSomethingQuery(proj, typ, dbvendor)...)
	code.Add(buildBuildGetAuditLogEntriesForSomethingQuery(proj, typ, dbvendor)...)

	return code
}

func buildDBQueryMethodConditionalClauses(p *models.Project, typ models.DataType, includeSelf bool) []jen.Code {
	n := typ.Name
	pn := n.Plural()
	uvn := n.UnexportedVarName()

	whereValues := []jen.Code{}
	for _, pt := range p.FindOwnerTypeChain(typ) {
		whereValues = append(
			whereValues,
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pt.Name.Plural())),
				jen.Qual(p.QuerybuildingPackage(), "IDColumn"),
			).MapAssign().IDf("%sID", pt.Name.UnexportedVarName()),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pt.Name.Plural())),
				jen.Qual(p.QuerybuildingPackage(), "ArchivedOnColumn"),
			).MapAssign().Nil(),
		)

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pt.Name.Plural())),
					jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pt.Name.Plural())),
				).MapAssign().ID("accountID"),
			)
		}

		if pt.BelongsToStruct != nil {
			whereValues = append(
				whereValues,
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pt.Name.Plural())),
					jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pt.Name.Plural(), pt.BelongsToStruct.Singular())),
				).MapAssign().IDf("%sID", pt.BelongsToStruct.UnexportedVarName()),
			)
		}
	}

	if includeSelf {
		whereValues = append(whereValues,
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)), jen.Qual(p.QuerybuildingPackage(), "IDColumn")).MapAssign().IDf("%sID", uvn),
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)), jen.Qual(p.QuerybuildingPackage(), "ArchivedOnColumn")).MapAssign().Nil(),
		)

		if typ.BelongsToStruct != nil {
			whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pn, typ.BelongsToStruct.Singular()))).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		}
		if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
			whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn))).MapAssign().ID("accountID"))
		}
	} else {
		whereValues = append(whereValues,
			jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)), jen.Qual(p.QuerybuildingPackage(), "ArchivedOnColumn")).MapAssign().Nil(),
		)
	}

	return whereValues
}

func modifyQueryBuildingStatementWithJoinClauses(p *models.Project, typ models.DataType, qbStmt *jen.Statement) *jen.Statement {
	if typ.BelongsToStruct != nil {
		qbStmt = qbStmt.Dotln("Join").Call(
			jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sOn%sJoinClause", typ.BelongsToStruct.Plural(), typ.Name.Plural())),
		)
	}

	owners := p.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			qbStmt = qbStmt.Dotln("Join").Call(
				jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sOn%sJoinClause", pt.BelongsToStruct.Plural(), pt.Name.Plural())),
			)
		}
	}

	return qbStmt
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
		params = append(params, jen.List(lp...).Uint64())
	}

	return params
}

func buildBuildSomethingExistsQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	params := buildDBQuerierExistenceQueryMethodParams(proj, typ)
	whereValues := buildDBQueryMethodConditionalClauses(proj, typ, true)

	queryBuilderDecl := jen.ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual("fmt", "Sprintf").Call(
		jen.Lit("%s.%s"),
		jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
		jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
	)).
		Dotln("Prefix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistencePrefix")).
		Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)))

	queryBuilderDecl = modifyQueryBuildingStatementWithJoinClauses(proj, typ, queryBuilderDecl)

	queryBuilderDecl = queryBuilderDecl.
		Dotln("Suffix").Call(jen.Qual(proj.QuerybuildingPackage(), "ExistenceSuffix")).
		Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(whereValues...))

	bodyLines := []jen.Code{
		jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		bodyLines = append(bodyLines,
			jen.ID("tracing").Dotf("Attach%sIDToSpan", owner.Name.Singular()).Call(
				jen.ID("span"),
				jen.IDf("%sID", owner.Name.UnexportedVarName()),
			),
		)
	}
	bodyLines = append(bodyLines,
		jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
			jen.ID("span"),
			jen.IDf("%sID", uvn),
		),
	)

	bodyLines = append(bodyLines,
		utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
		jen.Newline(),
		jen.Return(jen.ID("b").Dot("buildQuery").Callln(
			jen.ID("span"),
			queryBuilderDecl,
		)),
	)

	lines := []jen.Code{
		jen.Commentf("Build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a user with a given ID exists.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("Build%sExistsQuery", sn).Params(params...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			bodyLines...,
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
		params = append(params, jen.List(lp...).Uint64())
	}

	return params
}

func buildBuildGetSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	params := buildDBQuerierRetrievalQueryMethodParams(proj, typ)
	whereValues := buildDBQueryMethodConditionalClauses(proj, typ, true)

	queryBuilderDecl := jen.ID("b").Dot("sqlBuilder").
		Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
		Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)))

	queryBuilderDecl = modifyQueryBuildingStatementWithJoinClauses(proj, typ, queryBuilderDecl)

	queryBuilderDecl = queryBuilderDecl.Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(whereValues...))

	bodyLines := []jen.Code{
		jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
	}

	for _, owner := range proj.FindOwnerTypeChain(typ) {
		bodyLines = append(bodyLines,
			jen.ID("tracing").Dotf("Attach%sIDToSpan", owner.Name.Singular()).Call(
				jen.ID("span"),
				jen.IDf("%sID", owner.Name.UnexportedVarName()),
			),
		)
	}
	bodyLines = append(bodyLines,
		jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
			jen.ID("span"),
			jen.IDf("%sID", uvn),
		),
	)

	bodyLines = append(bodyLines,
		utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
		jen.Newline(),
		jen.Return().ID("b").Dot("buildQuery").Callln(
			jen.ID("span"),
			queryBuilderDecl,
		),
	)

	lines := []jen.Code{
		jen.Commentf("BuildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to a user with a given ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildGet%sQuery", sn).Params(params...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAllSomethingCountQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	lines := []jen.Code{
		jen.Commentf("BuildGetAll%sCountQuery returns a query that fetches the total number of %s in the database.", pn, pcn),
		jen.Newline(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildGetAll%sCountQuery", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.String()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
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
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
					).MapAssign().Nil(),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetBatchOfSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()

	lines := []jen.Code{
		jen.Commentf("BuildGetBatchOf%sQuery returns a query that fetches every %s in the database within a bucketed range.", pn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildGetBatchOf%sQuery", pn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.ID("beginID"), jen.ID("endID")).Uint64()).Params(jen.ID("query").String(),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Gt").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
					).MapAssign().ID("beginID"))).
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Lt").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(
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
		params = append(params, jen.List(lp...).Uint64())
	}

	params = append(params, jen.ID("includeArchived").Bool(), jen.ID("filter").PointerTo().Qual(p.TypesPackage(), "QueryFilter"))

	return params
}

func buildJoinsListForListQuery(p *models.Project, typ models.DataType) []jen.Code {
	out := []jen.Code{}

	if typ.BelongsToStruct != nil {
		out = append(out,
			jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sOn%sJoinClause", typ.BelongsToStruct.Plural(), typ.Name.Plural())),
		)
	}

	owners := p.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			out = append(out,
				jen.Qual(p.QuerybuildingPackage(), fmt.Sprintf("%sOn%sJoinClause", pt.BelongsToStruct.Plural(), pt.Name.Plural())),
			)
		}
	}

	return out
}

func buildBuildGetListOfSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	params := buildDBQuerierListRetrievalQueryBuildingMethodParams(proj, typ)
	whereValues := buildDBQueryMethodConditionalClauses(proj, typ, false)

	lines := []jen.Code{
		jen.Commentf("BuildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given account,", pn, pcn),
		jen.Newline(),
		jen.Comment("and returns both the query and the relevant args to pass to the query executor."),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildGet%sQuery", pn).Params(params...).Params(jen.ID("query").String(),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.If(jen.ID("filter").DoesNotEqual().Nil()).Body(
				jen.Qual(proj.InternalTracingPackage(), "AttachFilterToSpan").Call(jen.ID("span"), jen.ID("filter").Dot("Page"),
					jen.ID("filter").Dot("Limit"),
					jen.String().Call(jen.ID("filter").Dot("SortBy")),
				),
			),
			jen.Newline(),
			utils.ConditionalCode(len(proj.FindOwnerTypeChain(typ)) > 0, jen.ID("joins").Assign().Index().String().Valuesln(
				buildJoinsListForListQuery(proj, typ)...,
			)),
			utils.ConditionalCode(len(whereValues) > 0, jen.ID("where").Assign().Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(
				whereValues...,
			)),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildListQuery").Callln(
				jen.ID("ctx"),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
				func() jen.Code {
					if len(proj.FindOwnerTypeChain(typ)) > 0 {
						return jen.ID("joins")
					}
					return jen.Nil()
				}(),
				func() jen.Code {
					if len(whereValues) > 0 {
						return jen.ID("where")
					}
					return jen.Nil()
				}(),
				func() jen.Code {
					if typ.RestrictedToAccountMembers {
						return jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn))
					} else if typ.BelongsToStruct != nil {
						return jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pn, typ.BelongsToStruct.Singular()))
					}
					return jen.EmptyString()
				}(),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)),
				func() jen.Code {
					if typ.RestrictedToAccountMembers {
						return jen.ID("accountID")
					} else if typ.BelongsToStruct != nil {
						return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
					}
					return jen.Zero()
				}(),
				jen.ID("includeArchived"),
				jen.ID("filter"),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetSomethingWithIDsQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	prerequisiteIDs := []jen.Code{}

	if typ.BelongsToStruct != nil {
		prerequisiteIDs = append(prerequisiteIDs, jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	}
	if typ.BelongsToAccount {
		prerequisiteIDs = append(prerequisiteIDs, jen.ID("accountID"))
	}

	qbDecl := jen.ID("b").Dot("sqlBuilder").
		Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
		Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
		Dotln("Where").Call(jen.ID("where")).
		Dotln("OrderBy").Call(jen.Qual("fmt", "Sprintf").Call(
		jen.Lit("CASE %s.%s %s"),
		jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
		jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
		jen.ID("whenThenStatement"),
	)).Dotln("Limit").Call(jen.Uint64().Call(jen.ID("limit")))

	if dbvendor.SingularPackageName() == "postgres" {
		qbDecl = jen.ID("b").Dot("sqlBuilder").
			Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Op("...")).
			Dotln("FromSelect").Call(jen.ID("subqueryBuilder"), jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
			Dotln("Where").Call(jen.ID("where"))
	}

	bodyLines := []jen.Code{
		jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
		jen.Defer().ID("span").Dot("End").Call(),
		jen.Newline(),
		utils.ConditionalCode(typ.RestrictedToAccountAtSomeLevel(proj), jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("accountID"))),
		jen.Newline(),
		utils.ConditionalCode(dbvendor.SingularPackageName() != "postgres", jen.ID("whenThenStatement").Assign().ID("joinIDs").Call(jen.ID("ids"))),
		jen.ID("where").Assign().Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
				jen.Qual(proj.QuerybuildingPackage(), "IDColumn"),
			).MapAssign().ID("ids"),
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("%s.%s"),
				jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
				jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn"),
			).MapAssign().Nil(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.Qual("fmt", "Sprintf").Call(
						jen.Lit("%s.%s"),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
						jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pn, typ.BelongsToStruct.Singular())),
					).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
				}
				return jen.Null()
			}(),
		),
		jen.Newline(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.If(jen.ID("restrictToAccount")).Body(
			jen.ID("where").Index(
				jen.Qual("fmt", "Sprintf").Call(
					jen.Lit("%s.%s"),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)),
					jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)),
				),
			).Equals().ID("accountID"),
		)),
		jen.Newline(),
		utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.ID("subqueryBuilder").Assign().ID("b").Dot("sqlBuilder").Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableColumns", pn)).Spread()).
			Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).
			Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("unnest('{%s}'::int[])"), jen.ID("joinIDs").Call(jen.ID("ids")))).
			Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d"), jen.ID("limit"))),
		),
		jen.Newline(),
		jen.Return().ID("b").Dot("buildQuery").Callln(
			jen.ID("span"),
			qbDecl,
		),
	}

	lines := []jen.Code{
		jen.Commentf("BuildGet%sWithIDsQuery builds a SQL query selecting %s that belong to a given account,", pn, pcn).Newline(),
		jen.Comment("and have IDs that exist within a given set of IDs. Returns both the query and the relevant").Newline(),
		jen.Comment("args to pass to the query executor. This function is primarily intended for use with a search").Newline(),
		jen.Comment("index, which would provide a slice of string IDs to query against. This function accepts a").Newline(),
		jen.Comment("slice of uint64s instead of a slice of strings in order to ensure all the provided strings").Newline(),
		jen.Comment("are valid database IDs, because there's no way in squirrel to escape them in the unnest join,").Newline(),
		jen.Comment("and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.").Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildGet%sWithIDsQuery", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			func() jen.Code {
				if len(prerequisiteIDs) > 0 {
					return jen.List(prerequisiteIDs...).Uint64()
				}
				return jen.Null()
			}(),
			jen.ID("limit").ID("uint8"),
			jen.ID("ids").Index().Uint64(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("restrictToAccount").ID("bool")),
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Body(
			bodyLines...,
		),
		jen.Newline(),
	}

	return lines
}

func determineCreationColumns(proj *models.Project, typ models.DataType) []jen.Code {
	pn := typ.Name.Plural()

	creationColumns := []jen.Code{
		jen.Qual(proj.QuerybuildingPackage(), "ExternalIDColumn"),
	}

	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTable%sColumn", pn, field.Name.Singular())))
	}

	if typ.BelongsToStruct != nil {
		creationColumns = append(creationColumns, jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pn, typ.BelongsToStruct.Singular())))
	}
	if typ.BelongsToAccount {
		creationColumns = append(creationColumns, jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)))
	}

	return creationColumns
}

func determineCreationQueryValues(inputVarName string, typ models.DataType) []jen.Code {
	valuesColumns := []jen.Code{
		jen.ID("b").Dot("externalIDGenerator").Dot("NewExternalID").Call(),
	}

	for _, field := range typ.Fields {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dot(field.Name.Singular()))
	}

	if typ.BelongsToStruct != nil {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dot(constants.AccountOwnershipFieldName))
	}

	return valuesColumns
}

func buildBuildCreateSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	columns := determineCreationColumns(proj, typ)
	values := determineCreationQueryValues("input", typ)

	lines := []jen.Code{
		jen.Commentf("BuildCreate%sQuery takes %s and returns a creation query for that %s and the relevant arguments.", sn, scnwp, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildCreate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").PointerTo().ID("types").Dotf("%sCreationInput", sn)).Params(jen.ID("query").String(),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").Dot("Insert").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn))).Dotln("Columns").Callln(columns...).Dotln("Values").Callln(values...).Add(utils.ConditionalCode(dbvendor.SingularPackageName() == "postgres", jen.Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("RETURNING %s"), jen.Qual(proj.QuerybuildingPackage(), "IDColumn"))))),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildUpdateSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	queryBuilderDecl := jen.ID("b").Dot("sqlBuilder").Dot("Update").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableName", pn)))

	for _, field := range typ.Fields {
		queryBuilderDecl = queryBuilderDecl.Dotln("Set").Call(jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTable%sColumn", pn, field.Name.Singular())), jen.ID("input").Dot(field.Name.Singular()))
	}

	queryBuilderDecl = queryBuilderDecl.
		Dotln("Set").Call(
		jen.Qual(proj.QuerybuildingPackage(), "LastUpdatedOnColumn"),
		jen.ID("currentUnixTimeQuery"),
	).
		Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(
		jen.Qual(proj.QuerybuildingPackage(), "IDColumn").MapAssign().ID("input").Dot("ID"),
		jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").MapAssign().Nil(),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pn, typ.BelongsToStruct.Singular())).MapAssign().ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular())
			}
			return jen.Null()
		}(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)).MapAssign().ID("input").Dot("BelongsToAccount")),
	))

	lines := []jen.Code{
		jen.Commentf("BuildUpdate%sQuery takes %s and returns an update SQL query, with the relevant query parameters.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildUpdate%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").PointerTo().ID("types").Dot(sn)).Params(jen.ID("query").String(),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID("tracing").Dotf("Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(
						jen.ID("span"),
						jen.ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()),
					)
				}
				return jen.Null()
			}(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.ID("input").Dot("ID"),
			),
			utils.ConditionalCode(typ.BelongsToAccount, jen.Qual(proj.InternalTracingPackage(), "AttachAccountIDToSpan").Call(jen.ID("span"), jen.ID("input").Dot("BelongsToAccount"))),
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				queryBuilderDecl,
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

	params = append(params, jen.List(lp...).Uint64())

	return params
}

func buildBuildArchiveSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	uvn := typ.Name.UnexportedVarName()
	scn := typ.Name.SingularCommonName()

	params := buildDBQuerierArchiveQueryMethodParams(typ)

	lines := []jen.Code{
		jen.Commentf("BuildArchive%sQuery returns a SQL query which marks a given %s belonging to a given account as archived.", sn, scn),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildArchive%sQuery", sn).Params(params...).Params(jen.ID("query").String(),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.ID("tracing").Dotf("Attach%sIDToSpan", typ.BelongsToStruct.Singular()).Call(
						jen.ID("span"),
						jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()),
					)
				}
				return jen.Null()
			}(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(jen.ID("span"), jen.IDf("%sID", uvn)),
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
					Dotln("Where").Call(jen.Qual(constants.SQLGenerationLibrary, "Eq").Valuesln(
					jen.Qual(proj.QuerybuildingPackage(), "IDColumn").MapAssign().IDf("%sID", uvn),
					jen.Qual(proj.QuerybuildingPackage(), "ArchivedOnColumn").MapAssign().Nil(),
					func() jen.Code {
						if typ.BelongsToStruct != nil {
							return jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableBelongsTo%sColumn", pn, typ.BelongsToStruct.Singular())).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
						}
						return jen.Null()
					}(),
					utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.Qual(proj.QuerybuildingPackage(), fmt.Sprintf("%sTableAccountOwnershipColumn", pn)).MapAssign().ID("accountID")),
				)),
			),
		),
		jen.Newline(),
	}

	return lines
}

func buildBuildGetAuditLogEntriesForSomethingQuery(proj *models.Project, typ models.DataType, dbvendor wordsmith.SuperPalabra) []jen.Code {
	sn := typ.Name.Singular()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()

	keyDecl := jen.IDf("%sIDKey", typ.Name.UnexportedVarName()).Assign().Qual("fmt", "Sprintf").Callln(
		jen.ID("jsonPluckQuery"),
		jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName"),
		jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableContextColumn"),
		utils.ConditionalCode(dbvendor.SingularPackageName() == "mysql", jen.IDf("%sID", uvn)),
		jen.Qual(proj.InternalAuditPackage(), fmt.Sprintf("%sAssignmentKey", sn)),
	)

	keyUsage := jen.Null()

	switch dbvendor.Singular() {
	case string(models.MySQL):
		keyUsage = jen.Qual(constants.SQLGenerationLibrary, "Expr").Call(jen.IDf("%sIDKey", uvn))
	case string(models.Postgres), string(models.Sqlite):
		keyUsage = jen.Qual(constants.SQLGenerationLibrary, "Eq").Values(jen.IDf("%sIDKey", uvn).MapAssign().IDf("%sID", uvn))
	}

	lines := []jen.Code{
		jen.Commentf("BuildGetAuditLogEntriesFor%sQuery constructs a SQL query for fetching audit log entries relating to %s with a given ID.", sn, scnwp),
		jen.Newline(),
		jen.Func().Params(jen.ID("b").PointerTo().ID(dbvendor.Singular())).IDf("BuildGetAuditLogEntriesFor%sQuery", sn).Params(jen.ID("ctx").Qual("context", "Context"),
			jen.IDf("%sID", uvn).Uint64()).Params(jen.ID("query").String(),
			jen.ID("args").Index().Interface()).Body(
			jen.List(jen.Underscore(), jen.ID("span")).Assign().ID("b").Dot("tracer").Dot("StartSpan").Call(jen.ID("ctx")),
			jen.Defer().ID("span").Dot("End").Call(),
			jen.Newline(),
			jen.ID("tracing").Dotf("Attach%sIDToSpan", sn).Call(
				jen.ID("span"),
				jen.IDf("%sID", uvn),
			),
			jen.Newline(),
			keyDecl,
			jen.Newline(),
			jen.Return().ID("b").Dot("buildQuery").Callln(
				jen.ID("span"),
				jen.ID("b").Dot("sqlBuilder").
					Dot("Select").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableColumns").Op("...")).
					Dotln("From").Call(jen.Qual(proj.QuerybuildingPackage(), "AuditLogEntriesTableName")).
					Dotln("Where").Call(keyUsage).
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
