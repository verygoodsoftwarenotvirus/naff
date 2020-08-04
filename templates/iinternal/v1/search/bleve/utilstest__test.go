package bleve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_utilsTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := utilsTestDotGo(proj)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestEnsureQueryIsRestrictedToUser(T *testing.T) {
	T.Parallel()

	T.Run("leaves good queries alone", func(t *testing.T) {
		exampleUserID := fake.BuildFakeUser().ID

		exampleQuery := fmt.Sprintf("things +belongsToUser:%d", exampleUserID)
		expectation := fmt.Sprintf("things +belongsToUser:%d", exampleUserID)

		actual := ensureQueryIsRestrictedToUser(exampleQuery, exampleUserID)
		assert.Equal(t, expectation, actual, "expected %q to equal %q", expectation, actual)
	})

	T.Run("basic replacement", func(t *testing.T) {
		exampleUserID := fake.BuildFakeUser().ID

		exampleQuery := "things"
		expectation := fmt.Sprintf("things +belongsToUser:%d", exampleUserID)

		actual := ensureQueryIsRestrictedToUser(exampleQuery, exampleUserID)
		assert.Equal(t, expectation, actual, "expected %q to equal %q", expectation, actual)
	})

	T.Run("with invalid user restriction", func(t *testing.T) {
		exampleUserID := fake.BuildFakeUser().ID

		exampleQuery := fmt.Sprintf("stuff belongsToUser:%d", exampleUserID)
		expectation := fmt.Sprintf("stuff +belongsToUser:%d", exampleUserID)

		actual := ensureQueryIsRestrictedToUser(exampleQuery, exampleUserID)
		assert.Equal(t, expectation, actual, "expected %q to equal %q", expectation, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestEnsureQueryIsRestrictedToUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildTestEnsureQueryIsRestrictedToUser(proj)

		expected := `
package example

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	fake "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1/fake"
	"testing"
)

func TestEnsureQueryIsRestrictedToUser(T *testing.T) {
	T.Parallel()

	T.Run("leaves good queries alone", func(t *testing.T) {
		exampleUserID := fake.BuildFakeUser().ID

		exampleQuery := fmt.Sprintf("things +belongsToUser:%d", exampleUserID)
		expectation := fmt.Sprintf("things +belongsToUser:%d", exampleUserID)

		actual := ensureQueryIsRestrictedToUser(exampleQuery, exampleUserID)
		assert.Equal(t, expectation, actual, "expected %q to equal %q", expectation, actual)
	})

	T.Run("basic replacement", func(t *testing.T) {
		exampleUserID := fake.BuildFakeUser().ID

		exampleQuery := "things"
		expectation := fmt.Sprintf("things +belongsToUser:%d", exampleUserID)

		actual := ensureQueryIsRestrictedToUser(exampleQuery, exampleUserID)
		assert.Equal(t, expectation, actual, "expected %q to equal %q", expectation, actual)
	})

	T.Run("with invalid user restriction", func(t *testing.T) {
		exampleUserID := fake.BuildFakeUser().ID

		exampleQuery := fmt.Sprintf("stuff belongsToUser:%d", exampleUserID)
		expectation := fmt.Sprintf("stuff +belongsToUser:%d", exampleUserID)

		actual := ensureQueryIsRestrictedToUser(exampleQuery, exampleUserID)
		assert.Equal(t, expectation, actual, "expected %q to equal %q", expectation, actual)
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
