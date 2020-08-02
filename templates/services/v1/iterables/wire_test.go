package iterables

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
		typ := proj.DataTypes[0]
		x := wireDotGo(proj, typ)

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
		ProvideItemsService,
		ProvideItemDataManager,
		ProvideItemDataServer,
		ProvideItemsServiceSearchIndex,
	)
)

// ProvideItemDataManager turns a database into an ItemDataManager.
func ProvideItemDataManager(db v1.DataManager) v11.ItemDataManager {
	return db
}

// ProvideItemDataServer is an arbitrary function for dependency injection's sake.
func ProvideItemDataServer(s *Service) v11.ItemDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProviders(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildWireProviders(typ)

		expected := `
package example

import (
	wire "github.com/google/wire"
)

var (
	// Providers is our collection of what we provide to other services.
	Providers = wire.NewSet(
		ProvideItemsService,
		ProvideItemDataManager,
		ProvideItemDataServer,
		ProvideItemsServiceSearchIndex,
	)
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProvideSomethingDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildWireProvideSomethingDataManager(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	v11 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideItemDataManager turns a database into an ItemDataManager.
func ProvideItemDataManager(db v1.DataManager) v11.ItemDataManager {
	return db
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildWireProvideSomethingDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		typ := proj.DataTypes[0]
		x := buildWireProvideSomethingDataServer(proj, typ)

		expected := `
package example

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
)

// ProvideItemDataServer is an arbitrary function for dependency injection's sake.
func ProvideItemDataServer(s *Service) v1.ItemDataServer {
	return s
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
