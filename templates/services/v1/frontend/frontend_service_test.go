package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_frontendServiceDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := frontendServiceDotGo(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
)

const (
	serviceName = "frontend_service"
)

type (
	// Service is responsible for serving HTML (and other static resources)
	Service struct {
		logger v1.Logger
		config config.FrontendSettings
	}
)

// ProvideFrontendService provides the frontend service to dependency injection.
func ProvideFrontendService(logger v1.Logger, cfg config.FrontendSettings) *Service {
	svc := &Service{
		config: cfg,
		logger: logger.WithName(serviceName),
	}
	return svc
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFrontendConstantDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildFrontendConstantDefs()

		expected := `
package example

import ()

const (
	serviceName = "frontend_service"
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildFrontendTypeDefs(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildFrontendTypeDefs(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
)

type (
	// Service is responsible for serving HTML (and other static resources)
	Service struct {
		logger v1.Logger
		config config.FrontendSettings
	}
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideFrontendService(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildProvideFrontendService(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	config "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
)

// ProvideFrontendService provides the frontend service to dependency injection.
func ProvideFrontendService(logger v1.Logger, cfg config.FrontendSettings) *Service {
	svc := &Service{
		config: cfg,
		logger: logger.WithName(serviceName),
	}
	return svc
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
