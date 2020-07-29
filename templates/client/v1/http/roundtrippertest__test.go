package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_roundtripperTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := roundtripperTestDotGo(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestBuildDefaultTransport(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestBuildDefaultTransport()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestDefaultRoundTripperRoundTrip(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestDefaultRoundTripperRoundTrip()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestNewDefaultRoundTripper(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestNewDefaultRoundTripper()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
