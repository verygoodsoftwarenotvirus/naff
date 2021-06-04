package sqlite

import (
	"fmt"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/constants"
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

// applyFilterToSubCountQueryBuilder applies the query filter to a query builder.
func applyFilterToSubCountQueryBuilder(tableName string, queryBuilder squirrel.SelectBuilder) squirrel.SelectBuilder {
	return queryBuilder
}

func buildTotalCountQuery(tableName string, forAdmin bool) (query string, args []interface{}) {
	where := squirrel.Eq{}

	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	totalCountQueryBuilder := sqlBuilder.
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	if !forAdmin {
		where[fmt.Sprintf("%s.%s", tableName, "belongs_to_account")] = whateverValue
		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil
	}

	if len(where) > 0 {
		totalCountQueryBuilder = totalCountQueryBuilder.Where(where)
	}

	var err error

	query, args, err = totalCountQueryBuilder.ToSql()
	if err != nil {
		panic(err)
	}

	return query, args
}

func buildFilteredCountQuery(tableName string, forAdmin bool) (query string, args []interface{}) {
	where := squirrel.Eq{}

	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	filteredCountQueryBuilder := sqlBuilder.
		Select(fmt.Sprintf(columnCountQueryTemplate, tableName)).
		From(tableName)

	if !forAdmin {
		where[fmt.Sprintf("%s.belongs_to_account", tableName)] = whateverValue
		where[fmt.Sprintf("%s.%s", tableName, "archived_on")] = nil
	}

	if len(where) > 0 {
		filteredCountQueryBuilder = filteredCountQueryBuilder.Where(where)
	}

	filteredCountQueryBuilder = applyFleshedOutQueryFilterWithCode(filteredCountQueryBuilder, tableName, where)

	var err error

	query, args, err = filteredCountQueryBuilder.ToSql()
	if err != nil {
		panic(err)
	}

	return query, args
}

func applyFleshedOutQueryFilterWithCode(qb squirrel.SelectBuilder, tableName string, eq squirrel.Eq) squirrel.SelectBuilder {
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
		)}).
		OrderBy(fmt.Sprintf("%s.id", tableName)).
		Limit(20).
		Offset(180)

	return qb
}

func buildListQuery(tableName string, columns []string, forAdmin bool) (query string, args []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	filteredCountQuery, filteredCountQueryArgs := buildFilteredCountQuery(tableName, forAdmin)
	totalCountQuery, totalCountQueryArgs := buildTotalCountQuery(tableName, forAdmin)

	builder := sqlBuilder.
		Select(append(
			columns,
			fmt.Sprintf("(%s) as total_count", totalCountQuery),
			fmt.Sprintf("(%s) as filtered_count", filteredCountQuery),
		)...).
		From(tableName)

	if !forAdmin {
		w := squirrel.Eq{
			fmt.Sprintf("%s.%s", tableName, "archived_on"):  nil,
			fmt.Sprintf("%s.belongs_to_account", tableName): whateverValue,
		}

		builder = builder.Where(w)
	}

	builder = builder.GroupBy(fmt.Sprintf("%s.%s", tableName, "id"))

	builder = applyFleshedOutQueryFilterWithCode(builder, tableName, nil)

	query, selectArgs, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	return query, append(append(filteredCountQueryArgs, totalCountQueryArgs...), selectArgs...)
}
