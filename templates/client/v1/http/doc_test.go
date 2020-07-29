package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_docDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := docDotGo()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
