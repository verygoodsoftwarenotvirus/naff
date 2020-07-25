package client

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_usersDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp
		out := usersDotGo(proj)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildBuildGetUserRequest(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildGetUser(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildGetUsersRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildBuildGetUsersRequest(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildGetUsers(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildGetUsers(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildCreateUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildBuildCreateUserRequest(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildCreateUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildCreateUser(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildArchiveUserRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildBuildArchiveUserRequest(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildArchiveUser(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildLoginRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildBuildLoginRequest(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildLogin(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildLogin(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildBuildVerifyTOTPSecretRequest(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildBuildVerifyTOTPSecretRequest(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}

func Test_buildVerifyTOTPSecret(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.TodoApp

		out := jen.NewFile("main")
		out.Add(buildVerifyTOTPSecret(proj)...)

		var b bytes.Buffer
		require.NoError(t, out.Render(&b))

		expected := `

`
		actual := "\n" + b.String()

		assert.Equal(t, actual, expected, "expected and actual output do not match")
	})
}
