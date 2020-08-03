package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_counterDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := counterDotGo(proj)

		expected := `
package example

import (
	"context"
	"fmt"
	stats "go.opencensus.io/stats"
	view "go.opencensus.io/stats/view"
	"sync/atomic"
)

// opencensusCounter is a Counter that interfaces with opencensus.
type opencensusCounter struct {
	name        string
	actualCount uint64
	measure     *stats.Int64Measure
	v           *view.View
}

func (c *opencensusCounter) subtractFromCount(ctx context.Context, value uint64) {
	atomic.AddUint64(&c.actualCount, ^value+1)
	stats.Record(ctx, c.measure.M(int64(-value)))
}

func (c *opencensusCounter) addToCount(ctx context.Context, value uint64) {
	atomic.AddUint64(&c.actualCount, value)
	stats.Record(ctx, c.measure.M(int64(value)))
}

// Decrement satisfies our Counter interface.
func (c *opencensusCounter) Decrement(ctx context.Context) {
	c.subtractFromCount(ctx, 1)
}

// Increment satisfies our Counter interface.
func (c *opencensusCounter) Increment(ctx context.Context) {
	c.addToCount(ctx, 1)
}

// IncrementBy satisfies our Counter interface.
func (c *opencensusCounter) IncrementBy(ctx context.Context, value uint64) {
	c.addToCount(ctx, value)
}

// ProvideUnitCounter provides a new counter.
func ProvideUnitCounter(counterName CounterName, description string) (UnitCounter, error) {
	name := fmt.Sprintf("%s_count", string(counterName))
	// Counts/groups the lengths of lines read in.
	count := stats.Int64(name, description, "By")

	countView := &view.View{
		Name:        name,
		Description: description,
		Measure:     count,
		Aggregation: view.Count(),
	}

	if err := view.Register(countView); err != nil {
		return nil, fmt.Errorf("failed to register views: %w", err)
	}

	c := &opencensusCounter{
		name:    name,
		measure: count,
		v:       countView,
	}

	return c, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildOpencensusCounter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildOpencensusCounter()

		expected := `
package example

import (
	stats "go.opencensus.io/stats"
	view "go.opencensus.io/stats/view"
)

// opencensusCounter is a Counter that interfaces with opencensus.
type opencensusCounter struct {
	name        string
	actualCount uint64
	measure     *stats.Int64Measure
	v           *view.View
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildSubtractFromCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildSubtractFromCount()

		expected := `
package example

import (
	"context"
	stats "go.opencensus.io/stats"
	"sync/atomic"
)

func (c *opencensusCounter) subtractFromCount(ctx context.Context, value uint64) {
	atomic.AddUint64(&c.actualCount, ^value+1)
	stats.Record(ctx, c.measure.M(int64(-value)))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildAddToCount(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildAddToCount()

		expected := `
package example

import (
	"context"
	stats "go.opencensus.io/stats"
	"sync/atomic"
)

func (c *opencensusCounter) addToCount(ctx context.Context, value uint64) {
	atomic.AddUint64(&c.actualCount, value)
	stats.Record(ctx, c.measure.M(int64(value)))
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildDecrement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildDecrement()

		expected := `
package example

import (
	"context"
)

// Decrement satisfies our Counter interface.
func (c *opencensusCounter) Decrement(ctx context.Context) {
	c.subtractFromCount(ctx, 1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIncrement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildIncrement()

		expected := `
package example

import (
	"context"
)

// Increment satisfies our Counter interface.
func (c *opencensusCounter) Increment(ctx context.Context) {
	c.addToCount(ctx, 1)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildIncrementBy(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildIncrementBy()

		expected := `
package example

import (
	"context"
)

// IncrementBy satisfies our Counter interface.
func (c *opencensusCounter) IncrementBy(ctx context.Context, value uint64) {
	c.addToCount(ctx, value)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideUnitCounter(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildProvideUnitCounter()

		expected := `
package example

import (
	"fmt"
	stats "go.opencensus.io/stats"
	view "go.opencensus.io/stats/view"
)

// ProvideUnitCounter provides a new counter.
func ProvideUnitCounter(counterName CounterName, description string) (UnitCounter, error) {
	name := fmt.Sprintf("%s_count", string(counterName))
	// Counts/groups the lengths of lines read in.
	count := stats.Int64(name, description, "By")

	countView := &view.View{
		Name:        name,
		Description: description,
		Measure:     count,
		Aggregation: view.Count(),
	}

	if err := view.Register(countView); err != nil {
		return nil, fmt.Errorf("failed to register views: %w", err)
	}

	c := &opencensusCounter{
		name:    name,
		measure: count,
		v:       countView,
	}

	return c, nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
