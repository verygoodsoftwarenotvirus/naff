package querybuilders

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/Masterminds/squirrel"
)

const (
	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	whateverValue = "fart"

	// countQuery is a generic counter query used in a few query builders
	countQuery = "COUNT(%s.id)"
	// columnCountQueryTemplate is a generic counter query used in a few query builders.
	columnCountQueryTemplate = `COUNT(%s.id)`
	// allCountQuery is a generic counter query used in a few query builders.
	allCountQuery = `COUNT(*)`
)

func getIncIndex(dbv wordsmith.SuperPalabra, index uint) string {
	if dbv.LowercaseAbbreviation() == "p" {
		return fmt.Sprintf("$%d", index+1)
	}
	return "?"
}

// BuildQuery builds a given query, handles whatever errs and returns just the query and args.
func buildQuery(builder squirrel.Sqlizer) (query string, args []interface{}) {
	query, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	return query, args
}

func buildTotalCountQuery(sqlBuilder squirrel.StatementBuilderType, tableName string, joins []string, where squirrel.Eq, ownershipColumn string, userID uint64, forAdmin, includeArchived bool) (query string, args []interface{}) {
	if where == nil {
		where = squirrel.Eq{}
	}

	totalCountQueryBuilder := sqlBuilder.
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		totalCountQueryBuilder = totalCountQueryBuilder.Join(join)
	}

	if !forAdmin {
		if userID != 0 && ownershipColumn != "" {
			where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
		}

		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil
	} else if !includeArchived {
		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil
	}

	if len(where) > 0 {
		totalCountQueryBuilder = totalCountQueryBuilder.Where(where)
	}

	return buildQuery(totalCountQueryBuilder)
}

func buildFilteredCountQuery(sqlBuilder squirrel.StatementBuilderType, tableName string, joins []string, where squirrel.Eq, ownershipColumn string, userID uint64, forAdmin, includeArchived bool) (query string, args []interface{}) {
	if where == nil {
		where = squirrel.Eq{}
	}

	queryBuilder := sqlBuilder.
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	for _, join := range joins {
		queryBuilder = queryBuilder.Join(join)
	}

	if !forAdmin {
		if userID != 0 && ownershipColumn != "" {
			where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = userID
		}

		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil
	} else if !includeArchived {
		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil
	}

	if len(where) > 0 {
		queryBuilder = queryBuilder.Where(where)
	}

	queryBuilder = applyFleshedOutQueryFilterWithCode(queryBuilder, tableName, nil, false, false)

	return buildQuery(queryBuilder)
}

func applyFleshedOutQueryFilterWithCode(qb squirrel.SelectBuilder, tableName string, eq squirrel.Eq, includeOrdering, includeOffset bool) squirrel.SelectBuilder {
	if eq != nil {
		qb = qb.Where(eq)
	}

	qb = qb.
		Where(squirrel.Gt{fmt.Sprintf("%s.created_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("CreatedAfter"),
		)}).
		Where(squirrel.Lt{fmt.Sprintf("%s.created_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("CreatedBefore"),
		)}).
		Where(squirrel.Gt{fmt.Sprintf("%s.last_updated_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("UpdatedAfter"),
		)}).
		Where(squirrel.Lt{fmt.Sprintf("%s.last_updated_on", tableName): models.NewCodeWrapper(
			jen.ID(constants.FilterVarName).Dot("UpdatedBefore"),
		)})

	if includeOrdering {
		qb = qb.OrderBy(fmt.Sprintf("%s.id", tableName))
	}

	if includeOffset {
		qb = qb.
			Limit(20).
			Offset(180)
	}

	return qb
}

// BuildListQuery builds a SQL query selecting rows that adhere to a given QueryFilter and belong to a given account,
// and returns both the query and the relevant args to pass to the query executor.
func buildListQuery(sqlBuilder squirrel.StatementBuilderType, tableName string, joins []string, where squirrel.Eq, ownershipColumn string, columns []string, ownerID uint64, forAdmin bool) (query string, args []interface{}) {
	var includeArchived bool

	filteredCountQuery, filteredCountQueryArgs := buildFilteredCountQuery(sqlBuilder, tableName, joins, where, ownershipColumn, ownerID, forAdmin, includeArchived)
	totalCountQuery, totalCountQueryArgs := buildTotalCountQuery(sqlBuilder, tableName, joins, where, ownershipColumn, ownerID, forAdmin, includeArchived)

	builder := sqlBuilder.
		Select(append(
			columns,
			fmt.Sprintf("(%s) as total_count", totalCountQuery),
			fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
		)...).
		From(tableName)

	for _, join := range joins {
		builder = builder.Join(join)
	}

	if !forAdmin {
		if where == nil {
			where = squirrel.Eq{}
		}
		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil

		if ownershipColumn != "" && ownerID != 0 {
			where[fmt.Sprintf("%s.%s", tableName, ownershipColumn)] = ownerID
		}

		builder = builder.Where(where)
	}

	builder = builder.GroupBy(fmt.Sprintf("%s.%s", tableName, "id"))

	builder = applyFleshedOutQueryFilterWithCode(builder, tableName, nil, false, true)

	query, selectArgs := buildQuery(builder)

	return query, append(append(filteredCountQueryArgs, totalCountQueryArgs...), selectArgs...)
}
