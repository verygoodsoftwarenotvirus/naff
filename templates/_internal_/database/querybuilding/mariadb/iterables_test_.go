package mariadb

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

const timestampCall = "UNIX_TIMESTAMP()"

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
	sn := typ.Name.Singular()
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
		jen.Func().IDf("TestMariaDB_Build%sExistsQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.IDf("example%s", sn).Dot("BelongsToAccount"), jen.IDf("example%s", sn).Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("Build%sExistsQuery", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
	sn := typ.Name.Singular()

	whereValues := typ.BuildDBQuerierRetrievalQueryMethodQueryBuildingWhereClause(proj)
	cols := buildPrefixedStringColumns(typ)

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(cols...).
		From(tableName)
	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb).
		Where(whereValues)
	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildGet%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.IDf("example%s", sn).Dot("BelongsToAccount"), jen.IDf("example%s", sn).Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGet%sQuery", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
	pn := typ.Name.Plural()

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Select(fmt.Sprintf(countQuery, tableName)).
		From(tableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", tableName): nil,
		})

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildGetAll%sCountQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("actualQuery").Op(":=").ID("q").Dotf("BuildGetAll%sCountQuery", pn).Call(jen.ID("ctx")),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.Index().Interface().Values(),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
	pn := typ.Name.Plural()

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
		jen.Func().IDf("TestMariaDB_BuildGetBatchOf%sQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGetBatchOf%sQuery", pn).Call(
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
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
	pn := typ.Name.Plural()
	cols := buildPrefixedStringColumns(typ)

	query, _ := buildListQuery(tableName, cols, false)

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildGet%sQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("filter").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore"), jen.ID("exampleUser").Dot("ID"), jen.ID("exampleUser").Dot("ID"), jen.ID("filter").Dot("CreatedAfter"), jen.ID("filter").Dot("CreatedBefore"), jen.ID("filter").Dot("UpdatedAfter"), jen.ID("filter").Dot("UpdatedBefore")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGet%sQuery", pn).Call(
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
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
	pn := typ.Name.Plural()
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
		jen.Func().IDf("TestMariaDB_BuildGet%sWithIDsQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleUser").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeUser").Call(),
					jen.ID("exampleIDs").Op(":=").Index().ID("uint64").Valuesln(
						jen.Lit(789), jen.Lit(123), jen.Lit(456)),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleUser").Dot("ID"), jen.ID("exampleIDs").Index(jen.Lit(0)), jen.ID("exampleIDs").Index(jen.Lit(1)), jen.ID("exampleIDs").Index(jen.Lit(2))),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGet%sWithIDsQuery", pn).Call(
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
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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

func buildCreationStringColumnsAndArgs(typ models.DataType) (cols []string, args []jen.Code) {
	cols, args = []string{}, []jen.Code{}

	for _, field := range typ.Fields {
		if field.ValidForCreationInput {
			cols = append(cols, field.Name.RouteName())
			args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToStruct != nil {
		cols = append(cols, fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName()))
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	if typ.BelongsToAccount {
		cols = append(cols, "belongs_to_account")
		args = append(args, jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot(constants.AccountOwnershipFieldName))
	}

	return
}

func buildTestMariaDB_BuildCreateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	fieldCols, expectedArgs := buildCreationStringColumnsAndArgs(typ)
	valueArgs := []interface{}{}
	for range expectedArgs {
		valueArgs = append(valueArgs, whateverValue)
	}

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Insert(tableName).
		Columns(fieldCols...).
		Values(valueArgs...)

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildCreate%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
					jen.Newline(),
					jen.ID("exIDGen").Op(":=").Op("&").Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
					jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.IDf("example%s", sn).Dot("ExternalID")),
					jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.IDf("example%s", sn).Dot("ExternalID"), jen.IDf("example%s", sn).Dot("Name"), jen.IDf("example%s", sn).Dot("Details"), jen.IDf("example%s", sn).Dot("BelongsToAccount")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildCreate%sQuery", sn).Call(
						jen.ID("ctx"),
						jen.ID("exampleInput"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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

func buildUpdateQueryParts(typ models.DataType) []string {
	var out []string

	for _, field := range typ.Fields {
		out = append(out, fmt.Sprintf("%s = %s", field.Name.RouteName(), "?"))
	}

	return out
}

func buildTestMariaDB_BuildUpdateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	updateCols := buildUpdateQueryParts(typ)
	expectedArgs := []jen.Code{}

	valueArgs := []interface{}{}
	for range updateCols {
		valueArgs = append(valueArgs, whateverValue)
	}

	eq := squirrel.Eq{"id": whateverValue}
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).Update(tableName)
	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			qb = qb.Set(field.Name.RouteName(), jen.ID("input").Dot(field.Name.Singular()))
			expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToStruct != nil {
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount {
		eq["belongs_to_account"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.AccountOwnershipFieldName))
	}
	expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	qb = qb.Set("last_updated_on", squirrel.Expr(timestampCall)).Where(eq)

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildUpdate%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.IDf("example%s", sn).Dot("Name"), jen.IDf("example%s", sn).Dot("Details"), jen.IDf("example%s", sn).Dot("BelongsToAccount"), jen.IDf("example%s", sn).Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildUpdate%sQuery", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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

func buildTestBuildArchiveSomethingQueryBuilder(typ models.DataType) (qb squirrel.UpdateBuilder, expectedArgs []jen.Code, callArgs []jen.Code) {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	updateCols := buildUpdateQueryParts(typ)
	valueArgs := []interface{}{}
	for range updateCols {
		valueArgs = append(valueArgs, whateverValue)
	}

	eq := squirrel.Eq{
		"id":          whateverValue,
		"archived_on": nil,
	}
	if typ.BelongsToStruct != nil {
		btssn := typ.BelongsToStruct.Singular()
		eq[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"))
		callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(btssn)).Dot("ID"))
	}
	if typ.BelongsToAccount {
		eq["belongs_to_account"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName("User")).Dot("ID"))
	}
	callArgs = append(callArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).
		Update(tableName).
		Set("last_updated_on", squirrel.Expr(timestampCall)).
		Set("archived_on", squirrel.Expr(timestampCall)).
		Where(eq)

	return qb, expectedArgs, callArgs
}

func buildTestMariaDB_BuildArchiveSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	qb, _, _ := buildTestBuildArchiveSomethingQueryBuilder(typ)
	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildArchive%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.ID("exampleAccount").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeAccount").Call(),
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.IDf("example%s", sn).Dot("BelongsToAccount").Op("=").ID("exampleAccount").Dot("ID"),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Lit(query),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
						jen.ID("exampleAccount").Dot("ID"), jen.IDf("example%s", sn).Dot("ID")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildArchive%sQuery", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
						jen.ID("exampleAccount").Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
	sn := typ.Name.Singular()

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).Select(
		"audit_log.id",
		"audit_log.external_id",
		"audit_log.event_type",
		"audit_log.context",
		"audit_log.created_on",
	).
		From("audit_log").
		Where(squirrel.Expr(`JSON_CONTAINS(audit_log.context, '%d', '$.` + fmt.Sprintf(`%s_id')`, typ.Name.RouteName()))).
		OrderBy(fmt.Sprintf("%s.%s", "audit_log", "created_on"))

	query, _, _ := qb.ToSql()

	return []jen.Code{
		jen.Func().IDf("TestMariaDB_BuildGetAuditLogEntriesFor%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
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
					jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
					jen.Newline(),
					jen.ID("expectedQuery").Op(":=").Qual("fmt", "Sprintf").Call(
						jen.Lit(query),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.ID("expectedArgs").Op(":=").Index().Interface().Call(jen.ID("nil")),
					jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGetAuditLogEntriesFor%sQuery", sn).Call(
						jen.ID("ctx"),
						jen.IDf("example%s", sn).Dot("ID"),
					),
					jen.Newline(),
					jen.ID("assertArgCountMatchesQuery").Call(
						jen.ID("t"),
						jen.ID("actualQuery"),
						jen.ID("actualArgs"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
						jen.ID("t"),
						jen.ID("expectedQuery"),
						jen.ID("actualQuery"),
					),
					jen.Qual(constants.AssertionLibrary, "Equal").Call(
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
