package users

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
	// Providers is what we provide for dependency injectors.
	Providers = wire.NewSet(
		ProvideUsersService,
		ProvideUserDataServer,
		ProvideUserDataManager,
	)
)

// ProvideUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserDataManager(db v1.DataManager) v11.UserDataManager {
	return db
}

// ProvideUserDataServer is an arbitrary function for dependency injection's sake.
func ProvideUserDataServer(s *Service) v11.UserDataServer {
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
	// Providers is what we provide for dependency injectors.
	Providers = wire.NewSet(
		ProvideUsersService,
		ProvideUserDataServer,
		ProvideUserDataManager,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideUserDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideUserDataManager(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserDataManager(db v1.DataManager) v11.UserDataManager {
	return db
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideUserDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideUserDataServer(proj)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideUserDataServer is an arbitrary function for dependency injection's sake.
func ProvideUserDataServer(s *Service) v1.UserDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
