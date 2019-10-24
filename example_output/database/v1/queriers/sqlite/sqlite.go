package sqlite

import (
	"context"
	"database/sql"
	"sync"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
)

const (
	loggerName           = "sqlite"
	sqliteDriverName     = "wrapped-sqlite-driver"
	CountQuery           = "COUNT(id)"
	CurrentUnixTimeQuery = "(strftime('%s','now'))"
)

func init() {
	driver := ocsql.Wrap(&sqlite3.SQLiteDriver{}, ocsql.WithQuery(true), ocsql.WithAllowRoot(false), ocsql.WithRowsNext(true), ocsql.WithRowsClose(true), ocsql.WithQueryParams(true))
	sql.Register(sqliteDriverName, driver)
}

var _ database.Database = (*Sqlite)(nil)

type (
	Sqlite struct {
		logger      logging.Logger
		db          *sql.DB
		sqlBuilder  squirrel.StatementBuilderType
		migrateOnce sync.Once
		debug       bool
	}
	ConnectionDetails string
	Querier           interface {
		ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
	}
)

// ProvideSqliteDB provides an instrumented sqlite db
func ProvideSqliteDB(logger logging.Logger, connectionDetails database.ConnectionDetails) (*sql.DB, error) {
	logger.WithValue("connection_details", connectionDetails).Debug("Establishing connection to sqlite")
	return sql.Open(sqliteDriverName, string(connectionDetails))
}

// ProvideSqlite provides a sqlite db controller
func ProvideSqlite(debug bool, db *sql.DB, logger logging.Logger) database.Database {
	return &Sqlite{
		db:         db,
		debug:      debug,
		logger:     logger.WithName(loggerName),
		sqlBuilder: squirrel.StatementBuilder,
	}
}

// IsReady reports whether or not the db is ready
func (s *Sqlite) IsReady(ctx context.Context) (ready bool) {
	return true
}

// s.logQueryBuildingError logs errors that may occur during query construction.
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (s *Sqlite) logQueryBuildingError(err error) {
	if err != nil {
		s.logger.WithName("QUERY_ERROR").Error(err, "building query")
	}
}

// logCreationTimeRetrievalError logs errors that may occur during creation time retrieval
// Such errors should be few and far between, as the generally only occur with
// type discrepancies or other misuses of SQL. An alert should be set up for
// any log entries with the given name, and those alerts should be investigated
// with the utmost priority.
func (s *Sqlite) logCreationTimeRetrievalError(err error) {
	if err != nil {
		s.logger.WithName("CREATION_TIME_RETRIEVAL").Error(err, "building query")
	}
}

// buildError takes a given error and wraps it with a message, provided that it
// IS NOT sql.ErrNoRows, which we want to preserve and surface to the services.
func buildError(err error, msg string) error {
	if err == sql.ErrNoRows {
		return err
	}
	return errors.Wrap(err, msg)
}