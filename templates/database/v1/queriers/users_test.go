package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := usersDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	pq "github.com/lib/pq"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1/client"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	usersTableName                         = "users"
	usersTableUsernameColumn               = "username"
	usersTableHashedPasswordColumn         = "hashed_password"
	usersTableSaltColumn                   = "salt"
	usersTableRequiresPasswordChangeColumn = "requires_password_change"
	usersTablePasswordLastChangedOnColumn  = "password_last_changed_on"
	usersTableTwoFactorColumn              = "two_factor_secret"
	usersTableTwoFactorVerifiedOnColumn    = "two_factor_secret_verified_on"
	usersTableIsAdminColumn                = "is_admin"
)

var (
	usersTableColumns = []string{
		fmt.Sprintf("%s.%s", usersTableName, idColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableHashedPasswordColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableSaltColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableRequiresPasswordChangeColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTablePasswordLastChangedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableIsAdminColumn),
		fmt.Sprintf("%s.%s", usersTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn),
	}
)

// scanUser provides a consistent way to scan something like a *sql.Row into a User struct.
func (p *Postgres) scanUser(scan v1.Scanner) (*v11.User, error) {
	var (
		x = &v11.User{}
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Username,
		&x.HashedPassword,
		&x.Salt,
		&x.RequiresPasswordChange,
		&x.PasswordLastChangedOn,
		&x.TwoFactorSecret,
		&x.TwoFactorSecretVerifiedOn,
		&x.IsAdmin,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanUsers takes database rows and loads them into a slice of User structs.
func (p *Postgres) scanUsers(rows v1.ResultIterator) ([]v11.User, error) {
	var (
		list []v11.User
	)

	for rows.Next() {
		user, err := p.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user result: %w", err)
		}

		list = append(list, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		p.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID
func (p *Postgres) buildGetUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):         userID,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetUser fetches a user.
func (p *Postgres) GetUser(ctx context.Context, userID uint64) (*v11.User, error) {
	query, args := p.buildGetUserQuery(userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	u, err := p.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}

// buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user
// by their database ID, who has an unverified two factor secret
func (p *Postgres) buildGetUserWithUnverifiedTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):                            userID,
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):                    nil,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret
func (p *Postgres) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v11.User, error) {
	query, args := p.buildGetUserWithUnverifiedTwoFactorSecretQuery(userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	u, err := p.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}

// buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username
func (p *Postgres) buildGetUserByUsernameQuery(username string) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn): username,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):         nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetUserByUsername fetches a user by their username.
func (p *Postgres) GetUserByUsername(ctx context.Context, username string) (*v11.User, error) {
	query, args := p.buildGetUserByUsernameQuery(username)
	row := p.db.QueryRowContext(ctx, query, args...)

	u, err := p.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return u, nil
}

// buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere
// to a given filter's criteria.
func (p *Postgres) buildGetAllUsersCountQuery() (query string) {
	var err error

	builder := p.sqlBuilder.
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		})

	query, _, err = builder.ToSql()

	p.logQueryBuildingError(err)

	return query
}

// GetAllUsersCount fetches a count of users from the database.
func (p *Postgres) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	query := p.buildGetAllUsersCountQuery()
	err = p.db.QueryRowContext(ctx, query).Scan(&count)
	return
}

// buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere
// to a given filter's criteria.
func (p *Postgres) buildGetUsersQuery(filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", usersTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, usersTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)
	return query, args
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (p *Postgres) GetUsers(ctx context.Context, filter *v11.QueryFilter) (*v11.UserList, error) {
	query, args := p.buildGetUsersQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying for user")
	}

	userList, err := p.scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("loading response from database: %w", err)
	}

	x := &v11.UserList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Users: userList,
	}

	return x, nil
}

// buildCreateUserQuery returns a SQL query (and arguments) that would create a given User
func (p *Postgres) buildCreateUserQuery(input v11.UserDatabaseCreationInput) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(usersTableName).
		Columns(
			usersTableUsernameColumn,
			usersTableHashedPasswordColumn,
			usersTableSaltColumn,
			usersTableTwoFactorColumn,
			usersTableIsAdminColumn,
		).
		Values(
			input.Username,
			input.HashedPassword,
			input.Salt,
			input.TwoFactorSecret,
			false,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	// NOTE: we always default is_admin to false, on the assumption that
	// admins have DB access and will change that value via SQL query.
	// There should also be no way to update a user via this structure
	// such that they would have admin privileges.

	p.logQueryBuildingError(err)

	return query, args
}

// CreateUser creates a user.
func (p *Postgres) CreateUser(ctx context.Context, input v11.UserDatabaseCreationInput) (*v11.User, error) {
	x := &v11.User{
		Username:        input.Username,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
	}
	query, args := p.buildCreateUserQuery(input)

	// create the user.
	if err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn); err != nil {
		switch e := err.(type) {
		case *pq.Error:
			if e.Code == pq.ErrorCode(postgresRowExistsErrorCode) {
				return nil, client.ErrUserExists
			}
		default:
			return nil, fmt.Errorf("error executing user creation query: %w", err)
		}
	}

	return x, nil
}

// buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row
func (p *Postgres) buildUpdateUserQuery(input *v11.User) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(usersTableUsernameColumn, input.Username).
		Set(usersTableHashedPasswordColumn, input.HashedPassword).
		Set(usersTableSaltColumn, input.Salt).
		Set(usersTableTwoFactorColumn, input.TwoFactorSecret).
		Set(usersTableTwoFactorVerifiedOnColumn, input.TwoFactorSecretVerifiedOn).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateUser receives a complete User struct and updates its place in the db.
// NOTE this function uses the ID provided in the input to make its query. Pass in
// incomplete models at your peril.
func (p *Postgres) UpdateUser(ctx context.Context, input *v11.User) error {
	query, args := p.buildUpdateUserQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password.
func (p *Postgres) buildUpdateUserPasswordQuery(userID uint64, newHash string) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(usersTableHashedPasswordColumn, newHash).
		Set(usersTableRequiresPasswordChangeColumn, false).
		Set(usersTablePasswordLastChangedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateUserPassword updates a user's password.
func (p *Postgres) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	query, args := p.buildUpdateUserPasswordQuery(userID, newHash)

	_, err := p.db.ExecContext(ctx, query, args...)

	return err
}

// buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret
func (p *Postgres) buildVerifyUserTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(usersTableTwoFactorVerifiedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (p *Postgres) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	query, args := p.buildVerifyUserTwoFactorSecretQuery(userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveUserQuery builds a SQL query that marks a user as archived.
func (p *Postgres) buildArchiveUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveUser marks a user as archived.
func (p *Postgres) ArchiveUser(ctx context.Context, userID uint64) error {
	query, args := p.buildArchiveUserQuery(userID)
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
		x := usersDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	usersTableName                         = "users"
	usersTableUsernameColumn               = "username"
	usersTableHashedPasswordColumn         = "hashed_password"
	usersTableSaltColumn                   = "salt"
	usersTableRequiresPasswordChangeColumn = "requires_password_change"
	usersTablePasswordLastChangedOnColumn  = "password_last_changed_on"
	usersTableTwoFactorColumn              = "two_factor_secret"
	usersTableTwoFactorVerifiedOnColumn    = "two_factor_secret_verified_on"
	usersTableIsAdminColumn                = "is_admin"
)

var (
	usersTableColumns = []string{
		fmt.Sprintf("%s.%s", usersTableName, idColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableHashedPasswordColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableSaltColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableRequiresPasswordChangeColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTablePasswordLastChangedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableIsAdminColumn),
		fmt.Sprintf("%s.%s", usersTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn),
	}
)

// scanUser provides a consistent way to scan something like a *sql.Row into a User struct.
func (s *Sqlite) scanUser(scan v1.Scanner) (*v11.User, error) {
	var (
		x = &v11.User{}
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Username,
		&x.HashedPassword,
		&x.Salt,
		&x.RequiresPasswordChange,
		&x.PasswordLastChangedOn,
		&x.TwoFactorSecret,
		&x.TwoFactorSecretVerifiedOn,
		&x.IsAdmin,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanUsers takes database rows and loads them into a slice of User structs.
func (s *Sqlite) scanUsers(rows v1.ResultIterator) ([]v11.User, error) {
	var (
		list []v11.User
	)

	for rows.Next() {
		user, err := s.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user result: %w", err)
		}

		list = append(list, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		s.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID
func (s *Sqlite) buildGetUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):         userID,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetUser fetches a user.
func (s *Sqlite) GetUser(ctx context.Context, userID uint64) (*v11.User, error) {
	query, args := s.buildGetUserQuery(userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	u, err := s.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}

// buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user
// by their database ID, who has an unverified two factor secret
func (s *Sqlite) buildGetUserWithUnverifiedTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):                            userID,
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):                    nil,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret
func (s *Sqlite) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v11.User, error) {
	query, args := s.buildGetUserWithUnverifiedTwoFactorSecretQuery(userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	u, err := s.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}

// buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username
func (s *Sqlite) buildGetUserByUsernameQuery(username string) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn): username,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):         nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// GetUserByUsername fetches a user by their username.
func (s *Sqlite) GetUserByUsername(ctx context.Context, username string) (*v11.User, error) {
	query, args := s.buildGetUserByUsernameQuery(username)
	row := s.db.QueryRowContext(ctx, query, args...)

	u, err := s.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return u, nil
}

// buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere
// to a given filter's criteria.
func (s *Sqlite) buildGetAllUsersCountQuery() (query string) {
	var err error

	builder := s.sqlBuilder.
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		})

	query, _, err = builder.ToSql()

	s.logQueryBuildingError(err)

	return query
}

// GetAllUsersCount fetches a count of users from the database.
func (s *Sqlite) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	query := s.buildGetAllUsersCountQuery()
	err = s.db.QueryRowContext(ctx, query).Scan(&count)
	return
}

// buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere
// to a given filter's criteria.
func (s *Sqlite) buildGetUsersQuery(filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", usersTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, usersTableName)
	}

	query, args, err = builder.ToSql()
	s.logQueryBuildingError(err)
	return query, args
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (s *Sqlite) GetUsers(ctx context.Context, filter *v11.QueryFilter) (*v11.UserList, error) {
	query, args := s.buildGetUsersQuery(filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying for user")
	}

	userList, err := s.scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("loading response from database: %w", err)
	}

	x := &v11.UserList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Users: userList,
	}

	return x, nil
}

// buildCreateUserQuery returns a SQL query (and arguments) that would create a given User
func (s *Sqlite) buildCreateUserQuery(input v11.UserDatabaseCreationInput) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(usersTableName).
		Columns(
			usersTableUsernameColumn,
			usersTableHashedPasswordColumn,
			usersTableSaltColumn,
			usersTableTwoFactorColumn,
			usersTableIsAdminColumn,
		).
		Values(
			input.Username,
			input.HashedPassword,
			input.Salt,
			input.TwoFactorSecret,
			false,
		).
		ToSql()

	// NOTE: we always default is_admin to false, on the assumption that
	// admins have DB access and will change that value via SQL query.
	// There should also be no way to update a user via this structure
	// such that they would have admin privileges.

	s.logQueryBuildingError(err)

	return query, args
}

// CreateUser creates a user.
func (s *Sqlite) CreateUser(ctx context.Context, input v11.UserDatabaseCreationInput) (*v11.User, error) {
	x := &v11.User{
		Username:        input.Username,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
	}
	query, args := s.buildCreateUserQuery(input)

	// create the user.
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing user creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	s.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = s.timeTeller.Now()

	return x, nil
}

// buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row
func (s *Sqlite) buildUpdateUserQuery(input *v11.User) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(usersTableUsernameColumn, input.Username).
		Set(usersTableHashedPasswordColumn, input.HashedPassword).
		Set(usersTableSaltColumn, input.Salt).
		Set(usersTableTwoFactorColumn, input.TwoFactorSecret).
		Set(usersTableTwoFactorVerifiedOnColumn, input.TwoFactorSecretVerifiedOn).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateUser receives a complete User struct and updates its place in the db.
// NOTE this function uses the ID provided in the input to make its query. Pass in
// incomplete models at your peril.
func (s *Sqlite) UpdateUser(ctx context.Context, input *v11.User) error {
	query, args := s.buildUpdateUserQuery(input)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password.
func (s *Sqlite) buildUpdateUserPasswordQuery(userID uint64, newHash string) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(usersTableHashedPasswordColumn, newHash).
		Set(usersTableRequiresPasswordChangeColumn, false).
		Set(usersTablePasswordLastChangedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// UpdateUserPassword updates a user's password.
func (s *Sqlite) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	query, args := s.buildUpdateUserPasswordQuery(userID, newHash)

	_, err := s.db.ExecContext(ctx, query, args...)

	return err
}

// buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret
func (s *Sqlite) buildVerifyUserTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(usersTableTwoFactorVerifiedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (s *Sqlite) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	query, args := s.buildVerifyUserTwoFactorSecretQuery(userID)
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveUserQuery builds a SQL query that marks a user as archived.
func (s *Sqlite) buildArchiveUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	s.logQueryBuildingError(err)

	return query, args
}

// ArchiveUser marks a user as archived.
func (s *Sqlite) ArchiveUser(ctx context.Context, userID uint64) error {
	query, args := s.buildArchiveUserQuery(userID)
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
		x := usersDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

const (
	usersTableName                         = "users"
	usersTableUsernameColumn               = "username"
	usersTableHashedPasswordColumn         = "hashed_password"
	usersTableSaltColumn                   = "salt"
	usersTableRequiresPasswordChangeColumn = "requires_password_change"
	usersTablePasswordLastChangedOnColumn  = "password_last_changed_on"
	usersTableTwoFactorColumn              = "two_factor_secret"
	usersTableTwoFactorVerifiedOnColumn    = "two_factor_secret_verified_on"
	usersTableIsAdminColumn                = "is_admin"
)

var (
	usersTableColumns = []string{
		fmt.Sprintf("%s.%s", usersTableName, idColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableHashedPasswordColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableSaltColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableRequiresPasswordChangeColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTablePasswordLastChangedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableIsAdminColumn),
		fmt.Sprintf("%s.%s", usersTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn),
	}
)

// scanUser provides a consistent way to scan something like a *sql.Row into a User struct.
func (m *MariaDB) scanUser(scan v1.Scanner) (*v11.User, error) {
	var (
		x = &v11.User{}
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Username,
		&x.HashedPassword,
		&x.Salt,
		&x.RequiresPasswordChange,
		&x.PasswordLastChangedOn,
		&x.TwoFactorSecret,
		&x.TwoFactorSecretVerifiedOn,
		&x.IsAdmin,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanUsers takes database rows and loads them into a slice of User structs.
func (m *MariaDB) scanUsers(rows v1.ResultIterator) ([]v11.User, error) {
	var (
		list []v11.User
	)

	for rows.Next() {
		user, err := m.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user result: %w", err)
		}

		list = append(list, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		m.logger.Error(err, "closing rows")
	}

	return list, nil
}

// buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID
func (m *MariaDB) buildGetUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):         userID,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetUser fetches a user.
func (m *MariaDB) GetUser(ctx context.Context, userID uint64) (*v11.User, error) {
	query, args := m.buildGetUserQuery(userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	u, err := m.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}

// buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user
// by their database ID, who has an unverified two factor secret
func (m *MariaDB) buildGetUserWithUnverifiedTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):                            userID,
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):                    nil,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret
func (m *MariaDB) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v11.User, error) {
	query, args := m.buildGetUserWithUnverifiedTwoFactorSecretQuery(userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	u, err := m.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}

// buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username
func (m *MariaDB) buildGetUserByUsernameQuery(username string) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn): username,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):         nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetUserByUsername fetches a user by their username.
func (m *MariaDB) GetUserByUsername(ctx context.Context, username string) (*v11.User, error) {
	query, args := m.buildGetUserByUsernameQuery(username)
	row := m.db.QueryRowContext(ctx, query, args...)

	u, err := m.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return u, nil
}

// buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere
// to a given filter's criteria.
func (m *MariaDB) buildGetAllUsersCountQuery() (query string) {
	var err error

	builder := m.sqlBuilder.
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		})

	query, _, err = builder.ToSql()

	m.logQueryBuildingError(err)

	return query
}

// GetAllUsersCount fetches a count of users from the database.
func (m *MariaDB) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	query := m.buildGetAllUsersCountQuery()
	err = m.db.QueryRowContext(ctx, query).Scan(&count)
	return
}

// buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere
// to a given filter's criteria.
func (m *MariaDB) buildGetUsersQuery(filter *v11.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", usersTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, usersTableName)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)
	return query, args
}

// GetUsers fetches a list of users from the database that meet a particular filter.
func (m *MariaDB) GetUsers(ctx context.Context, filter *v11.QueryFilter) (*v11.UserList, error) {
	query, args := m.buildGetUsersQuery(filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying for user")
	}

	userList, err := m.scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("loading response from database: %w", err)
	}

	x := &v11.UserList{
		Pagination: v11.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Users: userList,
	}

	return x, nil
}

// buildCreateUserQuery returns a SQL query (and arguments) that would create a given User
func (m *MariaDB) buildCreateUserQuery(input v11.UserDatabaseCreationInput) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(usersTableName).
		Columns(
			usersTableUsernameColumn,
			usersTableHashedPasswordColumn,
			usersTableSaltColumn,
			usersTableTwoFactorColumn,
			usersTableIsAdminColumn,
		).
		Values(
			input.Username,
			input.HashedPassword,
			input.Salt,
			input.TwoFactorSecret,
			false,
		).
		ToSql()

	// NOTE: we always default is_admin to false, on the assumption that
	// admins have DB access and will change that value via SQL query.
	// There should also be no way to update a user via this structure
	// such that they would have admin privileges.

	m.logQueryBuildingError(err)

	return query, args
}

// CreateUser creates a user.
func (m *MariaDB) CreateUser(ctx context.Context, input v11.UserDatabaseCreationInput) (*v11.User, error) {
	x := &v11.User{
		Username:        input.Username,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
	}
	query, args := m.buildCreateUserQuery(input)

	// create the user.
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing user creation query: %w", err)
	}

	// fetch the last inserted ID.
	id, err := res.LastInsertId()
	m.logIDRetrievalError(err)
	x.ID = uint64(id)

	// this won't be completely accurate, but it will suffice.
	x.CreatedOn = m.timeTeller.Now()

	return x, nil
}

// buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row
func (m *MariaDB) buildUpdateUserQuery(input *v11.User) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(usersTableUsernameColumn, input.Username).
		Set(usersTableHashedPasswordColumn, input.HashedPassword).
		Set(usersTableSaltColumn, input.Salt).
		Set(usersTableTwoFactorColumn, input.TwoFactorSecret).
		Set(usersTableTwoFactorVerifiedOnColumn, input.TwoFactorSecretVerifiedOn).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateUser receives a complete User struct and updates its place in the db.
// NOTE this function uses the ID provided in the input to make its query. Pass in
// incomplete models at your peril.
func (m *MariaDB) UpdateUser(ctx context.Context, input *v11.User) error {
	query, args := m.buildUpdateUserQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password.
func (m *MariaDB) buildUpdateUserPasswordQuery(userID uint64, newHash string) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(usersTableHashedPasswordColumn, newHash).
		Set(usersTableRequiresPasswordChangeColumn, false).
		Set(usersTablePasswordLastChangedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateUserPassword updates a user's password.
func (m *MariaDB) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	query, args := m.buildUpdateUserPasswordQuery(userID, newHash)

	_, err := m.db.ExecContext(ctx, query, args...)

	return err
}

// buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret
func (m *MariaDB) buildVerifyUserTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(usersTableTwoFactorVerifiedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (m *MariaDB) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	query, args := m.buildVerifyUserTwoFactorSecretQuery(userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveUserQuery builds a SQL query that marks a user as archived.
func (m *MariaDB) buildArchiveUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveUser marks a user as archived.
func (m *MariaDB) ArchiveUser(ctx context.Context, userID uint64) error {
	query, args := m.buildArchiveUserQuery(userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersFileConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUsersFileConstDeclarations()

		expected := `
package example

import ()

const (
	usersTableName                         = "users"
	usersTableUsernameColumn               = "username"
	usersTableHashedPasswordColumn         = "hashed_password"
	usersTableSaltColumn                   = "salt"
	usersTableRequiresPasswordChangeColumn = "requires_password_change"
	usersTablePasswordLastChangedOnColumn  = "password_last_changed_on"
	usersTableTwoFactorColumn              = "two_factor_secret"
	usersTableTwoFactorVerifiedOnColumn    = "two_factor_secret_verified_on"
	usersTableIsAdminColumn                = "is_admin"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUsersFileVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUsersFileVarDeclarations()

		expected := `
package example

import (
	"fmt"
)

var (
	usersTableColumns = []string{
		fmt.Sprintf("%s.%s", usersTableName, idColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableHashedPasswordColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableSaltColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableRequiresPasswordChangeColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTablePasswordLastChangedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, usersTableIsAdminColumn),
		fmt.Sprintf("%s.%s", usersTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn),
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildScanUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildScanUser(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanUser provides a consistent way to scan something like a *sql.Row into a User struct.
func (p *Postgres) scanUser(scan v1.Scanner) (*v11.User, error) {
	var (
		x = &v11.User{}
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Username,
		&x.HashedPassword,
		&x.Salt,
		&x.RequiresPasswordChange,
		&x.PasswordLastChangedOn,
		&x.TwoFactorSecret,
		&x.TwoFactorSecretVerifiedOn,
		&x.IsAdmin,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildScanUser(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanUser provides a consistent way to scan something like a *sql.Row into a User struct.
func (s *Sqlite) scanUser(scan v1.Scanner) (*v11.User, error) {
	var (
		x = &v11.User{}
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Username,
		&x.HashedPassword,
		&x.Salt,
		&x.RequiresPasswordChange,
		&x.PasswordLastChangedOn,
		&x.TwoFactorSecret,
		&x.TwoFactorSecretVerifiedOn,
		&x.IsAdmin,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
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
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildScanUser(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanUser provides a consistent way to scan something like a *sql.Row into a User struct.
func (m *MariaDB) scanUser(scan v1.Scanner) (*v11.User, error) {
	var (
		x = &v11.User{}
	)

	targetVars := []interface{}{
		&x.ID,
		&x.Username,
		&x.HashedPassword,
		&x.Salt,
		&x.RequiresPasswordChange,
		&x.PasswordLastChangedOn,
		&x.TwoFactorSecret,
		&x.TwoFactorSecretVerifiedOn,
		&x.IsAdmin,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
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
}

func Test_buildScanUsers(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildScanUsers(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanUsers takes database rows and loads them into a slice of User structs.
func (p *Postgres) scanUsers(rows v1.ResultIterator) ([]v11.User, error) {
	var (
		list []v11.User
	)

	for rows.Next() {
		user, err := p.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user result: %w", err)
		}

		list = append(list, *user)
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
		x := buildScanUsers(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanUsers takes database rows and loads them into a slice of User structs.
func (s *Sqlite) scanUsers(rows v1.ResultIterator) ([]v11.User, error) {
	var (
		list []v11.User
	)

	for rows.Next() {
		user, err := s.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user result: %w", err)
		}

		list = append(list, *user)
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
		x := buildScanUsers(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// scanUsers takes database rows and loads them into a slice of User structs.
func (m *MariaDB) scanUsers(rows v1.ResultIterator) ([]v11.User, error) {
	var (
		list []v11.User
	)

	for rows.Next() {
		user, err := m.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user result: %w", err)
		}

		list = append(list, *user)
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

func Test_buildBuildGetUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetUserQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID
func (p *Postgres) buildGetUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):         userID,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetUserQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID
func (s *Sqlite) buildGetUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):         userID,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
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

		x := buildBuildGetUserQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserQuery returns a SQL query (and argument) for retrieving a user by their database ID
func (m *MariaDB) buildGetUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):         userID,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
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

func Test_buildGetUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUser fetches a user.
func (p *Postgres) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	query, args := p.buildGetUserQuery(userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	u, err := p.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUser fetches a user.
func (s *Sqlite) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	query, args := s.buildGetUserQuery(userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	u, err := s.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUser fetches a user.
func (m *MariaDB) GetUser(ctx context.Context, userID uint64) (*v1.User, error) {
	query, args := m.buildGetUserQuery(userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	u, err := m.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user
// by their database ID, who has an unverified two factor secret
func (p *Postgres) buildGetUserWithUnverifiedTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):                            userID,
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):                    nil,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user
// by their database ID, who has an unverified two factor secret
func (s *Sqlite) buildGetUserWithUnverifiedTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):                            userID,
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):                    nil,
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

		x := buildBuildGetUserWithUnverifiedTwoFactorSecretQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserWithUnverifiedTwoFactorSecretQuery returns a SQL query (and argument) for retrieving a user
// by their database ID, who has an unverified two factor secret
func (m *MariaDB) buildGetUserWithUnverifiedTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, idColumn):                            userID,
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):                    nil,
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

func Test_buildGetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetUserWithUnverifiedTwoFactorSecret(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret
func (p *Postgres) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	query, args := p.buildGetUserWithUnverifiedTwoFactorSecretQuery(userID)
	row := p.db.QueryRowContext(ctx, query, args...)

	u, err := p.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetUserWithUnverifiedTwoFactorSecret(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret
func (s *Sqlite) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	query, args := s.buildGetUserWithUnverifiedTwoFactorSecretQuery(userID)
	row := s.db.QueryRowContext(ctx, query, args...)

	u, err := s.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetUserWithUnverifiedTwoFactorSecret(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserWithUnverifiedTwoFactorSecret fetches a user with an unverified two factor secret
func (m *MariaDB) GetUserWithUnverifiedTwoFactorSecret(ctx context.Context, userID uint64) (*v1.User, error) {
	query, args := m.buildGetUserWithUnverifiedTwoFactorSecretQuery(userID)
	row := m.db.QueryRowContext(ctx, query, args...)

	u, err := m.scanUser(row)
	if err != nil {
		return nil, buildError(err, "fetching user from database")
	}

	return u, err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUserByUsernameQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetUserByUsernameQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username
func (p *Postgres) buildGetUserByUsernameQuery(username string) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn): username,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):         nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetUserByUsernameQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username
func (s *Sqlite) buildGetUserByUsernameQuery(username string) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn): username,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):         nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
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

		x := buildBuildGetUserByUsernameQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetUserByUsernameQuery returns a SQL query (and argument) for retrieving a user by their username
func (m *MariaDB) buildGetUserByUsernameQuery(username string) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableUsernameColumn): username,
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn):         nil,
		}).
		Where(squirrel.NotEq{
			fmt.Sprintf("%s.%s", usersTableName, usersTableTwoFactorVerifiedOnColumn): nil,
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

func Test_buildGetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetUserByUsername(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserByUsername fetches a user by their username.
func (p *Postgres) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	query, args := p.buildGetUserByUsernameQuery(username)
	row := p.db.QueryRowContext(ctx, query, args...)

	u, err := p.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return u, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildGetUserByUsername(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserByUsername fetches a user by their username.
func (s *Sqlite) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	query, args := s.buildGetUserByUsernameQuery(username)
	row := s.db.QueryRowContext(ctx, query, args...)

	u, err := s.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return u, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildGetUserByUsername(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUserByUsername fetches a user by their username.
func (m *MariaDB) GetUserByUsername(ctx context.Context, username string) (*v1.User, error) {
	query, args := m.buildGetUserByUsernameQuery(username)
	row := m.db.QueryRowContext(ctx, query, args...)

	u, err := m.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("fetching user from database: %w", err)
	}

	return u, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetAllUsersCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildGetAllUsersCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere
// to a given filter's criteria.
func (p *Postgres) buildGetAllUsersCountQuery() (query string) {
	var err error

	builder := p.sqlBuilder.
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		})

	query, _, err = builder.ToSql()

	p.logQueryBuildingError(err)

	return query
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildGetAllUsersCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere
// to a given filter's criteria.
func (s *Sqlite) buildGetAllUsersCountQuery() (query string) {
	var err error

	builder := s.sqlBuilder.
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		})

	query, _, err = builder.ToSql()

	s.logQueryBuildingError(err)

	return query
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildBuildGetAllUsersCountQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildGetAllUsersCountQuery returns a SQL query (and arguments) for retrieving the number of users who adhere
// to a given filter's criteria.
func (m *MariaDB) buildGetAllUsersCountQuery() (query string) {
	var err error

	builder := m.sqlBuilder.
		Select(fmt.Sprintf(countQuery, usersTableName)).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		})

	query, _, err = builder.ToSql()

	m.logQueryBuildingError(err)

	return query
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildGetAllUsersCount(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildGetAllUsersCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllUsersCount fetches a count of users from the database.
func (p *Postgres) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	query := p.buildGetAllUsersCountQuery()
	err = p.db.QueryRowContext(ctx, query).Scan(&count)
	return
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildGetAllUsersCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllUsersCount fetches a count of users from the database.
func (s *Sqlite) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	query := s.buildGetAllUsersCountQuery()
	err = s.db.QueryRowContext(ctx, query).Scan(&count)
	return
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		t.Parallel()

		dbvendor := buildMariaDBWord()

		x := buildGetAllUsersCount(dbvendor)

		expected := `
package example

import (
	"context"
)

// GetAllUsersCount fetches a count of users from the database.
func (m *MariaDB) GetAllUsersCount(ctx context.Context) (count uint64, err error) {
	query := m.buildGetAllUsersCountQuery()
	err = m.db.QueryRowContext(ctx, query).Scan(&count)
	return
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUsersQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildGetUsersQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere
// to a given filter's criteria.
func (p *Postgres) buildGetUsersQuery(filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", usersTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, usersTableName)
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
		x := buildBuildGetUsersQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere
// to a given filter's criteria.
func (s *Sqlite) buildGetUsersQuery(filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := s.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", usersTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, usersTableName)
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
		x := buildBuildGetUsersQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildGetUsersQuery returns a SQL query (and arguments) for retrieving a slice of users who adhere
// to a given filter's criteria.
func (m *MariaDB) buildGetUsersQuery(filter *v1.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := m.sqlBuilder.
		Select(usersTableColumns...).
		From(usersTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", usersTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", usersTableName, idColumn))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, usersTableName)
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

func Test_buildGetUsers(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildGetUsers(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUsers fetches a list of users from the database that meet a particular filter.
func (p *Postgres) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	query, args := p.buildGetUsersQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying for user")
	}

	userList, err := p.scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("loading response from database: %w", err)
	}

	x := &v1.UserList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Users: userList,
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
		x := buildGetUsers(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUsers fetches a list of users from the database that meet a particular filter.
func (s *Sqlite) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	query, args := s.buildGetUsersQuery(filter)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying for user")
	}

	userList, err := s.scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("loading response from database: %w", err)
	}

	x := &v1.UserList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Users: userList,
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
		x := buildGetUsers(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// GetUsers fetches a list of users from the database that meet a particular filter.
func (m *MariaDB) GetUsers(ctx context.Context, filter *v1.QueryFilter) (*v1.UserList, error) {
	query, args := m.buildGetUsersQuery(filter)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying for user")
	}

	userList, err := m.scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("loading response from database: %w", err)
	}

	x := &v1.UserList{
		Pagination: v1.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Users: userList,
	}

	return x, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildCreateUserQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildCreateUserQuery returns a SQL query (and arguments) that would create a given User
func (p *Postgres) buildCreateUserQuery(input v1.UserDatabaseCreationInput) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(usersTableName).
		Columns(
			usersTableUsernameColumn,
			usersTableHashedPasswordColumn,
			usersTableSaltColumn,
			usersTableTwoFactorColumn,
			usersTableIsAdminColumn,
		).
		Values(
			input.Username,
			input.HashedPassword,
			input.Salt,
			input.TwoFactorSecret,
			false,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
		ToSql()

	// NOTE: we always default is_admin to false, on the assumption that
	// admins have DB access and will change that value via SQL query.
	// There should also be no way to update a user via this structure
	// such that they would have admin privileges.

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
		x := buildBuildCreateUserQuery(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildCreateUserQuery returns a SQL query (and arguments) that would create a given User
func (s *Sqlite) buildCreateUserQuery(input v1.UserDatabaseCreationInput) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Insert(usersTableName).
		Columns(
			usersTableUsernameColumn,
			usersTableHashedPasswordColumn,
			usersTableSaltColumn,
			usersTableTwoFactorColumn,
			usersTableIsAdminColumn,
		).
		Values(
			input.Username,
			input.HashedPassword,
			input.Salt,
			input.TwoFactorSecret,
			false,
		).
		ToSql()

	// NOTE: we always default is_admin to false, on the assumption that
	// admins have DB access and will change that value via SQL query.
	// There should also be no way to update a user via this structure
	// such that they would have admin privileges.

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
		x := buildBuildCreateUserQuery(proj, dbvendor)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildCreateUserQuery returns a SQL query (and arguments) that would create a given User
func (m *MariaDB) buildCreateUserQuery(input v1.UserDatabaseCreationInput) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Insert(usersTableName).
		Columns(
			usersTableUsernameColumn,
			usersTableHashedPasswordColumn,
			usersTableSaltColumn,
			usersTableTwoFactorColumn,
			usersTableIsAdminColumn,
		).
		Values(
			input.Username,
			input.HashedPassword,
			input.Salt,
			input.TwoFactorSecret,
			false,
		).
		ToSql()

	// NOTE: we always default is_admin to false, on the assumption that
	// admins have DB access and will change that value via SQL query.
	// There should also be no way to update a user via this structure
	// such that they would have admin privileges.

	m.logQueryBuildingError(err)

	return query, args
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCreateUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildCreateUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	pq "github.com/lib/pq"
	client "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1/client"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateUser creates a user.
func (p *Postgres) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	x := &v1.User{
		Username:        input.Username,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
	}
	query, args := p.buildCreateUserQuery(input)

	// create the user.
	if err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn); err != nil {
		switch e := err.(type) {
		case *pq.Error:
			if e.Code == pq.ErrorCode(postgresRowExistsErrorCode) {
				return nil, client.ErrUserExists
			}
		default:
			return nil, fmt.Errorf("error executing user creation query: %w", err)
		}
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
		x := buildCreateUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateUser creates a user.
func (s *Sqlite) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	x := &v1.User{
		Username:        input.Username,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
	}
	query, args := s.buildCreateUserQuery(input)

	// create the user.
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing user creation query: %w", err)
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
		x := buildCreateUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// CreateUser creates a user.
func (m *MariaDB) CreateUser(ctx context.Context, input v1.UserDatabaseCreationInput) (*v1.User, error) {
	x := &v1.User{
		Username:        input.Username,
		HashedPassword:  input.HashedPassword,
		TwoFactorSecret: input.TwoFactorSecret,
	}
	query, args := m.buildCreateUserQuery(input)

	// create the user.
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing user creation query: %w", err)
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

func Test_buildBuildUpdateUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildBuildUpdateUserQuery(proj, dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row
func (p *Postgres) buildUpdateUserQuery(input *v1.User) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(usersTableUsernameColumn, input.Username).
		Set(usersTableHashedPasswordColumn, input.HashedPassword).
		Set(usersTableSaltColumn, input.Salt).
		Set(usersTableTwoFactorColumn, input.TwoFactorSecret).
		Set(usersTableTwoFactorVerifiedOnColumn, input.TwoFactorSecretVerifiedOn).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
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
		x := buildBuildUpdateUserQuery(proj, dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row
func (s *Sqlite) buildUpdateUserQuery(input *v1.User) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(usersTableUsernameColumn, input.Username).
		Set(usersTableHashedPasswordColumn, input.HashedPassword).
		Set(usersTableSaltColumn, input.Salt).
		Set(usersTableTwoFactorColumn, input.TwoFactorSecret).
		Set(usersTableTwoFactorVerifiedOnColumn, input.TwoFactorSecretVerifiedOn).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
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
		x := buildBuildUpdateUserQuery(proj, dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// buildUpdateUserQuery returns a SQL query (and arguments) that would update the given user's row
func (m *MariaDB) buildUpdateUserQuery(input *v1.User) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(usersTableUsernameColumn, input.Username).
		Set(usersTableHashedPasswordColumn, input.HashedPassword).
		Set(usersTableSaltColumn, input.Salt).
		Set(usersTableTwoFactorColumn, input.TwoFactorSecret).
		Set(usersTableTwoFactorVerifiedOnColumn, input.TwoFactorSecretVerifiedOn).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
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

func Test_buildUpdateUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildUpdateUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateUser receives a complete User struct and updates its place in the db.
// NOTE this function uses the ID provided in the input to make its query. Pass in
// incomplete models at your peril.
func (p *Postgres) UpdateUser(ctx context.Context, input *v1.User) error {
	query, args := p.buildUpdateUserQuery(input)
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
		x := buildUpdateUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateUser receives a complete User struct and updates its place in the db.
// NOTE this function uses the ID provided in the input to make its query. Pass in
// incomplete models at your peril.
func (s *Sqlite) UpdateUser(ctx context.Context, input *v1.User) error {
	query, args := s.buildUpdateUserQuery(input)
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
		x := buildUpdateUser(proj, dbvendor)

		expected := `
package example

import (
	"context"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// UpdateUser receives a complete User struct and updates its place in the db.
// NOTE this function uses the ID provided in the input to make its query. Pass in
// incomplete models at your peril.
func (m *MariaDB) UpdateUser(ctx context.Context, input *v1.User) error {
	query, args := m.buildUpdateUserQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildUpdateUserPasswordQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildUpdateUserPasswordQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password.
func (p *Postgres) buildUpdateUserPasswordQuery(userID uint64, newHash string) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(usersTableHashedPasswordColumn, newHash).
		Set(usersTableRequiresPasswordChangeColumn, false).
		Set(usersTablePasswordLastChangedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

		x := buildBuildUpdateUserPasswordQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password.
func (s *Sqlite) buildUpdateUserPasswordQuery(userID uint64, newHash string) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(usersTableHashedPasswordColumn, newHash).
		Set(usersTableRequiresPasswordChangeColumn, false).
		Set(usersTablePasswordLastChangedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

		x := buildBuildUpdateUserPasswordQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildUpdateUserPasswordQuery returns a SQL query (and arguments) that would update the given user's password.
func (m *MariaDB) buildUpdateUserPasswordQuery(userID uint64, newHash string) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(usersTableHashedPasswordColumn, newHash).
		Set(usersTableRequiresPasswordChangeColumn, false).
		Set(usersTablePasswordLastChangedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

func Test_buildUpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildUpdateUserPassword(dbvendor)

		expected := `
package example

import (
	"context"
)

// UpdateUserPassword updates a user's password.
func (p *Postgres) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	query, args := p.buildUpdateUserPasswordQuery(userID, newHash)

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

		x := buildUpdateUserPassword(dbvendor)

		expected := `
package example

import (
	"context"
)

// UpdateUserPassword updates a user's password.
func (s *Sqlite) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	query, args := s.buildUpdateUserPasswordQuery(userID, newHash)

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

		x := buildUpdateUserPassword(dbvendor)

		expected := `
package example

import (
	"context"
)

// UpdateUserPassword updates a user's password.
func (m *MariaDB) UpdateUserPassword(ctx context.Context, userID uint64, newHash string) error {
	query, args := m.buildUpdateUserPasswordQuery(userID, newHash)

	_, err := m.db.ExecContext(ctx, query, args...)

	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildVerifyUserTwoFactorSecretQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildVerifyUserTwoFactorSecretQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret
func (p *Postgres) buildVerifyUserTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(usersTableTwoFactorVerifiedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildBuildVerifyUserTwoFactorSecretQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret
func (s *Sqlite) buildVerifyUserTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(usersTableTwoFactorVerifiedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

		x := buildBuildVerifyUserTwoFactorSecretQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildVerifyUserTwoFactorSecretQuery returns a SQL query (and arguments) that would update a given user's two factor secret
func (m *MariaDB) buildVerifyUserTwoFactorSecretQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(usersTableTwoFactorVerifiedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

func Test_buildVerifyUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildVerifyUserTwoFactorSecret(dbvendor)

		expected := `
package example

import (
	"context"
)

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (p *Postgres) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	query, args := p.buildVerifyUserTwoFactorSecretQuery(userID)
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

		x := buildVerifyUserTwoFactorSecret(dbvendor)

		expected := `
package example

import (
	"context"
)

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (s *Sqlite) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	query, args := s.buildVerifyUserTwoFactorSecretQuery(userID)
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

		x := buildVerifyUserTwoFactorSecret(dbvendor)

		expected := `
package example

import (
	"context"
)

// VerifyUserTwoFactorSecret marks a user's two factor secret as validated.
func (m *MariaDB) VerifyUserTwoFactorSecret(ctx context.Context, userID uint64) error {
	query, args := m.buildVerifyUserTwoFactorSecretQuery(userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveUserQuery(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildBuildArchiveUserQuery(dbvendor)

		expected := `
package example

import (
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveUserQuery builds a SQL query that marks a user as archived.
func (p *Postgres) buildArchiveUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(usersTableName).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

		x := buildBuildArchiveUserQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveUserQuery builds a SQL query that marks a user as archived.
func (s *Sqlite) buildArchiveUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = s.sqlBuilder.
		Update(usersTableName).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

		x := buildBuildArchiveUserQuery(dbvendor)

		expected := `
package example

import (
	squirrel "github.com/Masterminds/squirrel"
)

// buildArchiveUserQuery builds a SQL query that marks a user as archived.
func (m *MariaDB) buildArchiveUserQuery(userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Update(usersTableName).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: userID,
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

func Test_buildArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		t.Parallel()

		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildArchiveUser(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveUser marks a user as archived.
func (p *Postgres) ArchiveUser(ctx context.Context, userID uint64) error {
	query, args := p.buildArchiveUserQuery(userID)
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

		x := buildArchiveUser(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveUser marks a user as archived.
func (s *Sqlite) ArchiveUser(ctx context.Context, userID uint64) error {
	query, args := s.buildArchiveUserQuery(userID)
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

		x := buildArchiveUser(dbvendor)

		expected := `
package example

import (
	"context"
)

// ArchiveUser marks a user as archived.
func (m *MariaDB) ArchiveUser(ctx context.Context, userID uint64) error {
	query, args := m.buildArchiveUserQuery(userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
