package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

func buildMockRowFromItem(item *models.Item) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(itemsTableColumns).AddRow(item.ID, item.Name, item.Details, item.CreatedOn, item.UpdatedOn, item.ArchivedOn, item.BelongsTo)
	return exampleRows
}

func buildErroneousMockRowFromItem(item *models.Item) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(itemsTableColumns).AddRow(item.ArchivedOn, item.Name, item.Details, item.CreatedOn, item.UpdatedOn, item.BelongsTo, item.ID)
	return exampleRows
}

func TestSqlite_buildGetItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleItemID := uint64(123)
		exampleUserID := uint64(321)
		expectedArgCount := 2
		expectedQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildGetItemQuery(exampleItemID, exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleItemID, args[1].(uint64))
	})
}

func TestSqlite_GetItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = ? AND id = ?"
		expected := &models.Item{
			ID:      123,
			Name:    "name",
			Details: "details",
		}
		expectedUserID := uint64(321)
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(expectedUserID, expected.ID).WillReturnRows(buildMockRowFromItem(expected))
		actual, err := s.GetItem(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE belongs_to = ? AND id = ?"
		expected := &models.Item{
			ID:      123,
			Name:    "name",
			Details: "details",
		}
		expectedUserID := uint64(321)
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(expectedUserID, expected.ID).WillReturnError(sql.ErrNoRows)
		actual, err := s.GetItem(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetItemCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)
		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := s.buildGetItemCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetItemCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WithArgs(expectedUserID).WillReturnRows(sqlmock.NewRows([]string{
			"count",
		}).AddRow(expectedCount))
		actualCount, err := s.GetItemCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllItemsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM items WHERE archived_on IS NULL"
		actualQuery := s.buildGetAllItemsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllItemsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM items WHERE archived_on IS NULL"
		expectedCount := uint64(666)
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).WillReturnRows(sqlmock.NewRows([]string{
			"count",
		}).AddRow(expectedCount))
		actualCount, err := s.GetAllItemsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetItemsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)
		expectedArgCount := 1
		expectedQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := s.buildGetItemsQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetItems(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedItem1 := &models.Item{
			Name: "name",
		}
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM items WHERE archived_on IS NULL"
		expectedCount := uint64(666)
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnRows(buildMockRowFromItem(expectedItem1))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).WillReturnRows(sqlmock.NewRows([]string{
			"count",
		}).AddRow(expectedCount))
		expected := &models.ItemList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			Items: []models.Item{
				*expectedItem1,
			},
		}
		actual, err := s.GetItems(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnError(sql.ErrNoRows)
		actual, err := s.GetItems(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnError(errors.New("blah"))
		actual, err := s.GetItems(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning item", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedItem1 := &models.Item{
			Name: "name",
		}
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnRows(buildErroneousMockRowFromItem(expectedItem1))
		actual, err := s.GetItems(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedItem1 := &models.Item{
			Name: "name",
		}
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM items WHERE archived_on IS NULL"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnRows(buildMockRowFromItem(expectedItem1))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).WillReturnError(errors.New("blah"))
		actual, err := s.GetItems(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_GetAllItemsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedItem := &models.Item{
			Name: "name",
		}
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnRows(buildMockRowFromItem(expectedItem))
		expected := []models.Item{
			*expectedItem,
		}
		actual, err := s.GetAllItemsForUser(context.Background(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnError(sql.ErrNoRows)
		actual, err := s.GetAllItemsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnError(errors.New("blah"))
		actual, err := s.GetAllItemsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedItem := &models.Item{
			Name: "name",
		}
		expectedListQuery := "SELECT id, name, details, created_on, updated_on, archived_on, belongs_to FROM items WHERE archived_on IS NULL AND belongs_to = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).WithArgs(expectedUserID).WillReturnRows(buildErroneousMockRowFromItem(expectedItem))
		actual, err := s.GetAllItemsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Item{
			Name:      "name",
			Details:   "details",
			BelongsTo: 123,
		}
		expectedArgCount := 3
		expectedQuery := "INSERT INTO items (name,details,belongs_to) VALUES (?,?,?)"
		actualQuery, args := s.buildCreateItemQuery(expected)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Details, args[1].(string))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
	})
}

func TestSqlite_CreateItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.ItemCreationInput{
			Name:      expected.Name,
			BelongsTo: expected.BelongsTo,
		}
		s, mockDB := buildTestService(t)
		expectedCreationQuery := "INSERT INTO items (name,details,belongs_to) VALUES (?,?,?)"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).WithArgs(expected.Name, expected.Details, expected.BelongsTo).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))
		expectedTimeQuery := "SELECT created_on FROM items WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).WithArgs(expected.ID).WillReturnRows(sqlmock.NewRows([]string{
			"created_on",
		}).AddRow(expected.CreatedOn))
		actual, err := s.CreateItem(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.ItemCreationInput{
			Name:      example.Name,
			BelongsTo: example.BelongsTo,
		}
		s, mockDB := buildTestService(t)
		expectedQuery := "INSERT INTO items (name,details,belongs_to) VALUES (?,?,?)"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(example.Name, example.Details, example.BelongsTo).WillReturnError(errors.New("blah"))
		actual, err := s.CreateItem(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdateItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Item{
			ID:        321,
			Name:      "name",
			Details:   "details",
			BelongsTo: 123,
		}
		expectedArgCount := 4
		expectedQuery := "UPDATE items SET name = ?, details = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildUpdateItemQuery(expected)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Name, args[0].(string))
		assert.Equal(t, expected.Details, args[1].(string))
		assert.Equal(t, expected.BelongsTo, args[2].(uint64))
		assert.Equal(t, expected.ID, args[3].(uint64))
	})
}

func TestSqlite_UpdateItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE items SET name = ?, details = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(expected.Name, expected.Details, expected.BelongsTo, expected.ID).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))
		err := s.UpdateItem(context.Background(), expected)
		assert.NoError(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE items SET name = ?, details = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(example.Name, example.Details, example.BelongsTo, example.ID).WillReturnError(errors.New("blah"))
		err := s.UpdateItem(context.Background(), example)
		assert.Error(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchiveItemQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.Item{
			ID:        321,
			Name:      "name",
			Details:   "details",
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE items SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := s.buildArchiveItemQuery(expected.ID, expected.BelongsTo)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestSqlite_ArchiveItem(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE items SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(expected.BelongsTo, expected.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		err := s.ArchiveItem(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.Item{
			ID:        123,
			Name:      "name",
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE items SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).WithArgs(example.BelongsTo, example.ID).WillReturnError(errors.New("blah"))
		err := s.ArchiveItem(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)
		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}