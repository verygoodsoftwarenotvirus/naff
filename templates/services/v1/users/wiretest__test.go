package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_wireTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := wireTestDotGo(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestProvideUserDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideUserDataManager(v1.BuildMockDatabase()))
	})
}

func TestProvideUserDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideUserDataServer(buildTestService(t)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideUserDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := buildTestProvideUserDataManager(proj)

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/database/v1"
	"testing"
)

func TestProvideUserDataManager(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideUserDataManager(v1.BuildMockDatabase()))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideUserDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestProvideUserDataServer()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestProvideUserDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		assert.NotNil(t, ProvideUserDataServer(buildTestService(t)))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
