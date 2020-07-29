package client

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_oauth2ClientsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := oauth2ClientsDotGo(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildGetOAuth2ClientRequest(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildGetOAuth2Client(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetOAuth2ClientsRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildGetOAuth2ClientsRequest(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildGetOAuth2Clients(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildCreateOAuth2ClientRequest(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildCreateOAuth2Client(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveOAuth2ClientRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildArchiveOAuth2ClientRequest(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildArchiveOAuth2Client(proj)

		expected := ``
		actual := testutils.RenderFunctionToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
