package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	accountInstrumentOwnershipsTableName = "account_instrument_ownerships"
)

func init() {
	registerTableName(accountInstrumentOwnershipsTableName)
}

var accountInstrumentOwnershipsColumns = []string{
	idColumn,
	notesColumn,
	"quantity",
	validInstrumentIDColumn,
	belongsToAccountColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildAccountInstrumentOwnershipQueries(database string) []*Query {
	switch database {
	case postgres:

		insertColumns := filterForInsert(accountInstrumentOwnershipsColumns)

		fullSelectColumns := mergeColumns(
			applyToEach(filterFromSlice(accountInstrumentOwnershipsColumns, validInstrumentIDColumn), func(i int, s string) string {
				return fmt.Sprintf("%s.%s", accountInstrumentOwnershipsTableName, s)
			}),
			applyToEach(validInstrumentsColumns, func(i int, s string) string {
				return fmt.Sprintf("%s.%s as valid_instrument_%s", validInstrumentsTableName, s, s)
			}),
			3,
		)

		return []*Query{
			{
				Annotation: QueryAnnotation{
					Name: "ArchiveAccountInstrumentOwnership",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
					accountInstrumentOwnershipsTableName,
					archivedAtColumn, currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
					belongsToAccountColumn, belongsToAccountColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CreateAccountInstrumentOwnership",
					Type: ExecType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
					accountInstrumentOwnershipsTableName,
					strings.Join(insertColumns, ",\n\t"),
					strings.Join(applyToEach(insertColumns, func(i int, s string) string {
						return fmt.Sprintf("sqlc.arg(%s)", s)
					}), ",\n\t"),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "CheckAccountInstrumentOwnershipExistence",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
	SELECT %s.id
	FROM %s
	WHERE %s.%s IS NULL
		AND %s.%s = sqlc.arg(%s)
		AND %s.%s = sqlc.arg(%s)
);`,
					accountInstrumentOwnershipsTableName,
					accountInstrumentOwnershipsTableName,
					accountInstrumentOwnershipsTableName,
					archivedAtColumn,
					accountInstrumentOwnershipsTableName, idColumn, idColumn,
					accountInstrumentOwnershipsTableName, belongsToAccountColumn, belongsToAccountColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetAccountInstrumentOwnerships",
					Type: ManyType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
INNER JOIN %s ON %s.%s = %s.%s
WHERE
	%s.%s IS NULL
	%s
GROUP BY
	%s.%s,
	%s.%s
%s;`,
					strings.Join(fullSelectColumns, ",\n\t"),
					buildFilterCountSelect(accountInstrumentOwnershipsTableName, true, true, []string{}, "account_instrument_ownerships.belongs_to_account = sqlc.arg(account_id)"),
					buildTotalCountSelect(accountInstrumentOwnershipsTableName, true, []string{}, "account_instrument_ownerships.belongs_to_account = sqlc.arg(account_id)"),
					accountInstrumentOwnershipsTableName,
					validInstrumentsTableName, accountInstrumentOwnershipsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
					accountInstrumentOwnershipsTableName, archivedAtColumn,
					buildFilterConditions(accountInstrumentOwnershipsTableName, true, true, "account_instrument_ownerships.belongs_to_account = sqlc.arg(account_id)"),
					accountInstrumentOwnershipsTableName, idColumn,
					validInstrumentsTableName, idColumn,
					buildCursorLimitClause(accountInstrumentOwnershipsTableName),
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "GetAccountInstrumentOwnership",
					Type: OneType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
INNER JOIN %s ON %s.%s = %s.%s
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
					strings.Join(fullSelectColumns, ",\n\t"),
					accountInstrumentOwnershipsTableName,
					validInstrumentsTableName, accountInstrumentOwnershipsTableName, validInstrumentIDColumn, validInstrumentsTableName, idColumn,
					accountInstrumentOwnershipsTableName,
					archivedAtColumn,
					accountInstrumentOwnershipsTableName, idColumn, idColumn,
					accountInstrumentOwnershipsTableName, belongsToAccountColumn, belongsToAccountColumn,
				)),
			},
			{
				Annotation: QueryAnnotation{
					Name: "UpdateAccountInstrumentOwnership",
					Type: ExecRowsType,
				},
				Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s,
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
					accountInstrumentOwnershipsTableName,
					strings.Join(applyToEach(filterForUpdate(accountInstrumentOwnershipsColumns, belongsToAccountColumn), func(i int, s string) string {
						return fmt.Sprintf("%s = sqlc.arg(%s)", s, s)
					}), ",\n\t"),
					lastUpdatedAtColumn,
					currentTimeExpression,
					archivedAtColumn,
					idColumn, idColumn,
					accountInstrumentOwnershipsTableName, belongsToAccountColumn, belongsToAccountColumn,
				)),
			},
		}
	default:
		return nil
	}
}
