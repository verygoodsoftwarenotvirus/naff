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

func buildTableColumns(typ models.DataType) []jen.Code {
	tableColumns := []jen.Code{jen.Lit("id")}

	for _, field := range typ.Fields {
		tableColumns = append(tableColumns, jen.Lit(field.Name.RouteName()))
	}

	tableColumns = append(tableColumns,
		jen.Lit("created_on"),
		jen.Lit("updated_on"),
		jen.Lit("archived_on"),
		jen.Lit("belongs_to"),
	)

	return tableColumns
}

func iterablesDotGo(pkg *models.Project, dbvendor wordsmith.SuperPalabra, typ models.DataType) *jen.File {
	ret := jen.NewFile(dbvendor.SingularPackageName())

	utils.AddImports(pkg.OutputPath, []models.DataType{typ}, ret)
	dbvsn := dbvendor.Singular()

	n := typ.Name
	sn := n.Singular()
	pn := n.Plural()
	scn := n.SingularCommonName()
	dbrn := dbvendor.RouteName()
	dbfl := strings.ToLower(string([]byte(dbvsn)[0]))
	pcn := n.PluralCommonName()
	uvn := n.UnexportedVarName()
	puvn := n.PluralUnexportedVarName()
	scnwp := n.SingularCommonNameWithPrefix()
	pscnwp := n.ProperSingularCommonNameWithPrefix()
	prn := n.PluralRouteName()

	isPostgres := dbrn == "postgres"
	isSqlite := dbrn == "sqlite"
	isMariaDB := dbrn == "mariadb" || dbrn == "maria_db"

	ret.Add(
		jen.Const().Defs(
			jen.IDf("%sTableName", puvn).Op("=").Lit(prn),
		),
		jen.Line(),
	)

	tableColumns := buildTableColumns(typ)

	ret.Add(
		jen.Var().Defs(
			jen.IDf("%sTableColumns", puvn).Op("=").Index().ID("string").Valuesln(tableColumns...),
		),
		jen.Line(),
	)

	buildScanFields := func(typ models.DataType) (scanFields, floatSets []jen.Code) {
		scanFields = []jen.Code{jen.Op("&").ID("x").Dot("ID")}

		for _, field := range typ.Fields {
			scanFields = append(scanFields, jen.Op("&").ID("x").Dot(field.Name.Singular()))
		}

		scanFields = append(scanFields,
			jen.Op("&").ID("x").Dot("CreatedOn"),
			jen.Op("&").ID("x").Dot("UpdatedOn"),
			jen.Op("&").ID("x").Dot("ArchivedOn"),
			jen.Op("&").ID("x").Dot("BelongsTo"),
		)

		return scanFields, floatSets
	}
	scanFields, floatSets := buildScanFields(typ)

	buildScanItemBody := func() []jen.Code {
		body := []jen.Code{
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Values(),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("scan").Dot("Scan").Callln(scanFields...), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
		}

		body = append(body, floatSets...)

		body = append(body,
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		)

		return body
	}

	ret.Add(
		jen.Commentf("scan%s takes a database Scanner (i.e. *sql.Row) and scans the result into %s struct", sn, pscnwp),
		jen.Line(),
		jen.Func().IDf("scan%s", sn).Params(jen.ID("scan").Qual(filepath.Join(pkg.OutputPath, "database/v1"), "Scanner")).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			buildScanItemBody()...,
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("scan%s takes a logger and some database rows and turns them into a slice of %s", pn, pcn),
		jen.Line(),
		jen.Func().IDf("scan%s", pn).Params(jen.ID("logger").Qual("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "Logger"), jen.ID("rows").Op("*").Qual("database/sql", "Rows")).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.Var().ID("list").Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn),
			jen.Line(),
			jen.For(jen.ID("rows").Dot("Next").Call()).Block(
				jen.List(jen.ID("x"), jen.ID("err")).Op(":=").IDf("scan%s", sn).Call(jen.ID("rows")),
				jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
					jen.Return().List(jen.ID("nil"), jen.ID("err")),
				),
				jen.ID("list").Op("=").ID("append").Call(jen.ID("list"), jen.Op("*").ID("x")),
			),
			jen.Line(),
			jen.If(jen.ID("err").Op(":=").ID("rows").Dot("Err").Call(), jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("err")),
			),
			jen.Line(),
			jen.If(jen.ID("closeErr").Op(":=").ID("rows").Dot("Close").Call(), jen.ID("closeErr").Op("!=").ID("nil")).Block(
				jen.ID("logger").Dot("Error").Call(jen.ID("closeErr"), jen.Lit("closing database rows")),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("buildGet%sQuery constructs a SQL query for fetching %s with a given ID belong to a user with a given ID.", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGet%sQuery", sn).Params(jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Op("...")).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("id").Op(":").IDf("%sID", uvn),
				jen.Lit("belongs_to").Op(":").ID("userID"),
			)).Dot("ToSql").Call(),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%s fetches %s from the %s database", sn, scnwp, dbvendor.SingularPackageName()),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Get%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64"),
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", sn).Call(jen.IDf("%sID", uvn), jen.ID("userID")),
			jen.ID("row").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().IDf("scan%s", sn).Call(jen.ID("row")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("buildGet%sCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for", sn),
		jen.Line(),
		jen.Commentf("fetching the number of %s belonging to a given user that meet a given query", pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGet%sCountQuery", sn).Params(jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.ID("CountQuery")).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("archived_on").Op(":").ID("nil"),
				jen.Lit("belongs_to").Op(":").ID("userID"),
			)),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%sCount will fetch the count of %s from the database that meet a particular filter and belong to a particular user.", sn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Get%sCount", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64"),
		).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sCountQuery", sn).Call(jen.ID("filter"), jen.ID("userID")),
			jen.ID("err").Op("=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Var().Defs(
			jen.IDf("all%sCountQueryBuilder", pn).Qual("sync", "Once"),
			jen.IDf("all%sCountQuery", pn).ID("string"),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("buildGetAll%sCountQuery returns a query that fetches the total number of %s in the database.", pn, pcn),
		jen.Line(),
		jen.Comment("This query only gets generated once, and is otherwise returned from cache."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGetAll%sCountQuery", pn).Params().Params(jen.ID("string")).Block(
			jen.IDf("all%sCountQueryBuilder", pn).Dot("Do").Call(jen.Func().Params().Block(
				jen.Var().ID("err").ID("error"),
				jen.List(jen.IDf("all%sCountQuery", pn), jen.ID("_"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.ID("CountQuery")).
					Dotln("From").Call(jen.IDf("%sTableName", puvn)).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("archived_on").Op(":").ID("nil"))).
					Dotln("ToSql").Call(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			)),
			jen.Line(),
			jen.Return().IDf("all%sCountQuery", pn),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("GetAll%sCount will fetch the count of %s from the database", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("GetAll%sCount", pn).Params(jen.ID("ctx").Qual("context", "Context")).Params(jen.ID("count").ID("uint64"), jen.ID("err").ID("error")).Block(
			jen.ID("err").Op("=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID(dbfl).Dotf("buildGetAll%sCountQuery", pn).Call()).Dot("Scan").Call(jen.Op("&").ID("count")),
			jen.Return().List(jen.ID("count"), jen.ID("err")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("buildGet%sQuery builds a SQL query selecting %s that adhere to a given QueryFilter and belong to a given user,", pn, pcn),
		jen.Line(),
		jen.Commentf("and returns both the query and the relevant args to pass to the query executor."),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildGet%sQuery", pn).Params(jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(

			jen.Var().ID("err").ID("error"),
			jen.ID("builder").Op(":=").ID(dbfl).Dot("sqlBuilder").
				Dotln("Select").Call(jen.IDf("%sTableColumns", puvn).Op("...")).
				Dotln("From").Call(jen.IDf("%sTableName", puvn)).
				Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
				jen.Lit("archived_on").Op(":").ID("nil"),
				jen.Lit("belongs_to").Op(":").ID("userID"),
			)),
			jen.Line(),
			jen.If(jen.ID("filter").Op("!=").ID("nil")).Block(
				jen.ID("builder").Op("=").ID("filter").Dot("ApplyToQueryBuilder").Call(jen.ID("builder")),
			),
			jen.Line(),
			jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID("builder").Dot("ToSql").Call(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("Get%s fetches a list of %s from the database that meet a particular filter", pn, pcn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Get%s", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("filter").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "QueryFilter"),
			jen.ID("userID").ID("uint64"),
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", pn).Call(jen.ID("filter"), jen.ID("userID")),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Litf("querying database for %s", pcn))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").IDf("scan%s", pn).Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("scanning response from database: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.List(jen.ID("count"), jen.ID("err")).Op(":=").ID(dbfl).Dotf("Get%sCount", sn).Call(jen.ID("ctx"), jen.ID("filter"), jen.ID("userID")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("fetching %s count: ", scn)+"%w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sList", sn)).Valuesln(
				jen.ID("Pagination").Op(":").Qual(filepath.Join(pkg.OutputPath, "models/v1"), "Pagination").Valuesln(
					jen.ID("Page").Op(":").ID("filter").Dot("Page"),
					jen.ID("Limit").Op(":").ID("filter").Dot("Limit"),
					jen.ID("TotalCount").Op(":").ID("count"),
				),
				jen.IDf(pn).Op(":").ID("list"),
			),
			jen.Line(),
			jen.Return().List(jen.ID("x"), jen.ID("nil")),
		),
		jen.Line(),
	)

	ret.Add(
		jen.Commentf("GetAll%sForUser fetches every %s belonging to a user", pn, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("GetAll%sForUser", pn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("userID").ID("uint64"),
		).Params(jen.Index().Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildGet%sQuery", pn).Call(jen.ID("nil"), jen.ID("userID")),
			jen.Line(),
			jen.List(jen.ID("rows"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("QueryContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.ID("buildError").Call(jen.ID("err"), jen.Litf("fetching %s for user", pcn))),
			),
			jen.Line(),
			jen.List(jen.ID("list"), jen.ID("err")).Op(":=").IDf("scan%s", pn).Call(jen.ID(dbfl).Dot("logger"), jen.ID("rows")),
			jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
				jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit("parsing database results: %w"), jen.ID("err"))),
			),
			jen.Line(),
			jen.Return().List(jen.ID("list"), jen.ID("nil")),
		),
		jen.Line(),
	)

	////////////

	var (
		creationColumns []jen.Code
		valuesColumns   []jen.Code
	)

	for _, field := range typ.Fields {
		creationColumns = append(creationColumns, jen.Lit(field.Name.RouteName()))
		valuesColumns = append(valuesColumns, jen.ID("input").Dot(field.Name.Singular()))
	}

	creationColumns = append(creationColumns, jen.Lit("belongs_to"))
	valuesColumns = append(valuesColumns, jen.ID("input").Dot("BelongsTo"))
	if isMariaDB {
		creationColumns = append(creationColumns, jen.Lit("created_on"))
		valuesColumns = append(valuesColumns, jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery")))
	}

	qb := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
		Dotln("Insert").Call(jen.IDf("%sTableName", puvn)).
		Dotln("Columns").Callln(creationColumns...).
		Dotln("Values").Callln(valuesColumns...)

	if isPostgres {
		qb.Dotln("Suffix").Call(jen.Lit("RETURNING id, created_on"))
	}

	qb.Dotln("ToSql").Call()

	createQueryFuncBody := []jen.Code{
		jen.Var().ID("err").ID("error"),
	}

	createQueryFuncBody = append(createQueryFuncBody,
		qb,
		jen.Line(),
		jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
		jen.Line(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	)

	ret.Add(
		jen.Commentf("buildCreate%sQuery takes %s and returns a creation query for that %s and the relevant arguments.", sn, scnwp, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildCreate%sQuery", sn).Params(jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			createQueryFuncBody...,
		),
		jen.Line(),
	)

	////////////

	if isSqlite || isMariaDB {
		ret.Add(
			jen.Commentf("build%sCreationTimeQuery takes %s and returns a creation query for that %s and the relevant arguments", sn, scnwp, scn),
			jen.Line(),
			jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("build%sCreationTimeQuery", sn).Params(jen.IDf("%sID", uvn).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
				jen.Var().ID("err").ID("error"),
				jen.Line(),
				jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
					Dotln("Select").Call(jen.Lit("created_on")).
					Dotln("From").Call(jen.IDf("%sTableName", puvn)).
					Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Values(jen.Lit("id").Op(":").IDf("%sID", uvn))).
					Dotln("ToSql").Call(),
				jen.Line(),
				jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
				jen.Line(),
				jen.Return().List(jen.ID("query"), jen.ID("args")),
			),
		)
	}

	////////////

	var (
		createInitColumns []jen.Code
	)

	for _, field := range typ.Fields {
		fn := field.Name.Singular()
		createInitColumns = append(createInitColumns, jen.ID(fn).Op(":").ID("input").Dot(fn))
	}
	createInitColumns = append(createInitColumns, jen.ID("BelongsTo").Op(":").ID("input").Dot("BelongsTo"))

	baseCreateFuncBody := []jen.Code{
		jen.ID("x").Op(":=").Op("&").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn).Valuesln(createInitColumns...),
		jen.Line(),
		jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildCreate%sQuery", sn).Call(jen.ID("x")),
		jen.Line(),
		jen.Commentf("create the %s", scn),
	}

	if isPostgres {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.ID("err").Op(":=").ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("ID"), jen.Op("&").ID("x").Dot("CreatedOn")),
		)
	} else if isSqlite || isMariaDB {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.List(jen.ID("res"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
		)
	} else {
		println("dbrn is weeird", dbrn)
	}

	baseCreateFuncBody = append(baseCreateFuncBody,
		jen.If(jen.ID("err").Op("!=").ID("nil")).Block(
			jen.Return().List(jen.ID("nil"), jen.Qual("fmt", "Errorf").Call(jen.Lit(fmt.Sprintf("error executing %s creation query: ", scn)+"%w"), jen.ID("err"))),
		),
		jen.Line(),
	)

	if isSqlite || isMariaDB {
		baseCreateFuncBody = append(baseCreateFuncBody,
			jen.Comment("fetch the last inserted ID"),
			jen.List(jen.ID("id"), jen.ID("idErr")).Op(":=").ID("res").Dot("LastInsertId").Call(),
			jen.If(jen.ID("idErr")).Op("==").ID("nil").Block(
				jen.ID("x").Dot("ID").Op("=").ID("uint64").Call(jen.ID("id")),
				jen.Line(),
				jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("build%sCreationTimeQuery", sn).Call(jen.ID("x").Dot("ID")),
				jen.ID(dbfl).Dot("logCreationTimeRetrievalError").Call(
					jen.ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("x").Dot("CreatedOn")),
				),
			),
		)
	}
	baseCreateFuncBody = append(baseCreateFuncBody, jen.Line(), jen.Return().List(jen.ID("x"), jen.ID("nil")))

	ret.Add(
		jen.Commentf("Create%s creates %s in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Create%s", sn).Params(
			jen.ID("ctx").Qual("context", "Context"),
			jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), fmt.Sprintf("%sCreationInput", sn)),
		).Params(jen.Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn), jen.ID("error")).Block(
			baseCreateFuncBody...,
		),
		jen.Line(),
	)

	////////////

	buildQueryFunc := func(typ models.DataType) jen.Code {
		x := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
			Dotln("Update").Call(jen.IDf("%sTableName", puvn))

		for _, field := range typ.Fields {
			if field.ValidForUpdateInput {
				x.Dotln("Set").Call(jen.Lit(field.Name.RouteName()), jen.ID("input").Dot(field.Name.Singular()))
			}
		}
		x.Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
			Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
			jen.Lit("id").Op(":").ID("input").Dot("ID"),
			jen.Lit("belongs_to").Op(":").ID("input").Dot("BelongsTo"),
		))

		if strings.ToLower(dbvsn) == "postgres" {
			x.Dotln("Suffix").Call(jen.Lit("RETURNING updated_on"))
		}

		x.Dotln("ToSql").Call()
		return x
	}

	////////////

	ret.Add(
		jen.Commentf("buildUpdate%sQuery takes %s and returns an update SQL query, with the relevant query parameters", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildUpdate%sQuery", sn).Params(jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			jen.Var().ID("err").ID("error"),
			buildQueryFunc(typ),
			jen.Line(),
			jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
			jen.Line(),
			jen.Return().List(jen.ID("query"), jen.ID("args")),
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		func() []jen.Code {

			var finalStatement jen.Code
			if isPostgres {
				finalStatement = jen.Return().ID(dbfl).Dot("db").Dot("QueryRowContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")).Dot("Scan").Call(jen.Op("&").ID("input").Dot("UpdatedOn"))
			} else if isSqlite || isMariaDB {
				_g := &jen.Group{}
				_g.Add(
					jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
					jen.Line(),
					jen.Return(jen.ID("err")),
				)
				finalStatement = _g
			}

			return []jen.Code{
				jen.Commentf("Update%s updates a particular %s. Note that Update%s expects the provided input to have a valid ID.", sn, scn, sn),
				jen.Line(),
				jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Update%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.ID("input").Op("*").Qual(filepath.Join(pkg.OutputPath, "models/v1"), sn)).Params(jen.ID("error")).Block(
					jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildUpdate%sQuery", sn).Call(jen.ID("input")),
					finalStatement,
				),
				jen.Line(),
			}
		}()...,
	)

	////////////

	_qs := jen.List(jen.ID("query"), jen.ID("args"), jen.ID("err")).Op("=").ID(dbfl).Dot("sqlBuilder").
		Dotln("Update").Call(jen.IDf("%sTableName", puvn)).
		Dotln("Set").Call(jen.Lit("updated_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
		Dotln("Set").Call(jen.Lit("archived_on"), jen.Qual("github.com/Masterminds/squirrel", "Expr").Call(jen.ID("CurrentUnixTimeQuery"))).
		Dotln("Where").Call(jen.Qual("github.com/Masterminds/squirrel", "Eq").Valuesln(
		jen.Lit("id").Op(":").IDf("%sID", uvn),
		jen.Lit("archived_on").Op(":").ID("nil"),
		jen.Lit("belongs_to").Op(":").ID("userID"),
	))

	if strings.ToLower(dbvsn) == "postgres" {
		_qs.Dotln("Suffix").Call(jen.Lit("RETURNING archived_on"))
	}
	_qs.Dotln("ToSql").Call()

	archiveFuncBody := []jen.Code{
		jen.Var().ID("err").ID("error"),
		_qs,
		jen.Line(),
		jen.ID(dbfl).Dot("logQueryBuildingError").Call(jen.ID("err")),
		jen.Line(),
		jen.Return().List(jen.ID("query"), jen.ID("args")),
	}

	ret.Add(
		jen.Commentf("buildArchive%sQuery returns a SQL query which marks a given %s belonging to a given user as archived.", sn, scn),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("buildArchive%sQuery", sn).Params(jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.ID("query").ID("string"), jen.ID("args").Index().Interface()).Block(
			archiveFuncBody...,
		),
		jen.Line(),
	)

	////////////

	ret.Add(
		jen.Commentf("Archive%s marks %s as archived in the database", sn, scnwp),
		jen.Line(),
		jen.Func().Params(jen.ID(dbfl).Op("*").ID(dbvsn)).IDf("Archive%s", sn).Params(jen.ID("ctx").Qual("context", "Context"), jen.List(jen.IDf("%sID", uvn), jen.ID("userID")).ID("uint64")).Params(jen.ID("error")).Block(
			jen.List(jen.ID("query"), jen.ID("args")).Op(":=").ID(dbfl).Dotf("buildArchive%sQuery", sn).Call(jen.IDf("%sID", uvn), jen.ID("userID")),
			jen.List(jen.ID("_"), jen.ID("err")).Op(":=").ID(dbfl).Dot("db").Dot("ExecContext").Call(jen.ID("ctx"), jen.ID("query"), jen.ID("args").Op("...")),
			jen.Return().ID("err"),
		),
		jen.Line(),
	)
	return ret
}
