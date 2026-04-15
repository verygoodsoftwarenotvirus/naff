package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	recipeListItemsTableName = "recipe_list_items"

	recipeListIDColumn = "recipe_list_id"
)

func init() {
	registerTableName(recipeListItemsTableName)
}

var recipeListItemsColumns = []string{
	idColumn,
	recipeIDColumn,
	notesColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
	"belongs_to_recipe_list",
}

func buildRecipeListItemsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(recipeListItemsColumns)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveRecipeListItem",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s) AND %s = sqlc.arg(%s);`,
					recipeListItemsTableName,
					archivedAtColumn, currentTimeExpression,
					archivedAtColumn,
					"belongs_to_recipe_list", "belongs_to_recipe_list",
					idColumn, idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateRecipeListItem",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					recipeListItemsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetRecipeListItems",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s.%s IS NULL
	%s
	AND %s.belongs_to_recipe_list = sqlc.arg(%s)
%s;`,
					strings.Join(applyToEach(recipeListItemsColumns, func(i int, s string) string {
						return fmt.Sprintf("%s.%s", recipeListItemsTableName, s)
					}), ",\n\t"),
					buildFilterCountSelect(recipeListItemsTableName, true, true, []string{}, fmt.Sprintf("%s.belongs_to_recipe_list = sqlc.arg(%s)", recipeListItemsTableName, recipeListIDColumn)),
					buildTotalCountSelect(recipeListItemsTableName, true, []string{}),
					recipeListItemsTableName,
					recipeListItemsTableName, archivedAtColumn,
					buildFilterConditions(recipeListItemsTableName, true, false, fmt.Sprintf("%s.belongs_to_recipe_list = sqlc.arg(%s)", recipeListItemsTableName, recipeListIDColumn)),
					recipeListItemsTableName, recipeListIDColumn,
					buildCursorLimitClause(recipeListItemsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateRecipeListItem",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					recipeListItemsTableName,
					strings.Join(applyToEach(filterForUpdate(recipeListItemsColumns, "belongs_to_recipe_list"), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn, currentTimeExpression,
					archivedAtColumn,
					"belongs_to_recipe_list", "belongs_to_recipe_list",
					idColumn, idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
