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

		expected := "github.com/verygoodsoftwarenotvirus/example/pkg/client/httpclient"
		actual := p.HTTPClientPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/pkg/client/httpclient/fart"
		actual := p.HTTPClientPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ModelsV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/pkg/types"
		actual := p.TypesPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/pkg/types/fart"
		actual := p.TypesPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_FakeModelsPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/pkg/types/fake"
		actual := p.FakeTypesPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/pkg/types/fake/fart"
		actual := p.FakeTypesPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_DatabaseV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/database"
		actual := p.DatabasePackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/database/fart"
		actual := p.DatabasePackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal"
		actual := p.InternalPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/fart"
		actual := p.InternalPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalAuthV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/auth"
		actual := p.InternalAuthenticationPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/auth/fart"
		actual := p.InternalAuthenticationPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalConfigV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/config"
		actual := p.InternalConfigPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/config/fart"
		actual := p.InternalConfigPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalEncodingV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/encoding"
		actual := p.EncodingPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/encoding/fart"
		actual := p.EncodingPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalMetricsV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/metrics"
		actual := p.MetricsPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/metrics/fart"
		actual := p.MetricsPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalTracingV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/tracing"
		actual := p.InternalTracingPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/tracing/fart"
		actual := p.InternalTracingPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_InternalSearchV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/search"
		actual := p.InternalSearchPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/search/fart"
		actual := p.InternalSearchPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services"
		actual := p.ServicePackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/fart"
		actual := p.ServicePackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1AuthPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/auth"
		actual := p.AuthServicePackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/auth/fart"
		actual := p.AuthServicePackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1FrontendPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/frontend"
		actual := p.FrontendServicePackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/frontend/fart"
		actual := p.FrontendServicePackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1OAuth2ClientsPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/oauth2clients"
		actual := p.ServiceOAuth2ClientsPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/oauth2clients/fart"
		actual := p.ServiceOAuth2ClientsPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1UsersPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/users"
		actual := p.UsersServicePackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/users/fart"
		actual := p.UsersServicePackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_ServiceV1WebhooksPackage(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/webhooks"
		actual := p.WebhooksServicePackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/internal/services/webhooks/fart"
		actual := p.WebhooksServicePackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}

func TestProject_TestUtilV1Package(T *testing.T) {
	T.Parallel()

	T.Run("empty input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/tests/utils"
		actual := p.TestUtilsPackage()

		assert.Equal(t, expected, actual)
	})

	T.Run("example input", func(t *testing.T) {
		t.Parallel()

		p := buildTestProjectForPathTests()

		expected := "github.com/verygoodsoftwarenotvirus/example/tests/utils/fart"
		actual := p.TestUtilsPackage(examplePathTestInput)

		assert.Equal(t, expected, actual)
	})
}
