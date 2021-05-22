package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_metaTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := metaTestDotGo(proj)

		expected := `
package example

import (
	require "github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestHoldOnForever(T *testing.T) {
	T.Parallel()

	if os.Getenv("WAIT_FOR_COVERAGE") == "yes" {
		// snooze for a year.
		time.Sleep(time.Hour * 24 * 365)
	}
}

func checkValueAndError(t *testing.T, i interface{}, err error) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, i)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestHoldOnForever(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestHoldOnForever()

		expected := `
package example

import (
	"os"
	"testing"
	"time"
)

func TestHoldOnForever(T *testing.T) {
	T.Parallel()

	if os.Getenv("WAIT_FOR_COVERAGE") == "yes" {
		// snooze for a year.
		time.Sleep(time.Hour * 24 * 365)
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildCheckValueAndError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildCheckValueAndError()

		expected := `
package example

import (
	require "github.com/stretchr/testify/require"
	"testing"
)

func checkValueAndError(t *testing.T, i interface{}, err error) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, i)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
