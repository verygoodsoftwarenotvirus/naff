package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_webhooksDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := webhooksDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
	"sync"
)

const (
	commaSeparator = ","

	eventsSeparator = commaSeparator
	typesSeparator  = commaSeparator
	topicsSeparator = commaSeparator

	webhooksTableName              = "webhooks"
	webhooksTableNameColumn        = "name"
	webhooksTableContentTypeColumn = "content_type"
	webhooksTableURLColumn         = "url"
	webhooksTableMethodColumn      = "method"
	webhooksTableEventsColumn      = "events"
	webhooksTableDataTypesColumn   = "data_types"
	webhooksTableTopicsColumn      = "topics"
	webhooksTableOwnershipColumn   = "belongs_to_user"
)

var (
	webhooksTableColumns = []string{
		fmt.Sprintf("%s.%s", webhooksTableName, idColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableNameColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableContentTypeColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableURLColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableMethodColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableEventsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableDataTypesColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableTopicsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn),
	}
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (p *Postgres) scanWebhook(scan v1.Scanner) (*v11.Webhook, error) {
	var (
		x = &v11.Webhook{}
		eventsStr,
		dataTypesStr,
		topicsStr string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ContentType,
		&x.URL,
		&x.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if events := strings.Split(eventsStr, eventsSeparator); len(events) >= 1 && events[0] != "" {
		x.Events = events
	}
	if dataTypes := strings.Split(dataTypesStr, typesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		x.DataTypes = dataTypes
	}
	if topics := strings.Split(topicsStr, topicsSeparator); len(topics) >= 1 && topics[0] != "" {
		x.Topics = topics
	}

	return x, nil
}

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (p *Postgres) scanWebhooks(rows v1.ResultIterator) ([]v11.Webhook, error) {
	var (
		list []v11.Webhook
	)

	for rows.Next() {
		webhook, err := p.scanWebhook(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		p.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook
func (p *Postgres) buildGetWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, idColumn):                     webhookID,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
		}).ToSql()

	p.logQueryBuildingError(err)
	return query, args
}

// GetWebhook fetches a webhook from the database.
func (p *Postgres) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v11.Webhook, error) {
	query, args := p.buildGetWebhookQuery(webhookID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	webhook, err := p.scanWebhook(row)
	if err != nil {
		return nil, buildError(err, "querying for webhook")
	}

	return webhook, nil
}

var (
	getAllWebhooksCountQueryBuilder sync.Once
	getAllWebhooksCountQuery        string
)

// buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership.
func (p *Postgres) buildGetAllWebhooksCountQuery() string {
	getAllWebhooksCountQueryBuilder.Do(func() {
		var err error

		getAllWebhooksCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, webhooksTableName)).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllWebhooksCountQuery
}

// GetAllWebhooksCount will fetch the count of every active webhook in the database.
func (p *Postgres) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllWebhooksCountQuery()).Scan(&count)
	return count, err
}

var (
	getAllWebhooksQueryBuilder sync.Once
	getAllWebhooksQuery        string
)

// buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership.
func (p *Postgres) buildGetAllWebhooksQuery() string {
	getAllWebhooksQueryBuilder.Do(func() {
		var err error

		getAllWebhooksQuery, _, err = p.sqlBuilder.
			Select(webhooksTableColumns...).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllWebhooksQuery
}

// GetAllWebhooks fetches a list of all webhooks from the database.
func (p *Postgres) GetAllWebhooks(ctx context.Context) (*v11.WebhookList, error) {
	rows, err := p.db.QueryContext(ctx, p.buildGetAllWebhooksQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for webhooks: %w", err)
	}

	list, err := p.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v11.WebhookList{
		Pagination: v11.Pagination{
			Page: 1,
		},
		Webhooks: list,
	}

	return x, err
}

// buildGetWebhooksQuery returns a SQL query (and arguments) that would return a
func (p *Postgres) buildGetWebhooksQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn):             nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", webhooksTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, webhooksTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (p *Postgres) GetWebhooks(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.WebhookList, error) {
	query, args := p.buildGetWebhooksQuery(userID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}

	list, err := p.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v11.WebhookList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Webhooks: list,
	}

	return x, err
}

// buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook
func (p *Postgres) buildWebhookCreationQuery(x *v11.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(webhooksTableName).
		Columns(
			webhooksTableNameColumn,
			webhooksTableContentTypeColumn,
			webhooksTableURLColumn,
			webhooksTableMethodColumn,
			webhooksTableEventsColumn,
			webhooksTableDataTypesColumn,
			webhooksTableTopicsColumn,
			webhooksTableOwnershipColumn,
		).
		Values(
			x.Name,
			x.ContentType,
			x.URL,
			x.Method,
			strings.Join(x.Events, eventsSeparator),
			strings.Join(x.DataTypes, typesSeparator),
			strings.Join(x.Topics, topicsSeparator),
			x.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateWebhook creates a webhook in the database.
func (p *Postgres) CreateWebhook(ctx context.Context, input *v11.WebhookCreationInput) (*v11.Webhook, error) {
	x := &v11.Webhook{
		Name:          input.Name,
		ContentType:   input.ContentType,
		URL:           input.URL,
		Method:        input.Method,
		Events:        input.Events,
		DataTypes:     input.DataTypes,
		Topics:        input.Topics,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildWebhookCreationQuery(x)
	if err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn); err != nil {
		return nil, fmt.Errorf("error executing webhook creation query: %w", err)
	}

	return x, nil
}

// buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update.
func (p *Postgres) buildUpdateWebhookQuery(input *v11.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(webhooksTableName).
		Set(webhooksTableNameColumn, input.Name).
		Set(webhooksTableContentTypeColumn, input.ContentType).
		Set(webhooksTableURLColumn, input.URL).
		Set(webhooksTableMethodColumn, input.Method).
		Set(webhooksTableEventsColumn, strings.Join(input.Events, topicsSeparator)).
		Set(webhooksTableDataTypesColumn, strings.Join(input.DataTypes, typesSeparator)).
		Set(webhooksTableTopicsColumn, strings.Join(input.Topics, topicsSeparator)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     input.ID,
			webhooksTableOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID.
func (p *Postgres) UpdateWebhook(ctx context.Context, input *v11.Webhook) error {
	query, args := p.buildUpdateWebhookQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.
func (p *Postgres) buildArchiveWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(webhooksTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     webhookID,
			webhooksTableOwnershipColumn: userID,
			archivedOnColumn:             nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveWebhook archives a webhook from the database by its ID.
func (p *Postgres) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	query, args := p.buildArchiveWebhookQuery(webhookID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := webhooksDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
	"sync"
)

const (
	commaSeparator = ","

	eventsSeparator = commaSeparator
	typesSeparator  = commaSeparator
	topicsSeparator = commaSeparator

	webhooksTableName              = "webhooks"
	webhooksTableNameColumn        = "name"
	webhooksTableContentTypeColumn = "content_type"
	webhooksTableURLColumn         = "url"
	webhooksTableMethodColumn      = "method"
	webhooksTableEventsColumn      = "events"
	webhooksTableDataTypesColumn   = "data_types"
	webhooksTableTopicsColumn      = "topics"
	webhooksTableOwnershipColumn   = "belongs_to_user"
)

var (
	webhooksTableColumns = []string{
		fmt.Sprintf("%s.%s", webhooksTableName, idColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableNameColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableContentTypeColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableURLColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableMethodColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableEventsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableDataTypesColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableTopicsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn),
	}
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (s *Sqlite) scanWebhook(scan v1.Scanner) (*v11.Webhook, error) {
	var (
		x = &v11.Webhook{}
		eventsStr,
		dataTypesStr,
		topicsStr string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ContentType,
		&x.URL,
		&x.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if events := strings.Split(eventsStr, eventsSeparator); len(events) >= 1 && events[0] != "" {
		x.Events = events
	}
	if dataTypes := strings.Split(dataTypesStr, typesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		x.DataTypes = dataTypes
	}
	if topics := strings.Split(topicsStr, topicsSeparator); len(topics) >= 1 && topics[0] != "" {
		x.Topics = topics
	}

	return x, nil
}

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (s *Sqlite) scanWebhooks(rows v1.ResultIterator) ([]v11.Webhook, error) {
	var (
		list []v11.Webhook
	)

	for rows.Next() {
		webhook, err := s.scanWebhook(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		s.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook
func (s *Sqlite) buildGetWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, idColumn):                     webhookID,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
		}).ToSql()

	s.logQueryBuildingError(err)
	return query, args
}

// GetWebhook fetches a webhook from the database.
func (s *Sqlite) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v11.Webhook, error) {
	query, args := s.buildGetWebhookQuery(webhookID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	webhook, err := s.scanWebhook(row)
	if err != nil {
		return nil, buildError(err, "querying for webhook")
	}

	return webhook, nil
}

var (
	getAllWebhooksCountQueryBuilder sync.Once
	getAllWebhooksCountQuery        string
)

// buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership.
func (s *Sqlite) buildGetAllWebhooksCountQuery() string {
	getAllWebhooksCountQueryBuilder.Do(func() {
		var err error

		getAllWebhooksCountQuery, _, err = s.sqlBuilder.
			Select(fmt.Sprintf(countQuery, webhooksTableName)).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllWebhooksCountQuery
}

// GetAllWebhooksCount will fetch the count of every active webhook in the database.
func (s *Sqlite) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllWebhooksCountQuery()).Scan(&count)
	return count, err
}

var (
	getAllWebhooksQueryBuilder sync.Once
	getAllWebhooksQuery        string
)

// buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership.
func (s *Sqlite) buildGetAllWebhooksQuery() string {
	getAllWebhooksQueryBuilder.Do(func() {
		var err error

		getAllWebhooksQuery, _, err = s.sqlBuilder.
			Select(webhooksTableColumns...).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllWebhooksQuery
}

// GetAllWebhooks fetches a list of all webhooks from the database.
func (s *Sqlite) GetAllWebhooks(ctx context.Context) (*v11.WebhookList, error) {
	rows, err := s.db.QueryContext(ctx, s.buildGetAllWebhooksQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for webhooks: %w", err)
	}

	list, err := s.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v11.WebhookList{
		Pagination: v11.Pagination{
			Page: 1,
		},
		Webhooks: list,
	}

	return x, err
}

// buildGetWebhooksQuery returns a SQL query (and arguments) that would return a
func (s *Sqlite) buildGetWebhooksQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn):             nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", webhooksTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, webhooksTableName)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (s *Sqlite) GetWebhooks(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.WebhookList, error) {
	query, args := s.buildGetWebhooksQuery(userID, filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}

	list, err := s.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v11.WebhookList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Webhooks: list,
	}

	return x, err
}

// buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook
func (s *Sqlite) buildWebhookCreationQuery(x *v11.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(webhooksTableName).
		Columns(
			webhooksTableNameColumn,
			webhooksTableContentTypeColumn,
			webhooksTableURLColumn,
			webhooksTableMethodColumn,
			webhooksTableEventsColumn,
			webhooksTableDataTypesColumn,
			webhooksTableTopicsColumn,
			webhooksTableOwnershipColumn,
		).
		Values(
			x.Name,
			x.ContentType,
			x.URL,
			x.Method,
			strings.Join(x.Events, eventsSeparator),
			strings.Join(x.DataTypes, typesSeparator),
			strings.Join(x.Topics, topicsSeparator),
			x.BelongsToUser,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateWebhook creates a webhook in the database.
func (s *Sqlite) CreateWebhook(ctx context.Context, input *v11.WebhookCreationInput) (*v11.Webhook, error) {
	x := &v11.Webhook{
		Name:          input.Name,
		ContentType:   input.ContentType,
		URL:           input.URL,
		Method:        input.Method,
		Events:        input.Events,
		DataTypes:     input.DataTypes,
		Topics:        input.Topics,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := s.buildWebhookCreationQuery(x)
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing webhook creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	s.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = s.timeTeller.Now()

	return x, nil
}

// buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update.
func (s *Sqlite) buildUpdateWebhookQuery(input *v11.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(webhooksTableName).
		Set(webhooksTableNameColumn, input.Name).
		Set(webhooksTableContentTypeColumn, input.ContentType).
		Set(webhooksTableURLColumn, input.URL).
		Set(webhooksTableMethodColumn, input.Method).
		Set(webhooksTableEventsColumn, strings.Join(input.Events, topicsSeparator)).
		Set(webhooksTableDataTypesColumn, strings.Join(input.DataTypes, typesSeparator)).
		Set(webhooksTableTopicsColumn, strings.Join(input.Topics, topicsSeparator)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     input.ID,
			webhooksTableOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID.
func (s *Sqlite) UpdateWebhook(ctx context.Context, input *v11.Webhook) error {
	query, args := s.buildUpdateWebhookQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.
func (s *Sqlite) buildArchiveWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(webhooksTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     webhookID,
			webhooksTableOwnershipColumn: userID,
			archivedOnColumn:             nil,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveWebhook archives a webhook from the database by its ID.
func (s *Sqlite) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	query, args := s.buildArchiveWebhookQuery(webhookID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := webhooksDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
	"sync"
)

const (
	commaSeparator = ","

	eventsSeparator = commaSeparator
	typesSeparator  = commaSeparator
	topicsSeparator = commaSeparator

	webhooksTableName              = "webhooks"
	webhooksTableNameColumn        = "name"
	webhooksTableContentTypeColumn = "content_type"
	webhooksTableURLColumn         = "url"
	webhooksTableMethodColumn      = "method"
	webhooksTableEventsColumn      = "events"
	webhooksTableDataTypesColumn   = "data_types"
	webhooksTableTopicsColumn      = "topics"
	webhooksTableOwnershipColumn   = "belongs_to_user"
)

var (
	webhooksTableColumns = []string{
		fmt.Sprintf("%s.%s", webhooksTableName, idColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableNameColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableContentTypeColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableURLColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableMethodColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableEventsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableDataTypesColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableTopicsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn),
	}
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (m *MariaDB) scanWebhook(scan v1.Scanner) (*v11.Webhook, error) {
	var (
		x = &v11.Webhook{}
		eventsStr,
		dataTypesStr,
		topicsStr string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ContentType,
		&x.URL,
		&x.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if events := strings.Split(eventsStr, eventsSeparator); len(events) >= 1 && events[0] != "" {
		x.Events = events
	}
	if dataTypes := strings.Split(dataTypesStr, typesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		x.DataTypes = dataTypes
	}
	if topics := strings.Split(topicsStr, topicsSeparator); len(topics) >= 1 && topics[0] != "" {
		x.Topics = topics
	}

	return x, nil
}

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (m *MariaDB) scanWebhooks(rows v1.ResultIterator) ([]v11.Webhook, error) {
	var (
		list []v11.Webhook
	)

	for rows.Next() {
		webhook, err := m.scanWebhook(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		m.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook
func (m *MariaDB) buildGetWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, idColumn):                     webhookID,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
		}).ToSql()

	m.logQueryBuildingError(err)
	return query, args
}

// GetWebhook fetches a webhook from the database.
func (m *MariaDB) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v11.Webhook, error) {
	query, args := m.buildGetWebhookQuery(webhookID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	webhook, err := m.scanWebhook(row)
	if err != nil {
		return nil, buildError(err, "querying for webhook")
	}

	return webhook, nil
}

var (
	getAllWebhooksCountQueryBuilder sync.Once
	getAllWebhooksCountQuery        string
)

// buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership.
func (m *MariaDB) buildGetAllWebhooksCountQuery() string {
	getAllWebhooksCountQueryBuilder.Do(func() {
		var err error

		getAllWebhooksCountQuery, _, err = m.sqlBuilder.
			Select(fmt.Sprintf(countQuery, webhooksTableName)).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllWebhooksCountQuery
}

// GetAllWebhooksCount will fetch the count of every active webhook in the database.
func (m *MariaDB) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllWebhooksCountQuery()).Scan(&count)
	return count, err
}

var (
	getAllWebhooksQueryBuilder sync.Once
	getAllWebhooksQuery        string
)

// buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership.
func (m *MariaDB) buildGetAllWebhooksQuery() string {
	getAllWebhooksQueryBuilder.Do(func() {
		var err error

		getAllWebhooksQuery, _, err = m.sqlBuilder.
			Select(webhooksTableColumns...).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllWebhooksQuery
}

// GetAllWebhooks fetches a list of all webhooks from the database.
func (m *MariaDB) GetAllWebhooks(ctx context.Context) (*v11.WebhookList, error) {
	rows, err := m.db.QueryContext(ctx, m.buildGetAllWebhooksQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for webhooks: %w", err)
	}

	list, err := m.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v11.WebhookList{
		Pagination: v11.Pagination{
			Page: 1,
		},
		Webhooks: list,
	}

	return x, err
}

// buildGetWebhooksQuery returns a SQL query (and arguments) that would return a
func (m *MariaDB) buildGetWebhooksQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn):             nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", webhooksTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, webhooksTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (m *MariaDB) GetWebhooks(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.WebhookList, error) {
	query, args := m.buildGetWebhooksQuery(userID, filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}

	list, err := m.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v11.WebhookList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Webhooks: list,
	}

	return x, err
}

// buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook
func (m *MariaDB) buildWebhookCreationQuery(x *v11.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(webhooksTableName).
		Columns(
			webhooksTableNameColumn,
			webhooksTableContentTypeColumn,
			webhooksTableURLColumn,
			webhooksTableMethodColumn,
			webhooksTableEventsColumn,
			webhooksTableDataTypesColumn,
			webhooksTableTopicsColumn,
			webhooksTableOwnershipColumn,
		).
		Values(
			x.Name,
			x.ContentType,
			x.URL,
			x.Method,
			strings.Join(x.Events, eventsSeparator),
			strings.Join(x.DataTypes, typesSeparator),
			strings.Join(x.Topics, topicsSeparator),
			x.BelongsToUser,
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateWebhook creates a webhook in the database.
func (m *MariaDB) CreateWebhook(ctx context.Context, input *v11.WebhookCreationInput) (*v11.Webhook, error) {
	x := &v11.Webhook{
		Name:          input.Name,
		ContentType:   input.ContentType,
		URL:           input.URL,
		Method:        input.Method,
		Events:        input.Events,
		DataTypes:     input.DataTypes,
		Topics:        input.Topics,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := m.buildWebhookCreationQuery(x)
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing webhook creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	m.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = m.timeTeller.Now()

	return x, nil
}

// buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update.
func (m *MariaDB) buildUpdateWebhookQuery(input *v11.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(webhooksTableName).
		Set(webhooksTableNameColumn, input.Name).
		Set(webhooksTableContentTypeColumn, input.ContentType).
		Set(webhooksTableURLColumn, input.URL).
		Set(webhooksTableMethodColumn, input.Method).
		Set(webhooksTableEventsColumn, strings.Join(input.Events, topicsSeparator)).
		Set(webhooksTableDataTypesColumn, strings.Join(input.DataTypes, typesSeparator)).
		Set(webhooksTableTopicsColumn, strings.Join(input.Topics, topicsSeparator)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     input.ID,
			webhooksTableOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID.
func (m *MariaDB) UpdateWebhook(ctx context.Context, input *v11.Webhook) error {
	query, args := m.buildUpdateWebhookQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.
func (m *MariaDB) buildArchiveWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(webhooksTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     webhookID,
			webhooksTableOwnershipColumn: userID,
			archivedOnColumn:             nil,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveWebhook archives a webhook from the database by its ID.
func (m *MariaDB) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	query, args := m.buildArchiveWebhookQuery(webhookID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhooksConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildWebhooksConstDeclarations()

		expected := `
package example

import ()

const (
	commaSeparator = ","

	eventsSeparator = commaSeparator
	typesSeparator  = commaSeparator
	topicsSeparator = commaSeparator

	webhooksTableName              = "webhooks"
	webhooksTableNameColumn        = "name"
	webhooksTableContentTypeColumn = "content_type"
	webhooksTableURLColumn         = "url"
	webhooksTableMethodColumn      = "method"
	webhooksTableEventsColumn      = "events"
	webhooksTableDataTypesColumn   = "data_types"
	webhooksTableTopicsColumn      = "topics"
	webhooksTableOwnershipColumn   = "belongs_to_user"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWebhooksVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildWebhooksVarDeclarations()

		expected := `
package example

import (
	"fmt"
)

var (
	webhooksTableColumns = []string{
		fmt.Sprintf("%s.%s", webhooksTableName, idColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableNameColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableContentTypeColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableURLColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableMethodColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableEventsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableDataTypesColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableTopicsColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn),
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanWebhook(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildScanWebhook(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (p *Postgres) scanWebhook(scan v1.Scanner) (*v11.Webhook, error) {
	var (
		x = &v11.Webhook{}
		eventsStr,
		dataTypesStr,
		topicsStr string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ContentType,
		&x.URL,
		&x.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if events := strings.Split(eventsStr, eventsSeparator); len(events) >= 1 && events[0] != "" {
		x.Events = events
	}
	if dataTypes := strings.Split(dataTypesStr, typesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		x.DataTypes = dataTypes
	}
	if topics := strings.Split(topicsStr, topicsSeparator); len(topics) >= 1 && topics[0] != "" {
		x.Topics = topics
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
		x := buildScanWebhook(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (s *Sqlite) scanWebhook(scan v1.Scanner) (*v11.Webhook, error) {
	var (
		x = &v11.Webhook{}
		eventsStr,
		dataTypesStr,
		topicsStr string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ContentType,
		&x.URL,
		&x.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if events := strings.Split(eventsStr, eventsSeparator); len(events) >= 1 && events[0] != "" {
		x.Events = events
	}
	if dataTypes := strings.Split(dataTypesStr, typesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		x.DataTypes = dataTypes
	}
	if topics := strings.Split(topicsStr, topicsSeparator); len(topics) >= 1 && topics[0] != "" {
		x.Topics = topics
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
		x := buildScanWebhook(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// scanWebhook is a consistent way to turn a *sql.Row into a webhook struct.
func (m *MariaDB) scanWebhook(scan v1.Scanner) (*v11.Webhook, error) {
	var (
		x = &v11.Webhook{}
		eventsStr,
		dataTypesStr,
		topicsStr string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ContentType,
		&x.URL,
		&x.Method,
		&eventsStr,
		&dataTypesStr,
		&topicsStr,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if events := strings.Split(eventsStr, eventsSeparator); len(events) >= 1 && events[0] != "" {
		x.Events = events
	}
	if dataTypes := strings.Split(dataTypesStr, typesSeparator); len(dataTypes) >= 1 && dataTypes[0] != "" {
		x.DataTypes = dataTypes
	}
	if topics := strings.Split(topicsStr, topicsSeparator); len(topics) >= 1 && topics[0] != "" {
		x.Topics = topics
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildScanWebhooks(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (p *Postgres) scanWebhooks(rows v1.ResultIterator) ([]v11.Webhook, error) {
	var (
		list []v11.Webhook
	)

	for rows.Next() {
		webhook, err := p.scanWebhook(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		p.logger.Error(err, "closing rows")
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
		x := buildScanWebhooks(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (s *Sqlite) scanWebhooks(rows v1.ResultIterator) ([]v11.Webhook, error) {
	var (
		list []v11.Webhook
	)

	for rows.Next() {
		webhook, err := s.scanWebhook(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		s.logger.Error(err, "closing rows")
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
		x := buildScanWebhooks(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanWebhooks provides a consistent way to turn sql rows into a slice of webhooks.
func (m *MariaDB) scanWebhooks(rows v1.ResultIterator) ([]v11.Webhook, error) {
	var (
		list []v11.Webhook
	)

	for rows.Next() {
		webhook, err := m.scanWebhook(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *webhook)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		m.logger.Error(err, "closing rows")
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetWebhookQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetWebhookQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook
func (p *Postgres) buildGetWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, idColumn):                     webhookID,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
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

		x := buildBuildGetWebhookQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook
func (s *Sqlite) buildGetWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, idColumn):                     webhookID,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
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

		x := buildBuildGetWebhookQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetWebhookQuery returns a SQL query (and arguments) for retrieving a given webhook
func (m *MariaDB) buildGetWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, idColumn):                     webhookID,
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
		}).ToSql()

	m.logQueryBuildingError(err)
	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_build_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := build_GetWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhook fetches a webhook from the database.
func (p *Postgres) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	query, args := p.buildGetWebhookQuery(webhookID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	webhook, err := p.scanWebhook(row)
	if err != nil {
		return nil, buildError(err, "querying for webhook")
	}

	return webhook, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := build_GetWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhook fetches a webhook from the database.
func (s *Sqlite) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	query, args := s.buildGetWebhookQuery(webhookID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	webhook, err := s.scanWebhook(row)
	if err != nil {
		return nil, buildError(err, "querying for webhook")
	}

	return webhook, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := build_GetWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhook fetches a webhook from the database.
func (m *MariaDB) GetWebhook(ctx context.Context, webhookID, userID uint64) (*v1.Webhook, error) {
	query, args := m.buildGetWebhookQuery(webhookID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	webhook, err := m.scanWebhook(row)
	if err != nil {
		return nil, buildError(err, "querying for webhook")
	}

	return webhook, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetAllWebhooksCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetAllWebhooksCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	getAllWebhooksCountQueryBuilder sync.Once
	getAllWebhooksCountQuery        string
)

// buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership.
func (p *Postgres) buildGetAllWebhooksCountQuery() string {
	getAllWebhooksCountQueryBuilder.Do(func() {
		var err error

		getAllWebhooksCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, webhooksTableName)).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllWebhooksCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetAllWebhooksCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	getAllWebhooksCountQueryBuilder sync.Once
	getAllWebhooksCountQuery        string
)

// buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership.
func (s *Sqlite) buildGetAllWebhooksCountQuery() string {
	getAllWebhooksCountQueryBuilder.Do(func() {
		var err error

		getAllWebhooksCountQuery, _, err = s.sqlBuilder.
			Select(fmt.Sprintf(countQuery, webhooksTableName)).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllWebhooksCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildBuildGetAllWebhooksCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	getAllWebhooksCountQueryBuilder sync.Once
	getAllWebhooksCountQuery        string
)

// buildGetAllWebhooksCountQuery returns a query which would return the count of webhooks regardless of ownership.
func (m *MariaDB) buildGetAllWebhooksCountQuery() string {
	getAllWebhooksCountQueryBuilder.Do(func() {
		var err error

		getAllWebhooksCountQuery, _, err = m.sqlBuilder.
			Select(fmt.Sprintf(countQuery, webhooksTableName)).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllWebhooksCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_build_GetAllWebhooksCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := build_GetAllWebhooksCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllWebhooksCount will fetch the count of every active webhook in the database.
func (p *Postgres) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllWebhooksCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := build_GetAllWebhooksCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllWebhooksCount will fetch the count of every active webhook in the database.
func (s *Sqlite) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	err = s.db.QueryRowContext(ctx, s.buildGetAllWebhooksCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := build_GetAllWebhooksCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllWebhooksCount will fetch the count of every active webhook in the database.
func (m *MariaDB) GetAllWebhooksCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllWebhooksCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetAllWebhooksQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetAllWebhooksQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	getAllWebhooksQueryBuilder sync.Once
	getAllWebhooksQuery        string
)

// buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership.
func (p *Postgres) buildGetAllWebhooksQuery() string {
	getAllWebhooksQueryBuilder.Do(func() {
		var err error

		getAllWebhooksQuery, _, err = p.sqlBuilder.
			Select(webhooksTableColumns...).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllWebhooksQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetAllWebhooksQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	getAllWebhooksQueryBuilder sync.Once
	getAllWebhooksQuery        string
)

// buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership.
func (s *Sqlite) buildGetAllWebhooksQuery() string {
	getAllWebhooksQueryBuilder.Do(func() {
		var err error

		getAllWebhooksQuery, _, err = s.sqlBuilder.
			Select(webhooksTableColumns...).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllWebhooksQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildBuildGetAllWebhooksQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	"sync"
)

var (
	getAllWebhooksQueryBuilder sync.Once
	getAllWebhooksQuery        string
)

// buildGetAllWebhooksQuery returns a SQL query which will return all webhooks, regardless of ownership.
func (m *MariaDB) buildGetAllWebhooksQuery() string {
	getAllWebhooksQueryBuilder.Do(func() {
		var err error

		getAllWebhooksQuery, _, err = m.sqlBuilder.
			Select(webhooksTableColumns...).
			From(webhooksTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllWebhooksQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_build_GetAllWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := build_GetAllWebhooks(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllWebhooks fetches a list of all webhooks from the database.
func (p *Postgres) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	rows, err := p.db.QueryContext(ctx, p.buildGetAllWebhooksQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for webhooks: %w", err)
	}

	list, err := p.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v1.WebhookList{
		Pagination: v1.Pagination{
			Page: 1,
		},
		Webhooks: list,
	}

	return x, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := build_GetAllWebhooks(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllWebhooks fetches a list of all webhooks from the database.
func (s *Sqlite) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	rows, err := s.db.QueryContext(ctx, s.buildGetAllWebhooksQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for webhooks: %w", err)
	}

	list, err := s.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v1.WebhookList{
		Pagination: v1.Pagination{
			Page: 1,
		},
		Webhooks: list,
	}

	return x, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := build_GetAllWebhooks(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllWebhooks fetches a list of all webhooks from the database.
func (m *MariaDB) GetAllWebhooks(ctx context.Context) (*v1.WebhookList, error) {
	rows, err := m.db.QueryContext(ctx, m.buildGetAllWebhooksQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for webhooks: %w", err)
	}

	list, err := m.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v1.WebhookList{
		Pagination: v1.Pagination{
			Page: 1,
		},
		Webhooks: list,
	}

	return x, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetWebhooksQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetWebhooksQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetWebhooksQuery returns a SQL query (and arguments) that would return a
func (p *Postgres) buildGetWebhooksQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn):             nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", webhooksTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, webhooksTableName)
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
		x := buildBuildGetWebhooksQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetWebhooksQuery returns a SQL query (and arguments) that would return a
func (s *Sqlite) buildGetWebhooksQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn):             nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", webhooksTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, webhooksTableName)
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
		x := buildBuildGetWebhooksQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetWebhooksQuery returns a SQL query (and arguments) that would return a
func (m *MariaDB) buildGetWebhooksQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(webhooksTableColumns...).
		From(webhooksTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", webhooksTableName, webhooksTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", webhooksTableName, archivedOnColumn):             nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", webhooksTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, webhooksTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetWebhooks(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetWebhooks(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (p *Postgres) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	query, args := p.buildGetWebhooksQuery(userID, filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}

	list, err := p.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v1.WebhookList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Webhooks: list,
	}

	return x, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetWebhooks(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (s *Sqlite) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	query, args := s.buildGetWebhooksQuery(userID, filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}

	list, err := s.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v1.WebhookList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Webhooks: list,
	}

	return x, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetWebhooks(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetWebhooks fetches a list of webhooks from the database that meet a particular filter.
func (m *MariaDB) GetWebhooks(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.WebhookList, error) {
	query, args := m.buildGetWebhooksQuery(userID, filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database: %w", err)
	}

	list, err := m.scanWebhooks(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	x := &v1.WebhookList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Webhooks: list,
	}

	return x, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildWebhookCreationQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildWebhookCreationQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook
func (p *Postgres) buildWebhookCreationQuery(x *v1.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(webhooksTableName).
		Columns(
			webhooksTableNameColumn,
			webhooksTableContentTypeColumn,
			webhooksTableURLColumn,
			webhooksTableMethodColumn,
			webhooksTableEventsColumn,
			webhooksTableDataTypesColumn,
			webhooksTableTopicsColumn,
			webhooksTableOwnershipColumn,
		).
		Values(
			x.Name,
			x.ContentType,
			x.URL,
			x.Method,
			strings.Join(x.Events, eventsSeparator),
			strings.Join(x.DataTypes, typesSeparator),
			strings.Join(x.Topics, topicsSeparator),
			x.BelongsToUser,
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
		x := buildBuildWebhookCreationQuery(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook
func (s *Sqlite) buildWebhookCreationQuery(x *v1.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(webhooksTableName).
		Columns(
			webhooksTableNameColumn,
			webhooksTableContentTypeColumn,
			webhooksTableURLColumn,
			webhooksTableMethodColumn,
			webhooksTableEventsColumn,
			webhooksTableDataTypesColumn,
			webhooksTableTopicsColumn,
			webhooksTableOwnershipColumn,
		).
		Values(
			x.Name,
			x.ContentType,
			x.URL,
			x.Method,
			strings.Join(x.Events, eventsSeparator),
			strings.Join(x.DataTypes, typesSeparator),
			strings.Join(x.Topics, topicsSeparator),
			x.BelongsToUser,
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
		x := buildBuildWebhookCreationQuery(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildWebhookCreationQuery returns a SQL query (and arguments) that would create a given webhook
func (m *MariaDB) buildWebhookCreationQuery(x *v1.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(webhooksTableName).
		Columns(
			webhooksTableNameColumn,
			webhooksTableContentTypeColumn,
			webhooksTableURLColumn,
			webhooksTableMethodColumn,
			webhooksTableEventsColumn,
			webhooksTableDataTypesColumn,
			webhooksTableTopicsColumn,
			webhooksTableOwnershipColumn,
		).
		Values(
			x.Name,
			x.ContentType,
			x.URL,
			x.Method,
			strings.Join(x.Events, eventsSeparator),
			strings.Join(x.DataTypes, typesSeparator),
			strings.Join(x.Topics, topicsSeparator),
			x.BelongsToUser,
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

func Test_buildCreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildCreateWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateWebhook creates a webhook in the database.
func (p *Postgres) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	x := &v1.Webhook{
		Name:          input.Name,
		ContentType:   input.ContentType,
		URL:           input.URL,
		Method:        input.Method,
		Events:        input.Events,
		DataTypes:     input.DataTypes,
		Topics:        input.Topics,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildWebhookCreationQuery(x)
	if err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn); err != nil {
		return nil, fmt.Errorf("error executing webhook creation query: %w", err)
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
		x := buildCreateWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateWebhook creates a webhook in the database.
func (s *Sqlite) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	x := &v1.Webhook{
		Name:          input.Name,
		ContentType:   input.ContentType,
		URL:           input.URL,
		Method:        input.Method,
		Events:        input.Events,
		DataTypes:     input.DataTypes,
		Topics:        input.Topics,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := s.buildWebhookCreationQuery(x)
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing webhook creation query: %w", err)
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
		x := buildCreateWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateWebhook creates a webhook in the database.
func (m *MariaDB) CreateWebhook(ctx context.Context, input *v1.WebhookCreationInput) (*v1.Webhook, error) {
	x := &v1.Webhook{
		Name:          input.Name,
		ContentType:   input.ContentType,
		URL:           input.URL,
		Method:        input.Method,
		Events:        input.Events,
		DataTypes:     input.DataTypes,
		Topics:        input.Topics,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := m.buildWebhookCreationQuery(x)
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing webhook creation query: %w", err)
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

func Test_buildBuildUpdateWebhookQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildUpdateWebhookQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update.
func (p *Postgres) buildUpdateWebhookQuery(input *v1.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(webhooksTableName).
		Set(webhooksTableNameColumn, input.Name).
		Set(webhooksTableContentTypeColumn, input.ContentType).
		Set(webhooksTableURLColumn, input.URL).
		Set(webhooksTableMethodColumn, input.Method).
		Set(webhooksTableEventsColumn, strings.Join(input.Events, topicsSeparator)).
		Set(webhooksTableDataTypesColumn, strings.Join(input.DataTypes, typesSeparator)).
		Set(webhooksTableTopicsColumn, strings.Join(input.Topics, topicsSeparator)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     input.ID,
			webhooksTableOwnershipColumn: input.BelongsToUser,
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
		x := buildBuildUpdateWebhookQuery(proj, dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update.
func (s *Sqlite) buildUpdateWebhookQuery(input *v1.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(webhooksTableName).
		Set(webhooksTableNameColumn, input.Name).
		Set(webhooksTableContentTypeColumn, input.ContentType).
		Set(webhooksTableURLColumn, input.URL).
		Set(webhooksTableMethodColumn, input.Method).
		Set(webhooksTableEventsColumn, strings.Join(input.Events, topicsSeparator)).
		Set(webhooksTableDataTypesColumn, strings.Join(input.DataTypes, typesSeparator)).
		Set(webhooksTableTopicsColumn, strings.Join(input.Topics, topicsSeparator)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     input.ID,
			webhooksTableOwnershipColumn: input.BelongsToUser,
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
		x := buildBuildUpdateWebhookQuery(proj, dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildUpdateWebhookQuery takes a given webhook and returns a SQL query to update.
func (m *MariaDB) buildUpdateWebhookQuery(input *v1.Webhook) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(webhooksTableName).
		Set(webhooksTableNameColumn, input.Name).
		Set(webhooksTableContentTypeColumn, input.ContentType).
		Set(webhooksTableURLColumn, input.URL).
		Set(webhooksTableMethodColumn, input.Method).
		Set(webhooksTableEventsColumn, strings.Join(input.Events, topicsSeparator)).
		Set(webhooksTableDataTypesColumn, strings.Join(input.DataTypes, typesSeparator)).
		Set(webhooksTableTopicsColumn, strings.Join(input.Topics, topicsSeparator)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     input.ID,
			webhooksTableOwnershipColumn: input.BelongsToUser,
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

func Test_buildUpdateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildUpdateWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID.
func (p *Postgres) UpdateWebhook(ctx context.Context, input *v1.Webhook) error {
	query, args := p.buildUpdateWebhookQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildUpdateWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID.
func (s *Sqlite) UpdateWebhook(ctx context.Context, input *v1.Webhook) error {
	query, args := s.buildUpdateWebhookQuery(input)
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
		x := buildUpdateWebhook(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateWebhook updates a particular webhook. Note that UpdateWebhook expects the provided input to have a valid ID.
func (m *MariaDB) UpdateWebhook(ctx context.Context, input *v1.Webhook) error {
	query, args := m.buildUpdateWebhookQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("panics on invalid database", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("invalid")
		proj := testprojects.BuildTodoApp()

		assert.Panics(t, func() { buildUpdateWebhook(proj, dbvendor) })
	})
}

func Test_buildBuildArchiveWebhookQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildArchiveWebhookQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.
func (p *Postgres) buildArchiveWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(webhooksTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     webhookID,
			webhooksTableOwnershipColumn: userID,
			archivedOnColumn:             nil,
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

		x := buildBuildArchiveWebhookQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.
func (s *Sqlite) buildArchiveWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(webhooksTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     webhookID,
			webhooksTableOwnershipColumn: userID,
			archivedOnColumn:             nil,
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

		x := buildBuildArchiveWebhookQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveWebhookQuery returns a SQL query (and arguments) that will mark a webhook as archived.
func (m *MariaDB) buildArchiveWebhookQuery(webhookID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(webhooksTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                     webhookID,
			webhooksTableOwnershipColumn: userID,
			archivedOnColumn:             nil,
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

func Test_buildArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildArchiveWebhook(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveWebhook archives a webhook from the database by its ID.
func (p *Postgres) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	query, args := p.buildArchiveWebhookQuery(webhookID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildArchiveWebhook(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveWebhook archives a webhook from the database by its ID.
func (s *Sqlite) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	query, args := s.buildArchiveWebhookQuery(webhookID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildArchiveWebhook(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveWebhook archives a webhook from the database by its ID.
func (m *MariaDB) ArchiveWebhook(ctx context.Context, webhookID, userID uint64) error {
	query, args := m.buildArchiveWebhookQuery(webhookID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
