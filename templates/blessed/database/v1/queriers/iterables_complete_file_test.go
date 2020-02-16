package queriers

import (
	"bytes"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PostgresIterables(T *testing.T) {
	T.Parallel()

	standardFields := []models.DataField{
		{
			Name: wordsmith.FromSingularPascalCase("First"),
		},
		{
			Name: wordsmith.FromSingularPascalCase("Second"),
		},
		{
			Name: wordsmith.FromSingularPascalCase("Third"),
		},
	}

	T.Run("belongs to user", func(t *testing.T) {
		p := &models.Project{
			OutputPath: "gitlab.com/verygoodsoftwarenotvirus/nafftesting",
		}

		dbv := wordsmith.FromSingularPascalCase("Postgres")
		dt := models.DataType{
			Name:          wordsmith.FromSingularPascalCase("Something"),
			BelongsToUser: true,
			Fields:        standardFields,
		}

		f := iterablesDotGo(p, dbv, dt)

		var b bytes.Buffer
		renderErr := f.Render(&b)
		require.NoError(t, renderErr)

		actual := b.String()
		expected := `package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	database "gitlab.com/verygoodsoftwarenotvirus/nafftesting/database/v1"
	models "gitlab.com/verygoodsoftwarenotvirus/nafftesting/models/v1"
	"sync"
)

const (
	somethingsTableName            = "somethings"
	somethingsTableOwnershipColumn = "belongs_to_user"
)

var (
	somethingsTableColumns = []string{
		"id",
		"first",
		"second",
		"third",
		"created_on",
		"updated_on",
		"archived_on",
		somethingsTableOwnershipColumn,
	}
)

// scanSomething takes a database Scanner (i.e. *sql.Row) and scans the result into a somethings struct
func scanSomething(scan database.Scanner) (*models.Something, error) {
	x := &models.Something{}

	if err := scan.Scan(
		&x.ID,
		&x.First,
		&x.Second,
		&x.Third,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanSomethings takes a logger and some database rows and turns them into a slice of somethings
func scanSomethings(logger logging.Logger, rows *sql.Rows) ([]models.Something, error) {
	var list []models.Something

	for rows.Next() {
		x, err := scanSomething(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetSomethingQuery constructs a SQL query for fetching a something with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetSomethingQuery(somethingID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"id":                           somethingID,
			somethingsTableOwnershipColumn: userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetSomething fetches a something from the postgres database
func (p *Postgres) GetSomething(ctx context.Context, somethingID, userID uint64) (*models.Something, error) {
	query, args := p.buildGetSomethingQuery(somethingID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanSomething(row)
}

// buildGetSomethingCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of somethings belonging to a given user that meet a given query
func (p *Postgres) buildGetSomethingCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"archived_on":                  nil,
			somethingsTableOwnershipColumn: userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetSomethingCount will fetch the count of somethings from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetSomethingCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetSomethingCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allSomethingsCountQueryBuilder sync.Once
	allSomethingsCountQuery        string
)

// buildGetAllSomethingsCountQuery returns a query that fetches the total number of somethings in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllSomethingsCountQuery() string {
	allSomethingsCountQueryBuilder.Do(func() {
		var err error
		allSomethingsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(somethingsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allSomethingsCountQuery
}

// GetAllSomethingsCount will fetch the count of somethings from the database
func (p *Postgres) GetAllSomethingsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllSomethingsCountQuery()).Scan(&count)
	return count, err
}

// buildGetSomethingsQuery builds a SQL query selecting somethings that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetSomethingsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"archived_on":                  nil,
			somethingsTableOwnershipColumn: userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetSomethings fetches a list of somethings from the database that meet a particular filter
func (p *Postgres) GetSomethings(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.SomethingList, error) {
	query, args := p.buildGetSomethingsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for somethings")
	}

	list, err := scanSomethings(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetSomethingCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching something count: %w", err)
	}

	x := &models.SomethingList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Somethings: list,
	}

	return x, nil
}

// GetAllSomethingsForUser fetches every something belonging to a user
func (p *Postgres) GetAllSomethingsForUser(ctx context.Context, userID uint64) ([]models.Something, error) {
	query, args := p.buildGetSomethingsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching somethings for user")
	}

	list, err := scanSomethings(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateSomethingQuery takes a something and returns a creation query for that something and the relevant arguments.
func (p *Postgres) buildCreateSomethingQuery(input *models.Something) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(somethingsTableName).
		Columns(
			"first",
			"second",
			"third",
			somethingsTableOwnershipColumn,
		).
		Values(
			input.First,
			input.Second,
			input.Third,
			input.BelongsToUser,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateSomething creates a something in the database
func (p *Postgres) CreateSomething(ctx context.Context, input *models.SomethingCreationInput) (*models.Something, error) {
	x := &models.Something{
		First:         input.First,
		Second:        input.Second,
		Third:         input.Third,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildCreateSomethingQuery(x)

	// create the something
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing something creation query: %w", err)
	}

	return x, nil
}

// buildUpdateSomethingQuery takes a something and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateSomethingQuery(input *models.Something) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(somethingsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           input.ID,
			somethingsTableOwnershipColumn: input.BelongsToUser,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateSomething updates a particular something. Note that UpdateSomething expects the provided input to have a valid ID.
func (p *Postgres) UpdateSomething(ctx context.Context, input *models.Something) error {
	query, args := p.buildUpdateSomethingQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveSomethingQuery returns a SQL query which marks a given something belonging to a given user as archived.
func (p *Postgres) buildArchiveSomethingQuery(somethingID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(somethingsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           somethingID,
			"archived_on":                  nil,
			somethingsTableOwnershipColumn: userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveSomething marks a something as archived in the database
func (p *Postgres) ArchiveSomething(ctx context.Context, somethingID, userID uint64) error {
	query, args := p.buildArchiveSomethingQuery(somethingID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to other type", func(t *testing.T) {
		p := &models.Project{
			OutputPath: "gitlab.com/verygoodsoftwarenotvirus/nafftesting",
		}

		dbv := wordsmith.FromSingularPascalCase("Postgres")
		dt := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Something"),
			BelongsToStruct: wordsmith.FromSingularPascalCase("Arbitrary"),
			Fields:          standardFields,
		}

		f := iterablesDotGo(p, dbv, dt)

		var b bytes.Buffer
		renderErr := f.Render(&b)
		require.NoError(t, renderErr)

		actual := b.String()
		expected := `package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	database "gitlab.com/verygoodsoftwarenotvirus/nafftesting/database/v1"
	models "gitlab.com/verygoodsoftwarenotvirus/nafftesting/models/v1"
	"sync"
)

const (
	somethingsTableName            = "somethings"
	somethingsTableOwnershipColumn = "belongs_to_arbitrary"
)

var (
	somethingsTableColumns = []string{
		"id",
		"first",
		"second",
		"third",
		"created_on",
		"updated_on",
		"archived_on",
		somethingsTableOwnershipColumn,
	}
)

// scanSomething takes a database Scanner (i.e. *sql.Row) and scans the result into a somethings struct
func scanSomething(scan database.Scanner) (*models.Something, error) {
	x := &models.Something{}

	if err := scan.Scan(
		&x.ID,
		&x.First,
		&x.Second,
		&x.Third,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToArbitrary,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanSomethings takes a logger and some database rows and turns them into a slice of somethings
func scanSomethings(logger logging.Logger, rows *sql.Rows) ([]models.Something, error) {
	var list []models.Something

	for rows.Next() {
		x, err := scanSomething(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetSomethingQuery constructs a SQL query for fetching a something with a given ID belong to an arbitrary with a given ID.
func (p *Postgres) buildGetSomethingQuery(somethingID, arbitraryID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"id":                           somethingID,
			somethingsTableOwnershipColumn: arbitraryID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetSomething fetches a something from the postgres database
func (p *Postgres) GetSomething(ctx context.Context, somethingID, arbitraryID uint64) (*models.Something, error) {
	query, args := p.buildGetSomethingQuery(somethingID, arbitraryID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanSomething(row)
}

// buildGetSomethingCountQuery takes a QueryFilter and an arbitrary ID and returns a SQL query (and the relevant arguments) for
// fetching the number of somethings belonging to a given arbitrary that meet a given query
func (p *Postgres) buildGetSomethingCountQuery(filter *models.QueryFilter, arbitraryID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"archived_on":                  nil,
			somethingsTableOwnershipColumn: arbitraryID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetSomethingCount will fetch the count of somethings from the database that meet a particular filter and belongs to a particular arbitrary.
func (p *Postgres) GetSomethingCount(ctx context.Context, filter *models.QueryFilter, arbitraryID uint64) (count uint64, err error) {
	query, args := p.buildGetSomethingCountQuery(filter, arbitraryID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allSomethingsCountQueryBuilder sync.Once
	allSomethingsCountQuery        string
)

// buildGetAllSomethingsCountQuery returns a query that fetches the total number of somethings in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllSomethingsCountQuery() string {
	allSomethingsCountQueryBuilder.Do(func() {
		var err error
		allSomethingsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(somethingsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allSomethingsCountQuery
}

// GetAllSomethingsCount will fetch the count of somethings from the database
func (p *Postgres) GetAllSomethingsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllSomethingsCountQuery()).Scan(&count)
	return count, err
}

// buildGetSomethingsQuery builds a SQL query selecting somethings that adhere to a given QueryFilter and belong to a given arbitrary,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetSomethingsQuery(filter *models.QueryFilter, arbitraryID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"archived_on":                  nil,
			somethingsTableOwnershipColumn: arbitrariesID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetSomethings fetches a list of somethings from the database that meet a particular filter
func (p *Postgres) GetSomethings(ctx context.Context, filter *models.QueryFilter, arbitraryID uint64) (*models.SomethingList, error) {
	query, args := p.buildGetSomethingsQuery(filter, arbitraryID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for somethings")
	}

	list, err := scanSomethings(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetSomethingCount(ctx, filter, arbitraryID)
	if err != nil {
		return nil, fmt.Errorf("fetching something count: %w", err)
	}

	x := &models.SomethingList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Somethings: list,
	}

	return x, nil
}

// buildCreateSomethingQuery takes a something and returns a creation query for that something and the relevant arguments.
func (p *Postgres) buildCreateSomethingQuery(input *models.Something) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(somethingsTableName).
		Columns(
			"first",
			"second",
			"third",
			somethingsTableOwnershipColumn,
		).
		Values(
			input.First,
			input.Second,
			input.Third,
			input.BelongsToArbitrary,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateSomething creates a something in the database
func (p *Postgres) CreateSomething(ctx context.Context, input *models.SomethingCreationInput) (*models.Something, error) {
	x := &models.Something{
		First:              input.First,
		Second:             input.Second,
		Third:              input.Third,
		BelongsToArbitrary: input.BelongsToArbitrary,
	}

	query, args := p.buildCreateSomethingQuery(x)

	// create the something
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing something creation query: %w", err)
	}

	return x, nil
}

// buildUpdateSomethingQuery takes a something and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateSomethingQuery(input *models.Something) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(somethingsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                           input.ID,
			somethingsTableOwnershipColumn: input.BelongsToArbitrary,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateSomething updates a particular something. Note that UpdateSomething expects the provided input to have a valid ID.
func (p *Postgres) UpdateSomething(ctx context.Context, input *models.Something) error {
	query, args := p.buildUpdateSomethingQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveSomethingQuery returns a SQL query which marks a given something belonging to a given arbitrary as archived.
func (p *Postgres) buildArchiveSomethingQuery(somethingID, arbitraryID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(somethingsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                            somethingID,
			"archived_on":                   nil,
			arbitrariesTableOwnershipColumn: arbitraryID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveSomething marks a something as archived in the database
func (p *Postgres) ArchiveSomething(ctx context.Context, somethingID, arbitraryID uint64) error {
	query, args := p.buildArchiveSomethingQuery(somethingID, arbitraryID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`

		assert.Equal(t, expected, actual)
	})

	T.Run("belongs to nobody", func(t *testing.T) {
		p := &models.Project{
			OutputPath: "gitlab.com/verygoodsoftwarenotvirus/nafftesting",
		}

		dbv := wordsmith.FromSingularPascalCase("Postgres")
		dt := models.DataType{
			Name:            wordsmith.FromSingularPascalCase("Something"),
			BelongsToNobody: true,
			Fields:          standardFields,
		}

		f := iterablesDotGo(p, dbv, dt)

		var b bytes.Buffer
		renderErr := f.Render(&b)
		require.NoError(t, renderErr)

		actual := b.String()
		expected := `package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	database "gitlab.com/verygoodsoftwarenotvirus/nafftesting/database/v1"
	models "gitlab.com/verygoodsoftwarenotvirus/nafftesting/models/v1"
	"sync"
)

const (
	somethingsTableName = "somethings"
)

var (
	somethingsTableColumns = []string{
		"id",
		"first",
		"second",
		"third",
		"created_on",
		"updated_on",
		"archived_on",
	}
)

// scanSomething takes a database Scanner (i.e. *sql.Row) and scans the result into a somethings struct
func scanSomething(scan database.Scanner) (*models.Something, error) {
	x := &models.Something{}

	if err := scan.Scan(
		&x.ID,
		&x.First,
		&x.Second,
		&x.Third,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanSomethings takes a logger and some database rows and turns them into a slice of somethings
func scanSomethings(logger logging.Logger, rows *sql.Rows) ([]models.Something, error) {
	var list []models.Something

	for rows.Next() {
		x, err := scanSomething(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetSomethingQuery constructs a SQL query for fetching a something with a given ID.
func (p *Postgres) buildGetSomethingQuery(somethingID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"id": somethingID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetSomething fetches a something from the postgres database
func (p *Postgres) GetSomething(ctx context.Context, somethingID uint64) (*models.Something, error) {
	query, args := p.buildGetSomethingQuery(somethingID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanSomething(row)
}

// buildGetSomethingCountQuery takes a QueryFilter and returns a SQL query (and the relevant arguments) for
// fetching the number of somethings that meet a given query
func (p *Postgres) buildGetSomethingCountQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetSomethingCount will fetch the count of somethings from the database that meet a particular filter.
func (p *Postgres) GetSomethingCount(ctx context.Context, filter *models.QueryFilter) (count uint64, err error) {
	query, args := p.buildGetSomethingCountQuery(filter)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allSomethingsCountQueryBuilder sync.Once
	allSomethingsCountQuery        string
)

// buildGetAllSomethingsCountQuery returns a query that fetches the total number of somethings in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllSomethingsCountQuery() string {
	allSomethingsCountQueryBuilder.Do(func() {
		var err error
		allSomethingsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(somethingsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allSomethingsCountQuery
}

// GetAllSomethingsCount will fetch the count of somethings from the database
func (p *Postgres) GetAllSomethingsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllSomethingsCountQuery()).Scan(&count)
	return count, err
}

// buildGetSomethingsQuery builds a SQL query selecting somethings that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetSomethingsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(somethingsTableColumns...).
		From(somethingsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetSomethings fetches a list of somethings from the database that meet a particular filter
func (p *Postgres) GetSomethings(ctx context.Context, filter *models.QueryFilter) (*models.SomethingList, error) {
	query, args := p.buildGetSomethingsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for somethings")
	}

	list, err := scanSomethings(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetSomethingCount(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("fetching something count: %w", err)
	}

	x := &models.SomethingList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Somethings: list,
	}

	return x, nil
}

// buildCreateSomethingQuery takes a something and returns a creation query for that something and the relevant arguments.
func (p *Postgres) buildCreateSomethingQuery(input *models.Something) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(somethingsTableName).
		Columns(
			"first",
			"second",
			"third",
		).
		Values(
			input.First,
			input.Second,
			input.Third,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateSomething creates a something in the database
func (p *Postgres) CreateSomething(ctx context.Context, input *models.SomethingCreationInput) (*models.Something, error) {
	x := &models.Something{
		First:  input.First,
		Second: input.Second,
		Third:  input.Third,
	}

	query, args := p.buildCreateSomethingQuery(x)

	// create the something
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing something creation query: %w", err)
	}

	return x, nil
}

// buildUpdateSomethingQuery takes a something and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateSomethingQuery(input *models.Something) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(somethingsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateSomething updates a particular something. Note that UpdateSomething expects the provided input to have a valid ID.
func (p *Postgres) UpdateSomething(ctx context.Context, input *models.Something) error {
	query, args := p.buildUpdateSomethingQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveSomethingQuery returns a SQL query which marks a given something as archived.
func (p *Postgres) buildArchiveSomethingQuery(somethingID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(somethingsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          somethingID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveSomething marks a something as archived in the database
func (p *Postgres) ArchiveSomething(ctx context.Context, somethingID uint64) error {
	query, args := p.buildArchiveSomethingQuery(somethingID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`

		assert.Equal(t, expected, actual)
	})
}
