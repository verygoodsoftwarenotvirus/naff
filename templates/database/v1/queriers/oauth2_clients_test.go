package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsDotGo(proj, dbvendor)

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
	scopesSeparator                      = ","
	oauth2ClientsTableName               = "oauth2_clients"
	oauth2ClientsTableNameColumn         = "name"
	oauth2ClientsTableClientIDColumn     = "client_id"
	oauth2ClientsTableScopesColumn       = "scopes"
	oauth2ClientsTableRedirectURIColumn  = "redirect_uri"
	oauth2ClientsTableClientSecretColumn = "client_secret"
	oauth2ClientsTableOwnershipColumn    = "belongs_to_user"
)

var (
	oauth2ClientsTableColumns = []string{
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableNameColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableScopesColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableRedirectURIColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientSecretColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn),
	}
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (p *Postgres) scanOAuth2Client(scan v1.Scanner) (*v11.OAuth2Client, error) {
	var (
		x      = &v11.OAuth2Client{}
		scopes string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ClientID,
		&scopes,
		&x.RedirectURI,
		&x.ClientSecret,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if scopes := strings.Split(scopes, scopesSeparator); len(scopes) >= 1 && scopes[0] != "" {
		x.Scopes = scopes
	}

	return x, nil
}

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients.
func (p *Postgres) scanOAuth2Clients(rows v1.ResultIterator) ([]*v11.OAuth2Client, error) {
	var (
		list []*v11.OAuth2Client
	)

	for rows.Next() {
		client, err := p.scanOAuth2Client(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		p.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID.
func (p *Postgres) buildGetOAuth2ClientByClientIDQuery(clientID string) (query string, args []interface{}) {
	var err error

	// This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't
	// care about ownership. It does still care about archived status
	query, args, err = p.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn): clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                 nil,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2ClientByClientID gets an OAuth2 client.
func (p *Postgres) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v11.OAuth2Client, error) {
	query, args := p.buildGetOAuth2ClientByClientIDQuery(clientID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanOAuth2Client(row)
}

var (
	getAllOAuth2ClientsQueryBuilder sync.Once
	getAllOAuth2ClientsQuery        string
)

// buildGetAllOAuth2ClientsQuery builds a SQL query.
func (p *Postgres) buildGetAllOAuth2ClientsQuery() (query string) {
	getAllOAuth2ClientsQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientsQuery, _, err = p.sqlBuilder.
			Select(oauth2ClientsTableColumns...).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientsQuery
}

// GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership.
func (p *Postgres) GetAllOAuth2Clients(ctx context.Context) ([]*v11.OAuth2Client, error) {
	rows, err := p.db.QueryContext(ctx, p.buildGetAllOAuth2ClientsQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := p.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}

// GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user.
func (p *Postgres) GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*v11.OAuth2Client, error) {
	query, args := p.buildGetOAuth2ClientsForUserQuery(userID, nil)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := p.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}

// buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID.
func (p *Postgres) buildGetOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn):                          clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2Client retrieves an OAuth2 client from the database.
func (p *Postgres) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v11.OAuth2Client, error) {
	query, args := p.buildGetOAuth2ClientQuery(clientID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	client, err := p.scanOAuth2Client(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 client: %w", err)
	}

	return client, nil
}

var (
	getAllOAuth2ClientCountQueryBuilder sync.Once
	getAllOAuth2ClientCountQuery        string
)

// buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients
// in the database, regardless of ownership.
func (p *Postgres) buildGetAllOAuth2ClientsCountQuery() string {
	getAllOAuth2ClientCountQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientCountQuery
}

// GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter.
func (p *Postgres) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := p.db.QueryRowContext(ctx, p.buildGetAllOAuth2ClientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that
// meet the given filter's criteria (if relevant) and belong to a given user.
func (p *Postgres) buildGetOAuth2ClientsForUserQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, oauth2ClientsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (p *Postgres) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.OAuth2ClientList, error) {
	query, args := p.buildGetOAuth2ClientsForUserQuery(userID, filter)
	rows, err := p.db.QueryContext(ctx, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 clients: %w", err)
	}

	list, err := p.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	ocl := &v11.OAuth2ClientList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	}

	// de-pointer-ize clients
	ocl.Clients = make([]v11.OAuth2Client, len(list))
	for i, t := range list {
		ocl.Clients[i] = *t
	}

	return ocl, nil
}

// buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database
func (p *Postgres) buildCreateOAuth2ClientQuery(input *v11.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(oauth2ClientsTableName).
		Columns(
			oauth2ClientsTableNameColumn,
			oauth2ClientsTableClientIDColumn,
			oauth2ClientsTableClientSecretColumn,
			oauth2ClientsTableScopesColumn,
			oauth2ClientsTableRedirectURIColumn,
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.ClientID,
			input.ClientSecret,
			strings.Join(input.Scopes, scopesSeparator),
			input.RedirectURI,
			input.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateOAuth2Client creates an OAuth2 client.
func (p *Postgres) CreateOAuth2Client(ctx context.Context, input *v11.OAuth2ClientCreationInput) (*v11.OAuth2Client, error) {
	x := &v11.OAuth2Client{
		Name:          input.Name,
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
	query, args := p.buildCreateOAuth2ClientQuery(x)

	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing client creation query: %w", err)
	}

	return x, nil
}

// buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database
func (p *Postgres) buildUpdateOAuth2ClientQuery(input *v11.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(oauth2ClientsTableClientIDColumn, input.ClientID).
		Set(oauth2ClientsTableClientSecretColumn, input.ClientSecret).
		Set(oauth2ClientsTableScopesColumn, strings.Join(input.Scopes, scopesSeparator)).
		Set(oauth2ClientsTableRedirectURIColumn, input.RedirectURI).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          input.ID,
			oauth2ClientsTableOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateOAuth2Client updates a OAuth2 client.
// NOTE: this function expects the input's ID field to be valid and non-zero.
func (p *Postgres) UpdateOAuth2Client(ctx context.Context, input *v11.OAuth2Client) error {
	query, args := p.buildUpdateOAuth2ClientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived.
func (p *Postgres) buildArchiveOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          clientID,
			oauth2ClientsTableOwnershipColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (p *Postgres) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	query, args := p.buildArchiveOAuth2ClientQuery(clientID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsDotGo(proj, dbvendor)

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
	scopesSeparator                      = ","
	oauth2ClientsTableName               = "oauth2_clients"
	oauth2ClientsTableNameColumn         = "name"
	oauth2ClientsTableClientIDColumn     = "client_id"
	oauth2ClientsTableScopesColumn       = "scopes"
	oauth2ClientsTableRedirectURIColumn  = "redirect_uri"
	oauth2ClientsTableClientSecretColumn = "client_secret"
	oauth2ClientsTableOwnershipColumn    = "belongs_to_user"
)

var (
	oauth2ClientsTableColumns = []string{
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableNameColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableScopesColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableRedirectURIColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientSecretColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn),
	}
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (s *Sqlite) scanOAuth2Client(scan v1.Scanner) (*v11.OAuth2Client, error) {
	var (
		x      = &v11.OAuth2Client{}
		scopes string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ClientID,
		&scopes,
		&x.RedirectURI,
		&x.ClientSecret,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if scopes := strings.Split(scopes, scopesSeparator); len(scopes) >= 1 && scopes[0] != "" {
		x.Scopes = scopes
	}

	return x, nil
}

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients.
func (s *Sqlite) scanOAuth2Clients(rows v1.ResultIterator) ([]*v11.OAuth2Client, error) {
	var (
		list []*v11.OAuth2Client
	)

	for rows.Next() {
		client, err := s.scanOAuth2Client(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		s.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID.
func (s *Sqlite) buildGetOAuth2ClientByClientIDQuery(clientID string) (query string, args []interface{}) {
	var err error

	// This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't
	// care about ownership. It does still care about archived status
	query, args, err = s.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn): clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                 nil,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2ClientByClientID gets an OAuth2 client.
func (s *Sqlite) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v11.OAuth2Client, error) {
	query, args := s.buildGetOAuth2ClientByClientIDQuery(clientID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return s.scanOAuth2Client(row)
}

var (
	getAllOAuth2ClientsQueryBuilder sync.Once
	getAllOAuth2ClientsQuery        string
)

// buildGetAllOAuth2ClientsQuery builds a SQL query.
func (s *Sqlite) buildGetAllOAuth2ClientsQuery() (query string) {
	getAllOAuth2ClientsQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientsQuery, _, err = s.sqlBuilder.
			Select(oauth2ClientsTableColumns...).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientsQuery
}

// GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership.
func (s *Sqlite) GetAllOAuth2Clients(ctx context.Context) ([]*v11.OAuth2Client, error) {
	rows, err := s.db.QueryContext(ctx, s.buildGetAllOAuth2ClientsQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := s.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}

// GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user.
func (s *Sqlite) GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*v11.OAuth2Client, error) {
	query, args := s.buildGetOAuth2ClientsForUserQuery(userID, nil)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := s.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}

// buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID.
func (s *Sqlite) buildGetOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn):                          clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2Client retrieves an OAuth2 client from the database.
func (s *Sqlite) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v11.OAuth2Client, error) {
	query, args := s.buildGetOAuth2ClientQuery(clientID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	client, err := s.scanOAuth2Client(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 client: %w", err)
	}

	return client, nil
}

var (
	getAllOAuth2ClientCountQueryBuilder sync.Once
	getAllOAuth2ClientCountQuery        string
)

// buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients
// in the database, regardless of ownership.
func (s *Sqlite) buildGetAllOAuth2ClientsCountQuery() string {
	getAllOAuth2ClientCountQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientCountQuery, _, err = s.sqlBuilder.
			Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientCountQuery
}

// GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter.
func (s *Sqlite) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := s.db.QueryRowContext(ctx, s.buildGetAllOAuth2ClientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that
// meet the given filter's criteria (if relevant) and belong to a given user.
func (s *Sqlite) buildGetOAuth2ClientsForUserQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, oauth2ClientsTableName)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (s *Sqlite) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.OAuth2ClientList, error) {
	query, args := s.buildGetOAuth2ClientsForUserQuery(userID, filter)
	rows, err := s.db.QueryContext(ctx, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 clients: %w", err)
	}

	list, err := s.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	ocl := &v11.OAuth2ClientList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	}

	// de-pointer-ize clients
	ocl.Clients = make([]v11.OAuth2Client, len(list))
	for i, t := range list {
		ocl.Clients[i] = *t
	}

	return ocl, nil
}

// buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database
func (s *Sqlite) buildCreateOAuth2ClientQuery(input *v11.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(oauth2ClientsTableName).
		Columns(
			oauth2ClientsTableNameColumn,
			oauth2ClientsTableClientIDColumn,
			oauth2ClientsTableClientSecretColumn,
			oauth2ClientsTableScopesColumn,
			oauth2ClientsTableRedirectURIColumn,
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.ClientID,
			input.ClientSecret,
			strings.Join(input.Scopes, scopesSeparator),
			input.RedirectURI,
			input.BelongsToUser,
		).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// CreateOAuth2Client creates an OAuth2 client.
func (s *Sqlite) CreateOAuth2Client(ctx context.Context, input *v11.OAuth2ClientCreationInput) (*v11.OAuth2Client, error) {
	x := &v11.OAuth2Client{
		Name:          input.Name,
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
	query, args := s.buildCreateOAuth2ClientQuery(x)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing client creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	s.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = s.timeTeller.Now()

	return x, nil
}

// buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database
func (s *Sqlite) buildUpdateOAuth2ClientQuery(input *v11.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(oauth2ClientsTableClientIDColumn, input.ClientID).
		Set(oauth2ClientsTableClientSecretColumn, input.ClientSecret).
		Set(oauth2ClientsTableScopesColumn, strings.Join(input.Scopes, scopesSeparator)).
		Set(oauth2ClientsTableRedirectURIColumn, input.RedirectURI).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          input.ID,
			oauth2ClientsTableOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateOAuth2Client updates a OAuth2 client.
// NOTE: this function expects the input's ID field to be valid and non-zero.
func (s *Sqlite) UpdateOAuth2Client(ctx context.Context, input *v11.OAuth2Client) error {
	query, args := s.buildUpdateOAuth2ClientQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived.
func (s *Sqlite) buildArchiveOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          clientID,
			oauth2ClientsTableOwnershipColumn: userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (s *Sqlite) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	query, args := s.buildArchiveOAuth2ClientQuery(clientID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := oauth2ClientsDotGo(proj, dbvendor)

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
	scopesSeparator                      = ","
	oauth2ClientsTableName               = "oauth2_clients"
	oauth2ClientsTableNameColumn         = "name"
	oauth2ClientsTableClientIDColumn     = "client_id"
	oauth2ClientsTableScopesColumn       = "scopes"
	oauth2ClientsTableRedirectURIColumn  = "redirect_uri"
	oauth2ClientsTableClientSecretColumn = "client_secret"
	oauth2ClientsTableOwnershipColumn    = "belongs_to_user"
)

var (
	oauth2ClientsTableColumns = []string{
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableNameColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableScopesColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableRedirectURIColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientSecretColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn),
	}
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (m *MariaDB) scanOAuth2Client(scan v1.Scanner) (*v11.OAuth2Client, error) {
	var (
		x      = &v11.OAuth2Client{}
		scopes string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ClientID,
		&scopes,
		&x.RedirectURI,
		&x.ClientSecret,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if scopes := strings.Split(scopes, scopesSeparator); len(scopes) >= 1 && scopes[0] != "" {
		x.Scopes = scopes
	}

	return x, nil
}

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients.
func (m *MariaDB) scanOAuth2Clients(rows v1.ResultIterator) ([]*v11.OAuth2Client, error) {
	var (
		list []*v11.OAuth2Client
	)

	for rows.Next() {
		client, err := m.scanOAuth2Client(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		m.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID.
func (m *MariaDB) buildGetOAuth2ClientByClientIDQuery(clientID string) (query string, args []interface{}) {
	var err error

	// This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't
	// care about ownership. It does still care about archived status
	query, args, err = m.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn): clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                 nil,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2ClientByClientID gets an OAuth2 client.
func (m *MariaDB) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v11.OAuth2Client, error) {
	query, args := m.buildGetOAuth2ClientByClientIDQuery(clientID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return m.scanOAuth2Client(row)
}

var (
	getAllOAuth2ClientsQueryBuilder sync.Once
	getAllOAuth2ClientsQuery        string
)

// buildGetAllOAuth2ClientsQuery builds a SQL query.
func (m *MariaDB) buildGetAllOAuth2ClientsQuery() (query string) {
	getAllOAuth2ClientsQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientsQuery, _, err = m.sqlBuilder.
			Select(oauth2ClientsTableColumns...).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientsQuery
}

// GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership.
func (m *MariaDB) GetAllOAuth2Clients(ctx context.Context) ([]*v11.OAuth2Client, error) {
	rows, err := m.db.QueryContext(ctx, m.buildGetAllOAuth2ClientsQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := m.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}

// GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user.
func (m *MariaDB) GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*v11.OAuth2Client, error) {
	query, args := m.buildGetOAuth2ClientsForUserQuery(userID, nil)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := m.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}

// buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID.
func (m *MariaDB) buildGetOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn):                          clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2Client retrieves an OAuth2 client from the database.
func (m *MariaDB) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v11.OAuth2Client, error) {
	query, args := m.buildGetOAuth2ClientQuery(clientID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	client, err := m.scanOAuth2Client(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 client: %w", err)
	}

	return client, nil
}

var (
	getAllOAuth2ClientCountQueryBuilder sync.Once
	getAllOAuth2ClientCountQuery        string
)

// buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients
// in the database, regardless of ownership.
func (m *MariaDB) buildGetAllOAuth2ClientsCountQuery() string {
	getAllOAuth2ClientCountQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientCountQuery, _, err = m.sqlBuilder.
			Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientCountQuery
}

// GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter.
func (m *MariaDB) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := m.db.QueryRowContext(ctx, m.buildGetAllOAuth2ClientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that
// meet the given filter's criteria (if relevant) and belong to a given user.
func (m *MariaDB) buildGetOAuth2ClientsForUserQuery(userID uint64, filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, oauth2ClientsTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (m *MariaDB) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v11.QueryFilter) (*v11.OAuth2ClientList, error) {
	query, args := m.buildGetOAuth2ClientsForUserQuery(userID, filter)
	rows, err := m.db.QueryContext(ctx, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 clients: %w", err)
	}

	list, err := m.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	ocl := &v11.OAuth2ClientList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	}

	// de-pointer-ize clients
	ocl.Clients = make([]v11.OAuth2Client, len(list))
	for i, t := range list {
		ocl.Clients[i] = *t
	}

	return ocl, nil
}

// buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database
func (m *MariaDB) buildCreateOAuth2ClientQuery(input *v11.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(oauth2ClientsTableName).
		Columns(
			oauth2ClientsTableNameColumn,
			oauth2ClientsTableClientIDColumn,
			oauth2ClientsTableClientSecretColumn,
			oauth2ClientsTableScopesColumn,
			oauth2ClientsTableRedirectURIColumn,
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.ClientID,
			input.ClientSecret,
			strings.Join(input.Scopes, scopesSeparator),
			input.RedirectURI,
			input.BelongsToUser,
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateOAuth2Client creates an OAuth2 client.
func (m *MariaDB) CreateOAuth2Client(ctx context.Context, input *v11.OAuth2ClientCreationInput) (*v11.OAuth2Client, error) {
	x := &v11.OAuth2Client{
		Name:          input.Name,
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
	query, args := m.buildCreateOAuth2ClientQuery(x)

	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing client creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	m.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = m.timeTeller.Now()

	return x, nil
}

// buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database
func (m *MariaDB) buildUpdateOAuth2ClientQuery(input *v11.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(oauth2ClientsTableClientIDColumn, input.ClientID).
		Set(oauth2ClientsTableClientSecretColumn, input.ClientSecret).
		Set(oauth2ClientsTableScopesColumn, strings.Join(input.Scopes, scopesSeparator)).
		Set(oauth2ClientsTableRedirectURIColumn, input.RedirectURI).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          input.ID,
			oauth2ClientsTableOwnershipColumn: input.BelongsToUser,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateOAuth2Client updates a OAuth2 client.
// NOTE: this function expects the input's ID field to be valid and non-zero.
func (m *MariaDB) UpdateOAuth2Client(ctx context.Context, input *v11.OAuth2Client) error {
	query, args := m.buildUpdateOAuth2ClientQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived.
func (m *MariaDB) buildArchiveOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          clientID,
			oauth2ClientsTableOwnershipColumn: userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveOAuth2Client archives an OAuth2 client.
func (m *MariaDB) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	query, args := m.buildArchiveOAuth2ClientQuery(clientID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientsConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2ClientsConstDeclarations()

		expected := `
package example

import ()

const (
	scopesSeparator                      = ","
	oauth2ClientsTableName               = "oauth2_clients"
	oauth2ClientsTableNameColumn         = "name"
	oauth2ClientsTableClientIDColumn     = "client_id"
	oauth2ClientsTableScopesColumn       = "scopes"
	oauth2ClientsTableRedirectURIColumn  = "redirect_uri"
	oauth2ClientsTableClientSecretColumn = "client_secret"
	oauth2ClientsTableOwnershipColumn    = "belongs_to_user"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOAuth2ClientsVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOAuth2ClientsVarDeclarations()

		expected := `
package example

import (
	"fmt"
)

var (
	oauth2ClientsTableColumns = []string{
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableNameColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableScopesColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableRedirectURIColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientSecretColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn),
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildScanOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (p *Postgres) scanOAuth2Client(scan v1.Scanner) (*v11.OAuth2Client, error) {
	var (
		x      = &v11.OAuth2Client{}
		scopes string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ClientID,
		&scopes,
		&x.RedirectURI,
		&x.ClientSecret,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if scopes := strings.Split(scopes, scopesSeparator); len(scopes) >= 1 && scopes[0] != "" {
		x.Scopes = scopes
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildScanOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (s *Sqlite) scanOAuth2Client(scan v1.Scanner) (*v11.OAuth2Client, error) {
	var (
		x      = &v11.OAuth2Client{}
		scopes string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ClientID,
		&scopes,
		&x.RedirectURI,
		&x.ClientSecret,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if scopes := strings.Split(scopes, scopesSeparator); len(scopes) >= 1 && scopes[0] != "" {
		x.Scopes = scopes
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildScanOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// scanOAuth2Client takes a Scanner (i.e. *sql.Row) and scans its results into an OAuth2Client struct.
func (m *MariaDB) scanOAuth2Client(scan v1.Scanner) (*v11.OAuth2Client, error) {
	var (
		x      = &v11.OAuth2Client{}
		scopes string
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.ClientID,
		&scopes,
		&x.RedirectURI,
		&x.ClientSecret,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	if scopes := strings.Split(scopes, scopesSeparator); len(scopes) >= 1 && scopes[0] != "" {
		x.Scopes = scopes
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildScanOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients.
func (p *Postgres) scanOAuth2Clients(rows v1.ResultIterator) ([]*v11.OAuth2Client, error) {
	var (
		list []*v11.OAuth2Client
	)

	for rows.Next() {
		client, err := p.scanOAuth2Client(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, client)
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildScanOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients.
func (s *Sqlite) scanOAuth2Clients(rows v1.ResultIterator) ([]*v11.OAuth2Client, error) {
	var (
		list []*v11.OAuth2Client
	)

	for rows.Next() {
		client, err := s.scanOAuth2Client(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, client)
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
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildScanOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanOAuth2Clients takes sql rows and turns them into a slice of OAuth2Clients.
func (m *MariaDB) scanOAuth2Clients(rows v1.ResultIterator) ([]*v11.OAuth2Client, error) {
	var (
		list []*v11.OAuth2Client
	)

	for rows.Next() {
		client, err := m.scanOAuth2Client(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, client)
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

func Test_buildBuildGetOAuth2ClientByClientIDQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetOAuth2ClientByClientIDQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID.
func (p *Postgres) buildGetOAuth2ClientByClientIDQuery(clientID string) (query string, args []interface{}) {
	var err error

	// This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't
	// care about ownership. It does still care about archived status
	query, args, err = p.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn): clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                 nil,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetOAuth2ClientByClientIDQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID.
func (s *Sqlite) buildGetOAuth2ClientByClientIDQuery(clientID string) (query string, args []interface{}) {
	var err error

	// This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't
	// care about ownership. It does still care about archived status
	query, args, err = s.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn): clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                 nil,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildGetOAuth2ClientByClientIDQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetOAuth2ClientByClientIDQuery builds a SQL query for fetching an OAuth2 client by its ClientID.
func (m *MariaDB) buildGetOAuth2ClientByClientIDQuery(clientID string) (query string, args []interface{}) {
	var err error

	// This query is more or less the same as the normal OAuth2 client retrieval query, only that it doesn't
	// care about ownership. It does still care about archived status
	query, args, err = m.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableClientIDColumn): clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                 nil,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientByClientID(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientByClientID gets an OAuth2 client.
func (p *Postgres) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v1.OAuth2Client, error) {
	query, args := p.buildGetOAuth2ClientByClientIDQuery(clientID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanOAuth2Client(row)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientByClientID(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientByClientID gets an OAuth2 client.
func (s *Sqlite) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v1.OAuth2Client, error) {
	query, args := s.buildGetOAuth2ClientByClientIDQuery(clientID)
	row := s.db.QueryRowContext(ctx, query, args...)
	return s.scanOAuth2Client(row)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientByClientID(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientByClientID gets an OAuth2 client.
func (m *MariaDB) GetOAuth2ClientByClientID(ctx context.Context, clientID string) (*v1.OAuth2Client, error) {
	query, args := m.buildGetOAuth2ClientByClientIDQuery(clientID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return m.scanOAuth2Client(row)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2ClientsQueryBuilderVarDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildGetAllOAuth2ClientsQueryBuilderVarDecls()

		expected := `
package example

import (
	"sync"
)

var (
	getAllOAuth2ClientsQueryBuilder sync.Once
	getAllOAuth2ClientsQuery        string
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetAllOAuth2ClientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetAllOAuth2ClientsQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllOAuth2ClientsQuery builds a SQL query.
func (p *Postgres) buildGetAllOAuth2ClientsQuery() (query string) {
	getAllOAuth2ClientsQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientsQuery, _, err = p.sqlBuilder.
			Select(oauth2ClientsTableColumns...).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientsQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetAllOAuth2ClientsQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllOAuth2ClientsQuery builds a SQL query.
func (s *Sqlite) buildGetAllOAuth2ClientsQuery() (query string) {
	getAllOAuth2ClientsQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientsQuery, _, err = s.sqlBuilder.
			Select(oauth2ClientsTableColumns...).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientsQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildGetAllOAuth2ClientsQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllOAuth2ClientsQuery builds a SQL query.
func (m *MariaDB) buildGetAllOAuth2ClientsQuery() (query string) {
	getAllOAuth2ClientsQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientsQuery, _, err = m.sqlBuilder.
			Select(oauth2ClientsTableColumns...).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientsQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership.
func (p *Postgres) GetAllOAuth2Clients(ctx context.Context) ([]*v1.OAuth2Client, error) {
	rows, err := p.db.QueryContext(ctx, p.buildGetAllOAuth2ClientsQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := p.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership.
func (s *Sqlite) GetAllOAuth2Clients(ctx context.Context) ([]*v1.OAuth2Client, error) {
	rows, err := s.db.QueryContext(ctx, s.buildGetAllOAuth2ClientsQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := s.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2Clients(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllOAuth2Clients gets a list of OAuth2 clients regardless of ownership.
func (m *MariaDB) GetAllOAuth2Clients(ctx context.Context) ([]*v1.OAuth2Client, error) {
	rows, err := m.db.QueryContext(ctx, m.buildGetAllOAuth2ClientsQuery())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := m.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user.
func (p *Postgres) GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*v1.OAuth2Client, error) {
	query, args := p.buildGetOAuth2ClientsForUserQuery(userID, nil)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := p.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user.
func (s *Sqlite) GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*v1.OAuth2Client, error) {
	query, args := s.buildGetOAuth2ClientsForUserQuery(userID, nil)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := s.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetAllOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetAllOAuth2ClientsForUser gets a list of OAuth2 clients belonging to a given user.
func (m *MariaDB) GetAllOAuth2ClientsForUser(ctx context.Context, userID uint64) ([]*v1.OAuth2Client, error) {
	query, args := m.buildGetOAuth2ClientsForUserQuery(userID, nil)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying database for oauth2 clients: %w", err)
	}

	list, err := m.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("fetching list of OAuth2Clients: %w", err)
	}

	return list, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetOAuth2ClientQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID.
func (p *Postgres) buildGetOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn):                          clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetOAuth2ClientQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID.
func (s *Sqlite) buildGetOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn):                          clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).ToSql()

	s.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildGetOAuth2ClientQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetOAuth2ClientQuery returns a SQL query which requests a given OAuth2 client by its database ID.
func (m *MariaDB) buildGetOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn):                          clientID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2Client retrieves an OAuth2 client from the database.
func (p *Postgres) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	query, args := p.buildGetOAuth2ClientQuery(clientID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	client, err := p.scanOAuth2Client(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 client: %w", err)
	}

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2Client retrieves an OAuth2 client from the database.
func (s *Sqlite) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	query, args := s.buildGetOAuth2ClientQuery(clientID, userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	client, err := s.scanOAuth2Client(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 client: %w", err)
	}

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2Client retrieves an OAuth2 client from the database.
func (m *MariaDB) GetOAuth2Client(ctx context.Context, clientID, userID uint64) (*v1.OAuth2Client, error) {
	query, args := m.buildGetOAuth2ClientQuery(clientID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	client, err := m.scanOAuth2Client(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 client: %w", err)
	}

	return client, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2ClientCountQueryBuilderVarDecls(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildGetAllOAuth2ClientCountQueryBuilderVarDecls()

		expected := `
package example

import (
	"sync"
)

var (
	getAllOAuth2ClientCountQueryBuilder sync.Once
	getAllOAuth2ClientCountQuery        string
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetAllOAuth2ClientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetAllOAuth2ClientsCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients
// in the database, regardless of ownership.
func (p *Postgres) buildGetAllOAuth2ClientsCountQuery() string {
	getAllOAuth2ClientCountQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		p.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetAllOAuth2ClientsCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients
// in the database, regardless of ownership.
func (s *Sqlite) buildGetAllOAuth2ClientsCountQuery() string {
	getAllOAuth2ClientCountQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientCountQuery, _, err = s.sqlBuilder.
			Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		s.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildGetAllOAuth2ClientsCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllOAuth2ClientsCountQuery returns a SQL query for the number of OAuth2 clients
// in the database, regardless of ownership.
func (m *MariaDB) buildGetAllOAuth2ClientsCountQuery() string {
	getAllOAuth2ClientCountQueryBuilder.Do(func() {
		var err error

		getAllOAuth2ClientCountQuery, _, err = m.sqlBuilder.
			Select(fmt.Sprintf(countQuery, oauth2ClientsTableName)).
			From(oauth2ClientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn): nil,
			}).
			ToSql()

		m.logQueryBuildingError(err)
	})

	return getAllOAuth2ClientCountQuery
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllOAuth2ClientCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildGetAllOAuth2ClientCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter.
func (p *Postgres) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := p.db.QueryRowContext(ctx, p.buildGetAllOAuth2ClientsCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildGetAllOAuth2ClientCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter.
func (s *Sqlite) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := s.db.QueryRowContext(ctx, s.buildGetAllOAuth2ClientsCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildGetAllOAuth2ClientCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllOAuth2ClientCount will get the count of OAuth2 clients that match the current filter.
func (m *MariaDB) GetAllOAuth2ClientCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := m.db.QueryRowContext(ctx, m.buildGetAllOAuth2ClientsCountQuery()).Scan(&count)
	return count, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetOAuth2ClientsForUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetOAuth2ClientsForUserQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that
// meet the given filter's criteria (if relevant) and belong to a given user.
func (p *Postgres) buildGetOAuth2ClientsForUserQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, oauth2ClientsTableName)
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetOAuth2ClientsForUserQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that
// meet the given filter's criteria (if relevant) and belong to a given user.
func (s *Sqlite) buildGetOAuth2ClientsForUserQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, oauth2ClientsTableName)
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
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetOAuth2ClientsForUserQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetOAuth2ClientsForUserQuery returns a SQL query (and arguments) that will retrieve a list of OAuth2 clients that
// meet the given filter's criteria (if relevant) and belong to a given user.
func (m *MariaDB) buildGetOAuth2ClientsForUserQuery(userID uint64, filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(oauth2ClientsTableColumns...).
		From(oauth2ClientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, oauth2ClientsTableOwnershipColumn): userID,
			fmt.Sprintf("%s.%s", oauth2ClientsTableName, archivedOnColumn):                  nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", oauth2ClientsTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, oauth2ClientsTableName)
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

func Test_buildGetOAuth2ClientsForUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (p *Postgres) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	query, args := p.buildGetOAuth2ClientsForUserQuery(userID, filter)
	rows, err := p.db.QueryContext(ctx, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 clients: %w", err)
	}

	list, err := p.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	ocl := &v1.OAuth2ClientList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	}

	// de-pointer-ize clients
	ocl.Clients = make([]v1.OAuth2Client, len(list))
	for i, t := range list {
		ocl.Clients[i] = *t
	}

	return ocl, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (s *Sqlite) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	query, args := s.buildGetOAuth2ClientsForUserQuery(userID, filter)
	rows, err := s.db.QueryContext(ctx, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 clients: %w", err)
	}

	list, err := s.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	ocl := &v1.OAuth2ClientList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	}

	// de-pointer-ize clients
	ocl.Clients = make([]v1.OAuth2Client, len(list))
	for i, t := range list {
		ocl.Clients[i] = *t
	}

	return ocl, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetOAuth2ClientsForUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetOAuth2ClientsForUser gets a list of OAuth2 clients.
func (m *MariaDB) GetOAuth2ClientsForUser(ctx context.Context, userID uint64, filter *v1.QueryFilter) (*v1.OAuth2ClientList, error) {
	query, args := m.buildGetOAuth2ClientsForUserQuery(userID, filter)
	rows, err := m.db.QueryContext(ctx, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("querying for oauth2 clients: %w", err)
	}

	list, err := m.scanOAuth2Clients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	ocl := &v1.OAuth2ClientList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	}

	// de-pointer-ize clients
	ocl.Clients = make([]v1.OAuth2Client, len(list))
	for i, t := range list {
		ocl.Clients[i] = *t
	}

	return ocl, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database
func (p *Postgres) buildCreateOAuth2ClientQuery(input *v1.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(oauth2ClientsTableName).
		Columns(
			oauth2ClientsTableNameColumn,
			oauth2ClientsTableClientIDColumn,
			oauth2ClientsTableClientSecretColumn,
			oauth2ClientsTableScopesColumn,
			oauth2ClientsTableRedirectURIColumn,
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.ClientID,
			input.ClientSecret,
			strings.Join(input.Scopes, scopesSeparator),
			input.RedirectURI,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database
func (s *Sqlite) buildCreateOAuth2ClientQuery(input *v1.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(oauth2ClientsTableName).
		Columns(
			oauth2ClientsTableNameColumn,
			oauth2ClientsTableClientIDColumn,
			oauth2ClientsTableClientSecretColumn,
			oauth2ClientsTableScopesColumn,
			oauth2ClientsTableRedirectURIColumn,
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.ClientID,
			input.ClientSecret,
			strings.Join(input.Scopes, scopesSeparator),
			input.RedirectURI,
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
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildCreateOAuth2ClientQuery returns a SQL query (and args) that will create the given OAuth2Client in the database
func (m *MariaDB) buildCreateOAuth2ClientQuery(input *v1.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(oauth2ClientsTableName).
		Columns(
			oauth2ClientsTableNameColumn,
			oauth2ClientsTableClientIDColumn,
			oauth2ClientsTableClientSecretColumn,
			oauth2ClientsTableScopesColumn,
			oauth2ClientsTableRedirectURIColumn,
			oauth2ClientsTableOwnershipColumn,
		).
		Values(
			input.Name,
			input.ClientID,
			input.ClientSecret,
			strings.Join(input.Scopes, scopesSeparator),
			input.RedirectURI,
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

func Test_buildCreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildCreateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateOAuth2Client creates an OAuth2 client.
func (p *Postgres) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	x := &v1.OAuth2Client{
		Name:          input.Name,
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
	query, args := p.buildCreateOAuth2ClientQuery(x)

	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing client creation query: %w", err)
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildCreateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateOAuth2Client creates an OAuth2 client.
func (s *Sqlite) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	x := &v1.OAuth2Client{
		Name:          input.Name,
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
	query, args := s.buildCreateOAuth2ClientQuery(x)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing client creation query: %w", err)
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
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildCreateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateOAuth2Client creates an OAuth2 client.
func (m *MariaDB) CreateOAuth2Client(ctx context.Context, input *v1.OAuth2ClientCreationInput) (*v1.OAuth2Client, error) {
	x := &v1.OAuth2Client{
		Name:          input.Name,
		ClientID:      input.ClientID,
		ClientSecret:  input.ClientSecret,
		RedirectURI:   input.RedirectURI,
		Scopes:        input.Scopes,
		BelongsToUser: input.BelongsToUser,
	}
	query, args := m.buildCreateOAuth2ClientQuery(x)

	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing client creation query: %w", err)
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

func Test_buildBuildUpdateOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildUpdateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database
func (p *Postgres) buildUpdateOAuth2ClientQuery(input *v1.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(oauth2ClientsTableClientIDColumn, input.ClientID).
		Set(oauth2ClientsTableClientSecretColumn, input.ClientSecret).
		Set(oauth2ClientsTableScopesColumn, strings.Join(input.Scopes, scopesSeparator)).
		Set(oauth2ClientsTableRedirectURIColumn, input.RedirectURI).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          input.ID,
			oauth2ClientsTableOwnershipColumn: input.BelongsToUser,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildBuildUpdateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database
func (s *Sqlite) buildUpdateOAuth2ClientQuery(input *v1.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(oauth2ClientsTableClientIDColumn, input.ClientID).
		Set(oauth2ClientsTableClientSecretColumn, input.ClientSecret).
		Set(oauth2ClientsTableScopesColumn, strings.Join(input.Scopes, scopesSeparator)).
		Set(oauth2ClientsTableRedirectURIColumn, input.RedirectURI).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          input.ID,
			oauth2ClientsTableOwnershipColumn: input.BelongsToUser,
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
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildBuildUpdateOAuth2ClientQuery(proj, dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"strings"
)

// buildUpdateOAuth2ClientQuery returns a SQL query (and args) that will update a given OAuth2 client in the database
func (m *MariaDB) buildUpdateOAuth2ClientQuery(input *v1.OAuth2Client) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(oauth2ClientsTableClientIDColumn, input.ClientID).
		Set(oauth2ClientsTableClientSecretColumn, input.ClientSecret).
		Set(oauth2ClientsTableScopesColumn, strings.Join(input.Scopes, scopesSeparator)).
		Set(oauth2ClientsTableRedirectURIColumn, input.RedirectURI).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          input.ID,
			oauth2ClientsTableOwnershipColumn: input.BelongsToUser,
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

func Test_buildUpdateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildUpdateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateOAuth2Client updates a OAuth2 client.
// NOTE: this function expects the input's ID field to be valid and non-zero.
func (p *Postgres) UpdateOAuth2Client(ctx context.Context, input *v1.OAuth2Client) error {
	query, args := p.buildUpdateOAuth2ClientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildUpdateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateOAuth2Client updates a OAuth2 client.
// NOTE: this function expects the input's ID field to be valid and non-zero.
func (s *Sqlite) UpdateOAuth2Client(ctx context.Context, input *v1.OAuth2Client) error {
	query, args := s.buildUpdateOAuth2ClientQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildUpdateOAuth2Client(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateOAuth2Client updates a OAuth2 client.
// NOTE: this function expects the input's ID field to be valid and non-zero.
func (m *MariaDB) UpdateOAuth2Client(ctx context.Context, input *v1.OAuth2Client) error {
	query, args := m.buildUpdateOAuth2ClientQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveOAuth2ClientQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildArchiveOAuth2ClientQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived.
func (p *Postgres) buildArchiveOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          clientID,
			oauth2ClientsTableOwnershipColumn: userID,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildArchiveOAuth2ClientQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived.
func (s *Sqlite) buildArchiveOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          clientID,
			oauth2ClientsTableOwnershipColumn: userID,
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
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildArchiveOAuth2ClientQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveOAuth2ClientQuery returns a SQL query (and arguments) that will mark an OAuth2 client as archived.
func (m *MariaDB) buildArchiveOAuth2ClientQuery(clientID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(oauth2ClientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                          clientID,
			oauth2ClientsTableOwnershipColumn: userID,
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

func Test_buildArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildArchiveOAuth2Client(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveOAuth2Client archives an OAuth2 client.
func (p *Postgres) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	query, args := p.buildArchiveOAuth2ClientQuery(clientID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildArchiveOAuth2Client(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveOAuth2Client archives an OAuth2 client.
func (s *Sqlite) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	query, args := s.buildArchiveOAuth2ClientQuery(clientID, userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildArchiveOAuth2Client(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveOAuth2Client archives an OAuth2 client.
func (m *MariaDB) ArchiveOAuth2Client(ctx context.Context, clientID, userID uint64) error {
	query, args := m.buildArchiveOAuth2ClientQuery(clientID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
