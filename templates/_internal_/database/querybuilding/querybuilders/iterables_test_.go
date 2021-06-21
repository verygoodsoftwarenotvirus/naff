package querybuilders

import (
	"fmt"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

func convertArgsToCode(args []interface{}) (code []jen.Code) {
	for _, arg := range args {
		if c, ok := arg.(models.Coder); ok {
			code = append(code, c.Code())
		}
	}
	return
}

func unixTimeForDatabase(db wordsmith.SuperPalabra) string {
	if db == nil {
		return "(strftime('%s','now'))"
	}

	switch db.LowercaseAbbreviation() {
	case "m":
		return "UNIX_TIMESTAMP()"
	case "p":
		return "extract(epoch FROM NOW())"
	case "s":
		return "(strftime('%s','now'))"
	default:
		panic(fmt.Sprintf("invalid database type! %q", db.LowercaseAbbreviation()))
	}
}

func queryBuilderForDatabase(db wordsmith.SuperPalabra) squirrel.StatementBuilderType {
	if db == nil {
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}

	switch db.LowercaseAbbreviation() {
	case "p":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	default:
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
}

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

	code.Add(buildTestSqlite_BuildSomethingExistsQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildGetSomethingQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildGetAllSomethingsCountQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildGetBatchOfSomethingsQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildGetSomethingsQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildGetSomethingsWithIDsQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildCreateSomethingQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildUpdateSomethingQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildArchiveSomethingQuery(proj, typ)...)
	code.Add(buildTestSqlite_BuildGetAuditLogEntriesForSomethingQuery(proj, typ)...)

	return code
}

func buildPrerequisiteIDs(proj *models.Project, typ models.DataType, includeSelf bool) []jen.Code {
	lines := []jen.Code{}

	if typ.RestrictedToAccountMembers {
		lines = append(lines, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
	}

	for _, dep := range proj.FindOwnerTypeChain(typ) {
		lines = append(lines, jen.IDf("example%sID", dep.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call())
	}

	if includeSelf {
		lines = append(lines, jen.IDf("example%s", typ.Name.Singular()).Assign().Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", typ.Name.Singular())).Call())
	}

	return lines
}

func buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	owners := p.FindOwnerTypeChain(typ)
	sn := typ.Name.Singular()

	for _, pt := range owners {
		pts := pt.Name.Singular()
		params = append(params, jen.IDf("example%sID", pts))
	}
	params = append(params, jen.IDf("example%s", sn).Dot("ID"))

	if typ.RestrictedToAccountAtSomeLevel(p) {
		params = append(params, jen.ID("exampleAccountID"))
	}

	return params
}

func buildDBQuerierSingleInstanceQueryMethodWhereClauses(p *models.Project, typ models.DataType) squirrel.Eq {
	n := typ.Name
	sn := n.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.id", tableName): models.NewCodeWrapper(jen.IDf("example%s", sn).Dot("ID")),
	}

	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		pTableName := pt.Name.PluralRouteName()

		whereValues[fmt.Sprintf("%s.id", pTableName)] = models.NewCodeWrapper(jen.IDf("example%sID", pt.Name.Singular()))
		whereValues[fmt.Sprintf("%s.archived_on", pTableName)] = nil

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues[fmt.Sprintf("%s.belongs_to_account", pTableName)] = models.NewCodeWrapper(jen.ID("exampleAccountID"))
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s", pTableName, pt.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("example%sID", pt.BelongsToStruct.Singular()))
		}
	}

	whereValues[fmt.Sprintf("%s.archived_on", tableName)] = nil

	if typ.BelongsToStruct != nil {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers {
		whereValues[fmt.Sprintf("%s.belongs_to_account", tableName)] = models.NewCodeWrapper(jen.ID("exampleAccountID"))
	}

	return whereValues
}

func buildTestSqlite_BuildSomethingExistsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	whereValues := buildDBQuerierSingleInstanceQueryMethodWhereClauses(proj, typ)

	qb := queryBuilderForDatabase(nil).Select(fmt.Sprintf("%s.id", tableName)).
		Prefix(existencePrefix).
		From(tableName)

	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb)

	qb = qb.Suffix(existenceSuffix).
		Where(whereValues)

	query, queryArgs, _ := qb.ToSql()
	expectedArgs := convertArgsToCode(queryArgs)

	bodyLines := []jen.Code{jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildPrerequisiteIDs(proj, typ, true)...)

	callArgs := buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
			expectedArgs...,
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("Build%sExistsQuery", sn).Call(
			callArgs...,
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_Build%sExistsQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSqlite_BuildGetSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	sn := typ.Name.Singular()

	whereValues := buildDBQuerierSingleInstanceQueryMethodWhereClauses(proj, typ)
	cols := buildPrefixedStringColumns(typ)

	qb := queryBuilderForDatabase(nil).
		Select(cols...).
		From(tableName)
	qb = typ.ModifyQueryBuilderWithJoinClauses(proj, qb).
		Where(whereValues)

	query, queryArgs, _ := qb.ToSql()
	expectedArgs := convertArgsToCode(queryArgs)

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildPrerequisiteIDs(proj, typ, true)...)

	callArgs := buildSingleInstanceQueryTestCallArgsWithoutOwnerVar(proj, typ)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
			expectedArgs...,
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGet%sQuery", sn).Call(
			callArgs...,
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildGet%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSqlite_BuildGetAllSomethingsCountQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	pn := typ.Name.Plural()

	qb := queryBuilderForDatabase(nil).
		Select(fmt.Sprintf(countQuery, tableName)).
		From(tableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", tableName): nil,
		})

	query, _, _ := qb.ToSql()

	bodyLines := []jen.Code{
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
	}

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildGetAll%sCountQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSqlite_BuildGetBatchOfSomethingsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	pn := typ.Name.Plural()

	cols := buildPrefixedStringColumns(typ)

	qb := queryBuilderForDatabase(nil).
		Select(cols...).
		From(tableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", tableName, "id"): whateverValue,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", tableName, "id"): whateverValue,
		})

	query, _, _ := qb.ToSql()

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
		jen.List(jen.ID("beginID"), jen.ID("endID")).Op(":=").List(jen.ID("uint64").Call(jen.Lit(1)),
			jen.ID("uint64").Call(jen.Lit(1000)),
		),
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(jen.ID("beginID"),
			jen.ID("endID"),
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGetBatchOf%sQuery", pn).Call(
			jen.ID("ctx"),
			jen.ID("beginID"),
			jen.ID("endID"),
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	}

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildGetBatchOf%sQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildArgsForDBQuerierTestOfListRetrievalQueryBuilder(p *models.Project, typ models.DataType) []jen.Code {
	params := []jen.Code{constants.CtxVar()}

	lp := []jen.Code{}
	owners := p.FindOwnerTypeChain(typ)
	for _, pt := range owners {
		lp = append(lp, jen.IDf("example%sID", pt.Name.Singular()))
	}

	if typ.RestrictedToAccountAtSomeLevel(p) {
		lp = append(lp, jen.ID("exampleAccountID"))
	}
	lp = append(lp, jen.False(), jen.ID(constants.FilterVarName))

	params = append(params, lp...)

	return params
}

func buildJoinsListForRealListQuery(p *models.Project, typ models.DataType) []string {
	out := []string{}
	prn := typ.Name.PluralRouteName()

	if typ.BelongsToStruct != nil {
		btsrn := typ.BelongsToStruct.RouteName()
		btprn := typ.BelongsToStruct.PluralRouteName()

		out = append(out, fmt.Sprintf("%s ON %s.belongs_to_%s=%s.id", btprn, prn, btsrn, btprn))
	}

	owners := p.FindOwnerTypeChain(typ)
	for i := len(owners) - 1; i >= 0; i-- {
		pt := owners[i]

		if pt.BelongsToStruct != nil {
			btsrn := pt.BelongsToStruct.RouteName()
			btprn := pt.BelongsToStruct.PluralRouteName()

			out = append(out, fmt.Sprintf("%s ON %s.belongs_to_%s=%s.id", btprn, pt.Name.PluralRouteName(), btsrn, btprn))
		}
	}

	return out
}

func buildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(p *models.Project, typ models.DataType) squirrel.Eq {
	tableName := typ.Name.PluralRouteName()

	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.archived_on", tableName): nil,
	}
	for _, pt := range p.FindOwnerTypeChain(typ) {
		pTableName := pt.Name.PluralRouteName()

		whereValues[fmt.Sprintf("%s.id", pTableName)] = models.NewCodeWrapper(jen.IDf("example%sID", pt.Name.Singular()))
		whereValues[fmt.Sprintf("%s.archived_on", pTableName)] = nil

		if pt.BelongsToAccount && pt.RestrictedToAccountMembers {
			whereValues[fmt.Sprintf("%s.belongs_to_account", pTableName)] = models.NewCodeWrapper(jen.ID("exampleAccountID"))
		}

		if pt.BelongsToStruct != nil {
			whereValues[fmt.Sprintf("%s.belongs_to_%s", pTableName, pt.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("example%sID", pt.BelongsToStruct.Singular()))
		}
	}

	whereValues[fmt.Sprintf("%s.archived_on", tableName)] = nil

	if typ.BelongsToStruct != nil && !typ.IsEnumeration {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = models.NewCodeWrapper(jen.IDf("example%sID", typ.BelongsToStruct.Singular()))
	}
	if typ.BelongsToAccount && typ.RestrictedToAccountMembers && !typ.IsEnumeration {
		whereValues[fmt.Sprintf("%s.belongs_to_account", tableName)] = models.NewCodeWrapper(jen.ID("exampleAccountID"))
	}

	return whereValues
}

func buildTestSqlite_BuildGetSomethingsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	pn := typ.Name.Plural()

	cols := buildPrefixedStringColumns(typ)
	joins := buildJoinsListForRealListQuery(proj, typ)
	where := buildDBQuerierListRetrievalQueryMethodQueryBuildingWhereClause(proj, typ)

	query, queryArgs := buildListQuery(queryBuilderForDatabase(nil), tableName, joins, where, "", cols, 0, false)
	expectedArgs := convertArgsToCode(queryArgs)

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines, buildPrerequisiteIDs(proj, typ, false)...)
	callArgs := buildArgsForDBQuerierTestOfListRetrievalQueryBuilder(proj, typ)

	bodyLines = append(bodyLines,
		jen.ID("filter").Op(":=").Qual(proj.FakeTypesPackage(), "BuildFleshedOutQueryFilter").Call(),
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
			expectedArgs...,
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGet%sQuery", pn).Call(callArgs...),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildGet%sQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSqlite_BuildGetSomethingsWithIDsQuery(proj *models.Project, typ models.DataType) []jen.Code {
	tableName := typ.Name.PluralRouteName()
	pn := typ.Name.Plural()
	cols := buildPrefixedStringColumns(typ)

	var qb squirrel.SelectBuilder
	whereValues := squirrel.Eq{
		fmt.Sprintf("%s.%s", tableName, "id"):          []string{whateverValue, whateverValue, whateverValue},
		fmt.Sprintf("%s.%s", tableName, "archived_on"): nil,
	}
	if typ.BelongsToStruct != nil {
		whereValues[fmt.Sprintf("%s.belongs_to_%s", tableName, typ.BelongsToStruct.RouteName())] = whateverValue
	}

	if typ.BelongsToAccount {
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

	qb = queryBuilderForDatabase(nil).
		Select(cols...).
		From(tableName).
		Where(whereValues).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", tableName, "id", whenThenStatement)).
		Limit(20)

	expectedQuery, _, _ := qb.ToSql()

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
		utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()),
		func() jen.Code {
			if typ.BelongsToStruct != nil {
				return jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Op(":=").Qual(proj.FakeTypesPackage(), "BuildFakeID").Call()
			}
			return jen.Null()
		}(),
		jen.ID("exampleIDs").Op(":=").Index().ID("uint64").Valuesln(
			jen.Lit(789),
			jen.Lit(123),
			jen.Lit(456),
		),
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(expectedQuery),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			jen.ID("exampleIDs").Index(jen.Lit(0)),
			jen.ID("exampleIDs").Index(jen.Lit(1)),
			jen.ID("exampleIDs").Index(jen.Lit(2)),
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGet%sWithIDsQuery", pn).Call(
			jen.ID("ctx"),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.BelongsToAccount, jen.ID("exampleAccountID")),
			jen.ID("defaultLimit"),
			jen.ID("exampleIDs"),
			func() jen.Code {
				if typ.BelongsToAccount {
					return jen.True()
				}
				return jen.False()
			}(),
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	}

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildGet%sWithIDsQuery", pn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildCreationStringColumnsAndArgs(typ models.DataType) (cols []string, args []jen.Code) {
	cols, args = []string{"external_id"}, []jen.Code{jen.ID(utils.BuildFakeVarName(typ.Name.Singular())).Dot("ExternalID")}

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

func buildTestSqlite_BuildCreateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	fieldCols, fieldArgs := buildCreationStringColumnsAndArgs(typ)
	valueArgs := []interface{}{}
	for range fieldArgs {
		valueArgs = append(valueArgs, whateverValue)
	}

	qb := queryBuilderForDatabase(nil).
		Insert(tableName).
		Columns(fieldCols...).
		Values(valueArgs...)

	query, _, _ := qb.ToSql()

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines,
		jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
	)

	bodyLines = append(bodyLines,
		jen.ID("exampleInput").Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%sCreationInputFrom%s", sn, sn)).Call(jen.IDf("example%s", sn)),
		jen.Newline(),
		jen.ID("exIDGen").Op(":=").Op("&").Qual(proj.QuerybuildingPackage(), "MockExternalIDGenerator").Values(),
		jen.ID("exIDGen").Dot("On").Call(jen.Lit("NewExternalID")).Dot("Return").Call(jen.IDf("example%s", sn).Dot("ExternalID")),
		jen.ID("q").Dot("externalIDGenerator").Op("=").ID("exIDGen"),
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		func() jen.Code {
			if len(fieldArgs) == 0 {
				return jen.ID("expectedArgs").Op(":=").Index().Interface().Values(fieldArgs...)
			}
			return jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
				fieldArgs...,
			)
		}(),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildCreate%sQuery", sn).Call(
			jen.ID("ctx"),
			jen.ID("exampleInput"),
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
		jen.Newline(),
		jen.ID("mock").Dot("AssertExpectationsForObjects").Call(
			jen.ID("t"),
			jen.ID("exIDGen"),
		),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildCreate%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
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

func buildTestSqlite_BuildUpdateSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()
	tableName := typ.Name.PluralRouteName()

	updateCols := buildUpdateQueryParts(typ)
	expectedArgs := []jen.Code{}

	valueArgs := []interface{}{}
	for range updateCols {
		valueArgs = append(valueArgs, whateverValue)
	}

	where := squirrel.Eq{"id": whateverValue, "archived_on": nil}
	qb := queryBuilderForDatabase(nil).Update(tableName)

	for _, field := range typ.Fields {
		if field.ValidForUpdateInput {
			qb = qb.Set(field.Name.RouteName(), jen.ID("input").Dot(field.Name.Singular()))
			expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(field.Name.Singular()))
		}
	}

	if typ.BelongsToAccount {
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot(constants.AccountOwnershipFieldName))
		where["belongs_to_account"] = whateverValue
	}

	if typ.BelongsToStruct != nil {
		where[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dotf("BelongsTo%s", typ.BelongsToStruct.Singular()))
	}

	expectedArgs = append(expectedArgs, jen.ID(utils.BuildFakeVarName(sn)).Dot("ID"))

	qb = qb.Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(nil))).Where(where)

	query, _, _ := qb.ToSql()

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines,
		jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		func() jen.Code {
			if len(expectedArgs) == 0 {
				return jen.ID("expectedArgs").Op(":=").Index().Interface().Values(expectedArgs...)
			}
			return jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
				expectedArgs...,
			)
		}(),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildUpdate%sQuery", sn).Call(
			jen.ID("ctx"),
			jen.IDf("example%s", sn),
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildUpdate%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
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

	where := squirrel.Eq{
		"id":          whateverValue,
		"archived_on": nil,
	}
	if typ.BelongsToStruct != nil {
		btssn := typ.BelongsToStruct.Singular()
		where[fmt.Sprintf("belongs_to_%s", typ.BelongsToStruct.RouteName())] = whateverValue
		expectedArgs = append(expectedArgs, jen.IDf("example%sID", btssn))
		callArgs = append(callArgs, jen.IDf("example%sID", btssn))
	}
	if typ.RestrictedToAccountMembers {
		where["belongs_to_account"] = whateverValue
		expectedArgs = append(expectedArgs, jen.ID("exampleAccountID"))
	}
	callArgs = append(callArgs, jen.IDf("example%sID", sn))

	expectedArgs = append(expectedArgs, jen.IDf("example%sID", sn))

	qb = queryBuilderForDatabase(nil).
		Update(tableName).
		Set("last_updated_on", squirrel.Expr(unixTimeForDatabase(nil))).
		Set("archived_on", squirrel.Expr(unixTimeForDatabase(nil))).
		Where(where)

	return qb, expectedArgs, callArgs
}

func buildTestSqlite_BuildArchiveSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	qb, _, _ := buildTestBuildArchiveSomethingQueryBuilder(typ)
	query, _, _ := qb.ToSql()

	bodyLines := []jen.Code{
		jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	if typ.BelongsToStruct != nil {
		bodyLines = append(bodyLines,
			jen.IDf("example%sID", typ.BelongsToStruct.Singular()).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
		)
	}

	if typ.RestrictedToAccountMembers {
		bodyLines = append(bodyLines,
			jen.IDf("exampleAccountID").Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
		)
	}

	bodyLines = append(bodyLines,
		jen.IDf("example%sID", sn).Assign().Qual(proj.FakeTypesPackage(), "BuildFakeID").Call(),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.ID("exampleAccountID")),
			jen.IDf("example%sID", sn),
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildArchive%sQuery", sn).Call(
			jen.ID("ctx"),
			func() jen.Code {
				if typ.BelongsToStruct != nil {
					return jen.IDf("example%sID", typ.BelongsToStruct.Singular())
				}
				return jen.Null()
			}(),
			jen.IDf("example%sID", sn),
			utils.ConditionalCode(typ.RestrictedToAccountMembers, jen.ID("exampleAccountID")),
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildArchive%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}

func buildTestSqlite_BuildGetAuditLogEntriesForSomethingQuery(proj *models.Project, typ models.DataType) []jen.Code {
	sn := typ.Name.Singular()

	qb := queryBuilderForDatabase(nil).Select(
		"audit_log.id",
		"audit_log.external_id",
		"audit_log.event_type",
		"audit_log.context",
		"audit_log.created_on",
	).
		From("audit_log").
		Where(squirrel.Eq{fmt.Sprintf(`json_extract(audit_log.context, '$.%s_id')`, typ.Name.RouteName()): whateverValue}).
		OrderBy(fmt.Sprintf("%s.%s", "audit_log", "created_on"))

	query, _, _ := qb.ToSql()

	bodyLines := []jen.Code{jen.ID("t").Dot("Parallel").Call(),
		jen.Newline(),
		jen.List(jen.ID("q"), jen.ID("_")).Op(":=").ID("buildTestService").Call(jen.ID("t")),
		jen.ID("ctx").Op(":=").Qual("context", "Background").Call(),
		jen.Newline(),
	}

	bodyLines = append(bodyLines,
		jen.IDf("example%s", sn).Op(":=").Qual(proj.FakeTypesPackage(), fmt.Sprintf("BuildFake%s", sn)).Call(),
	)

	bodyLines = append(bodyLines,
		jen.Newline(),
		jen.ID("expectedQuery").Op(":=").Lit(query),
		jen.ID("expectedArgs").Op(":=").Index().Interface().Valuesln(
			jen.IDf("example%s", sn).Dot("ID"),
		),
		jen.List(jen.ID("actualQuery"), jen.ID("actualArgs")).Op(":=").ID("q").Dotf("BuildGetAuditLogEntriesFor%sQuery", sn).Call(
			jen.ID("ctx"),
			jen.IDf("example%s", sn).Dot("ID"),
		),
		jen.Newline(),
		jen.ID("assertArgCountMatchesQuery").Call(jen.ID("t"), jen.ID("actualQuery"), jen.ID("actualArgs")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedQuery"), jen.ID("actualQuery")),
		jen.Qual(constants.AssertionLibrary, "Equal").Call(jen.ID("t"), jen.ID("expectedArgs"), jen.ID("actualArgs")),
	)

	return []jen.Code{
		jen.Func().IDf("TestSqlite_BuildGetAuditLogEntriesFor%sQuery", sn).Params(jen.ID("T").Op("*").Qual("testing", "T")).Body(
			jen.ID("T").Dot("Parallel").Call(),
			jen.Newline(),
			jen.ID("T").Dot("Run").Call(jen.Lit("standard"),
				jen.Func().Params(jen.ID("t").Op("*").Qual("testing", "T")).Body(
					bodyLines...,
				),
			),
		),
		jen.Newline(),
	}
}
