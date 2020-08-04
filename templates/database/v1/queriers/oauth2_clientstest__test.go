package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
	"time"
)

func buildMockRowsFromOAuth2Client(clients ...*v1.OAuth2Client) *gosqlmock.Rows {
	columns := oauth2ClientsTableColumns
	exampleRows := gosqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			strings.Join(c.Scopes, scopesSeparator),
			c.RedirectURI,
			c.ClientSecret,
			c.CreatedOn,
			c.LastUpdatedOn,
			c.ArchivedOn,
			c.BelongsToUser,
		}
		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromOAuth2Client(c *v1.OAuth2Client) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(oauth2ClientsTableColumns).AddRow(
		c.ArchivedOn,
		c.Name,
		c.ClientID,
		strings.Join(c.Scopes, scopesSeparator),
		c.RedirectURI,
		c.ClientSecret,
		c.CreatedOn,
		c.LastUpdatedOn,
		c.BelongsToUser,
		c.ID,
	)

	return exampleRows
}

func TestPostgres_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanOAuth2Clients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanOAuth2Clients(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = $1"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
		}
		actualQuery, actualArgs := p.buildGetOAuth2ClientByClientIDQuery(exampleOAuth2Client.ClientID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous row", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := p.buildGetAllOAuth2ClientsQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := []*v1.OAuth2Client{
			exampleOAuth2Client,
			exampleOAuth2Client,
			exampleOAuth2Client,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				exampleOAuth2Client,
				exampleOAuth2Client,
				exampleOAuth2Client,
			),
		)

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 ORDER BY oauth2_clients.id"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		expected := []*v1.OAuth2Client{
			&exampleOAuth2ClientList.Clients[0],
			&exampleOAuth2ClientList.Clients[1],
			&exampleOAuth2ClientList.Clients[2],
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildMockRowsFromOAuth2Client(expected...))

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 AND oauth2_clients.id = $2"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := p.buildGetOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 AND oauth2_clients.id = $2"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := p.buildGetAllOAuth2ClientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 AND oauth2_clients.created_on > $2 AND oauth2_clients.created_on < $3 AND oauth2_clients.last_updated_on > $4 AND oauth2_clients.last_updated_on < $5 ORDER BY oauth2_clients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetOAuth2ClientsForUserQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 ORDER BY oauth2_clients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				&exampleOAuth2ClientList.Clients[0],
				&exampleOAuth2ClientList.Clients[1],
				&exampleOAuth2ClientList.Clients[2],
			),
		)

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error reading from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(fake.BuildFakeOAuth2Client()))

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		exampleRows := gosqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleOAuth2Client.ID, exampleOAuth2Client.CreatedOn)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnRows(exampleRows)

		actual, err := p.CreateOAuth2Client(ctx, expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnError(errors.New("blah"))

		actual, err := p.CreateOAuth2Client(ctx, expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET client_id = $1, client_secret = $2, scopes = $3, redirect_uri = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $5 AND id = $6 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := p.buildUpdateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET client_id = $1, client_secret = $2, scopes = $3, redirect_uri = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $5 AND id = $6 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"last_updated_on"}).AddRow(time.Now().Unix()))

		err := p.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		err := p.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := p.buildArchiveOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := p.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(errors.New("blah"))

		err := p.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func buildMockRowsFromOAuth2Client(clients ...*v1.OAuth2Client) *gosqlmock.Rows {
	columns := oauth2ClientsTableColumns
	exampleRows := gosqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			strings.Join(c.Scopes, scopesSeparator),
			c.RedirectURI,
			c.ClientSecret,
			c.CreatedOn,
			c.LastUpdatedOn,
			c.ArchivedOn,
			c.BelongsToUser,
		}
		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromOAuth2Client(c *v1.OAuth2Client) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(oauth2ClientsTableColumns).AddRow(
		c.ArchivedOn,
		c.Name,
		c.ClientID,
		strings.Join(c.Scopes, scopesSeparator),
		c.RedirectURI,
		c.ClientSecret,
		c.CreatedOn,
		c.LastUpdatedOn,
		c.BelongsToUser,
		c.ID,
	)

	return exampleRows
}

func TestSqlite_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := s.scanOAuth2Clients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := s.scanOAuth2Clients(mockRows)
		assert.NoError(t, err)
	})
}

func TestSqlite_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
		}
		actualQuery, actualArgs := s.buildGetOAuth2ClientByClientIDQuery(exampleOAuth2Client.ClientID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous row", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := s.buildGetAllOAuth2ClientsQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := []*v1.OAuth2Client{
			exampleOAuth2Client,
			exampleOAuth2Client,
			exampleOAuth2Client,
		}

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				exampleOAuth2Client,
				exampleOAuth2Client,
				exampleOAuth2Client,
			),
		)

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing query", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		expected := []*v1.OAuth2Client{
			&exampleOAuth2ClientList.Clients[0],
			&exampleOAuth2ClientList.Clients[1],
			&exampleOAuth2ClientList.Clients[2],
		}

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildMockRowsFromOAuth2Client(expected...))

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := s.buildGetOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := s.buildGetAllOAuth2ClientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		expectedCount := uint64(123)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.created_on > ? AND oauth2_clients.created_on < ? AND oauth2_clients.last_updated_on > ? AND oauth2_clients.last_updated_on < ? ORDER BY oauth2_clients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := s.buildGetOAuth2ClientsForUserQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				&exampleOAuth2ClientList.Clients[0],
				&exampleOAuth2ClientList.Clients[1],
				&exampleOAuth2ClientList.Clients[2],
			),
		)

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error reading from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(fake.BuildFakeOAuth2Client()))

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"
		expectedArgs := []interface{}{
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		}
		actualQuery, actualArgs := s.buildCreateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		exampleRows := gosqlmock.NewResult(int64(exampleOAuth2Client.ID), 1)

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnResult(exampleRows)

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleOAuth2Client.CreatedOn)
		s.timeTeller = mtt

		actual, err := s.CreateOAuth2Client(ctx, expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnError(errors.New("blah"))

		actual, err := s.CreateOAuth2Client(ctx, expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := s.buildUpdateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		err := s.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := s.buildArchiveOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(errors.New("blah"))

		err := s.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsTestDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func buildMockRowsFromOAuth2Client(clients ...*v1.OAuth2Client) *gosqlmock.Rows {
	columns := oauth2ClientsTableColumns
	exampleRows := gosqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			strings.Join(c.Scopes, scopesSeparator),
			c.RedirectURI,
			c.ClientSecret,
			c.CreatedOn,
			c.LastUpdatedOn,
			c.ArchivedOn,
			c.BelongsToUser,
		}
		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromOAuth2Client(c *v1.OAuth2Client) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(oauth2ClientsTableColumns).AddRow(
		c.ArchivedOn,
		c.Name,
		c.ClientID,
		strings.Join(c.Scopes, scopesSeparator),
		c.RedirectURI,
		c.ClientSecret,
		c.CreatedOn,
		c.LastUpdatedOn,
		c.BelongsToUser,
		c.ID,
	)

	return exampleRows
}

func TestMariaDB_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := m.scanOAuth2Clients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := m.scanOAuth2Clients(mockRows)
		assert.NoError(t, err)
	})
}

func TestMariaDB_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
		}
		actualQuery, actualArgs := m.buildGetOAuth2ClientByClientIDQuery(exampleOAuth2Client.ClientID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous row", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := m.buildGetAllOAuth2ClientsQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := []*v1.OAuth2Client{
			exampleOAuth2Client,
			exampleOAuth2Client,
			exampleOAuth2Client,
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				exampleOAuth2Client,
				exampleOAuth2Client,
				exampleOAuth2Client,
			),
		)

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing query", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		expected := []*v1.OAuth2Client{
			&exampleOAuth2ClientList.Clients[0],
			&exampleOAuth2ClientList.Clients[1],
			&exampleOAuth2ClientList.Clients[2],
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildMockRowsFromOAuth2Client(expected...))

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := m.buildGetOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := m.buildGetAllOAuth2ClientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		expectedCount := uint64(123)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.created_on > ? AND oauth2_clients.created_on < ? AND oauth2_clients.last_updated_on > ? AND oauth2_clients.last_updated_on < ? ORDER BY oauth2_clients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := m.buildGetOAuth2ClientsForUserQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				&exampleOAuth2ClientList.Clients[0],
				&exampleOAuth2ClientList.Clients[1],
				&exampleOAuth2ClientList.Clients[2],
			),
		)

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error reading from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(fake.BuildFakeOAuth2Client()))

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"
		expectedArgs := []interface{}{
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		}
		actualQuery, actualArgs := m.buildCreateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		exampleRows := gosqlmock.NewResult(int64(exampleOAuth2Client.ID), 1)

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnResult(exampleRows)

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleOAuth2Client.CreatedOn)
		m.timeTeller = mtt

		actual, err := m.CreateOAuth2Client(ctx, expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnError(errors.New("blah"))

		actual, err := m.CreateOAuth2Client(ctx, expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := m.buildUpdateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		err := m.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := m.buildArchiveOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(errors.New("blah"))

		err := m.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildMockRowsFromOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildMockRowsFromOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"database/sql/driver"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

func buildMockRowsFromOAuth2Client(clients ...*v1.OAuth2Client) *gosqlmock.Rows {
	columns := oauth2ClientsTableColumns
	exampleRows := gosqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			strings.Join(c.Scopes, scopesSeparator),
			c.RedirectURI,
			c.ClientSecret,
			c.CreatedOn,
			c.LastUpdatedOn,
			c.ArchivedOn,
			c.BelongsToUser,
		}
		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildBuildMockRowsFromOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"database/sql/driver"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

func buildMockRowsFromOAuth2Client(clients ...*v1.OAuth2Client) *gosqlmock.Rows {
	columns := oauth2ClientsTableColumns
	exampleRows := gosqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			strings.Join(c.Scopes, scopesSeparator),
			c.RedirectURI,
			c.ClientSecret,
			c.CreatedOn,
			c.LastUpdatedOn,
			c.ArchivedOn,
			c.BelongsToUser,
		}
		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildBuildMockRowsFromOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"database/sql/driver"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

func buildMockRowsFromOAuth2Client(clients ...*v1.OAuth2Client) *gosqlmock.Rows {
	columns := oauth2ClientsTableColumns
	exampleRows := gosqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			strings.Join(c.Scopes, scopesSeparator),
			c.RedirectURI,
			c.ClientSecret,
			c.CreatedOn,
			c.LastUpdatedOn,
			c.ArchivedOn,
			c.BelongsToUser,
		}
		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildErroneousMockRowFromOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildErroneousMockRowFromOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

func buildErroneousMockRowFromOAuth2Client(c *v1.OAuth2Client) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(oauth2ClientsTableColumns).AddRow(
		c.ArchivedOn,
		c.Name,
		c.ClientID,
		strings.Join(c.Scopes, scopesSeparator),
		c.RedirectURI,
		c.ClientSecret,
		c.CreatedOn,
		c.LastUpdatedOn,
		c.BelongsToUser,
		c.ID,
	)

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildBuildErroneousMockRowFromOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

func buildErroneousMockRowFromOAuth2Client(c *v1.OAuth2Client) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(oauth2ClientsTableColumns).AddRow(
		c.ArchivedOn,
		c.Name,
		c.ClientID,
		strings.Join(c.Scopes, scopesSeparator),
		c.RedirectURI,
		c.ClientSecret,
		c.CreatedOn,
		c.LastUpdatedOn,
		c.BelongsToUser,
		c.ID,
	)

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildBuildErroneousMockRowFromOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

func buildErroneousMockRowFromOAuth2Client(c *v1.OAuth2Client) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(oauth2ClientsTableColumns).AddRow(
		c.ArchivedOn,
		c.Name,
		c.ClientID,
		strings.Join(c.Scopes, scopesSeparator),
		c.RedirectURI,
		c.ClientSecret,
		c.CreatedOn,
		c.LastUpdatedOn,
		c.BelongsToUser,
		c.ID,
	)

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestScanOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestPostgres_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanOAuth2Clients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanOAuth2Clients(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestScanOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestSqlite_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := s.scanOAuth2Clients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := s.scanOAuth2Clients(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestScanOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestMariaDB_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := m.scanOAuth2Clients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := m.scanOAuth2Clients(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientByClientIDQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = $1"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
		}
		actualQuery, actualArgs := p.buildGetOAuth2ClientByClientIDQuery(exampleOAuth2Client.ClientID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientByClientIDQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
		}
		actualQuery, actualArgs := s.buildGetOAuth2ClientByClientIDQuery(exampleOAuth2Client.ClientID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientByClientIDQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
		}
		actualQuery, actualArgs := m.buildGetOAuth2ClientByClientIDQuery(exampleOAuth2Client.ClientID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2ClientByClientID(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous row", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2ClientByClientID(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous row", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2ClientByClientID(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.client_id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous row", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.ClientID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetAllOAuth2ClientsQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := p.buildGetAllOAuth2ClientsQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetAllOAuth2ClientsQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := s.buildGetAllOAuth2ClientsQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetAllOAuth2ClientsQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_buildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := m.buildGetAllOAuth2ClientsQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetAllOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := []*v1.OAuth2Client{
			exampleOAuth2Client,
			exampleOAuth2Client,
			exampleOAuth2Client,
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				exampleOAuth2Client,
				exampleOAuth2Client,
				exampleOAuth2Client,
			),
		)

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetAllOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := []*v1.OAuth2Client{
			exampleOAuth2Client,
			exampleOAuth2Client,
			exampleOAuth2Client,
		}

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				exampleOAuth2Client,
				exampleOAuth2Client,
				exampleOAuth2Client,
			),
		)

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing query", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetAllOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expected := []*v1.OAuth2Client{
			exampleOAuth2Client,
			exampleOAuth2Client,
			exampleOAuth2Client,
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				exampleOAuth2Client,
				exampleOAuth2Client,
				exampleOAuth2Client,
			),
		)

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing query", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetAllOAuth2Clients(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetAllOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 ORDER BY oauth2_clients.id"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		expected := []*v1.OAuth2Client{
			&exampleOAuth2ClientList.Clients[0],
			&exampleOAuth2ClientList.Clients[1],
			&exampleOAuth2ClientList.Clients[2],
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildMockRowsFromOAuth2Client(expected...))

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetAllOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		expected := []*v1.OAuth2Client{
			&exampleOAuth2ClientList.Clients[0],
			&exampleOAuth2ClientList.Clients[1],
			&exampleOAuth2ClientList.Clients[2],
		}

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildMockRowsFromOAuth2Client(expected...))

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetAllOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		expected := []*v1.OAuth2Client{
			&exampleOAuth2ClientList.Clients[0],
			&exampleOAuth2ClientList.Clients[1],
			&exampleOAuth2ClientList.Clients[2],
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildMockRowsFromOAuth2Client(expected...))

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetAllOAuth2ClientsForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 AND oauth2_clients.id = $2"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := p.buildGetOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := s.buildGetOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := m.buildGetOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 AND oauth2_clients.id = $2"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := p.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := s.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildMockRowsFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(exampleOAuth2Client))

		actual, err := m.GetOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetAllOAuth2ClientsCountQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := p.buildGetAllOAuth2ClientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetAllOAuth2ClientsCountQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := s.buildGetAllOAuth2ClientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetAllOAuth2ClientsCountQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_buildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		actualQuery := m.buildGetAllOAuth2ClientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildTestDB_GetAllOAuth2ClientCount(dbvendor)

		expected := `
package example

import (
	"context"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildTestDB_GetAllOAuth2ClientCount(dbvendor)

		expected := `
package example

import (
	"context"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		expectedCount := uint64(123)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildTestDB_GetAllOAuth2ClientCount(dbvendor)

		expected := `
package example

import (
	"context"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_GetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(oauth2_clients.id) FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL"
		expectedCount := uint64(123)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllOAuth2ClientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientsForUserQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 AND oauth2_clients.created_on > $2 AND oauth2_clients.created_on < $3 AND oauth2_clients.last_updated_on > $4 AND oauth2_clients.last_updated_on < $5 ORDER BY oauth2_clients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetOAuth2ClientsForUserQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientsForUserQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.created_on > ? AND oauth2_clients.created_on < ? AND oauth2_clients.last_updated_on > ? AND oauth2_clients.last_updated_on < ? ORDER BY oauth2_clients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := s.buildGetOAuth2ClientsForUserQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildGetOAuth2ClientsForUserQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? AND oauth2_clients.created_on > ? AND oauth2_clients.created_on < ? AND oauth2_clients.last_updated_on > ? AND oauth2_clients.last_updated_on < ? ORDER BY oauth2_clients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := m.buildGetOAuth2ClientsForUserQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = $1 ORDER BY oauth2_clients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				&exampleOAuth2ClientList.Clients[0],
				&exampleOAuth2ClientList.Clients[1],
				&exampleOAuth2ClientList.Clients[2],
			),
		)

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error reading from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(fake.BuildFakeOAuth2Client()))

		actual, err := p.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				&exampleOAuth2ClientList.Clients[0],
				&exampleOAuth2ClientList.Clients[1],
				&exampleOAuth2ClientList.Clients[2],
			),
		)

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error reading from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(fake.BuildFakeOAuth2Client()))

		actual, err := s.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_GetOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT oauth2_clients.id, oauth2_clients.name, oauth2_clients.client_id, oauth2_clients.scopes, oauth2_clients.redirect_uri, oauth2_clients.client_secret, oauth2_clients.created_on, oauth2_clients.last_updated_on, oauth2_clients.archived_on, oauth2_clients.belongs_to_user FROM oauth2_clients WHERE oauth2_clients.archived_on IS NULL AND oauth2_clients.belongs_to_user = ? ORDER BY oauth2_clients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		exampleOAuth2ClientList := fake.BuildFakeOAuth2ClientList()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(
			buildMockRowsFromOAuth2Client(
				&exampleOAuth2ClientList.Clients[0],
				&exampleOAuth2ClientList.Clients[1],
				&exampleOAuth2ClientList.Clients[2],
			),
		)

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error reading from database", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		filter := v1.DefaultQueryFilter()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromOAuth2Client(fake.BuildFakeOAuth2Client()))

		actual, err := m.GetOAuth2ClientsForUser(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildCreateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestPostgres_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildCreateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestSqlite_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"
		expectedArgs := []interface{}{
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		}
		actualQuery, actualArgs := s.buildCreateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildCreateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestMariaDB_buildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"
		expectedArgs := []interface{}{
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		}
		actualQuery, actualArgs := m.buildCreateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_CreateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestPostgres_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		exampleRows := gosqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleOAuth2Client.ID, exampleOAuth2Client.CreatedOn)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnRows(exampleRows)

		actual, err := p.CreateOAuth2Client(ctx, expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnError(errors.New("blah"))

		actual, err := p.CreateOAuth2Client(ctx, expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_CreateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestSqlite_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		exampleRows := gosqlmock.NewResult(int64(exampleOAuth2Client.ID), 1)

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnResult(exampleRows)

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleOAuth2Client.CreatedOn)
		s.timeTeller = mtt

		actual, err := s.CreateOAuth2Client(ctx, expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnError(errors.New("blah"))

		actual, err := s.CreateOAuth2Client(ctx, expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_CreateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestMariaDB_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "INSERT INTO oauth2_clients (name,client_id,client_secret,scopes,redirect_uri,belongs_to_user) VALUES (?,?,?,?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)
		exampleRows := gosqlmock.NewResult(int64(exampleOAuth2Client.ID), 1)

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnResult(exampleRows)

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleOAuth2Client.CreatedOn)
		m.timeTeller = mtt

		actual, err := m.CreateOAuth2Client(ctx, expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()
		expectedInput := fake.BuildFakeOAuth2ClientCreationInputFromClient(exampleOAuth2Client)

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(
			exampleOAuth2Client.Name,
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
		).WillReturnError(errors.New("blah"))

		actual, err := m.CreateOAuth2Client(ctx, expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildUpdateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestPostgres_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET client_id = $1, client_secret = $2, scopes = $3, redirect_uri = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $5 AND id = $6 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := p.buildUpdateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildUpdateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestSqlite_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := s.buildUpdateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildUpdateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"strings"
	"testing"
)

func TestMariaDB_buildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.ClientID,
			exampleOAuth2Client.ClientSecret,
			strings.Join(exampleOAuth2Client.Scopes, scopesSeparator),
			exampleOAuth2Client.RedirectURI,
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := m.buildUpdateOAuth2ClientQuery(exampleOAuth2Client)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_UpdateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func TestPostgres_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET client_id = $1, client_secret = $2, scopes = $3, redirect_uri = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $5 AND id = $6 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"last_updated_on"}).AddRow(time.Now().Unix()))

		err := p.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		err := p.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_UpdateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		err := s.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_UpdateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_UpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET client_id = ?, client_secret = ?, scopes = ?, redirect_uri = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		err := m.UpdateOAuth2Client(ctx, exampleOAuth2Client)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildArchiveOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := p.buildArchiveOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildArchiveOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := s.buildArchiveOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_buildArchiveOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		expectedQuery := "UPDATE oauth2_clients SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleOAuth2Client.BelongsToUser,
			exampleOAuth2Client.ID,
		}
		actualQuery, actualArgs := m.buildArchiveOAuth2ClientQuery(exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDB_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_ArchiveOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := p.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(errors.New("blah"))

		err := p.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_ArchiveOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(errors.New("blah"))

		err := s.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildTestDB_ArchiveOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE oauth2_clients SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		exampleOAuth2Client := fake.BuildFakeOAuth2Client()

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleOAuth2Client.BelongsToUser, exampleOAuth2Client.ID).
			WillReturnError(errors.New("blah"))

		err := m.ArchiveOAuth2Client(ctx, exampleOAuth2Client.ID, exampleOAuth2Client.BelongsToUser)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
