package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	validPrepTaskConfigsTableName = "valid_prep_task_configs"
)

func init() {
	registerTableName(validPrepTaskConfigsTableName)
}

var validPrepTaskConfigsColumns = []string{
	idColumn,
	validIngredientIDColumn,
	validPreparationIDColumn,
	"minimum_storage_duration_in_seconds",
	"maximum_storage_duration_in_seconds",
	"storage_container_type",
	"minimum_storage_temperature_in_celsius",
	"maximum_storage_temperature_in_celsius",
	"storage_instructions",
	notesColumn,
	"source",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildValidPrepTaskConfigsQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(validPrepTaskConfigsColumns)

		fullSelectColumns := mergeColumns(
			mergeColumns(
				applyToEach(filterFromSlice(validPrepTaskConfigsColumns, "valid_preparation_id", "valid_ingredient_id"), func(i int, s string) string {
					return fmt.Sprintf("%s.%s as valid_prep_task_config_%s", validPrepTaskConfigsTableName, s, s)
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
					Name: "ArchiveValidPrepTaskConfig",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = %s WHERE %s IS NULL AND %s = sqlc.arg(%s);`,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateValidPrepTaskConfig",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					validPrepTaskConfigsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						switch s {
						case "maximum_storage_duration_in_seconds",
							"minimum_storage_temperature_in_celsius",
							"maximum_storage_temperature_in_celsius":
							return fmt.Sprintf("sqlc.narg(%s)", s)
						default:
							return fmt.Sprintf("sqlc.arg(%s)", s)
						}
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckValidPrepTaskConfigExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.%s
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
);`,
					validPrepTaskConfigsTableName, idColumn,
					validPrepTaskConfigsTableName,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					validPrepTaskConfigsTableName,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidPrepTaskConfigsForIngredient",
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
					buildFilterCountSelect(validPrepTaskConfigsTableName, true, true, []string{}),
					buildTotalCountSelect(validPrepTaskConfigsTableName, true, []string{}),
					validPrepTaskConfigsTableName,
					validIngredientsTableName,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					idColumn,
					buildFilterConditions(validPrepTaskConfigsTableName, true, false),
					buildCursorLimitClause(validPrepTaskConfigsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidPrepTaskConfigsForPreparation",
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
					buildFilterCountSelect(validPrepTaskConfigsTableName, true, true, []string{}),
					buildTotalCountSelect(validPrepTaskConfigsTableName, true, []string{}),
					validPrepTaskConfigsTableName,
					validIngredientsTableName,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					idColumn,
					buildFilterConditions(validPrepTaskConfigsTableName, true, false),
					buildCursorLimitClause(validPrepTaskConfigsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidPrepTaskConfigsForIngredientAndPreparation",
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
	AND %s.%s = sqlc.arg(%s)
	%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(validPrepTaskConfigsTableName, true, true, []string{}),
					buildTotalCountSelect(validPrepTaskConfigsTableName, true, []string{}),
					validPrepTaskConfigsTableName,
					validIngredientsTableName,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					validIngredientIDColumn,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					validPreparationIDColumn,
					buildFilterConditions(validPrepTaskConfigsTableName, true, false),
					buildCursorLimitClause(validPrepTaskConfigsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidPrepTaskConfigs",
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
					buildFilterCountSelect(validPrepTaskConfigsTableName, true, true, []string{}),
					buildTotalCountSelect(validPrepTaskConfigsTableName, true, []string{}),
					validPrepTaskConfigsTableName,
					validIngredientsTableName,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					buildFilterConditions(validPrepTaskConfigsTableName, true, false),
					buildCursorLimitClause(validPrepTaskConfigsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetValidPrepTaskConfig",
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
					validPrepTaskConfigsTableName,
					validIngredientsTableName,
					validPrepTaskConfigsTableName,
					validIngredientIDColumn,
					validIngredientsTableName,
					idColumn,
					validPreparationsTableName,
					validPrepTaskConfigsTableName,
					validPreparationIDColumn,
					validPreparationsTableName,
					idColumn,
					validPrepTaskConfigsTableName,
					archivedAtColumn,
					validIngredientsTableName,
					archivedAtColumn,
					validPreparationsTableName,
					archivedAtColumn,
					validPrepTaskConfigsTableName,
					idColumn,
					idColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateValidPrepTaskConfig",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
					validPrepTaskConfigsTableName,
					strings.Join(applyToEach(filterForUpdate(validPrepTaskConfigsColumns), func(i int, s string) string {
						switch s {
						case "maximum_storage_duration_in_seconds",
							"minimum_storage_temperature_in_celsius",
							"maximum_storage_temperature_in_celsius":
							return fmt.Sprintf("%s = sqlc.narg(%s)", s, s)
						default:
							return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
						}
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
