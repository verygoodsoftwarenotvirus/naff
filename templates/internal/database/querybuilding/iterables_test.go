package querybuilding

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_iterablesDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesDotGo(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"sync"
)

const (
	itemsTableName           = "items"
	itemsTableNameColumn     = "name"
	itemsTableDetailsColumn  = "details"
	itemsUserOwnershipColumn = "belongs_to_user"
)

var (
	itemsTableColumns = []string{
		fmt.Sprintf("%s.%s", itemsTableName, idColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableNameColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableDetailsColumn),
		fmt.Sprintf("%s.%s", itemsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn),
	}
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (p *Postgres) scanItem(scan v1.Scanner) (*v11.Item, error) {
	x := &v11.Item{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (p *Postgres) scanItems(rows v1.ResultIterator) ([]v11.Item, error) {
	var (
		list []v11.Item
	)

	for rows.Next() {
		x, err := p.scanItem(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists
func (p *Postgres) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ItemExists queries the database to see if a given item belonging to a given user exists.
func (p *Postgres) ItemExists(ctx context.Context, itemID, userID uint64) (exists bool, err error) {
	query, args := p.buildItemExistsQuery(itemID, userID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetItem fetches an item from the database.
func (p *Postgres) GetItem(ctx context.Context, itemID, userID uint64) (*v11.Item, error) {
	query, args := p.buildGetItemQuery(itemID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanItem(row)
}

var (
	allItemsCountQueryBuilder sync.Once
	allItemsCountQuery        string
)

// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllItemsCountQuery() string {
	allItemsCountQueryBuilder.Do(func() {
		var err error

		allItemsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, itemsTableName)).
			From(itemsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allItemsCountQuery
}

// GetAllItemsCount will fetch the count of items from the database.
func (p *Postgres) GetAllItemsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllItemsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfItemsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllItems fetches every item from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllItems(ctx context.Context, resultChannel chan []v11.Item) error {
	count, err := p.GetAllItemsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfItemsQuery(begin, end)
			logger := p.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := p.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			items, err := p.scanItems(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- items
		}(beginID, endID)
	}

	return nil
}

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetItemsQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetItems fetches a list of items from the database that meet a particular filter.
func (p *Postgres) GetItems(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.ItemList, error) {
	query, args := p.buildGetItemsQuery(userID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := p.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v11.ItemList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Items: items,
	}

	return list, nil
}

// buildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given user,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetItemsWithIDsQuery(userID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(itemsTableColumns...).
		From(itemsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(itemsTableColumns...).
		FromSelect(subqueryBuilder, itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (p *Postgres) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v11.Item, error) {
	if limit == 0 {
		limit = uint8(v11.DefaultLimit)
	}

	query, args := p.buildGetItemsWithIDsQuery(userID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := p.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}

// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (p *Postgres) buildCreateItemQuery(input *v11.Item) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			itemsTableNameColumn,
			itemsTableDetailsColumn,
			itemsUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateItem creates an item in the database.
func (p *Postgres) CreateItem(ctx context.Context, input *v11.ItemCreationInput) (*v11.Item, error) {
	x := &v11.Item{
		Name:          input.Name,
		Details:       input.Details,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildCreateItemQuery(x)

	// create the item.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing item creation query: %w", err)
	}

	return x, nil
}

// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateItemQuery(input *v11.Item) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(itemsTableName).
		Set(itemsTableNameColumn, input.Name).
		Set(itemsTableDetailsColumn, input.Details).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 input.ID,
			itemsUserOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.
func (p *Postgres) UpdateItem(ctx context.Context, input *v11.Item) error {
	query, args := p.buildUpdateItemQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.
func (p *Postgres) buildArchiveItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 itemID,
			archivedOnColumn:         nil,
			itemsUserOwnershipColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveItem marks an item as archived in the database.
func (p *Postgres) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	query, args := p.buildArchiveItemQuery(itemID, userID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesDotGo(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"sync"
)

const (
	itemsTableName           = "items"
	itemsTableNameColumn     = "name"
	itemsTableDetailsColumn  = "details"
	itemsUserOwnershipColumn = "belongs_to_user"
)

var (
	itemsTableColumns = []string{
		fmt.Sprintf("%s.%s", itemsTableName, idColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableNameColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableDetailsColumn),
		fmt.Sprintf("%s.%s", itemsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn),
	}
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (s *Sqlite) scanItem(scan v1.Scanner) (*v11.Item, error) {
	x := &v11.Item{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (s *Sqlite) scanItems(rows v1.ResultIterator) ([]v11.Item, error) {
	var (
		list []v11.Item
	)

	for rows.Next() {
		x, err := s.scanItem(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		s.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists
func (s *Sqlite) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ItemExists queries the database to see if a given item belonging to a given user exists.
func (s *Sqlite) ItemExists(ctx context.Context, itemID, userID uint64) (exists bool, err error) {
	query, args := s.buildItemExistsQuery(itemID, userID)

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.
func (s *Sqlite) buildGetItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetItem fetches an item from the database.
func (s *Sqlite) GetItem(ctx context.Context, itemID, userID uint64) (*v11.Item, error) {
	query, args := s.buildGetItemQuery(itemID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return s.scanItem(row)
}

var (
	allItemsCountQueryBuilder sync.Once
	allItemsCountQuery        string
)

// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (s *Sqlite) buildGetAllItemsCountQuery() string {
	allItemsCountQueryBuilder.Do(func() {
		var err error

		allItemsCountQuery, _, err = s.sqlBuilder.
			Select(fmt.Sprintf(countQuery, itemsTableName)).
			From(itemsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		s.logQueryBuildingError(err)
	})

	return allItemsCountQuery
}

// GetAllItemsCount will fetch the count of items from the database.
func (s *Sqlite) GetAllItemsCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllItemsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range.
func (s *Sqlite) buildGetBatchOfItemsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): endID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetAllItems fetches every item from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (s *Sqlite) GetAllItems(ctx context.Context, resultChannel chan []v11.Item) error {
	count, err := s.GetAllItemsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := s.buildGetBatchOfItemsQuery(begin, end)
			logger := s.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := s.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			items, err := s.scanItems(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- items
		}(beginID, endID)
	}

	return nil
}

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetItemsQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetItems fetches a list of items from the database that meet a particular filter.
func (s *Sqlite) GetItems(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.ItemList, error) {
	query, args := s.buildGetItemsQuery(userID, filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := s.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v11.ItemList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Items: items,
	}

	return list, nil
}

// buildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given user,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (s *Sqlite) buildGetItemsWithIDsQuery(userID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 ids,
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", itemsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (s *Sqlite) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v11.Item, error) {
	if limit == 0 {
		limit = uint8(v11.DefaultLimit)
	}

	query, args := s.buildGetItemsWithIDsQuery(userID, limit, ids)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := s.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}

// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (s *Sqlite) buildCreateItemQuery(input *v11.Item) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			itemsTableNameColumn,
			itemsTableDetailsColumn,
			itemsUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsToUser,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateItem creates an item in the database.
func (s *Sqlite) CreateItem(ctx context.Context, input *v11.ItemCreationInput) (*v11.Item, error) {
	x := &v11.Item{
		Name:          input.Name,
		Details:       input.Details,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := s.buildCreateItemQuery(x)

	// create the item.
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing item creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	s.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = s.timeTeller.Now()

	return x, nil
}

// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters.
func (s *Sqlite) buildUpdateItemQuery(input *v11.Item) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(itemsTableName).
		Set(itemsTableNameColumn, input.Name).
		Set(itemsTableDetailsColumn, input.Details).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 input.ID,
			itemsUserOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.
func (s *Sqlite) UpdateItem(ctx context.Context, input *v11.Item) error {
	query, args := s.buildUpdateItemQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.
func (s *Sqlite) buildArchiveItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 itemID,
			archivedOnColumn:         nil,
			itemsUserOwnershipColumn: userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveItem marks an item as archived in the database.
func (s *Sqlite) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	query, args := s.buildArchiveItemQuery(itemID, userID)

	res, err := s.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := iterablesDotGo(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
	"sync"
)

const (
	itemsTableName           = "items"
	itemsTableNameColumn     = "name"
	itemsTableDetailsColumn  = "details"
	itemsUserOwnershipColumn = "belongs_to_user"
)

var (
	itemsTableColumns = []string{
		fmt.Sprintf("%s.%s", itemsTableName, idColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableNameColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableDetailsColumn),
		fmt.Sprintf("%s.%s", itemsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn),
	}
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (m *MariaDB) scanItem(scan v1.Scanner) (*v11.Item, error) {
	x := &v11.Item{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (m *MariaDB) scanItems(rows v1.ResultIterator) ([]v11.Item, error) {
	var (
		list []v11.Item
	)

	for rows.Next() {
		x, err := m.scanItem(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		m.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists
func (m *MariaDB) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ItemExists queries the database to see if a given item belonging to a given user exists.
func (m *MariaDB) ItemExists(ctx context.Context, itemID, userID uint64) (exists bool, err error) {
	query, args := m.buildItemExistsQuery(itemID, userID)

	err = m.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetItem fetches an item from the database.
func (m *MariaDB) GetItem(ctx context.Context, itemID, userID uint64) (*v11.Item, error) {
	query, args := m.buildGetItemQuery(itemID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return m.scanItem(row)
}

var (
	allItemsCountQueryBuilder sync.Once
	allItemsCountQuery        string
)

// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllItemsCountQuery() string {
	allItemsCountQueryBuilder.Do(func() {
		var err error

		allItemsCountQuery, _, err = m.sqlBuilder.
			Select(fmt.Sprintf(countQuery, itemsTableName)).
			From(itemsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allItemsCountQuery
}

// GetAllItemsCount will fetch the count of items from the database.
func (m *MariaDB) GetAllItemsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllItemsCountQuery()).Scan(&count)
	return count, err
}

// buildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range.
func (m *MariaDB) buildGetBatchOfItemsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): endID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetAllItems fetches every item from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (m *MariaDB) GetAllItems(ctx context.Context, resultChannel chan []v11.Item) error {
	count, err := m.GetAllItemsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := m.buildGetBatchOfItemsQuery(begin, end)
			logger := m.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := m.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			items, err := m.scanItems(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- items
		}(beginID, endID)
	}

	return nil
}

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetItemsQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetItems fetches a list of items from the database that meet a particular filter.
func (m *MariaDB) GetItems(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.ItemList, error) {
	query, args := m.buildGetItemsQuery(userID, filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := m.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v11.ItemList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Items: items,
	}

	return list, nil
}

// buildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given user,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (m *MariaDB) buildGetItemsWithIDsQuery(userID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 ids,
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", itemsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (m *MariaDB) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v11.Item, error) {
	if limit == 0 {
		limit = uint8(v11.DefaultLimit)
	}

	query, args := m.buildGetItemsWithIDsQuery(userID, limit, ids)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := m.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}

// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (m *MariaDB) buildCreateItemQuery(input *v11.Item) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			itemsTableNameColumn,
			itemsTableDetailsColumn,
			itemsUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsToUser,
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateItem creates an item in the database.
func (m *MariaDB) CreateItem(ctx context.Context, input *v11.ItemCreationInput) (*v11.Item, error) {
	x := &v11.Item{
		Name:          input.Name,
		Details:       input.Details,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := m.buildCreateItemQuery(x)

	// create the item.
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing item creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	m.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = m.timeTeller.Now()

	return x, nil
}

// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters.
func (m *MariaDB) buildUpdateItemQuery(input *v11.Item) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(itemsTableName).
		Set(itemsTableNameColumn, input.Name).
		Set(itemsTableDetailsColumn, input.Details).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 input.ID,
			itemsUserOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.
func (m *MariaDB) UpdateItem(ctx context.Context, input *v11.Item) error {
	query, args := m.buildUpdateItemQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.
func (m *MariaDB) buildArchiveItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 itemID,
			archivedOnColumn:         nil,
			itemsUserOwnershipColumn: userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveItem marks an item as archived in the database.
func (m *MariaDB) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	query, args := m.buildArchiveItemQuery(itemID, userID)

	res, err := m.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIterableConstants(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildIterableConstants(typ)

		expected := `
package example

import ()

const (
	itemsTableName           = "items"
	itemsTableNameColumn     = "name"
	itemsTableDetailsColumn  = "details"
	itemsUserOwnershipColumn = "belongs_to_user"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildIterableConstants(proj.LastDataType())

		expected := `
package example

import ()

const (
	yetAnotherThingsTableName            = "yet_another_things"
	yetAnotherThingsTableOwnershipColumn = "belongs_to_another_thing"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIterableVariableDecs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildIterableVariableDecs(proj, typ)

		expected := `
package example

import (
	"fmt"
)

var (
	itemsTableColumns = []string{
		fmt.Sprintf("%s.%s", itemsTableName, idColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableNameColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableDetailsColumn),
		fmt.Sprintf("%s.%s", itemsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn),
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildIterableVariableDecs(proj, proj.DataTypes[0])

		expected := `
package example

import (
	"fmt"
)

var (
	thingsTableColumns = []string{
		fmt.Sprintf("%s.%s", thingsTableName, idColumn),
		fmt.Sprintf("%s.%s", thingsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", thingsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", thingsTableName, archivedOnColumn),
	}

	thingsOnAnotherThingsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.%s", thingsTableName, anotherThingsTableName, anotherThingsTableOwnershipColumn, thingsTableName, idColumn)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTableColumns(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildTableColumns(typ)

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleFunction(
		fmt.Sprintf("%s.%s", itemsTableName, idColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableNameColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsTableDetailsColumn),
		fmt.Sprintf("%s.%s", itemsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn),
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildTableColumns(proj.LastDataType())

		expected := `
package main

import (
	"fmt"
)

func main() {
	exampleFunction(
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn),
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn),
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanFields(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanFields(typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildScanFields(proj.LastDataType())

		expected := `
package main

import ()

func main() {
	exampleFunction(
		&x.ID,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToAnotherThing,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (p *Postgres) scanItem(scan v1.Scanner) (*v11.Item, error) {
	x := &v11.Item{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (s *Sqlite) scanItem(scan v1.Scanner) (*v11.Item, error) {
	x := &v11.Item{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (m *MariaDB) scanItem(scan v1.Scanner) (*v11.Item, error) {
	x := &v11.Item{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
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

		x := buildScanSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItem takes a database Scanner (i.e. *sql.Row) and scans the result into an Item struct
func (p *Postgres) scanItem(scan v1.Scanner, includeCount bool) (*v11.Item, uint64, error) {
	x := &v11.Item{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Details,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanListOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (p *Postgres) scanItems(rows v1.ResultIterator) ([]v11.Item, error) {
	var (
		list []v11.Item
	)

	for rows.Next() {
		x, err := p.scanItem(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (s *Sqlite) scanItems(rows v1.ResultIterator) ([]v11.Item, error) {
	var (
		list []v11.Item
	)

	for rows.Next() {
		x, err := s.scanItem(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		s.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildScanListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (m *MariaDB) scanItems(rows v1.ResultIterator) ([]v11.Item, error) {
	var (
		list []v11.Item
	)

	for rows.Next() {
		x, err := m.scanItem(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		m.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
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

		x := buildScanListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// scanItems takes a logger and some database rows and turns them into a slice of items.
func (p *Postgres) scanItems(rows v1.ResultIterator) ([]v11.Item, uint64, error) {
	var (
		list  []v11.Item
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanItem(rows, true)
		if err != nil {
			return nil, 0, err
		}

		if count == 0 {
			count = c
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, count, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingExistsQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists
func (p *Postgres) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildYetAnotherThingExistsQuery constructs a SQL query for checking if a yet another thing with a given ID belong to a an another thing with a given ID exists
func (p *Postgres) buildYetAnotherThingExistsQuery(thingID, anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn)).
		Prefix(existencePrefix).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
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
		typ.BelongsToStruct = nil
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID exists
func (p *Postgres) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists
func (s *Sqlite) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildYetAnotherThingExistsQuery constructs a SQL query for checking if a yet another thing with a given ID belong to a an another thing with a given ID exists
func (s *Sqlite) buildYetAnotherThingExistsQuery(thingID, anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn)).
		Prefix(existencePrefix).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
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
		typ.BelongsToStruct = nil
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID exists
func (s *Sqlite) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID belong to a user with a given ID exists
func (m *MariaDB) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildYetAnotherThingExistsQuery constructs a SQL query for checking if a yet another thing with a given ID belong to a an another thing with a given ID exists
func (m *MariaDB) buildYetAnotherThingExistsQuery(thingID, anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn)).
		Prefix(existencePrefix).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
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
		typ.BelongsToStruct = nil
		x := buildSomethingExistsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildItemExistsQuery constructs a SQL query for checking if an item with a given ID exists
func (m *MariaDB) buildItemExistsQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", itemsTableName, idColumn)).
		Prefix(existencePrefix).
		From(itemsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingExistsFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExistsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
)

// ItemExists queries the database to see if a given item belonging to a given user exists.
func (p *Postgres) ItemExists(ctx context.Context, itemID, userID uint64) (exists bool, err error) {
	query, args := p.buildItemExistsQuery(itemID, userID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExistsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
)

// ItemExists queries the database to see if a given item belonging to a given user exists.
func (s *Sqlite) ItemExists(ctx context.Context, itemID, userID uint64) (exists bool, err error) {
	query, args := s.buildItemExistsQuery(itemID, userID)

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingExistsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
)

// ItemExists queries the database to see if a given item belonging to a given user exists.
func (m *MariaDB) ItemExists(ctx context.Context, itemID, userID uint64) (exists bool, err error) {
	query, args := m.buildItemExistsQuery(itemID, userID)

	err = m.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetYetAnotherThingQuery constructs a SQL query for fetching a yet another thing with a given ID belong to an another thing with a given ID.
func (p *Postgres) buildGetYetAnotherThingQuery(thingID, anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToStruct = nil
		typ.BelongsToUser = false
		typ.RestrictedToUser = false
		typ.IsEnumeration = true
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID.
func (p *Postgres) buildGetItemQuery(itemID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): itemID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.
func (s *Sqlite) buildGetItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetYetAnotherThingQuery constructs a SQL query for fetching a yet another thing with a given ID belong to an another thing with a given ID.
func (s *Sqlite) buildGetYetAnotherThingQuery(thingID, anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToStruct = nil
		typ.BelongsToUser = false
		typ.RestrictedToUser = false
		typ.IsEnumeration = true
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID.
func (s *Sqlite) buildGetItemQuery(itemID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): itemID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 itemID,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetYetAnotherThingQuery constructs a SQL query for fetching a yet another thing with a given ID belong to an another thing with a given ID.
func (m *MariaDB) buildGetYetAnotherThingQuery(thingID, anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):                             yetAnotherThingID,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with enumeration", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToStruct = nil
		typ.BelongsToUser = false
		typ.RestrictedToUser = false
		typ.IsEnumeration = true
		x := buildGetSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemQuery constructs a SQL query for fetching an item with a given ID.
func (m *MariaDB) buildGetItemQuery(itemID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): itemID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem fetches an item from the database.
func (p *Postgres) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	query, args := p.buildGetItemQuery(itemID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanItem(row)
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

		x := buildGetSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem fetches an item from the database.
func (p *Postgres) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	query, args := p.buildGetItemQuery(itemID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	i, _, err := p.scanItem(row, false)
	return i, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem fetches an item from the database.
func (s *Sqlite) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	query, args := s.buildGetItemQuery(itemID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return s.scanItem(row)
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

		x := buildGetSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem fetches an item from the database.
func (s *Sqlite) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	query, args := s.buildGetItemQuery(itemID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)
	i, _, err := s.scanItem(row, false)
	return i, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem fetches an item from the database.
func (m *MariaDB) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	query, args := m.buildGetItemQuery(itemID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return m.scanItem(row)
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

		x := buildGetSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItem fetches an item from the database.
func (m *MariaDB) GetItem(ctx context.Context, itemID, userID uint64) (*v1.Item, error) {
	query, args := m.buildGetItemQuery(itemID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	i, _, err := m.scanItem(row, false)
	return i, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSomethingAllCountQueryDecls(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingAllCountQueryDecls(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	allItemsCountQueryBuilder sync.Once
	allItemsCountQuery        string
)

// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllItemsCountQuery() string {
	allItemsCountQueryBuilder.Do(func() {
		var err error

		allItemsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, itemsTableName)).
			From(itemsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allItemsCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingAllCountQueryDecls(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	allItemsCountQueryBuilder sync.Once
	allItemsCountQuery        string
)

// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (s *Sqlite) buildGetAllItemsCountQuery() string {
	allItemsCountQueryBuilder.Do(func() {
		var err error

		allItemsCountQuery, _, err = s.sqlBuilder.
			Select(fmt.Sprintf(countQuery, itemsTableName)).
			From(itemsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		s.logQueryBuildingError(err)
	})

	return allItemsCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildSomethingAllCountQueryDecls(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	allItemsCountQueryBuilder sync.Once
	allItemsCountQuery        string
)

// buildGetAllItemsCountQuery returns a query that fetches the total number of items in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllItemsCountQuery() string {
	allItemsCountQueryBuilder.Do(func() {
		var err error

		allItemsCountQuery, _, err = m.sqlBuilder.
			Select(fmt.Sprintf(countQuery, itemsTableName)).
			From(itemsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
			}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allItemsCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllSomethingCountFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllSomethingCountFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"context"
)

// GetAllItemsCount will fetch the count of items from the database.
func (p *Postgres) GetAllItemsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllItemsCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllSomethingCountFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"context"
)

// GetAllItemsCount will fetch the count of items from the database.
func (s *Sqlite) GetAllItemsCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllItemsCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllSomethingCountFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"context"
)

// GetAllItemsCount will fetch the count of items from the database.
func (m *MariaDB) GetAllItemsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllItemsCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetBatchOfSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetBatchOfSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfItemsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetBatchOfSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range.
func (s *Sqlite) buildGetBatchOfItemsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): endID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetBatchOfSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetBatchOfItemsQuery returns a query that fetches every item in the database within a bucketed range.
func (m *MariaDB) buildGetBatchOfItemsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn): endID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetAllItems fetches every item from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllItems(ctx context.Context, resultChannel chan []v1.Item) error {
	count, err := p.GetAllItemsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfItemsQuery(begin, end)
			logger := p.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := p.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			items, err := p.scanItems(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- items
		}(beginID, endID)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetAllItems fetches every item from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (s *Sqlite) GetAllItems(ctx context.Context, resultChannel chan []v1.Item) error {
	count, err := s.GetAllItemsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := s.buildGetBatchOfItemsQuery(begin, end)
			logger := s.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := s.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			items, err := s.scanItems(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- items
		}(beginID, endID)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetAllOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetAllItems fetches every item from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (m *MariaDB) GetAllItems(ctx context.Context, resultChannel chan []v1.Item) error {
	count, err := m.GetAllItemsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := m.buildGetBatchOfItemsQuery(begin, end)
			logger := m.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := m.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			items, err := m.scanItems(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- items
		}(beginID, endID)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetItemsQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
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

		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetItemsQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(itemsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllItemsCountQuery()))...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetItemsQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
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

		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (s *Sqlite) buildGetItemsQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(append(itemsTableColumns, fmt.Sprintf("(%s)", s.buildGetAllItemsCountQuery()))...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetItemsQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
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

		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetItemsQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(append(itemsTableColumns, fmt.Sprintf("(%s)", m.buildGetAllItemsCountQuery()))...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	// note: these tests don't test the other databases, because as of this writing,
	//  the only thing that these modifications affect is the comment

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetYetAnotherThingsQuery builds a SQL query selecting yet another things that adhere to a given QueryFilter and belong to a given another thing,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetYetAnotherThingsQuery(thingID, anotherThingID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn):                     nil,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, yetAnotherThingsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain and user restriction", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		typ.BelongsToUser = true
		typ.RestrictedToUser = true
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetYetAnotherThingsQuery builds a SQL query selecting yet another things that adhere to a given QueryFilter, and belong to a given user and another thing,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetYetAnotherThingsQuery(thingID, anotherThingID, userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Join(anotherThingsOnYetAnotherThingsJoinClause).
		Join(thingsOnAnotherThingsJoinClause).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn):                     nil,
			fmt.Sprintf("%s.%s", thingsTableName, idColumn):                                       thingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, idColumn):                                anotherThingID,
			fmt.Sprintf("%s.%s", anotherThingsTableName, anotherThingsTableOwnershipColumn):       thingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsTableOwnershipColumn): anotherThingID,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, yetAnotherThingsUserOwnershipColumn):  userID,
		}).
		OrderBy(fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, yetAnotherThingsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres belonging to a user and no other", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		typ.BelongsToStruct = nil
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetYetAnotherThingsQuery builds a SQL query selecting yet another things that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetYetAnotherThingsQuery(filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, yetAnotherThingsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		x := buildGetListOfSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildGetItemsQuery builds a SQL query selecting items that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetItemsQuery(filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", itemsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, itemsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems fetches a list of items from the database that meet a particular filter.
func (p *Postgres) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	query, args := p.buildGetItemsQuery(userID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := p.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v1.ItemList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Items: items,
	}

	return list, nil
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

		x := buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems fetches a list of items from the database that meet a particular filter.
func (p *Postgres) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	query, args := p.buildGetItemsQuery(userID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, count, err := p.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v1.ItemList{
		Pagination: v1.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Items: items,
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems fetches a list of items from the database that meet a particular filter.
func (s *Sqlite) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	query, args := s.buildGetItemsQuery(userID, filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := s.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v1.ItemList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Items: items,
	}

	return list, nil
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

		x := buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems fetches a list of items from the database that meet a particular filter.
func (s *Sqlite) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	query, args := s.buildGetItemsQuery(userID, filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, count, err := s.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v1.ItemList{
		Pagination: v1.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Items: items,
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems fetches a list of items from the database that meet a particular filter.
func (m *MariaDB) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	query, args := m.buildGetItemsQuery(userID, filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := m.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v1.ItemList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Items: items,
	}

	return list, nil
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

		x := buildGetListOfSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItems fetches a list of items from the database that meet a particular filter.
func (m *MariaDB) GetItems(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.ItemList, error) {
	query, args := m.buildGetItemsQuery(userID, filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, count, err := m.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &v1.ItemList{
		Pagination: v1.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Items: items,
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetListOfSomethingWithIDsQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given user,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetItemsWithIDsQuery(userID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(itemsTableColumns...).
		From(itemsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(itemsTableColumns...).
		FromSelect(subqueryBuilder, itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetYetAnotherThingsWithIDsQuery builds a SQL query selecting yetAnotherThings that belong to a given another thing,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetYetAnotherThingsWithIDsQuery(thingID, anotherThingID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		FromSelect(subqueryBuilder, yetAnotherThingsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemsWithIDsQuery builds a SQL query selecting items
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetItemsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(itemsTableColumns...).
		From(itemsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(itemsTableColumns...).
		FromSelect(subqueryBuilder, itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given user,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (s *Sqlite) buildGetItemsWithIDsQuery(userID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 ids,
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", itemsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetYetAnotherThingsWithIDsQuery builds a SQL query selecting yetAnotherThings that belong to a given another thing,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (s *Sqlite) buildGetYetAnotherThingsWithIDsQuery(thingID, anotherThingID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := s.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):         ids,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", yetAnotherThingsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemsWithIDsQuery builds a SQL query selecting items
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (s *Sqlite) buildGetItemsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := s.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):         ids,
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", itemsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemsWithIDsQuery builds a SQL query selecting items that belong to a given user,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (m *MariaDB) buildGetItemsWithIDsQuery(userID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):                 ids,
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn):         nil,
			fmt.Sprintf("%s.%s", itemsTableName, itemsUserOwnershipColumn): userID,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", itemsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, proj.LastDataType())

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetYetAnotherThingsWithIDsQuery builds a SQL query selecting yetAnotherThings that belong to a given another thing,
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (m *MariaDB) buildGetYetAnotherThingsWithIDsQuery(thingID, anotherThingID uint64, limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := m.sqlBuilder.
		Select(yetAnotherThingsTableColumns...).
		From(yetAnotherThingsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, idColumn):         ids,
			fmt.Sprintf("%s.%s", yetAnotherThingsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", yetAnotherThingsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with enumeration", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		x := buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetItemsWithIDsQuery builds a SQL query selecting items
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (m *MariaDB) buildGetItemsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	var whenThenStatement string
	for i, id := range ids {
		if i != 0 {
			whenThenStatement += " "
		}
		whenThenStatement += fmt.Sprintf("WHEN %d THEN %d", id, i)
	}
	whenThenStatement += " END"

	builder := m.sqlBuilder.
		Select(itemsTableColumns...).
		From(itemsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", itemsTableName, idColumn):         ids,
			fmt.Sprintf("%s.%s", itemsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("CASE %s.%s %s", itemsTableName, idColumn, whenThenStatement)).
		Limit(uint64(limit))

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("panics with invalid database", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("invalid")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		assert.Panics(t, func() { buildGetListOfSomethingWithIDsQueryFuncDecl(proj, dbvendor, typ) })
	})
}

func Test_buildGetListOfSomethingWithIDsFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (p *Postgres) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v1.Item, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := p.buildGetItemsWithIDsQuery(userID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := p.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
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
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetYetAnotherThingsWithIDs fetches a list of yet another things from the database that exist within a given set of IDs.
func (p *Postgres) GetYetAnotherThingsWithIDs(ctx context.Context, thingID, anotherThingID uint64, limit uint8, ids []uint64) ([]v1.YetAnotherThing, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := p.buildGetYetAnotherThingsWithIDsQuery(thingID, anotherThingID, limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for yet another things")
	}

	yetAnotherThings, err := p.scanYetAnotherThings(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return yetAnotherThings, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToStruct = nil
		typ.BelongsToUser = false
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (p *Postgres) GetItemsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]v1.Item, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := p.buildGetItemsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := p.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (s *Sqlite) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v1.Item, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := s.buildGetItemsWithIDsQuery(userID, limit, ids)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := s.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetYetAnotherThingsWithIDs fetches a list of yet another things from the database that exist within a given set of IDs.
func (s *Sqlite) GetYetAnotherThingsWithIDs(ctx context.Context, thingID, anotherThingID uint64, limit uint8, ids []uint64) ([]v1.YetAnotherThing, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := s.buildGetYetAnotherThingsWithIDsQuery(thingID, anotherThingID, limit, ids)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for yet another things")
	}

	yetAnotherThings, err := s.scanYetAnotherThings(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return yetAnotherThings, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToStruct = nil
		typ.BelongsToUser = false
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (s *Sqlite) GetItemsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]v1.Item, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := s.buildGetItemsWithIDsQuery(limit, ids)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := s.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (m *MariaDB) GetItemsWithIDs(ctx context.Context, userID uint64, limit uint8, ids []uint64) ([]v1.Item, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := m.buildGetItemsWithIDsQuery(userID, limit, ids)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := m.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetYetAnotherThingsWithIDs fetches a list of yet another things from the database that exist within a given set of IDs.
func (m *MariaDB) GetYetAnotherThingsWithIDs(ctx context.Context, thingID, anotherThingID uint64, limit uint8, ids []uint64) ([]v1.YetAnotherThing, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := m.buildGetYetAnotherThingsWithIDsQuery(thingID, anotherThingID, limit, ids)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for yet another things")
	}

	yetAnotherThings, err := m.scanYetAnotherThings(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return yetAnotherThings, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with enumeration", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToStruct = nil
		typ.BelongsToUser = false
		x := buildGetListOfSomethingWithIDsFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// GetItemsWithIDs fetches a list of items from the database that exist within a given set of IDs.
func (m *MariaDB) GetItemsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]v1.Item, error) {
	if limit == 0 {
		limit = uint8(v1.DefaultLimit)
	}

	query, args := m.buildGetItemsWithIDsQuery(limit, ids)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for items")
	}

	items, err := m.scanItems(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return items, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_determineCreationColumns(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := determineCreationColumns(typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		itemsTableNameColumn,
		itemsTableDetailsColumn,
		itemsUserOwnershipColumn,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := determineCreationColumns(typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		yetAnotherThingsTableOwnershipColumn,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_determineCreationQueryValues(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := determineCreationQueryValues("example", typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		example.Name,
		example.Details,
		example.BelongsToUser,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("with ownership chain", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := determineCreationQueryValues("example", typ)

		expected := `
package main

import ()

func main() {
	exampleFunction(
		example.BelongsToAnotherThing,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (p *Postgres) buildCreateItemQuery(input *v1.Item) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			itemsTableNameColumn,
			itemsTableDetailsColumn,
			itemsUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (s *Sqlite) buildCreateItemQuery(input *v1.Item) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			itemsTableNameColumn,
			itemsTableDetailsColumn,
			itemsUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsToUser,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildCreateItemQuery takes an item and returns a creation query for that item and the relevant arguments.
func (m *MariaDB) buildCreateItemQuery(input *v1.Item) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(itemsTableName).
		Columns(
			itemsTableNameColumn,
			itemsTableDetailsColumn,
			itemsUserOwnershipColumn,
		).
		Values(
			input.Name,
			input.Details,
			input.BelongsToUser,
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateItem creates an item in the database.
func (p *Postgres) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (*v1.Item, error) {
	x := &v1.Item{
		Name:          input.Name,
		Details:       input.Details,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildCreateItemQuery(x)

	// create the item.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing item creation query: %w", err)
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateItem creates an item in the database.
func (s *Sqlite) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (*v1.Item, error) {
	x := &v1.Item{
		Name:          input.Name,
		Details:       input.Details,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := s.buildCreateItemQuery(x)

	// create the item.
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing item creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	s.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = s.timeTeller.Now()

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateItem creates an item in the database.
func (m *MariaDB) CreateItem(ctx context.Context, input *v1.ItemCreationInput) (*v1.Item, error) {
	x := &v1.Item{
		Name:          input.Name,
		Details:       input.Details,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := m.buildCreateItemQuery(x)

	// create the item.
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing item creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	m.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = m.timeTeller.Now()

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("panics with invalid database", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("invalid")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]

		assert.Panics(t, func() { buildCreateSomethingFuncDecl(proj, dbvendor, typ) })
	})

	T.Run("postgres with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateYetAnotherThing creates a yet another thing in the database.
func (p *Postgres) CreateYetAnotherThing(ctx context.Context, input *v1.YetAnotherThingCreationInput) (*v1.YetAnotherThing, error) {
	x := &v1.YetAnotherThing{
		BelongsToAnotherThing: input.BelongsToAnotherThing,
	}

	query, args := p.buildCreateYetAnotherThingQuery(x)

	// create the yet another thing.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing yet another thing creation query: %w", err)
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateYetAnotherThing creates a yet another thing in the database.
func (s *Sqlite) CreateYetAnotherThing(ctx context.Context, input *v1.YetAnotherThingCreationInput) (*v1.YetAnotherThing, error) {
	x := &v1.YetAnotherThing{
		BelongsToAnotherThing: input.BelongsToAnotherThing,
	}

	query, args := s.buildCreateYetAnotherThingQuery(x)

	// create the yet another thing.
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing yet another thing creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	s.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = s.timeTeller.Now()

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildCreateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// CreateYetAnotherThing creates a yet another thing in the database.
func (m *MariaDB) CreateYetAnotherThing(ctx context.Context, input *v1.YetAnotherThingCreationInput) (*v1.YetAnotherThing, error) {
	x := &v1.YetAnotherThing{
		BelongsToAnotherThing: input.BelongsToAnotherThing,
	}

	query, args := m.buildCreateYetAnotherThingQuery(x)

	// create the yet another thing.
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing yet another thing creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	m.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = m.timeTeller.Now()

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateItemQuery(input *v1.Item) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(itemsTableName).
		Set(itemsTableNameColumn, input.Name).
		Set(itemsTableDetailsColumn, input.Details).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 input.ID,
			itemsUserOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
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
		x := buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildUpdateYetAnotherThingQuery takes a yet another thing and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateYetAnotherThingQuery(input *v1.YetAnotherThing) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             input.ID,
			yetAnotherThingsTableOwnershipColumn: input.BelongsToAnotherThing,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters.
func (s *Sqlite) buildUpdateItemQuery(input *v1.Item) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(itemsTableName).
		Set(itemsTableNameColumn, input.Name).
		Set(itemsTableDetailsColumn, input.Details).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 input.ID,
			itemsUserOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildUpdateYetAnotherThingQuery takes a yet another thing and returns an update SQL query, with the relevant query parameters.
func (s *Sqlite) buildUpdateYetAnotherThingQuery(input *v1.YetAnotherThing) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             input.ID,
			yetAnotherThingsTableOwnershipColumn: input.BelongsToAnotherThing,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildUpdateItemQuery takes an item and returns an update SQL query, with the relevant query parameters.
func (m *MariaDB) buildUpdateItemQuery(input *v1.Item) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(itemsTableName).
		Set(itemsTableNameColumn, input.Name).
		Set(itemsTableDetailsColumn, input.Details).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 input.ID,
			itemsUserOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildUpdateSomethingQueryFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// buildUpdateYetAnotherThingQuery takes a yet another thing and returns an update SQL query, with the relevant query parameters.
func (m *MariaDB) buildUpdateYetAnotherThingQuery(input *v1.YetAnotherThing) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             input.ID,
			yetAnotherThingsTableOwnershipColumn: input.BelongsToAnotherThing,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUpdateSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.
func (p *Postgres) UpdateItem(ctx context.Context, input *v1.Item) error {
	query, args := p.buildUpdateItemQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.
func (s *Sqlite) UpdateItem(ctx context.Context, input *v1.Item) error {
	query, args := s.buildUpdateItemQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildUpdateSomethingFuncDecl(proj, dbvendor, typ)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/pkg/types"
)

// UpdateItem updates a particular item. Note that UpdateItem expects the provided input to have a valid ID.
func (m *MariaDB) UpdateItem(ctx context.Context, input *v1.Item) error {
	query, args := m.buildUpdateItemQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveSomethingQueryFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.
func (p *Postgres) buildArchiveItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 itemID,
			archivedOnColumn:         nil,
			itemsUserOwnershipColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
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
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveYetAnotherThingQuery returns a SQL query which marks a given yet another thing belonging to a given another thing as archived.
func (p *Postgres) buildArchiveYetAnotherThingQuery(anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             yetAnotherThingID,
			archivedOnColumn:                     nil,
			yetAnotherThingsTableOwnershipColumn: anotherThingID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with ownership chain and user ownership", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		typ.BelongsToUser = true
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveYetAnotherThingQuery returns a SQL query which marks a given yet another thing belonging to a given another thing and a given user as archived.
func (p *Postgres) buildArchiveYetAnotherThingQuery(anotherThingID, yetAnotherThingID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             yetAnotherThingID,
			archivedOnColumn:                     nil,
			yetAnotherThingsTableOwnershipColumn: anotherThingID,
			yetAnotherThingsUserOwnershipColumn:  userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("postgres with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		typ.IsEnumeration = true
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveItemQuery returns a SQL query which marks a given item as archived.
func (p *Postgres) buildArchiveItemQuery(itemID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         itemID,
			archivedOnColumn: nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.
func (s *Sqlite) buildArchiveItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 itemID,
			archivedOnColumn:         nil,
			itemsUserOwnershipColumn: userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with ownership chain", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveYetAnotherThingQuery returns a SQL query which marks a given yet another thing belonging to a given another thing as archived.
func (s *Sqlite) buildArchiveYetAnotherThingQuery(anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             yetAnotherThingID,
			archivedOnColumn:                     nil,
			yetAnotherThingsTableOwnershipColumn: anotherThingID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite with enumeration", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		typ.IsEnumeration = true
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveItemQuery returns a SQL query which marks a given item as archived.
func (s *Sqlite) buildArchiveItemQuery(itemID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         itemID,
			archivedOnColumn: nil,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveItemQuery returns a SQL query which marks a given item belonging to a given user as archived.
func (m *MariaDB) buildArchiveItemQuery(itemID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                 itemID,
			archivedOnColumn:         nil,
			itemsUserOwnershipColumn: userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with ownership chain", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		proj.DataTypes = models.BuildOwnershipChain("Thing", "AnotherThing", "YetAnotherThing")
		typ := proj.LastDataType()
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveYetAnotherThingQuery returns a SQL query which marks a given yet another thing belonging to a given another thing as archived.
func (m *MariaDB) buildArchiveYetAnotherThingQuery(anotherThingID, yetAnotherThingID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(yetAnotherThingsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                             yetAnotherThingID,
			archivedOnColumn:                     nil,
			yetAnotherThingsTableOwnershipColumn: anotherThingID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb with enumeration", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		typ.BelongsToUser = false
		typ.BelongsToStruct = nil
		typ.IsEnumeration = true
		x := buildArchiveSomethingQueryFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveItemQuery returns a SQL query which marks a given item as archived.
func (m *MariaDB) buildArchiveItemQuery(itemID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(itemsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         itemID,
			archivedOnColumn: nil,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildArchiveSomethingFuncDecl(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
)

// ArchiveItem marks an item as archived in the database.
func (p *Postgres) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	query, args := p.buildArchiveItemQuery(itemID, userID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
)

// ArchiveItem marks an item as archived in the database.
func (s *Sqlite) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	query, args := s.buildArchiveItemQuery(itemID, userID)

	res, err := s.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildArchiveSomethingFuncDecl(dbvendor, typ)

		expected := `
package example

import (
	"context"
	"database/sql"
)

// ArchiveItem marks an item as archived in the database.
func (m *MariaDB) ArchiveItem(ctx context.Context, itemID, userID uint64) error {
	query, args := m.buildArchiveItemQuery(itemID, userID)

	res, err := m.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
