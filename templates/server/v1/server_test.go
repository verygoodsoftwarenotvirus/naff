package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_serverDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := serverDotGo(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1/http"
)

type (
	// Server is the structure responsible for hosting all available protocols.
	// In the events we adopted a gRPC implementation of the surface, this is
	// the structure that would contain it and be responsible for calling its
	// serve method.
	Server struct {
		config     *config.ServerConfig
		httpServer *http.Server
	}
)

var (
	// Providers is our wire superset of providers this package offers.
	Providers = wire.NewSet(
		ProvideServer,
	)
)

// ProvideServer builds a new Server instance.
func ProvideServer(cfg *config.ServerConfig, httpServer *http.Server) (*Server, error) {
	srv := &Server{
		config:     cfg,
		httpServer: httpServer,
	}

	return srv, nil
}

// Serve serves HTTP traffic.
func (s *Server) Serve() {
	s.httpServer.Serve()
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTypeDefs(proj)

		expected := `
package example

import (
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1/http"
)

type (
	// Server is the structure responsible for hosting all available protocols.
	// In the events we adopted a gRPC implementation of the surface, this is
	// the structure that would contain it and be responsible for calling its
	// serve method.
	Server struct {
		config     *config.ServerConfig
		httpServer *http.Server
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildVarDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildVarDefs()

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is our wire superset of providers this package offers.
	Providers = wire.NewSet(
		ProvideServer,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideServer(proj)

		expected := `
package example

import (
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/server/v1/http"
)

// ProvideServer builds a new Server instance.
func ProvideServer(cfg *config.ServerConfig, httpServer *http.Server) (*Server, error) {
	srv := &Server{
		config:     cfg,
		httpServer: httpServer,
	}

	return srv, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildServe(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildServe()

		expected := `
package example

import ()

// Serve serves HTTP traffic.
func (s *Server) Serve() {
	s.httpServer.Serve()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
