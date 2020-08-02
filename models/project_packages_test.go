package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	examplePathForPathTests = "github.com/verygoodsoftwarenotvirus/example"
	examplePathTestInput    = "fart"
)

func buildTestProjectForPathTests() *Project {
	return &Project{
		OutputPath: examplePathForPathTests,
	}
}

func TestProject_RelativePath(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := examplePathForPathTests
		actual := p.RelativePath()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/fart"
		actual := p.RelativePath("fart")

		assert.Equal(t, expected, actual)
	})
}

func TestProject_HTTPClientV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/client/v1/http"
		actual := p.HTTPClientV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/client/v1/http/fart"
		actual := p.HTTPClientV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ModelsV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/models/v1"
		actual := p.ModelsV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/models/v1/fart"
		actual := p.ModelsV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_FakeModelsPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/models/v1/fake"
		actual := p.FakeModelsPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/models/v1/fake/fart"
		actual := p.FakeModelsPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_DatabaseV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/database/v1"
		actual := p.DatabaseV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/database/v1/fart"
		actual := p.DatabaseV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1"
		actual := p.InternalV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/fart"
		actual := p.InternalV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalAuthV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/auth"
		actual := p.InternalAuthV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/auth/fart"
		actual := p.InternalAuthV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalConfigV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/config"
		actual := p.InternalConfigV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/config/fart"
		actual := p.InternalConfigV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalEncodingV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/encoding"
		actual := p.InternalEncodingV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/encoding/fart"
		actual := p.InternalEncodingV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalMetricsV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/metrics"
		actual := p.InternalMetricsV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/metrics/fart"
		actual := p.InternalMetricsV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalTracingV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/tracing"
		actual := p.InternalTracingV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/tracing/fart"
		actual := p.InternalTracingV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalSearchV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/search"
		actual := p.InternalSearchV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/v1/search/fart"
		actual := p.InternalSearchV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1"
		actual := p.ServiceV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/fart"
		actual := p.ServiceV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1AuthPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/auth"
		actual := p.ServiceV1AuthPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/auth/fart"
		actual := p.ServiceV1AuthPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1FrontendPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/frontend"
		actual := p.ServiceV1FrontendPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/frontend/fart"
		actual := p.ServiceV1FrontendPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1OAuth2ClientsPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/oauth2clients"
		actual := p.ServiceV1OAuth2ClientsPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/oauth2clients/fart"
		actual := p.ServiceV1OAuth2ClientsPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1UsersPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/users"
		actual := p.ServiceV1UsersPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/users/fart"
		actual := p.ServiceV1UsersPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1WebhooksPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/webhooks"
		actual := p.ServiceV1WebhooksPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/services/v1/webhooks/fart"
		actual := p.ServiceV1WebhooksPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_TestUtilV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/tests/v1/testutil"
		actual := p.TestUtilV1Package()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/tests/v1/testutil/fart"
		actual := p.TestUtilV1Package(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}
