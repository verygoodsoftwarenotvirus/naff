package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_runtimeTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := runtimeTestDotGo(proj)

		expected := `
package example

import (
	require "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRecordRuntimeStats(T *testing.T) {
	T.Parallel()

	// this is sort of an obligatory test for coverage's sake.

	d := time.Second
	sf := RecordRuntimeStats(d / 5)
	time.Sleep(d)
	sf()
}

func TestRegisterDefaultViews(t *testing.T) {
	t.Parallel()

	// obligatory
	require.NoError(t, RegisterDefaultViews())
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestRecordRuntimeStats(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestRecordRuntimeStats()

		expected := `
package example

import (
	"testing"
	"time"
)

func TestRecordRuntimeStats(T *testing.T) {
	T.Parallel()

	// this is sort of an obligatory test for coverage's sake.

	d := time.Second
	sf := RecordRuntimeStats(d / 5)
	time.Sleep(d)
	sf()
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestRegisterDefaultViews(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestRegisterDefaultViews()

		expected := `
package example

import (
	require "github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterDefaultViews(t *testing.T) {
	t.Parallel()

	// obligatory
	require.NoError(t, RegisterDefaultViews())
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
