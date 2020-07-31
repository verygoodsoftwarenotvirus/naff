package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_spansDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := spansDotGo(proj)

		expected := `
package example

import (
	"context"
	trace "go.opencensus.io/trace"
)

// StartSpan starts a span.
func StartSpan(ctx context.Context, funcName string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, funcName)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildStartSpan(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildStartSpan()

		expected := `
package example

import (
	"context"
	trace "go.opencensus.io/trace"
)

// StartSpan starts a span.
func StartSpan(ctx context.Context, funcName string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, funcName)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
