package indexinitializer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func TestRenderPackage(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = os.TempDir()

		assert.NoError(t, RenderPackage(proj))
	})

	T.Run("with invalid output directory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.OutputPath = `/\0/\0/\0`

		assert.Error(t, RenderPackage(proj))
	})
}

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"context"
	pflag "github.com/spf13/pflag"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	bleve "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/bleve"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"time"
)

var (
	indexOutputPath string
	typeName        string

	dbConnectionDetails string
	databaseType        string

	deadline time.Duration

	validTypeNames = map[string]struct{}{
		"item": {},
	}

	validDatabaseTypes = map[string]struct{}{
		config.PostgresProviderKey: {},
		config.MariaDBProviderKey:  {},
		config.SqliteProviderKey:   {},
	}
)

const (
	outputPathVerboseFlagName   = "output"
	dbConnectionVerboseFlagName = "db_connection"
	dbTypeVerboseFlagName       = "db_type"
)

func init() {
	pflag.StringVarP(&indexOutputPath, outputPathVerboseFlagName, "o", "", "output path for bleve index")
	pflag.StringVarP(&typeName, "type", "t", "", "which type to create bleve index for")

	pflag.StringVarP(&dbConnectionDetails, dbConnectionVerboseFlagName, "c", "", "connection string for the relevant database")
	pflag.StringVarP(&databaseType, dbTypeVerboseFlagName, "b", "", "which type of database to connect to")

	pflag.DurationVarP(&deadline, "deadline", "d", time.Minute, "amount of time to spend adding to the index")
}

func main() {
	pflag.Parse()
	logger := zerolog.NewZeroLogger().WithName("search_index_initializer")
	ctx := context.Background()

	if indexOutputPath == "" {
		log.Fatalf("No output path specified, please provide one via the --%s flag", outputPathVerboseFlagName)
		return
	} else if _, ok := validTypeNames[typeName]; !ok {
		log.Fatalf("Invalid type name %q specified, one of [ 'item' ] expected", typeName)
		return
	} else if dbConnectionDetails == "" {
		log.Fatalf("No database connection details %q specified, please provide one via the --%s flag", dbConnectionDetails, dbConnectionVerboseFlagName)
		return
	} else if _, ok := validDatabaseTypes[databaseType]; !ok {
		log.Fatalf("Invalid database type %q specified, please provide one via the --%s flag", databaseType, dbTypeVerboseFlagName)
		return
	}

	im, err := bleve.NewBleveIndexManager(search.IndexPath(indexOutputPath), search.IndexName(typeName), logger)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.ServerConfig{
		Database: config.DatabaseSettings{
			Provider:          databaseType,
			ConnectionDetails: v1.ConnectionDetails(dbConnectionDetails),
		},
		Metrics: config.MetricsSettings{
			DBMetricsCollectionInterval: time.Second,
		},
	}

	// connect to our database.
	logger.Debug("connecting to database")
	rawDB, err := cfg.ProvideDatabaseConnection(logger)
	if err != nil {
		log.Fatalf("error establishing connection to database: %v", err)
	}

	// establish the database client.
	logger.Debug("setting up database client")
	dbClient, err := cfg.ProvideDatabaseClient(ctx, logger, rawDB)
	if err != nil {
		log.Fatalf("error initializing database client: %v", err)
	}

	switch typeName {
	case "item":
		outputChan := make(chan []v11.Item)
		if queryErr := dbClient.GetAllItems(ctx, outputChan); queryErr != nil {
			log.Fatalf("error fetching items from database: %v", err)
		}

		for {
			select {
			case items := <-outputChan:
				for _, x := range items {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	default:
		log.Fatal("this should never occur")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildVarDeclarations(proj)

		expected := `
package example

import (
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"time"
)

var (
	indexOutputPath string
	typeName        string

	dbConnectionDetails string
	databaseType        string

	deadline time.Duration

	validTypeNames = map[string]struct{}{
		"item": {},
	}

	validDatabaseTypes = map[string]struct{}{
		config.PostgresProviderKey: {},
		config.MariaDBProviderKey:  {},
		config.SqliteProviderKey:   {},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConstDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildConstDeclarations()

		expected := `
package example

import ()

const (
	outputPathVerboseFlagName   = "output"
	dbConnectionVerboseFlagName = "db_connection"
	dbTypeVerboseFlagName       = "db_type"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildInit()

		expected := `
package example

import (
	pflag "github.com/spf13/pflag"
	"time"
)

func init() {
	pflag.StringVarP(&indexOutputPath, outputPathVerboseFlagName, "o", "", "output path for bleve index")
	pflag.StringVarP(&typeName, "type", "t", "", "which type to create bleve index for")

	pflag.StringVarP(&dbConnectionDetails, dbConnectionVerboseFlagName, "c", "", "connection string for the relevant database")
	pflag.StringVarP(&databaseType, dbTypeVerboseFlagName, "b", "", "which type of database to connect to")

	pflag.DurationVarP(&deadline, "deadline", "d", time.Minute, "amount of time to spend adding to the index")
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildMain(proj)

		expected := `
package example

import (
	"context"
	pflag "github.com/spf13/pflag"
	zerolog "gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	bleve "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/bleve"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"log"
	"time"
)

func main() {
	pflag.Parse()
	logger := zerolog.NewZeroLogger().WithName("search_index_initializer")
	ctx := context.Background()

	if indexOutputPath == "" {
		log.Fatalf("No output path specified, please provide one via the --%s flag", outputPathVerboseFlagName)
		return
	} else if _, ok := validTypeNames[typeName]; !ok {
		log.Fatalf("Invalid type name %q specified, one of [ 'item' ] expected", typeName)
		return
	} else if dbConnectionDetails == "" {
		log.Fatalf("No database connection details %q specified, please provide one via the --%s flag", dbConnectionDetails, dbConnectionVerboseFlagName)
		return
	} else if _, ok := validDatabaseTypes[databaseType]; !ok {
		log.Fatalf("Invalid database type %q specified, please provide one via the --%s flag", databaseType, dbTypeVerboseFlagName)
		return
	}

	im, err := bleve.NewBleveIndexManager(search.IndexPath(indexOutputPath), search.IndexName(typeName), logger)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.ServerConfig{
		Database: config.DatabaseSettings{
			Provider:          databaseType,
			ConnectionDetails: v1.ConnectionDetails(dbConnectionDetails),
		},
		Metrics: config.MetricsSettings{
			DBMetricsCollectionInterval: time.Second,
		},
	}

	// connect to our database.
	logger.Debug("connecting to database")
	rawDB, err := cfg.ProvideDatabaseConnection(logger)
	if err != nil {
		log.Fatalf("error establishing connection to database: %v", err)
	}

	// establish the database client.
	logger.Debug("setting up database client")
	dbClient, err := cfg.ProvideDatabaseClient(ctx, logger, rawDB)
	if err != nil {
		log.Fatalf("error initializing database client: %v", err)
	}

	switch typeName {
	case "item":
		outputChan := make(chan []v11.Item)
		if queryErr := dbClient.GetAllItems(ctx, outputChan); queryErr != nil {
			log.Fatalf("error fetching items from database: %v", err)
		}

		for {
			select {
			case items := <-outputChan:
				for _, x := range items {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	default:
		log.Fatal("this should never occur")
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
