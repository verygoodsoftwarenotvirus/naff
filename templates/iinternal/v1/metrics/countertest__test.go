package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_counterTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := counterTestDotGo(proj)

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"testing"
)

func Test_opencensusCounter_Increment(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ct, err := ProvideUnitCounter("v", "description")
		c, typOK := ct.(*opencensusCounter)
		require.NotNil(t, c)
		require.True(t, typOK)
		require.NoError(t, err)

		assert.Equal(t, c.actualCount, uint64(0))

		c.Increment(ctx)
		assert.Equal(t, c.actualCount, uint64(1))
	})
}

func Test_opencensusCounter_IncrementBy(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ct, err := ProvideUnitCounter("v", "description")
		c, typOK := ct.(*opencensusCounter)
		require.NotNil(t, c)
		require.True(t, typOK)
		require.NoError(t, err)

		assert.Equal(t, c.actualCount, uint64(0))

		c.IncrementBy(ctx, 666)
		assert.Equal(t, c.actualCount, uint64(666))
	})
}

func Test_opencensusCounter_Decrement(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ct, err := ProvideUnitCounter("v", "description")
		c, typOK := ct.(*opencensusCounter)
		require.NotNil(t, c)
		require.True(t, typOK)
		require.NoError(t, err)

		assert.Equal(t, c.actualCount, uint64(0))

		c.Increment(ctx)
		assert.Equal(t, c.actualCount, uint64(1))

		c.Decrement(ctx)
		assert.Equal(t, c.actualCount, uint64(0))
	})
}

func TestProvideUnitCounterProvider(t *testing.T) {
	t.Parallel()

	// obligatory.
	assert.NotNil(t, ProvideUnitCounterProvider())
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_opencensusCounter_Increment(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTest_opencensusCounter_Increment()

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"testing"
)

func Test_opencensusCounter_Increment(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ct, err := ProvideUnitCounter("v", "description")
		c, typOK := ct.(*opencensusCounter)
		require.NotNil(t, c)
		require.True(t, typOK)
		require.NoError(t, err)

		assert.Equal(t, c.actualCount, uint64(0))

		c.Increment(ctx)
		assert.Equal(t, c.actualCount, uint64(1))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_opencensusCounter_IncrementBy(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTest_opencensusCounter_IncrementBy()

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"testing"
)

func Test_opencensusCounter_IncrementBy(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ct, err := ProvideUnitCounter("v", "description")
		c, typOK := ct.(*opencensusCounter)
		require.NotNil(t, c)
		require.True(t, typOK)
		require.NoError(t, err)

		assert.Equal(t, c.actualCount, uint64(0))

		c.IncrementBy(ctx, 666)
		assert.Equal(t, c.actualCount, uint64(666))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTest_opencensusCounter_Decrement(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTest_opencensusCounter_Decrement()

		expected := `
package example

import (
	"context"
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"testing"
)

func Test_opencensusCounter_Decrement(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		ct, err := ProvideUnitCounter("v", "description")
		c, typOK := ct.(*opencensusCounter)
		require.NotNil(t, c)
		require.True(t, typOK)
		require.NoError(t, err)

		assert.Equal(t, c.actualCount, uint64(0))

		c.Increment(ctx)
		assert.Equal(t, c.actualCount, uint64(1))

		c.Decrement(ctx)
		assert.Equal(t, c.actualCount, uint64(0))
	})
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildTestProvideUnitCounterProvider(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildTestProvideUnitCounterProvider()

		expected := `
package example

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestProvideUnitCounterProvider(t *testing.T) {
	t.Parallel()

	// obligatory.
	assert.NotNil(t, ProvideUnitCounterProvider())
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
