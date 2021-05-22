package querybuilding

import (
	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_iterablesTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres todo", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func buildMockRowsFromItems(items ...*v1.Item) *gosqlmock.Rows {
	columns := itemsTableColumns

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range items {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Details,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromItem(x *v1.Item) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(itemsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Details,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanItems(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildItemExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = $1 AND items.id = $2 )"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildItemExistsQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ItemExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = $1 AND items.id = $2 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = $1 AND items.id = $2"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildGetItemQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetItem(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = $1 AND items.id = $2"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := p.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		actualQuery := p.buildGetAllItemsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllItemsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > $1 AND items.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfItemsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllItems(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
	expectedGetQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > $1 AND items.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleItemList := fake.BuildFakeItemList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		out := make(chan []v1.Item)
		doneChan := make(chan bool, 1)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleItem := fake.BuildFakeItem()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1 AND items.created_on > $2 AND items.created_on < $3 AND items.last_updated_on > $4 AND items.last_updated_on < $5 ORDER BY items.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetItemsQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1 ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetItemsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}
		exampleIDsAsStrings := joinUint64s(exampleIDs)

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", exampleIDsAsStrings, defaultLimit)
		expectedArgs := []interface{}{
			exampleUser.ID,
		}
		actualQuery, actualArgs := p.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateItemQuery(exampleItem)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateItem(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		exampleRows := gosqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleItem.ID, exampleItem.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateItem(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET name = $1, details = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildUpdateItemQuery(exampleItem)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET name = $1, details = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		exampleRows := gosqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateItem(ctx, exampleItem)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildArchiveItemQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := p.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := p.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres every type", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func buildMockRowsFromEveryTypes(everyTypes ...*v1.EveryType) *gosqlmock.Rows {
	includeCount := len(everyTypes) > 1
	columns := everyTypesTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range everyTypes {
		rowValues := []driver.Value{
			x.ID,
			x.String,
			x.PointerToString,
			x.Bool,
			x.PointerToBool,
			x.Int,
			x.PointerToInt,
			x.Int8,
			x.PointerToInt8,
			x.Int16,
			x.PointerToInt16,
			x.Int32,
			x.PointerToInt32,
			x.Int64,
			x.PointerToInt64,
			x.Uint,
			x.PointerToUint,
			x.Uint8,
			x.PointerToUint8,
			x.Uint16,
			x.PointerToUint16,
			x.Uint32,
			x.PointerToUint32,
			x.Uint64,
			x.PointerToUint64,
			x.Float32,
			x.PointerToFloat32,
			x.Float64,
			x.PointerToFloat64,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(everyTypes))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromEveryType(x *v1.EveryType) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(everyTypesTableColumns).AddRow(
		x.ArchivedOn,
		x.String,
		x.PointerToString,
		x.Bool,
		x.PointerToBool,
		x.Int,
		x.PointerToInt,
		x.Int8,
		x.PointerToInt8,
		x.Int16,
		x.PointerToInt16,
		x.Int32,
		x.PointerToInt32,
		x.Int64,
		x.PointerToInt64,
		x.Uint,
		x.PointerToUint,
		x.Uint8,
		x.PointerToUint8,
		x.Uint16,
		x.PointerToUint16,
		x.Uint32,
		x.PointerToUint32,
		x.Uint64,
		x.PointerToUint64,
		x.Float32,
		x.PointerToFloat32,
		x.Float64,
		x.PointerToFloat64,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanEveryTypes(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanEveryTypes(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanEveryTypes(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildEveryTypeExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "SELECT EXISTS ( SELECT every_types.id FROM every_types WHERE every_types.id = $1 )"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := p.buildEveryTypeExistsQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_EveryTypeExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT every_types.id FROM every_types WHERE every_types.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.EveryTypeExists(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.EveryTypeExists(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id = $1"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := p.buildGetEveryTypeQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnRows(buildMockRowsFromEveryTypes(exampleEveryType))

		actual, err := p.GetEveryType(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleEveryType, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllEveryTypesCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
		actualQuery := p.buildGetAllEveryTypesCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllEveryTypesCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllEveryTypesCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfEveryTypesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id > $1 AND every_types.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfEveryTypesQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllEveryTypes(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
	expectedGetQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id > $1 AND every_types.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		out := make(chan []v1.EveryType)
		doneChan := make(chan bool, 1)

		err := p.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.EveryType)

		err := p.GetAllEveryTypes(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.EveryType)

		err := p.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.EveryType)

		err := p.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleEveryType := fake.BuildFakeEveryType()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		out := make(chan []v1.EveryType)

		err := p.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetEveryTypesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL AND every_types.created_on > $1 AND every_types.created_on < $2 AND every_types.last_updated_on > $3 AND every_types.last_updated_on < $4 ORDER BY every_types.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetEveryTypesQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetEveryTypes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL ORDER BY every_types.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := p.GetEveryTypes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := p.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetEveryTypesWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}
		exampleIDsAsStrings := joinUint64s(exampleIDs)

		expectedQuery := fmt.Sprintf("SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM (SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS every_types WHERE every_types.archived_on IS NULL", exampleIDsAsStrings, defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetEveryTypesWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM (SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS every_types WHERE every_types.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := p.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList.EveryTypes, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM (SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS every_types WHERE every_types.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM (SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS every_types WHERE every_types.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM (SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS every_types WHERE every_types.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := p.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "INSERT INTO every_types (string,pointer_to_string,bool,pointer_to_bool,int,pointer_to_int,int8,pointer_to_int8,int16,pointer_to_int16,int32,pointer_to_int32,int64,pointer_to_int64,uint,pointer_to_uint,uint8,pointer_to_uint8,uint16,pointer_to_uint16,uint32,pointer_to_uint32,uint64,pointer_to_uint64,float32,pointer_to_float32,float64,pointer_to_float64) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleEveryType.String,
			exampleEveryType.PointerToString,
			exampleEveryType.Bool,
			exampleEveryType.PointerToBool,
			exampleEveryType.Int,
			exampleEveryType.PointerToInt,
			exampleEveryType.Int8,
			exampleEveryType.PointerToInt8,
			exampleEveryType.Int16,
			exampleEveryType.PointerToInt16,
			exampleEveryType.Int32,
			exampleEveryType.PointerToInt32,
			exampleEveryType.Int64,
			exampleEveryType.PointerToInt64,
			exampleEveryType.Uint,
			exampleEveryType.PointerToUint,
			exampleEveryType.Uint8,
			exampleEveryType.PointerToUint8,
			exampleEveryType.Uint16,
			exampleEveryType.PointerToUint16,
			exampleEveryType.Uint32,
			exampleEveryType.PointerToUint32,
			exampleEveryType.Uint64,
			exampleEveryType.PointerToUint64,
			exampleEveryType.Float32,
			exampleEveryType.PointerToFloat32,
			exampleEveryType.Float64,
			exampleEveryType.PointerToFloat64,
		}
		actualQuery, actualArgs := p.buildCreateEveryTypeQuery(exampleEveryType)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateEveryType(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO every_types (string,pointer_to_string,bool,pointer_to_bool,int,pointer_to_int,int8,pointer_to_int8,int16,pointer_to_int16,int32,pointer_to_int32,int64,pointer_to_int64,uint,pointer_to_uint,uint8,pointer_to_uint8,uint16,pointer_to_uint16,uint32,pointer_to_uint32,uint64,pointer_to_uint64,float32,pointer_to_float32,float64,pointer_to_float64) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()
		exampleInput := fake.BuildFakeEveryTypeCreationInputFromEveryType(exampleEveryType)

		exampleRows := gosqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleEveryType.ID, exampleEveryType.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateEveryType(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleEveryType, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()
		exampleInput := fake.BuildFakeEveryTypeCreationInputFromEveryType(exampleEveryType)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateEveryType(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "UPDATE every_types SET string = $1, pointer_to_string = $2, bool = $3, pointer_to_bool = $4, int = $5, pointer_to_int = $6, int8 = $7, pointer_to_int8 = $8, int16 = $9, pointer_to_int16 = $10, int32 = $11, pointer_to_int32 = $12, int64 = $13, pointer_to_int64 = $14, uint = $15, pointer_to_uint = $16, uint8 = $17, pointer_to_uint8 = $18, uint16 = $19, pointer_to_uint16 = $20, uint32 = $21, pointer_to_uint32 = $22, uint64 = $23, pointer_to_uint64 = $24, float32 = $25, pointer_to_float32 = $26, float64 = $27, pointer_to_float64 = $28, last_updated_on = extract(epoch FROM NOW()) WHERE id = $29 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleEveryType.String,
			exampleEveryType.PointerToString,
			exampleEveryType.Bool,
			exampleEveryType.PointerToBool,
			exampleEveryType.Int,
			exampleEveryType.PointerToInt,
			exampleEveryType.Int8,
			exampleEveryType.PointerToInt8,
			exampleEveryType.Int16,
			exampleEveryType.PointerToInt16,
			exampleEveryType.Int32,
			exampleEveryType.PointerToInt32,
			exampleEveryType.Int64,
			exampleEveryType.PointerToInt64,
			exampleEveryType.Uint,
			exampleEveryType.PointerToUint,
			exampleEveryType.Uint8,
			exampleEveryType.PointerToUint8,
			exampleEveryType.Uint16,
			exampleEveryType.PointerToUint16,
			exampleEveryType.Uint32,
			exampleEveryType.PointerToUint32,
			exampleEveryType.Uint64,
			exampleEveryType.PointerToUint64,
			exampleEveryType.Float32,
			exampleEveryType.PointerToFloat32,
			exampleEveryType.Float64,
			exampleEveryType.PointerToFloat64,
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := p.buildUpdateEveryTypeQuery(exampleEveryType)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE every_types SET string = $1, pointer_to_string = $2, bool = $3, pointer_to_bool = $4, int = $5, pointer_to_int = $6, int8 = $7, pointer_to_int8 = $8, int16 = $9, pointer_to_int16 = $10, int32 = $11, pointer_to_int32 = $12, int64 = $13, pointer_to_int64 = $14, uint = $15, pointer_to_uint = $16, uint8 = $17, pointer_to_uint8 = $18, uint16 = $19, pointer_to_uint16 = $20, uint32 = $21, pointer_to_uint32 = $22, uint64 = $23, pointer_to_uint64 = $24, float32 = $25, pointer_to_float32 = $26, float64 = $27, pointer_to_float64 = $28, last_updated_on = extract(epoch FROM NOW()) WHERE id = $29 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		exampleRows := gosqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
				exampleEveryType.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateEveryType(ctx, exampleEveryType)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
				exampleEveryType.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateEveryType(ctx, exampleEveryType)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "UPDATE every_types SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := p.buildArchiveEveryTypeQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE every_types SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := p.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := p.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite todo", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, dbvendor, typ)

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
	"testing"
	"time"
)

func buildMockRowsFromItems(items ...*v1.Item) *gosqlmock.Rows {
	columns := itemsTableColumns

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range items {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Details,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromItem(x *v1.Item) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(itemsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Details,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}

func TestSqlite_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := s.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := s.scanItems(mockRows)
		assert.NoError(t, err)
	})
}

func TestSqlite_buildItemExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildItemExistsQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_ItemExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := s.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildGetItemQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetItem(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := s.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		actualQuery := s.buildGetAllItemsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		expectedCount := uint64(123)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllItemsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetBatchOfItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := s.buildGetBatchOfItemsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetAllItems(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
	expectedGetQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		exampleItemList := fake.BuildFakeItemList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		out := make(chan []v1.Item)
		doneChan := make(chan bool, 1)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		exampleItem := fake.BuildFakeItem()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.created_on > ? AND items.created_on < ? AND items.last_updated_on > ? AND items.last_updated_on < ? ORDER BY items.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := s.buildGetItemsQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetItemsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.id IN (?,?,?) ORDER BY CASE items.id WHEN 789 THEN 0 WHEN 123 THEN 1 WHEN 456 THEN 2 END LIMIT 20"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
		}
		actualQuery, actualArgs := s.buildCreateItemQuery(exampleItem)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_CreateItem(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnResult(gosqlmock.NewResult(int64(exampleItem.ID), 1))

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleItem.CreatedOn)
		s.timeTeller = mtt

		actual, err := s.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := s.CreateItem(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildUpdateItemQuery(exampleItem)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_UpdateItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		exampleRows := gosqlmock.NewResult(int64(exampleItem.ID), 1)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnResult(exampleRows)

		err := s.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := s.UpdateItem(ctx, exampleItem)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildArchiveItemQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_ArchiveItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := s.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := s.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite every type", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, dbvendor, typ)

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
	"testing"
	"time"
)

func buildMockRowsFromEveryTypes(everyTypes ...*v1.EveryType) *gosqlmock.Rows {
	includeCount := len(everyTypes) > 1
	columns := everyTypesTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range everyTypes {
		rowValues := []driver.Value{
			x.ID,
			x.String,
			x.PointerToString,
			x.Bool,
			x.PointerToBool,
			x.Int,
			x.PointerToInt,
			x.Int8,
			x.PointerToInt8,
			x.Int16,
			x.PointerToInt16,
			x.Int32,
			x.PointerToInt32,
			x.Int64,
			x.PointerToInt64,
			x.Uint,
			x.PointerToUint,
			x.Uint8,
			x.PointerToUint8,
			x.Uint16,
			x.PointerToUint16,
			x.Uint32,
			x.PointerToUint32,
			x.Uint64,
			x.PointerToUint64,
			x.Float32,
			x.PointerToFloat32,
			x.Float64,
			x.PointerToFloat64,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(everyTypes))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromEveryType(x *v1.EveryType) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(everyTypesTableColumns).AddRow(
		x.ArchivedOn,
		x.String,
		x.PointerToString,
		x.Bool,
		x.PointerToBool,
		x.Int,
		x.PointerToInt,
		x.Int8,
		x.PointerToInt8,
		x.Int16,
		x.PointerToInt16,
		x.Int32,
		x.PointerToInt32,
		x.Int64,
		x.PointerToInt64,
		x.Uint,
		x.PointerToUint,
		x.Uint8,
		x.PointerToUint8,
		x.Uint16,
		x.PointerToUint16,
		x.Uint32,
		x.PointerToUint32,
		x.Uint64,
		x.PointerToUint64,
		x.Float32,
		x.PointerToFloat32,
		x.Float64,
		x.PointerToFloat64,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestSqlite_ScanEveryTypes(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := s.scanEveryTypes(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := s.scanEveryTypes(mockRows)
		assert.NoError(t, err)
	})
}

func TestSqlite_buildEveryTypeExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "SELECT EXISTS ( SELECT every_types.id FROM every_types WHERE every_types.id = ? )"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := s.buildEveryTypeExistsQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_EveryTypeExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT every_types.id FROM every_types WHERE every_types.id = ? )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := s.EveryTypeExists(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.EveryTypeExists(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id = ?"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := s.buildGetEveryTypeQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnRows(buildMockRowsFromEveryTypes(exampleEveryType))

		actual, err := s.GetEveryType(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleEveryType, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllEveryTypesCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
		actualQuery := s.buildGetAllEveryTypesCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllEveryTypesCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
		expectedCount := uint64(123)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllEveryTypesCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetBatchOfEveryTypesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id > ? AND every_types.id < ?"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := s.buildGetBatchOfEveryTypesQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetAllEveryTypes(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
	expectedGetQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id > ? AND every_types.id < ?"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		out := make(chan []v1.EveryType)
		doneChan := make(chan bool, 1)

		err := s.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.EveryType)

		err := s.GetAllEveryTypes(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.EveryType)

		err := s.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.EveryType)

		err := s.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		exampleEveryType := fake.BuildFakeEveryType()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		out := make(chan []v1.EveryType)

		err := s.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetEveryTypesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL AND every_types.created_on > ? AND every_types.created_on < ? AND every_types.last_updated_on > ? AND every_types.last_updated_on < ? ORDER BY every_types.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := s.buildGetEveryTypesQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetEveryTypes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL ORDER BY every_types.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := s.GetEveryTypes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := s.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetEveryTypesWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.archived_on IS NULL AND every_types.id IN (?,?,?) ORDER BY CASE every_types.id WHEN 789 THEN 0 WHEN 123 THEN 1 WHEN 456 THEN 2 END LIMIT 20"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := s.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_GetEveryTypesWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := s.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := s.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList.EveryTypes, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := s.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := s.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := s.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := s.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreateEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "INSERT INTO every_types (string,pointer_to_string,bool,pointer_to_bool,int,pointer_to_int,int8,pointer_to_int8,int16,pointer_to_int16,int32,pointer_to_int32,int64,pointer_to_int64,uint,pointer_to_uint,uint8,pointer_to_uint8,uint16,pointer_to_uint16,uint32,pointer_to_uint32,uint64,pointer_to_uint64,float32,pointer_to_float32,float64,pointer_to_float64) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		expectedArgs := []interface{}{
			exampleEveryType.String,
			exampleEveryType.PointerToString,
			exampleEveryType.Bool,
			exampleEveryType.PointerToBool,
			exampleEveryType.Int,
			exampleEveryType.PointerToInt,
			exampleEveryType.Int8,
			exampleEveryType.PointerToInt8,
			exampleEveryType.Int16,
			exampleEveryType.PointerToInt16,
			exampleEveryType.Int32,
			exampleEveryType.PointerToInt32,
			exampleEveryType.Int64,
			exampleEveryType.PointerToInt64,
			exampleEveryType.Uint,
			exampleEveryType.PointerToUint,
			exampleEveryType.Uint8,
			exampleEveryType.PointerToUint8,
			exampleEveryType.Uint16,
			exampleEveryType.PointerToUint16,
			exampleEveryType.Uint32,
			exampleEveryType.PointerToUint32,
			exampleEveryType.Uint64,
			exampleEveryType.PointerToUint64,
			exampleEveryType.Float32,
			exampleEveryType.PointerToFloat32,
			exampleEveryType.Float64,
			exampleEveryType.PointerToFloat64,
		}
		actualQuery, actualArgs := s.buildCreateEveryTypeQuery(exampleEveryType)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_CreateEveryType(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO every_types (string,pointer_to_string,bool,pointer_to_bool,int,pointer_to_int,int8,pointer_to_int8,int16,pointer_to_int16,int32,pointer_to_int32,int64,pointer_to_int64,uint,pointer_to_uint,uint8,pointer_to_uint8,uint16,pointer_to_uint16,uint32,pointer_to_uint32,uint64,pointer_to_uint64,float32,pointer_to_float32,float64,pointer_to_float64) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()
		exampleInput := fake.BuildFakeEveryTypeCreationInputFromEveryType(exampleEveryType)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
			).WillReturnResult(gosqlmock.NewResult(int64(exampleEveryType.ID), 1))

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleEveryType.CreatedOn)
		s.timeTeller = mtt

		actual, err := s.CreateEveryType(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleEveryType, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()
		exampleInput := fake.BuildFakeEveryTypeCreationInputFromEveryType(exampleEveryType)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
			).WillReturnError(errors.New("blah"))

		actual, err := s.CreateEveryType(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdateEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "UPDATE every_types SET string = ?, pointer_to_string = ?, bool = ?, pointer_to_bool = ?, int = ?, pointer_to_int = ?, int8 = ?, pointer_to_int8 = ?, int16 = ?, pointer_to_int16 = ?, int32 = ?, pointer_to_int32 = ?, int64 = ?, pointer_to_int64 = ?, uint = ?, pointer_to_uint = ?, uint8 = ?, pointer_to_uint8 = ?, uint16 = ?, pointer_to_uint16 = ?, uint32 = ?, pointer_to_uint32 = ?, uint64 = ?, pointer_to_uint64 = ?, float32 = ?, pointer_to_float32 = ?, float64 = ?, pointer_to_float64 = ?, last_updated_on = (strftime('%s','now')) WHERE id = ?"
		expectedArgs := []interface{}{
			exampleEveryType.String,
			exampleEveryType.PointerToString,
			exampleEveryType.Bool,
			exampleEveryType.PointerToBool,
			exampleEveryType.Int,
			exampleEveryType.PointerToInt,
			exampleEveryType.Int8,
			exampleEveryType.PointerToInt8,
			exampleEveryType.Int16,
			exampleEveryType.PointerToInt16,
			exampleEveryType.Int32,
			exampleEveryType.PointerToInt32,
			exampleEveryType.Int64,
			exampleEveryType.PointerToInt64,
			exampleEveryType.Uint,
			exampleEveryType.PointerToUint,
			exampleEveryType.Uint8,
			exampleEveryType.PointerToUint8,
			exampleEveryType.Uint16,
			exampleEveryType.PointerToUint16,
			exampleEveryType.Uint32,
			exampleEveryType.PointerToUint32,
			exampleEveryType.Uint64,
			exampleEveryType.PointerToUint64,
			exampleEveryType.Float32,
			exampleEveryType.PointerToFloat32,
			exampleEveryType.Float64,
			exampleEveryType.PointerToFloat64,
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := s.buildUpdateEveryTypeQuery(exampleEveryType)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_UpdateEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE every_types SET string = ?, pointer_to_string = ?, bool = ?, pointer_to_bool = ?, int = ?, pointer_to_int = ?, int8 = ?, pointer_to_int8 = ?, int16 = ?, pointer_to_int16 = ?, int32 = ?, pointer_to_int32 = ?, int64 = ?, pointer_to_int64 = ?, uint = ?, pointer_to_uint = ?, uint8 = ?, pointer_to_uint8 = ?, uint16 = ?, pointer_to_uint16 = ?, uint32 = ?, pointer_to_uint32 = ?, uint64 = ?, pointer_to_uint64 = ?, float32 = ?, pointer_to_float32 = ?, float64 = ?, pointer_to_float64 = ?, last_updated_on = (strftime('%s','now')) WHERE id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		exampleRows := gosqlmock.NewResult(int64(exampleEveryType.ID), 1)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
				exampleEveryType.ID,
			).WillReturnResult(exampleRows)

		err := s.UpdateEveryType(ctx, exampleEveryType)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
				exampleEveryType.ID,
			).WillReturnError(errors.New("blah"))

		err := s.UpdateEveryType(ctx, exampleEveryType)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchiveEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "UPDATE every_types SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND id = ?"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := s.buildArchiveEveryTypeQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSqlite_ArchiveEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE every_types SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := s.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnError(errors.New("blah"))

		err := s.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb todo", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, dbvendor, typ)

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
	"testing"
	"time"
)

func buildMockRowsFromItems(items ...*v1.Item) *gosqlmock.Rows {
	columns := itemsTableColumns

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range items {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Details,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromItem(x *v1.Item) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(itemsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Details,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := m.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := m.scanItems(mockRows)
		assert.NoError(t, err)
	})
}

func TestMariaDB_buildItemExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildItemExistsQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_ItemExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := m.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildGetItemQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetItem(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := m.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		actualQuery := m.buildGetAllItemsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		expectedCount := uint64(123)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllItemsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetBatchOfItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := m.buildGetBatchOfItemsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetAllItems(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
	expectedGetQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		exampleItemList := fake.BuildFakeItemList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		out := make(chan []v1.Item)
		doneChan := make(chan bool, 1)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		exampleItem := fake.BuildFakeItem()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.created_on > ? AND items.created_on < ? AND items.last_updated_on > ? AND items.last_updated_on < ? ORDER BY items.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := m.buildGetItemsQuery(exampleUser.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetItemsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.id IN (?,?,?) ORDER BY CASE items.id WHEN 789 THEN 0 WHEN 123 THEN 1 WHEN 456 THEN 2 END LIMIT 20"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
		}
		actualQuery, actualArgs := m.buildCreateItemQuery(exampleItem)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_CreateItem(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnResult(gosqlmock.NewResult(int64(exampleItem.ID), 1))

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleItem.CreatedOn)
		m.timeTeller = mtt

		actual, err := m.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateItem(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildUpdateItemQuery(exampleItem)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_UpdateItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		exampleRows := gosqlmock.NewResult(int64(exampleItem.ID), 1)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateItem(ctx, exampleItem)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildArchiveItemQuery(exampleItem.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_ArchiveItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := m.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb every type", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]
		x := iterablesTestDotGo(proj, dbvendor, typ)

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
	"testing"
	"time"
)

func buildMockRowsFromEveryTypes(everyTypes ...*v1.EveryType) *gosqlmock.Rows {
	includeCount := len(everyTypes) > 1
	columns := everyTypesTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range everyTypes {
		rowValues := []driver.Value{
			x.ID,
			x.String,
			x.PointerToString,
			x.Bool,
			x.PointerToBool,
			x.Int,
			x.PointerToInt,
			x.Int8,
			x.PointerToInt8,
			x.Int16,
			x.PointerToInt16,
			x.Int32,
			x.PointerToInt32,
			x.Int64,
			x.PointerToInt64,
			x.Uint,
			x.PointerToUint,
			x.Uint8,
			x.PointerToUint8,
			x.Uint16,
			x.PointerToUint16,
			x.Uint32,
			x.PointerToUint32,
			x.Uint64,
			x.PointerToUint64,
			x.Float32,
			x.PointerToFloat32,
			x.Float64,
			x.PointerToFloat64,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(everyTypes))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromEveryType(x *v1.EveryType) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(everyTypesTableColumns).AddRow(
		x.ArchivedOn,
		x.String,
		x.PointerToString,
		x.Bool,
		x.PointerToBool,
		x.Int,
		x.PointerToInt,
		x.Int8,
		x.PointerToInt8,
		x.Int16,
		x.PointerToInt16,
		x.Int32,
		x.PointerToInt32,
		x.Int64,
		x.PointerToInt64,
		x.Uint,
		x.PointerToUint,
		x.Uint8,
		x.PointerToUint8,
		x.Uint16,
		x.PointerToUint16,
		x.Uint32,
		x.PointerToUint32,
		x.Uint64,
		x.PointerToUint64,
		x.Float32,
		x.PointerToFloat32,
		x.Float64,
		x.PointerToFloat64,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_ScanEveryTypes(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := m.scanEveryTypes(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v11.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := m.scanEveryTypes(mockRows)
		assert.NoError(t, err)
	})
}

func TestMariaDB_buildEveryTypeExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "SELECT EXISTS ( SELECT every_types.id FROM every_types WHERE every_types.id = ? )"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := m.buildEveryTypeExistsQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_EveryTypeExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT every_types.id FROM every_types WHERE every_types.id = ? )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := m.EveryTypeExists(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.EveryTypeExists(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id = ?"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := m.buildGetEveryTypeQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnRows(buildMockRowsFromEveryTypes(exampleEveryType))

		actual, err := m.GetEveryType(ctx, exampleEveryType.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleEveryType, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleEveryType := fake.BuildFakeEveryType()

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllEveryTypesCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
		actualQuery := m.buildGetAllEveryTypesCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllEveryTypesCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
		expectedCount := uint64(123)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllEveryTypesCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetBatchOfEveryTypesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id > ? AND every_types.id < ?"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := m.buildGetBatchOfEveryTypesQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetAllEveryTypes(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL"
	expectedGetQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.id > ? AND every_types.id < ?"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		out := make(chan []v1.EveryType)
		doneChan := make(chan bool, 1)

		err := m.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.EveryType)

		err := m.GetAllEveryTypes(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.EveryType)

		err := m.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.EveryType)

		err := m.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		exampleEveryType := fake.BuildFakeEveryType()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		out := make(chan []v1.EveryType)

		err := m.GetAllEveryTypes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetEveryTypesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL AND every_types.created_on > ? AND every_types.created_on < ? AND every_types.last_updated_on > ? AND every_types.last_updated_on < ? ORDER BY every_types.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := m.buildGetEveryTypesQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetEveryTypes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL ORDER BY every_types.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := m.GetEveryTypes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := m.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetEveryTypesWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types WHERE every_types.archived_on IS NULL AND every_types.id IN (?,?,?) ORDER BY CASE every_types.id WHEN 789 THEN 0 WHEN 123 THEN 1 WHEN 456 THEN 2 END LIMIT 20"
		expectedArgs := []interface{}{
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := m.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_GetEveryTypesWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := m.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := m.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList.EveryTypes, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := m.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := m.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()
		var exampleIDs []uint64
		for _, everyType := range exampleEveryTypeList.EveryTypes {
			exampleIDs = append(exampleIDs, everyType.ID)
		}

		expectedQuery, expectedArgs := m.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := m.GetEveryTypesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "INSERT INTO every_types (string,pointer_to_string,bool,pointer_to_bool,int,pointer_to_int,int8,pointer_to_int8,int16,pointer_to_int16,int32,pointer_to_int32,int64,pointer_to_int64,uint,pointer_to_uint,uint8,pointer_to_uint8,uint16,pointer_to_uint16,uint32,pointer_to_uint32,uint64,pointer_to_uint64,float32,pointer_to_float32,float64,pointer_to_float64) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		expectedArgs := []interface{}{
			exampleEveryType.String,
			exampleEveryType.PointerToString,
			exampleEveryType.Bool,
			exampleEveryType.PointerToBool,
			exampleEveryType.Int,
			exampleEveryType.PointerToInt,
			exampleEveryType.Int8,
			exampleEveryType.PointerToInt8,
			exampleEveryType.Int16,
			exampleEveryType.PointerToInt16,
			exampleEveryType.Int32,
			exampleEveryType.PointerToInt32,
			exampleEveryType.Int64,
			exampleEveryType.PointerToInt64,
			exampleEveryType.Uint,
			exampleEveryType.PointerToUint,
			exampleEveryType.Uint8,
			exampleEveryType.PointerToUint8,
			exampleEveryType.Uint16,
			exampleEveryType.PointerToUint16,
			exampleEveryType.Uint32,
			exampleEveryType.PointerToUint32,
			exampleEveryType.Uint64,
			exampleEveryType.PointerToUint64,
			exampleEveryType.Float32,
			exampleEveryType.PointerToFloat32,
			exampleEveryType.Float64,
			exampleEveryType.PointerToFloat64,
		}
		actualQuery, actualArgs := m.buildCreateEveryTypeQuery(exampleEveryType)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_CreateEveryType(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO every_types (string,pointer_to_string,bool,pointer_to_bool,int,pointer_to_int,int8,pointer_to_int8,int16,pointer_to_int16,int32,pointer_to_int32,int64,pointer_to_int64,uint,pointer_to_uint,uint8,pointer_to_uint8,uint16,pointer_to_uint16,uint32,pointer_to_uint32,uint64,pointer_to_uint64,float32,pointer_to_float32,float64,pointer_to_float64) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()
		exampleInput := fake.BuildFakeEveryTypeCreationInputFromEveryType(exampleEveryType)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
			).WillReturnResult(gosqlmock.NewResult(int64(exampleEveryType.ID), 1))

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleEveryType.CreatedOn)
		m.timeTeller = mtt

		actual, err := m.CreateEveryType(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleEveryType, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()
		exampleInput := fake.BuildFakeEveryTypeCreationInputFromEveryType(exampleEveryType)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateEveryType(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "UPDATE every_types SET string = ?, pointer_to_string = ?, bool = ?, pointer_to_bool = ?, int = ?, pointer_to_int = ?, int8 = ?, pointer_to_int8 = ?, int16 = ?, pointer_to_int16 = ?, int32 = ?, pointer_to_int32 = ?, int64 = ?, pointer_to_int64 = ?, uint = ?, pointer_to_uint = ?, uint8 = ?, pointer_to_uint8 = ?, uint16 = ?, pointer_to_uint16 = ?, uint32 = ?, pointer_to_uint32 = ?, uint64 = ?, pointer_to_uint64 = ?, float32 = ?, pointer_to_float32 = ?, float64 = ?, pointer_to_float64 = ?, last_updated_on = UNIX_TIMESTAMP() WHERE id = ?"
		expectedArgs := []interface{}{
			exampleEveryType.String,
			exampleEveryType.PointerToString,
			exampleEveryType.Bool,
			exampleEveryType.PointerToBool,
			exampleEveryType.Int,
			exampleEveryType.PointerToInt,
			exampleEveryType.Int8,
			exampleEveryType.PointerToInt8,
			exampleEveryType.Int16,
			exampleEveryType.PointerToInt16,
			exampleEveryType.Int32,
			exampleEveryType.PointerToInt32,
			exampleEveryType.Int64,
			exampleEveryType.PointerToInt64,
			exampleEveryType.Uint,
			exampleEveryType.PointerToUint,
			exampleEveryType.Uint8,
			exampleEveryType.PointerToUint8,
			exampleEveryType.Uint16,
			exampleEveryType.PointerToUint16,
			exampleEveryType.Uint32,
			exampleEveryType.PointerToUint32,
			exampleEveryType.Uint64,
			exampleEveryType.PointerToUint64,
			exampleEveryType.Float32,
			exampleEveryType.PointerToFloat32,
			exampleEveryType.Float64,
			exampleEveryType.PointerToFloat64,
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := m.buildUpdateEveryTypeQuery(exampleEveryType)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_UpdateEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE every_types SET string = ?, pointer_to_string = ?, bool = ?, pointer_to_bool = ?, int = ?, pointer_to_int = ?, int8 = ?, pointer_to_int8 = ?, int16 = ?, pointer_to_int16 = ?, int32 = ?, pointer_to_int32 = ?, int64 = ?, pointer_to_int64 = ?, uint = ?, pointer_to_uint = ?, uint8 = ?, pointer_to_uint8 = ?, uint16 = ?, pointer_to_uint16 = ?, uint32 = ?, pointer_to_uint32 = ?, uint64 = ?, pointer_to_uint64 = ?, float32 = ?, pointer_to_float32 = ?, float64 = ?, pointer_to_float64 = ?, last_updated_on = UNIX_TIMESTAMP() WHERE id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		exampleRows := gosqlmock.NewResult(int64(exampleEveryType.ID), 1)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
				exampleEveryType.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateEveryType(ctx, exampleEveryType)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.String,
				exampleEveryType.PointerToString,
				exampleEveryType.Bool,
				exampleEveryType.PointerToBool,
				exampleEveryType.Int,
				exampleEveryType.PointerToInt,
				exampleEveryType.Int8,
				exampleEveryType.PointerToInt8,
				exampleEveryType.Int16,
				exampleEveryType.PointerToInt16,
				exampleEveryType.Int32,
				exampleEveryType.PointerToInt32,
				exampleEveryType.Int64,
				exampleEveryType.PointerToInt64,
				exampleEveryType.Uint,
				exampleEveryType.PointerToUint,
				exampleEveryType.Uint8,
				exampleEveryType.PointerToUint8,
				exampleEveryType.Uint16,
				exampleEveryType.PointerToUint16,
				exampleEveryType.Uint32,
				exampleEveryType.PointerToUint32,
				exampleEveryType.Uint64,
				exampleEveryType.PointerToUint64,
				exampleEveryType.Float32,
				exampleEveryType.PointerToFloat32,
				exampleEveryType.Float64,
				exampleEveryType.PointerToFloat64,
				exampleEveryType.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateEveryType(ctx, exampleEveryType)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveEveryTypeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		expectedQuery := "UPDATE every_types SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND id = ?"
		expectedArgs := []interface{}{
			exampleEveryType.ID,
		}
		actualQuery, actualArgs := m.buildArchiveEveryTypeQuery(exampleEveryType.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestMariaDB_ArchiveEveryType(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE every_types SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := m.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleEveryType.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveEveryType(ctx, exampleEveryType.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildMockRowsFromSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildMockRowsFromSomething(proj, typ)

		expected := `
package example

import (
	"database/sql/driver"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

func buildMockRowsFromItems(items ...*v1.Item) *gosqlmock.Rows {
	columns := itemsTableColumns

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range items {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Details,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with enumeration", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.IsEnumeration = true

		x := buildBuildMockRowsFromSomething(proj, typ)

		expected := `
package example

import (
	"database/sql/driver"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

func buildMockRowsFromItems(items ...*v1.Item) *gosqlmock.Rows {
	includeCount := len(items) > 1
	columns := itemsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := gosqlmock.NewRows(columns)

	for _, x := range items {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Details,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		if includeCount {
			rowValues = append(rowValues, len(items))
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

func Test_buildBuildErroneousMockRowFromSomething(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBuildErroneousMockRowFromSomething(proj, typ)

		expected := `
package example

import (
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

func buildErroneousMockRowFromItem(x *v1.Item) *gosqlmock.Rows {
	exampleRows := gosqlmock.NewRows(itemsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Details,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_applyFleshedOutQueryFilter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		qb := squirrel.Select("*")
		x := applyFleshedOutQueryFilter(qb, "example")

		expected := `SELECT * WHERE example.created_on > ? AND example.created_on < ? AND example.last_updated_on > ? AND example.last_updated_on < ? ORDER BY example.id LIMIT 20 OFFSET 180`
		actual, _, err := x.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_applyFleshedOutQueryFilterWithCode(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		qb := squirrel.Select("*")
		x := applyFleshedOutQueryFilterWithCode(qb, "example", nil)

		expected := `SELECT * WHERE example.created_on > ? AND example.created_on < ? AND example.last_updated_on > ? AND example.last_updated_on < ? ORDER BY example.id LIMIT 20 OFFSET 180`
		actual, _, err := x.ToSql()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_appendFleshedOutQueryFilterArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := appendFleshedOutQueryFilterArgs([]jen.Code{})

		expected := `
package main

import ()

func main() {
	exampleFunction(filter.CreatedAfter, filter.CreatedBefore, filter.UpdatedAfter, filter.UpdatedBefore)
}
`
		actual := testutils.RenderCallArgsToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGeneralFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGeneralFields("exampleVarName", typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		exampleVarName.ID,
		exampleVarName.Name,
		exampleVarName.Details,
		exampleVarName.CreatedOn,
		exampleVarName.LastUpdatedOn,
		exampleVarName.ArchivedOn,
		exampleVarName.BelongsToUser,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with chain ownership", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGeneralFields("exampleVarName", proj.LastDataType())

		expected := `
package main

import ()

func main() {
	exampleFunction(
		exampleVarName.ID,
		exampleVarName.CreatedOn,
		exampleVarName.LastUpdatedOn,
		exampleVarName.ArchivedOn,
		exampleVarName.BelongsToAnotherThing,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBadFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildBadFields("exampleVarName", typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		exampleVarName.ArchivedOn,
		exampleVarName.Name,
		exampleVarName.Details,
		exampleVarName.CreatedOn,
		exampleVarName.LastUpdatedOn,
		exampleVarName.BelongsToUser,
		exampleVarName.ID,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with chain ownership", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildBadFields("exampleVarName", proj.LastDataType())

		expected := `
package main

import ()

func main() {
	exampleFunction(
		exampleVarName.ArchivedOn,
		exampleVarName.CreatedOn,
		exampleVarName.LastUpdatedOn,
		exampleVarName.BelongsToAnotherThing,
		exampleVarName.ID,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildPrefixedStringColumns(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expected := []string{"items.id", "items.name", "items.details", "items.created_on", "items.last_updated_on", "items.archived_on", "items.belongs_to_user"}
		actual := buildPrefixedStringColumns(typ)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with chain ownership", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expected := []string{"yet_another_things.id", "yet_another_things.created_on", "yet_another_things.last_updated_on", "yet_another_things.archived_on", "yet_another_things.belongs_to_another_thing"}
		actual := buildPrefixedStringColumns(proj.LastDataType())

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildPrefixedStringColumnsAsString(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expected := `items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user`
		actual := buildPrefixedStringColumnsAsString(typ)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreationStringColumnsAndArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expectedCols := []string{"name", "details", "belongs_to_user"}
		actualCols, args := buildCreationStringColumnsAndArgs(typ)

		expectedRenderedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.Name, exampleItem.Details, exampleItem.BelongsToUser)
}
`
		actualRenderedArgs := testutils.RenderCallArgsToString(t, args)

		assert.Equal(t, expectedCols, actualCols, "expected and actual columns do not match")
		assert.Equal(t, expectedRenderedArgs, actualRenderedArgs, "expected and actual rendered arguments do not match")
	})

	T.Run("with chain ownership", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")

		expectedCols := []string{"belongs_to_another_thing"}
		actualCols, args := buildCreationStringColumnsAndArgs(proj.LastDataType())

		expectedRenderedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleYetAnotherThing.BelongsToAnotherThing)
}
`
		actualRenderedArgs := testutils.RenderCallArgsToString(t, args)

		assert.Equal(t, expectedCols, actualCols, "expected and actual columns do not match")
		assert.Equal(t, expectedRenderedArgs, actualRenderedArgs, "expected and actual rendered arguments do not match")
	})

}

func Test_buildUpdateQueryParts(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expected := []string{"name = $1", "details = $2"}
		actual := buildUpdateQueryParts(dbvendor, typ)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expected := []string{"name = ?", "details = ?"}
		actual := buildUpdateQueryParts(dbvendor, typ)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		expected := []string{"name = ?", "details = ?"}
		actual := buildUpdateQueryParts(dbvendor, typ)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_getIncIndex(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		exampleIndex := uint(1)

		expected := `$2`
		actual := getIncIndex(dbvendor, exampleIndex)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		exampleIndex := uint(1)

		expected := `?`
		actual := getIncIndex(dbvendor, exampleIndex)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		exampleIndex := uint(1)

		expected := `?`
		actual := getIncIndex(dbvendor, exampleIndex)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("invalid dbvendor", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("invalid")
		exampleIndex := uint(1)

		expected := ``
		actual := getIncIndex(dbvendor, exampleIndex)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_getTimeQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		expected := `extract(epoch FROM NOW())`
		actual := getTimeQuery(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		expected := `(strftime('%s','now'))`
		actual := getTimeQuery(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		expected := `UNIX_TIMESTAMP()`
		actual := getTimeQuery(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with invalid db", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("invalid")

		expected := ``
		actual := getTimeQuery(dbvendor)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestScanListOfThings(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestScanListOfThings(proj, dbvendor, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestPostgres_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanItems(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.IsEnumeration = true

		x := buildTestScanListOfThings(proj, dbvendor, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestPostgres_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanItems(mockRows)
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
		typ := proj.DataTypes[0]
		x := buildTestScanListOfThings(proj, dbvendor, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestSqlite_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := s.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := s.scanItems(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.IsEnumeration = true

		x := buildTestScanListOfThings(proj, dbvendor, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestSqlite_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := s.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		s, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := s.scanItems(mockRows)
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
		typ := proj.DataTypes[0]
		x := buildTestScanListOfThings(proj, dbvendor, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestMariaDB_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := m.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := m.scanItems(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with enumeration", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.IsEnumeration = true

		x := buildTestScanListOfThings(proj, dbvendor, typ)

		expected := `
package example

import (
	"errors"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestMariaDB_ScanItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := m.scanItems(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		m, _ := buildTestService(t)
		mockRows := &v1.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := m.scanItems(mockRows)
		assert.NoError(t, err)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBBuildSomethingExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBBuildSomethingExistsQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildItemExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = $1 AND items.id = $2 )"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildItemExistsQuery(exampleItem.ID, exampleUser.ID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBBuildSomethingExistsQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildItemExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildItemExistsQuery(exampleItem.ID, exampleUser.ID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBBuildSomethingExistsQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildItemExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildItemExistsQuery(exampleItem.ID, exampleUser.ID)

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

func Test_buildTestDBSomethingExists(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBSomethingExists(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_ItemExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = $1 AND items.id = $2 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

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
		typ := proj.DataTypes[0]
		x := buildTestDBSomethingExists(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_ItemExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := s.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

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
		typ := proj.DataTypes[0]
		x := buildTestDBSomethingExists(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_ItemExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT items.id FROM items WHERE items.belongs_to_user = ? AND items.id = ? )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := m.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.ItemExists(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildTestDBSomethingExists(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"context"
	"database/sql"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_YetAnotherThingExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT yet_another_things.id FROM yet_another_things JOIN another_things ON yet_another_things.belongs_to_another_thing=another_things.id JOIN things ON another_things.belongs_to_thing=things.id WHERE another_things.belongs_to_thing = $1 AND another_things.id = $2 AND things.id = $3 AND yet_another_things.belongs_to_another_thing = $4 AND yet_another_things.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = exampleThing.ID
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleThing.ID,
				exampleAnotherThing.ID,
				exampleThing.ID,
				exampleAnotherThing.ID,
				exampleYetAnotherThing.ID,
			).
			WillReturnRows(gosqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.YetAnotherThingExists(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = exampleThing.ID
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleThing.ID,
				exampleAnotherThing.ID,
				exampleThing.ID,
				exampleAnotherThing.ID,
				exampleYetAnotherThing.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.YetAnotherThingExists(ctx, exampleThing.ID, exampleAnotherThing.ID, exampleYetAnotherThing.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBBuildGetSomethingQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBBuildGetSomethingQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = $1 AND items.id = $2"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildGetItemQuery(exampleItem.ID, exampleUser.ID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBBuildGetSomethingQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildGetItemQuery(exampleItem.ID, exampleUser.ID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBBuildGetSomethingQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"
		expectedArgs := []interface{}{
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildGetItemQuery(exampleItem.ID, exampleUser.ID)

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

func Test_buildTestDBGetSomething(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetSomething(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetItem(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = $1 AND items.id = $2"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := p.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetSomething(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetItem(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := s.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetSomething(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetItem(T *testing.T) {
	T.Parallel()

	exampleUser := fake.BuildFakeUser()
	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.belongs_to_user = ? AND items.id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := m.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres without belonging to user", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		x := buildTestDBGetSomething(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM items WHERE items.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.ID,
			).
			WillReturnRows(buildMockRowsFromItems(exampleItem))

		actual, err := p.GetItem(ctx, exampleItem.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleItem := fake.BuildFakeItem()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItem(ctx, exampleItem.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBBuildGetAllSomethingCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBBuildGetAllSomethingCountQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		actualQuery := p.buildGetAllItemsCountQuery()

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
		typ := proj.DataTypes[0]
		x := buildTestDBBuildGetAllSomethingCountQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		actualQuery := s.buildGetAllItemsCountQuery()

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
		typ := proj.DataTypes[0]
		x := buildTestDBBuildGetAllSomethingCountQuery(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		actualQuery := m.buildGetAllItemsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBGetAllSomethingCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetAllSomethingCount(dbvendor, typ)

		expected := `
package example

import (
	"context"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllItemsCount(ctx)
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
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetAllSomethingCount(dbvendor, typ)

		expected := `
package example

import (
	"context"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		expectedCount := uint64(123)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllItemsCount(ctx)
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
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetAllSomethingCount(dbvendor, typ)

		expected := `
package example

import (
	"context"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
		expectedCount := uint64(123)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllItemsCount(ctx)
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

func Test_buildTestDBGetBatchOfSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetBatchOfSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_buildGetBatchOfItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > $1 AND items.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfItemsQuery(beginID, endID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetBatchOfSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlite_buildGetBatchOfItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := s.buildGetBatchOfItemsQuery(beginID, endID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetBatchOfSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestMariaDB_buildGetBatchOfItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := m.buildGetBatchOfItemsQuery(beginID, endID)

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

func Test_buildTestDBGetAllOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetAllOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func TestPostgres_GetAllItems(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
	expectedGetQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > $1 AND items.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleItemList := fake.BuildFakeItemList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		out := make(chan []v1.Item)
		doneChan := make(chan bool, 1)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleItem := fake.BuildFakeItem()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		out := make(chan []v1.Item)

		err := p.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetAllOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func TestSqlite_GetAllItems(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
	expectedGetQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		exampleItemList := fake.BuildFakeItemList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		out := make(chan []v1.Item)
		doneChan := make(chan bool, 1)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		s, mockDB := buildTestService(t)
		exampleItem := fake.BuildFakeItem()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		out := make(chan []v1.Item)

		err := s.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetAllOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
	"time"
)

func TestMariaDB_GetAllItems(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(items.id) FROM items WHERE items.archived_on IS NULL"
	expectedGetQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.id > ? AND items.id < ?"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		exampleItemList := fake.BuildFakeItemList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		out := make(chan []v1.Item)
		doneChan := make(chan bool, 1)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		m, mockDB := buildTestService(t)
		exampleItem := fake.BuildFakeItem()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(gosqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		out := make(chan []v1.Item)

		err := m.GetAllItems(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestDBGetListOfSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1 AND items.created_on > $2 AND items.created_on < $3 AND items.last_updated_on > $4 AND items.last_updated_on < $5 ORDER BY items.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetItemsQuery(exampleUser.ID, filter)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.created_on > ? AND items.created_on < ? AND items.last_updated_on > ? AND items.last_updated_on < ? ORDER BY items.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := s.buildGetItemsQuery(exampleUser.ID, filter)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		filter := fake.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.created_on > ? AND items.created_on < ? AND items.last_updated_on > ? AND items.last_updated_on < ? ORDER BY items.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := m.buildGetItemsQuery(exampleUser.ID, filter)

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

func Test_buildTestDBGetListOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

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

func TestPostgres_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1 ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := p.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]

		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

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

func TestPostgres_GetEveryTypes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL ORDER BY every_types.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := p.GetEveryTypes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := p.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres while not belonging to user", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.RestrictedToUser = false
		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM items WHERE items.archived_on IS NULL ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := p.GetItems(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItems(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		actual, err := p.GetItems(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildErroneousMockRowFromItem(exampleItem),
			)

		actual, err := p.GetItems(ctx, filter)
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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

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

func TestSqlite_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := s.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]

		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

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

func TestSqlite_GetEveryTypes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL ORDER BY every_types.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := s.GetEveryTypes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := s.GetEveryTypes(ctx, filter)
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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

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

func TestMariaDB_GetItems(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? ORDER BY items.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleItemList := fake.BuildFakeItemList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
			).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := m.GetItems(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with enumeration", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]

		x := buildTestDBGetListOfSomethingFuncDecl(proj, dbvendor, typ)

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

func TestMariaDB_GetEveryTypes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on, (SELECT COUNT(every_types.id) FROM every_types WHERE every_types.archived_on IS NULL) FROM every_types WHERE every_types.archived_on IS NULL ORDER BY every_types.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryTypeList := fake.BuildFakeEveryTypeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromEveryTypes(
					&exampleEveryTypeList.EveryTypes[0],
					&exampleEveryTypeList.EveryTypes[1],
					&exampleEveryTypeList.EveryTypes[2],
				),
			)

		actual, err := m.GetEveryTypes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleEveryTypeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetEveryTypes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning every type", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)
		filter := v1.DefaultQueryFilter()

		exampleEveryType := fake.BuildFakeEveryType()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromEveryType(exampleEveryType))

		actual, err := m.GetEveryTypes(ctx, filter)
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

func Test_buildTestDBGetListOfSomethingWithIDsQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildGetItemsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}
		exampleIDsAsStrings := joinUint64s(exampleIDs)

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", exampleIDsAsStrings, defaultLimit)
		expectedArgs := []interface{}{
			exampleUser.ID,
		}
		actualQuery, actualArgs := p.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres every type", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgres_buildGetEveryTypesWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}
		exampleIDsAsStrings := joinUint64s(exampleIDs)

		expectedQuery := fmt.Sprintf("SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM (SELECT every_types.id, every_types.string, every_types.pointer_to_string, every_types.bool, every_types.pointer_to_bool, every_types.int, every_types.pointer_to_int, every_types.int8, every_types.pointer_to_int8, every_types.int16, every_types.pointer_to_int16, every_types.int32, every_types.pointer_to_int32, every_types.int64, every_types.pointer_to_int64, every_types.uint, every_types.pointer_to_uint, every_types.uint8, every_types.pointer_to_uint8, every_types.uint16, every_types.pointer_to_uint16, every_types.uint32, every_types.pointer_to_uint32, every_types.uint64, every_types.pointer_to_uint64, every_types.float32, every_types.pointer_to_float32, every_types.float64, every_types.pointer_to_float64, every_types.created_on, every_types.last_updated_on, every_types.archived_on FROM every_types JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS every_types WHERE every_types.archived_on IS NULL", exampleIDsAsStrings, defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetEveryTypesWithIDsQuery(defaultLimit, exampleIDs)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildGetItemsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.id IN (?,?,?) ORDER BY CASE items.id WHEN 789 THEN 0 WHEN 123 THEN 1 WHEN 456 THEN 2 END LIMIT 20"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildGetItemsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := "SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items WHERE items.archived_on IS NULL AND items.belongs_to_user = ? AND items.id IN (?,?,?) ORDER BY CASE items.id WHEN 789 THEN 0 WHEN 123 THEN 1 WHEN 456 THEN 2 END LIMIT 20"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleIDs[0],
			exampleIDs[1],
			exampleIDs[2],
		}
		actualQuery, actualArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

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

func Test_buildTestDBGetListOfSomethingWithIDsFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on, items.belongs_to_user FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL AND items.belongs_to_user = $1", joinUint64s(exampleIDs), defaultLimit)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(exampleUser.ID).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := p.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		s, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := s.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := s.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

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
		typ := proj.DataTypes[0]
		x := buildTestDBGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fake.BuildFakeUser()

		m, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery, expectedArgs := m.buildGetItemsWithIDsQuery(exampleUser.ID, defaultLimit, exampleIDs)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(interfacesToDriverValues(expectedArgs)...).
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := m.GetItemsWithIDs(ctx, exampleUser.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres while not belonging to user", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		x := buildTestDBGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_GetItemsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleItemList := fake.BuildFakeItemList()
		var exampleIDs []uint64
		for _, item := range exampleItemList.Items {
			exampleIDs = append(exampleIDs, item.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromItems(
					&exampleItemList.Items[0],
					&exampleItemList.Items[1],
					&exampleItemList.Items[2],
				),
			)

		actual, err := p.GetItemsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleItemList.Items, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetItemsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetItemsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM (SELECT items.id, items.name, items.details, items.created_on, items.last_updated_on, items.archived_on FROM items JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS items WHERE items.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleItem := fake.BuildFakeItem()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromItem(exampleItem))

		actual, err := p.GetItemsWithIDs(ctx, defaultLimit, exampleIDs)

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

func Test_buildTestDBCreateSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBCreateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateItemQuery(exampleItem)

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
		typ := proj.DataTypes[0]
		x := buildTestDBCreateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
		}
		actualQuery, actualArgs := s.buildCreateItemQuery(exampleItem)

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
		typ := proj.DataTypes[0]
		x := buildTestDBCreateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
		}
		actualQuery, actualArgs := m.buildCreateItemQuery(exampleItem)

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

func Test_buildTestDBCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBCreateSomethingFuncDecl(proj, dbvendor, typ)

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

func TestPostgres_CreateItem(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		exampleRows := gosqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleItem.ID, exampleItem.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateItem(ctx, exampleInput)
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
		typ := proj.DataTypes[0]
		x := buildTestDBCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_CreateItem(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnResult(gosqlmock.NewResult(int64(exampleItem.ID), 1))

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleItem.CreatedOn)
		s.timeTeller = mtt

		actual, err := s.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := s.CreateItem(ctx, exampleInput)
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
		typ := proj.DataTypes[0]
		x := buildTestDBCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_CreateItem(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO items (name,details,belongs_to_user) VALUES (?,?,?)"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnResult(gosqlmock.NewResult(int64(exampleItem.ID), 1))

		mtt := &mockTimeTeller{}
		mtt.On("Now").Return(exampleItem.CreatedOn)
		m.timeTeller = mtt

		actual, err := m.CreateItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleItem, actual)

		mock.AssertExpectationsForObjects(t, mtt)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID
		exampleInput := fake.BuildFakeItemCreationInputFromItem(exampleItem)

		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateItem(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildTestDBCreateSomethingFuncDecl(proj, dbvendor, typ)

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

func TestPostgres_CreateYetAnotherThing(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO yet_another_things (belongs_to_another_thing) VALUES ($1) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = exampleThing.ID
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		exampleRows := gosqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleYetAnotherThing.ID, exampleYetAnotherThing.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleYetAnotherThing.BelongsToAnotherThing,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateYetAnotherThing(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleYetAnotherThing, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleThing := fake.BuildFakeThing()
		exampleAnotherThing := fake.BuildFakeAnotherThing()
		exampleAnotherThing.BelongsToThing = exampleThing.ID
		exampleYetAnotherThing := fake.BuildFakeYetAnotherThing()
		exampleYetAnotherThing.BelongsToAnotherThing = exampleAnotherThing.ID
		exampleInput := fake.BuildFakeYetAnotherThingCreationInputFromYetAnotherThing(exampleYetAnotherThing)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleYetAnotherThing.BelongsToAnotherThing,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateYetAnotherThing(ctx, exampleInput)
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

func Test_buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		qb, args := buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE items SET name = $1, details = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on`
		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)

		expectedRenderedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.Name, exampleItem.Details, exampleItem.BelongsToUser, exampleItem.ID)
}
`
		actualRenderedArgs := testutils.RenderCallArgsToString(t, args)

		assert.Equal(t, expectedQuery, actualQuery, "expected and actual query do not match")
		assert.Equal(t, expectedRenderedArgs, actualRenderedArgs, "expected and actual rendered args do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		qb, args := buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE items SET name = ?, details = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?`
		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)

		expectedRenderedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.Name, exampleItem.Details, exampleItem.BelongsToUser, exampleItem.ID)
}
`
		actualRenderedArgs := testutils.RenderCallArgsToString(t, args)

		assert.Equal(t, expectedQuery, actualQuery, "expected and actual query do not match")
		assert.Equal(t, expectedRenderedArgs, actualRenderedArgs, "expected and actual rendered args do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		qb, args := buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE items SET name = ?, details = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?`
		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)

		expectedRenderedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.Name, exampleItem.Details, exampleItem.BelongsToUser, exampleItem.ID)
}
`
		actualRenderedArgs := testutils.RenderCallArgsToString(t, args)

		assert.Equal(t, expectedQuery, actualQuery, "expected and actual query do not match")
		assert.Equal(t, expectedRenderedArgs, actualRenderedArgs, "expected and actual rendered args do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		qb, args := buildTestBuildUpdateSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE yet_another_things SET last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_another_thing = $1 AND id = $2 RETURNING last_updated_on`
		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)

		expectedRenderedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleYetAnotherThing.BelongsToAnotherThing, exampleYetAnotherThing.ID)
}
`
		actualRenderedArgs := testutils.RenderCallArgsToString(t, args)

		assert.Equal(t, expectedQuery, actualQuery, "expected and actual query do not match")
		assert.Equal(t, expectedRenderedArgs, actualRenderedArgs, "expected and actual rendered args do not match")
	})
}

func Test_buildTestBuildUpdateSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestBuildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET name = $1, details = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildUpdateItemQuery(exampleItem)

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
		typ := proj.DataTypes[0]
		x := buildTestBuildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildUpdateItemQuery(exampleItem)

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
		typ := proj.DataTypes[0]
		x := buildTestBuildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleItem.Name,
			exampleItem.Details,
			exampleItem.BelongsToUser,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildUpdateItemQuery(exampleItem)

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

func Test_buildTestDBUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBUpdateSomethingFuncDecl(proj, dbvendor, typ)

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

func TestPostgres_UpdateItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET name = $1, details = $2, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $3 AND id = $4 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		exampleRows := gosqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateItem(ctx, exampleItem)
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
		typ := proj.DataTypes[0]
		x := buildTestDBUpdateSomethingFuncDecl(proj, dbvendor, typ)

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

func TestSqlite_UpdateItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = (strftime('%s','now')) WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		exampleRows := gosqlmock.NewResult(int64(exampleItem.ID), 1)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnResult(exampleRows)

		err := s.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := s.UpdateItem(ctx, exampleItem)
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
		typ := proj.DataTypes[0]
		x := buildTestDBUpdateSomethingFuncDecl(proj, dbvendor, typ)

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

func TestMariaDB_UpdateItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET name = ?, details = ?, last_updated_on = UNIX_TIMESTAMP() WHERE belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		exampleRows := gosqlmock.NewResult(int64(exampleItem.ID), 1)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateItem(ctx, exampleItem)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleItem.Name,
				exampleItem.Details,
				exampleItem.BelongsToUser,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateItem(ctx, exampleItem)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		qb, returnedExpectedArgs, returnedCallArgs := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE items SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on`
		expectedRenderedExpectedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleUser.ID, exampleItem.ID)
}
`
		expectedRenderedCallArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.ID)
}
`

		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, expectedQuery, actualQuery, "expected and actual output do not match")

		actualRenderedExpectedArgs := testutils.RenderCallArgsToString(t, returnedExpectedArgs)
		assert.Equal(t, expectedRenderedExpectedArgs, actualRenderedExpectedArgs, "expected and actual output do not match")

		actualRenderedCallArgs := testutils.RenderCallArgsToString(t, returnedCallArgs)
		assert.Equal(t, expectedRenderedCallArgs, actualRenderedCallArgs, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		qb, returnedExpectedArgs, returnedCallArgs := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE items SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?`
		expectedRenderedExpectedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleUser.ID, exampleItem.ID)
}
`
		expectedRenderedCallArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.ID)
}
`

		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, expectedQuery, actualQuery, "expected and actual output do not match")

		actualRenderedExpectedArgs := testutils.RenderCallArgsToString(t, returnedExpectedArgs)
		assert.Equal(t, expectedRenderedExpectedArgs, actualRenderedExpectedArgs, "expected and actual output do not match")

		actualRenderedCallArgs := testutils.RenderCallArgsToString(t, returnedCallArgs)
		assert.Equal(t, expectedRenderedCallArgs, actualRenderedCallArgs, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		qb, returnedExpectedArgs, returnedCallArgs := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE items SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?`
		expectedRenderedExpectedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleUser.ID, exampleItem.ID)
}
`
		expectedRenderedCallArgs := `
package main

import ()

func main() {
	exampleFunction(exampleItem.ID)
}
`

		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, expectedQuery, actualQuery, "expected and actual output do not match")

		actualRenderedExpectedArgs := testutils.RenderCallArgsToString(t, returnedExpectedArgs)
		assert.Equal(t, expectedRenderedExpectedArgs, actualRenderedExpectedArgs, "expected and actual output do not match")

		actualRenderedCallArgs := testutils.RenderCallArgsToString(t, returnedCallArgs)
		assert.Equal(t, expectedRenderedCallArgs, actualRenderedCallArgs, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		qb, returnedExpectedArgs, returnedCallArgs := buildTestBuildArchiveSomethingQueryFuncDeclQueryBuilder(dbvendor, typ)

		expectedQuery := `UPDATE yet_another_things SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_another_thing = $1 AND id = $2 RETURNING archived_on`
		expectedRenderedExpectedArgs := `
package main

import ()

func main() {
	exampleFunction(exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`
		expectedRenderedCallArgs := `
package main

import ()

func main() {
	exampleFunction(exampleAnotherThing.ID, exampleYetAnotherThing.ID)
}
`

		actualQuery, _, err := qb.ToSql()
		assert.NoError(t, err)
		assert.Equal(t, expectedQuery, actualQuery, "expected and actual output do not match")

		actualRenderedExpectedArgs := testutils.RenderCallArgsToString(t, returnedExpectedArgs)
		assert.Equal(t, expectedRenderedExpectedArgs, actualRenderedExpectedArgs, "expected and actual output do not match")

		actualRenderedCallArgs := testutils.RenderCallArgsToString(t, returnedCallArgs)
		assert.Equal(t, expectedRenderedCallArgs, actualRenderedCallArgs, "expected and actual output do not match")
	})
}

func Test_buildTestDBArchiveSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBArchiveSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleItem.ID,
		}
		actualQuery, actualArgs := p.buildArchiveItemQuery(exampleItem.ID, exampleUser.ID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBArchiveSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleItem.ID,
		}
		actualQuery, actualArgs := s.buildArchiveItemQuery(exampleItem.ID, exampleUser.ID)

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
		typ := proj.DataTypes[0]
		x := buildTestDBArchiveSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE items SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleItem.ID,
		}
		actualQuery, actualArgs := m.buildArchiveItemQuery(exampleItem.ID, exampleUser.ID)

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

func Test_buildTestDBArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTestDBArchiveSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestPostgres_ArchiveItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := p.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := p.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
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
		typ := proj.DataTypes[0]
		x := buildTestDBArchiveSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestSqlite_ArchiveItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET last_updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := s.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := s.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		s, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := s.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
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
		typ := proj.DataTypes[0]
		x := buildTestDBArchiveSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"errors"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestMariaDB_ArchiveItem(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE items SET last_updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to_user = ? AND id = ?"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(1, 1))

		err := m.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnResult(gosqlmock.NewResult(0, 0))

		err := m.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		m, mockDB := buildTestService(t)

		exampleUser := fake.BuildFakeUser()
		exampleItem := fake.BuildFakeItem()
		exampleItem.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleItem.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveItem(ctx, exampleItem.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
