package v1

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"fmt"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"log"
	"time"
)

const (
	defaultPort                      = 8888
	oneDay                           = 24 * time.Hour
	debugCookieSecret                = "HEREISA32CHARSECRETWHICHISMADEUP"
	defaultFrontendFilepath          = "/frontend"
	postgresDBConnDetails            = "postgres://dbuser:hunter2@database:5432/todo?sslmode=disable"
	metaDebug                        = "meta.debug"
	metaRunMode                      = "meta.run_mode"
	metaStartupDeadline              = "meta.startup_deadline"
	serverHTTPPort                   = "server.http_port"
	serverDebug                      = "server.debug"
	frontendDebug                    = "frontend.debug"
	frontendStaticFilesDir           = "frontend.static_files_directory"
	frontendCacheStatics             = "frontend.cache_static_files"
	authDebug                        = "auth.debug"
	authCookieDomain                 = "auth.cookie_domain"
	authCookieSecret                 = "auth.cookie_secret"
	authCookieLifetime               = "auth.cookie_lifetime"
	authSecureCookiesOnly            = "auth.secure_cookies_only"
	authEnableUserSignup             = "auth.enable_user_signup"
	metricsProvider                  = "metrics.metrics_provider"
	metricsTracer                    = "metrics.tracing_provider"
	metricsDBCollectionInterval      = "metrics.database_metrics_collection_interval"
	metricsRuntimeCollectionInterval = "metrics.runtime_metrics_collection_interval"
	dbDebug                          = "database.debug"
	dbProvider                       = "database.provider"
	dbDeets                          = "database.connection_details"
	itemsSearchIndexPath             = "search.items_index_path"

	// run modes
	developmentEnv = "development"
	testingEnv     = "testing"

	// database providers
	postgres = "postgres"
	sqlite   = "sqlite"
	mariadb  = "mariadb"

	// search index paths
	defaultItemsSearchIndexPath = "items.bleve"
)

type configFunc func(filePath string) error

var (
	files = map[string]configFunc{
		"environments/local/config.toml":                                    developmentConfig,
		"environments/testing/config_files/frontend-tests.toml":             frontendTestsConfig,
		"environments/testing/config_files/coverage.toml":                   coverageConfig,
		"environments/testing/config_files/integration-tests-postgres.toml": buildIntegrationTestForDBImplementation(postgres, postgresDBConnDetails),
		"environments/testing/config_files/integration-tests-sqlite.toml":   buildIntegrationTestForDBImplementation(sqlite, "/tmp/db"),
		"environments/testing/config_files/integration-tests-mariadb.toml":  buildIntegrationTestForDBImplementation(mariadb, "dbuser:hunter2@tcp(database:3306)/todo"),
	}
)

func developmentConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, developmentEnv)
	cfg.Set(metaDebug, true)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, true)
	cfg.Set(authCookieDomain, "localhost")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)

	cfg.Set(dbDebug, true)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing developmentEnv config: %w", writeErr)
	}

	return nil
}

func frontendTestsConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, developmentEnv)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, true)
	cfg.Set(authCookieDomain, "localhost")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)

	cfg.Set(dbDebug, true)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing developmentEnv config: %w", writeErr)
	}

	return nil
}

func coverageConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, testingEnv)
	cfg.Set(metaDebug, true)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, false)
	cfg.Set(authCookieSecret, debugCookieSecret)

	cfg.Set(dbDebug, false)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing coverage config: %w", writeErr)
	}

	return nil
}

func buildIntegrationTestForDBImplementation(dbVendor, dbDetails string) configFunc {
	return func(filePath string) error {
		cfg := config.BuildConfig()

		cfg.Set(metaRunMode, testingEnv)
		cfg.Set(metaDebug, false)

		sd := time.Minute
		if dbVendor == mariadb {
			sd = 5 * time.Minute
		}
		cfg.Set(metaStartupDeadline, sd)

		cfg.Set(serverHTTPPort, defaultPort)
		cfg.Set(serverDebug, true)

		cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
		cfg.Set(authCookieSecret, debugCookieSecret)

		cfg.Set(metricsProvider, "prometheus")
		cfg.Set(metricsTracer, "jaeger")

		cfg.Set(dbDebug, false)
		cfg.Set(dbProvider, dbVendor)
		cfg.Set(dbDeets, dbDetails)

		cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

		if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
			return fmt.Errorf("error writing integration test config for %s: %w", dbVendor, writeErr)
		}

		return nil
	}
}

func main() {
	for filePath, fun := range files {
		if err := fun(filePath); err != nil {
			log.Fatalf("error rendering %s: %v", filePath, err)
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_renderFileMap(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := renderFileMap(proj)

		expected := `
package example

import ()

var (
	files = map[string]configFunc{
		"environments/local/config.toml":                                    developmentConfig,
		"environments/testing/config_files/frontend-tests.toml":             frontendTestsConfig,
		"environments/testing/config_files/coverage.toml":                   coverageConfig,
		"environments/testing/config_files/integration-tests-postgres.toml": buildIntegrationTestForDBImplementation(postgres, postgresDBConnDetails),
		"environments/testing/config_files/integration-tests-sqlite.toml":   buildIntegrationTestForDBImplementation(sqlite, "/tmp/db"),
		"environments/testing/config_files/integration-tests-mariadb.toml":  buildIntegrationTestForDBImplementation(mariadb, "dbuser:hunter2@tcp(database:3306)/todo"),
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})

	T.Run("with no databases enabled", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.DisableDatabase(models.Postgres)
		proj.DisableDatabase(models.Sqlite)
		proj.DisableDatabase(models.MariaDB)
		x := renderFileMap(proj)

		expected := `
package example

import ()

var (
	files = map[string]configFunc{
		"environments/local/config.toml":                        developmentConfig,
		"environments/testing/config_files/frontend-tests.toml": frontendTestsConfig,
		"environments/testing/config_files/coverage.toml":       coverageConfig,
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_determineConstants(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := determineConstants(proj)

		expected := `
package example

import (
	"time"
)

const (
	defaultPort                      = 8888
	oneDay                           = 24 * time.Hour
	debugCookieSecret                = "HEREISA32CHARSECRETWHICHISMADEUP"
	defaultFrontendFilepath          = "/frontend"
	postgresDBConnDetails            = "postgres://dbuser:hunter2@database:5432/todo?sslmode=disable"
	metaDebug                        = "meta.debug"
	metaRunMode                      = "meta.run_mode"
	metaStartupDeadline              = "meta.startup_deadline"
	serverHTTPPort                   = "server.http_port"
	serverDebug                      = "server.debug"
	frontendDebug                    = "frontend.debug"
	frontendStaticFilesDir           = "frontend.static_files_directory"
	frontendCacheStatics             = "frontend.cache_static_files"
	authDebug                        = "auth.debug"
	authCookieDomain                 = "auth.cookie_domain"
	authCookieSecret                 = "auth.cookie_secret"
	authCookieLifetime               = "auth.cookie_lifetime"
	authSecureCookiesOnly            = "auth.secure_cookies_only"
	authEnableUserSignup             = "auth.enable_user_signup"
	metricsProvider                  = "metrics.metrics_provider"
	metricsTracer                    = "metrics.tracing_provider"
	metricsDBCollectionInterval      = "metrics.database_metrics_collection_interval"
	metricsRuntimeCollectionInterval = "metrics.runtime_metrics_collection_interval"
	dbDebug                          = "database.debug"
	dbProvider                       = "database.provider"
	dbDeets                          = "database.connection_details"
	itemsSearchIndexPath             = "search.items_index_path"

	// run modes
	developmentEnv = "development"
	testingEnv     = "testing"

	// database providers
	postgres = "postgres"
	sqlite   = "sqlite"
	mariadb  = "mariadb"

	// search index paths
	defaultItemsSearchIndexPath = "items.bleve"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildDevelopmentConfig(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildDevelopmentConfig(proj)

		expected := `
package example

import (
	"fmt"
	"time"
)

func developmentConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, developmentEnv)
	cfg.Set(metaDebug, true)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, true)
	cfg.Set(authCookieDomain, "localhost")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)

	cfg.Set(dbDebug, true)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing developmentEnv config: %w", writeErr)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildFrontendTestsConfig(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildFrontendTestsConfig(proj)

		expected := `
package example

import (
	"fmt"
	"time"
)

func frontendTestsConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, developmentEnv)
	cfg.Set(metaStartupDeadline, time.Minute)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, true)
	cfg.Set(authCookieDomain, "localhost")
	cfg.Set(authCookieSecret, debugCookieSecret)
	cfg.Set(authCookieLifetime, oneDay)
	cfg.Set(authSecureCookiesOnly, false)
	cfg.Set(authEnableUserSignup, true)

	cfg.Set(metricsProvider, "prometheus")
	cfg.Set(metricsTracer, "jaeger")
	cfg.Set(metricsDBCollectionInterval, time.Second)
	cfg.Set(metricsRuntimeCollectionInterval, time.Second)

	cfg.Set(dbDebug, true)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing developmentEnv config: %w", writeErr)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCoverageConfig(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildCoverageConfig(proj)

		expected := `
package example

import (
	"fmt"
)

func coverageConfig(filePath string) error {
	cfg := config.BuildConfig()

	cfg.Set(metaRunMode, testingEnv)
	cfg.Set(metaDebug, true)

	cfg.Set(serverHTTPPort, defaultPort)
	cfg.Set(serverDebug, true)

	cfg.Set(frontendDebug, true)
	cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
	cfg.Set(frontendCacheStatics, false)

	cfg.Set(authDebug, false)
	cfg.Set(authCookieSecret, debugCookieSecret)

	cfg.Set(dbDebug, false)
	cfg.Set(dbProvider, postgres)
	cfg.Set(dbDeets, postgresDBConnDetails)

	cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

	if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
		return fmt.Errorf("error writing coverage config: %w", writeErr)
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildIntegrationTestForDBImplementation(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildIntegrationTestForDBImplementation(proj)

		expected := `
package example

import (
	"fmt"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"time"
)

func buildIntegrationTestForDBImplementation(dbVendor, dbDetails string) configFunc {
	return func(filePath string) error {
		cfg := config.BuildConfig()

		cfg.Set(metaRunMode, testingEnv)
		cfg.Set(metaDebug, false)

		sd := time.Minute
		if dbVendor == mariadb {
			sd = 5 * time.Minute
		}
		cfg.Set(metaStartupDeadline, sd)

		cfg.Set(serverHTTPPort, defaultPort)
		cfg.Set(serverDebug, true)

		cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
		cfg.Set(authCookieSecret, debugCookieSecret)

		cfg.Set(metricsProvider, "prometheus")
		cfg.Set(metricsTracer, "jaeger")

		cfg.Set(dbDebug, false)
		cfg.Set(dbProvider, dbVendor)
		cfg.Set(dbDeets, dbDetails)

		cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

		if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
			return fmt.Errorf("error writing integration test config for %s: %w", dbVendor, writeErr)
		}

		return nil
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})

	T.Run("without MariaDB", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		proj.DisableDatabase(models.MariaDB)
		x := buildBuildIntegrationTestForDBImplementation(proj)

		expected := `
package example

import (
	"fmt"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	"time"
)

func buildIntegrationTestForDBImplementation(dbVendor, dbDetails string) configFunc {
	return func(filePath string) error {
		cfg := config.BuildConfig()

		cfg.Set(metaRunMode, testingEnv)
		cfg.Set(metaDebug, false)

		sd := time.Minute
		cfg.Set(metaStartupDeadline, sd)

		cfg.Set(serverHTTPPort, defaultPort)
		cfg.Set(serverDebug, true)

		cfg.Set(frontendStaticFilesDir, defaultFrontendFilepath)
		cfg.Set(authCookieSecret, debugCookieSecret)

		cfg.Set(metricsProvider, "prometheus")
		cfg.Set(metricsTracer, "jaeger")

		cfg.Set(dbDebug, false)
		cfg.Set(dbProvider, dbVendor)
		cfg.Set(dbDeets, dbDetails)

		cfg.Set(itemsSearchIndexPath, defaultItemsSearchIndexPath)

		if writeErr := cfg.WriteConfigAs(filePath); writeErr != nil {
			return fmt.Errorf("error writing integration test config for %s: %w", dbVendor, writeErr)
		}

		return nil
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMain()

		expected := `
package example

import (
	"log"
)

func main() {
	for filePath, fun := range files {
		if err := fun(filePath); err != nil {
			log.Fatalf("error rendering %s: %v", filePath, err)
		}
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
