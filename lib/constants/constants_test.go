package constants

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"

	"github.com/stretchr/testify/assert"
)

// end helper funcs

func TestCreateCtx(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := CreateCtx()

		expected := `
package main

import (
	"context"
)

func main() {
	ctx := context.Background()
}
`
		actual := testutils.RenderIndependentStatementToString(t, result)

		assert.Equal(t, expected, actual)
	})
}

func TestCtxParam(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := CtxParam()

		expected := `
package main

import (
	"context"
)

func example(ctx context.Context) {}
`
		actual := testutils.RenderFunctionParamsToString(t, []jen.Code{result})

		assert.Equal(t, expected, actual)
	})
}

func TestCtxVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := CtxVar()

		expected := `
package main

import ()

func main() {
	exampleFunction(ctx)
}
`
		actual := testutils.RenderCallArgsToString(t, []jen.Code{result})

		assert.Equal(t, expected, actual)
	})
}

func TestInlineCtx(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := InlineCtx()

		expected := `
package main

import (
	"context"
)

func main() {
	context.Background()
}
`
		actual := testutils.RenderIndependentStatementToString(t, result)

		assert.Equal(t, expected, actual)
	})
}

func TestLoggerParam(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := models.LoggerParam()

		expected := `
package main

import (
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

func example(logger v1.Logger) {}
`
		actual := testutils.RenderFunctionParamsToString(t, []jen.Code{result})

		assert.Equal(t, expected, actual)
	})
}

func TestObligatoryError(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := ObligatoryError()

		expected := `
package main

import (
	"errors"
)

func main() {
	errors.New("blah")
}
`
		actual := testutils.RenderIndependentStatementToString(t, result)

		assert.Equal(t, expected, actual)
	})
}

func TestUserIDParam(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := UserIDParam()

		expected := `
package main

import ()

func example(userID uint64) {}
`
		actual := testutils.RenderFunctionParamsToString(t, []jen.Code{result})

		assert.Equal(t, expected, actual)
	})
}

func TestUserIDVar(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := UserIDVar()

		expected := `
package main

import ()

func main() {
	exampleFunction(userID)
}
`
		actual := testutils.RenderCallArgsToString(t, []jen.Code{result})

		assert.Equal(t, expected, actual)
	})
}

func Test_err(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		result := err("example")

		expected := `
package main

import (
	"errors"
)

func main() {
	errors.New("example")
}
`
		actual := testutils.RenderIndependentStatementToString(t, result)

		assert.Equal(t, expected, actual)
	})
}
