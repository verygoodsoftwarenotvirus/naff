package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_configDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := configDotGo(proj)

		expected := `
package example

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	viper "github.com/spf13/viper"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	"io/ioutil"
	"time"
)

const (
	// DevelopmentRunMode is the run mode for a development environment
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment
	ProductionRunMode runMode = "production"

	defaultStartupDeadline                   = time.Minute
	defaultRunMode                           = DevelopmentRunMode
	defaultCookieLifetime                    = 24 * time.Hour
	defaultMetricsCollectionInterval         = 2 * time.Second
	defaultDatabaseMetricsCollectionInterval = 2 * time.Second
	randStringSize                           = 32
)

var (
	validModes = map[runMode]struct{}{
		DevelopmentRunMode: {},
		TestingRunMode:     {},
		ProductionRunMode:  {},
	}
)

func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}

type (
	runMode string

	// MetaSettings is primarily used for development.
	MetaSettings struct {
		// Debug enables debug mode service-wide
		// NOTE: this debug should override all other debugs, which is to say, if this is enabled, all of them are enabled.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// StartupDeadline indicates how long the service can take to spin up. This includes database migrations, configuring services, etc.
		StartupDeadline time.Duration ` + "`" + `json:"startup_deadline" mapstructure:"startup_deadline" toml:"startup_deadline,omitempty"` + "`" + `
		// RunMode indicates the current run mode
		RunMode runMode ` + "`" + `json:"run_mode" mapstructure:"run_mode" toml:"run_mode,omitempty"` + "`" + `
	}

	// ServerSettings describes the settings pertinent to the HTTP serving portion of the service.
	ServerSettings struct {
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// HTTPPort indicates which port to serve HTTP traffic on.
		HTTPPort uint16 ` + "`" + `json:"http_port" mapstructure:"http_port" toml:"http_port,omitempty"` + "`" + `
	}

	// FrontendSettings describes the settings pertinent to the frontend.
	FrontendSettings struct {
		// StaticFilesDirectory indicates which directory contains our static files for the frontend (i.e. CSS/JS/HTML files)
		StaticFilesDirectory string ` + "`" + `json:"static_files_directory" mapstructure:"static_files_directory" toml:"static_files_directory,omitempty"` + "`" + `
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// CacheStaticFiles indicates whether or not to load the static files directory into memory via afero's MemMapFs.
		CacheStaticFiles bool ` + "`" + `json:"cache_static_files" mapstructure:"cache_static_files" toml:"cache_static_files,omitempty"` + "`" + `
	}

	// AuthSettings represents our authentication configuration.
	AuthSettings struct {
		// CookieDomain indicates what domain the cookies will have set for them.
		CookieDomain string ` + "`" + `json:"cookie_domain" mapstructure:"cookie_domain" toml:"cookie_domain,omitempty"` + "`" + `
		// CookieSecret indicates the secret the cookie builder should use.
		CookieSecret string ` + "`" + `json:"cookie_secret" mapstructure:"cookie_secret" toml:"cookie_secret,omitempty"` + "`" + `
		// CookieLifetime indicates how long the cookies built should last.
		CookieLifetime time.Duration ` + "`" + `json:"cookie_lifetime" mapstructure:"cookie_lifetime" toml:"cookie_lifetime,omitempty"` + "`" + `
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// SecureCookiesOnly indicates if the cookies built should be marked as HTTPS only.
		SecureCookiesOnly bool ` + "`" + `json:"secure_cookies_only" mapstructure:"secure_cookies_only" toml:"secure_cookies_only,omitempty"` + "`" + `
		// EnableUserSignup enables user signups.
		EnableUserSignup bool ` + "`" + `json:"enable_user_signup" mapstructure:"enable_user_signup" toml:"enable_user_signup,omitempty"` + "`" + `
	}

	// DatabaseSettings represents our database configuration.
	DatabaseSettings struct {
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// Provider indicates what database we'll connect to (postgres, mysql, etc.)
		Provider string ` + "`" + `json:"provider" mapstructure:"provider" toml:"provider,omitempty"` + "`" + `
		// ConnectionDetails indicates how our database driver should connect to the instance.
		ConnectionDetails v1.ConnectionDetails ` + "`" + `json:"connection_details" mapstructure:"connection_details" toml:"connection_details,omitempty"` + "`" + `
	}

	// MetricsSettings contains settings about how we report our metrics.
	MetricsSettings struct {
		// MetricsProvider indicates where our metrics should go.
		MetricsProvider metricsProvider ` + "`" + `json:"metrics_provider" mapstructure:"metrics_provider" toml:"metrics_provider,omitempty"` + "`" + `
		// TracingProvider indicates where our traces should go.
		TracingProvider tracingProvider ` + "`" + `json:"tracing_provider" mapstructure:"tracing_provider" toml:"tracing_provider,omitempty"` + "`" + `
		// DBMetricsCollectionInterval is the interval we collect database statistics at.
		DBMetricsCollectionInterval time.Duration ` + "`" + `json:"database_metrics_collection_interval" mapstructure:"database_metrics_collection_interval" toml:"database_metrics_collection_interval,omitempty"` + "`" + `
		// RuntimeMetricsCollectionInterval  is the interval we collect runtime statistics at.
		RuntimeMetricsCollectionInterval time.Duration ` + "`" + `json:"runtime_metrics_collection_interval" mapstructure:"runtime_metrics_collection_interval" toml:"runtime_metrics_collection_interval,omitempty"` + "`" + `
	}

	// SearchSettings contains settings regarding search indices.
	SearchSettings struct {
		// ItemsIndexPath indicates where our items search index files should go.
		ItemsIndexPath search.IndexPath ` + "`" + `json:"items_index_path" mapstructure:"items_index_path" toml:"items_index_path,omitempty"` + "`" + `
	}

	// ServerConfig is our server configuration struct. It is comprised of all the other setting structs
	// For information on this structs fields, refer to their definitions.
	ServerConfig struct {
		Meta     MetaSettings     ` + "`" + `json:"meta" mapstructure:"meta" toml:"meta,omitempty"` + "`" + `
		Frontend FrontendSettings ` + "`" + `json:"frontend" mapstructure:"frontend" toml:"frontend,omitempty"` + "`" + `
		Auth     AuthSettings     ` + "`" + `json:"auth" mapstructure:"auth" toml:"auth,omitempty"` + "`" + `
		Server   ServerSettings   ` + "`" + `json:"server" mapstructure:"server" toml:"server,omitempty"` + "`" + `
		Database DatabaseSettings ` + "`" + `json:"database" mapstructure:"database" toml:"database,omitempty"` + "`" + `
		Metrics  MetricsSettings  ` + "`" + `json:"metrics" mapstructure:"metrics" toml:"metrics,omitempty"` + "`" + `
		Search   SearchSettings   ` + "`" + `json:"search" mapstructure:"search" toml:"search,omitempty"` + "`" + `
	}
)

// EncodeToFile renders your config to a file given your favorite encoder.
func (cfg *ServerConfig) EncodeToFile(path string, marshaler func(v interface{}) ([]byte, error)) error {
	byteSlice, err := marshaler(*cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, byteSlice, 0600)
}

// BuildConfig is a constructor function that initializes a viper config.
func BuildConfig() *viper.Viper {
	cfg := viper.New()

	// meta stuff.
	cfg.SetDefault("meta.run_mode", defaultRunMode)
	cfg.SetDefault("meta.startup_deadline", defaultStartupDeadline)

	// auth stuff.
	cfg.SetDefault("auth.cookie_lifetime", defaultCookieLifetime)
	cfg.SetDefault("auth.enable_user_signup", true)

	// metrics stuff.
	cfg.SetDefault("metrics.database_metrics_collection_interval", defaultMetricsCollectionInterval)
	cfg.SetDefault("metrics.runtime_metrics_collection_interval", defaultDatabaseMetricsCollectionInterval)

	// server stuff.
	cfg.SetDefault("server.http_port", 80)

	return cfg
}

// ParseConfigFile parses a configuration file.
func ParseConfigFile(filename string) (*ServerConfig, error) {
	cfg := BuildConfig()
	cfg.SetConfigFile(filename)

	if err := cfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("trying to read the config file: %w", err)
	}

	var serverConfig *ServerConfig
	if err := cfg.Unmarshal(&serverConfig); err != nil {
		return nil, fmt.Errorf("trying to unmarshal the config: %w", err)
	}

	if _, ok := validModes[serverConfig.Meta.RunMode]; !ok {
		return nil, fmt.Errorf("invalid run mode: %q", serverConfig.Meta.RunMode)
	}

	// set the cookie secret to something (relatively) secure if not provided
	if serverConfig.Auth.CookieSecret == "" {
		serverConfig.Auth.CookieSecret = randString(randStringSize)
	}

	return serverConfig, nil
}

// randString produces a random string.
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func randString(size uint) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConfigConstantVariableDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildConfigConstantVariableDeclarations()

		expected := `
package example

import (
	"time"
)

const (
	// DevelopmentRunMode is the run mode for a development environment
	DevelopmentRunMode runMode = "development"
	// TestingRunMode is the run mode for a testing environment
	TestingRunMode runMode = "testing"
	// ProductionRunMode is the run mode for a production environment
	ProductionRunMode runMode = "production"

	defaultStartupDeadline                   = time.Minute
	defaultRunMode                           = DevelopmentRunMode
	defaultCookieLifetime                    = 24 * time.Hour
	defaultMetricsCollectionInterval         = 2 * time.Second
	defaultDatabaseMetricsCollectionInterval = 2 * time.Second
	randStringSize                           = 32
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildConfigVariableDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildConfigVariableDeclarations()

		expected := `
package example

import ()

var (
	validModes = map[runMode]struct{}{
		DevelopmentRunMode: {},
		TestingRunMode:     {},
		ProductionRunMode:  {},
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildInit(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildInit()

		expected := `
package example

import (
	"crypto/rand"
)

func init() {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTypeDefinitions(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTypeDefinitions(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	search "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search"
	"time"
)

type (
	runMode string

	// MetaSettings is primarily used for development.
	MetaSettings struct {
		// Debug enables debug mode service-wide
		// NOTE: this debug should override all other debugs, which is to say, if this is enabled, all of them are enabled.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// StartupDeadline indicates how long the service can take to spin up. This includes database migrations, configuring services, etc.
		StartupDeadline time.Duration ` + "`" + `json:"startup_deadline" mapstructure:"startup_deadline" toml:"startup_deadline,omitempty"` + "`" + `
		// RunMode indicates the current run mode
		RunMode runMode ` + "`" + `json:"run_mode" mapstructure:"run_mode" toml:"run_mode,omitempty"` + "`" + `
	}

	// ServerSettings describes the settings pertinent to the HTTP serving portion of the service.
	ServerSettings struct {
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// HTTPPort indicates which port to serve HTTP traffic on.
		HTTPPort uint16 ` + "`" + `json:"http_port" mapstructure:"http_port" toml:"http_port,omitempty"` + "`" + `
	}

	// FrontendSettings describes the settings pertinent to the frontend.
	FrontendSettings struct {
		// StaticFilesDirectory indicates which directory contains our static files for the frontend (i.e. CSS/JS/HTML files)
		StaticFilesDirectory string ` + "`" + `json:"static_files_directory" mapstructure:"static_files_directory" toml:"static_files_directory,omitempty"` + "`" + `
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// CacheStaticFiles indicates whether or not to load the static files directory into memory via afero's MemMapFs.
		CacheStaticFiles bool ` + "`" + `json:"cache_static_files" mapstructure:"cache_static_files" toml:"cache_static_files,omitempty"` + "`" + `
	}

	// AuthSettings represents our authentication configuration.
	AuthSettings struct {
		// CookieDomain indicates what domain the cookies will have set for them.
		CookieDomain string ` + "`" + `json:"cookie_domain" mapstructure:"cookie_domain" toml:"cookie_domain,omitempty"` + "`" + `
		// CookieSecret indicates the secret the cookie builder should use.
		CookieSecret string ` + "`" + `json:"cookie_secret" mapstructure:"cookie_secret" toml:"cookie_secret,omitempty"` + "`" + `
		// CookieLifetime indicates how long the cookies built should last.
		CookieLifetime time.Duration ` + "`" + `json:"cookie_lifetime" mapstructure:"cookie_lifetime" toml:"cookie_lifetime,omitempty"` + "`" + `
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// SecureCookiesOnly indicates if the cookies built should be marked as HTTPS only.
		SecureCookiesOnly bool ` + "`" + `json:"secure_cookies_only" mapstructure:"secure_cookies_only" toml:"secure_cookies_only,omitempty"` + "`" + `
		// EnableUserSignup enables user signups.
		EnableUserSignup bool ` + "`" + `json:"enable_user_signup" mapstructure:"enable_user_signup" toml:"enable_user_signup,omitempty"` + "`" + `
	}

	// DatabaseSettings represents our database configuration.
	DatabaseSettings struct {
		// Debug determines if debug logging or other development conditions are active.
		Debug bool ` + "`" + `json:"debug" mapstructure:"debug" toml:"debug,omitempty"` + "`" + `
		// Provider indicates what database we'll connect to (postgres, mysql, etc.)
		Provider string ` + "`" + `json:"provider" mapstructure:"provider" toml:"provider,omitempty"` + "`" + `
		// ConnectionDetails indicates how our database driver should connect to the instance.
		ConnectionDetails v1.ConnectionDetails ` + "`" + `json:"connection_details" mapstructure:"connection_details" toml:"connection_details,omitempty"` + "`" + `
	}

	// MetricsSettings contains settings about how we report our metrics.
	MetricsSettings struct {
		// MetricsProvider indicates where our metrics should go.
		MetricsProvider metricsProvider ` + "`" + `json:"metrics_provider" mapstructure:"metrics_provider" toml:"metrics_provider,omitempty"` + "`" + `
		// TracingProvider indicates where our traces should go.
		TracingProvider tracingProvider ` + "`" + `json:"tracing_provider" mapstructure:"tracing_provider" toml:"tracing_provider,omitempty"` + "`" + `
		// DBMetricsCollectionInterval is the interval we collect database statistics at.
		DBMetricsCollectionInterval time.Duration ` + "`" + `json:"database_metrics_collection_interval" mapstructure:"database_metrics_collection_interval" toml:"database_metrics_collection_interval,omitempty"` + "`" + `
		// RuntimeMetricsCollectionInterval  is the interval we collect runtime statistics at.
		RuntimeMetricsCollectionInterval time.Duration ` + "`" + `json:"runtime_metrics_collection_interval" mapstructure:"runtime_metrics_collection_interval" toml:"runtime_metrics_collection_interval,omitempty"` + "`" + `
	}

	// SearchSettings contains settings regarding search indices.
	SearchSettings struct {
		// ItemsIndexPath indicates where our items search index files should go.
		ItemsIndexPath search.IndexPath ` + "`" + `json:"items_index_path" mapstructure:"items_index_path" toml:"items_index_path,omitempty"` + "`" + `
	}

	// ServerConfig is our server configuration struct. It is comprised of all the other setting structs
	// For information on this structs fields, refer to their definitions.
	ServerConfig struct {
		Meta     MetaSettings     ` + "`" + `json:"meta" mapstructure:"meta" toml:"meta,omitempty"` + "`" + `
		Frontend FrontendSettings ` + "`" + `json:"frontend" mapstructure:"frontend" toml:"frontend,omitempty"` + "`" + `
		Auth     AuthSettings     ` + "`" + `json:"auth" mapstructure:"auth" toml:"auth,omitempty"` + "`" + `
		Server   ServerSettings   ` + "`" + `json:"server" mapstructure:"server" toml:"server,omitempty"` + "`" + `
		Database DatabaseSettings ` + "`" + `json:"database" mapstructure:"database" toml:"database,omitempty"` + "`" + `
		Metrics  MetricsSettings  ` + "`" + `json:"metrics" mapstructure:"metrics" toml:"metrics,omitempty"` + "`" + `
		Search   SearchSettings   ` + "`" + `json:"search" mapstructure:"search" toml:"search,omitempty"` + "`" + `
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildEncodeToFile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildEncodeToFile()

		expected := `
package example

import (
	"io/ioutil"
)

// EncodeToFile renders your config to a file given your favorite encoder.
func (cfg *ServerConfig) EncodeToFile(path string, marshaler func(v interface{}) ([]byte, error)) error {
	byteSlice, err := marshaler(*cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, byteSlice, 0600)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildConfig(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildBuildConfig()

		expected := `
package example

import (
	viper "github.com/spf13/viper"
)

// BuildConfig is a constructor function that initializes a viper config.
func BuildConfig() *viper.Viper {
	cfg := viper.New()

	// meta stuff.
	cfg.SetDefault("meta.run_mode", defaultRunMode)
	cfg.SetDefault("meta.startup_deadline", defaultStartupDeadline)

	// auth stuff.
	cfg.SetDefault("auth.cookie_lifetime", defaultCookieLifetime)
	cfg.SetDefault("auth.enable_user_signup", true)

	// metrics stuff.
	cfg.SetDefault("metrics.database_metrics_collection_interval", defaultMetricsCollectionInterval)
	cfg.SetDefault("metrics.runtime_metrics_collection_interval", defaultDatabaseMetricsCollectionInterval)

	// server stuff.
	cfg.SetDefault("server.http_port", 80)

	return cfg
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildParseConfigFile(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildParseConfigFile()

		expected := `
package example

import (
	"fmt"
)

// ParseConfigFile parses a configuration file.
func ParseConfigFile(filename string) (*ServerConfig, error) {
	cfg := BuildConfig()
	cfg.SetConfigFile(filename)

	if err := cfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("trying to read the config file: %w", err)
	}

	var serverConfig *ServerConfig
	if err := cfg.Unmarshal(&serverConfig); err != nil {
		return nil, fmt.Errorf("trying to unmarshal the config: %w", err)
	}

	if _, ok := validModes[serverConfig.Meta.RunMode]; !ok {
		return nil, fmt.Errorf("invalid run mode: %q", serverConfig.Meta.RunMode)
	}

	// set the cookie secret to something (relatively) secure if not provided
	if serverConfig.Auth.CookieSecret == "" {
		serverConfig.Auth.CookieSecret = randString(randStringSize)
	}

	return serverConfig, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildRandString(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildRandString()

		expected := `
package example

import (
	"crypto/rand"
	"encoding/base32"
)

// randString produces a random string.
// https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/
func randString(size uint) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
