package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mockUserDataServerDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		proj := testprojects.BuildTodoApp()
		x := mockUserDataServerDotGo(proj)

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
	v1 "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/models/v1"
	"net/http"
)

var _ v1.UserDataServer = (*UserDataServer)(nil)

// UserDataServer is a mocked models.UserDataServer for testing
type UserDataServer struct {
	mock.Mock
}

// UserLoginInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) UserLoginInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// UserInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) UserInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// PasswordUpdateInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) PasswordUpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// TOTPSecretVerificationInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) TOTPSecretVerificationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// TOTPSecretRefreshInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) TOTPSecretRefreshInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

// ListHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// CreateHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ReadHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// TOTPSecretVerificationHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) TOTPSecretVerificationHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// NewTOTPSecretHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) NewTOTPSecretHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// UpdatePasswordHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) UpdatePasswordHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

// ArchiveHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserDataServer(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserDataServer()

		expected := `
package example

import (
	mock "github.com/stretchr/testify/mock"
)

// UserDataServer is a mocked models.UserDataServer for testing
type UserDataServer struct {
	mock.Mock
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserLoginInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserLoginInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// UserLoginInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) UserLoginInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// UserInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) UserInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserPasswordUpdateInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserPasswordUpdateInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// PasswordUpdateInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) PasswordUpdateInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserTOTPSecretVerificationInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserTOTPSecretVerificationInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// TOTPSecretVerificationInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) TOTPSecretVerificationInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserTOTPSecretRefreshInputMiddleware(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserTOTPSecretRefreshInputMiddleware()

		expected := `
package example

import (
	"net/http"
)

// TOTPSecretRefreshInputMiddleware is a mock method to satisfy our interface requirements.
func (m *UserDataServer) TOTPSecretRefreshInputMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserListHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserListHandler()

		expected := `
package example

import (
	"net/http"
)

// ListHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) ListHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserCreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserCreateHandler()

		expected := `
package example

import (
	"net/http"
)

// CreateHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) CreateHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserReadHandler()

		expected := `
package example

import (
	"net/http"
)

// ReadHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) ReadHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserTOTPSecretVerificationHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserTOTPSecretVerificationHandler()

		expected := `
package example

import (
	"net/http"
)

// TOTPSecretVerificationHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) TOTPSecretVerificationHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserNewTOTPSecretHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserNewTOTPSecretHandler()

		expected := `
package example

import (
	"net/http"
)

// NewTOTPSecretHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) NewTOTPSecretHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserUpdatePasswordHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserUpdatePasswordHandler()

		expected := `
package example

import (
	"net/http"
)

// UpdatePasswordHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) UpdatePasswordHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildUserArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		t.Parallel()

		x := buildUserArchiveHandler()

		expected := `
package example

import (
	"net/http"
)

// ArchiveHandler is a mock method to satisfy our interface requirements.
func (m *UserDataServer) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
