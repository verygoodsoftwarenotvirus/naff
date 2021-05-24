package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_counterDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := counterDotGo(proj)

		expected := `
package example

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
)

var _ metrics.UnitCounter = (*UnitCounter)(nil)

// UnitCounter is a mock metrics.UnitCounter
type UnitCounter struct {
	mock.Mock
}

// Increment implements our UnitCounter interface.
func (m *UnitCounter) Increment(ctx context.Context) {
	m.Called(ctx)
}

// IncrementBy implements our UnitCounter interface.
func (m *UnitCounter) IncrementBy(ctx context.Context, val uint64) {
	m.Called(ctx, val)
}

// Decrement implements our UnitCounter interface.
func (m *UnitCounter) Decrement(ctx context.Context) {
	m.Called(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUnitCounter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildUnitCounter(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
)

var _ metrics.UnitCounter = (*UnitCounter)(nil)

// UnitCounter is a mock metrics.UnitCounter
type UnitCounter struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIncrement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildIncrement()

		expected := `
package example

import (
	"context"
)

// Increment implements our UnitCounter interface.
func (m *UnitCounter) Increment(ctx context.Context) {
	m.Called(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIncrementBy(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildIncrementBy()

		expected := `
package example

import (
	"context"
)

// IncrementBy implements our UnitCounter interface.
func (m *UnitCounter) IncrementBy(ctx context.Context, val uint64) {
	m.Called(ctx, val)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDecrement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildDecrement()

		expected := `
package example

import (
	"context"
)

// Decrement implements our UnitCounter interface.
func (m *UnitCounter) Decrement(ctx context.Context) {
	m.Called(ctx)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
