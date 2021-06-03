package mariadb

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

func buildPrefixedStringColumns(typ models.DataType) []string {
	tableName := typ.Name.PluralRouteName()
	out := []string{
		fmt.Sprintf("%s.id", tableName),
		fmt.Sprintf("%s.external_id", tableName),
	}

	for _, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s.%s", tableName, field.Name.RouteName()))
	}

	out = append(out, fmt.Sprintf("%s.created_on", tableName), fmt.Sprintf("%s.last_updated_on", tableName), fmt.Sprintf("%s.archived_on", tableName))
	if typ.BelongsToAccount {
		out = append(out, fmt.Sprintf("%s.belongs_to_account", tableName))
	}
	if typ.BelongsToStruct != nil {
		out = append(out, fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName()))
	}

	return out
}

func iterablesTestDotGo(proj *models.Project, typ models.DataType) *jen.File {
	code := jen.NewFile(packageName)

	utils.AddImports(proj, code, false)

	code.Add(buildTestMariaDB_BuildSomethingExistsQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildGetSomethingQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildGetAllSomethingsCountQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildGetBatchOfSomethingsQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildGetSomethingsQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildGetSomethingsWithIDsQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildCreateSomethingQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildUpdateSomethingQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildArchiveSomethingQuery(proj, typ)...)
	code.Add(buildTestMariaDB_BuildGetAuditLogEntriesForSomethingQuery(proj, typ)...)

	return code
}

func buildTestMariaDB_BuildSomethingExistsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()

	whereValues := typ.BuildDBQuerierExistenceQueryMethodQueryBuildingWhereClause(proj)

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildItemExistsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.IDf("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.IDf("exampleItem").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.IDf("exampleItem").Dot("BelongsToAccount"), jen.IDf("exampleItem").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildItemExistsQuery").Call(
						jen.ID("ctx"),
						jen.IDf("exampleItem").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildGetSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	whereValues := typ.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(proj)
	cols := buildPrefixedStringColumns(typ)

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cols...).
		From(tableName)
	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb).
		Where(whereValues)
	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildGetItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleItem").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleItem").Dot("BelongsToAccount"), jen.ID("exampleItem").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetItemQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildGetAllSomethingsCountQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(countQuery, tableName)).
		From(tableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", tableName): nil,
		})

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildGetAllItemsCountQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("actualQuery").Op(":=").ID("q").Dot("BuildGetAllItemsCountQuery").Call(jen.ID("ctx")),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildGetBatchOfSomethingsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumns(typ)

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cols...).
		From(tableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", tableName, "id"): whateverValue,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", tableName, "id"): whateverValue,
		})

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildGetBatchOfItemsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.List(jen.ID("beginID"), jen.ID("endID")).Op(":=").List(jen.ID("uint64").Call(jen.Lit(1)), jen.ID("uint64").Call(jen.Lit(1000))),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("beginID"), jen.ID("endID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetBatchOfItemsQuery").Call(
						jen.ID("ctx"),
						jen.ID("beginID"),
						jen.ID("endID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildGetSomethingsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumns(typ)

	query, _ := buildListQuery(tableName, "belongs_to_account", cols, 0, false)

	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildGetItemsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("filter").Op(":=").ID("fakes").Dot("BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetItemsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("false"),
						jen.ID("filter"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildGetSomethingsWithIDsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	cols := buildPrefixedStringColumns(typ)

	var qb squirrel.SelectBuilder
	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.%s", tableName, "id"):          []string{whateverValue, whateverValue, whateverValue},
		fmt.Sprintf("%s.%s", tableName, "archived_on"): nil,
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		whereValues[fmt.Sprintf("%s.%s", tableName, "belongs_to_account")] = whateverValue
	}

	var whenThenStatement string
	for i, id := range []uint64{789, 123, 456} {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cols...).
		From(tableName).
		Where(whereValues).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", tableName, "id", whenThenStatement)).
		Limit(20)

	expectedQuery, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildGetItemsWithIDsQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleUser").Op(":=").ID("fakes").Dot("BuildFakeUser").Call(),
					jen.ID("exampleIDs").Op(":=").Index().ID("uint64").Valuesln(
						jen.Lit(789), jen.Lit(123), jen.Lit(456)),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleUser").Dot("ID"), jen.ID("exampleIDs").Index(jen.Lit(0)), jen.ID("exampleIDs").Index(jen.Lit(1)), jen.ID("exampleIDs").Index(jen.Lit(2))),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetItemsWithIDsQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleUser").Dot("ID"),
						jen.ID("defaultLimit"),
						jen.ID("exampleIDs"),
						jen.ID("false"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildCreateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildCreateItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleItem").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").ID("fakes").Dot("BuildFakeItemCreationInputFromItem").Call(jen.ID("exampleItem")),
					jen.Newline(),
					jen.ID("exIDGen").Op(":=").Op("&").ID("querybuilding").Dot("MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.ID("exampleItem").Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("INSERT INTO items (external_id,name,details,belongs_to_account) VALUES (?,?,?,?)"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleItem").Dot("ExternalID"), jen.ID("exampleItem").Dot("Name"), jen.ID("exampleItem").Dot("Details"), jen.ID("exampleItem").Dot("BelongsToAccount")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildCreateItemQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
					jen.Newline(),
					jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
						jen.ID("t"),
						jen.ID("exIDGen"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildUpdateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildUpdateItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleItem").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET name = ?, details = ?, last_updated_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_account = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleItem").Dot("Name"), jen.ID("exampleItem").Dot("Details"), jen.ID("exampleItem").Dot("BelongsToAccount"), jen.ID("exampleItem").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildUpdateItemQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildArchiveSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildArchiveItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleAccount").Op(":=").ID("fakes").Dot("BuildFakeAccount").Call(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.ID("exampleItem").Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit("UPDATE items SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_account = ? AND id = ?"),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccount").Dot("ID"), jen.ID("exampleItem").Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildArchiveItemQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem").Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestMariaDB_BuildGetAuditLogEntriesForSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	return []jen.Code{
		jen.Func().ID("TestMariaDB_BuildGetAuditLogEntriesForItemQuery").Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(
				jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					jen.ID("t").Dot("Parallel").Call(),
					jen.Newline(),
					jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
					jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
					jen.Newline(),
					jen.ID("exampleItem").Op(":=").ID("fakes").Dot("BuildFakeItem").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit("SELECT audit_log.id, audit_log.external_id, audit_log.event_type, audit_log.context, audit_log.created_on FROM audit_log WHERE JSON_CONTAINS(audit_log.context, '%d', '$.item_id') ORDER BY audit_log.created_on"),
						jen.ID("exampleItem").Dot("ID"),
					),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Call(jen.ID("nil")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dot("BuildGetAuditLogEntriesForItemQuery").Call(
						jen.ID("ctx"),
						jen.ID("exampleItem").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.ID("assert").Dot("Equal").Call(
						jen.ID("t"),
						jen.ID("expectedArgs"),
						jen.ID("actualArgs"),
					),
				),
			),
		),
		jen.Newline(),
	}
}
