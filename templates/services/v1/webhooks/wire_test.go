package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := wireDotGo(proj)

		expected := `
package example

import (
	wire "github.com/google/wire"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideWebhooksService,
		ProvideWebhookDataManager,
		ProvideWebhookDataServer,
	)
)

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataManager(db v1.DataManager) v11.WebhookDataManager {
	return db
}

// ProvideWebhookDataServer is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataServer(s *Service) v11.WebhookDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProviders(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProviders()

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideWebhooksService,
		ProvideWebhookDataManager,
		ProvideWebhookDataServer,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideWebhookDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideWebhookDataManager(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataManager(db v1.DataManager) v11.WebhookDataManager {
	return db
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideWebhookDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideWebhookDataServer(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideWebhookDataServer is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataServer(s *Service) v1.WebhookDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
