package utils

import "gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

func AddImports(file *jen.File) {
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/client/v1/http", "client")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1", "database")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth", "auth")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/auth/mock", "mockauth")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/config", "config")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding", "encoding")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/encoding/mock", "mockencoding")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics", "metrics")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/internal/v1/metrics/mock", "mockmetrics")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/client", "dbclient")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/mariadb", "mariadb")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/postgres", "postgres")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/database/v1/queriers/sqlite", "sqlite")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1", "models")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/models/v1/mock", "mockmodels")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1", "server")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/server/v1/http", "httpserver")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/auth", "auth")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/frontend", "frontend")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/items", "items")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/oauth2clients", "oauth2clients")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/users", "users")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/services/v1/webhooks", "webhooks")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/frontend", "frontend")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/integration", "integration")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/load", "load")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil", "testutil")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/mock", "mockutil")
	file.ImportAlias("gitlab.com/verygoodsoftwarenotvirus/todo/tests/v1/testutil/rand/model", "randmodel")

	file.ImportName("golang.org/x/oauth2", "oauth2")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/logging/v1", "logging")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop", "noop")
	file.ImportName("gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog", "zerolog")
	file.ImportName("github.com/stretchr/testify/assert", "assert")
	file.ImportName("github.com/stretchr/testify/require", "require")
	file.ImportName("github.com/stretchr/testify/mock", "mock")
	file.ImportName("go.opencensus.io/trace", "trace")

	file.ImportAlias("gopkg.in/oauth2.v3/models", "oauth2models")
	file.ImportAlias("gopkg.in/oauth2.v3/errors", "oauth2errors")
	file.ImportAlias("gopkg.in/oauth2.v3/server", "oauth2server")
	file.ImportAlias("gopkg.in/oauth2.v3/store", "oauth2store")

	file.ImportNames(map[string]string{
		"context":           "context",
		"fmt":               "fmt",
		"net/http":          "http",
		"net/http/httputil": "httputil",
		"errors":            "errors",
		"net/url":           "url",
		"path":              "path",
		"strings":           "strings",
		"time":              "time",
		"bytes":             "bytes",
		"encoding/json":     "json",
		"io":                "io",
		"io/ioutil":         "ioutil",
		"reflect":           "reflect",

		"contrib.go.opencensus.io/exporter/jaeger":     "jaeger",
		"contrib.go.opencensus.io/exporter/prometheus": "prometheus",
		"contrib.go.opencensus.io/integrations/ocsql":  "ocsql",
		"github.com/DATA-DOG/go-sqlmock":               "sqlmock",
		"github.com/GuiaBolso/darwin":                  "darwin",
		"github.com/Masterminds/squirrel":              "squirrel",
		"github.com/boombuler/barcode":                 "barcode",
		"github.com/emicklei/hazana":                   "hazana",
		"github.com/go-chi/chi":                        "chi",
		"github.com/go-chi/cors":                       "cors",
		"github.com/google/wire":                       "wire",
		"github.com/gorilla/securecookie":              "securecookie",
		"github.com/heptiolabs/healthcheck":            "healthcheck",
		"github.com/icrowley/fake":                     "fake",
		"github.com/lib/pq":                            "pq",
		"github.com/mattn/go-sqlite3":                  "sqlite3",
		"github.com/moul/http2curl":                    "http2curl",
		"github.com/pquerna/otp":                       "otp",
		"github.com/spf13/afero":                       "afero",
		"github.com/spf13/viper":                       "viper",
		"github.com/tebeka/selenium":                   "selenium",
		"gitlab.com/verygoodsoftwarenotvirus/newsman":  "newsman",
		"go.opencensus.io":                             "opencensus",
		"golang.org/x/crypto":                          "crypto",
		"gopkg.in/oauth2.v3":                           "oauth2",
		"go.opencensus.io/plugin/ochttp":               "ochttp",
		"github.com/pquerna/otp/totp":                  "totp",
		"golang.org/x/oauth2/clientcredentials":        "clientcredentials",
	})

	file.Line()
}
