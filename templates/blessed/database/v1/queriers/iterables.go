package queriers

import (
	"fmt"
	"path/filepath"
	"strings"

	jen "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	utils "gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

func isPostgres(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.RouteName() == "postgres"
}

func isSqlite(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.RouteName() == "sqlite"
}

func isMariaDB(dbvendor wordsmith.SuperPalabra) bool {
	return dbvendor.RouteName() == "mariadb" || dbvendor.RouteName() == "maria_db"
}

func buildIterableConstants(typ models.DataType) []jen.Code {
	n := typ.Name
	prn := n.PluralRouteName()
	puvn := n.PluralUnexportedVarName()

	consts := []jen.Code{
		jen.IDf("%sTableName", puvn).Equals().Lit(prn),
	}

	if typ.BelongsToUser {
		consts = append(consts, jen.IDf("%sTableOwnershipColumn", puvn).Equals().Lit("belongs_to_user"))
	}
	if typ.BelongsToStruct != nil {
		consts = append(consts, jen.IDf("%sTableOwnershipColumn", puvn).Equals().Litf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
	}

	return consts
}

func buildIterableVariableDecs(typ models.DataType) []jen.Code {
	puvn := typ.Name.PluralUnexportedVarName()

	vars := []jen.Code{
		jen.Var().Defs(
			jen.IDf("%sTableColumns", puvn).Equals().Index().ID("string").Valuesln(buildTableColumns(typ)...),
		),
		jen.Line(),
	}

	return vars
}

func iterablesDotGo(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(pkg, ret)

	ret.Add(jen.Const().Defs(buildIterableConstants(typ)...), jen.Line())
	ret.Add(buildIterableVariableDecs(typ)...)
	ret.Add(buildScanSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildScanListOfSomethingFuncDecl(pkg, typ)...)
	ret.Add(buildSomethingExistsQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildSomethingExistsFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildGetSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildGetSomethingFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildGetSomethingCountQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildGetSomethingCountFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildSomethingAllCountQueryDecls(dbvendor, typ)...)
	ret.Add(buildGetAllSomethingCountFuncDecl(dbvendor, typ)...)
	ret.Add(buildGetListOfSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildGetListOfSomethingFuncDecl(pkg, dbvendor, typ)...)

	if typ.BelongsToUser {
		ret.Add(buildGetAllSomethingForUserFuncDecl(pkg, dbvendor, typ)...)
	}
	if typ.BelongsToStruct != nil {
		ret.Add(buildGetAllSomethingForSomethingElseFuncDecl(pkg, dbvendor, typ)...)
	}

	ret.Add(buildCreateSomethingQueryFuncDecl(pkg, dbvendor, typ)...)

	if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		ret.Add(buildSomethingCreationTimeQueryFuncDecl(dbvendor, typ)...)
	}

	ret.Add(buildCreateSomethingFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildUpdateSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildUpdateSomethingFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildArchiveSomethingQueryFuncDecl(pkg, dbvendor, typ)...)
	ret.Add(buildArchiveSomethingFuncDecl(pkg, dbvendor, typ)...)

	return ret
}

func buildTableColumns(typ models.DataType) []jen.Code {
	puvn := typ.Name.PluralUnexportedVarName()
	tableNameVar := fmt.Sprintf("%sTableName", puvn)

	tableColumns := []jen.Code{jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.Lit("id"))}

	for _, field := range typ.Fields {
		tableColumns = append(tableColumns, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.Lit(field.Name.RouteName())))
	}

	tableColumns = append(tableColumns,
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.Lit("created_on")),
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.Lit("updated_on")),
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.Lit("archived_on")),
	)

	if typ.BelongsToUser {
		tableColumns = append(tableColumns, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.IDf("%sTableOwnershipColumn", puvn)))
	}
	if typ.BelongsToStruct != nil {
		tableColumns = append(tableColumns, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf(tableNameVar), jen.IDf("%sTableOwnershipColumn", puvn)))
	}

	return tableColumns
}

func buildScanFields(typ models.DataType) (scanFields []jen.Code) {
	scanFields = []jen.Code{jen.VarPointer().ID("x").Dot("ID")}

	for _, field := range typ.Fields {
		scanFields = append(scanFields, jen.VarPointer().ID("x").Dot(field.Name.Singular()))
	}

	scanFields = append(scanFields,
		jen.VarPointer().ID("x").Dot("CreatedOn"),
		jen.VarPointer().ID("x").Dot("UpdatedOn"),
		jen.VarPointer().ID("x").Dot("ArchivedOn"),
	)

	if typ.BelongsToUser {
		scanFields = append(scanFields, jen.VarPointer().ID("x").Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		scanFields = append(scanFields, jen.VarPointer().ID("x").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	return scanFields
}

// what's the difference between these two things
func buildScanSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pscnwp := typ.Name.ProperSingularCommonNameWithPrefix()

	return []jen.Code{
		jen.Commentf("scan%s takes a database Scanner (i.e. *sql.Row) and scans the result into %s struct", sn, pscnwp),
		jen.Line(),
		jen.Func().IDf("scan%s", sn).Params(jen.ID("scan").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Scanner")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			func() []jen.Code {
				body := []jen.Code{
					jen.ID("x").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
					jen.Line(),
					jen.If(jen.Err().Assign().ID("scan").Dot("Scan").Callln(buildScanFields(typ)...), jen.Err().DoesNotEqual().ID("nil")).Block(
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

func buildScanListOfSomethingFuncDecl(pkg *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()

	return []jen.Code{
		jen.Commentf("scan%s takes a logger and some database rows and turns them into a slice of %s", pn, pcn),
		jen.Line(),
		jen.Func().IDf("scan%s", pn).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("rows").ParamPointer().Qual("database/sql", "Rows")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("x"), jen.Err()).Assign().IDf("scan%s", sn).Call(jen.ID("rows")),
				jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
					jen.Return().List(jen.Nil(), jen.Err()),
				),
				jen.ID("list").Equals().ID("append").Call(jen.ID("list"), jen.Op("*").ID("x")),
			),
			jen.Line(),
			jen.If(jen.Err().Assign().ID("rows").Dot("Err").Call(), jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Err()),
			),
			jen.Line(),
			jen.If(jen.ID("closeErr").Assign().ID("rows").Dot("Close").Call(), jen.ID("closeErr").DoesNotEqual().ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("closeErr"), jen.Lit("closing database rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}
}

func buildSomethingExistsQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	var (
		comment string
	)

	dbvsn := dbvendor.Singular()
	n := typ.Name
	sn := n.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	scnwp := n.SingularCommonNameWithPrefix()
	uvn := n.UnexportedVarName()
	puvn := n.PluralUnexportedVarName()

	params := typ.BuildGetSomethingParams(pkg)[1:]
	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.IDf("%sTableName", puvn)).MapAssign().IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		comment = fmt.Sprintf("build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a user with a given ID exists.", sn, scnwp)
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID("userID"))
	}
	if typ.BelongsToStruct != nil {
		comment = fmt.Sprintf("build%sExistsQuery constructs a SQL query for checking if %s with a given ID belong to a %s with a given ID exists.", sn, scnwp, typ.BelongsToStruct.SingularCommonNameWithPrefix())
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID(fmt.Sprintf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	} else if typ.BelongsToNobody {
		comment = fmt.Sprintf("build%sExistsQuery constructs a SQL query for checking if %s with a given ID exists.", sn, scnwp)
	}

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("build%sExistsQuery", sn).Params(params...).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.IDf("%sTableName", puvn))).
				Dotln("Prefix").Call(jen.ID("existencePrefix")).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Suffix").Call(jen.ID("existenceSuffix")).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(whereValues...)).
				Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildSomethingExistsFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	cn := typ.Name.SingularCommonName()
	sn := typ.Name.Singular()

	const existenceVarName = "exists"

	params := typ.BuildGetSomethingParams(pkg)
	buildQueryParams := typ.BuildGetSomethingArgs(pkg)[1:]

	return []jen.Code{
		jen.Commentf("%sExists queries the database to see if a given %s belonging to a given user exists", sn, cn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("%sExists", sn).Params(params...).Params(jen.Bool(), jen.ID("error")).Block(
			jen.Var().ID(existenceVarName).Bool(),
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("build%sExistsQuery", sn).Call(buildQueryParams...),
			jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.VarPointer().ID(existenceVarName)),
			jen.Return().List(jen.ID(existenceVarName), jen.Err()),
		),
		jen.Line(),
	}
}

func buildGetSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	var (
		comment string
	)

	dbvsn := dbvendor.Singular()
	n := typ.Name
	sn := n.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	uvn := n.UnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()
	puvn := n.PluralUnexportedVarName()

	params := typ.BuildGetSomethingParams(pkg)[1:]
	whereValues := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.IDf("%sTableName", puvn)).MapAssign().IDf("%sID", uvn),
	}

	if typ.BelongsToUser {
		comment = fmt.Sprintf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to a user with a given ID.", sn, scnwp)
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID("userID"))
	} else if typ.BelongsToStruct != nil {
		tsnwp := typ.BelongsToStruct.SingularCommonNameWithPrefix()
		comment = fmt.Sprintf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to %s with a given ID.", sn, scnwp, tsnwp)
		whereValues = append(whereValues, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID(fmt.Sprintf("%sID", typ.BelongsToStruct.UnexportedVarName())))
	} else if typ.BelongsToNobody {
		comment = fmt.Sprintf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID.", sn, scnwp)
	}

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGet%sQuery", sn).Params(params...).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Op("...")).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(whereValues...)).
				Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	}
}

func buildGetSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	sn := typ.Name.Singular()

	params := typ.BuildGetSomethingParams(pkg)
	buildQueryParams := typ.BuildGetSomethingArgs(pkg)[1:]

	return []jen.Code{
		jen.Commentf("Get%s fetches %s from the %s database", sn, scnwp, dbvendor.SingularPackageName()),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Get%s", sn).Params(params...).
			Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sQuery", sn).Call(buildQueryParams...),
			jen.ID("row").Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().IDf("scan%s", sn).Call(jen.ID("row")),
		),
		jen.Line(),
	}
}

func buildGetSomethingCountQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()

	var (
		commentOne string
		commentTwo string
	)

	params := typ.BuildGetListOfSomethingParams(pkg, false)[1:]
	vals := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.IDf("%sTableName", puvn)).MapAssign().Nil(),
	}

	if typ.BelongsToUser && typ.BelongsToStruct != nil {
		vals = append(vals, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID("userID"))
		vals = append(vals, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))

		btsscnwp := typ.BelongsToStruct.SingularCommonNameWithPrefix()
		btspcn := typ.BelongsToStruct.PluralCommonName()

		commentOne = fmt.Sprintf("buildGet%sCountQuery takes a QueryFilter, a user ID, and %s ID and returns a SQL query", sn, btsscnwp)
		commentTwo = fmt.Sprintf("(and the relevant arguments) for fetching the number of %s belonging to a given %s", pcn, btspcn)
	} else if typ.BelongsToUser {
		vals = append(vals, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID("userID"))

		commentOne = fmt.Sprintf("buildGet%sCountQuery takes a QueryFilter and a user ID, and returns a SQL query", sn)
		commentTwo = fmt.Sprintf("(and the relevant arguments) for fetching the number of %s belonging to a given user", pcn)
	} else if typ.BelongsToStruct != nil {
		vals = append(vals, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))

		commentOne = fmt.Sprintf("buildGet%sCountQuery takes a QueryFilter, %s ID and returns a SQL query", sn, scn)
		commentTwo = fmt.Sprintf("(and the relevant arguments) for fetching the number of %s belonging to a given %s", pcn, typ.BelongsToStruct.SingularCommonName())
	} else {
		commentOne = fmt.Sprintf("buildGet%sCountQuery takes a QueryFilter, and returns a SQL query", sn)
		commentTwo = fmt.Sprintf("(and the relevant arguments) for fetching the number of %s", pcn)
	}

	return []jen.Code{
		jen.Comment(commentOne),
		jen.Line(),
		jen.Comment(commentTwo),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGet%sCountQuery", sn).Params(params...).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(jen.ID("countQuery"), jen.IDf("%sTableName", puvn))).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(vals...)),
			jen.Line(),
			jen.If(jen.ID(utils.FilterVarName).DoesNotEqual().ID("nil")).Block(
				jen.ID("builder").Equals().ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
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

func buildGetSomethingCountFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	pcn := typ.Name.PluralCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	var comment string
	params := typ.BuildGetListOfSomethingParams(pkg, false)
	queryBuildingParams := typ.BuildGetListOfSomethingArgs(pkg)[1:]

	if typ.BelongsToUser {
		comment = fmt.Sprintf("Get%sCount will fetch the count of %s from the database that meet a particular filter and belong to a particular user.", sn, pcn)
	}
	if typ.BelongsToStruct != nil {
		comment = fmt.Sprintf("Get%sCount will fetch the count of %s from the database that meet a particular filter and belongs to a particular %s.", sn, pcn, typ.BelongsToStruct.SingularCommonName())
	} else if typ.BelongsToNobody {
		comment = fmt.Sprintf("Get%sCount will fetch the count of %s from the database that meet a particular filter.", sn, pcn)
	}

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Get%sCount", sn).Params(
			params...,
		).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sCountQuery", sn).Call(queryBuildingParams...),
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.VarPointer().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
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
			jen.IDf("all%sCountQuery", pn).ID("string"),
		),
		jen.Line(),
		jen.Line(),
		jen.Commentf("buildGetAll%sCountQuery returns a query that fetches the total number of %s in the database.", pn, pcn),
		jen.Line(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGetAll%sCountQuery", pn).Params().Params(jen.ID("string")).Block(
			jen.IDf("all%sCountQueryBuilder", pn).Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.IDf("all%sCountQuery", pn), jen.ID("_"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(jen.ID("countQuery"), jen.IDf("%sTableName", puvn))).
					Dotln("From").Call(jen.IDf("%sTableName", puvn)).
					Dotln("Where").Call(
					jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
						jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.IDf("%sTableName", puvn)).MapAssign().ID("nil"),
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
		jen.Commentf("GetAll%sCount will fetch the count of %s from the database", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("GetAll%sCount", pn).Params(utils.CtxVar().Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.Err().ID("error")).Block(
			jen.Err().Equals().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID(dbfl).Dotf("buildGetAll%sCountQuery", pn).Call()).Dot("Scan").Call(jen.VarPointer().ID("count")),
			jen.Return().List(jen.ID("count"), jen.Err()),
		),
		jen.Line(),
	}
}

func buildGetListOfSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	pn := typ.Name.Plural()
	pcn := typ.Name.PluralCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	/*
		CREATE TABLE IF NOT EXISTS forums (
			"id" BIGSERIAL NOT NULL PRIMARY KEY,
			"name" CHARACTER VARYING NOT NULL,
			"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
			"updated_on" BIGINT DEFAULT NULL,
			"archived_on" BIGINT DEFAULT NULL
		);


		CREATE TABLE IF NOT EXISTS threads (
			"id" BIGSERIAL NOT NULL PRIMARY KEY,
			"title" CHARACTER VARYING NOT NULL,
			"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
			"updated_on" BIGINT DEFAULT NULL,
			"archived_on" BIGINT DEFAULT NULL,
			"belongs_to_forum" BIGINT NOT NULL,
			FOREIGN KEY ("belongs_to_forum") REFERENCES "forums"("id")
		);


		CREATE TABLE IF NOT EXISTS comments (
			"id" BIGSERIAL NOT NULL PRIMARY KEY,
			"content" CHARACTER VARYING NOT NULL,
			"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
			"updated_on" BIGINT DEFAULT NULL,
			"archived_on" BIGINT DEFAULT NULL,
			"belongs_to_thread" BIGINT NOT NULL,
			FOREIGN KEY ("belongs_to_thread") REFERENCES "threads"("id")
		);

		INSERT INTO forums (name) VALUES ('forum_a');
		INSERT INTO forums (name) VALUES ('forum_b');

		INSERT INTO threads (title, belongs_to_forum) VALUES('thread_a', 2);
		INSERT INTO threads (title, belongs_to_forum) VALUES('thread_b', 2);
		INSERT INTO threads (title, belongs_to_forum) VALUES('thread_c', 2);

		INSERT INTO comments (content, belongs_to_thread) VALUES ('hello world', 3);
		INSERT INTO comments (content, belongs_to_thread) VALUES ('hello world', 3);
		INSERT INTO comments (content, belongs_to_thread) VALUES ('hello world', 3);
		INSERT INTO comments (content, belongs_to_thread) VALUES ('hello world', 3);
		INSERT INTO comments (content, belongs_to_thread) VALUES ('hello world', 3);
	*/

	/*
		SELECT comments.content, comments.created_on, comments.belongs_to_thread FROM comments
			INNER JOIN threads ON comments.belongs_to_thread=threads.id
			INNER JOIN forums ON threads.belongs_to_forum=forums.id
			WHERE threads.ID = 3 AND forums.id = 2;
	*/

	/*
		queryParts := []string{
			`SELECT comments.content, comments.created_on, comments.belongs_to_thread FROM comments`,
			`JOIN threads ON comments.belongs_to_thread=threads.id`,
			`JOIN forums ON threads.belongs_to_forum=forums.id`,
			`WHERE forums.id = $1 AND threads.id = $2`,
		}

		sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		q, _, err := sqlBuilder.
			Select(
				"comments.content",
				"comments.created_on",
				"comments.belongs_to_thread",
			).
			From("comments").
			Join("threads ON comments.belongs_to_thread=threads.id").
			Join("forums ON threads.belongs_to_forum=forums.id").
			Where(squirrel.Eq{
				"threads.id": 3,
				"forums.id":  2,
			}).
			ToSql()

		if err != nil {
			log.Fatal(err)
		}
		x := strings.Join(queryParts, " ")

		fmt.Println(q == x)
		fmt.Println(q)
		fmt.Println(x)
	*/

	vals := []jen.Code{
		jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.archived_on"), jen.IDf("%sTableName", puvn)).MapAssign().ID("nil"),
	}
	params := typ.BuildGetListOfSomethingParams(pkg, false)[1:]

	var firstCommentLine string
	if typ.BelongsToUser && typ.BelongsToStruct != nil {
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter, and belong to a given user and %s,", pn, pcn, typ.BelongsToStruct.SingularCommonName())
	} else if typ.BelongsToUser {
		vals = append(vals, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().ID("userID"))
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given user,", pn, pcn)
	} else if typ.BelongsToStruct != nil {
		vals = append(vals, jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.%s"), jen.IDf("%sTableName", puvn), jen.IDf("%sTableOwnershipColumn", puvn)).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given %s,", pn, pcn, typ.BelongsToStruct.SingularCommonName())
	} else {
		firstCommentLine = fmt.Sprintf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter,", pn, pcn)
	}

	return []jen.Code{
		jen.Comment(firstCommentLine),
		jen.Line(),
		jen.Commentf("and returns both the query and the relevant args to pass to the query executor."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGet%sQuery", pn).Params(
			params...,
		).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Assign().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Op("...")).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(vals...)),
			jen.Line(),
			jen.If(jen.ID(utils.FilterVarName).DoesNotEqual().ID("nil")).Block(
				jen.ID("builder").Equals().ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
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

func buildGetListOfSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	dbvsn := dbvendor.Singular()
	pcn := typ.Name.PluralCommonName()
	scn := typ.Name.SingularCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildGetListOfSomethingParams(pkg, false)
	queryBuildingParams := typ.BuildGetListOfSomethingArgs(pkg)[1:]
	countRetrievalParams := typ.BuildGetListOfSomethingArgs(pkg)

	return []jen.Code{
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Get%s", pn).Params(
			params...,
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sQuery", pn).Call(queryBuildingParams...),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Litf("querying database for %s", pcn))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().IDf("scan%s", pn).Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.Err())),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.Err()).Assign().ID(dbfl).Dotf("Get%sCount", sn).Call(countRetrievalParams...),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("fetching %s count: ", scn)+"%w"), jen.Err())),
			),
			jen.Line(),
			jen.ID("x").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").MapAssign().Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").MapAssign().ID("filter").Dot("Page"),
					jen.ID("Limit").MapAssign().ID("filter").Dot("Limit"),
					jen.ID("TotalCount").MapAssign().ID("count"),
				),
				jen.IDf(pn).MapAssign().ID("list"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.Nil()),
		),
		jen.Line(),
	}
}

func buildGetAllSomethingForUserFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildGetListOfSomethingParams(pkg, false)
	params = params[:len(params)-1]

	queryBuildingArgs := typ.BuildGetListOfSomethingArgs(pkg)[1:]
	queryBuildingArgs = queryBuildingArgs[:len(queryBuildingArgs)-1]
	queryBuildingArgs = append(queryBuildingArgs, jen.Nil())

	return []jen.Code{
		jen.Commentf("GetAll%sForUser fetches every %s belonging to a user", pn, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("GetAll%sForUser", pn).Params(
			params...,
		).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sQuery", pn).Call(
				queryBuildingArgs...,
			),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Litf("fetching %s for user", pcn))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().IDf("scan%s", pn).Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("parsing database results: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}
}

func buildGetAllSomethingForSomethingElseFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	pn := typ.Name.Plural()
	scn := typ.Name.SingularCommonName()
	pcn := typ.Name.PluralCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildGetSomethingForSomethingElseParams(pkg)
	queryBuildingArgs := typ.BuildGetSomethingForSomethingElseArgs(pkg)[1:]

	return []jen.Code{
		jen.Commentf("GetAll%sFor%s fetches every %s belonging to %s", pn, typ.BelongsToStruct.Singular(), scn, typ.BelongsToStruct.SingularCommonNameWithPrefix()),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("GetAll%sFor%s", pn, typ.BelongsToStruct.Singular()).Params(
			params...,
		).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildGet%sQuery", pn).Call(
				queryBuildingArgs...,
			),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("QueryContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.ID("buildError").Call(jen.Err(), jen.Litf("fetching %s for %s", pcn, typ.BelongsToStruct.SingularCommonName()))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.Err()).Assign().IDf("scan%s", pn).Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.Err().DoesNotEqual().ID("nil")).Block(
				jen.Return().List(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("parsing database results: %w"), jen.Err())),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.Nil()),
		),
		jen.Line(),
	}
}

func determineCreationColumns(dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	puvn := typ.Name.PluralUnexportedVarName()
	var creationColumns []jen.Code

	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, jen.Lit(field.Name.RouteName()))
	}

	if typ.BelongsToUser {
		creationColumns = append(creationColumns, jen.IDf("%sTableOwnershipColumn", puvn))
	}
	if typ.BelongsToStruct != nil {
		creationColumns = append(creationColumns, jen.IDf("%sTableOwnershipColumn", puvn))
	}

	if isMariaDB(dbvendor) {
		creationColumns = append(creationColumns, jen.Lit("created_on"))
	}

	return creationColumns
}

func determineValuesColumns(dbvendor wordsmith.SuperPalabra, inputVarName string, typ models.DataType) []jen.Code {
	var valuesColumns []jen.Code

	for _, field := range typ.Fields {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dot(field.Name.Singular()))
	}

	if typ.BelongsToUser {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		valuesColumns = append(valuesColumns, jen.ID(inputVarName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if isMariaDB(dbvendor) {
		valuesColumns = append(valuesColumns, jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery")))
	}

	return valuesColumns
}

func buildSomethingCreationTimeQueryFuncDecl(dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	uvn := typ.Name.UnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	return []jen.Code{
		jen.Commentf("build%sCreationTimeQuery takes %s and returns a creation query for that %s and the relevant arguments", sn, scnwp, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("build%sCreationTimeQuery", sn).Params(jen.IDf("%sID", uvn).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.created_on"), jen.IDf("%sTableName", puvn))).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Qual("fmt", "Sprintf").Call(jen.Lit("%s.id"), jen.IDf("%sTableName", puvn)).MapAssign().IDf("%sID", uvn))).
				Dotln("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
	}
}

func buildCreateSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	puvn := typ.Name.PluralUnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	qb := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Insert").Call(jen.IDf("%sTableName", puvn)).
		Dotln("Columns").Callln(determineCreationColumns(dbvendor, typ)...).
		Dotln("Values").Callln(determineValuesColumns(dbvendor, "input", typ)...)

	if isPostgres(dbvendor) {
		qb.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
	}
	qb.Dotln("ToSql").Call()

	params := typ.BuildCreateSomethingQueryParams(pkg, false)[1:]

	createQueryFuncBody := []jen.Code{
		jen.Var().ID("err").ID("error"),
		qb,
		jen.Line(),
		jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
		jen.Line(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	}

	return []jen.Code{
		jen.Commentf("buildCreate%sQuery takes %s and returns a creation query for that %s and the relevant arguments.", sn, scnwp, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildCreate%sQuery", sn).Params(
			params...,
		).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			createQueryFuncBody...,
		),
		jen.Line(),
	}
}

func buildCreateSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	dbrn := dbvendor.RouteName()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	var (
		createInitColumns []jen.Code
	)

	params := typ.BuildCreateSomethingParams(pkg, false)
	queryBuildingArgs := typ.BuildCreateSomethingArgs(pkg)[1:]
	queryBuildingArgs = queryBuildingArgs[:len(queryBuildingArgs)-1]
	queryBuildingArgs = append(queryBuildingArgs, jen.ID("x"))

	for _, field := range typ.Fields {
		fn := field.Name.Singular()
		createInitColumns = append(createInitColumns, jen.ID(fn).MapAssign().ID("input").Dot(fn))
	}

	if typ.BelongsToUser {
		createInitColumns = append(createInitColumns, jen.ID("BelongsToUser").MapAssign().ID("input").Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		createInitColumns = append(createInitColumns, jen.IDf("BelongsTo%s", typ.BelongsToStruct.Singular()).MapAssign().ID("input").Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	baseCreateFuncBody := []jen.Code{
		jen.ID("x").Assign().VarPointer().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(createInitColumns...),
		jen.Line(),
		jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildCreate%sQuery", sn).Call(
			queryBuildingArgs...,
		),
		jen.Line(),
		jen.Commentf("create the %s", scn),
	}

	if isPostgres(dbvendor) {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.Err().Assign().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.VarPointer().ID("x").Dot("ID"), jen.VarPointer().ID("x").Dot("CreatedOn")),
		)
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.List(jen.ID("res"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
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
			jen.Comment("fetch the last inserted ID"),
			jen.List(jen.ID("id"), jen.ID("idErr")).Assign().ID("res").Dot("LastInsertId").Call(),
			jen.If(jen.ID("idErr")).Op("==").ID("nil").Block(
				jen.ID("x").Dot("ID").Equals().ID("uint64").Call(jen.ID("id")),
				jen.Line(),
				jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("build%sCreationTimeQuery", sn).Call(jen.ID("x").Dot("ID")),
				jen.ID(dbfl).Dot("logCreationTimeRetrievalError").Call(
					jen.ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.VarPointer().ID("x").Dot("CreatedOn")),
				),
			),
		)
	}
	baseCreateFuncBody = append(baseCreateFuncBody, jen.Line(), jen.Return().List(jen.ID("x"), jen.Nil()))

	return []jen.Code{
		jen.Commentf("Create%s creates %s in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Create%s", sn).Params(
			params...,
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			baseCreateFuncBody...,
		),
		jen.Line(),
	}
}

func buildUpdateSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	puvn := typ.Name.PluralUnexportedVarName()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	const inputVarName = "updated"
	params := typ.BuildUpdateSomethingParams(pkg, inputVarName, false)[1:]

	vals := []jen.Code{
		jen.Lit("id").MapAssign().ID(inputVarName).Dot("ID"),
	}
	if typ.BelongsToUser {
		vals = append(vals, jen.IDf("%sTableOwnershipColumn", puvn).MapAssign().ID(inputVarName).Dot("BelongsToUser"))
	}
	if typ.BelongsToStruct != nil {
		vals = append(vals, jen.IDf("%sTableOwnershipColumn", puvn).MapAssign().ID(inputVarName).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	return []jen.Code{
		jen.Commentf("buildUpdate%sQuery takes %s and returns an update SQL query, with the relevant query parameters", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildUpdate%sQuery", sn).Params(
			params...,
		).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			func(typ models.DataType) jen.Code {
				x := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
					Dotln("Update").Call(jen.IDf("%sTableName", puvn))

				for _, field := range typ.Fields {
					if field.ValidForUpdateInput {
						x.Dotln("Set").Call(jen.Lit(field.Name.RouteName()), jen.ID(inputVarName).Dot(field.Name.Singular()))
					}
				}

				x.Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(vals...))

				if strings.ToLower(dbvsn) == "postgres" {
					x.Dotln("Suffix").Call(jen.Lit("RETURNING updated_on"))
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

func buildUpdateSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	n := typ.Name
	sn := n.Singular()
	scn := n.SingularCommonName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	const updatedVarName = "updated"

	var finalStatement jen.Code
	if isPostgres(dbvendor) {
		finalStatement = jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.VarPointer().ID(updatedVarName).Dot("UpdatedOn"))
	} else if isSqlite(dbvendor) || isMariaDB(dbvendor) {
		_g := &jen.Group{}
		_g.Add(
			jen.List(jen.ID("_"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Line(),
			jen.Return(jen.Err()),
		)
		finalStatement = _g
	}

	params := typ.BuildUpdateSomethingParams(pkg, updatedVarName, false)
	args := typ.BuildUpdateSomethingArgs(pkg, updatedVarName)[1:]

	return []jen.Code{
		jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the provided input to have a valid ID.", sn, scn, sn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Update%s", sn).Params(
			params...,
		).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildUpdate%sQuery", sn).Call(
				args...,
			),
			finalStatement,
		),
		jen.Line(),
	}
}

func buildArchiveSomethingQueryFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scn := typ.Name.SingularCommonName()
	uvn := typ.Name.UnexportedVarName()
	puvn := typ.Name.PluralUnexportedVarName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	var comment string
	vals := []jen.Code{
		jen.Lit("id").MapAssign().IDf("%sID", uvn),
		jen.Lit("archived_on").MapAssign().ID("nil"),
	}
	paramsList := typ.BuildGetSomethingParams(pkg)[1:]

	if typ.BelongsToUser {
		comment = fmt.Sprintf("buildArchive%sQuery returns a SQL query which marks a given %s belonging to a given user as archived.", sn, scn)
		vals = append(vals, jen.IDf("%sTableOwnershipColumn", puvn).MapAssign().ID("userID"))
	} else if typ.BelongsToStruct != nil {
		comment = fmt.Sprintf("buildArchive%sQuery returns a SQL query which marks a given %s belonging to a given %s as archived.", sn, scn, typ.BelongsToStruct.SingularCommonName())
		vals = append(vals, jen.IDf("%sTableOwnershipColumn", puvn).MapAssign().IDf("%sID", typ.BelongsToStruct.UnexportedVarName()))
	} else {
		comment = fmt.Sprintf("buildArchive%sQuery returns a SQL query which marks a given %s as archived.", sn, scn)
	}

	_qs := jen.List(jen.ID("query"), jen.ID("args"), jen.Err()).Equals().ID(dbfl).Dot("sqlBuilder").
		Dotln("Update").Call(jen.IDf("%sTableName", puvn)).
		Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
		Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("currentUnixTimeQuery"))).
		Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(vals...))

	if strings.ToLower(dbvsn) == "postgres" {
		_qs.Dotln("Suffix").Call(jen.Lit("RETURNING archived_on"))
	}
	_qs.Dotln("ToSql").Call()

	archiveFuncBody := []jen.Code{
		jen.Var().ID("err").ID("error"),
		_qs,
		jen.Line(),
		jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.Err()),
		jen.Line(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	}

	return []jen.Code{
		jen.Comment(comment),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildArchive%sQuery", sn).Params(
			jen.List(paramsList...)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			archiveFuncBody...,
		),
		jen.Line(),
	}
}

func buildArchiveSomethingFuncDecl(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) []jen.Code {
	dbvsn := dbvendor.Singular()
	sn := typ.Name.Singular()
	scnwp := typ.Name.SingularCommonNameWithPrefix()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))

	params := typ.BuildGetSomethingParams(pkg)
	queryBuildingArgs := typ.BuildGetSomethingArgs(pkg)[1:]

	return []jen.Code{
		jen.Commentf("Archive%s marks %s as archived in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Archive%s", sn).Params(
			params...,
		).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Assign().ID(dbfl).Dotf("buildArchive%sQuery", sn).Call(
				queryBuildingArgs...,
			),
			jen.List(jen.ID("_"), jen.Err()).Assign().ID(dbfl).Dot("db").Dot("ExecContext").Call(utils.CtxVar(), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	}
}
