package queriers

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

const (
	squirrelPkg = "github.com/Masterminds/squirrel"
)

func isPostgres(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.Singular() == "Postgres"
}

func isSqlite(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.Singular() == "Sqlite"
}

func isMariaDB(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.Singular() == "MariaDB" || dbvendor.RouteName() == "maria_db"
}

func buildIterableConstants(typ models.DataType) []jen.Code {
	n := typ.Name
	prn := n.PluralRouteName()
	puvn := n.PluralUnexportedVarName()

	consts := []jen.Code{
		jen.IDf("%sTableName", puvn).Equals().Lit(prn),
	}

	for _, field := range typ.Fields {
		consts = append(consts,
			jen.IDf("%sTable%sColumn", puvn, field.Name.Singular()).Equals().Lit(field.Name.RouteName()),
		)
	}

	if typ.BelongsToUser {
		consts = append(consts, jen.IDf("%sUserOwnershipColumn", puvn).Equals().Lit("belongs_to_user"))
	}
	if typ.BelongsToStruct != nil {
		consts = append(consts, jen.IDf("%sTableOwnershipColumn", puvn).Equals().Litf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
	}

	return consts
}

func buildIterableVariableDecs(proj *models.Project, typ models.DataType) []jen.Code {
	puvn := typ.Name.PluralUnexportedVarName()

	vars := []jen.Code{
		jen.IDf("%sTableColumns", puvn).Equals().Index().String().Valuesln(
			buildTableColumns(proj, typ)...,
		),
		jen.Line(),
	}

	for _, ct := range proj.FindDependentsOfType(typ) {
		vars = append(
			vars,
			jen.IDf("%sOn%sJoinClause", puvn, ct.Name.Plural()).Equals().Qual("fmt", "Sprintf").Call(
				jen.Lit("%s ON %s.%s=%s.%s"),
				jen.IDf("%sTableName", puvn),
				jen.IDf("%sTableName", ct.Name.PluralUnexportedVarName()),
				jen.IDf("%sTableOwnershipColumn", ct.Name.PluralUnexportedVarName()),
				jen.IDf("%sTableName", puvn),
				jen.ID("idColumn"),
			),
		)
	}

	return []jen.Code{
		jen.Var().Defs(vars...),
	}
}

func iterablesDotGo(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	spn := dbvendor.SingularPackageName()

	code := jen.NewFilePathName(proj.DatabaseV1Package("queriers", "v1", spn), spn)

	utils.AddImports(proj, code)

	code.Add(jen.Const().Defs(buildIterableConstants(typ)...), jen.Line())
	code.Add(buildIterableVariableDecs(proj, typ)...)
	code.Add(buildScanSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildScanListOfSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildSomethingExistsFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildSomethingAllCountQueryDecls(dbvendor, typ)...)
	code.Add(buildGetAllSomethingCountFuncDecl(dbvendor, typ)...)
	code.Add(buildGetBatchOfSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetAllOfSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildCreateSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildCreateSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildUpdateSomethingFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildArchiveSomethingQueryFuncDecl(proj, dbvendor, typ)...)
	code.Add(buildArchiveSomethingFuncDecl(proj, dbvendor, typ)...)

	return code
}

func buildTableColumns(_ *models.Project, typ models.DataType) []jen.Code {
	puvn := typ.Name.PluralUnexportedVarName()
	tableNameVar := fmt.Sprintf("%sTableName", puvn)

	tableColumns := []jen.Code{utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.ID("idColumn"))}

	for _, field := range typ.Fields {
		tableColumns = append(tableColumns, utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.IDf("%sTable%sColumn", puvn, field.Name.Singular())))
	}

	tableColumns = append(tableColumns,
		utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.ID("createdOnColumn")),
		utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.ID("lastUpdatedOnColumn")),
		utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.ID("archivedOnColumn")),
	)

	if typ.BelongsToUser {
		tableColumns = append(tableColumns, utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.IDf("%sUserOwnershipColumn", puvn)))
	}
	if typ.BelongsToStruct != nil {
		tableColumns = append(tableColumns, utils.FormatString("%s.%s", jen.ID(tableNameVar), jen.IDf("%sTableOwnershipColumn", puvn)))
	}

	return tableColumns
}

func buildScanFields(_ *models.Project, typ models.DataType) (scanFields []jen.Code) {
	scanFields = []jen.Code{jen.AddressOf().ID("x").Dot("ID")}

	for _, field := range typ.Fields {
		scanFields = append(scanFields, jen.AddressOf().ID("x").Dot(field.Name.Singular()))
	}

	scanFields = append(scanFields,
		jen.AddressOf().ID("x").Dot("CreatedOn"),
		jen.AddressOf().ID("x").Dot("LastUpdatedOn"),
		jen.AddressOf().ID("x").Dot("ArchivedOn"),
	)

	if typ.BelongsToUser {
		scanFields = append(scanFields, jen.AddressOf().ID("x").Dot(constants.UserOwnershipFieldName))
	}
	if typ.BelongsToStruct != nil {
		scanFields = append(scanFields, jen.AddressOf().ID("x").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	return scanFields
}

// what's the difference between these two things
func buildScanSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pscnwp := typ.Name.ProperSingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string(dbvendor.Abbreviation()[0]))
	dbvsn := dbvendor.Singular()

	return []jen.Code{
		jen.Commentf("scan%s takes a database Scanner (i.e. *sql.Row) and scans the result into %s struct", sn, pscnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("scan%s", sn).Params(
			jen.ID("scan").Qual(proj.DatabaseV1Package(), "Scanner"),
		).Params(
			jen.PointerTo().Qual(proj.ModelsV1Package(), sn),
			jen.Error(),
		).Block(
			func() []jen.Code {
				body := []jen.Code{
					jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), sn).Values(),
					jen.Line(),
					jen.ID("targetVars").Assign().Index().Interface().Valuesln(buildScanFields(proj, typ)...),
					jen.Line(),
					jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Call(jen.ID("targetVars").Spread()), jen.Err().DoesNotEqual().ID("nil")).Block(
						jen.Return().List(jen.Nil(), jen.Err()),
					),
					jen.Line(),
				}

				body = append(body,
					jen.Line(),
					jen.Return().List(jen.ID("x"), jen.Nil()),
				)

				return body
			}()...,
		),
		jen.Line(),
	}
}

func buildScanListOfSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	dbfl := strings.ToLower(string(dbvendor.Abbreviation()[0]))
	dbvsn := dbvendor.Singular()

	return []jen.Code{
		jen.Commentf("scan%s takes a logger and some database rows and turns them into a slice of %s.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("scan%s", pn).Params(
			jen.ID("rows").Qual(proj.DatabaseV1Package(), "ResultIterator"),
		).Params(
			jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.Error(),
		).Block(
			jen.Var().Defs(
				jen.ID("list").Index().Qual(proj.ModelsV1Package(), sn),
			),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("x"), jen.Err()).Assign().ID(dbfl).Dotf("scan%s", sn).Call(jen.ID("rows")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.Line(),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.PointerTo().ID("x")),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.ID("closeErr").Assign().ID("rows").Dot("Close").Call(), jen.ID("closeErr").DoesNotEqual().ID("nil")).Block(
				jen.ID(dbfl).Dot(constants.LoggerVarName).Dot("Error").Call(jen.ID("closeErr"), jen.Lit("closing database rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}
}

func buildSomethingExistsQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	var (
		comment string
	)

	dbfl := strings.ToLower(string(dbvendor.Abbreviation()[0]))
	dbvsn := dbvendor.Singular()
	n := typ.Name
	sn := n.Singular()
	scnwp := n.SingularCommonNameWithPrefix()
	puvn := n.PluralUnexportedVarName()

	params := typ.BuildDBQuerierExistenceQueryMethodParams(proj)
	whereValues := typ.BuildDBQuerierExistenceQueryMethodConditionalClauses(proj)

	if typ.BelongsToUser && typ.RestrictedToUser {
		comment = fmt.Sprintf("build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a user with a given ID exists", sn, scnwp)
	}
	if typ.BelongsToStruct != nil {
		comment = fmt.Sprintf("build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a %s with a given ID exists", sn, scnwp, typ.BelongsToStruct.SingularCommonNameWithPrefix())
	} else if typ.BelongsToNobody {
		comment = fmt.Sprintf("build%sExistsQuery constructs a SQL query for checking if %s with a given ID exists", sn, scnwp)
	}

	qbStmt := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Select").Call(utils.FormatString("%s.%s", jen.IDf("%sTableName", puvn), jen.ID("idColumn"))).
		Dotln("Prefix").Call(jen.ID("existencePrefix")).
		Dotln("From").Call(jen.IDf("%sTableName", puvn))

	qbStmt = typ.ModifyQueryBuildingStatementWithJoinClauses(proj, qbStmt)

	qbStmt = qbStmt.Dotln("Suffix").Call(jen.ID("existenceSuffix")).
		Dotln("Where").Call(jen.Qual(squirrelPkg, "Eq").Valuesln(whereValues...)).
		Dot("ToSql").Call()

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("build%sExistsQuery", sn).Params(params...).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			qbStmt,
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildSomethingExistsFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	cn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()

	const existenceVarName = "exists"

	params := typ.BuildDBQuerierExistenceMethodParams(proj)
	buildQueryParams := typ.BuildDBQuerierExistenceQueryBuildingArgs(proj)

	return []jen.Code{
		jen.Commentf("%sExists queries the database to see if a given %s belonging to a given user exists.", sn, cn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("%sExists", sn).Params(
			params...,
		).Params(
			jen.ID(existenceVarName).Bool(),
			jen.Err().Error(),
		).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("build%sExistsQuery", sn).Call(buildQueryParams...),
			jen.Line(),
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID(existenceVarName)),
			jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
				jen.Return(jen.False(), jen.Nil()),
			),
			jen.Line(),
			jen.Return().List(jen.ID(existenceVarName), jen.Err()),
		),
		jen.Line(),
	}
}

func buildGetSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	var (
		comment string
	)

	dbvsn := dbvendor.Singular()
	n := typ.Name
	sn := n.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	scnwp := n.SingularCommonNameWithPrefix()
	puvn := n.PluralUnexportedVarName()

	params := typ.BuildDBQuerierRetrievalMethodParams(proj)
	whereValues := typ.BuildDBQuerierRetrievalQueryMethodConditionalClauses(proj)

	if typ.BelongsToUser && typ.RestrictedToUser {
		comment = fmt.Sprintf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to a user with a given ID.", sn, scnwp)
	} else if typ.BelongsToStruct != nil {
		tsnwp := typ.BelongsToStruct.SingularCommonNameWithPrefix()
		comment = fmt.Sprintf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to %s with a given ID.", sn, scnwp, tsnwp)
	} else if typ.BelongsToNobody {
		comment = fmt.Sprintf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID.", sn, scnwp)
	}

	qbStmt := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Spread()).
		Dotln("From").Call(jen.IDf("%sTableName", puvn))

	qbStmt = typ.ModifyQueryBuildingStatementWithJoinClauses(proj, qbStmt)

	qbStmt = qbStmt.Dotln("Where").Call(jen.Qual(squirrelPkg, "Eq").Valuesln(whereValues...)).
		Dotln("ToSql").Call()

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildGet%sQuery", sn).Params(
			params...,
		).Params(
			jen.ID("query").String(),
			jen.ID("args").Index().Interface(),
		).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			qbStmt,
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildGetSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()

	params := typ.BuildDBQuerierRetrievalQueryMethodParams(proj)
	buildQueryParams := typ.BuildDBQuerierRetrievalQueryBuildingArgs(proj)

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("Get%s", sn).Params(params...).
			Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sQuery", sn).Call(buildQueryParams...),
			jen.Line(),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return(jen.ID(dbfl).Dotf("scan%s", sn).Call(jen.ID("row"))),
		),
		jen.Line(),
	}
}

func buildSomethingAllCountQueryDecls(dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Var().Defs(
			jen.IDf("all%sCountQueryBuilder", pn).Qual("sync", "Once"),
			jen.IDf("all%sCountQuery", pn).String(),
		),
		jen.Line(),
		jen.Commentf("buildGetAll%sCountQuery returns a query that fetches the total number of %s in the database.", pn, pcn),
		jen.Line(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildGetAll%sCountQuery", pn).Params().Params(jen.String()).Block(
			jen.IDf("all%sCountQueryBuilder", pn).Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().Err().Error(),
				jen.Line(),
				jen.List(jen.IDf("all%sCountQuery", pn), jen.Underscore(), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(utils.FormatStringWithArg(jen.ID("countQuery"), jen.IDf("%sTableName", puvn))).
					Dotln("From").Call(jen.IDf("%sTableName", puvn)).
					Dotln("Where").Call(
					jen.Qual(squirrelPkg, "Eq").Valuesln(
						utils.FormatString("%s.%s", jen.IDf("%sTableName", puvn), jen.ID("archivedOnColumn")).MapAssign().ID("nil"),
					)).
					Dotln("ToSql").Call(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			)),
			jen.Line(),
			jen.Return().IDf("all%sCountQuery", pn),
		),
		jen.Line(),
	}
}

func buildGetAllSomethingCountFuncDecl(dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Commentf("GetAll%sCount will fetch the count of %s from the database.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("GetAll%sCount", pn).Params(constants.CtxParam()).Params(jen.ID("count").Uint64(), jen.Err().Error()).Block(
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(
				constants.CtxVar(),
				jen.ID(dbfl).Dotf("buildGetAll%sCountQuery", pn).Call(),
			).Dot("Scan").Call(jen.AddressOf().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
		),
		jen.Line(),
	}
}

func buildGetBatchOfSomethingQueryFuncDecl(_ *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	pn := typ.Name.Plural()
	sn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Commentf("buildGetBatchOf%sQuery returns a query that fetches every %s in the database within a bucketed range.", sn, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildGetBatchOf%sQuery", pn).Params(
			jen.List(jen.ID("beginID"), jen.ID("endID")).Uint64(),
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Spread()).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(
				jen.Qual(squirrelPkg, "Gt").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("idColumn")).MapAssign().ID("beginID"),
				)).
				Dotln("Where").Call(
				jen.Qual(squirrelPkg, "Lt").Valuesln(
					jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("idColumn")).MapAssign().ID("endID"),
				)).Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildGetAllOfSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Commentf("GetAll%s fetches every %s from the database and writes them to a channel. This method primarily exists", pn, scn),
		jen.Line(),
		jen.Comment("to aid in administrative data tasks."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("GetAll%s", pn).Params(
			constants.CtxParam(),
			jen.ID("resultChannel").Chan().Index().Qual(proj.ModelsV1Package(), sn),
		).Params(jen.Error()).Block(
			jen.List(jen.ID("count"), jen.Err()).Assign().ID(dbfl).Dotf("GetAll%sCount", pn).Call(constants.CtxVar()),
			jen.If(jen.Err().DoesNotEqual().Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Line(),
			jen.For(
				jen.ID("beginID").Assign().Uint64().Call(jen.One()),
				jen.ID("beginID").LessThanOrEqual().ID("count"),
				jen.ID("beginID").PlusEquals().ID("defaultBucketSize"),
			).Block(
				jen.ID("endID").Assign().ID("beginID").Plus().ID("defaultBucketSize"),
				jen.Go().Func().Params(jen.List(jen.ID("begin"), jen.ID("end")).Uint64()).Block(
					jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGetBatchOf%sQuery", pn).Call(jen.ID("begin"), jen.ID("end")),
					jen.ID(constants.LoggerVarName).Assign().ID(dbfl).Dot(constants.LoggerVarName).Dot("WithValues").Call(jen.Map(jen.String()).Interface().Valuesln(
						jen.Lit("query").MapAssign().ID("query"),
						jen.Lit("begin").MapAssign().ID("begin"),
						jen.Lit("end").MapAssign().ID("end"),
					)),
					jen.Line(),
					jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("Query").Call(jen.ID("query"), jen.ID("args").Spread()),
					jen.If(jen.Err().IsEqualTo().Qual("database/sql", "ErrNoRows")).Block(
						jen.Return(),
					).Else().If(jen.Err().DoesNotEqual().Nil()).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("querying for database rows")),
						jen.Return(),
					),
					jen.Line(),
					jen.List(jen.ID(puvn), jen.Err()).Assign().ID(dbfl).Dotf("scan%s", pn).Call(jen.ID("rows")),
					jen.If(jen.Err().DoesNotEqual().Nil()).Block(
						jen.ID(constants.LoggerVarName).Dot("Error").Call(jen.Err(), jen.Lit("scanning database rows")),
						jen.Return(),
					),
					jen.Line(),
					jen.ID("resultChannel").ReceiveFromChannel().ID(puvn),
				).Call(jen.ID("beginID"), jen.ID("endID")),
			),
			jen.Line(),
			jen.Return(jen.Nil()),
		),
		jen.Line(),
	}
}

func buildGetListOfSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildDBQuerierListRetrievalQueryBuildingMethodParams(proj)

	var firstCommentLine string
	if (typ.BelongsToUser && typ.RestrictedToUser) && typ.BelongsToStruct != nil {
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter, and belong to a given user and %s,", pn, pcn, typ.BelongsToStruct.SingularCommonName())
	} else if typ.BelongsToUser && typ.RestrictedToUser {
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given user,", pn, pcn)
	} else if typ.BelongsToStruct != nil {
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given %s,", pn, pcn, typ.BelongsToStruct.SingularCommonName())
	} else {
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter,", pn, pcn)
	}

	whereValues := typ.BuildDBQuerierListRetrievalQueryMethodConditionalClauses(proj)
	qbStmt := jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
		Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Spread()).
		Dotln("From").Call(jen.IDf("%sTableName", puvn))

	qbStmt = typ.ModifyQueryBuildingStatementWithJoinClauses(proj, qbStmt)

	if len(whereValues) > 0 {
		qbStmt = qbStmt.Dotln("Where").Call(jen.Qual(squirrelPkg, "Eq").Valuesln(whereValues...))
	}
	qbStmt = qbStmt.Dotln("OrderBy").Call(utils.FormatString("%s.%s", jen.IDf("%sTableName", puvn), jen.ID("idColumn")))

	return []jen.Code{
		jen.Comment(firstCommentLine),
		jen.Line(),
		jen.Commentf("and returns both the query and the relevant args to pass to the query executor."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildGet%sQuery", pn).Params(
			params...,
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			qbStmt,
			jen.Line(),
			jen.If(jen.ID(constants.FilterVarName).DoesNotEqual().ID("nil")).Block(
				jen.ID("builder").Equals().ID(constants.FilterVarName).Dot("ApplyToQueryBuilder").Call(jen.ID("builder"), jen.IDf("%sTableName", puvn)),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildGetListOfSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildDBQuerierListRetrievalMethodParams(proj)
	queryBuildingParams := typ.BuildDBQuerierListRetrievalMethodArgs(proj)

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("Get%s", pn).Params(
			params...,
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)), jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sQuery", pn).Call(queryBuildingParams...),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Litf("querying database for %s", pcn))),
			),
			jen.Line(),
			jen.List(jen.ID(puvn), jen.Err()).Assign().ID(dbfl).Dotf("scan%s", pn).Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("list").Assign().AddressOf().Qual(proj.ModelsV1Package(), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").MapAssign().Qual(proj.ModelsV1Package(), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID(constants.FilterVarName).Dot("Page"),
					jen.ID("Limit").MapAssign().ID(constants.FilterVarName).Dot("Limit"),
				),
				jen.ID(pn).MapAssign().ID(puvn),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}
}

func buildGetListOfSomethingWithIDsQueryFuncDecl(_ *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	pn := typ.Name.Plural()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	var firstCommentLine string
	if typ.BelongsToUser && typ.RestrictedToUser {
		firstCommentLine = fmt.Sprintf("buildGet%sWithIDsQuery builds a SQL query selecting %s that belong to a given user,", pn, puvn)
	} else if typ.BelongsToStruct != nil {
		firstCommentLine = fmt.Sprintf("buildGet%sWithIDsQuery builds a SQL query selecting %s that belong to a given %s,", pn, puvn, typ.BelongsToStruct.SingularCommonName())
	} else {
		firstCommentLine = fmt.Sprintf("buildGet%sWithIDsQuery builds a SQL query selecting %s", pn, puvn)
	}

	var queryBuilderStmt jen.Code
	if isPostgres(dbvendor) {
		queryBuilderStmt = jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
			Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Spread()).
			Dotln("FromSelect").Call(jen.ID("subqueryBuilder"), jen.IDf("%sTableName", puvn)).
			Dotln("Where").Call(
			jen.Qual(squirrelPkg, "Eq").Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				func() jen.Code {
					if typ.BelongsToUser {
						return jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sUserOwnershipColumn", puvn)).MapAssign().ID("userID")
					}
					return jen.Null()
				}(),
			))
	} else if isMariaDB(dbvendor) || isSqlite(dbvendor) {
		queryBuilderStmt = jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
			Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Spread()).
			Dotln("From").Call(jen.IDf("%sTableName", puvn)).
			Dotln("Where").Call(
			jen.Qual(squirrelPkg, "Eq").Valuesln(
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("idColumn")).MapAssign().ID("ids"),
				jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.ID("archivedOnColumn")).MapAssign().Nil(),
				func() jen.Code {
					if typ.BelongsToUser {
						return jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sUserOwnershipColumn", puvn)).MapAssign().ID("userID")
					}
					return jen.Null()
				}(),
			)).
			Dotln("OrderBy").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("CASE %s.%s %s"),
				jen.IDf("%sTableName", puvn),
				jen.ID("idColumn"),
				jen.ID("whenThenStatement"),
			)).
			Dotln("Limit").Call(jen.Uint64().Call(jen.ID("limit")))
	}

	return []jen.Code{
		jen.Comment(firstCommentLine),
		jen.Line(),
		jen.Comment("and have IDs that exist within a given set of IDs. Returns both the query and the relevant"),
		jen.Line(),
		jen.Comment("args to pass to the query executor. This function is primarily intended for use with a search"),
		jen.Line(),
		jen.Comment("index, which would provide a slice of string IDs to query against. This function accepts a"),
		jen.Line(),
		jen.Comment("slice of uint64s instead of a slice of strings in order to ensure all the provided strings"),
		jen.Line(),
		jen.Comment("are valid database IDs, because there's no way in squirrel to escape them in the unnest join,"),
		jen.Line(),
		jen.Comment("and if we accept strings we could leave ourselves vulnerable to SQL injection attacks."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildGet%sWithIDsQuery", pn).Params(
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.ID("userID").Uint64()
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Uint64()
				}
				return jen.Null()
			}(),
			jen.ID("limit").Uint8(),
			jen.ID("ids").Index().Uint64(),
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			func() jen.Code {
				if isPostgres(dbvendor) {
					return jen.ID("subqueryBuilder").Assign().ID(dbfl).Dot("sqlBuilder").Dot("Select").Call(jen.IDf("%sTableColumns", puvn).Spread()).
						Dotln("From").Call(jen.IDf("%sTableName", puvn)).
						Dotln("Join").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("unnest('{%s}'::int[])"), jen.ID("joinUint64s").Call(jen.ID("ids")))).
						Dotln("Suffix").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d"), jen.ID("limit")))
				} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
					return jen.Var().ID("whenThenStatement").String()
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if isSqlite(dbvendor) || isMariaDB(dbvendor) {
					return jen.Line()
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if isSqlite(dbvendor) || isMariaDB(dbvendor) {
					return jen.For(jen.List(jen.ID("i"), jen.ID("id")).Assign().Range().ID("ids")).Block(
						jen.If(jen.ID("i").DoesNotEqual().Zero()).Block(
							jen.ID("whenThenStatement").PlusEquals().Lit(" "),
						),
						jen.ID("whenThenStatement").PlusEquals().Qual("fmt", "Sprintf").Call(
							jen.Lit("WHEN %d THEN %d"),
							jen.ID("id"),
							jen.ID("i"),
						),
					)
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if isSqlite(dbvendor) || isMariaDB(dbvendor) {
					return jen.ID("whenThenStatement").PlusEquals().Lit(" END")
				}
				return jen.Null()
			}(),
			queryBuilderStmt,
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildGetListOfSomethingWithIDsFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Commentf("Get%sWithIDs fetches a list of %s from the database that exist within a given set of IDs.", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("Get%sWithIDs", pn).Params(
			constants.CtxParam(),
			func() jen.Code {
				if typ.BelongsToUser {
					return jen.ID("userID").Uint64()
				}
				return jen.Null()
			}(),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName()).Uint64()
				}
				return jen.Null()
			}(),
			jen.ID("limit").Uint8(),
			jen.ID("ids").Index().Uint64(),
		).Params(
			jen.Index().Qual(proj.ModelsV1Package(), sn),
			jen.Error(),
		).Block(
			jen.If(jen.ID("limit").IsEqualTo().Zero()).Block(
				jen.ID("limit").Equals().Uint8().Call(jen.Qual(proj.ModelsV1Package(), "DefaultLimit")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sWithIDsQuery", pn).Call(
				func() jen.Code {
					if typ.BelongsToUser {
						return jen.ID("userID")
					}
					return jen.Null()
				}(),
				func() jen.Code {
					if typ.BelongsToStruct != nil {
						return jen.IDf("%sID", typ.BelongsToStruct.UnexportedVarName())
					}
					return jen.Null()
				}(),
				jen.ID("limit"),
				jen.ID("ids"),
			),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Litf("querying database for %s", pcn))),
			),
			jen.Line(),
			jen.List(jen.ID(puvn), jen.Err()).Assign().ID(dbfl).Dotf("scan%s", pn).Call(jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID(puvn), jen.Nil()),
		),
		jen.Line(),
	}
}

func determineCreationColumns(_ *models.Project, typ models.DataType) []jen.Code {
	puvn := typ.Name.PluralUnexportedVarName()
	var creationColumns []jen.Code

	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, jen.IDf("%sTable%sColumn", puvn, field.Name.Singular()))
	}

	if typ.BelongsToStruct != nil {
		creationColumns = append(creationColumns, jen.IDf("%sTableOwnershipColumn", puvn))
	}
	if typ.BelongsToUser {
		creationColumns = append(creationColumns, jen.IDf("%sUserOwnershipColumn", puvn))
	}

	return creationColumns
}

func determineCreationQueryValues(_ *models.Project, inputVarName string, typ models.DataType) []jen.Code {
	var valuesColumns []jen.Code

	for _, field := range typ.Fields {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dot(field.Name.Singular()))
	}

	if typ.BelongsToStruct != nil {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dot(constants.UserOwnershipFieldName))
	}

	return valuesColumns
}

func buildCreateSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	qb := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Insert").Call(jen.IDf("%sTableName", puvn)).
		Dotln("Columns").Callln(determineCreationColumns(proj, typ)...).
		Dotln("Values").Callln(determineCreationQueryValues(proj, "input", typ)...)

	if isPostgres(dbvendor) {
		qb.Dotln("Suffix").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("RETURNING %s, %s"),
				jen.ID("idColumn"),
				jen.ID("createdOnColumn"),
			),
		)
	}
	qb.Dotln("ToSql").Call()

	params := typ.BuildDBQuerierCreationQueryBuildingMethodParams(proj, false)[1:]

	createQueryFuncBody := []jen.Code{
		jen.Var().Err().Error(),
		jen.Line(),
		qb,
		jen.Line(),
		jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
		jen.Line(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	}

	return []jen.Code{
		jen.Commentf("buildCreate%sQuery takes %s and returns a creation query for that %s and the relevant arguments.", sn, scnwp, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildCreate%sQuery", sn).Params(
			params...,
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			createQueryFuncBody...,
		),
		jen.Line(),
	}
}

func buildCreateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	var (
		createInitColumns []jen.Code
	)

	params := typ.BuildDBQuerierCreationMethodParams(proj)
	queryBuildingArgs := typ.BuildDBQuerierCreationMethodQueryBuildingArgs(proj)
	queryBuildingArgs = queryBuildingArgs[:len(queryBuildingArgs)-1]
	queryBuildingArgs = append(queryBuildingArgs, jen.ID("x"))

	for _, field := range typ.Fields {
		fn := field.Name.Singular()
		createInitColumns = append(createInitColumns, jen.ID(fn).MapAssign().ID("input").Dot(fn))
	}

	if typ.BelongsToStruct != nil {
		createInitColumns = append(createInitColumns, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser {
		createInitColumns = append(createInitColumns, jen.ID(constants.UserOwnershipFieldName).MapAssign().ID("input").Dot(constants.UserOwnershipFieldName))
	}

	baseCreateFuncBody := []jen.Code{
		jen.ID("x").Assign().AddressOf().Qual(proj.ModelsV1Package(), sn).Valuesln(createInitColumns...),
		jen.Line(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildCreate%sQuery", sn).Call(
			queryBuildingArgs...,
		),
		jen.Line(),
		jen.Commentf("create the %s.", scn),
	}

	if isPostgres(dbvendor) {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()).Dot("Scan").Call(jen.AddressOf().ID("x").Dot("ID"), jen.AddressOf().ID("x").Dot("CreatedOn")),
		)
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.List(jen.ID(constants.ResponseVarName), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
		)
	} else {
		panic(fmt.Sprintf("dbrn is weird: %q", dbrn))
	}

	baseCreateFuncBody = append(baseCreateFuncBody,
		jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
			jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("error executing %s creation query: ", scn)+"%w"), jen.Err())),
		),
		jen.Line(),
	)

	if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.Comment("fetch the last inserted ID."),
			jen.List(jen.ID("id"), jen.ID("err")).Assign().ID(constants.ResponseVarName).Dot("LastInsertId").Call(),
			jen.ID(dbfl).Dot("logIDRetrievalError").Call(jen.Err()),
			jen.ID("x").Dot("ID").Equals().Uint64().Call(jen.ID("id")),
			jen.Line(),
			jen.Comment("this won't be completely accurate, but it will suffice."),
			jen.ID("x").Dot("CreatedOn").Equals().ID(dbfl).Dot("timeTeller").Dot("Now").Call(),
		)
	}
	baseCreateFuncBody = append(baseCreateFuncBody, jen.Line(), jen.Return().List(jen.ID("x"), jen.Nil()))

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.PointerTo().Qual(proj.ModelsV1Package(), sn), jen.Error()).Block(
			baseCreateFuncBody...,
		),
		jen.Line(),
	}
}

func buildUpdateSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	puvn := typ.Name.PluralUnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	const inputVarName = "input"
	params := typ.BuildDBQuerierUpdateQueryBuildingMethodParams(proj, inputVarName)

	vals := []jen.Code{
		jen.ID("idColumn").MapAssign().ID(inputVarName).Dot("ID"),
	}
	if typ.BelongsToStruct != nil {
		vals = append(vals, jen.IDf("%sTableOwnershipColumn", puvn).MapAssign().ID(inputVarName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToUser {
		vals = append(vals, jen.IDf("%sUserOwnershipColumn", puvn).MapAssign().ID(inputVarName).Dot(constants.UserOwnershipFieldName))
	}

	return []jen.Code{
		jen.Commentf("buildUpdate%sQuery takes %s and returns an update SQL query, with the relevant query parameters.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildUpdate%sQuery", sn).Params(
			params...,
		).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			jen.Var().Err().Error(),
			jen.Line(),
			func(typ models.DataType) jen.Code {
				x := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Update").Call(jen.IDf("%sTableName", puvn))

				for _, field := range typ.Fields {
					if field.ValidForUpdateInput {
						x.Dotln("Set").Call(jen.IDf("%sTable%sColumn", puvn, field.Name.Singular()), jen.ID(inputVarName).Dot(field.Name.Singular()))
					}
				}

				x.Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual(squirrelPkg, "Expr").Call(jen.ID("currentUnixTimeQuery"))).
					Dotln("Where").Call(jen.Qual(squirrelPkg, "Eq").Valuesln(vals...))

				if strings.ToLower(dbvsn) == "postgres" {
					x.Dotln("Suffix").Call(
						jen.Qual("fmt", "Sprintf").Call(
							jen.Lit("RETURNING %s"),
							jen.ID("lastUpdatedOnColumn"),
						),
					)
				}

				x.Dotln("ToSql").Call()
				return x
			}(typ),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildUpdateSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	n := typ.Name
	sn := n.Singular()
	scn := n.SingularCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	const updatedVarName = "input"

	args := typ.BuildDBQuerierUpdateMethodArgs(proj, updatedVarName)
	block := []jen.Code{
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildUpdate%sQuery", sn).Call(args...),
	}

	if isPostgres(dbvendor) {
		block = append(block,
			jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(
				constants.CtxVar(),
				jen.ID("query"),
				jen.ID("args").Spread(),
			).Dot("Scan").Call(
				jen.AddressOf().ID(updatedVarName).Dot("LastUpdatedOn"),
			),
		)
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		block = append(block,
			jen.List(jen.Underscore(), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.Return(jen.Err()),
		)
	}

	params := typ.BuildDBQuerierUpdateMethodParams(proj, updatedVarName)

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the provided input to have a valid ID.", sn, scn, sn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("Update%s", sn).Params(params...).Params(jen.Error()).Block(block...),
	}
}

func buildArchiveSomethingQueryFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	var comment string
	vals := []jen.Code{
		jen.ID("idColumn").MapAssign().IDf("%sID", uvn),
		jen.ID("archivedOnColumn").MapAssign().ID("nil"),
	}
	paramsList := typ.BuildDBQuerierArchiveQueryMethodParams(proj)

	if typ.BelongsToNobody {
		comment = fmt.Sprintf("buildArchive%sQuery returns a SQL query which marks a given %s ", sn, scn)
	} else {
		comment = fmt.Sprintf("buildArchive%sQuery returns a SQL query which marks a given %s belonging to ", sn, scn)
		if typ.BelongsToStruct != nil {
			comment += fmt.Sprintf("a given %s ", typ.BelongsToStruct.SingularCommonName())
			vals = append(vals, jen.IDf("%sTableOwnershipColumn", puvn).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		}
		if typ.BelongsToUser {
			if typ.BelongsToStruct != nil {
				comment += "and a given user "
			} else {
				comment += "a given user "
			}
			vals = append(vals, jen.IDf("%sUserOwnershipColumn", puvn).MapAssign().ID(constants.UserIDVarName))
		}
	}
	comment += "as archived."

	_qs := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Update").Call(jen.IDf("%sTableName", puvn)).
		Dotln("Set").Call(jen.ID("lastUpdatedOnColumn"), jen.Qual(squirrelPkg, "Expr").Call(jen.ID("currentUnixTimeQuery"))).
		Dotln("Set").Call(jen.ID("archivedOnColumn"), jen.Qual(squirrelPkg, "Expr").Call(jen.ID("currentUnixTimeQuery"))).
		Dotln("Where").Call(jen.Qual(squirrelPkg, "Eq").Valuesln(vals...))

	if strings.ToLower(dbvsn) == "postgres" {
		_qs.Dotln("Suffix").Call(
			jen.Qual("fmt", "Sprintf").Call(
				jen.Lit("RETURNING %s"),
				jen.ID("archivedOnColumn"),
			),
		)
	}
	_qs.Dotln("ToSql").Call()

	archiveFuncBody := []jen.Code{
		jen.Var().Err().Error(),
		jen.Line(),
		_qs,
		jen.Line(),
		jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
		jen.Line(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	}

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("buildArchive%sQuery", sn).Params(
			jen.List(paramsList...)).Params(jen.ID("query").String(), jen.ID("args").Index().Interface()).Block(
			archiveFuncBody...,
		),
		jen.Line(),
	}
}

func buildArchiveSomethingFuncDecl(proj *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildDBQuerierArchiveMethodParams(proj)
	queryBuildingArgs := typ.BuildDBQuerierArchiveQueryBuildingArgs(proj)

	return []jen.Code{
		jen.Commentf("Archive%s marks %s as archived in the database.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).PointerTo().ID(dbvsn)).IDf("Archive%s", sn).Params(
			params...,
		).Params(jen.Error()).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildArchive%sQuery", sn).Call(
				queryBuildingArgs...,
			),
			jen.Line(),
			jen.List(jen.ID("res"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(constants.CtxVar(), jen.ID("query"), jen.ID("args").Spread()),
			jen.If(jen.ID("res").DoesNotEqual().Nil()).Block(
				jen.If(
					jen.List(jen.ID("rowCount"), jen.ID("rowCountErr")).Assign().ID("res").Dot("RowsAffected").Call(),
					jen.ID("rowCountErr").IsEqualTo().Nil().And().ID("rowCount").IsEqualTo().Zero(),
				).Block(
					jen.Return(jen.Qual("database/sql", "ErrNoRows")),
				),
			),
			jen.Line(),
			jen.Return(jen.Err()),
		),
		jen.Line(),
	}
}
