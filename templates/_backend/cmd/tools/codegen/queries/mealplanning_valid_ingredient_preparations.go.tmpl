package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validIngredientPreparationsTableName = "valid_ingredient_preparations"
)

func init() {
	registerTableName(validIngredientPreparationsTableName)
}

var validIngredientPreparationsColumns = []string{
	idColumn,
	notesColumn,
	validPreparationIDColumn,
	validIngredientIDColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidIngredientPreparationsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(validIngredientPreparationsColumns)

		fullSelectColumns := mergeColumns(
			mergeColumns(
				applyToEach(filterFromSlice(validIngredientPreparationsColumns, "valid_preparation_id", "valid_ingredient_id"), func(i int, s string) string {
					return fmt.Sprintf("%s.%s as valid_ingredient_preparation_%s", validIngredientPreparationsTableName, s, s)
				}),
				applyToEach(validIngredientsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s as valid_ingredient_%s", validIngredientsTableName, s, s)
				}),
				2,
			),
			applyToEach(validPreparationsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_preparation_%s", validPreparationsTableName, s, s)
			}),
			2,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveValidIngredientPreparation",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateValidIngredientPreparation",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					validIngredientPreparationsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckValidIngredientPreparationExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					validIngredientPreparationsTableName, idColumn,
					validIngredientPreparationsTableName,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					validIngredientPreparationsTableName,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidIngredientPreparationsForIngredient",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(validIngredientPreparationsTableName, true, true, []string{}),
					buildTotalCountSelect(validIngredientPreparationsTableName, true, []string{}),
					validIngredientPreparationsTableName,
					validIngredientsTableName,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					idColumn,
					buildFilterConditions(validIngredientPreparationsTableName, true, false),
					buildCursorLimitClause(validIngredientPreparationsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidIngredientPreparationsForPreparation",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(validIngredientPreparationsTableName, true, true, []string{}),
					buildTotalCountSelect(validIngredientPreparationsTableName, true, []string{}),
					validIngredientPreparationsTableName,
					validIngredientsTableName,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					idColumn,
					buildFilterConditions(validIngredientPreparationsTableName, true, false),
					buildCursorLimitClause(validIngredientPreparationsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidIngredientPreparations",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(validIngredientPreparationsTableName, true, true, []string{}),
					buildTotalCountSelect(validIngredientPreparationsTableName, true, []string{}),
					validIngredientPreparationsTableName,
					validIngredientsTableName,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					buildFilterConditions(validIngredientPreparationsTableName, true, false),
					buildCursorLimitClause(validIngredientPreparationsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidIngredientPreparation",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					validIngredientPreparationsTableName,
					validIngredientsTableName,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					validIngredientsTableName,
					archivedAtColumn,
					validPreparationsTableName,
					archivedAtColumn,
					validIngredientPreparationsTableName,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidIngredientPreparationsByIDs",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = ANY(sqlc.arg(ids)::text[]);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					validIngredientPreparationsTableName,
					validIngredientsTableName,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					validIngredientsTableName,
					archivedAtColumn,
					validPreparationsTableName,
					archivedAtColumn,
					validIngredientPreparationsTableName,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "ValidIngredientPreparationPairIsValid",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s
	FROM %s
	WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s)
	AND %s IS NULL
);`,
					idColumn,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientIDColumn,
					validPreparationIDColumn,
					validPreparationIDColumn,
					archivedAtColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "SearchValidIngredientPreparationsByPreparationAndIngredientName",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
	JOIN %s ON %s.%s = %s.%s
	JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s %s
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(validIngredientPreparationsTableName, true, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", validPreparationsTableName, idColumn, idColumn)),
					buildTotalCountSelect(validIngredientPreparationsTableName, true, []string{}, fmt.Sprintf("%s.%s = sqlc.arg(%s)", validPreparationsTableName, idColumn, idColumn)),
					validIngredientPreparationsTableName,
					validIngredientsTableName,
					validIngredientPreparationsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validIngredientPreparationsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validIngredientPreparationsTableName,
					archivedAtColumn,
					validIngredientsTableName,
					archivedAtColumn,
					validPreparationsTableName,
					archivedAtColumn,
					validPreparationsTableName,
					idColumn,
					idColumn,
					validIngredientsTableName,
					nameColumn,
					buildILIKEForArgument("name_query"),
					buildFilterConditions(
						validIngredientPreparationsTableName,
						true,
						true,
						fmt.Sprintf("%s.%s = sqlc.arg(%s)", validPreparationsTableName, idColumn, idColumn),
					),
					buildCursorLimitClause(validIngredientPreparationsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateValidIngredientPreparation",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					validIngredientPreparationsTableName,
					strings.Join(applyToEach(filterForUpdate(validIngredientPreparationsColumns), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					idColumn,
					idColumn,
				)),
			},
		}
	default:
		return nil
	}
}
