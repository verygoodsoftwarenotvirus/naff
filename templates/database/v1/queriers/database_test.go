package queriers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_databaseDotGo(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := databaseDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	pq "github.com/lib/pq"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	loggerName         = "postgres"
	postgresDriverName = "wrapped-postgres-driver"

	postgresRowExistsErrorCode = "23505"

	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	idColumn            = "id"
	createdOnColumn     = "created_on"
	lastUpdatedOnColumn = "last_updated_on"
	archivedOnColumn    = "archived_on"

	// countQuery is a generic counter query used in a few query builders.
	countQuery = "COUNT(%s.id)"

	// currentUnixTimeQuery is the query postgres uses to determine the current unix time.
	currentUnixTimeQuery = "extract(epoch FROM NOW())"

	defaultBucketSize = uint64(1000)
)

func init() {
	// Explicitly wrap the Postgres driver with ocsql.
	driver := ocsql.Wrap(
		&pq.Driver{},
		ocsql.WithQuery(true),
		ocsql.WithAllowRoot(false),
		ocsql.WithRowsNext(true),
		ocsql.WithRowsClose(true),
		ocsql.WithQueryParams(true),
	)

	// Register our ocsql wrapper as a db driver.
	sql.Register(postgresDriverName, driver)
}

var _ v1.DataManager = (*Postgres)(nil)

type (
	// Postgres is our main Postgres interaction db.
	Postgres struct {
		logger      v11.Logger
		db          *sql.DB
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}

	// ConnectionDetails is a string alias for a Postgres url.
	ConnectionDetails string

	// Querier is a subset interface for sql.{DB|Tx|Stmt} objects
	Querier interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)

// ProvidePostgresDB provides an instrumented postgres db.
func ProvidePostgresDB(logger v11.Logger, connectionDetails v1.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to postgres")
	return sql.Open(postgresDriverName, string(connectionDetails))
}

// ProvidePostgres provides a postgres db controller.
func ProvidePostgres(debug bool, db *sql.DB, logger v11.Logger) v1.DataManager {
	return &Postgres{
		db:         db,
		debug:      debug,
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// IsReady reports whether or not the db is ready.
func (p *Postgres) IsReady(ctx context.Context) (ready bool) {
	numberOfUnsuccessfulAttempts := 0

	p.logger.WithValues(map[string]interface{}{
		"interval":     time.Second,
		"max_attempts": 50,
	}).Debug("IsReady called")

	for !ready {
		err := p.db.PingContext(ctx)
		if err != nil {
			p.logger.Debug("ping failed, waiting for db")
			time.Sleep(time.Second)

			numberOfUnsuccessfulAttempts++
			if numberOfUnsuccessfulAttempts >= 50 {
				return false
			}
		} else {
			ready = true
			return ready
		}
	}
	return false
}

// logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (p *Postgres) logQueryBuildingError(err error) {
	if err != nil {
		p.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}

// buildError takes a given error and wraps it with a message, provided that it
// IS NOT sql.ErrNoRows, which we want to preserve and surface to the services.
func buildError(err error, msg string) error {
	if err == sql.ErrNoRows {
		return err
	}

	if !strings.Contains(msg, ` + "`" + `%w` + "`" + `) {
		msg += ": %w"
	}

	return fmt.Errorf(msg, err)
}

func joinUint64s(in []uint64) string {
	out := []string{}

	for _, x := range in {
		out = append(out, strconv.FormatUint(x, 10))
	}

	return strings.Join(out, ",")
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := databaseDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	gosqlite3 "github.com/mattn/go-sqlite3"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"strings"
	"sync"
)

const (
	loggerName       = "sqlite"
	sqliteDriverName = "wrapped-sqlite-driver"

	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	idColumn            = "id"
	createdOnColumn     = "created_on"
	lastUpdatedOnColumn = "last_updated_on"
	archivedOnColumn    = "archived_on"

	// countQuery is a generic counter query used in a few query builders.
	countQuery = "COUNT(%s.id)"

	// currentUnixTimeQuery is the query sqlite uses to determine the current unix time.
	currentUnixTimeQuery = "(strftime('%s','now'))"

	defaultBucketSize = uint64(1000)
)

func init() {
	// Explicitly wrap the Sqlite driver with ocsql.
	driver := ocsql.Wrap(
		&gosqlite3.SQLiteDriver{},
		ocsql.WithQuery(true),
		ocsql.WithAllowRoot(false),
		ocsql.WithRowsNext(true),
		ocsql.WithRowsClose(true),
		ocsql.WithQueryParams(true),
	)

	// Register our ocsql wrapper as a db driver.
	sql.Register(sqliteDriverName, driver)
}

var _ v1.DataManager = (*Sqlite)(nil)

type (
	// Sqlite is our main Sqlite interaction db.
	Sqlite struct {
		logger      v11.Logger
		db          *sql.DB
		timeTeller  timeTeller
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}

	// ConnectionDetails is a string alias for a Sqlite url.
	ConnectionDetails string

	// Querier is a subset interface for sql.{DB|Tx|Stmt} objects
	Querier interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)

// ProvideSqliteDB provides an instrumented sqlite db.
func ProvideSqliteDB(logger v11.Logger, connectionDetails v1.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to sqlite")
	return sql.Open(sqliteDriverName, string(connectionDetails))
}

// ProvideSqlite provides a sqlite db controller.
func ProvideSqlite(debug bool, db *sql.DB, logger v11.Logger) v1.DataManager {
	return &Sqlite{
		db:         db,
		debug:      debug,
		timeTeller: &stdLibTimeTeller{},
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

// IsReady reports whether or not the db is ready.
func (s *Sqlite) IsReady(_ context.Context) (ready bool) {
	return true
}

// logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (s *Sqlite) logQueryBuildingError(err error) {
	if err != nil {
		s.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}

// logIDRetrievalError logs errors that may occur during created db row ID retrieval.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (s *Sqlite) logIDRetrievalError(err error) {
	if err != nil {
		s.logger.WithName("ROW_ID_ERROR").Error(err, "fetching row ID")
	}
}

// buildError takes a given error and wraps it with a message, provided that it
// IS NOT sql.ErrNoRows, which we want to preserve and surface to the services.
func buildError(err error, msg string) error {
	if err == sql.ErrNoRows {
		return err
	}

	if !strings.Contains(msg, ` + "`" + `%w` + "`" + `) {
		msg += ": %w"
	}

	return fmt.Errorf(msg, err)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := databaseDotGo(proj, dbvendor)

		expected := `
package example

import (
	"context"
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	"fmt"
	squirrel "github.com/Masterminds/squirrel"
	mysql "github.com/go-sql-driver/mysql"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"strings"
	"sync"
	"time"
)

const (
	loggerName        = "mariadb"
	mariaDBDriverName = "wrapped-mariadb-driver"

	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	idColumn            = "id"
	createdOnColumn     = "created_on"
	lastUpdatedOnColumn = "last_updated_on"
	archivedOnColumn    = "archived_on"

	// countQuery is a generic counter query used in a few query builders.
	countQuery = "COUNT(%s.id)"

	// currentUnixTimeQuery is the query maria DB uses to determine the current unix time.
	currentUnixTimeQuery = "UNIX_TIMESTAMP()"

	defaultBucketSize = uint64(1000)
)

func init() {
	// Explicitly wrap the MariaDB driver with ocsql.
	driver := ocsql.Wrap(
		&mysql.MySQLDriver{},
		ocsql.WithQuery(true),
		ocsql.WithAllowRoot(false),
		ocsql.WithRowsNext(true),
		ocsql.WithRowsClose(true),
		ocsql.WithQueryParams(true),
	)

	// Register our ocsql wrapper as a db driver.
	sql.Register(mariaDBDriverName, driver)
}

var _ v1.DataManager = (*MariaDB)(nil)

type (
	// MariaDB is our main MariaDB interaction db.
	MariaDB struct {
		logger      v11.Logger
		db          *sql.DB
		timeTeller  timeTeller
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}

	// ConnectionDetails is a string alias for a MariaDB url.
	ConnectionDetails string

	// Querier is a subset interface for sql.{DB|Tx|Stmt} objects
	Querier interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)

// ProvideMariaDBConnection provides an instrumented maria DB db.
func ProvideMariaDBConnection(logger v11.Logger, connectionDetails v1.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to maria DB")
	return sql.Open(mariaDBDriverName, string(connectionDetails))
}

// ProvideMariaDB provides a maria DB controller.
func ProvideMariaDB(debug bool, db *sql.DB, logger v11.Logger) v1.DataManager {
	return &MariaDB{
		db:         db,
		debug:      debug,
		timeTeller: &stdLibTimeTeller{},
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

// IsReady reports whether or not the db is ready.
func (m *MariaDB) IsReady(ctx context.Context) (ready bool) {
	numberOfUnsuccessfulAttempts := 0

	m.logger.WithValues(map[string]interface{}{
		"interval":     time.Second,
		"max_attempts": 50,
	}).Debug("IsReady called")

	for !ready {
		err := m.db.PingContext(ctx)
		if err != nil {
			m.logger.Debug("ping failed, waiting for db")
			time.Sleep(time.Second)

			numberOfUnsuccessfulAttempts++
			if numberOfUnsuccessfulAttempts >= 50 {
				return false
			}
		} else {
			ready = true
			return ready
		}
	}
	return false
}

// logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (m *MariaDB) logQueryBuildingError(err error) {
	if err != nil {
		m.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}

// logIDRetrievalError logs errors that may occur during created db row ID retrieval.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (m *MariaDB) logIDRetrievalError(err error) {
	if err != nil {
		m.logger.WithName("ROW_ID_ERROR").Error(err, "fetching row ID")
	}
}

// buildError takes a given error and wraps it with a message, provided that it
// IS NOT sql.ErrNoRows, which we want to preserve and surface to the services.
func buildError(err error, msg string) error {
	if err == sql.ErrNoRows {
		return err
	}

	if !strings.Contains(msg, ` + "`" + `%w` + "`" + `) {
		msg += ": %w"
	}

	return fmt.Errorf(msg, err)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDBDotGoConsts(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildDBDotGoConsts(dbvendor)

		expected := `
package example

import ()

const (
	loggerName         = "postgres"
	postgresDriverName = "wrapped-postgres-driver"

	postgresRowExistsErrorCode = "23505"

	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	idColumn            = "id"
	createdOnColumn     = "created_on"
	lastUpdatedOnColumn = "last_updated_on"
	archivedOnColumn    = "archived_on"

	// countQuery is a generic counter query used in a few query builders.
	countQuery = "COUNT(%s.id)"

	// currentUnixTimeQuery is the query postgres uses to determine the current unix time.
	currentUnixTimeQuery = "extract(epoch FROM NOW())"

	defaultBucketSize = uint64(1000)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildDBDotGoConsts(dbvendor)

		expected := `
package example

import ()

const (
	loggerName       = "sqlite"
	sqliteDriverName = "wrapped-sqlite-driver"

	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	idColumn            = "id"
	createdOnColumn     = "created_on"
	lastUpdatedOnColumn = "last_updated_on"
	archivedOnColumn    = "archived_on"

	// countQuery is a generic counter query used in a few query builders.
	countQuery = "COUNT(%s.id)"

	// currentUnixTimeQuery is the query sqlite uses to determine the current unix time.
	currentUnixTimeQuery = "(strftime('%s','now'))"

	defaultBucketSize = uint64(1000)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildDBDotGoConsts(dbvendor)

		expected := `
package example

import ()

const (
	loggerName        = "mariadb"
	mariaDBDriverName = "wrapped-mariadb-driver"

	existencePrefix, existenceSuffix = "SELECT EXISTS (", ")"

	idColumn            = "id"
	createdOnColumn     = "created_on"
	lastUpdatedOnColumn = "last_updated_on"
	archivedOnColumn    = "archived_on"

	// countQuery is a generic counter query used in a few query builders.
	countQuery = "COUNT(%s.id)"

	// currentUnixTimeQuery is the query maria DB uses to determine the current unix time.
	currentUnixTimeQuery = "UNIX_TIMESTAMP()"

	defaultBucketSize = uint64(1000)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDBDotGoInit(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildDBDotGoInit(dbvendor)

		expected := `
package example

import (
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	pq "github.com/lib/pq"
)

func init() {
	// Explicitly wrap the Postgres driver with ocsql.
	driver := ocsql.Wrap(
		&pq.Driver{},
		ocsql.WithQuery(true),
		ocsql.WithAllowRoot(false),
		ocsql.WithRowsNext(true),
		ocsql.WithRowsClose(true),
		ocsql.WithQueryParams(true),
	)

	// Register our ocsql wrapper as a db driver.
	sql.Register(postgresDriverName, driver)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildDBDotGoInit(dbvendor)

		expected := `
package example

import (
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	gosqlite3 "github.com/mattn/go-sqlite3"
)

func init() {
	// Explicitly wrap the Sqlite driver with ocsql.
	driver := ocsql.Wrap(
		&gosqlite3.SQLiteDriver{},
		ocsql.WithQuery(true),
		ocsql.WithAllowRoot(false),
		ocsql.WithRowsNext(true),
		ocsql.WithRowsClose(true),
		ocsql.WithQueryParams(true),
	)

	// Register our ocsql wrapper as a db driver.
	sql.Register(sqliteDriverName, driver)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildDBDotGoInit(dbvendor)

		expected := `
package example

import (
	ocsql "contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	mysql "github.com/go-sql-driver/mysql"
)

func init() {
	// Explicitly wrap the MariaDB driver with ocsql.
	driver := ocsql.Wrap(
		&mysql.MySQLDriver{},
		ocsql.WithQuery(true),
		ocsql.WithAllowRoot(false),
		ocsql.WithRowsNext(true),
		ocsql.WithRowsClose(true),
		ocsql.WithQueryParams(true),
	)

	// Register our ocsql wrapper as a db driver.
	sql.Register(mariaDBDriverName, driver)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDBDotGoVarDecls(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildDBDotGoVarDecls(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	squirrel "github.com/Masterminds/squirrel"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"sync"
)

var _ v1.DataManager = (*Postgres)(nil)

type (
	// Postgres is our main Postgres interaction db.
	Postgres struct {
		logger      v11.Logger
		db          *sql.DB
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}

	// ConnectionDetails is a string alias for a Postgres url.
	ConnectionDetails string

	// Querier is a subset interface for sql.{DB|Tx|Stmt} objects
	Querier interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildDBDotGoVarDecls(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	squirrel "github.com/Masterminds/squirrel"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"sync"
)

var _ v1.DataManager = (*Sqlite)(nil)

type (
	// Sqlite is our main Sqlite interaction db.
	Sqlite struct {
		logger      v11.Logger
		db          *sql.DB
		timeTeller  timeTeller
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}

	// ConnectionDetails is a string alias for a Sqlite url.
	ConnectionDetails string

	// Querier is a subset interface for sql.{DB|Tx|Stmt} objects
	Querier interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildDBDotGoVarDecls(proj, dbvendor)

		expected := `
package example

import (
	"context"
	"database/sql"
	squirrel "github.com/Masterminds/squirrel"
	v11 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"sync"
)

var _ v1.DataManager = (*MariaDB)(nil)

type (
	// MariaDB is our main MariaDB interaction db.
	MariaDB struct {
		logger      v11.Logger
		db          *sql.DB
		timeTeller  timeTeller
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}

	// ConnectionDetails is a string alias for a MariaDB url.
	ConnectionDetails string

	// Querier is a subset interface for sql.{DB|Tx|Stmt} objects
	Querier interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideDatabaseConn(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseConn(proj, dbvendor)

		expected := `
package example

import (
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvidePostgresDB provides an instrumented postgres db.
func ProvidePostgresDB(logger v1.Logger, connectionDetails v11.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to postgres")
	return sql.Open(postgresDriverName, string(connectionDetails))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseConn(proj, dbvendor)

		expected := `
package example

import (
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvideSqliteDB provides an instrumented sqlite db.
func ProvideSqliteDB(logger v1.Logger, connectionDetails v11.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to sqlite")
	return sql.Open(sqliteDriverName, string(connectionDetails))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseConn(proj, dbvendor)

		expected := `
package example

import (
	"database/sql"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvideMariaDBConnection provides an instrumented maria DB db.
func ProvideMariaDBConnection(logger v1.Logger, connectionDetails v11.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to maria DB")
	return sql.Open(mariaDBDriverName, string(connectionDetails))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideDatabaseClient(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseClient(proj, dbvendor)

		expected := `
package example

import (
	"database/sql"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvidePostgres provides a postgres db controller.
func ProvidePostgres(debug bool, db *sql.DB, logger v1.Logger) v11.DataManager {
	return &Postgres{
		db:         db,
		debug:      debug,
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseClient(proj, dbvendor)

		expected := `
package example

import (
	"database/sql"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvideSqlite provides a sqlite db controller.
func ProvideSqlite(debug bool, db *sql.DB, logger v1.Logger) v11.DataManager {
	return &Sqlite{
		db:         db,
		debug:      debug,
		timeTeller: &stdLibTimeTeller{},
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()
		proj := testprojects.BuildTodoApp()
		x := buildProvideDatabaseClient(proj, dbvendor)

		expected := `
package example

import (
	"database/sql"
	squirrel "github.com/Masterminds/squirrel"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

// ProvideMariaDB provides a maria DB controller.
func ProvideMariaDB(debug bool, db *sql.DB, logger v1.Logger) v11.DataManager {
	return &MariaDB{
		db:         db,
		debug:      debug,
		timeTeller: &stdLibTimeTeller{},
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIsReady(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildIsReady(dbvendor)

		expected := `
package example

import (
	"context"
	"time"
)

// IsReady reports whether or not the db is ready.
func (p *Postgres) IsReady(ctx context.Context) (ready bool) {
	numberOfUnsuccessfulAttempts := 0

	p.logger.WithValues(map[string]interface{}{
		"interval":     time.Second,
		"max_attempts": 50,
	}).Debug("IsReady called")

	for !ready {
		err := p.db.PingContext(ctx)
		if err != nil {
			p.logger.Debug("ping failed, waiting for db")
			time.Sleep(time.Second)

			numberOfUnsuccessfulAttempts++
			if numberOfUnsuccessfulAttempts >= 50 {
				return false
			}
		} else {
			ready = true
			return ready
		}
	}
	return false
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildIsReady(dbvendor)

		expected := `
package example

import (
	"context"
)

// IsReady reports whether or not the db is ready.
func (s *Sqlite) IsReady(_ context.Context) (ready bool) {
	return true
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildIsReady(dbvendor)

		expected := `
package example

import (
	"context"
	"time"
)

// IsReady reports whether or not the db is ready.
func (m *MariaDB) IsReady(ctx context.Context) (ready bool) {
	numberOfUnsuccessfulAttempts := 0

	m.logger.WithValues(map[string]interface{}{
		"interval":     time.Second,
		"max_attempts": 50,
	}).Debug("IsReady called")

	for !ready {
		err := m.db.PingContext(ctx)
		if err != nil {
			m.logger.Debug("ping failed, waiting for db")
			time.Sleep(time.Second)

			numberOfUnsuccessfulAttempts++
			if numberOfUnsuccessfulAttempts >= 50 {
				return false
			}
		} else {
			ready = true
			return ready
		}
	}
	return false
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("invalid", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("invalid")

		assert.Panics(t, func() { buildIsReady(dbvendor) })
	})
}

func Test_buildLogQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildLogQueryBuildingError(dbvendor)

		expected := `
package example

import ()

// logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (p *Postgres) logQueryBuildingError(err error) {
	if err != nil {
		p.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildLogQueryBuildingError(dbvendor)

		expected := `
package example

import ()

// logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (s *Sqlite) logQueryBuildingError(err error) {
	if err != nil {
		s.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildLogQueryBuildingError(dbvendor)

		expected := `
package example

import ()

// logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (m *MariaDB) logQueryBuildingError(err error) {
	if err != nil {
		m.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildLogIDRetrievalError(T *testing.T) {
	T.Parallel()

	T.Run("postgres", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Postgres")

		x := buildLogIDRetrievalError(dbvendor)

		expected := `
package example

import ()

// logIDRetrievalError logs errors that may occur during created db row ID retrieval.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (p *Postgres) logIDRetrievalError(err error) {
	if err != nil {
		p.logger.WithName("ROW_ID_ERROR").Error(err, "fetching row ID")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("sqlite", func(t *testing.T) {
		dbvendor := wordsmith.FromSingularPascalCase("Sqlite")

		x := buildLogIDRetrievalError(dbvendor)

		expected := `
package example

import ()

// logIDRetrievalError logs errors that may occur during created db row ID retrieval.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (s *Sqlite) logIDRetrievalError(err error) {
	if err != nil {
		s.logger.WithName("ROW_ID_ERROR").Error(err, "fetching row ID")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})

	T.Run("mariadb", func(t *testing.T) {
		dbvendor := buildMariaDBWord()

		x := buildLogIDRetrievalError(dbvendor)

		expected := `
package example

import ()

// logIDRetrievalError logs errors that may occur during created db row ID retrieval.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (m *MariaDB) logIDRetrievalError(err error) {
	if err != nil {
		m.logger.WithName("ROW_ID_ERROR").Error(err, "fetching row ID")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildError()

		expected := `
package example

import (
	"database/sql"
	"fmt"
	"strings"
)

// buildError takes a given error and wraps it with a message, provided that it
// IS NOT sql.ErrNoRows, which we want to preserve and surface to the services.
func buildError(err error, msg string) error {
	if err == sql.ErrNoRows {
		return err
	}

	if !strings.Contains(msg, ` + "`" + `%w` + "`" + `) {
		msg += ": %w"
	}

	return fmt.Errorf(msg, err)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildJoinUint64s(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildJoinUint64s()

		expected := `
package example

import (
	"strconv"
	"strings"
)

func joinUint64s(in []uint64) string {
	out := []string{}

	for _, x := range in {
		out = append(out, strconv.FormatUint(x, 10))
	}

	return strings.Join(out, ",")
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
