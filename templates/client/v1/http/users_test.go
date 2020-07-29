package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_usersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := usersDotGo(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildGetUserRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildGetUser(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUsersRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildGetUsersRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildGetUsers(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildCreateUserRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildCreateUser(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildArchiveUserRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildArchiveUser(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildLoginRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildLoginRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildLogin(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildVerifyTOTPSecretRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildBuildVerifyTOTPSecretRequest(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildVerifyTOTPSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		x := buildVerifyTOTPSecret(proj)

		expected := ``
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
