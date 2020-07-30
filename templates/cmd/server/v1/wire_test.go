package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	wire "github.com/google/wire"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	bleve "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/bleve"
	v12 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1/http"
	auth1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day.
func ProvideReporter(n *newsman.Newsman) newsman.Reporter {
	return n
}

// BuildServer builds a server.
func BuildServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	logger v1.Logger,
	database v11.DataManager,
	db *sql.DB,
) (*v12.Server, error) {
	wire.Build(
		config.Providers,
		auth.Providers,
		// server things,
		bleve.Providers,
		v12.Providers,
		encoding.Providers,
		http.Providers,
		// metrics,
		metrics.Providers,
		// external libs,
		newsman.NewNewsman,
		ProvideReporter,
		// services,
		auth1.Providers,
		users.Providers,
		items.Providers,
		frontend.Providers,
		webhooks.Providers,
		oauth2clients.Providers,
	)
	return nil, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideReporter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildProvideReporter()

		expected := `
package example

import (
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day.
func ProvideReporter(n *newsman.Newsman) newsman.Reporter {
	return n
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildBuildServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildBuildServer(proj)

		expected := `
package example

import (
	"context"
	"database/sql"
	wire "github.com/google/wire"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	bleve "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/bleve"
	v12 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1/http"
	auth1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// BuildServer builds a server.
func BuildServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	logger v1.Logger,
	database v11.DataManager,
	db *sql.DB,
) (*v12.Server, error) {
	wire.Build(
		config.Providers,
		auth.Providers,
		// server things,
		bleve.Providers,
		v12.Providers,
		encoding.Providers,
		http.Providers,
		// metrics,
		metrics.Providers,
		// external libs,
		newsman.NewNewsman,
		ProvideReporter,
		// services,
		auth1.Providers,
		users.Providers,
		items.Providers,
		frontend.Providers,
		webhooks.Providers,
		oauth2clients.Providers,
	)
	return nil, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireBuildCallArgs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildWireBuildCallArgs(proj)

		expected := `
package main

import (
	auth "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/auth"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	encoding "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/encoding"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	bleve "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/search/bleve"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1/http"
	auth1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/auth"
	frontend "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/frontend"
	items "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/items"
	oauth2clients "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/oauth2clients"
	users "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/users"
	webhooks "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/services/v1/webhooks"
	newsman "gitlab.com/verygoodsoftwarenotvirus/newsman"
)

func main() {
	exampleFunction(
		config.Providers,
		auth.Providers,
		// server things,
		bleve.Providers,
		v1.Providers,
		encoding.Providers,
		http.Providers,
		// metrics,
		metrics.Providers,
		// external libs,
		newsman.NewNewsman,
		ProvideReporter,
		// services,
		auth1.Providers,
		users.Providers,
		items.Providers,
		frontend.Providers,
		webhooks.Providers,
		oauth2clients.Providers,
	)
}
`
		actual := testutils.RenderCallArgsPerLineToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
