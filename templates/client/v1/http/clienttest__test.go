package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mainTestDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := mainTestDotGo(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildClientTestConstants(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildClientTestConstants()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildClientTestTypes(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildClientTestTypes()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildClientTestValuer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildClientTestValuer()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildMustParseURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildMustParseURL()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildTestClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildBuildTestClient()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildTestClientWithInvalidURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildBuildTestClientWithInvalidURL()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_AuthenticatedClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_AuthenticatedClient()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_PlainClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_PlainClient()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_TokenSource(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_TokenSource()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestNewClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestNewClient()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestNewSimpleClient(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestNewSimpleClient()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_CloseRequestBody(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_CloseRequestBody()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestBuildURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestBuildURL()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestBuildVersionlessURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestBuildVersionlessURL()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildWebsocketURL(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_BuildWebsocketURL()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_BuildHealthCheckRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_BuildHealthCheckRequest()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_IsUp(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_IsUp()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_buildDataRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_buildDataRequest()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_checkExistence(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_checkExistence()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_retrieve(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_retrieve()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_executeRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_executeRequest()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_executeRawRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_executeRawRequest()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildTestV1Client_executeUnauthenticatedDataRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildTestV1Client_executeUnauthenticatedDataRequest()

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
