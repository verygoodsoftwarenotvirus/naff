package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeListsTableName = "recipe_lists"
)

func init() {
	registerTableName(recipeListsTableName)
}

var recipeListsColumns = []string{
	idColumn,
	nameColumn,
	descriptionColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	belongsToUserColumn,
}

func buildRecipeListsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(recipeListsColumns)
		fullSelectColumns := mergeColumns(
			applyToEach(recipeListsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s", recipeListsTableName, s)
			}),
			applyToEach(recipeListItemsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as recipe_list_item_%s", recipeListItemsTableName, s, s)
			}),
			2,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveRecipeList",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
					recipeListsTableName,
					archivedAtColumn, currentTimeExpression,
					archivedAtColumn,
					belongsToUserColumn, belongsToUserColumn,
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateRecipeList",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					recipeListsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeLists",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	LEFT JOIN %s ON %s.%s = %s.%s AND %s.%s IS NULL
	WHERE %s.%s IS NULL
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(recipeListsTableName, true, true, []string{}),
					buildTotalCountSelect(recipeListsTableName, true, []string{}),
					recipeListsTableName,
					recipeListItemsTableName, recipeListItemsTableName, "belongs_to_recipe_list", recipeListsTableName, idColumn, recipeListItemsTableName, archivedAtColumn,
					recipeListsTableName, archivedAtColumn,
					buildFilterConditions(recipeListsTableName, true, false),
					buildCursorLimitClause(recipeListsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateRecipeList",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					recipeListsTableName,
					strings.Join(applyToEach(filterForUpdate(recipeListsColumns, belongsToUserColumn), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					belongsToUserColumn, belongsToUserColumn,
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
